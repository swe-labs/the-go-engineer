# Mockery

## When you need this

Testing Go code that depends on external interfaces (database, HTTP clients, message queues) requires test doubles that implement those interfaces. Writing mock implementations by hand is tedious, error-prone, and must be updated when the interface changes. Mockery solves this by generating mock implementations from Go interfaces automatically, using go:generate directives or CLI commands.

## When you do not need this

When simpler concrete types or standard core library primitives suffice without adding abstraction overhead.

## Prerequisites

- elective-07

## Mental model

Mockery reads a Go interface and generates a struct that implements it, with methods for setting expectations (On, Return), verifying calls (AssertExpectations), and tracking call count. The generated mock can be configured to return specific values, panic, or call through to a real implementation. Tests use the mock in place of the real dependency and assert that expected methods were called with the correct arguments.

## How it works

Mockery uses go/packages to load and parse Go source files, extracts interface definitions, and generates a Go source file with a mock struct that implements the interface. Each mock method records call arguments and returns configured values. Expectations are stored in a map keyed by method name + argument values. The generated code uses sync.Mutex for thread-safe expectation tracking.

## Go example

```go
package main

import "fmt"

func main() {
    fmt.Println("Elective: Mockery")
}
```

## Tradeoffs

- Mock method calls use testify's reflection-based matching — each expectation check is O(n) in the number of expectations for that method
- Creating mock instances allocates: testify stores call expectations in slices — large test suites with many mocks can have significant setup overhead
- Mock assertions at test teardown (AssertExpectations) iterate over all expectations — negligible for typical tests (<100 expectations)
- For performance-critical test execution (1000+ test cases), consider hand-written stub implementations instead of mockery-generated mocks

## Common mistakes

- Testing implementation details — asserting that a specific method was called when only the result matters, making tests brittle to refactoring.
- Not calling AssertExpectations in tests — expectations are never verified, and the test passes even if the mock was never called.
- Over-mocking — mocking interfaces that could use real implementations, creating fragile tests that fail when mock expectations are incorrect.
- Setting expectations that are too strict — using .Once() and exact argument matchers when the implementation could call the method multiple times or with slightly different arguments.
- Not regenerating mocks when interfaces change — the generated mock becomes stale and does not compile with the updated interface.

## Practice task

The learner must define a Go interface, generate a mock using Mockery with a //go:generate directive, write a test that injects the mock, sets expectations, exercises the code, and calls AssertExpectations — including a negative test that proves AssertExpectations fails when expected calls are not made.

## Production notes

Mockery is used in every Go project that follows interface-based testing: repository pattern tests (mock the database interface), HTTP handler tests (mock the service interface), event processing tests (mock the publisher interface). Production Go services use Mockery with go:generate directives so mocks are always in sync with interfaces. CI runs go generate ./... to regenerate mocks on interface changes.

## Related core lessons

- Module 16.07 — go generate
- Module 17.06 — Build tags
