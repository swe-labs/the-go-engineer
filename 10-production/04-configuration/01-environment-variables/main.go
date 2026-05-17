// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 10: Production Operations
// Title: Environment variables
// Level: Core
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Learn how environment variables shape runtime configuration without rebuilding the binary.
//
// WHY THIS MATTERS:
//   - Environment variables are process-level inputs provided by the runtime environment, not by the source code itself.
//
// RUN:
//   go run ./10-production/04-configuration/01-environment-variables
//
// KEY TAKEAWAY:
//   - Environment variables are late-bound process configuration.
//   - Missing or malformed values should fail fast.
//   - Keep names stable and documented.
// ============================================================================

package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

// Config (Struct): aggregates configuration values read from environment variables.
type Config struct {
	Port   int
	DBHost string
	DBUser string
	Debug  bool
}

// loadConfig (Function): reads configuration from environment variables with defaults.
// Uses os.Getenv for optional values and os.LookupEnv for values that must be present.
func loadConfig() Config {
	portStr := os.Getenv("APP_PORT")
	port := 8080
	if portStr != "" {
		if p, err := strconv.Atoi(portStr); err == nil {
			port = p
		}
	}

	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "localhost"
	}

	// os.LookupEnv distinguishes empty from unset
	dbUser, userSet := os.LookupEnv("DB_USER")
	if !userSet {
		dbUser = "default_user"
	}

	debug := os.Getenv("DEBUG") == "true"

	return Config{
		Port:   port,
		DBHost: dbHost,
		DBUser: dbUser,
		Debug:  debug,
	}
}

func main() {
	cfg := loadConfig()

	fmt.Println("=== CFG.1 Environment variables ===")
	fmt.Println()
	fmt.Println("Config loaded from environment:")
	fmt.Printf("  APP_PORT=%d      (os.Getenv with Atoi, default 8080)\n", cfg.Port)
	fmt.Printf("  DB_HOST=%s       (os.Getenv with empty check, default localhost)\n", cfg.DBHost)
	fmt.Printf("  DB_USER=%s       (os.LookupEnv distinguishes empty from unset)\n", cfg.DBUser)
	fmt.Printf("  DEBUG=%v         (os.Getenv == \"true\")\n", cfg.Debug)
	fmt.Println()
	fmt.Println("Try setting environment variables before running:")
	log.Println("Example: APP_PORT=9090 DB_HOST=prod.example.com DB_USER=admin DEBUG=true go run ./10-production/04-configuration/01-environment-variables")
	fmt.Println()
	fmt.Println("---------------------------------------------------")
	fmt.Println("NEXT UP: CFG.2 -> 10-production/04-configuration/02-configuration-files")
	fmt.Println("Current: CFG.1 (environment variables)")
	fmt.Println("---------------------------------------------------")
}
