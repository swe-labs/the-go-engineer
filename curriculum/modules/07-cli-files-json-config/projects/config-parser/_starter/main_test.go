package main

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestConfigParser_Defaults(t *testing.T) {
	parser := NewConfigParserWithLookup(func(s string) (string, bool) {
		return "", false
	})

	cfg, err := parser.Load("")
	if err != nil {
		t.Fatalf("failed to load default config: %v", err)
	}

	if cfg.Environment != "development" {
		t.Fatalf("expected environment 'development', got %q", cfg.Environment)
	}
	if cfg.Server.Port != 8080 {
		t.Fatalf("expected port 8080, got %d", cfg.Server.Port)
	}
	if cfg.Database.Name != "app_dev" {
		t.Fatalf("expected db name 'app_dev', got %q", cfg.Database.Name)
	}
}

func TestConfigParser_JSONFileMerge(t *testing.T) {
	tempDir := t.TempDir()
	jsonPath := filepath.Join(tempDir, "config.json")

	jsonContent := `{
		"environment": "staging",
		"server": {
			"port": 3000
		},
		"database": {
			"host": "staging-db.internal",
			"name": "staging_db"
		}
	}`
	_ = os.WriteFile(jsonPath, []byte(jsonContent), 0644)

	parser := NewConfigParserWithLookup(func(s string) (string, bool) {
		return "", false
	})

	cfg, err := parser.Load(jsonPath)
	if err != nil {
		t.Fatalf("failed to load json config: %v", err)
	}

	if cfg.Environment != "staging" {
		t.Fatalf("expected environment 'staging', got %q", cfg.Environment)
	}
	if cfg.Server.Port != 3000 {
		t.Fatalf("expected server port 3000, got %d", cfg.Server.Port)
	}
	// Unspecified in file should retain default
	if cfg.Server.Host != "0.0.0.0" {
		t.Fatalf("expected server host '0.0.0.0', got %q", cfg.Server.Host)
	}
	if cfg.Database.Host != "staging-db.internal" {
		t.Fatalf("expected db host 'staging-db.internal', got %q", cfg.Database.Host)
	}
	if cfg.Database.Name != "staging_db" {
		t.Fatalf("expected db name 'staging_db', got %q", cfg.Database.Name)
	}
}

func TestConfigParser_EnvOverridePrecedence(t *testing.T) {
	tempDir := t.TempDir()
	jsonPath := filepath.Join(tempDir, "config.json")

	jsonContent := `{
		"server": { "port": 4000 }
	}`
	_ = os.WriteFile(jsonPath, []byte(jsonContent), 0644)

	// Mock env variables: APP_SERVER_PORT overrides file port (4000 -> 9090)
	envMap := map[string]string{
		"APP_SERVER_PORT": "9090",
		"APP_DB_NAME":     "production_cluster",
	}

	parser := NewConfigParserWithLookup(func(key string) (string, bool) {
		val, ok := envMap[key]
		return val, ok
	})

	cfg, err := parser.Load(jsonPath)
	if err != nil {
		t.Fatalf("failed to load config with overrides: %v", err)
	}

	if cfg.Server.Port != 9090 {
		t.Fatalf("expected env override port 9090, got %d", cfg.Server.Port)
	}
	if cfg.Database.Name != "production_cluster" {
		t.Fatalf("expected env override db name 'production_cluster', got %q", cfg.Database.Name)
	}
}

func TestConfigParser_ValidationErrors(t *testing.T) {
	parser := NewConfigParserWithLookup(func(s string) (string, bool) {
		return "", false
	})

	badConfig := DefaultConfig()
	badConfig.Server.Port = 99999 // Exceeds 65535
	badConfig.Database.Name = ""  // Empty db name

	err := parser.Validate(badConfig)
	if !errors.Is(err, ErrValidation) {
		t.Fatalf("expected ErrValidation, got %v", err)
	}
}

func TestConfigParser_UnsupportedFormat(t *testing.T) {
	tempDir := t.TempDir()
	xmlPath := filepath.Join(tempDir, "config.xml")
	_ = os.WriteFile(xmlPath, []byte("<config></config>"), 0644)

	parser := NewConfigParser()
	_, err := parser.Load(xmlPath)
	if !errors.Is(err, ErrUnsupportedFormat) {
		t.Fatalf("expected ErrUnsupportedFormat, got %v", err)
	}
}
