// Copyright (c) 2026 Rasel Hossen

// ============================================================================
// Section 04: Types and Design
// Title: Functional Options
// Level: Core
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Implementing the functional options pattern for flexible constructors.
//   - Managing default states in complex data structures.
//   - Utilizing variadic functions and closures for configuration logic.
//   - Building extensible and type-safe APIs for library development.
//
// WHY THIS MATTERS:
//   - When a type possesses numerous optional configuration fields,
//     traditional constructors become unwieldy and fragile. The
//     functional options pattern provides a declarative mechanism
//     for configuring entities, allowing callers to specify only the
//     attributes they care about while maintaining sensible
//     defaults and high readability.
//
// RUN:
//   go run ./04-types-design/12-functional-options
//
// KEY TAKEAWAY:
//   - Functional options decoupling configuration logic from entity instantiation.
// ============================================================================

// See LICENSE for usage terms.

package main

import "fmt"

// Section 04: Types & Design - Functional Options

// Server represents a configurable compute resource.
type Server struct {
	Name     string
	Region   string
	CPUs     int
	RAMGB    int
	Firewall bool
	SSL      bool
}

// Option defines the signature for functions that configure a Server.
type Option func(*Server)

// WithName configures the identifier of the server.
func WithName(name string) Option {
	return func(s *Server) {
		s.Name = name
	}
}

// WithRegion configures the geographical location of the server.
func WithRegion(region string) Option {
	return func(s *Server) {
		s.Region = region
	}
}

// WithCPUs configures the number of CPU cores allocated to the server.
func WithCPUs(cpus int) Option {
	return func(s *Server) {
		s.CPUs = cpus
	}
}

// WithRAM configures the amount of memory in Gigabytes.
func WithRAM(gb int) Option {
	return func(s *Server) {
		s.RAMGB = gb
	}
}

// WithFirewall enables or disables the network firewall.
func WithFirewall(enabled bool) Option {
	return func(s *Server) {
		s.Firewall = enabled
	}
}

// WithSSL enables or disables SSL encryption for the server.
func WithSSL(enabled bool) Option {
	return func(s *Server) {
		s.SSL = enabled
	}
}

// NewServer initializes a Server with default settings and applies the provided options.
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
	fmt.Println("=== Functional Options: Configurable APIs ===")
	fmt.Println()

	// 1. Sensible Defaults.
	// Calling the constructor without arguments initializes the target
	// with a predefined production-safe state.
	fmt.Println("--- Default Configuration ---")
	s1 := NewServer()
	fmt.Printf("  Defaults: Name=%s, Region=%s, CPUs=%d\n", s1.Name, s1.Region, s1.CPUs)
	fmt.Println()

	// 2. Targeted Customization.
	// The caller only specifies the fields requiring modification.
	// All other fields remain at their default values.
	fmt.Println("--- Partial Configuration ---")
	s2 := NewServer(WithName("web-tier"))
	fmt.Printf("  Custom:   Name=%s, CPUs=%d (Defaults applied)\n", s2.Name, s2.CPUs)
	fmt.Println()

	// 3. Declarative Composition.
	// Options can be combined and reused to create standardized configurations.
	fmt.Println("--- Full Configuration ---")
	s3 := NewServer(
		WithName("api-tier"),
		WithRegion("us-west"),
		WithCPUs(8),
		WithRAM(32),
		WithSSL(true),
	)
	fmt.Printf("  Advanced: %+v\n", s3)

	fmt.Println()
	fmt.Println("---------------------------------------------------")
	fmt.Println("NEXT UP: TI.13 -> 04-types-design/13-method-values")
	fmt.Println("Run    : go run ./04-types-design/13-method-values")
	fmt.Println("Current: TI.12 (functional-options)")
	fmt.Println("---------------------------------------------------")
}
