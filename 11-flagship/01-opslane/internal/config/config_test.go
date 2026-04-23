package config

import (
	"log/slog"
	"testing"
	"time"
)

func TestLoadFromLookupDefaults(t *testing.T) {
	t.Parallel()

	cfg, err := LoadFromLookup(func(string) (string, bool) {
		return "", false
	})
	if err != nil {
		t.Fatalf("LoadFromLookup returned error: %v", err)
	}

	if cfg.App.Name != "opslane" {
		t.Fatalf("name = %q, want %q", cfg.App.Name, "opslane")
	}

	if cfg.App.Env != "development" {
		t.Fatalf("env = %q, want %q", cfg.App.Env, "development")
	}

	if cfg.App.LogLevel != slog.LevelInfo {
		t.Fatalf("log level = %v, want %v", cfg.App.LogLevel, slog.LevelInfo)
	}

	if cfg.HTTP.Address != ":8080" {
		t.Fatalf("address = %q, want %q", cfg.HTTP.Address, ":8080")
	}
}

func TestLoadFromLookupOverrides(t *testing.T) {
	t.Parallel()

	values := map[string]string{
		"OPSLANE_ENV":                      "staging",
		"OPSLANE_HTTP_ADDR":                ":9090",
		"OPSLANE_LOG_LEVEL":                "debug",
		"OPSLANE_HTTP_READ_HEADER_TIMEOUT": "7s",
		"OPSLANE_HTTP_READ_TIMEOUT":        "12s",
		"OPSLANE_HTTP_WRITE_TIMEOUT":       "18s",
		"OPSLANE_HTTP_IDLE_TIMEOUT":        "90s",
		"OPSLANE_HTTP_SHUTDOWN_TIMEOUT":    "25s",
	}

	cfg, err := LoadFromLookup(func(key string) (string, bool) {
		value, ok := values[key]
		return value, ok
	})
	if err != nil {
		t.Fatalf("LoadFromLookup returned error: %v", err)
	}

	if cfg.App.Env != "staging" {
		t.Fatalf("env = %q, want %q", cfg.App.Env, "staging")
	}

	if cfg.HTTP.Address != ":9090" {
		t.Fatalf("address = %q, want %q", cfg.HTTP.Address, ":9090")
	}

	if cfg.App.LogLevel != slog.LevelDebug {
		t.Fatalf("log level = %v, want %v", cfg.App.LogLevel, slog.LevelDebug)
	}

	if cfg.HTTP.ReadHeaderTimeout != 7*time.Second {
		t.Fatalf("read header timeout = %v, want %v", cfg.HTTP.ReadHeaderTimeout, 7*time.Second)
	}
}

func TestLoadFromLookupRejectsInvalidLogLevel(t *testing.T) {
	t.Parallel()

	_, err := LoadFromLookup(func(key string) (string, bool) {
		if key == "OPSLANE_LOG_LEVEL" {
			return "loud", true
		}
		return "", false
	})
	if err == nil {
		t.Fatal("expected error for invalid log level")
	}
}

func TestLoadFromLookupRejectsInvalidEnv(t *testing.T) {
	t.Parallel()

	_, err := LoadFromLookup(func(key string) (string, bool) {
		if key == "OPSLANE_ENV" {
			return "preview", true
		}
		return "", false
	})
	if err == nil {
		t.Fatal("expected error for invalid env")
	}
}

func TestLoadFromLookupRejectsInvalidDuration(t *testing.T) {
	t.Parallel()

	_, err := LoadFromLookup(func(key string) (string, bool) {
		if key == "OPSLANE_HTTP_READ_TIMEOUT" {
			return "soon", true
		}
		return "", false
	})
	if err == nil {
		t.Fatal("expected error for invalid duration")
	}
}
