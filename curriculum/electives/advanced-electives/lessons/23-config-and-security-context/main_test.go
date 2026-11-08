package main

import (
	"os"
	"path/filepath"
	"testing"
)

func mapEnv(m map[string]string) func(string) string {
	return func(k string) string { return m[k] }
}

func writeSecret(t *testing.T, body string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "db-url")
	if err := os.WriteFile(path, []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestLoadConfigValid(t *testing.T) {
	path := writeSecret(t, "postgres://db\n")
	env := mapEnv(map[string]string{
		"APP_ADDR":              ":9090",
		"APP_LOG_LEVEL":         "debug",
		"APP_MAX_CONNS":         "25",
		"APP_DATABASE_URL_FILE": path,
	})

	cfg, err := LoadConfig(env, os.ReadFile)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Addr != ":9090" || cfg.LogLevel != "debug" || cfg.MaxConns != 25 {
		t.Fatalf("unexpected config: %+v", cfg)
	}
	if cfg.DatabaseURL != "postgres://db" {
		t.Fatalf("DatabaseURL = %q, want trimmed value", cfg.DatabaseURL)
	}
}

func TestLoadConfigDefaults(t *testing.T) {
	path := writeSecret(t, "postgres://db")
	env := mapEnv(map[string]string{
		"APP_MAX_CONNS":         "5",
		"APP_DATABASE_URL_FILE": path,
	})

	cfg, err := LoadConfig(env, os.ReadFile)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Addr != ":8080" || cfg.LogLevel != "info" {
		t.Fatalf("defaults not applied: %+v", cfg)
	}
}

func TestLoadConfigMissingMaxConns(t *testing.T) {
	if _, err := LoadConfig(mapEnv(nil), os.ReadFile); err == nil {
		t.Fatal("expected error for missing APP_MAX_CONNS")
	}
}

func TestLoadConfigInvalidMaxConns(t *testing.T) {
	env := mapEnv(map[string]string{"APP_MAX_CONNS": "zero"})
	if _, err := LoadConfig(env, os.ReadFile); err == nil {
		t.Fatal("expected error for invalid APP_MAX_CONNS")
	}
}

func TestLoadConfigInvalidLogLevel(t *testing.T) {
	env := mapEnv(map[string]string{
		"APP_LOG_LEVEL": "trace",
		"APP_MAX_CONNS": "5",
	})
	if _, err := LoadConfig(env, os.ReadFile); err == nil {
		t.Fatal("expected error for invalid APP_LOG_LEVEL")
	}
}

func TestLoadConfigSecretRotates(t *testing.T) {
	path := writeSecret(t, "postgres://a")
	env := mapEnv(map[string]string{"APP_MAX_CONNS": "5", "APP_DATABASE_URL_FILE": path})

	first, err := LoadConfig(env, os.ReadFile)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte("postgres://b"), 0o600); err != nil {
		t.Fatal(err)
	}
	second, err := LoadConfig(env, os.ReadFile)
	if err != nil {
		t.Fatal(err)
	}
	if first.DatabaseURL == second.DatabaseURL {
		t.Fatalf("expected rotated value, both %q", first.DatabaseURL)
	}
}
