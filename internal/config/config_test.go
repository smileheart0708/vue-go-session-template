package config_test

import (
	"testing"

	"main/internal/config"
)

func TestLoadRejectsInvalidSessionEncKeyLength(t *testing.T) {
	t.Setenv("SESSION_ENC_KEY", "too-short")

	cfg, err := config.Load()
	if err == nil {
		t.Fatalf("expected load error, got cfg=%+v", cfg)
	}
}

func TestLoadAcceptsValidSessionEncKeyLength(t *testing.T) {
	t.Setenv("SESSION_ENC_KEY", "1234567890abcdef")

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if cfg.SessionEncKey != "1234567890abcdef" {
		t.Fatalf("unexpected SessionEncKey: %q", cfg.SessionEncKey)
	}
}
