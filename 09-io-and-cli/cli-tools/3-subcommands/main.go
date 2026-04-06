// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

// RUN: go run ./09-io-and-cli/cli-tools/3-subcommands
package main

import (
	"flag"
	"fmt"
	"os"
)

// ============================================================================
// Section 19: CLI Tools — Subcommands
// Level: Advanced
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Building multi-command CLIs like git, docker, kubectl
//   - flag.NewFlagSet for independent flag sets per subcommand
//   - Routing subcommands with switch on os.Args[1]
//   - Error handling and usage messages for invalid subcommands
//
// ENGINEERING DEPTH:
//   Massive multi-tool CLIs like `kubectl` or `docker` do not parse global flags.
//   Instead, they act as "Routers". By slicing `os.Args[2:]`, they mathematically
//   sever the initial subcommand string from the underlying array, passing only the
//   remaining flags down into an isolated runtime scope. This prevents fatal namespace
//   collisions where `docker run --help` might overlap with a global `--help` flag,
//   allowing independent validation lifetimes for every subcommand.
//
// RUN:
//   go run ./19-cli-tools/3-subcommands greet -name="Gopher"
//   go run ./19-cli-tools/3-subcommands version
//   go run ./19-cli-tools/3-subcommands calc -a=10 -b=20
// ============================================================================

func main() {
	// --- SUBCOMMAND ROUTING ---
	// Production CLIs use this pattern:
	//   myapp <subcommand> [flags]
	//   git commit -m "message"
	//   docker run --name mycontainer
	//   kubectl get pods -n production
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	// os.Args[1] is the subcommand name
	switch os.Args[1] {
	case "greet":
		cmdGreet(os.Args[2:])
	case "calc":
		cmdCalc(os.Args[2:])
	case "version":
		cmdVersion()
	default:
		fmt.Fprintf(os.Stderr, "Error: unknown subcommand %q\n\n", os.Args[1])
		printUsage()
		os.Exit(1)
	}
	fmt.Println("\n---------------------------------------------------")
	fmt.Println("🚀 NEXT UP: DB.1 connecting")
	fmt.Println("   Current: CL.3 (subcommands)")
	fmt.Println("---------------------------------------------------")
}

// --- SUBCOMMAND: greet ---
func cmdGreet(args []string) {
	// flag.NewFlagSet creates an INDEPENDENT flag set for this subcommand.
	// Each subcommand has its own flags that don't interfere with others.
	fs := flag.NewFlagSet("greet", flag.ExitOnError)
	name := fs.String("name", "World", "Name to greet")
	loud := fs.Bool("loud", false, "Shout the greeting")

	// Parse ONLY the arguments for this subcommand (not os.Args)
	fs.Parse(args)

	greeting := fmt.Sprintf("Hello, %s!", *name)
	if *loud {
		greeting = fmt.Sprintf("HELLO, %s!!!", *name)
	}
	fmt.Println(greeting)
}

// --- SUBCOMMAND: calc ---
func cmdCalc(args []string) {
	fs := flag.NewFlagSet("calc", flag.ExitOnError)
	a := fs.Int("a", 0, "First number")
	b := fs.Int("b", 0, "Second number")
	op := fs.String("op", "add", "Operation: add, sub, mul")

	fs.Parse(args)

	var result int
	switch *op {
	case "add":
		result = *a + *b
	case "sub":
		result = *a - *b
	case "mul":
		result = *a * *b
	default:
		fmt.Fprintf(os.Stderr, "Unknown operation: %s\n", *op)
		os.Exit(1)
	}
	fmt.Printf("%d %s %d = %d\n", *a, *op, *b, result)
}

// --- SUBCOMMAND: version ---
func cmdVersion() {
	fmt.Println("The Go Engineer CLI v1.0.0")
	fmt.Printf("Built with: %s\n", "go1.26")
}

func printUsage() {
	fmt.Println("Usage: program <subcommand> [flags]")
	fmt.Println()
	fmt.Println("Subcommands:")
	fmt.Println("  greet     Greet someone by name")
	fmt.Println("  calc      Perform arithmetic operations")
	fmt.Println("  version   Print version information")
	fmt.Println()
	fmt.Println("Run 'program <subcommand> -help' for subcommand-specific flags.")
}
