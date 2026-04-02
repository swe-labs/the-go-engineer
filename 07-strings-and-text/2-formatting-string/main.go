// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

package main

import (
	"errors"
	"fmt"
	"math"
	"strings"
	"time"
)

// ============================================================================
// Section 7: Strings & Text — Formatting with fmt
// Level: Beginner
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - All fmt verbs: %v, %+v, %#v, %T, %d, %f, %s, %q, %t, %p, %b, %x
//   - fmt.Printf vs fmt.Sprintf vs fmt.Fprintf
//   - Width and precision formatting (alignment, decimal places)
//   - fmt.Errorf for wrapping errors with context
//   - Padding and alignment for table output
//
// ANALOGY:
//   Format verbs are like fill-in-the-blank templates:
//     "Hello %s, you are %d years old" ← template with placeholders
//     fmt.Printf(template, "Rasel", 25) → "Hello Rasel, you are 25 years old"
//   Each % placeholder specifies HOW to format that argument.
//
// RUN: go run ./07-strings-and-text/2-formatting-string
// ============================================================================

// ServerConfig demonstrates struct formatting with different verbs.
type ServerConfig struct {
	Host     string
	Port     int
	TLS      bool
	MaxConns int
}

func main() {
	fmt.Println("=== Format Verbs Reference ===")
	fmt.Println()

	// --- GENERAL VERBS (work with any type) ---
	config := ServerConfig{Host: "0.0.0.0", Port: 8080, TLS: true, MaxConns: 100}

	fmt.Println("--- General Verbs ---")
	fmt.Printf("  %%v  (default):     %v\n", config)  // {0.0.0.0 8080 true 100}
	fmt.Printf("  %%+v (with names):  %+v\n", config) // {Host:0.0.0.0 Port:8080 ...}
	fmt.Printf("  %%#v (Go syntax):   %#v\n", config) // main.ServerConfig{Host:"0.0.0.0", ...}
	fmt.Printf("  %%T  (type):        %T\n", config)  // main.ServerConfig
	fmt.Println()

	// --- STRING VERBS ---
	name := "Go Mastery"
	fmt.Println("--- String Verbs ---")
	fmt.Printf("  %%s  (string):      %s\n", name) // Go Mastery
	fmt.Printf("  %%q  (quoted):      %q\n", name) // "Go Mastery"
	fmt.Printf("  %%x  (hex bytes):   %x\n", name) // 476f204d6173746572...
	fmt.Println()

	// --- INTEGER VERBS ---
	port := 8080
	fmt.Println("--- Integer Verbs ---")
	fmt.Printf("  %%d  (decimal):     %d\n", port) // 8080
	fmt.Printf("  %%b  (binary):      %b\n", port) // 1111110010000
	fmt.Printf("  %%o  (octal):       %o\n", port) // 17620
	fmt.Printf("  %%x  (hex lower):   %x\n", port) // 1f90
	fmt.Printf("  %%X  (hex upper):   %X\n", port) // 1F90
	fmt.Printf("  %%c  (character):   %c\n", 65)   // A (ASCII 65)
	fmt.Println()

	// --- FLOAT VERBS ---
	pi := math.Pi
	fmt.Println("--- Float Verbs ---")
	fmt.Printf("  %%f  (decimal):     %f\n", pi)   // 3.141593
	fmt.Printf("  %%.2f (2 decimals): %.2f\n", pi) // 3.14
	fmt.Printf("  %%e  (scientific):  %e\n", pi)   // 3.141593e+00
	fmt.Printf("  %%g  (compact):     %g\n", pi)   // 3.141592653589793
	fmt.Println()

	// --- BOOLEAN ---
	fmt.Println("--- Boolean Verb ---")
	fmt.Printf("  %%t  (bool):        %t\n", true) // true
	fmt.Println()

	// --- POINTER ---
	x := 42
	fmt.Println("--- Pointer Verb ---")
	fmt.Printf("  %%p  (pointer):     %p\n", &x) // 0xc000014098
	fmt.Println()

	// --- WIDTH & ALIGNMENT ---
	fmt.Println("--- Width & Alignment ---")
	// %Nd = right-aligned in N characters
	// %-Ns = left-aligned in N characters
	fmt.Printf("  Right-align: [%10d]\n", 42)       // [        42]
	fmt.Printf("  Left-align:  [%-10d]\n", 42)      // [42        ]
	fmt.Printf("  Zero-pad:    [%010d]\n", 42)      // [0000000042]
	fmt.Printf("  Right str:   [%20s]\n", "hello")  // [               hello]
	fmt.Printf("  Left str:    [%-20s]\n", "hello") // [hello               ]
	fmt.Println()

	// --- FORMATTED TABLE OUTPUT ---
	// Real-world use: formatting command-line output as aligned tables
	fmt.Println("--- Table Output (practical use) ---")
	services := []struct {
		name   string
		port   int
		status string
	}{
		{"api-gateway", 8080, "running"},
		{"auth-service", 9090, "running"},
		{"db-postgres", 5432, "running"},
		{"redis-cache", 6379, "stopped"},
	}

	fmt.Printf("  %-15s %-8s %s\n", "SERVICE", "PORT", "STATUS")
	fmt.Printf("  %s\n", strings.Repeat("─", 35))
	for _, svc := range services {
		fmt.Printf("  %-15s %-8d %s\n", svc.name, svc.port, svc.status)
	}
	fmt.Println()

	// --- fmt.Sprintf vs fmt.Printf ---
	// Printf writes to stdout (screen)
	// Sprintf RETURNS a string (for variables, logging, error messages)
	// Fprintf writes to any io.Writer (files, HTTP responses, etc.)
	timestamp := time.Now().Format("15:04:05")
	logMsg := fmt.Sprintf("[%s] Server started on port %d", timestamp, port)
	fmt.Printf("  Sprintf result: %s\n", logMsg)
	fmt.Println()

	// --- fmt.Errorf — Error wrapping ---
	// The %w verb wraps an existing error with additional context.
	// errors.Is() and errors.As() can traverse the wrapped chain.
	originalErr := errors.New("connection refused")
	wrappedErr := fmt.Errorf("failed to connect to database: %w", originalErr)
	fmt.Printf("  Wrapped error: %v\n", wrappedErr)
	fmt.Printf("  Is original?:  %t\n", errors.Is(wrappedErr, originalErr))

	fmt.Println()
	fmt.Println("KEY TAKEAWAY:")
	fmt.Println("  - General: v, +v, #v, T verbs work on any type")
	fmt.Println("  - Types: s (string), d (int), f (float), t (bool), p (pointer)")
	fmt.Println("  - Width: right-align, left-align (with -), zero-pad (with 0)")
	fmt.Println("  - fmt.Sprintf returns a string, fmt.Printf prints to stdout")
	fmt.Println("  - fmt.Errorf with the w verb wraps errors for chain traversal")
}
