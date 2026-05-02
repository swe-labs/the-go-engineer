// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

package main

import "fmt"

// ============================================================================
// Section 04: Types and Design - Config Parser (Exercise Starter)
// Level: Core
// ============================================================================
//
// EXERCISE: Build an application config (.env) parser
//
// REQUIREMENTS:
//  1. [ ] Define a `parseConfig` function returning a `map[string]string`
//  2. [ ] Compile a regex pattern to match key-value config lines like:
//         APP_NAME="My Cool App"
//  3. [ ] Ignore empty lines and comments (`#`)
//  4. [ ] Read a multi-line config string using `bufio.Scanner`
//  5. [ ] Parse each line and collect matches into the resulting map
//  6. [ ] Render a stable summary using `text/template`
//  7. [ ] Make `go test ./04-types-design/24-config-parser-project` pass
//
// HINTS:
//   - Use regexp.MustCompile to compile the pattern at startup
//   - Use strings.NewReader to create an io.Reader from a string
//   - bufio.NewScanner(reader) reads line-by-line
//   - Sort keys before rendering so the template output is stable
//
// RUN: go run ./04-types-design/24-config-parser-project/_starter
// TEST: go test ./04-types-design/24-config-parser-project
// SOLUTION: See the main.go file in the parent directory
// ============================================================================

// Learner task: prepare your mock .env file string.

// Learner task: compile your regex pattern.

// Learner task: implement the parseConfig function.

// Learner task: render a stable summary with text/template.

func main() {
	fmt.Println("=== Config Parser Exercise ===")
	fmt.Println()
	fmt.Println("Implement the config parser and rendered summary.")
	fmt.Println("Use the tests to confirm your parsing logic.")
	fmt.Println()
	fmt.Println("When finished, compare your solution with ../main.go")
}
