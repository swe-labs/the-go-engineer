// Copyright (c) 2026 Rasel Hossen
// See LICENSE for usage terms.

package main

import (
	"bytes"
	"fmt"
)

// ============================================================================
// Section 6: Types & Interfaces — Interface Embedding
// Level: Intermediate
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Embedding one interface inside another
//   - How embedded interfaces combine contracts
//   - The io.ReadWriter pattern from the standard library
//
// RUN: go run ./04-types-design/4-interface-embedding
// ============================================================================

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
	fmt.Println("NEXT UP: TI.5 Stringer")
	fmt.Println("   Current: TI.4 (interface embedding)")
	fmt.Println("---------------------------------------------------")
}
