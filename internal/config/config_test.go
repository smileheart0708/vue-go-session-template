package config_test

import (
	"testing"

	"main/internal/config"
)

func TestLoadGeneratesAuthKeyWhenMissing(t *testing.T) {
	t.Setenv("AUTH_KEY", "")

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if cfg.AuthKey == "" {
		t.Fatal("expected auto generated auth key")
	}
	if !cfg.IsAutoAuthKey {
		t.Fatal("expected IsAutoAuthKey=true when AUTH_KEY is empty")
	}
}

func TestLoadUsesProvidedAuthKey(t *testing.T) {
	t.Setenv("AUTH_KEY", "fixed-auth-key")

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if cfg.AuthKey != "fixed-auth-key" {
		t.Fatalf("unexpected AuthKey: %q", cfg.AuthKey)
	}
	if cfg.IsAutoAuthKey {
		t.Fatal("expected IsAutoAuthKey=false when AUTH_KEY is provided")
	}
}
