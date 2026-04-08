// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

// RUN: go run ./08-modules-and-packages/1-module-basics
package main

import "fmt"

// ============================================================================
// Section 08: Modules and Packages — Module Basics
// Level: Beginner → Intermediate
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - What go.mod is and why it exists
//   - How Go resolves import paths
//   - The anatomy of go.mod and go.sum
//
// KEY COMMANDS:
//   go mod init <module-path>    — Create a new module
//   go mod tidy                  — Remove unused, add missing dependencies
//   go list -m all               — List all direct + indirect dependencies
//   go mod why <module>          — Explain why a dependency is needed
//
// REFERENCES:
//   - https://go.dev/blog/using-go-modules
//   - https://go.dev/ref/mod
// ============================================================================

// GoModAnatomy explains the structure of go.mod using this repo as an example.
//
// File: go.mod
//
//	module github.com/rasel9t6/the-go-engineer ← Module path: the root import path for all packages
//	go 1.24                                    ← Minimum Go version required to build
//
//	require (
//	    github.com/mattn/go-sqlite3 v1.14.28  ← Direct dependency
//	    github.com/stretchr/testify v1.11.1   ← Direct dependency
//	    golang.org/x/crypto v0.39.0           ← Direct dependency
//	)
//
//	require (
//	    github.com/davecgh/go-spew v1.1.1 // indirect ← Transitive dependency
//	    github.com/stretchr/objx v0.5.2 // indirect   ← Pulled by testify/mock
//	)
//
// IMPORTANT CONCEPTS:
//
// 1. MODULE PATH is the identity of your module.
//    - For published modules: use your repo URL (e.g., github.com/user/project)
//    - For local-only modules: any name works (e.g., "the-go-engineer")
//
// 2. SEMANTIC VERSIONING: vMAJOR.MINOR.PATCH
//    - MAJOR: breaking changes (v2 requires import path change: /v2)
//    - MINOR: new features, backward compatible
//    - PATCH: bug fixes only
//
// 3. go.sum contains checksums for reproducible builds.
//    - NEVER edit go.sum manually — it's auto-generated
//    - ALWAYS commit go.sum to version control
//
// 4. // indirect means YOUR code doesn't import it directly.
//    - It came from one of your dependencies.
//    - go mod tidy adds this annotation automatically.

func main() {
	fmt.Println("=== Go Module Basics ===")
	fmt.Println()

	fmt.Println("Module path: github.com/rasel9t6/the-go-engineer")
	fmt.Println("This means packages in this repo are importable as:")
	fmt.Println("  github.com/rasel9t6/the-go-engineer/10-web-and-database/databases/6-repository/models")
	fmt.Println("  github.com/rasel9t6/the-go-engineer/10-web-and-database/databases/6-repository/repository")
	fmt.Println()

	commands := []struct {
		cmd  string
		desc string
	}{
		{"go mod init <path>", "Initialize a new module"},
		{"go mod tidy", "Sync go.mod/go.sum with actual imports"},
		{"go get pkg@version", "Add or update a dependency"},
		{"go get pkg@none", "Remove a dependency"},
		{"go list -m all", "List all dependencies (direct + indirect)"},
		{"go list -m -versions pkg", "List available versions of a module"},
		{"go mod why pkg", "Explain why a dependency is needed"},
		{"go mod graph", "Print the module dependency graph"},
		{"go mod verify", "Verify dependencies match go.sum checksums"},
		{"go mod download", "Download modules to local cache"},
	}

	fmt.Println("Essential Module Commands:")
	for _, c := range commands {
		fmt.Printf("  %-32s — %s\n", c.cmd, c.desc)
	}

	fmt.Println("\n---------------------------------------------------")
	fmt.Println("🚀 NEXT UP: MP.2 managing deps")
	fmt.Println("   Current: MP.1 (module basics)")
	fmt.Println("---------------------------------------------------")
}
