// Copyright (c) 2026 Rasel Hossen

// ============================================================================
// Section 04: Types and Design
// Title: Interface Embedding
// Level: Core
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Composing complex behavioral contracts from simpler interfaces.
//   - The mechanics of method set aggregation during embedding.
//   - Analyzing standard library patterns like `io.ReadWriter`.
//
// WHY THIS MATTERS:
//   - Go promotes the use of small, focused interfaces. Embedding
//     allows you to combine these atomic behaviors into sophisticated
//     contracts without duplicating method signatures, adhering to the
//     Interface Segregation Principle.
//
// RUN:
//   go run ./04-types-design/4-interface-embedding
//
// KEY TAKEAWAY:
//   - Interface embedding is the primary tool for behavioral composition.
// ============================================================================

// See LICENSE for usage terms.

package main

import (
	"bytes"
	"fmt"
)

// Section 04: Types & Design - Interface Embedding
//   - Embedding one interface inside another
//   - How embedded interfaces combine contracts
//   - The io.ReadWriter pattern from the standard library
//

// Reader defines a basic data retrieval contract.
type Reader interface {
	Read(p []byte) (n int, err error)
}

// Writer defines a basic data persistence contract.
type Writer interface {
	Write(p []byte) (n int, err error)
}

// ReadWriter embeds both Reader and Writer to create a composite behavioral contract.
type ReadWriter interface {
	Reader
	Writer
}

type Buffer struct {
	buf bytes.Buffer
}

func (b *Buffer) Read(p []byte) (n int, err error) {
	return b.buf.Read(p)
}

func (b *Buffer) Write(p []byte) (n int, err error) {
	return b.buf.Write(p)
}

func processReadWrite(rw ReadWriter) {
	buf := make([]byte, 10)
	rw.Write([]byte("Hello"))
	n, _ := rw.Read(buf)
	fmt.Printf("Read %d bytes: %s\n", n, string(buf[:n]))
}

func main() {
	fmt.Println("=== Interface Embedding: Behavioral Composition ===")
	fmt.Println()

	// 1. Implementation of embedded behaviors.
	// Buffer satisfies Reader and Writer individually.
	fmt.Println("--- Satisfying Embedded Contracts ---")
	b := &Buffer{}
	processReadWrite(b)

	// 2. Composed Interfaces (Standard Library Pattern).
	// Types like bytes.Buffer satisfy io.ReadWriter by implementing the individual
	// methods of io.Reader and io.Writer.
	fmt.Println()
	fmt.Println("--- io.ReadWriter Composition Pattern ---")
	rb := &bytes.Buffer{}
	rb.WriteString("Composed behavior")
	data := make([]byte, 8)
	n, _ := rb.Read(data)
	fmt.Printf("  Read %d bytes: %s\n", n, string(data))

	fmt.Println()
	fmt.Println("---------------------------------------------------")
	fmt.Println("NEXT UP: TI.6 -> 04-types-design/6-type-switch")
	fmt.Println("Run    : go run ./04-types-design/6-type-switch")
	fmt.Println("Current: TI.4 (interface-embedding)")
	fmt.Println("---------------------------------------------------")
}
