// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

package main

import "fmt"

// ============================================================================
// Section 09: Encoding - Config File Parser (Exercise Starter)
// Level: Intermediate
// ============================================================================
//
// EXERCISE: Build a JSON Config File Parser
//
// REQUIREMENTS:
//  1. [ ] Define an `AppConfig` struct with json tags for: app_name, port,
//         database_url, debug, max_workers, allowed_hosts
//  2. [ ] Implement `validate()` on AppConfig to check required fields
//  3. [ ] Implement `loadConfig(path string) (*AppConfig, error)` that:
//         - Opens a file using os.Open
//         - Uses json.NewDecoder(file).Decode() to parse
//         - Validates after parsing
//  4. [ ] In main(), create a sample JSON config file, load it, and print the values
//
// RUN: go run ./09-io-and-cli/encoding/6-config-parser/_starter
// SOLUTION: See the main.go file in the parent directory
// ============================================================================

// TODO: Define AppConfig struct with json tags

// TODO: Implement validate method

// TODO: Implement loadConfig function

func main() {
	fmt.Println("=== Config File Parser Exercise ===")
	fmt.Println()
	fmt.Println("TODO: Implement your config parser!")
	fmt.Println("See the REQUIREMENTS above for what to build.")
	fmt.Println()
	fmt.Println("When finished, compare your solution with ../main.go")
}
