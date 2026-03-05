package session

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const authKeyFingerprintFileName = ".session_auth_key"

type CleanupStats struct {
	Scanned int
	Deleted int
	Kept    int
	Failed  int
}

func (s *CleanupStats) add(other CleanupStats) {
	s.Scanned += other.Scanned
	s.Deleted += other.Deleted
	s.Kept += other.Kept
	s.Failed += other.Failed
}

// Bootstrap prepares session storage, handles AUTH_KEY changes and runs one cleanup pass.
func Bootstrap(dataDir, authKey string, now time.Time) (CleanupStats, error) {
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return CleanupStats{}, fmt.Errorf("create data dir: %w", err)
	}

	sessionDir := filepath.Join(dataDir, SessionDirectoryName)
	if err := os.MkdirAll(sessionDir, 0755); err != nil {
		return CleanupStats{}, fmt.Errorf("create session dir: %w", err)
	}

	markerPath := filepath.Join(dataDir, authKeyFingerprintFileName)
	currentFingerprint := buildAuthKeyFingerprint(authKey)

	totalStats := CleanupStats{}
	authKeyChanged := false
	keyPurgeDeleted := 0
	existingFingerprint, err := readFingerprint(markerPath)
	if err != nil {
		return CleanupStats{}, fmt.Errorf("read auth key fingerprint: %w", err)
	}
	if existingFingerprint != "" && existingFingerprint != currentFingerprint {
		authKeyChanged = true
		purgeStats, purgeErr := purgeSessionFiles(sessionDir)
		if purgeErr != nil {
			return purgeStats, fmt.Errorf("purge sessions on auth key change: %w", purgeErr)
		}
		totalStats.add(purgeStats)
		keyPurgeDeleted = purgeStats.Deleted
	}

	if existingFingerprint == "" || authKeyChanged {
		if err := os.WriteFile(markerPath, []byte(currentFingerprint), 0600); err != nil {
			return CleanupStats{}, fmt.Errorf("write auth key fingerprint: %w", err)
		}
	}

	expiredStats, err := CleanupExpired(sessionDir, now)
	if err != nil {
		return expiredStats, fmt.Errorf("cleanup expired sessions: %w", err)
	}
	totalStats.add(expiredStats)

	slog.Info(
		"session bootstrap completed",
		"auth_key_changed",
		authKeyChanged,
		"purged_count",
		keyPurgeDeleted,
		"scanned",
		totalStats.Scanned,
		"deleted",
		totalStats.Deleted,
		"kept",
		totalStats.Kept,
		"errors",
		totalStats.Failed,
	)

	return totalStats, nil
}

// RunJanitor periodically removes expired session files.
func RunJanitor(ctx context.Context, dataDir string, nowFn func() time.Time) {
	if nowFn == nil {
		nowFn = time.Now
	}

	ticker := time.NewTicker(CleanupInterval)
	defer ticker.Stop()

	sessionDir := filepath.Join(dataDir, SessionDirectoryName)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			stats, err := CleanupExpired(sessionDir, nowFn())
			if err != nil {
				slog.Warn("session janitor cleanup failed", "error", err)
				continue
			}
			if stats.Failed > 0 {
				slog.Warn(
					"session janitor cleanup completed with errors",
					"scanned",
					stats.Scanned,
					"deleted",
					stats.Deleted,
					"kept",
					stats.Kept,
					"errors",
					stats.Failed,
				)
				continue
			}
			if stats.Deleted > 0 {
				slog.Info(
					"session janitor cleanup completed",
					"scanned",
					stats.Scanned,
					"deleted",
					stats.Deleted,
					"kept",
					stats.Kept,
					"errors",
					stats.Failed,
				)
			}
		}
	}
}

// CleanupExpired removes session files older than SessionTTL plus CleanupGrace.
func CleanupExpired(sessionDir string, now time.Time) (CleanupStats, error) {
	entries, err := os.ReadDir(sessionDir)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return CleanupStats{}, nil
		}
		return CleanupStats{}, fmt.Errorf("read session dir: %w", err)
	}

	threshold := SessionTTL + CleanupGrace
	stats := CleanupStats{}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasPrefix(entry.Name(), SessionFilePrefix) {
			continue
		}

		stats.Scanned++

		info, infoErr := entry.Info()
		if infoErr != nil {
			stats.Failed++
			slog.Warn(
				"failed to stat session file",
				"file",
				entry.Name(),
				"error",
				infoErr,
			)
			continue
		}

		if now.Sub(info.ModTime()) > threshold {
			filePath := filepath.Join(sessionDir, entry.Name())
			if removeErr := os.Remove(filePath); removeErr != nil && !errors.Is(removeErr, os.ErrNotExist) {
				stats.Failed++
				slog.Warn(
					"failed to remove expired session file",
					"file",
					entry.Name(),
					"error",
					removeErr,
				)
				continue
			}
			stats.Deleted++
			continue
		}

		stats.Kept++
	}

	return stats, nil
}

func purgeSessionFiles(sessionDir string) (CleanupStats, error) {
	entries, err := os.ReadDir(sessionDir)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return CleanupStats{}, nil
		}
		return CleanupStats{}, fmt.Errorf("read session dir: %w", err)
	}

	stats := CleanupStats{}
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasPrefix(entry.Name(), SessionFilePrefix) {
			continue
		}

		stats.Scanned++
		filePath := filepath.Join(sessionDir, entry.Name())
		if removeErr := os.Remove(filePath); removeErr != nil && !errors.Is(removeErr, os.ErrNotExist) {
			stats.Failed++
			slog.Warn(
				"failed to purge session file",
				"file",
				entry.Name(),
				"error",
				removeErr,
			)
			continue
		}
		stats.Deleted++
	}

	return stats, nil
}

func readFingerprint(markerPath string) (string, error) {
	content, err := os.ReadFile(markerPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return "", nil
		}
		return "", err
	}
	return strings.TrimSpace(string(content)), nil
}

func buildAuthKeyFingerprint(authKey string) string {
	hash := sha256.Sum256([]byte(authKey))
	return hex.EncodeToString(hash[:])
}
