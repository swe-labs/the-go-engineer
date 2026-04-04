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
// Section 14: Testing — Testable Design (io.Writer injection)
// Level: Intermediate
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Why hardcoding stdout/stderr makes code un-testable
//   - Refactoring to use `io.Writer` for Dependency Injection
//   - Testing standard output by overriding os.Stdout (the hard way)
//
// ENGINEERING DEPTH:
//   "Testable Design" means structuring your code so it can be easily mocked
//   or observed. If a function prints directly to `os.Stdout` (via `fmt.Printf`),
//   it is extremely painful to assert what was printed.
//
//   Instead of `fmt.Println(text)`, accept an `io.Writer` interface parameter!
//   In production, pass `os.Stdout`. In tests, pass a `strings.Builder`.
// ============================================================================

// --- The Bad Way (Hardcoded dependency on os.Stdout) ---

// GreetingHardcoded prints directly to the terminal.
func GreetingHardcoded(prefix, name string) {
	fmt.Printf("Hello, %s %s!\n", prefix, name)
}

// TestGreetingHardcoded is a painful test.
// We must hijack the entire operating system's standard output stream to capture
// what prints. This is highly NOT recommended, not concurrent-safe, and brittle.
func TestGreetingHardcoded(t *testing.T) {
	// Save the real stdout
	originalStdout := os.Stdout

	// Create a pipe (simulates a terminal connection)
	r, w, err := os.Pipe()
	assert.NoError(t, err)

	// Replace the global standard out with our pipe
	os.Stdout = w
	defer func() {
		os.Stdout = originalStdout // Restore it immediately when done!
	}()

	// Execute function
	GreetingHardcoded("Mr.", "Joseph")
	w.Close() // Finished writing

	// Read the output from the pipe
	var buf strings.Builder
	_, err = io.Copy(&buf, r)
	assert.NoError(t, err)

	// Assert
	assert.Equal(t, "Hello, Mr. Joseph!\n", buf.String())
}

// --- The Good Way (Testable Design via Dependency Injection) ---

// GreetingTestable accepts any io.Writer (a file, a network connection, stdout, a buffer).
func GreetingTestable(out io.Writer, prefix, name string) {
	// Fprintf writes strictly to the provided interface
	fmt.Fprintf(out, "Hello, %s %s!\n", prefix, name)
}

// TestGreetingTestable is incredibly clean and fast.
func TestGreetingTestable(t *testing.T) {
	// 1. Create an in-memory buffer (which implements io.Writer)
	var buf strings.Builder

	// 2. Inject it!
	GreetingTestable(&buf, "Ms.", "Alice")

	// 3. Assert
	assert.Equal(t, "Hello, Ms. Alice!\n", buf.String())
}
