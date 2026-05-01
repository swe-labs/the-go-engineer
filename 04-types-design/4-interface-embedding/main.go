// Copyright (c) 2026 Rasel Hossen

// ============================================================================
// Section 04: Types and Design
// Title: Interface Embedding
// Level: Core
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Learn how to embed one interface into another to build larger contracts from smaller pieces.
//
// WHY THIS MATTERS:
//   - Think of a universal remote. It does not have buttons for every function directly-it embeds the capabilities of a TV remote, a DVD remote, and a soundbar remote.
//
// RUN:
//   go run ./04-types-design/4-interface-embedding
//
// KEY TAKEAWAY:
//   - Learn how to embed one interface into another to build larger contracts from smaller pieces.
// ============================================================================

// See LICENSE for usage terms.

package main

import (
	"bytes"
	"fmt"
)

//
//   - Embedding one interface inside another
//   - How embedded interfaces combine contracts
//   - The io.ReadWriter pattern from the standard library
//

type Reader interface {
	Read(p []byte) (n int, err error)
}

type Writer interface {
	Write(p []byte) (n int, err error)
}

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
	fmt.Println("=== Interface Embedding ===")
	fmt.Println()

	fmt.Println("--- Embedded Interfaces ---")
	b := &Buffer{}
	processReadWrite(b)

	fmt.Println()
	fmt.Println("--- io.ReadWriter Pattern ---")
	rb := &bytes.Buffer{}
	rb.WriteString("Testing ReadWriter")
	data := make([]byte, 7)
	n, _ := rb.Read(data)
	fmt.Printf("Read %d bytes: %s\n", n, string(data))

	fmt.Println()
	fmt.Println("KEY TAKEAWAY:")
	fmt.Println("  - Interface embedding combines contracts without copying methods")
	fmt.Println("  - io.ReadWriter = io.Reader + io.Writer")
	fmt.Println("  - Embedding is static: compiler verifies all methods exist")
	fmt.Println("\n---------------------------------------------------")
	fmt.Println("NEXT UP: TI.5 -> 04-types-design/5-stringer")
	fmt.Println("Current: TI.4 (interface-embedding)")
	fmt.Println("Previous: TI.3 (interfaces)")
	fmt.Println("---------------------------------------------------")
}
