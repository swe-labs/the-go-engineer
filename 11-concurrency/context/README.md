# Section 17: Context

## Learning Objectives

`context.Context` is the backbone of every production Go application. It flows through every function call that does I/O, carrying cancellation signals, deadlines, and request-scoped values.

By the end of this section, you will understand:

- What `context.Context` is and why it exists
- `context.Background()` and `context.TODO()`
- `context.WithCancel()` for manual cancellation
- `context.WithTimeout()` and `context.WithDeadline()` for automatic deadlines
- `context.WithValue()` for request-scoped data
- Propagating context through function chains
- Context in HTTP handlers (`r.Context()`)

## Beginner → Expert Mapping

| Topic | Level | Importance | Engineering Concept |
| ----- | ----- | ---------- | ------------------- |
| Background & TODO | Beginner | High | Root context creation |
| WithCancel | Intermediate | **Critical** | Manual cancellation propagation |
| WithTimeout | Intermediate | **Critical** | Automatic deadline enforcement |
| WithValue | Advanced | Medium | Request-scoped metadata |
| HTTP Context | Advanced | **Critical** | Server request lifecycle |

## Engineering Depth

Google's internal Go style guide mandates context as the first parameter of any function that does I/O or may be long-running. The signature `func DoSomething(ctx context.Context, ...)` is not optional — it's the law.

Context solves the "cancellation propagation" problem: when a user cancels an HTTP request, how do you stop the database query, the API call, and the file write that are all happening downstream? Context creates a tree of cancellation signals that flow from parent to child automatically.

## Contents

| Directory | Topic | Level |
| --------- | ----- | ----- |
| `1-background/` | `context.Background()`, `context.TODO()` | Beginner |
| `2-with-cancel/` | Manual cancellation with `context.WithCancel` | Intermediate |
| `3-with-timeout/` | Automatic deadlines with `context.WithTimeout` | Intermediate |
| `4-with-value/` | Request-scoped data with `context.WithValue` | Advanced |

## How to Run

```bash
go run ./17-context/1-background
go run ./17-context/2-with-cancel
go run ./17-context/3-with-timeout
go run ./17-context/4-with-value
```

---

## 🏗 Exercise: Timeout-Aware API Client (`5-timeout-client`)

Build an HTTP client that uses `context.WithTimeout` to enforce request deadlines. Try it yourself first!

```bash
go run ./17-context/5-timeout-client/_starter   # Try the exercise
go run ./17-context/5-timeout-client            # See the solution
```

## References

- [Go Blog: Go Concurrency Patterns: Context](https://go.dev/blog/context)
- [Package context documentation](https://pkg.go.dev/context)
