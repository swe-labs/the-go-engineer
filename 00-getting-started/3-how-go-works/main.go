// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

package main

import (
	"fmt"
	"math"
	"strings"
)

// ============================================================================
// Section 0: Getting Started — How Go Works
// Level: Beginner
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - How Go compiles and runs your code
//   - What "package" and "import" really mean
//   - The difference between compiled and interpreted languages
//   - How Go organizes code with packages from the standard library
//   - How to use multiple packages in a single program
//
// ENGINEERING DEPTH:
//   Go binaries are "statically linked". When you compile your code, the Go compiler
//   actually bundles the necessary C-libraries, runtime orchestrator, and garbage
//   collector straight into your single executable file. This means you can drop
//   a Go binary onto an empty Linux server and it will run perfectly without
//   needing to install Go on that server.
//
// RUN: go run ./00-getting-started/3-how-go-works
// ============================================================================

// --- THE GO COMPILATION MODEL ---
//
// When you run "go run main.go", three things happen:
//
//   Step 1: COMPILE — Go reads your .go files and translates them to
//           machine code (binary). This catches ALL syntax errors.
//           If there's a bug, you find out NOW, not when a user hits it.
//
//   Step 2: LINK — Go bundles your code with any imported packages
//           into a single, self-contained binary file.
//           No runtime dependencies. No "pip install". No "npm install".
//
//   Step 3: EXECUTE — The binary runs directly on your CPU.
//           No interpreter, no virtual machine, no JIT.
//           This is why Go programs start in milliseconds.
//
// "go run" does all 3 steps in memory (doesn't save the binary).
// "go build" does steps 1-2 and saves the binary to disk.
//
// TRY IT:
//   go build -o hello ./00-getting-started/3-how-go-works
//   ./hello                    # Run the binary directly
//   file hello                 # Shows it's a native executable
//   ls -lh hello               # Shows the binary size (~2MB)
//   rm hello                   # Clean up
//
// That ~2MB binary contains EVERYTHING needed to run. You can copy it
// to any machine with the same OS/architecture and it will just work.
// No Go installation required on the target machine.

// --- PACKAGES: GO'S ORGANIZATION SYSTEM ---
//
// A "package" is a directory of .go files that work together.
// Every .go file declares which package it belongs to on line 1.
//
// Go ships with a rich STANDARD LIBRARY of ~150 packages:
//   fmt      — Formatted text output (print, scan)
//   strings  — String manipulation (split, join, replace, trim)
//   math     — Mathematical functions (sqrt, pi, abs, floor)
//   os       — Operating system interaction (files, env vars, exit)
//   net/http — HTTP client and server (web development)
//   time     — Time, duration, and timer functions
//   errors   — Error creation and inspection
//
// These are always available — no download needed.
// Full list: https://pkg.go.dev/std

func main() {
	fmt.Println("=== How Go Works ===")
	fmt.Println()

	// --- USING THE "strings" PACKAGE ---
	// We imported "strings" at the top. Now we can use its functions.
	// Every function in a package is accessed with the dot notation:
	//   packageName.FunctionName()

	greeting := "hello, go developer!"

	// strings.ToUpper converts all characters to uppercase.
	upper := strings.ToUpper(greeting)
	fmt.Printf("Original:    %s\n", greeting)
	fmt.Printf("Uppercased:  %s\n", upper)

	// strings.Contains checks if a string contains a substring.
	// Returns true or false (a "bool" in Go).
	hasGo := strings.Contains(greeting, "go")
	fmt.Printf("Contains 'go': %t\n", hasGo) // %t = boolean format

	// strings.Replace replaces occurrences of a substring.
	// The last argument (-1) means "replace all occurrences".
	replaced := strings.Replace(greeting, "go", "Go", -1)
	fmt.Printf("Replaced:    %s\n", replaced)

	// strings.Split breaks a string into a slice (like an array) at a delimiter.
	words := strings.Split("one,two,three", ",")
	fmt.Printf("Split:       %v\n", words)      // %v = default format for any type
	fmt.Printf("Word count:  %d\n", len(words)) // len() returns the length

	fmt.Println()

	// --- USING THE "math" PACKAGE ---
	// Mathematical constants and functions.

	// math.Pi is a constant: 3.141592653589793
	fmt.Printf("Pi:          %.4f\n", math.Pi) // %.4f = float with 4 decimal places

	// math.Sqrt computes the square root.
	fmt.Printf("√144:        %.0f\n", math.Sqrt(144)) // %.0f = float with 0 decimals

	// math.Pow computes x^y (x raised to the power of y).
	fmt.Printf("2^10:        %.0f\n", math.Pow(2, 10)) // = 1024

	// math.Abs returns the absolute (positive) value.
	fmt.Printf("|-42|:       %.0f\n", math.Abs(-42))

	fmt.Println()

	// --- WHY PACKAGES MATTER ---
	//
	// Packages solve three problems:
	//
	// 1. ORGANIZATION — Code is grouped by purpose (strings, math, net).
	//    You never have a 10,000 line file. Each package is small and focused.
	//
	// 2. REUSABILITY — Write once, import anywhere. The "fmt" package is
	//    used by every Go program ever written. You don't copy-paste print code.
	//
	// 3. ENCAPSULATION — Each package controls what it exports (makes public).
	//    In Go, exported names start with an UPPERCASE letter:
	//      fmt.Println  ← Exported (public). You can use it.
	//      fmt.println  ← Would NOT be exported (private). You can't use it.
	//    This is one of Go's most distinctive rules. No "public" keyword needed.

	fmt.Println("KEY TAKEAWAYS:")
	fmt.Println("  1. Go compiles to a single binary — no runtime dependencies")
	fmt.Println("  2. Packages organize code — import only what you need")
	fmt.Println("  3. Uppercase = exported (public), lowercase = unexported (private)")
	fmt.Println("  4. The standard library is massive — learn it before reaching for external packages")
	fmt.Println()
	fmt.Println("   Next step: go run ./00-getting-started/4-dev-environment")
}
