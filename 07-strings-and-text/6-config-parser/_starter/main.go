// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

package main

import "fmt"

// ============================================================================
// Section 7: Strings & Text — Config Parser (Exercise Starter)
// Level: Intermediate
// ============================================================================
//
// EXERCISE: Build an Application Config (.env) Parser
//
// REQUIREMENTS:
//  1. [ ] Define a `parseConfig` function returning a `map[string]string`
//  2. [ ] Compile a regex pattern to match key-value config lines like:
//         APP_NAME="My Cool App"
//  3. [ ] Ignore empty lines and comments (`#`)
//  4. [ ] Read a multi-line config string using `bufio.Scanner`
//  5. [ ] Parse each line, collect matches, and print the resulting map
//
// HINTS:
//   - Use regexp.MustCompile to compile the pattern at startup
//   - Use strings.NewReader to create an io.Reader from a string
//   - bufio.NewScanner(reader) reads line-by-line
//
// RUN: go run ./07-strings-and-text/6-config-parser/_starter
// SOLUTION: See the main.go file in the parent directory
// ============================================================================

// TODO: Prepare your mock .env file string

// TODO: Compile your regex pattern

// TODO: Implement parseConfig function

func main() {
	fmt.Println("=== Config Parser Exercise ===")
	fmt.Println()
	fmt.Println("TODO: Implement your config parser!")
	fmt.Println("See the REQUIREMENTS above for what to build.")
	fmt.Println()
	fmt.Println("When finished, compare your solution with ../main.go")
}
