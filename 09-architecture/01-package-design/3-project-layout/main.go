// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 09: Architecture & Security
// Title: Project Layout
// Level: Core
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Decide when a Go project should stay flat, when `internal/` is useful,
//     and when `cmd/` or `pkg/` have earned their place.
//
// WHY THIS MATTERS:
//   - Project layout communicates ownership, import boundaries, and build
//     intent before a reader opens the first source file.
//
// RUN:
//   go run ./09-architecture/01-package-design/3-project-layout
//
// KEY TAKEAWAY:
//   - Start simple, add structure when it protects boundaries, and avoid
//     copying large layouts before the project needs them.
// ============================================================================

// Commercial use is prohibited without permission.

package main

import "fmt"

// Stage 09: Application Architecture - Package Design: Standard Go Project Layout
//
//   - The common directory conventions for Go projects
//   - When to use cmd/, internal/, pkg/, and other directories
//   - Simple vs complex project layouts
//   - Anti-patterns to avoid
//
// ENGINEERING DEPTH:
//   Go does not enforce one project structure, but the community has converged
//   on a small set of practical conventions. The most important one is
//   `internal/`, because the compiler uses it to block external imports.
//

func main() {
	fmt.Println("=== Standard Go Project Layout ===")
	fmt.Println()

	fmt.Println("SIMPLE PROJECT (library or small app)")
	fmt.Println("   Perfect for: libraries, small tools, learning")
	fmt.Println()
	fmt.Println("   mylib/")
	fmt.Println("   |-- go.mod")
	fmt.Println("   |-- go.sum")
	fmt.Println("   |-- mylib.go          <- main package code")
	fmt.Println("   |-- mylib_test.go     <- tests alongside code")
	fmt.Println("   `-- README.md")
	fmt.Println()

	fmt.Println("MEDIUM PROJECT (single binary plus packages)")
	fmt.Println("   Perfect for: web servers, CLI tools, microservices")
	fmt.Println()
	fmt.Println("   myapp/")
	fmt.Println("   |-- go.mod")
	fmt.Println("   |-- main.go           <- entry point (package main)")
	fmt.Println("   |-- internal/         <- private packages (compiler-enforced)")
	fmt.Println("   |   |-- auth/")
	fmt.Println("   |   |   |-- auth.go")
	fmt.Println("   |   |   `-- auth_test.go")
	fmt.Println("   |   |-- store/")
	fmt.Println("   |   |   |-- store.go")
	fmt.Println("   |   |   `-- store_test.go")
	fmt.Println("   |   `-- handler/")
	fmt.Println("   |       |-- handler.go")
	fmt.Println("   |       `-- handler_test.go")
	fmt.Println("   |-- Makefile")
	fmt.Println("   `-- README.md")
	fmt.Println()

	fmt.Println("LARGE PROJECT (multiple binaries, shared code)")
	fmt.Println("   Perfect for: monorepos, complex systems, multi-service apps")
	fmt.Println()
	fmt.Println("   platform/")
	fmt.Println("   |-- go.mod")
	fmt.Println("   |-- cmd/              <- each subdirectory is a binary")
	fmt.Println("   |   |-- api/")
	fmt.Println("   |   |   `-- main.go   <- go build ./cmd/api")
	fmt.Println("   |   |-- worker/")
	fmt.Println("   |   |   `-- main.go   <- go build ./cmd/worker")
	fmt.Println("   |   `-- migrate/")
	fmt.Println("   |       `-- main.go   <- go build ./cmd/migrate")
	fmt.Println("   |-- internal/         <- private shared packages")
	fmt.Println("   |   |-- auth/")
	fmt.Println("   |   |-- store/")
	fmt.Println("   |   `-- email/")
	fmt.Println("   |-- pkg/              <- public shared packages")
	fmt.Println("   |   `-- middleware/")
	fmt.Println("   |-- migrations/       <- SQL migration files")
	fmt.Println("   |-- Makefile")
	fmt.Println("   |-- Dockerfile")
	fmt.Println("   `-- README.md")
	fmt.Println()

	fmt.Println("=== Directory Purpose Guide ===")
	dirs := []struct {
		dir  string
		use  string
		when string
	}{
		{"cmd/", "Entry points (package main)", "Multiple binaries in one repo"},
		{"internal/", "Private packages", "Code that must not be imported externally"},
		{"pkg/", "Public shared packages", "Code designed for reuse by other modules"},
		{"migrations/", "Database migration files", "Apps with SQL databases"},
		{"testdata/", "Test fixtures", "Go test tooling ignores this directory"},
		{"scripts/", "Build/deploy/CI scripts", "Automation beyond Makefile"},
		{"docs/", "Additional documentation", "Complex projects needing detailed docs"},
	}

	for _, d := range dirs {
		fmt.Printf("  %-15s - %s\n", d.dir, d.use)
		fmt.Printf("  %15s   When: %s\n", "", d.when)
		fmt.Println()
	}

	fmt.Println("=== Anti-Patterns (Avoid These) ===")
	fmt.Println("  [avoid] src/                - Go does not use src/ the way Java does")
	fmt.Println("  [avoid] models/             - Too vague. Name by domain: user/, order/")
	fmt.Println("  [avoid] utils/ or helpers/  - Junk drawer. Split by responsibility")
	fmt.Println("  [avoid] Over-engineering    - Do not use cmd/internal/pkg/ for a 200-line app")
	fmt.Println()

	fmt.Println("KEY TAKEAWAYS:")
	fmt.Println("  1. Start simple and add structure as the project grows")
	fmt.Println("  2. cmd/ is for multiple binaries; internal/ is for private packages")
	fmt.Println("  3. pkg/ is optional - use it only for code meant to be shared publicly")
	fmt.Println("  4. Tests live next to the code they verify (user.go -> user_test.go)")
	fmt.Println("  5. Do not copy a complex layout for a small project")
	fmt.Println("\n---------------------------------------------------")
	fmt.Println("NEXT UP: ARCH.1")
	fmt.Println("Current: PD.3 (project layout)")
	fmt.Println("Previous: PD.2 (visibility)")
	fmt.Println("---------------------------------------------------")
}
