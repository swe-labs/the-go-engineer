// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// ============================================================================
// Section 0: Getting Started — Development Environment
// Level: Beginner
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Essential Go CLI tools every developer must know
//   - What "go fmt", "go vet", and "go build" do
//   - How to verify your editor is configured correctly
//   - The Go workspace model
//
// RUN: go run ./00-getting-started/4-dev-environment
// ============================================================================

func main() {
	fmt.Println("=== Go Development Environment Check ===")
	fmt.Println()

	// --- 1. SYSTEM INFORMATION ---
	fmt.Println("📦 System Information")
	fmt.Println(strings.Repeat("─", 50))
	fmt.Printf("  Go Version:     %s\n", runtime.Version())
	fmt.Printf("  OS/Arch:        %s/%s\n", runtime.GOOS, runtime.GOARCH)
	fmt.Printf("  Num CPUs:       %d\n", runtime.NumCPU())
	//lint:ignore SA1019 runtime.GOROOT is deprecated in Go 1.24 but helpful for beginners to see.
	fmt.Printf("  GOROOT:         %s\n", runtime.GOROOT())
	fmt.Printf("  GOPATH:         %s\n", os.Getenv("GOPATH"))
	fmt.Println()

	// --- 2. ESSENTIAL GO COMMANDS ---
	// These are the commands you'll use every day as a Go engineer.
	fmt.Println("🔧 Essential Go Commands")
	fmt.Println(strings.Repeat("─", 50))

	commands := []struct {
		cmd  string
		desc string
	}{
		// go fmt — The Formatter
		// Automatically formats your code to Go's standard style.
		// This is NOT optional. In Go, EVERYONE uses the same formatting.
		// No tabs vs spaces debate. No brace style arguments. Just run go fmt.
		// Your editor should run this automatically on save.
		{"go fmt ./...", "Format all Go files (standard style, non-negotiable)"},

		// go vet — The Static Analyzer
		// Finds bugs that the compiler doesn't catch: suspicious constructs,
		// unreachable code, incorrect format strings. It's like a spell-checker
		// for your code logic.
		{"go vet ./...", "Find suspicious code patterns (bugs the compiler misses)"},

		// go build — The Compiler
		// Compiles your code into a native binary. If it succeeds with no output,
		// your code is free of syntax and type errors.
		{"go build ./...", "Compile all packages (verify everything compiles)"},

		// go run — The Quick Runner
		// Compiles and runs in one step. Great for development.
		// Does NOT save the binary — it's temporary.
		{"go run ./path/to/pkg", "Compile and run (development shortcut)"},

		// go test — The Test Runner
		// Runs all files ending in _test.go. Go has testing built into the language.
		// No need for Jest, pytest, or JUnit — it's all in the box.
		{"go test ./...", "Run all tests in the project"},

		// go test -race — The Race Detector
		// Finds data races in concurrent code. ALWAYS use this in CI/CD.
		{"go test -race ./...", "Run tests with race condition detection"},

		// go mod tidy — The Dependency Cleaner
		// Adds missing imports to go.mod and removes unused ones.
		// Run this whenever you add or remove an import.
		{"go mod tidy", "Sync go.mod with actual imports"},
	}

	for _, c := range commands {
		fmt.Printf("  %-28s — %s\n", c.cmd, c.desc)
	}
	fmt.Println()

	// --- 3. CHECK EDITOR TOOLS ---
	fmt.Println("🔍 Checking Go Tools")
	fmt.Println(strings.Repeat("─", 50))

	// These tools power your editor's Go features.
	// They should have been installed when you set up the Go extension.
	tools := []struct {
		name string
		desc string
	}{
		{"gopls", "Go Language Server (autocomplete, go-to-definition, errors)"},
		{"gofmt", "Code formatter (built into Go)"},
	}

	for _, tool := range tools {
		path, err := exec.LookPath(tool.name)
		if err != nil {
			fmt.Printf("  ❌ %s — NOT FOUND (%s)\n", tool.name, tool.desc)
			fmt.Printf("     Install: go install golang.org/x/tools/gopls@latest\n")
		} else {
			fmt.Printf("  ✅ %s — Found at %s\n", tool.name, path)
		}
	}
	fmt.Println()

	// --- 4. THE "./..." PATTERN ---
	// You'll see "./..." used constantly in Go commands. Here's what it means:
	//
	//   ./          = the current directory
	//   ...         = and all subdirectories, recursively
	//   ./...       = every Go package from here downward
	//
	// So "go build ./..." means "compile every package in the entire project".
	// And "go test ./..." means "run all tests in the entire project".
	//
	// This is one of Go's most useful conventions. One command, entire project.

	fmt.Println("KEY TAKEAWAYS:")
	fmt.Println("  1. go fmt — Format code (your editor should do this on save)")
	fmt.Println("  2. go vet — Find bugs (run before every commit)")
	fmt.Println("  3. go build — Verify code compiles (catches all type errors)")
	fmt.Println("  4. go test — Run tests (the foundation of reliable software)")
	fmt.Println("  5. ./... — Means 'everything in this project, recursively'")
	fmt.Println()
	fmt.Println("🎉 Section 0 complete! You're ready to learn Go.")
	fmt.Println("   Next: go run ./01-language-basics/1-variables")
}
