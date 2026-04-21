// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

// RUN: go run ./05-packages-io/02-io-and-cli/cli-tools/3-subcommands
package main

import (
	"flag"
	"fmt"
	"os"
)

// ============================================================================
// Stage 05: I/O and CLI — Subcommands
// Level: Advanced
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Building multi-command CLIs like git, docker, and kubectl
//   - flag.NewFlagSet for independent flag sets per subcommand
//   - Routing subcommands with switch on os.Args[1]
//   - Error handling and usage messages for invalid subcommands
//
// ENGINEERING DEPTH:
//   Large CLIs rarely parse every option globally. They route to a subcommand,
//   then hand only the remaining arguments to a command-specific flag set.
//   That keeps each command's options isolated and avoids accidental flag-name
//   collisions across the whole tool.
//
// RUN:
//   go run ./05-packages-io/02-io-and-cli/cli-tools/3-subcommands greet -name="Gopher"
//   go run ./05-packages-io/02-io-and-cli/cli-tools/3-subcommands version
//   go run ./05-packages-io/02-io-and-cli/cli-tools/3-subcommands calc -a=10 -b=20
// ============================================================================

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

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
	fmt.Println("🚀 NEXT UP: CL.4 file organizer")
	fmt.Println("   Current: CL.3 (subcommands)")
	fmt.Println("---------------------------------------------------")
}

func cmdGreet(args []string) {
	fs := flag.NewFlagSet("greet", flag.ExitOnError)
	name := fs.String("name", "World", "Name to greet")
	loud := fs.Bool("loud", false, "Shout the greeting")

	fs.Parse(args)

	greeting := fmt.Sprintf("Hello, %s!", *name)
	if *loud {
		greeting = fmt.Sprintf("HELLO, %s!!!", *name)
	}

	fmt.Println(greeting)
}

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
