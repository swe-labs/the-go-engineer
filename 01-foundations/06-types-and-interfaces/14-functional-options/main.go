// Copyright (c) 2026 Rasel Hossen
// See LICENSE for usage terms.

package main

import "fmt"

// ============================================================================
// Section 6: Types & Interfaces — Functional Options
// Level: Intermediate
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - The functional options pattern
//   - Defining Option type and functions
//   - Building clean, extensible APIs
//
// RUN: go run ./01-foundations/06-types-and-interfaces/14-functional-options
// ============================================================================

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
	fmt.Println("NEXT UP: TI.15 method values")
	fmt.Println("   Current: TI.14 (functional options)")
	fmt.Println("---------------------------------------------------")
}
