// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

package main

import "fmt"

// ============================================================================
// Section 18: Package Design — Standard Go Project Layout
// Level: Intermediate → Advanced
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - The standard directory conventions for Go projects
//   - When to use cmd/, internal/, pkg/, and other directories
//   - Simple vs complex project layouts
//   - Anti-patterns to avoid
//
// ENGINEERING DEPTH:
//   Unlike Django or Spring Boot, Go does NOT enforce a project structure.
//   However, the community standardized the "golang-standards/project-layout".
//   The critical mechanics here involve the `internal/` directory. If a package
//   is inside `internal/`, the Go Compiler's import graph algorithm explicitly
//   blocks *any* repository from importing it outside of the direct parent tree.
//   This gives you absolute guarantees that internal microservices or database
//   layer structs won't leak into the public API boundaries of users importing
//   your module!
//
// RUN: go run ./18-package-design/3-project-layout
// ============================================================================

func main() {
	fmt.Println("=== Standard Go Project Layout ===")
	fmt.Println()

	// --- SIMPLE PROJECT (Most projects start here) ---
	fmt.Println("📁 SIMPLE PROJECT (library or small app)")
	fmt.Println("   Perfect for: libraries, small tools, learning")
	fmt.Println()
	fmt.Println("   mylib/")
	fmt.Println("   ├── go.mod")
	fmt.Println("   ├── go.sum")
	fmt.Println("   ├── mylib.go          ← Main package code")
	fmt.Println("   ├── mylib_test.go     ← Tests alongside code")
	fmt.Println("   └── README.md")
	fmt.Println()

	// --- MEDIUM PROJECT (cmd + internal) ---
	fmt.Println("📁 MEDIUM PROJECT (single binary + some packages)")
	fmt.Println("   Perfect for: web servers, CLI tools, microservices")
	fmt.Println()
	fmt.Println("   myapp/")
	fmt.Println("   ├── go.mod")
	fmt.Println("   ├── main.go           ← Entry point (package main)")
	fmt.Println("   ├── internal/         ← Private packages (compiler-enforced)")
	fmt.Println("   │   ├── auth/")
	fmt.Println("   │   │   ├── auth.go")
	fmt.Println("   │   │   └── auth_test.go")
	fmt.Println("   │   ├── store/")
	fmt.Println("   │   │   ├── store.go")
	fmt.Println("   │   │   └── store_test.go")
	fmt.Println("   │   └── handler/")
	fmt.Println("   │       ├── handler.go")
	fmt.Println("   │       └── handler_test.go")
	fmt.Println("   ├── Makefile")
	fmt.Println("   └── README.md")
	fmt.Println()

	// --- LARGE PROJECT (multiple binaries) ---
	fmt.Println("📁 LARGE PROJECT (multiple binaries, shared code)")
	fmt.Println("   Perfect for: monorepos, complex systems, multi-service apps")
	fmt.Println()
	fmt.Println("   platform/")
	fmt.Println("   ├── go.mod")
	fmt.Println("   ├── cmd/              ← Each subdirectory is a binary")
	fmt.Println("   │   ├── api/")
	fmt.Println("   │   │   └── main.go   ← go build ./cmd/api")
	fmt.Println("   │   ├── worker/")
	fmt.Println("   │   │   └── main.go   ← go build ./cmd/worker")
	fmt.Println("   │   └── migrate/")
	fmt.Println("   │       └── main.go   ← go build ./cmd/migrate")
	fmt.Println("   ├── internal/         ← Private shared packages")
	fmt.Println("   │   ├── auth/")
	fmt.Println("   │   ├── store/")
	fmt.Println("   │   └── email/")
	fmt.Println("   ├── pkg/              ← Public shared packages (importable by others)")
	fmt.Println("   │   └── middleware/")
	fmt.Println("   ├── migrations/       ← SQL migration files")
	fmt.Println("   ├── Makefile")
	fmt.Println("   ├── Dockerfile")
	fmt.Println("   └── README.md")
	fmt.Println()

	// --- KEY DIRECTORIES EXPLAINED ---
	fmt.Println("=== Directory Purpose Guide ===")
	dirs := []struct {
		dir  string
		use  string
		when string
	}{
		{"cmd/", "Entry points (package main)", "Multiple binaries in one repo"},
		{"internal/", "Private packages", "Code that MUST NOT be imported externally"},
		{"pkg/", "Public shared packages", "Code designed for reuse by other modules"},
		{"migrations/", "Database migration files", "Apps with SQL databases"},
		{"testdata/", "Test fixtures", "Go test tooling ignores this directory"},
		{"scripts/", "Build/deploy/CI scripts", "Automation beyond Makefile"},
		{"docs/", "Additional documentation", "Complex projects needing detailed docs"},
	}

	for _, d := range dirs {
		fmt.Printf("  %-15s — %s\n", d.dir, d.use)
		fmt.Printf("  %15s   When: %s\n", "", d.when)
		fmt.Println()
	}

	// --- ANTI-PATTERNS ---
	fmt.Println("=== Anti-Patterns (Avoid These) ===")
	fmt.Println("  ❌ src/               — This isn't Java. Go doesn't use src/")
	fmt.Println("  ❌ models/            — Too vague. Name by domain: user/, order/")
	fmt.Println("  ❌ utils/ or helpers/  — Junk drawer. Split by responsibility")
	fmt.Println("  ❌ Over-engineering   — Don't use cmd/internal/pkg/ for a 200-line app")
	fmt.Println()

	fmt.Println("KEY TAKEAWAYS:")
	fmt.Println("  1. Start simple (flat) — add structure as the project grows")
	fmt.Println("  2. cmd/ for multiple binaries, internal/ for private packages")
	fmt.Println("  3. pkg/ is optional — only for code you intend to share publicly")
	fmt.Println("  4. Tests live NEXT TO the code they test (user.go → user_test.go)")
	fmt.Println("  5. Don't cargo-cult complex layouts for small projects")
}
