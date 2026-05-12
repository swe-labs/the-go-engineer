// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

package main

import (
	"io"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestGreetingHardcoded demonstrates the painful way — hijacking os.Stdout
// This test is brittle: not concurrent-safe, not portable, and slow.
func TestGreetingHardcoded(t *testing.T) {
	orig := os.Stdout
	r, w, err := os.Pipe()
	assert.NoError(t, err)

	os.Stdout = w
	defer func() { os.Stdout = orig }()

	GreetingHardcoded("Mr.", "Joseph")
	w.Close()

	var buf strings.Builder
	_, err = io.Copy(&buf, r)
	assert.NoError(t, err)
	assert.Equal(t, "Hello, Mr. Joseph!\n", buf.String())
}

// TestGreetingTestable demonstrates the clean way — dependency injection via io.Writer
// This test is fast, concurrent-safe, and readable.
func TestGreetingTestable(t *testing.T) {
	var buf strings.Builder
	GreetingTestable(&buf, "Ms.", "Alice")
	assert.Equal(t, "Hello, Ms. Alice!\n", buf.String())
}
