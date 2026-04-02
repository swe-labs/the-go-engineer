// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

package main

import (
	"context"
	"fmt"
)

// ============================================================================
// Section 17: Context — WithValue
// Level: Advanced
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - context.WithValue stores request-scoped data in the context
//   - Using custom key types to prevent collisions
//   - When to use (and NOT use) context values
//   - The request ID / user ID pattern used in production
//
// IMPORTANT RULE:
//   Context values are for REQUEST-SCOPED data only.
//   DO NOT use context to pass function parameters.
//   Bad:  ctx = context.WithValue(ctx, "db", database)     ← Anti-pattern
//   Good: ctx = context.WithValue(ctx, requestIDKey, "abc") ← Request metadata
//
// ENGINEERING DEPTH:
//   `WithValue` creates a `valueCtx` struct containing exactly *one* key-value pair
//   and a pointer to its parent. Because Contexts are immutable, adding 5 values
//   creates a chain of 5 nested Context wrappers. When you call `ctx.Value(key)`,
//   Go walks UP the parent chain comparing keys until it finds a match or reaches
//   the root. This is O(depth) — proportional to how many values were added. This
//   is why you must NEVER store large amounts of data in Context; the linear lookup
//   cost at high concurrency will degrade your server's performance.
//
// RUN: go run ./17-context/4-with-value
// ============================================================================

// --- CUSTOM KEY TYPES ---
// Always use unexported custom types for context keys.
// This prevents collisions between packages that might use the same string key.
//
// If package A uses context.WithValue(ctx, "userID", 1)
// and package B uses context.WithValue(ctx, "userID", "admin")
// they would overwrite each other! Custom types prevent this.
type contextKey string

const (
	requestIDKey contextKey = "requestID"
	userIDKey    contextKey = "userID"
)

func main() {
	fmt.Println("=== Context: WithValue ===")
	fmt.Println()

	// Start with a background context
	ctx := context.Background()

	// Add request-scoped values
	// Each WithValue call wraps the parent context with a new value.
	// The original context is NOT modified (contexts are immutable).
	ctx = context.WithValue(ctx, requestIDKey, "req-abc-123")
	ctx = context.WithValue(ctx, userIDKey, 42)

	// Pass the enriched context to handlers
	handleRequest(ctx)

	fmt.Println()

	// --- VALUES ARE INHERITED BY CHILDREN ---
	// Any derived context (WithCancel, WithTimeout) inherits parent values.
	childCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	fmt.Println("=== Child Context Inherits Values ===")
	reqID := childCtx.Value(requestIDKey)
	fmt.Printf("  Child sees requestID: %v\n", reqID) // Inherited from parent!

	fmt.Println()
	fmt.Println("KEY TAKEAWAYS:")
	fmt.Println("  1. Use custom key types (not strings) to prevent collisions")
	fmt.Println("  2. Context values are for REQUEST-SCOPED metadata only")
	fmt.Println("  3. Good uses: request ID, user ID, trace ID, auth token")
	fmt.Println("  4. Bad uses: database connections, loggers, config (use DI instead)")
	fmt.Println("  5. Values are inherited by child contexts automatically")
}

// handleRequest demonstrates extracting values from context.
// In a real HTTP server, the framework populates these values in middleware
// and handlers downstream read them.
func handleRequest(ctx context.Context) {
	// Extract values with ctx.Value(key).
	// Value() returns `any` — you must type-assert to the expected type.
	// Always handle the case where the value is nil (key not found).
	requestID, ok := ctx.Value(requestIDKey).(string)
	if !ok {
		requestID = "unknown"
	}

	userID, ok := ctx.Value(userIDKey).(int)
	if !ok {
		userID = -1
	}

	fmt.Printf("  Handling request: ID=%s, UserID=%d\n", requestID, userID)

	// Pass context further downstream — the values travel with it
	logAction(ctx, "processing order")
}

// logAction shows context values flowing through the entire call chain.
// Every function in the chain has access to request-scoped data
// without needing it as an explicit parameter.
func logAction(ctx context.Context, action string) {
	requestID, _ := ctx.Value(requestIDKey).(string)
	fmt.Printf("  [%s] Action: %s\n", requestID, action)
}
