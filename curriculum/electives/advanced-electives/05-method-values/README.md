# Method values

## When you need this

Go developers often need to pass a method as a callback — for example, passing an HTTP handler's method to a router or passing a struct method to a sort.Interface. Method values (bound methods) and method expressions (unbound methods) provide type-safe ways to reference and call methods as first-class function values.

## When you do not need this

When simpler concrete types or standard core library primitives suffice without adding abstraction overhead.

## Prerequisites

- elective-04

## Mental model

A method value is a function closure that has a receiver bound to it. When you write f := obj.Method, f is a function value that, when called, calls Method on obj. The receiver is captured at the point of assignment — if you change obj later, f still calls Method on the original obj. A method expression is the unbound function: T.Method returns a function whose first parameter is the receiver.

## How it works

A method value m := obj.Method creates a closure that holds a pointer to obj. When called, it invokes Method on the captured pointer. This is equivalent to func(args...) { return obj.Method(args...) }. A method expression T.Method creates a function where the receiver is the first explicit parameter: func(receiver T, args...) returns. This allows passing methods as callbacks where the caller supplies the receiver, useful for sort.Interface and http.Handler patterns.

## Go example

```go
package main

import "fmt"

func main() {
    fmt.Println("Elective: Method values")
}
```

## Tradeoffs

- Method values allocate: the closure that captures the receiver is allocated on the heap if it escapes the creating function
- Method value call is a single indirect function call — comparable to calling a function pointer
- In hot paths, method values can prevent inlining because the compiler cannot see through the closure — prefer direct calls or method expressions where performance matters
- For stateless functions (no receiver), prefer function values instead of method values to avoid the closure allocation

## Common mistakes

- Confusing method values (bound) with method expressions (unbound) — a method value captures the receiver, while a method expression requires the receiver as the first argument.
- Capturing a method value inside a loop — the method captures the loop variable by reference, and all callbacks operate on the last iteration's value.
- Using a method value when a method expression is needed — passing obj.Method to sort.Sort does not work because sort.Interface expects Len, Less, Swap as separate methods, not a single function.
- Assuming method values are free — each method value allocation creates a closure, adding GC pressure in hot paths.
- Using method values with value receivers on large structs — the entire struct is copied into the closure.

## Practice task

The learner must demonstrate the difference between method values and method expressions, write a sort.Interface implementation using method expressions, and use a method value as a callback in an event emitter pattern.

## Production notes

Method values are used in every Go HTTP service (http.HandlerFunc), every sort implementation (sort.Sort), every event emitter (registering callbacks), and every test suite (registering test helpers). Production Go code uses method values for dependency injection (passing service methods to HTTP handlers) and middleware chaining.

## Related core lessons

- Module 16.04 — Functional options
- Module 16.15 — Demo preparation
