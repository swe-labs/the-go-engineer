// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

package user

import (
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// ============================================================================
// Stage 08: Quality and Performance - Testable Design (io.Writer Injection)
// Level: Intermediate
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Why hardcoding stdout/stderr makes code untestable
//   - Refactoring to use `io.Writer` for dependency injection
//   - Testing standard output by overriding os.Stdout (the hard way)
//
// ENGINEERING DEPTH:
//   "Testable Design" means structuring your code so it can be easily mocked
//   or observed. If a function prints directly to `os.Stdout` (via `fmt.Printf`),
//   it is extremely painful to assert what was printed.
//
//   Instead of `fmt.Println(text)`, accept an `io.Writer` interface parameter.
//   In production, pass `os.Stdout`. In tests, pass a `strings.Builder`.
// ============================================================================

// GreetingHardcoded prints directly to the terminal.
func GreetingHardcoded(prefix, name string) {
	fmt.Printf("Hello, %s %s!\n", prefix, name)
}

// TestGreetingHardcoded is a painful test.
// We must hijack the entire operating system's standard output stream to capture
// what prints. This is highly not recommended, not concurrent-safe, and brittle.
func TestGreetingHardcoded(t *testing.T) {
	originalStdout := os.Stdout

	r, w, err := os.Pipe()
	assert.NoError(t, err)

	os.Stdout = w
	defer func() {
		os.Stdout = originalStdout
	}()

	GreetingHardcoded("Mr.", "Joseph")
	w.Close()

	var buf strings.Builder
	_, err = io.Copy(&buf, r)
	assert.NoError(t, err)

	assert.Equal(t, "Hello, Mr. Joseph!\n", buf.String())
}

// GreetingTestable accepts any io.Writer (a file, a network connection, stdout, or a buffer).
func GreetingTestable(out io.Writer, prefix, name string) {
	fmt.Fprintf(out, "Hello, %s %s!\n", prefix, name)
}

// TestGreetingTestable is incredibly clean and fast.
func TestGreetingTestable(t *testing.T) {
	var buf strings.Builder

	GreetingTestable(&buf, "Ms.", "Alice")

	assert.Equal(t, "Hello, Ms. Alice!\n", buf.String())
}
