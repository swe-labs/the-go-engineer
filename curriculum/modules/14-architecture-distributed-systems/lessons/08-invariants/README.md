# Invariants

## Learning objective

Identify and enforce business invariants at system boundaries, implement consistency checks in Go, and prevent invalid state transitions through encapsulation and validation.

## Why this matters

An invariant is a condition that must always be true for the system to be correct. "A user's email must be unique." "An order total must equal the sum of its line items." "An account balance must never go negative." When invariants are not enforced, data corruptions silently accumulate. A user gets two accounts with the same email. An order total does not match its items. A bank account goes negative. Enforcing invariants is how you keep data trustworthy.

## Mental model

Think of invariants as guardrails on a highway. The guardrails do not tell you where to drive, but they prevent you from driving into a ravine. You can still change lanes, speed up, slow down, or take an exit. But you cannot leave the road.

In code, invariants work the same way. You can create orders, add items, change addresses, and process payments. But you cannot create an order without items. You cannot ship an unpaid order. You cannot have a negative account balance. The guardrails (invariant checks) prevent the system from entering invalid states.

## Core idea

There are three levels of invariant enforcement:

1. **Constructor/Factory invariants**: enforced when creating an object. A `NewAccount` constructor rejects a negative initial balance. A `NewOrder` constructor rejects empty customer names. These invariants ensure every object is born in a valid state.

2. **Method invariants**: enforced when transitioning state. `account.Withdraw(amount)` checks that the balance is sufficient. `order.Ship()` checks that the order is pending. These invariants ensure every state transition is valid.

3. **Cross-object invariants**: enforced across multiple objects or aggregates. A transfer between accounts must debit one and credit the other atomically. An order cancellation must release inventory. These invariants require coordinated actions, often within a transaction or saga.

Invariants should be enforced at system boundaries -- the edges where data enters the system. Inside the boundary, you can assume the data is valid. This is the "validate once" principle: validate when data enters the system, and trust it afterward.

## Under the hood

Invariant enforcement in Go is done through:

- **Validation in constructors**: `NewAddress` returns an error if fields are empty. The caller cannot obtain an invalid `Address` value.
- **Guard clauses in methods**: `Withdraw` checks `amount > balance` and returns an error before mutating state. The guard clause is the first thing in the method.
- **Unexported fields**: preventing direct field access forces callers to use methods that enforce invariants.
- **Sentinel errors**: `var ErrInsufficientBalance = errors.New("insufficient balance")` lets callers distinguish between invariant violations and infrastructure errors.

Go has no built-in contract programming (preconditions, postconditions, invariants). Instead, these are implemented as explicit checks at the beginning and end of methods. The compiler does not enforce them; the programmer does.

## How Go uses it

Invariant enforcement is pervasive in the Go standard library:

- `bytes.Buffer.Write` returns an error if the buffer has been closed.
- `os.File.Close` can be called multiple times but returns an error on subsequent calls.
- `sql.DB.Close` is safe to call multiple times.
- `json.Decoder.Decode` returns an error if the input is malformed.

The standard library uses the Go idiom: return an error for invalid operations rather than panicking or silently ignoring the violation. Sentinel errors like `io.EOF` and `sql.ErrNoRows` let callers handle specific invariant violations.

## Go example

```go
package main

import (
	"errors"
	"fmt"
)

type Account struct {
	ID      string
	Owner   string
	Balance float64
}

func NewAccount(id, owner string, initialBalance float64) (*Account, error) {
	if id == "" {
		return nil, errors.New("account ID is required")
	}
	if owner == "" {
		return nil, errors.New("owner is required")
	}
	if initialBalance < 0 {
		return nil, errors.New("initial balance cannot be negative")
	}
	return &Account{ID: id, Owner: owner, Balance: initialBalance}, nil
}

func (a *Account) Deposit(amount float64) error {
	if amount <= 0 {
		return errors.New("deposit amount must be positive")
	}
	a.Balance += amount
	return nil
}

func (a *Account) Withdraw(amount float64) error {
	if amount <= 0 {
		return errors.New("withdrawal amount must be positive")
	}
	if amount > a.Balance {
		return fmt.Errorf("insufficient balance: have %.2f, need %.2f", a.Balance, amount)
	}
	a.Balance -= amount
	return nil
}

func Transfer(from, to *Account, amount float64) error {
	if from == nil || to == nil {
		return errors.New("both accounts must be provided")
	}
	if from.ID == to.ID {
		return errors.New("cannot transfer to the same account")
	}
	if amount <= 0 {
		return errors.New("transfer amount must be positive")
	}
	if err := from.Withdraw(amount); err != nil {
		return fmt.Errorf("withdrawal failed: %w", err)
	}
	if err := to.Deposit(amount); err != nil {
		from.Balance += amount
		return fmt.Errorf("deposit failed, rollback: %w", err)
	}
	return nil
}

func main() {
	alice, _ := NewAccount("acc-1", "Alice", 500.00)
	bob, _ := NewAccount("acc-2", "Bob", 100.00)

	fmt.Printf("Before: Alice=%.2f, Bob=%.2f\n", alice.Balance, bob.Balance)
	if err := Transfer(alice, bob, 200.00); err != nil {
		fmt.Println("Transfer failed:", err)
		return
	}
	fmt.Printf("After:  Alice=%.2f, Bob=%.2f\n", alice.Balance, bob.Balance)
}
```

## Step-by-step execution

For `Transfer(alice, bob, 200.00)` with Alice having 500.00 and Bob having 100.00:

1. `Transfer` checks both accounts are non-nil. Pass.
2. `Transfer` checks the accounts are different (`from.ID != to.ID`). Pass (acc-1 vs acc-2).
3. `Transfer` checks the amount is positive. 200.00 passes.
4. `from.Withdraw(200.00)` is called.
5. Inside `Withdraw`: checks amount positive (pass), checks `200.00 <= 500.00` (pass), deducts 200.00 from Alice's balance (now 300.00).
6. `to.Deposit(200.00)` is called.
7. Inside `Deposit`: checks amount positive (pass), adds 200.00 to Bob's balance (now 300.00).
8. `Transfer` returns nil.

Now consider `Transfer(alice, bob, 600.00)`:

1-3. Same checks pass.
4. `from.Withdraw(600.00)` is called.
5. Inside `Withdraw`: checks `600.00 <= 500.00` fails. Returns error "insufficient balance".
6. `Transfer` returns the error wrapped. Alice's balance is untouched (still 500.00).

The invariant "account balance must never be negative" is enforced by `Withdraw`. The invariant "transfer must be atomic" is enforced by `Transfer` -- if the deposit fails, the withdrawal is rolled back.

## Common mistakes

- **Skipping validation in constructors**: creating an `Account` with a negative balance because the constructor just assigns fields. The invariant is only checked elsewhere, inconsistently. Always validate in the constructor.
- **Enforcing invariants only in the service layer**: the handler validates, the service validates, the repository validates -- the same check in three places. The invariant should be enforced at the domain level (in the constructor or method) and only there.
- **Panicking on invariant violations**: `if amount > balance { panic("negative balance") }`. This crashes the entire process. Return an error instead. The caller decides how to handle it (retry, reject, log).
- **Ignoring rollback on partial failure**: the transfer succeeds but the deposit fails due to a network error. The money is lost. `Transfer` rolls back the withdrawal on deposit failure. Consider this in every multi-step operation.

## Debugging walkthrough

```go
func (a *Account) Withdraw(amount float64) {
    a.Balance -= amount
    if a.Balance < 0 {
        log.Printf("WARNING: account %s is negative: %.2f", a.ID, a.Balance)
    }
}
```

**Symptom**: Account balances are occasionally negative. The logs show warnings but no errors.

**Root cause**: The invariant check is after the mutation. The balance goes negative briefly, and the function only logs a warning instead of preventing the mutation.

**Fix**: Check before mutating:

```go
func (a *Account) Withdraw(amount float64) error {
    if amount <= 0 { return errors.New("...") }
    if amount > a.Balance { return errors.New("...") }
    a.Balance -= amount
    return nil
}
```

The guard clause ensures the mutation only happens when the invariant holds.

## Production notes

- **Validate once at boundaries**: validate input at the API boundary (HTTP handler, gRPC interceptor, message consumer). The domain model trusts that the validated data is correct. This prevents redundant validation.
- **Idempotency invariants**: some invariants must hold even under retries. "A payment must be processed at most once." Use idempotency keys to detect duplicate operations.
- **Concurrent invariant enforcement**: the in-memory `Account.Withdraw` is not safe for concurrent use. In production, use database transactions with row-level locking. The invariant is enforced by the database: `UPDATE accounts SET balance = balance - $1 WHERE id = $2 AND balance >= $1`.
- **Audit logging**: when an invariant enforcement fails (e.g., attempt to withdraw from an empty account), log the event. It may indicate a bug, a malicious actor, or a race condition.

## Performance implications

- **Validation overhead**: each invariant check is a comparison and a branch. In normal flow, this is immeasurable. In high-frequency trading or real-time systems, minimize checks in the hot path by validating before entering the hot loop.
- **Rollback cost**: in the `Transfer` example, rolling back involves re-adding the withdrawn amount. For in-memory operations this is negligible. For database operations, rollback means an additional UPDATE statement. Consider using database transactions for atomic multi-step operations.
- **Error allocation**: returning `fmt.Errorf(...)` allocates a new error value. For operations that frequently fail (e.g., withdrawing from an account that is often empty), consider using sentinel errors that are pre-allocated.

## Practice task

Extend the account example with a `Freeze()` and `Unfreeze()` feature:

1. Add a `frozen` boolean field to `Account` (unexported).
2. `Freeze()` sets `frozen = true`. Can only be called if the account is not already frozen.
3. `Unfreeze()` sets `frozen = false`. Can only be called if the account is frozen.
4. `Withdraw` and `Deposit` return an error if the account is frozen.
5. Write table-driven tests for all valid and invalid state transitions.

## Tests / verification

```bash
go run ./curriculum/modules/14-architecture-distributed-systems/lessons/08-invariants
go test ./curriculum/modules/14-architecture-distributed-systems/lessons/08-invariants
```

The tests verify account creation validation, deposit and withdrawal rules, error handling for invalid amounts and insufficient balance, and the transfer function with rollback. After completing the practice task, add tests for freeze/unfreeze operations.

## Review questions

1. What is a business invariant? Give three examples from a banking domain.
2. Why should invariants be enforced in constructors and methods rather than in the service layer?
3. How does the `Transfer` function ensure atomicity across two accounts?
4. What is the difference between checking invariants before a mutation vs after?
5. How would you enforce the invariant "a user must not have more than 5 active orders" at the domain level?

## NEXT UP

Event-driven basics -- decoupling producers and consumers with an in-memory event bus in Go.
