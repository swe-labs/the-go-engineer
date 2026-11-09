# Bank Account

## Learning objective

Build a thread-safe, polymorphic bank account ledger system in Go. Master custom types, method receivers (value vs pointer semantics), implicit interface satisfaction, generic transaction ring buffers, nil interface safety, and race-free concurrent state mutations.

## Why this matters

Go does not have class inheritance. Instead, Go architectures rely on composition, method sets, and interfaces satisfied implicitly. Choosing between value and pointer receivers dictates whether state mutations persist or trigger silent data bugs. Understanding the difference between `nil` and an interface holding a typed nil pointer (`(*T)(nil)`) is a hallmark of an experienced Go backend engineer.

## Prerequisites

Complete all Module 05 lessons (01–23) on struct definitions, value vs pointer receivers, interfaces, method sets, type switches, nil interfaces, and generics.

## Available concepts

This project draws on structs, custom types, method sets, pointer receivers for state mutation, `Account` interface polymorphism, `sync.RWMutex`, generic collections (`[T any]`), and `go test -race`.

## Project overview

The `_starter/` directory provides type signatures, error definitions, and interface declarations. Your job is to implement the concrete account types, generic transaction ring buffer, and transfer logic so that:

- `CheckingAccount` allows withdrawals up to its balance plus an overdraft limit.
- `SavingsAccount` restricts withdrawals that would breach its minimum balance threshold and supports interest accrual.
- Mutating methods use pointer receivers to mutate internal balance and append to transaction history.
- The generic `TransactionBuffer[T]` stores records up to its capacity, evicting the oldest element on overflow.
- All account operations are protected with `sync.RWMutex` to guarantee race condition freedom under concurrent access.

## Architecture

```
                       +-------------------+
                       | Account Interface |
                       +-------------------+
                       | ID() string       |
                       | Owner() string    |
                       | Balance() int     |
                       | Deposit(int) err  |
                       | Withdraw(int) err |
                       | Transactions()    |
                       +-------------------+
                                ^
               +----------------+----------------+
               |                                 |
     +-------------------+             +-------------------+
     |  CheckingAccount  |             |  SavingsAccount   |
     | - overdraftLimit  |             | - minimumBalance  |
     | - buffer[Tx]      |             | - AccrueInterest  |
     +-------------------+             +-------------------+
```

## Tasks

### Task 1 — Implement generic transaction buffer

In `_starter/main.go`:
- Implement `TransactionBuffer[T]`:
  - `Push(item T)`: appends item under write lock; if buffer exceeds capacity, slices off the oldest element.
  - `All() []T`: returns a copied slice snapshot under read lock.
  - `Len() int`: returns count of stored items.

### Task 2 — Implement CheckingAccount

- Implement `NewCheckingAccount(id, owner string, initialDeposit, overdraftLimit int) (*CheckingAccount, error)`.
- Implement `Deposit(amountCents int) error`: requires positive amount; updates balance and logs transaction.
- Implement `Withdraw(amountCents int) error`: verifies `balance + overdraftLimit >= amount`; updates balance and logs transaction.
- Ensure all mutating methods use pointer receivers `*CheckingAccount` and thread-safe locking.

### Task 3 — Implement SavingsAccount

- Implement `NewSavingsAccount(id, owner string, initialDeposit, minBalance int) (*SavingsAccount, error)`: rejects initial deposits below `minBalance`.
- Implement `Withdraw(amountCents int) error`: verifies remaining balance after withdrawal does not breach `minimumBalance`.
- Implement `AccrueInterest(ratePercent float64) error`: calculates interest and adds to balance.

### Task 4 — Implement atomic Transfer

- Implement `Transfer(from Account, to Account, amountCents int) error`:
  - Validates non-nil accounts and distinct IDs.
  - Withdraws from sender; if deposit into receiver fails, rolls back withdrawal.

## Verification

Run all unit tests with the Go race detector enabled:

```bash
go test -v -race ./curriculum/modules/05-types-interfaces-packages-modules/projects/bank-account/_starter
```

Expected output:
```
=== RUN   TestCheckingAccount_DepositAndWithdrawal
--- PASS: TestCheckingAccount_DepositAndWithdrawal (0.00s)
=== RUN   TestSavingsAccount_MinimumBalanceAndInterest
--- PASS: TestSavingsAccount_MinimumBalanceAndInterest (0.00s)
=== RUN   TestAccountPolymorphism_Transfer
--- PASS: TestAccountPolymorphism_Transfer (0.00s)
=== RUN   TestConcurrentAccountAccess
--- PASS: TestConcurrentAccountAccess (0.01s)
=== RUN   TestGenericTransactionBuffer
--- PASS: TestGenericTransactionBuffer (0.00s)
=== RUN   TestNilInterfaceTrap
--- PASS: TestNilInterfaceTrap (0.00s)
PASS
```

## Rubric

| Criterion | Excellent (5) | Good (3) | Needs Work (0) |
|---|---|---|---|
| **Interface Semantics** | `Account` interface satisfied implicitly by both accounts; polymorphic transfer works smoothly | Interface defined but concrete types leak into helper functions | Interface bypassed with type casting |
| **Receiver Correctness** | Pointer receivers used for state mutation; value receiver semantics understood | Pointer receivers used but lock copied by value | Value receivers used causing mutations to be lost |
| **Concurrency & Thread Safety** | All access protected by `sync.RWMutex`; zero data races detected by race detector | Mutex used but read operations unsynchronized | Unsynchronized state mutations causing data races |
| **Generics & Buffering** | `TransactionBuffer[T]` properly parameterized and handles capacity eviction safely | Buffer works but tied strictly to `Transaction` | Generic syntax errors or slicing out-of-bounds panics |

## NEXT UP

[Payroll Processor](../payroll-processor/README.md)
