# Context values with caveats

## Learning objective

Use `context.WithValue` to carry request-scoped data, implement type-safe context keys, and understand when context values are appropriate — and when they are not.

## Why this matters

Context values are Go's mechanism for request-scoped data that spans the entire call graph. The canonical example is a request ID that every function in the handler chain logs. Without context values, you would pass the request ID as an explicit parameter through every function signature — a change that touches hundreds of functions. Context values solve this, but they come with sharp edges: they are immutable in practice, untyped by default, and easily abused as a replacement for proper function parameters.

## Mental model

Think of context as a tree of immutable key-value stores. Each `WithValue` call creates a new node in the tree. When you call `ctx.Value(key)`, Go walks up the tree from the current node to the root, looking for the first node that has the key. This is a linear search — not a hash map. The values are meant to be read-only metadata about the request (request ID, user ID, auth token), not mutable state or function arguments.

## Core idea

`context.WithValue(parent, key, val)` returns a copy of the parent context with `key` mapped to `val`. The key must be comparable (`==` works on it). The value is retrieved with `ctx.Value(key)`.

```go
func WithValue(parent Context, key, val any) Context
func (ctx Context) Value(key any) any
```

Type-safe key pattern — always define a custom unexported key type:

```go
type contextKey string

const keyRequestID contextKey = "request_id"

ctx = context.WithValue(ctx, keyRequestID, "req-abc-123")
reqID, ok := ctx.Value(keyRequestID).(string)
```

## Under the hood

`context.WithValue` creates a `valueCtx` struct with three fields: the parent context, the key, and the value. When `Value(key)` is called, the method checks if the current node's key matches. If yes, it returns the value. If no, it delegates to the parent. This is a singly linked list traversal — O(depth) per lookup.

There is no hashing, no concurrency-safe mutation, and no iteration. The only way to "change" a value is to call `WithValue` again, creating a new node that shadows the previous one for the same key.

## How Go uses it

- **Request IDs**: Middleware extracts or generates a request ID and stores it in the context. Every log line in the handler chain includes it.
- **Authentication**: Auth middleware stores the user ID and claims in the context. Downstream handlers read `ctx.Value(authKey)`.
- **gRPC metadata**: The `google.golang.org/grpc/metadata` package uses `context.WithValue` to attach incoming and outgoing metadata.
- **OpenTelemetry**: The span context is propagated through context values.

## Go example

```go
package main

import (
	"context"
	"fmt"
)

type contextKey string

const (
	keyUserID    contextKey = "user_id"
	keyRequestID contextKey = "request_id"
)

func handler(ctx context.Context) {
	userID, _ := ctx.Value(keyUserID).(string)
	reqID, _ := ctx.Value(keyRequestID).(string)
	fmt.Printf("Handling request %s for user %s\n", reqID, userID)
}

func main() {
	ctx := context.Background()
	ctx = context.WithValue(ctx, keyUserID, "user-42")
	ctx = context.WithValue(ctx, keyRequestID, "req-abc-123")
	handler(ctx)
}
```

## Step-by-step execution

1. `context.Background()` returns an `emptyCtx` — no values, no keys.
2. `context.WithValue(ctx, keyUserID, "user-42")` creates a `valueCtx` node with `key = keyUserID`, `val = "user-42"`, `parent = emptyCtx`.
3. `context.WithValue(ctx, keyRequestID, "req-abc-123")` creates a second `valueCtx` node with `key = keyRequestID`, `parent = first valueCtx`.
4. In `handler`, `ctx.Value(keyUserID)` walks: current node → key is `keyRequestID` (no match) → parent node → key is `keyUserID` (match) → returns `"user-42"`.
5. The type assertion `.(string)` extracts the string. If the key had not been set, `Value` returns `nil` and the assertion yields `("", false)`.
6. The handler prints the request and user IDs.

If a different package used a bare string `"user_id"` as a key, it would collide with this package's key. The custom type `contextKey` prevents this — even if the underlying string is the same, the types are different.

## Common mistakes

- Mistake: Using `string` as the key type directly (e.g., `ctx.Value("user_id")`).
  - Why: Any package can use the same string key, causing collisions and unexpected value overwrites.
  - Fix: Define an unexported type (`type contextKey string`) and use that as the key.

- Mistake: Storing mutable data (slices, maps, pointers) in context values.
  - Why: The context is shared across goroutines. If one goroutine modifies the slice, another goroutine reading the same context value sees the mutation — a data race.
  - Fix: Store only immutable values (strings, ints, structs by value) or defensively copy.

- Mistake: Using context values as optional function parameters.
  - Why: This hides dependencies. A function that accepts `ctx context.Context` and reads values implicitly depends on those values being set by the caller — but the signature does not document this.
  - Fix: Use explicit function parameters for configuration and dependencies. Use context values only for request-scoped metadata.

- Mistake: Assuming `ctx.Value(key)` is a constant-time lookup.
  - Why: `Value` walks the context tree linearly. Deep context trees (dozens of `WithValue` calls) make lookups O(depth).
  - Fix: Keep the number of context values small (1–5 per request) to maintain shallow trees.

## Debugging walkthrough

Consider a handler that intermittently panics with a nil pointer dereference when reading user claims from context:

```go
func handler(ctx context.Context) {
    claims := ctx.Value(keyClaims).(*Claims)
    fmt.Println("User:", claims.UserID) // panic: nil pointer
}
```

**Symptom**: Panics under certain request paths but not others.

**Investigation**: Check where context values are set:

```go
// In middleware
func authMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        token := r.Header.Get("Authorization")
        if token == "" {
            http.Error(w, "unauthorized", 401)
            return
        }
        claims := parseClaims(token)
        ctx := context.WithValue(r.Context(), keyClaims, claims)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
```

**Root cause**: When the `Authorization` header is missing, the middleware returns early without calling `next`. The handler downstream never has `keyClaims` in its context. `ctx.Value(keyClaims)` returns `nil`, and the type assertion `.(*Claims)` returns `nil`. The handler dereferences it and panics.

**Fix**: Either ensure all paths set the context value, or check the value before using it:

```go
claims, ok := ctx.Value(keyClaims).(*Claims)
if !ok {
    http.Error(w, "unauthorized", 401)
    return
}
```

## Production notes

- Limit context values to request-scoped metadata: request IDs, user IDs, auth tokens, trace IDs.
- Never use context values for: database connections, logger instances, configuration, or any mutable state.
- Always use a custom unexported key type to prevent collisions between packages.
- Document every context key your package expects. Preferably, provide accessor functions:
  ```go
  func RequestIDFromContext(ctx context.Context) (string, bool) {
      id, ok := ctx.Value(keyRequestID).(string)
      return id, ok
  }
  ```
- In middleware, prefer returning early without calling the next handler if the context value cannot be set.

## Performance implications

- `WithValue` allocates a `valueCtx` struct (~24 bytes + interface header for key and value) per call.
- `Value(key)` is O(depth) — each node comparison is an interface equality check (~2ns).
- For trees of depth 5–10, the cost is negligible (~20–50ns). For depths > 50, it adds up.
- Type assertions in hot paths (millions of requests per second) add measurable cost. Store values as concrete types and assert once early.

## Practice task

Write a middleware chain: a function `addRequestID(next http.Handler) http.Handler` that generates a UUID (or use `fmt.Sprintf("req-%d", time.Now().UnixNano())`) and stores it in the context. Write a handler that reads the request ID from the context and prints it. Use a custom unexported key type. Run it as a simple HTTP server on a random port and verify the request ID appears in the output.

## Tests / verification

```bash
go run ./curriculum/modules/11-lifecycle-context-concurrency/lessons/09-context-values-with-caveats
go test ./curriculum/modules/11-lifecycle-context-concurrency/lessons/09-context-values-with-caveats
```

The existing tests verify value storage, missing key behavior, multiple keys, type safety across different key types, and the immutability caveat. After completing the practice task, add tests for the `addRequestID` middleware.

## Review questions

1. Why should you use a custom unexported type for context keys rather than a plain string?
2. What does `ctx.Value(key)` return if the key was never set?
3. Why should you not store a slice in a context value?
4. How does `context.WithValue` perform lookup internally?
5. When is it appropriate to use context values, and when should you use explicit function parameters instead?

## NEXT UP

Why concurrency exists — the fundamental reasons for concurrency and Go's model for expressing it.
