// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 06: Backend, APIs & Databases
// Title: gRPC Interceptors
// Level: Core
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - How to use Interceptors to add middleware logic to gRPC calls.
//   - The difference between Unary and Stream interceptors.
//   - How to implement cross-cutting concerns like logging and auth.
//
// WHY THIS MATTERS:
//   - Just like HTTP middleware, gRPC interceptors allow you to keep your
//     business logic clean by moving infrastructure concerns to a shared layer.
//
// RUN:
//   go run ./06-backend-db/01-web-and-database/apis/7-grpc-interceptors
//
// KEY TAKEAWAY:
//   - Interceptors are the "Middleware" of the gRPC world.
// ============================================================================

package main

import "fmt"

// Stage 06: APIs - gRPC Interceptors
//
//   - Unary Interceptor: Wraps a single request-response call
//   - Stream Interceptor: Wraps the initialization of a stream
//   - Middleware chain: Logging -> Auth -> Metrics
//
// ENGINEERING DEPTH:
//   A gRPC Unary Interceptor is a function that sits between the network
//   and your handler. It has access to the `FullMethod` name (e.g.,
//   "/user.UserService/GetUser"), the request object, and the context.
//   This makes it the perfect place to inject trace IDs or verify
//   JWT tokens before your business logic even sees the request.

func main() {
	fmt.Println("=== gRPC Interceptors (Conceptual) ===")
	fmt.Println()

	fmt.Println("  1. UNARY INTERCEPTOR SIGNATURE")
	fmt.Println("     func(ctx, req, info, handler) (res, error)")
	fmt.Println()

	fmt.Println("  2. EXAMPLE: LOGGING INTERCEPTOR")
	fmt.Println("     - Record the start time.")
	fmt.Println("     - Log the method: info.FullMethod.")
	fmt.Println("     - Call the actual handler.")
	fmt.Println("     - Log the result and duration.")
	fmt.Println()

	fmt.Println("  3. REGISTRATION")
	fmt.Println("     s := grpc.NewServer(")
	fmt.Println("       grpc.UnaryInterceptor(myLoggingInterceptor),")
	fmt.Println("     )")
	fmt.Println()

	fmt.Println("  4. COMMON USES")
	fmt.Println("     - Authentication (JWT validation).")
	fmt.Println("     - Observability (Tracing and Metrics).")
	fmt.Println("     - Error Mapping (Converting internal errors to gRPC codes).")
	fmt.Println("     - Rate Limiting.")

	fmt.Println("\n---------------------------------------------------")
	fmt.Println("NEXT UP: API.8 rest-vs-grpc")
	fmt.Println("Current: API.7 (grpc-interceptors)")
	fmt.Println("Previous: API.6 (grpc-streaming)")
	fmt.Println("---------------------------------------------------")
}
