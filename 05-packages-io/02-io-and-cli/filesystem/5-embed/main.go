// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 05: Packages and I/O
// Title: Embed
// Level: Core
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - How to use the '//go:embed' directive to bundle static files into your binary.
//   - How to access embedded assets as strings, byte slices, or a virtual filesystem.
//
// WHY THIS MATTERS:
//   - Embedding allows you to ship a single executable file that includes
//     all its necessary assets (templates, config, data), making deployment
//     and distribution simple and bulletproof.
//
// RUN:
//   go run ./05-packages-io/02-io-and-cli/filesystem/5-embed
//
// KEY TAKEAWAY:
//   - Embedding eliminates the 'missing file' errors common in production deployments.
// ============================================================================

// Commercial use is prohibited without permission.

package main

import (
	"embed"
	"fmt"
	"log"
)

// Stage 05: I/O and CLI - The embed Directive
//
//   - //go:embed for compiling static files directly into a binary
//   - Embedding as string, []byte, or embed.FS
//   - Why embedding helps with single-binary deployment
//
// ENGINEERING DEPTH:
//   The embed package lets Go ship assets without runtime file lookups.
//   That is especially useful for CLIs, internal tools, and compact services
//   that want one binary without a fragile "copy these extra files too" story.

//go:embed public/data.txt
var greeting string

//go:embed public/data.txt
var rawData []byte

//go:embed public/*
var staticFiles embed.FS

func main() {
	fmt.Println("=== Embedding Files inside the Binary ===")
	fmt.Println()

	fmt.Println("1) Embedded String:")
	fmt.Printf("   Greeting: %q\n", greeting)
	fmt.Println()

	fmt.Println("2) Embedded []byte:")
	fmt.Printf("   Bytes: %v\n", rawData)
	fmt.Printf("   Size:  %d bytes\n", len(rawData))
	fmt.Println()

	fmt.Println("3) Reading from embed.FS:")
	content, err := staticFiles.ReadFile("public/data.txt")
	if err != nil {
		log.Fatal("Failed reading from embedded FS:", err)
	}

	fmt.Printf("   public/data.txt content:\n   %s\n", string(content))
	fmt.Println()

	fmt.Println("   Directory contents of 'public':")
	entries, err := staticFiles.ReadDir("public")
	if err != nil {
		log.Fatal(err)
	}

	for _, entry := range entries {
		fmt.Printf("    - %s (IsDir: %t)\n", entry.Name(), entry.IsDir())
	}

	fmt.Println()
	fmt.Println("KEY TAKEAWAYS:")
	fmt.Println("  - //go:embed compiles files directly into your executable")
	fmt.Println("  - The embed package must be imported to use the directive")
	fmt.Println("  - Use string or []byte for single files")
	fmt.Println("  - Use embed.FS for directories and grouped assets")
	fmt.Println("  - Embedding avoids missing-runtime-asset problems")

	fmt.Println("\n---------------------------------------------------")
	fmt.Println("NEXT UP: FS.6 io-patterns")
	fmt.Println("Current: FS.5 (embed)")
	fmt.Println("Previous: FS.4 (temp)")
	fmt.Println("---------------------------------------------------")
}
