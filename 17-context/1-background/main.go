// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

package main

import (
	"context"
	"fmt"
)

// ============================================================================
// Section 17: Context — Background & TODO
// Level: Beginner
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - What context.Context IS and why every Go program uses it
//   - context.Background() — the root of all context trees
//   - context.TODO() — placeholder when you don't know which context to use
//   - How context forms a TREE structure (parent → child → grandchild)
//
// ENGINEERING DEPTH:
//   At its core, `context.Context` is just an interface containing an empty data
//   type `type emptyCtx int`. Calling `context.Background()` simply returns
//   a hardcoded pointer `new(emptyCtx)`. Because all derived Contexts (like timeouts)
//   wrap their parent in a linked list, this empty struct acts as the immutable,
//   zero-byte invisible root anchor for the entire concurrency tree.
//
// RUN: go run ./17-context/1-background
// ============================================================================

func main() {
	fmt.Println("=== Context: Background & TODO ===")
	fmt.Println()

	// --- context.Background() ---
	// This creates the ROOT context — it's the starting point for all context trees.
	// It is NEVER cancelled, has no deadline, and carries no values.
	//
	// Use Background() in:
	//   - main() function
	//   - Initialization code
	//   - Top-level test functions
	//   - As the parent for all derived contexts
	ctx := context.Background()
	fmt.Printf("Background context: %v\n", ctx)

	// Context is an INTERFACE with 4 methods:
	//   Deadline() (deadline time.Time, ok bool)  — when this context expires
	//   Done() <-chan struct{}                      — closed when context is cancelled
	//   Err() error                                — why the context was cancelled
	//   Value(key any) any                         — request-scoped values

	// Background has no deadline
	deadline, ok := ctx.Deadline()
	fmt.Printf("Has deadline: %v (deadline: %v)\n", ok, deadline)

	// Background is never done (Done() channel is nil)
	fmt.Printf("Done channel: %v (nil = never cancels)\n", ctx.Done())

	// Background has no error
	fmt.Printf("Error: %v\n", ctx.Err())
	fmt.Println()

	// --- context.TODO() ---
	// TODO() is identical to Background() in behavior — it never cancels.
	// The difference is SEMANTIC: it signals to other developers that
	// "I know I need a context here but haven't decided which one yet."
	//
	// Use TODO() when:
	//   - You're prototyping and will add proper context later
	//   - A function requires context but you haven't wired it through yet
	//   - You're refactoring and will replace it in the next commit
	//
	// In production code, every TODO() should eventually be replaced
	// with a proper derived context.
	todoCtx := context.TODO()
	fmt.Printf("TODO context: %v\n", todoCtx)
	fmt.Println()

	// --- THE CONTEXT TREE ---
	// Contexts form a tree:
	//
	//   Background()  ← Root (never cancels)
	//       │
	//       ├── WithTimeout(parent, 5s)  ← Cancels after 5 seconds
	//       │       │
	//       │       └── WithValue(parent, "userID", 42)  ← Carries data
	//       │
	//       └── WithCancel(parent)  ← Cancels when you call cancel()
	//
	// When a PARENT context is cancelled, ALL CHILDREN are cancelled too.
	// This is how cancellation propagates through your entire call stack.

	// Let's demonstrate using context in a function
	processRequest(ctx, "order-123")

	fmt.Println()
	fmt.Println("KEY TAKEAWAYS:")
	fmt.Println("  1. context.Background() is the root — use in main() and init")
	fmt.Println("  2. context.TODO() is a placeholder — replace before shipping")
	fmt.Println("  3. Context is the FIRST parameter: func Do(ctx context.Context, ...)")
	fmt.Println("  4. Cancelling a parent cancels ALL children automatically")
	fmt.Println()
	fmt.Println("   Next: go run ./17-context/2-with-cancel")
}

// processRequest demonstrates the Go convention: context is ALWAYS the first parameter.
// This isn't just a style choice — it's enforced by linters (revive, golangci-lint)
// and is part of Google's Go style guide.
func processRequest(ctx context.Context, orderID string) {
	fmt.Printf("Processing order: %s\n", orderID)
	fmt.Printf("  Context error: %v (nil = still active)\n", ctx.Err())

	// In a real application, you would pass ctx down to:
	//   db.QueryContext(ctx, "SELECT ...")
	//   http.NewRequestWithContext(ctx, "GET", url, nil)
	//   grpcClient.Call(ctx, request)
}
