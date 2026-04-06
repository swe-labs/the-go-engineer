// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

// RUN: go run ./04-functions-and-errors/9-functional-options
package main

import (
	"fmt"
	"time"
)

// ============================================================================
// Section 04: Functions & Errors — Functional Options Pattern
// Level: Intermediate
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - The Functional Options design pattern in Go (popularized by Dave Cheney).
//   - How to avoid constructor functions with dozens of parameters.
//   - How to safely provide default configurations with optional overrides.
//
// ENGINEERING DEPTH:
//   When building APIs or structs with many configurable fields, you generally
//   have three choices in Go:
//   1. Pass every parameter to the constructor: `NewServer("127.0.0.1", 8080, 10, true, ...)` // Hard to read, breaks backward compatibility when adding fields.
//   2. Pass a Config struct: `NewServer(ServerConfig{Port: 8080})` // Better, but distinguishing zero-values (e.g., 0 port vs default port) requires pointers in the struct.
//   3. Functional Options: `NewServer(WithPort(8080), WithTimeout(time.Second))` // Best! Highly readable, backward compatible, and safely handles defaults.
// ============================================================================

// --- 1. The Core Struct ---
// Server is a simulated HTTP server that requires configuration.
type Server struct {
	host         string
	port         int
	timeout      time.Duration
	maxBodyBytes int64
	tlsEnabled   bool
}

// --- 2. The Option Type ---
// Option is a function that takes a pointer to a Server and modifies it.
type Option func(*Server)

// --- 3. The Functional Option Builders ---

// WithHost configures the server hostname
func WithHost(host string) Option {
	// We return an anonymous function matching the Option signature.
	return func(s *Server) {
		s.host = host
	}
}

// WithPort configures the server port
func WithPort(port int) Option {
	return func(s *Server) {
		s.port = port
	}
}

// WithTimeout configures the connection timeout
func WithTimeout(d time.Duration) Option {
	return func(s *Server) {
		s.timeout = d
	}
}

// WithTLS enables or disables TLS
func WithTLS(enabled bool) Option {
	return func(s *Server) {
		s.tlsEnabled = enabled
	}
}

// --- 4. The Constructor ---

// NewServer creates a Server. It takes variadic Options.
// This means you can pass 0, 1, or 100 options cleanly.
func NewServer(opts ...Option) *Server {
	// Step A: Set sensible defaults
	server := &Server{
		host:         "localhost",
		port:         8080,
		timeout:      30 * time.Second,
		maxBodyBytes: 1048576, // 1MB
		tlsEnabled:   false,
	}

	// Step B: Loop over every option provided by the caller
	for _, opt := range opts {
		// Step C: Execute the option function, mutating our Server pointer
		opt(server)
	}

	return server
}

func main() {
	// Scenario A: Start server with defaults
	// We pass no options. It uses localhost:8080.
	defaultServer := NewServer()
	fmt.Printf("Default Server:  %+v\n", defaultServer)

	// Scenario B: Start server with specific overrides
	// Very readable. Order does not matter!
	customServer := NewServer(
		WithPort(443),
		WithTLS(true),
		WithTimeout(5*time.Second),
	)
	fmt.Printf("Custom Server:   %+v\n", customServer)

	// Scenario C: Dynamic configuration driven by config/env
	var myOpts []Option
	myOpts = append(myOpts, WithHost("0.0.0.0"))

	// Imagine we read an env variable here:
	useSecureMode := true
	if useSecureMode {
		myOpts = append(myOpts, WithTLS(true))
		myOpts = append(myOpts, WithPort(8443))
	}

	envServer := NewServer(myOpts...) // Unpack the slice into the variadic parameter
	fmt.Printf("Env-based Server: %+v\n", envServer)
}
