# TE.1-TE.3 Testing Fundamentals

## Mission

Learn the foundational Go testing patterns through one small user-focused package.

This surface covers:

- `TE.1` unit testing
- `TE.2` table-driven tests
- `TE.3` HTTP handler testing

## Files

- [user.go](./user.go): code under test plus teaching comments
- [user_test.go](./user_test.go): basic tests and table-driven tests
- [greeting_test.go](./greeting_test.go): testable design with `io.Writer`
- [http_handler_test.go](./http_handler_test.go): handler testing with `httptest`

## Run Instructions

```bash
go test ./08-quality-test/01-quality-and-performance/testing/user
```

## Success Criteria

You should be able to:

- write small focused tests with `testing.T`
- structure table-driven tests with `t.Run`
- test handler behavior without starting a real server
- explain why dependency injection makes code easier to test


## Prerequisites

You should be comfortable with Go syntax, basic data structures, and the control flow mechanics covered in earlier sections.

## Mental Model

Think of this as the conceptual blueprint. The components interact by exchanging state, defining clear boundaries between what is requested and what is provided.

## Visual Model

Visualizing this process involves tracing the execution path from the input entry point, through the processing layers, and out to the final output or side effect.

## Machine View

At the hardware level, this translates into specific memory allocations, CPU instruction cycles, and OS-level system calls to manage resources efficiently.

## Code Walkthrough

We step through the code sequentially, examining how the interfaces are satisfied, where the errors are checked, and how the core loop manages control flow.

## Try It

Run the code locally. Modify the inputs, toggle the conditions, and observe how the output shifts. Experimentation is the fastest way to cement your understanding.

## Prerequisites

You should be comfortable with Go syntax, basic data structures, and the control flow mechanics covered in earlier sections.

## Mental Model

Think of this as the conceptual blueprint. The components interact by exchanging state, defining clear boundaries between what is requested and what is provided.

## Visual Model

Visualizing this process involves tracing the execution path from the input entry point, through the processing layers, and out to the final output or side effect.

## Machine View

At the hardware level, this translates into specific memory allocations, CPU instruction cycles, and OS-level system calls to manage resources efficiently.

## Code Walkthrough

We step through the code sequentially, examining how the interfaces are satisfied, where the errors are checked, and how the core loop manages control flow.

## Try It

Run the code locally. Modify the inputs, toggle the conditions, and observe how the output shifts. Experimentation is the fastest way to cement your understanding.

## Prerequisites

You should be comfortable with Go syntax, basic data structures, and the control flow mechanics covered in earlier sections.

## Mental Model

Think of this as the conceptual blueprint. The components interact by exchanging state, defining clear boundaries between what is requested and what is provided.

## Visual Model

Visualizing this process involves tracing the execution path from the input entry point, through the processing layers, and out to the final output or side effect.

## Machine View

At the hardware level, this translates into specific memory allocations, CPU instruction cycles, and OS-level system calls to manage resources efficiently.

## Code Walkthrough

We step through the code sequentially, examining how the interfaces are satisfied, where the errors are checked, and how the core loop manages control flow.

## Try It

Run the code locally. Modify the inputs, toggle the conditions, and observe how the output shifts. Experimentation is the fastest way to cement your understanding.

## In Production

Tests are the structural foundation of refactoring. In a production codebase, engineers rarely write code once and never touch it again; they constantly modify it to add features, fix bugs, or optimize performance. Without a reliable test suite, every change carries the risk of silent regressions. The patterns taught here — specifically table-driven tests and `httptest` — are ubiquitous in professional Go. Table-driven tests keep test files concise and make it trivial to add new edge cases when production bugs are discovered. The `httptest` package allows testing HTTP handlers directly as Go functions, entirely bypassing the network stack, which keeps tests extremely fast and stable. Teams that design for testability by using interfaces (like `io.Writer`) instead of hardcoding concrete dependencies (like `os.Stdout`) find that their code is naturally more modular and easier to maintain.

## Thinking Questions

1. Why are table-driven tests preferred in Go over writing a separate `Test...` function for every single input variation?
2. What happens to a test suite if an HTTP handler test binds to a real port like `:8080` instead of using `httptest.NewRecorder`?
3. If a function's logic is perfectly correct but it is impossible to test without a live database, what is wrong with the function's design?
4. How do you balance between writing too many granular unit tests (which make refactoring difficult) and too few (which allow bugs to slip through)?

## Next Step

After this surface, continue to [TE.4 benchmarking](../benchmarks) or back to the
[Testing track](../README.md).

