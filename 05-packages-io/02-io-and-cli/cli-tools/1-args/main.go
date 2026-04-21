// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

// RUN: go run ./05-packages-io/02-io-and-cli/cli-tools/1-args
package main

import (
	"fmt"
	"os"
	"strings"
)

// ============================================================================
// Stage 05: CLI Tools — Command-Line Arguments
// Level: Beginner
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - os.Args — raw access to command-line arguments
//   - os.Getenv — reading environment variables
//   - os.Exit — exiting with status codes (0=success, 1=error)
//   - Building a simple CLI tool from scratch
//
// ENGINEERING DEPTH:
//   Command line arguments aren't just strings; they are the fundamental UNIX
//   inter-process communication mechanism. When the shell (bash/zsh) launches your
//   compiled Go binary, it invokes the `execve()` system call. The OS kernel reads
//   the arguments you typed, allocates them into the new process's stack memory space,
//   and points Go's `os.Args` string slice directly at those memory addresses. This
//   is why `os.Args[0]` is incredibly rigid and always contains the executable path.
//
// RUN:
//   go run ./05-packages-io/02-io-and-cli/cli-tools/1-args
//   go run ./05-packages-io/02-io-and-cli/cli-tools/1-args hello world
//   GREETING=Hi go run ./05-packages-io/02-io-and-cli/cli-tools/1-args gopher
// ============================================================================

func main() {
	fmt.Println("=== Command-Line Arguments ===")
	fmt.Println()

	// --- os.Args ---
	// os.Args is a []string slice containing all command-line arguments.
	//   os.Args[0] = the program name (or path to the binary)
	//   os.Args[1] = first argument
	//   os.Args[2] = second argument, etc.
	//
	// When you run: go run ./05-packages-io/02-io-and-cli/cli-tools/1-args hello world
	//   os.Args = ["/tmp/go-build.../main", "hello", "world"]
	fmt.Printf("  Program: %s\n", os.Args[0])
	fmt.Printf("  Total args: %d\n", len(os.Args))

	// --- SAFE ARGUMENT ACCESS ---
	// Always check len(os.Args) before accessing arguments.
	// Accessing os.Args[1] when no arguments exist causes a runtime PANIC.
	if len(os.Args) > 1 {
		args := os.Args[1:] // Everything after the program name
		fmt.Printf("  Arguments: %s\n", strings.Join(args, ", "))
	} else {
		fmt.Println("  No arguments provided.")
		fmt.Println("  Try: go run ./05-packages-io/02-io-and-cli/cli-tools/1-args hello world")
	}

	fmt.Println()

	// --- ENVIRONMENT VARIABLES ---
	// os.Getenv reads environment variables — returns "" if not set.
	// Environment variables are the standard way to configure production apps:
	//   DATABASE_URL, API_KEY, PORT, LOG_LEVEL, etc.
	fmt.Println("=== Environment Variables ===")
	greeting := os.Getenv("GREETING")
	if greeting == "" {
		greeting = "Hello" // Default value
	}

	user := os.Getenv("USER") // Most systems set this automatically
	fmt.Printf("  %s, %s!\n", greeting, user)
	fmt.Println("  Try: GREETING=Howdy go run ./05-packages-io/02-io-and-cli/cli-tools/1-args")

	fmt.Println()

	// --- EXIT CODES ---
	// Unix convention:
	//   0 = success (program completed normally)
	//   1 = general error
	//   2 = misuse (wrong arguments)
	//
	// os.Exit(code) terminates the program IMMEDIATELY.
	// Deferred functions DO NOT RUN after os.Exit.
	// Use it at the very end, or for fatal errors only.
	fmt.Println("=== Exit Codes ===")
	fmt.Println("  os.Exit(0) — success")
	fmt.Println("  os.Exit(1) — error")
	fmt.Println("  os.Exit(2) — usage error")

	// Example: validate minimum arguments
	if len(os.Args) > 1 && os.Args[1] == "--fail" {
		fmt.Fprintln(os.Stderr, "Error: intentional failure")
		os.Exit(1) // Non-zero = error
	}

	fmt.Println()
	fmt.Println("KEY TAKEAWAYS:")
	fmt.Println("  1. os.Args[0] is the program name, os.Args[1:] are your arguments")
	fmt.Println("  2. Always check len(os.Args) before accessing specific indices")
	fmt.Println("  3. os.Getenv for config, with sensible defaults for missing vars")
	fmt.Println("  4. Exit 0 = success, Exit 1 = error (defers DON'T run after os.Exit)")
	fmt.Println()
	fmt.Println("   Next: go run ./05-packages-io/02-io-and-cli/cli-tools/2-flags -name='The Go Engineer' -count=3")
	fmt.Println("\n---------------------------------------------------")
	fmt.Println("🚀 NEXT UP: CL.2 flags")
	fmt.Println("   Current: CL.1 (args)")
	fmt.Println("---------------------------------------------------")
}
