package main

import (
	"strings"
	"testing"
)

func TestValidateConfigValid(t *testing.T) {
	cfg := ServerConfig{
		Port:     8080,
		Host:     "0.0.0.0",
		DBURL:    "postgres://localhost:5432/db",
		LogLevel: "info",
		Workers:  4,
	}
	if err := ValidateConfig(cfg); err != nil {
		t.Fatalf("unexpected error for valid config: %v", err)
	}
}

func TestValidateConfigPortRange(t *testing.T) {
	cfg := ServerConfig{
		Port:     99999,
		Host:     "localhost",
		DBURL:    "postgres://localhost:5432/db",
		LogLevel: "info",
		Workers:  1,
	}
	err := ValidateConfig(cfg)
	if err == nil {
		t.Fatal("expected error for invalid port")
	}
	if !strings.Contains(err.Error(), "port") {
		t.Errorf("error should mention port: %v", err)
	}
}

func TestValidateConfigMissingHost(t *testing.T) {
	cfg := ServerConfig{
		Port:     8080,
		Host:     "",
		DBURL:    "postgres://localhost:5432/db",
		LogLevel: "info",
		Workers:  1,
	}
	err := ValidateConfig(cfg)
	if err == nil {
		t.Fatal("expected error for missing host")
	}
	if !strings.Contains(err.Error(), "host") {
		t.Errorf("error should mention host: %v", err)
	}
}

func TestValidateConfigBadURL(t *testing.T) {
	cfg := ServerConfig{
		Port:     8080,
		Host:     "localhost",
		DBURL:    ":::invalid",
		LogLevel: "info",
		Workers:  1,
	}
	err := ValidateConfig(cfg)
	if err == nil {
		t.Fatal("expected error for bad DB URL")
	}
}

func TestValidateConfigBadLogLevel(t *testing.T) {
	cfg := ServerConfig{
		Port:     8080,
		Host:     "localhost",
		DBURL:    "postgres://localhost:5432/db",
		LogLevel: "critical",
		Workers:  1,
	}
	err := ValidateConfig(cfg)
	if err == nil {
		t.Fatal("expected error for bad log level")
	}
}

func TestValidateConfigWorkers(t *testing.T) {
	cfg := ServerConfig{
		Port:     8080,
		Host:     "localhost",
		DBURL:    "postgres://localhost:5432/db",
		LogLevel: "info",
		Workers:  0,
	}
	err := ValidateConfig(cfg)
	if err == nil {
		t.Fatal("expected error for zero workers")
	}
}

func TestValidateConfigMultipleErrors(t *testing.T) {
	cfg := ServerConfig{
		Port:     0,
		Host:     "",
		DBURL:    "",
		LogLevel: "",
		Workers:  0,
	}
	err := ValidateConfig(cfg)
	if err == nil {
		t.Fatal("expected multiple errors")
	}
	msg := err.Error()
	// Should contain at least 3 distinct error messages
	count := 0
	for _, keyword := range []string{"port", "host", "database_url", "log_level", "workers"} {
		if strings.Contains(msg, keyword) {
			count++
		}
	}
	if count < 3 {
		t.Errorf("expected at least 3 field errors, found %d in: %s", count, msg)
	}
}

func TestCompiles(t *testing.T) {
}
