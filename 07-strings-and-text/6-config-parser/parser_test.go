// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

package main

import "testing"

// ============================================================================
// Tests for: Config Parser (Strings Exercise)
// ============================================================================
//
// These tests verify the parseConfig function correctly handles
// key-value pairs, quoted strings, comments, and edge cases.
//
// RUN: go test -v ./07-strings-and-text/6-config-parser
// ============================================================================

func TestParseConfigBasic(t *testing.T) {
	input := `
host = localhost
port = 8080
debug = true
`
	config, err := parseConfig(input)
	if err != nil {
		t.Fatalf("parseConfig failed: %v", err)
	}

	tests := []struct {
		key  string
		want string
	}{
		{"host", "localhost"},
		{"port", "8080"},
		{"debug", "true"},
	}

	for _, tt := range tests {
		got, exists := config[tt.key]
		if !exists {
			t.Errorf("key %q not found in config", tt.key)
			continue
		}
		if got != tt.want {
			t.Errorf("config[%q] = %q, want %q", tt.key, got, tt.want)
		}
	}
}

func TestParseConfigSkipsComments(t *testing.T) {
	input := `
# This is a comment
name = Bob
# Another comment
`
	config, err := parseConfig(input)
	if err != nil {
		t.Fatalf("parseConfig failed: %v", err)
	}

	if len(config) != 1 {
		t.Errorf("expected 1 key, got %d", len(config))
	}
	if config["name"] != "Bob" {
		t.Errorf("config[name] = %q, want %q", config["name"], "Bob")
	}
}

func TestParseConfigQuotedValues(t *testing.T) {
	input := `
greeting = "hello world"
path = '/usr/local/bin'
`
	config, err := parseConfig(input)
	if err != nil {
		t.Fatalf("parseConfig failed: %v", err)
	}

	if config["greeting"] != "hello world" {
		t.Errorf("config[greeting] = %q, want %q", config["greeting"], "hello world")
	}
	if config["path"] != "/usr/local/bin" {
		t.Errorf("config[path] = %q, want %q", config["path"], "/usr/local/bin")
	}
}

func TestParseConfigEmptyInput(t *testing.T) {
	config, err := parseConfig("")
	if err != nil {
		t.Fatalf("parseConfig failed: %v", err)
	}

	if len(config) != 0 {
		t.Errorf("expected 0 keys for empty input, got %d", len(config))
	}
}
