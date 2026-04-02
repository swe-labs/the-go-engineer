// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

package main

import (
	"os"
	"path/filepath"
	"testing"
)

// ============================================================================
// Tests for: Config File Parser (Encoding Exercise)
// ============================================================================
//
// These tests verify loadConfig and validate handle valid configs,
// missing fields, invalid JSON, and missing files correctly.
//
// RUN: go test -v ./11-encoding/6-config-parser
// ============================================================================

func TestLoadConfigValid(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.json")
	validJSON := `{
		"app_name": "TestApp",
		"port": 3000,
		"database_url": "postgres://localhost/test",
		"debug": false,
		"max_workers": 2,
		"allowed_hosts": ["localhost"]
	}`
	os.WriteFile(configPath, []byte(validJSON), 0644)

	config, err := loadConfig(configPath)
	if err != nil {
		t.Fatalf("loadConfig failed: %v", err)
	}

	if config.AppName != "TestApp" {
		t.Errorf("AppName = %q, want %q", config.AppName, "TestApp")
	}
	if config.Port != 3000 {
		t.Errorf("Port = %d, want 3000", config.Port)
	}
	if config.MaxWorkers != 2 {
		t.Errorf("MaxWorkers = %d, want 2", config.MaxWorkers)
	}
}

func TestLoadConfigMissingField(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.json")
	// Missing "port" field
	incompleteJSON := `{"app_name": "TestApp", "database_url": "postgres://localhost", "max_workers": 1}`
	os.WriteFile(configPath, []byte(incompleteJSON), 0644)

	_, err := loadConfig(configPath)
	if err == nil {
		t.Error("loadConfig should return error for missing port, got nil")
	}
}

func TestLoadConfigInvalidJSON(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.json")
	os.WriteFile(configPath, []byte("{invalid json}"), 0644)

	_, err := loadConfig(configPath)
	if err == nil {
		t.Error("loadConfig should return error for invalid JSON, got nil")
	}
}

func TestLoadConfigFileMissing(t *testing.T) {
	_, err := loadConfig("/nonexistent/path/config.json")
	if err == nil {
		t.Error("loadConfig should return error for missing file, got nil")
	}
}

func TestValidateAllFieldsPresent(t *testing.T) {
	config := &AppConfig{
		AppName:     "App",
		Port:        8080,
		DatabaseURL: "postgres://localhost",
		MaxWorkers:  4,
	}

	err := config.validate()
	if err != nil {
		t.Errorf("validate should pass with all fields, got: %v", err)
	}
}

func TestValidateMissingAppName(t *testing.T) {
	config := &AppConfig{Port: 8080, DatabaseURL: "postgres://localhost", MaxWorkers: 4}

	err := config.validate()
	if err == nil {
		t.Error("validate should fail with empty app_name")
	}
}
