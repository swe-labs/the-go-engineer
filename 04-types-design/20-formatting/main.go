// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 04: Types and Design
// Title: Formatting
// Level: Core
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Utilizing 'fmt' verbs for type-specific and general value formatting.
//   - Controlling width, precision, and alignment for tabular output.
//   - Implementing error wrapping using the %w verb for context propagation.
//   - Understanding the differences between Printf, Sprintf, and Fprintf.
//
// WHY THIS MATTERS:
//   - Professional Go applications require consistent and readable output
//     for logs, CLI tools, and debugging. The 'fmt' package is the
//     standard mechanism for transforming internal data structures
//     into human-readable text. Mastering its template-driven verbs
//     allows you to build sophisticated logging systems and user
//     interfaces while maintaining low overhead and high precision.
//
// RUN:
//   go run ./04-types-design/20-formatting
//
// KEY TAKEAWAY:
//   - The 'fmt' package enables precise, template-driven text generation.
// ============================================================================

// Commercial use is prohibited without permission.

package main

import (
	"errors"
	"fmt"
	"math"
	"strings"
)

//   - Reflection-based formatting verbs (%v, %+v, %#v, %T).
//   - Numeric and string specialization verbs (%d, %f, %q, %x).
//   - Width, precision, and alignment controls for structured output.
//   - Dynamic error wrapping using the %w verb for context propagation.
//
// TECHNICAL RATIONALE:
//   The 'fmt' package utilizes reflection to inspect the type and
//   structure of data at runtime. While this introduces a small
//   performance overhead compared to manual string concatenation,
//   it provides a powerful, template-driven mechanism for
//   generating human-readable state. For high-performance logging
//   or error handling, understanding the cost of these operations
//   and the mechanics of the %w wrapping verb is essential for
//   maintaining system transparency and debuggability.
//

// ServerConfig (Struct): aggregates configuration state for demonstration of reflection-based formatting verbs.
type ServerConfig struct {
	Host     string
	Port     int
	TLS      bool
	MaxConns int
}

func main() {
	fmt.Println("=== Formatting: Template-Driven Output ===")
	fmt.Println()

	// 1. General Reflection Verbs.
	// %v is the default format. %+v adds field names. %#v adds type information.
	fmt.Println("--- Reflection & Debugging ---")
	config := ServerConfig{Host: "0.0.0.0", Port: 8080, TLS: true, MaxConns: 100}
	fmt.Printf("  Default (%%v):  %v\n", config)
	fmt.Printf("  Struct (%%+v): %+v\n", config)
	fmt.Printf("  Syntax (%%#v): %#v\n", config)
	fmt.Printf("  Type   (%%T):  %T\n", config)
	fmt.Println()

	// 2. Type-Specific Verbs.
	// Verbs allow precise control over numeric and string representation.
	fmt.Println("--- Type Specialization ---")
	pi := math.Pi
	fmt.Printf("  Decimal:   %d (int)\n", 8080)
	fmt.Printf("  Precision: %.2f (float)\n", pi)
	fmt.Printf("  Quoted:    %q (string)\n", "the-go-engineer")
	fmt.Printf("  Pointer:   %p\n", &config)
	fmt.Println()

	// 3. Alignment & Padding.
	// Width specifications enable the construction of aligned tables.
	fmt.Println("--- Tabular Alignment ---")
	services := []struct {
		name string
		port int
	}{
		{"api-gateway", 8080},
		{"auth-service", 9090},
	}
	fmt.Printf("  %-15s | %s\n", "SERVICE", "PORT")
	fmt.Printf("  %s\n", strings.Repeat("-", 25))
	for _, s := range services {
		fmt.Printf("  %-15s | %d\n", s.name, s.port)
	}
	fmt.Println()

	// 4. Error Wrapping.
	// %w allows errors to be wrapped with context while preserving
	// the ability to inspect the original error via errors.Is.
	fmt.Println("--- Error Context Propagation ---")
	baseErr := errors.New("network timeout")
	wrappedErr := fmt.Errorf("service request failed: %w", baseErr)
	fmt.Printf("  Wrapped: %v\n", wrappedErr)
	fmt.Printf("  Unwrap:  %t (errors.Is)\n", errors.Is(wrappedErr, baseErr))

	fmt.Println()
	fmt.Println("---------------------------------------------------")
	fmt.Println("NEXT UP: ST.3 -> 04-types-design/21-unicode")
	fmt.Println("Run    : go run ./04-types-design/21-unicode")
	fmt.Println("Current: ST.2 (formatting)")
	fmt.Println("---------------------------------------------------")
}
