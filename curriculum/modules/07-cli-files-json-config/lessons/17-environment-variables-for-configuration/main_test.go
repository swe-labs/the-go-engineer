package main

import (
	"os"
	"testing"
)

func TestConfigFromEnvDefaults(t *testing.T) {
	os.Unsetenv("PORT")
	os.Unsetenv("DATABASE_URL")
	os.Unsetenv("VERBOSE")

	cfg := ConfigFromEnv()
	if cfg.Port != 8080 {
		t.Errorf("expected default port 8080, got %d", cfg.Port)
	}
	if cfg.DBURL != "postgres://localhost:5432/app?sslmode=disable" {
		t.Errorf("unexpected default DBURL: %s", cfg.DBURL)
	}
	if cfg.Verbose {
		t.Error("expected verbose false by default")
	}
}

func TestConfigFromEnvOverride(t *testing.T) {
	os.Setenv("PORT", "9999")
	os.Setenv("DATABASE_URL", "postgres://test:5432/testdb")
	os.Setenv("VERBOSE", "true")
	defer os.Unsetenv("PORT")
	defer os.Unsetenv("DATABASE_URL")
	defer os.Unsetenv("VERBOSE")

	cfg := ConfigFromEnv()
	if cfg.Port != 9999 {
		t.Errorf("expected port 9999, got %d", cfg.Port)
	}
	if cfg.DBURL != "postgres://test:5432/testdb" {
		t.Errorf("unexpected DBURL: %s", cfg.DBURL)
	}
	if !cfg.Verbose {
		t.Error("expected verbose true")
	}
}

func TestConfigFromEnvEmptyVerbose(t *testing.T) {
	os.Setenv("VERBOSE", "")
	defer os.Unsetenv("VERBOSE")

	cfg := ConfigFromEnv()
	if cfg.Verbose {
		t.Error("expected verbose false when env var is empty")
	}
}

func TestConfigFromEnvBadPort(t *testing.T) {
	os.Setenv("PORT", "not-a-number")
	defer os.Unsetenv("PORT")

	cfg := ConfigFromEnv()
	if cfg.Port != 8080 {
		t.Errorf("expected port to stay at default 8080, got %d", cfg.Port)
	}
}

func TestCompiles(t *testing.T) {
}
