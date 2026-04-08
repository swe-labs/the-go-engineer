// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// ============================================================================
// Section 09: Encoding — Config File Parser (Exercise)
// Level: Intermediate
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Reading a JSON file from disk and decoding it into a struct
//   - Validating required fields after parsing
//   - Combining filesystem I/O with JSON decoding
//   - Using json.NewDecoder for streaming file reads
//
// ENGINEERING DEPTH:
//   Using `json.NewDecoder(file)` instead of `os.ReadFile` + `json.Unmarshal`
//   is more memory-efficient for large config files. The Decoder reads directly
//   from the `io.Reader` (the open file handle) in buffered chunks, avoiding
//   the need to load the entire file into a `[]byte` first. For a 10KB config
//   file this doesn't matter, but it establishes the correct habit for when
//   you're processing 100MB JSON API responses in production.
//
// RUN: go run ./09-io-and-cli/encoding/6-config-parser
// ============================================================================

// AppConfig represents the application configuration structure.
// Each field's `json` tag maps it to a key in the JSON config file.
type AppConfig struct {
	AppName      string   `json:"app_name"`
	Port         int      `json:"port"`
	DatabaseURL  string   `json:"database_url"`
	Debug        bool     `json:"debug"`
	MaxWorkers   int      `json:"max_workers"`
	AllowedHosts []string `json:"allowed_hosts"`
}

// validate checks that all required fields are non-zero.
// In Go, a missing JSON field is set to its zero value (empty string, 0, false).
// We use this knowledge to detect incomplete configurations.
func (c *AppConfig) validate() error {
	if c.AppName == "" {
		return fmt.Errorf("config validation: 'app_name' is required")
	}
	if c.Port == 0 {
		return fmt.Errorf("config validation: 'port' is required")
	}
	if c.DatabaseURL == "" {
		return fmt.Errorf("config validation: 'database_url' is required")
	}
	if c.MaxWorkers == 0 {
		return fmt.Errorf("config validation: 'max_workers' is required")
	}
	return nil
}

// loadConfig reads and parses a JSON config file from disk.
func loadConfig(path string) (*AppConfig, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("could not open config file: %w", err)
	}
	defer file.Close()

	var config AppConfig
	// json.NewDecoder reads directly from the file handle (io.Reader)
	// instead of loading the entire file into memory first.
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, fmt.Errorf("could not parse config JSON: %w", err)
	}

	// Validate after parsing
	if err := config.validate(); err != nil {
		return nil, err
	}

	return &config, nil
}

func main() {
	fmt.Println("=== Config File Parser ===")
	fmt.Println()

	// Create a sample config file for the demo
	tmpDir, err := os.MkdirTemp("", "config-demo-*")
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	defer os.RemoveAll(tmpDir)

	configPath := filepath.Join(tmpDir, "config.json")
	sampleConfig := `{
  "app_name": "The Go Engineer API",
  "port": 8080,
  "database_url": "postgres://localhost:5432/mydb",
  "debug": true,
  "max_workers": 4,
  "allowed_hosts": ["localhost", "api.example.com"]
}`
	os.WriteFile(configPath, []byte(sampleConfig), 0644)

	// Load and display the config
	config, err := loadConfig(configPath)
	if err != nil {
		fmt.Printf("❌ Failed to load config: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✅ Config loaded successfully!")
	fmt.Println()
	fmt.Printf("  App Name:      %s\n", config.AppName)
	fmt.Printf("  Port:          %d\n", config.Port)
	fmt.Printf("  Database:      %s\n", config.DatabaseURL)
	fmt.Printf("  Debug Mode:    %v\n", config.Debug)
	fmt.Printf("  Max Workers:   %d\n", config.MaxWorkers)
	fmt.Printf("  Allowed Hosts: %v\n", config.AllowedHosts)

	fmt.Println()
	fmt.Println("KEY TAKEAWAY:")
	fmt.Println("  - json.NewDecoder streams from io.Reader — use for files & HTTP")
	fmt.Println("  - Always validate config after parsing — zero values hide missing fields")
	fmt.Println("  - Use %w in fmt.Errorf to wrap errors for context")
}
