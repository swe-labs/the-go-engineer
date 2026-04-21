// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

// RUN: go run ./05-packages-io/02-io-and-cli/cli-tools/2-flags
package main

import (
	"flag"
	"fmt"
	"strings"
)

// ============================================================================
// Stage 05: CLI Tools — Flag Package
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
//   the shell input yet. Returning a pointer allows the `flag` package to allocate
//   the destination memory immediately, then overwrite that data later when you
//   invoke `flag.Parse()`.
//
// RUN:
//   go run ./05-packages-io/02-io-and-cli/cli-tools/2-flags
//   go run ./05-packages-io/02-io-and-cli/cli-tools/2-flags -name="The Go Engineer" -count=3 -verbose
//   go run ./05-packages-io/02-io-and-cli/cli-tools/2-flags -help
// ============================================================================

func main() {
	name := flag.String("name", "World", "Name to greet")
	count := flag.Int("count", 1, "Number of times to greet")
	verbose := flag.Bool("verbose", false, "Enable verbose output")
	separator := flag.String("sep", "-", "Separator character")

	flag.Parse()

	if *verbose {
		fmt.Println("=== Verbose Mode Enabled ===")
		fmt.Printf("  Name:      %s\n", *name)
		fmt.Printf("  Count:     %d\n", *count)
		fmt.Printf("  Verbose:   %t\n", *verbose)
		fmt.Printf("  Separator: %q\n", *separator)
		fmt.Println()
	}

	line := strings.Repeat(*separator, 30)
	for i := 0; i < *count; i++ {
		if i > 0 {
			fmt.Println(line)
		}
		fmt.Printf("Hello, %s! (greeting #%d)\n", *name, i+1)
	}

	remaining := flag.Args()
	if len(remaining) > 0 {
		fmt.Printf("\nExtra arguments: %v\n", remaining)
	}

	fmt.Println("\n---------------------------------------------------")
	fmt.Println("🚀 NEXT UP: CL.3 subcommands")
	fmt.Println("   Current: CL.2 (flags)")
	fmt.Println("---------------------------------------------------")
}
