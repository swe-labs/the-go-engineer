// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

package main

import (
	"embed"
	"fmt"
	"log"
)

// ============================================================================
// Section 10: Filesystem — The embed Directive
// Level: Advanced
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - //go:embed: compiling static files directly into your Go binary
//   - Embedding as a string or []byte
//   - embed.FS (File System) for embedding whole directories
//   - Why and when to use embedding (single-binary deployments)
//
// ANALOGY:
//   Normal files are like putting a textbook in your backpack — if you
//   forget it, you can't read it.
//
//   Embedded files are like memorizing the textbook. The data is stored
//   inside your brain (the compiled binary). You never need the external
//   file again because it's baked right in.
//
// ENGINEERING DEPTH:
//   Added in Go 1.16, `embed` solved a massive pain point. Previously,
//   developers used 3rd party tools (like `packr` or `go-bindata`) to convert
//   files into giant Go string constants.
//   Now, `go build` does this natively. The resulting binary is completely
//   standalone. You can ship a web server with all HTML, CSS, and JS
//   inside a single executable! No more "missing template" errors in production.
//
// RUN: go run ./10-filesystem/5-embed
// ============================================================================

// --- 1. Embed as a string ---
// The directive MUST be immediately above the variable declaration
// with NO empty lines between them.

//go:embed public/data.txt
var greeting string

// --- 2. Embed as a byte slice ---
// Useful for images, zip files, or any binary data.

//go:embed public/data.txt
var rawData []byte

// --- 3. Embed a whole directory (embed.FS) ---
// This embeds the entire "public" directory and all its contents.
// embed.FS implements the standard `fs.FS` interface, meaning it can be passed
// directly to functions like `http.FileServer`!

//go:embed public/*
var staticFiles embed.FS

func main() {
	fmt.Println("=== Embedding Files inside the Binary ===")
	fmt.Println()

	// =====================================================================
	// 1. Reading an embedded string
	// =====================================================================
	fmt.Println("1️⃣  Embedded String:")
	// The file public/data.txt was read at COMPILE time.
	// We don't need os.Open() — the data is just a variable!
	fmt.Printf("   Greeting: %q\n", greeting)
	fmt.Println()

	// =====================================================================
	// 2. Reading an embedded []byte
	// =====================================================================
	fmt.Println("2️⃣  Embedded []byte:")
	fmt.Printf("   Bytes: %v\n", rawData)
	fmt.Printf("   Size:  %d bytes\n", len(rawData))
	fmt.Println()

	// =====================================================================
	// 3. Using embed.FS (Virtual Filesystem)
	// =====================================================================
	fmt.Println("3️⃣  Reading from embed.FS (Directory):")

	// Read a specific file from the embedded filesystem
	content, err := staticFiles.ReadFile("public/data.txt")
	if err != nil {
		log.Fatal("Failed reading from embedded FS:", err)
	}

	fmt.Printf("   public/data.txt content:\n   %s\n", string(content))
	fmt.Println()

	// List directory contents using ReadDir
	fmt.Println("   Directory contents of 'public':")
	entries, err := staticFiles.ReadDir("public")
	if err != nil {
		log.Fatal(err)
	}

	for _, entry := range entries {
		fmt.Printf("    - %s (IsDir: %t)\n", entry.Name(), entry.IsDir())
	}

	fmt.Println()
	fmt.Println("KEY TAKEAWAY:")
	fmt.Println("  - //go:embed compiles files directly into your Go executable.")
	fmt.Println("  - The `embed` package must be imported to use the directive.")
	fmt.Println("  - Use `string` or `[]byte` for single files.")
	fmt.Println("  - Use `embed.FS` for directories (great for web server assets).")
	fmt.Println("  - Solves the 'works on my machine but missing files in prod' problem.")
}
