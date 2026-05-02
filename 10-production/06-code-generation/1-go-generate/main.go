// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 10: Production Operations
// Title: go generate Primer
// Level: Foundation
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - What the //go:generate directive actually does
//   - Why code generation is a build-time workflow, not runtime reflection
//   - When to generate code vs keep things handwritten
//   - How mockery, stringer, sqlc fit into the same workflow
//
// WHY THIS MATTERS:
//   - Reflection: flexible but fails at runtime.
//   - Generation: heavier build but fails early.
//
// RUN:
//   go run ./10-production/06-code-generation/1-go-generate
// KEY TAKEAWAY:
//   - Go scans for //go:generate, runs tool, writes generated files.
//   - You review and commit generated output intentionally.
// ============================================================================

package main

import "fmt"

// Stage 10: Code Generation - go generate Primer
//
//   - What the //go:generate directive actually does
//   - Why code generation is a build-time workflow, not runtime reflection
//   - When to generate code and when to keep things handwritten
//   - How mockery, stringer, and sqlc fit into the same workflow
//
// ENGINEERING DEPTH:
//   Code generation is one of Go's cleanest scaling tools because it moves
//   repetitive work into build time while still producing ordinary Go source.
//   That means the generated output can be reviewed, benchmarked, and debugged
//   like any other code in the repository.
//
//   The key trade-off is this:
//   - reflection keeps flexibility at runtime but hides failures until runtime
//   - generation makes the build step a little heavier but keeps failures early
//

func main() {
	fmt.Println("=== go generate Primer ===")
	fmt.Println()

	fmt.Println("The //go:generate directive records a tool command next to the code")
	fmt.Println("that depends on the generated output.")
	fmt.Println()

	fmt.Println("Example directives:")
	fmt.Println("  //go:generate mockery --name=Storer --output=../../mocks")
	fmt.Println("  //go:generate stringer -type=Direction")
	fmt.Println("  //go:generate sqlc generate")
	fmt.Println()

	fmt.Println("What happens when you run `go generate ./...`?")
	steps := []string{
		"Go scans source files for //go:generate directives",
		"Each directive is executed as a normal shell command",
		"The tool writes or updates generated Go files",
		"You review and commit the generated output intentionally",
	}
	for i, step := range steps {
		fmt.Printf("  %d. %s\n", i+1, step)
	}

	fmt.Println()
	fmt.Println("When code generation is a good fit:")
	goodFits := []string{
		"Mocks that mirror interfaces and change when the interface changes",
		"String methods for enums that would otherwise be repetitive switches",
		"Typed query code generated from real SQL schemas and queries",
	}
	for _, item := range goodFits {
		fmt.Printf("  - %s\n", item)
	}

	fmt.Println()
	fmt.Println("When NOT to reach for generation first:")
	badFits := []string{
		"A tiny helper you can write clearly by hand in two minutes",
		"Logic that would become harder to review once hidden behind a generator",
		"Cases where the generated output is never read, tested, or committed carefully",
	}
	for _, item := range badFits {
		fmt.Printf("  - %s\n", item)
	}

	fmt.Println()
	fmt.Println("KEY TAKEAWAYS:")
	fmt.Println("  1. `go generate` is an explicit build-time workflow, not part of `go build`")
	fmt.Println("  2. Generated code should still be readable, reviewable, and committed intentionally")
	fmt.Println("  3. Good generators remove repetition without hiding the actual runtime behavior")
	fmt.Println("  4. Mockery, stringer, and sqlc are different tools with the same build-time philosophy")
	fmt.Println()
	fmt.Println("---------------------------------------------------")
	fmt.Println("NEXT UP: CG.2 -> 10-production/06-code-generation/2-mockery")
	fmt.Println("   Current: CG.1 (go generate primer)")
	fmt.Println("---------------------------------------------------")
}
