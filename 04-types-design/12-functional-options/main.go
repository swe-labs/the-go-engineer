// Copyright (c) 2026 Rasel Hossen

// ============================================================================
// Section 04: Types and Design
// Title: Functional Options
// Level: Core
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Learn the functional options pattern-a common Go pattern for building configurable APIs without requiring many constructor parameters.
//
// WHY THIS MATTERS:
//   - Think of ordering a pizza. You could have a constructor with 20 parameters (crust, sauce, cheese, toppings, size, etc.). Or you could have `WithExt...
//
// RUN:
//   go run ./04-types-design/12-functional-options
//
// KEY TAKEAWAY:
//   - Learn the functional options pattern-a common Go pattern for building configurable APIs without requiring many constructor parameters.
// ============================================================================

// See LICENSE for usage terms.

package main

import "fmt"

//
//   - The functional options pattern
//   - Defining Option type and functions
//   - Building clean, extensible APIs
//

type Server struct {
	Name     string
	Region   string
	CPUs     int
	RAMGB    int
	Firewall bool
	SSL      bool
}

type Option func(*Server)

func WithName(name string) Option {
	return func(s *Server) {
		s.Name = name
	}
}

func WithRegion(region string) Option {
	return func(s *Server) {
		s.Region = region
	}
}

func WithCPUs(cpus int) Option {
	return func(s *Server) {
		s.CPUs = cpus
	}
}

func WithRAM(gb int) Option {
	return func(s *Server) {
		s.RAMGB = gb
	}
}

func WithFirewall(enabled bool) Option {
	return func(s *Server) {
		s.Firewall = enabled
	}
}

func WithSSL(enabled bool) Option {
	return func(s *Server) {
		s.SSL = enabled
	}
}

func NewServer(opts ...Option) *Server {
	s := &Server{
		Name:     "default",
		Region:   "us-east-1",
		CPUs:     2,
		RAMGB:    4,
		Firewall: true,
		SSL:      false,
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func main() {
	fmt.Println("=== Functional Options Pattern ===")
	fmt.Println()

	fmt.Println("--- Server with default values ---")
	s1 := NewServer()
	fmt.Printf("  %+v\n", s1)

	fmt.Println()
	fmt.Println("--- Server with custom name only ---")
	s2 := NewServer(WithName("web-server"))
	fmt.Printf("  Name: %s, CPUs: %d, RAM: %dGB\n", s2.Name, s2.CPUs, s2.RAMGB)

	fmt.Println()
	fmt.Println("--- Server with multiple options ---")
	s3 := NewServer(
		WithName("api-server"),
		WithRegion("eu-west"),
		WithCPUs(8),
		WithRAM(32),
		WithSSL(true),
	)
	fmt.Printf("  %+v\n", s3)

	fmt.Println()
	fmt.Println("--- Reusing option combinations ---")
	withProdSettings := []Option{
		WithCPUs(16),
		WithRAM(64),
		WithFirewall(true),
		WithSSL(true),
	}
	withName := WithName("prod-db")
	allOpts := append([]Option{withName}, withProdSettings...)
	s4 := NewServer(allOpts...)
	fmt.Printf("  %+v\n", s4)

	fmt.Println()
	fmt.Println("KEY TAKEAWAY:")
	fmt.Println("  - Functional options allow flexible configuration")
	fmt.Println("  - Callers specify only what they need")
	fmt.Println("  - Options are composable and reusable")
	fmt.Println("  - Standard Go pattern for extensible APIs")
	fmt.Println("\n---------------------------------------------------")
	fmt.Println("NEXT UP: TI.13 method-values")
	fmt.Println("Current: TI.12 (functional-options)")
	fmt.Println("Previous: TI.11 (dynamic-typing-with-any)")
	fmt.Println("---------------------------------------------------")
}
