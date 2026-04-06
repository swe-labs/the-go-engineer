// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

// RUN: go run ./09-io-and-cli/cli-tools/2-flags
package main

import (
	"flag"
	"fmt"
	"strings"
)

// ============================================================================
// Section 19: CLI Tools — Flag Package
// Level: Intermediate
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - The flag package for typed argument parsing
//   - String, int, bool, and duration flags
//   - Default values and usage help
//   - flag.Parse() — MUST call before accessing flag values
//   - Accessing remaining (non-flag) arguments
//
// ENGINEERING DEPTH:
//   Why does the `flag` package return memory pointers (`*int`, `*string`) instead
//   of raw values? When you define `count := flag.Int(...)`, the struct hasn't parsed
//   the shell input yet! Returning a pointer allows the `flag` package to allocate
//   the destination memory *immediately*, and then asynchronously overwrite the data
//   at those memory addresses later when you invoke `flag.Parse()`. Without pointers,
//   you would be forced to capture every return value from the parse step individually.
//
// RUN:
//   go run ./19-cli-tools/2-flags
//   go run ./19-cli-tools/2-flags -name="The Go Engineer" -count=3 -verbose
//   go run ./19-cli-tools/2-flags -help
// ============================================================================

func main() {
	// --- DEFINING FLAGS ---
	// flag.String returns a *string (pointer). You must dereference with * to get the value.
	// The three arguments are: flag name, default value, usage description.
	name := flag.String("name", "World", "Name to greet")
	count := flag.Int("count", 1, "Number of times to greet")
	verbose := flag.Bool("verbose", false, "Enable verbose output")
	separator := flag.String("sep", "-", "Separator character")

	// --- PARSING FLAGS ---
	// flag.Parse() reads os.Args and sets all flag values.
	// You MUST call this before accessing any flag values.
	// If you forget, all flags will have their default values.
	flag.Parse()

	// --- ACCESSING FLAG VALUES ---
	// Since flag functions return pointers, dereference with *
	if *verbose {
		fmt.Println("=== Verbose Mode Enabled ===")
		fmt.Printf("  Name:      %s\n", *name)
		fmt.Printf("  Count:     %d\n", *count)
		fmt.Printf("  Verbose:   %t\n", *verbose)
		fmt.Printf("  Separator: %q\n", *separator)
		fmt.Println()
	}

	// Use the parsed flags
	line := strings.Repeat(*separator, 30)
	for i := 0; i < *count; i++ {
		if i > 0 {
			fmt.Println(line)
		}
		fmt.Printf("Hello, %s! (greeting #%d)\n", *name, i+1)
	}

	// --- REMAINING ARGUMENTS ---
	// Any arguments after the flags are available via flag.Args()
	// Example: go run . -name=Go extra1 extra2
	//   → flag.Args() = ["extra1", "extra2"]
	remaining := flag.Args()
	if len(remaining) > 0 {
		fmt.Printf("\nExtra arguments: %v\n", remaining)
	}

	// --- AUTO-GENERATED HELP ---
	// Running with -help or -h automatically prints usage information
	// generated from the flag definitions above.
	// Try: go run ./19-cli-tools/2-flags -help

	// KEY TAKEAWAY:
	// - flag package gives you typed CLI arguments with defaults and help text
	// - ALWAYS call flag.Parse() before accessing flag values
	// - flag functions return pointers — dereference with *
	// - -help is auto-generated from your flag definitions
	// - For complex CLIs with subcommands, see 3-subcommands/
	fmt.Println("\n---------------------------------------------------")
	fmt.Println("🚀 NEXT UP: CL.3 subcommands")
	fmt.Println("   Current: CL.2 (flags)")
	fmt.Println("---------------------------------------------------")
}
