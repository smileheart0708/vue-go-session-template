package session

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestBootstrapKeepsSessionsWhenAuthKeyUnchanged(t *testing.T) {
	dataDir := t.TempDir()
	now := time.Unix(1_700_000_000, 0)

	if _, err := Bootstrap(dataDir, "auth-key", now); err != nil {
		t.Fatalf("bootstrap should succeed: %v", err)
	}

	sessionDir := filepath.Join(dataDir, SessionDirectoryName)
	sessionFile := filepath.Join(sessionDir, SessionFilePrefix+"ACTIVE")
	if err := os.WriteFile(sessionFile, []byte("value"), 0600); err != nil {
		t.Fatalf("write session file: %v", err)
	}
	if err := os.Chtimes(sessionFile, now, now); err != nil {
		t.Fatalf("set session file mtime: %v", err)
	}

	stats, err := Bootstrap(dataDir, "auth-key", now)
	if err != nil {
		t.Fatalf("bootstrap should succeed with unchanged auth key: %v", err)
	}
	if stats.Deleted != 0 {
		t.Fatalf("expected no deleted sessions, got %d", stats.Deleted)
	}

	if _, err := os.Stat(sessionFile); err != nil {
		t.Fatalf("session file should still exist: %v", err)
	}
}

func TestBootstrapPurgesSessionsWhenAuthKeyChanges(t *testing.T) {
	dataDir := t.TempDir()
	now := time.Unix(1_700_000_000, 0)

	if _, err := Bootstrap(dataDir, "auth-key-old", now); err != nil {
		t.Fatalf("initial bootstrap should succeed: %v", err)
	}

	sessionDir := filepath.Join(dataDir, SessionDirectoryName)
	sessionFile := filepath.Join(sessionDir, SessionFilePrefix+"STALE")
	if err := os.WriteFile(sessionFile, []byte("value"), 0600); err != nil {
		t.Fatalf("write session file: %v", err)
	}
	nonSessionFile := filepath.Join(sessionDir, "keep.txt")
	if err := os.WriteFile(nonSessionFile, []byte("keep"), 0600); err != nil {
		t.Fatalf("write non-session file: %v", err)
	}

	if _, err := Bootstrap(dataDir, "auth-key-new", now); err != nil {
		t.Fatalf("bootstrap should succeed after auth key change: %v", err)
	}

	if _, err := os.Stat(sessionFile); !os.IsNotExist(err) {
		t.Fatalf("session file should be removed, got err=%v", err)
	}
	if _, err := os.Stat(nonSessionFile); err != nil {
		t.Fatalf("non-session file should be kept: %v", err)
	}
}

func TestCleanupExpiredUsesGraceWindowBoundary(t *testing.T) {
	sessionDir := t.TempDir()
	now := time.Unix(1_700_000_000, 0)

	oldFile := filepath.Join(sessionDir, SessionFilePrefix+"OLD")
	if err := os.WriteFile(oldFile, []byte("old"), 0600); err != nil {
		t.Fatalf("write old session file: %v", err)
	}
	oldMTime := now.Add(-(SessionTTL + CleanupGrace + time.Second))
	if err := os.Chtimes(oldFile, oldMTime, oldMTime); err != nil {
		t.Fatalf("set old file mtime: %v", err)
	}

	boundaryFile := filepath.Join(sessionDir, SessionFilePrefix+"BOUNDARY")
	if err := os.WriteFile(boundaryFile, []byte("boundary"), 0600); err != nil {
		t.Fatalf("write boundary session file: %v", err)
	}
	boundaryMTime := now.Add(-(SessionTTL + CleanupGrace))
	if err := os.Chtimes(boundaryFile, boundaryMTime, boundaryMTime); err != nil {
		t.Fatalf("set boundary file mtime: %v", err)
	}

	freshFile := filepath.Join(sessionDir, SessionFilePrefix+"FRESH")
	if err := os.WriteFile(freshFile, []byte("fresh"), 0600); err != nil {
		t.Fatalf("write fresh session file: %v", err)
	}
	freshMTime := now.Add(-time.Hour)
	if err := os.Chtimes(freshFile, freshMTime, freshMTime); err != nil {
		t.Fatalf("set fresh file mtime: %v", err)
	}

	nonSessionFile := filepath.Join(sessionDir, "plain.txt")
	if err := os.WriteFile(nonSessionFile, []byte("plain"), 0600); err != nil {
		t.Fatalf("write non-session file: %v", err)
	}
	if err := os.Chtimes(nonSessionFile, oldMTime, oldMTime); err != nil {
		t.Fatalf("set non-session file mtime: %v", err)
	}

	stats, err := CleanupExpired(sessionDir, now)
	if err != nil {
		t.Fatalf("cleanup should succeed: %v", err)
	}
	if stats.Scanned != 3 {
		t.Fatalf("expected scanned=3, got %d", stats.Scanned)
	}
	if stats.Deleted != 1 {
		t.Fatalf("expected deleted=1, got %d", stats.Deleted)
	}
	if stats.Kept != 2 {
		t.Fatalf("expected kept=2, got %d", stats.Kept)
	}
	if stats.Failed != 0 {
		t.Fatalf("expected failed=0, got %d", stats.Failed)
	}

	if _, err := os.Stat(oldFile); !os.IsNotExist(err) {
		t.Fatalf("old file should be removed, got err=%v", err)
	}
	if _, err := os.Stat(boundaryFile); err != nil {
		t.Fatalf("boundary file should be kept: %v", err)
	}
	if _, err := os.Stat(freshFile); err != nil {
		t.Fatalf("fresh file should be kept: %v", err)
	}
	if _, err := os.Stat(nonSessionFile); err != nil {
		t.Fatalf("non-session file should be kept: %v", err)
	}
}
