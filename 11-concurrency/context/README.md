# Section 11: Context

## Learning Objectives

`context.Context` is the backbone of production Go code that performs I/O. It carries cancellation
signals, deadlines, and request-scoped values through function chains.

By the end of this section, you should understand:

- what `context.Context` is and why it exists
- `context.Background()` and `context.TODO()`
- `context.WithCancel()` for manual cancellation
- `context.WithTimeout()` and `context.WithDeadline()` for automatic deadlines
- `context.WithValue()` for request-scoped data
- how context propagates through function chains
- how HTTP handlers use `r.Context()`

## Beginner → Expert Mapping

| Topic | Level | Importance | Engineering Concept |
| --- | --- | --- | --- |
| Background & TODO | Beginner | High | Root context creation |
| WithCancel | Intermediate | Critical | Manual cancellation propagation |
| WithTimeout | Intermediate | Critical | Automatic deadline enforcement |
| WithValue | Advanced | Medium | Request-scoped metadata |
| HTTP Context | Advanced | Critical | Server request lifecycle |

## Engineering Depth

Google's internal Go style guidance treats context as the first parameter of any function that does
I/O or may run for a meaningful amount of time. The signature
`func DoSomething(ctx context.Context, ...)` is not decoration — it is the contract that allows
callers to stop downstream work cleanly.

Context solves the cancellation-propagation problem: if a user cancels an HTTP request, the
database query, API call, and file write happening underneath it should stop too. Context creates a
tree of cancellation signals that flows from parent to child automatically.

## Contents

| Directory | Topic | Level |
| --- | --- | --- |
| `1-background/` | `context.Background()`, `context.TODO()` | Beginner |
| `2-with-cancel/` | Manual cancellation with `context.WithCancel` | Intermediate |
| `3-with-timeout/` | Automatic deadlines with `context.WithTimeout` | Intermediate |
| `4-with-value/` | Request-scoped data with `context.WithValue` | Advanced |

## How to Run

```bash
go run ./11-concurrency/context/1-background
go run ./11-concurrency/context/2-with-cancel
go run ./11-concurrency/context/3-with-timeout
go run ./11-concurrency/context/4-with-value
```

## Exercise: Timeout-Aware API Client (`5-timeout-client`)

Build an HTTP client that uses `context.WithTimeout` to enforce request deadlines.

```bash
go run ./11-concurrency/context/5-timeout-client/_starter
go run ./11-concurrency/context/5-timeout-client
```

## References

- [Go Concurrency Patterns: Context](https://go.dev/blog/context)
- [Package context documentation](https://pkg.go.dev/context)

## Learning Path

| ID | Lesson | Concept | Requires |
| --- | --- | --- | --- |
| CT.1 | [Background & TODO](./1-background) | Root context · Context interface (`Deadline`, `Done`, `Err`, `Value`) | entry |
| CT.2 | [WithCancel](./2-with-cancel) | cancel func · `ctx.Done()` · goroutine leak prevention · tree propagation | CT.1 |
| CT.3 | [WithTimeout](./3-with-timeout) | duration-based auto-cancel · `WithDeadline` · `DeadlineExceeded` | CT.1, CT.2 |
| CT.4 | [WithValue](./4-with-value) | private key type · request-scoped metadata · `O(depth)` lookup | CT.1, CT.2, CT.3 |
