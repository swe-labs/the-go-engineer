# Project Assessment — Bank Account

## Purpose

Evaluate your ability to build a polymorphic, thread-safe system using Go interfaces, pointer receivers, generics, and `sync.RWMutex`. This assessment verifies that you understand implicit interface satisfaction, value vs pointer receiver semantics for state mutation, generic type parameters, and concurrent access safety.

## Targets

- `project-bank-account`

## Scope

This assessment covers:

- Custom struct types with private and public fields
- Pointer receivers for state mutation vs value receivers for read-only access
- Implicit interface satisfaction (no `implements` keyword)
- Interface polymorphism: processing mixed concrete types through a single interface
- Generic type parameters (`[T any]`)
- `sync.RWMutex` for concurrent read/write safety
- The `go test -race` flag for data race detection
- Nil interface traps: `(*T)(nil)` vs `nil`

## Evidence required

Submit:

1. Your completed `_starter/main.go` with all types, methods, and `Transfer` function implemented
2. Output of `go test -v -race .` from the `_starter/` directory showing all tests passing with zero race conditions
3. A brief explanation of why `Deposit` and `Withdraw` must use pointer receivers (`*CheckingAccount`) — what would happen with value receivers?

## Time expectation

Plan for 90 to 120 minutes.

## Rubric

| Criterion | Weight | Excellent (5) | Good (3) | Needs Work (0) |
|-----------|--------|--------------|----------|----------------|
| **Interface semantics** | 25% | `Account` interface satisfied implicitly by both `CheckingAccount` and `SavingsAccount`; polymorphic `Transfer` works with any `Account` | Interface defined but concrete type leaks into helper functions via type assertions | Interface bypassed; type casting used instead of polymorphism |
| **Receiver correctness** | 25% | Pointer receivers used for all state-mutating methods; can explain why value receivers would cause silent data loss | Pointer receivers used but rationale unclear; some methods inconsistent | Value receivers used for mutations causing lost updates |
| **Concurrency safety** | 25% | All access protected by `sync.RWMutex` (`RLock` for reads, `Lock` for writes); zero data races detected by `go test -race` | Mutex used but read operations use `Lock` instead of `RLock`, or lock copied by value | Unsynchronized state access causing data races |
| **Generics and buffering** | 15% | `TransactionBuffer[T]` properly parameterized with capacity eviction; ring buffer behavior correct | Buffer works but tied to a specific type instead of using generics | Generic syntax errors or out-of-bounds panics |
| **Transfer atomicity** | 10% | `Transfer` validates distinct accounts, withdraws atomically, and rolls back on deposit failure | Transfer works for happy path but no rollback on failure | Transfer panics or allows negative balances |

**Passing score**: 75%

## Review questions

Be prepared to answer these questions during a review session:

1. What is the method set of `CheckingAccount` vs `*CheckingAccount`? Why does this matter for interface satisfaction?
2. If you change `func (c *CheckingAccount) Deposit(...)` to `func (c CheckingAccount) Deposit(...)`, what breaks and why? Demonstrate with a test.
3. Explain the difference between a nil `Account` interface and an `Account` holding a `(*CheckingAccount)(nil)`. Why does the nil check `if acc == nil` fail for the second case?
4. Your `TransactionBuffer[T]` uses `sync.RWMutex`. Why do you use `RLock` for `All()` and `Len()` but `Lock` for `Push()`?
5. Walk through your `Transfer` function: what happens if `to.Deposit()` fails after `from.Withdraw()` succeeded? How do you recover?

## Passing standard

You pass when all tests pass with `-race`, your interface polymorphism works without type assertions, and you can explain receiver semantics and mutex usage mechanically.

## Retake policy

If you score below 75%, review Module 05 lessons on interfaces, method sets, receiver types, and sync primitives, then resubmit with corrections.

## NEXT UP

[Project Assessment: Payroll Processor](../payroll-processor/README.md)
