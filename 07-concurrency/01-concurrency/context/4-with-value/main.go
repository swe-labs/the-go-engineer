// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

package main

import (
	"context"
	"fmt"
)

// ============================================================================
// Stage 07: Context - WithValue
// Level: Advanced
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - context.WithValue stores request-scoped data in the context
//   - Using custom key types to prevent collisions
//   - When to use (and not use) context values
//   - The request ID / user ID pattern used in production
//
// IMPORTANT RULE:
//   Context values are for request-scoped data only.
//   Do not use context to pass normal function parameters.
//   Bad:  ctx = context.WithValue(ctx, "db", database)       // Anti-pattern
//   Good: ctx = context.WithValue(ctx, requestIDKey, "abc")  // Request metadata
//
// ENGINEERING DEPTH:
//   `WithValue` creates a `valueCtx` struct containing exactly one key-value pair
//   and a pointer to its parent. Because Contexts are immutable, adding 5 values
//   creates a chain of 5 nested Context wrappers. When you call `ctx.Value(key)`,
//   Go walks up the parent chain comparing keys until it finds a match or reaches
//   the root. This is O(depth) - proportional to how many values were added.
//
// RUN: go run ./07-concurrency/01-concurrency/context/4-with-value
// ============================================================================

type contextKey string

const (
	requestIDKey contextKey = "requestID"
	userIDKey    contextKey = "userID"
)

func main() {
	fmt.Println("=== Context: WithValue ===")
	fmt.Println()

	ctx := context.Background()
	ctx = context.WithValue(ctx, requestIDKey, "req-abc-123")
	ctx = context.WithValue(ctx, userIDKey, 42)

	handleRequest(ctx)
	fmt.Println()

	childCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	fmt.Println("=== Child Context Inherits Values ===")
	reqID := childCtx.Value(requestIDKey)
	fmt.Printf("  Child sees requestID: %v\n", reqID)

	fmt.Println()
	fmt.Println("KEY TAKEAWAYS:")
	fmt.Println("  1. Use custom key types (not strings) to prevent collisions")
	fmt.Println("  2. Context values are for request-scoped metadata only")
	fmt.Println("  3. Good uses: request ID, user ID, trace ID, auth token")
	fmt.Println("  4. Bad uses: database connections, loggers, config (use DI instead)")
	fmt.Println("  5. Values are inherited by child contexts automatically")
	fmt.Println("\n---------------------------------------------------")
	fmt.Println("NEXT UP: CT.5 timeout-aware API client")
	fmt.Println("   Current: CT.4 (WithValue)")
	fmt.Println("---------------------------------------------------")
}

func handleRequest(ctx context.Context) {
	requestID, ok := ctx.Value(requestIDKey).(string)
	if !ok {
		requestID = "unknown"
	}

	userID, ok := ctx.Value(userIDKey).(int)
	if !ok {
		userID = -1
	}

	fmt.Printf("  Handling request: ID=%s, UserID=%d\n", requestID, userID)
	logAction(ctx, "processing order")
}

func logAction(ctx context.Context, action string) {
	requestID, _ := ctx.Value(requestIDKey).(string)
	fmt.Printf("  [%s] Action: %s\n", requestID, action)
}
