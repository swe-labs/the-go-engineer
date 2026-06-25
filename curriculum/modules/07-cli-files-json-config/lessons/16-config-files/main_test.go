package main

import (
	"os"
	"testing"
)

func TestLoadAppConfigDefaults(t *testing.T) {
	cfg, err := LoadAppConfig("_nonexistent_file_xyz.json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.ServerPort != 8080 {
		t.Errorf("expected default port 8080, got %d", cfg.ServerPort)
	}
	if cfg.LogLevel != "info" {
		t.Errorf("expected default log_level info, got %s", cfg.LogLevel)
	}
}

func TestLoadAppConfigMerge(t *testing.T) {
	os.WriteFile("_test_merge_a.json", []byte(`{"server_port":9000,"database":{"host":"db1"}}`), 0644)
	os.WriteFile("_test_merge_b.json", []byte(`{"debug":true,"database":{"port":6543}}`), 0644)
	defer os.Remove("_test_merge_a.json")
	defer os.Remove("_test_merge_b.json")

	cfg, err := LoadAppConfig("_test_merge_a.json", "_test_merge_b.json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.ServerPort != 9000 {
		t.Errorf("expected 9000, got %d", cfg.ServerPort)
	}
	if cfg.Database.Host != "db1" {
		t.Errorf("expected db1, got %s", cfg.Database.Host)
	}
	if cfg.Database.Port != 6543 {
		t.Errorf("expected 6543, got %d", cfg.Database.Port)
	}
	if !cfg.Debug {
		t.Error("expected debug true")
	}
}

func TestLoadAppConfigSkipMissing(t *testing.T) {
	os.WriteFile("_test_only.json", []byte(`{"server_port":3000}`), 0644)
	defer os.Remove("_test_only.json")

	cfg, err := LoadAppConfig("_test_missing.json", "_test_only.json", "_test_also_missing.json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.ServerPort != 3000 {
		t.Errorf("expected 3000, got %d", cfg.ServerPort)
	}
}

func TestCompiles(t *testing.T) {
}
