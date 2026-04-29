# FE.10 panic and recover

## Mission

Learn when panic is appropriate, when it is not, and how recover turns a crash into an explicit boundary decision.

## Why This Lesson Exists Now

You have learned how to handle ordinary, expected failures using the `error` interface. But sometimes, a program encounters a state that is fundamentally "wrong" or "impossible"-like trying to access a database that *must* be connected, or discovering a corrupted configuration.

In these rare cases, Go provides `panic`. It stops the normal execution of the program and starts "unwinding" the stack. `recover` is the mechanism to catch that panic and stop the program from crashing entirely.

> **Backward Reference:** In [Lesson 7: Order Summary](../7-order-summary/README.md), you handled validation errors gracefully. `panic` is the opposite of that: it is the nuclear option for when things have gone so wrong that the program cannot reasonably continue its current job.

## Prerequisites

- `FE.7` order summary
- `FE.9` closures - mechanics (for understanding how `defer` with an anonymous function works)

## Mental Model

Errors describe expected failure. Panic describes broken assumptions. Recover belongs at process or request boundaries, not in ordinary business flow.

## Visual Model

```mermaid
graph TD
    A["Normal Flow"] --> B["panic('boom')"]
    B --> C["Unwind Stack"]
    C --> D["Run Deferred Functions"]
    D --> E{recover() called?}
    E -- Yes --> F["Stop Unwinding, Continue Program"]
    E -- No --> G["Crash Process"]
```

```text
defer func() {
    if r := recover(); r != nil {
        // Safe again!
    }
}()

panic("something impossible happened")
```

## Machine View

A `panic` starts stack unwinding. Go stops running the normal lines of your code and immediately starts running all `defer` statements in the current function, then the caller's `defer` statements, and so on up the stack.

The `recover` function only works inside a `defer` call on the same goroutine. If `recover` is called while the program is panicking, it returns the value that was passed to `panic` and stops the unwinding process.

## Run Instructions

```bash
go run ./03-functions-errors/10-panic-and-recover
```

## Code Walkthrough

### `defer func() { if r := recover(); r != nil { ... } }()`

This is the recovery pattern. It must be `defer`red at the very start of the function you want to protect. If a panic happens anywhere in `accessDatabase`, this function will run during the stack unwind and catch the panic value.

### `panic("database connection lost...")`

This line triggers the panic. Normal execution stops here. Go skips the `fmt.Println` below it and jumps straight to the deferred recovery function.

### `if !connected { ... }`

Notice that we only panic when the state is truly unrecoverable. For ordinary problems (like a missing user), you should still use `error`.

### `Program continued running...`

Because we recovered, the `main` function continues to run after `accessDatabase` finishes. If we hadn't recovered, the whole program would have crashed.

## Try It

1. Comment out the `defer` block in `accessDatabase` and run the program. Observe the crash and the stack trace.
2. Change the panic message and see it reflected in the recovery output.
3. Add a second `panic` after the recovery and see if it is caught (hint: you'd need another `defer` for the second one).

## Common Questions

- Should I use `panic` for validation?
  No. Use `error` for anything that is expected to happen (like bad user input). Use `panic` only for programmer errors or "impossible" system states.
  
- Can I recover from a panic in a different goroutine?
  No. `recover` only works on the goroutine where the panic started.

## In Production
Use panic sparingly for programmer bugs or impossible states, then recover only at boundaries (like a web server's request handler) where you can translate the crash into a log entry and a "500 Internal Server Error" response instead of letting the whole server die.

## Thinking Questions
1. What problem does this topic solve?
2. What breaks if this boundary is handled implicitly instead of explicitly?
3. Where would you expect to use this topic in production Go code?

> **Forward Reference:** You have completed the Functions and Errors section! You are now ready to move into data modeling and system design. In the next section, [Section 04: Types and Design](../../04-types-design/README.md), you will learn how to group data together using **Structs**.

## Next Step

Continue to [Section 04: Types and Design](../../04-types-design/README.md) and start with [Lesson 1: Structs](../../04-types-design/1-struct/README.md).
