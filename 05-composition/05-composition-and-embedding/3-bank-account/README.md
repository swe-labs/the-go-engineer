# CO.3 Bank Account Project

## Mission

Build a small bank-account model that proves the difference between ordinary composition,
embedding, promoted methods, and method shadowing.

This exercise is the Section 06 milestone.
It is where named-field composition and embedding come together in one runnable artifact with
tests.

## Prerequisites

Complete these first:

- `CO.1` composition
- `CO.2` embedding

## What You Will Build

Implement a small bank-account model that:

1. defines a reusable `Account` type with shared deposit and withdrawal behavior
2. embeds `Account` inside `SavingsAccount`
3. adds interest through promoted fields and methods on `SavingsAccount`
4. embeds `Account` inside `OverdraftAccount`
5. shadows `Withdraw` on `OverdraftAccount` to allow overdraft behavior
6. demonstrates the final behavior in `main()`
7. passes the provided tests

## Files

- [main.go](./main.go): complete solution with teaching comments
- [account_test.go](./account_test.go): tests for deposit, withdraw, interest, and overdraft rules
- [_starter/main.go](./_starter/main.go): starter file with TODOs and requirements

## Run Instructions

Run the completed solution:

```bash
go run ./05-composition./05-composition-and-embedding/3-bank-account
```

Run the tests:

```bash
go test ./05-composition./05-composition-and-embedding/3-bank-account
```

Run the starter:

```bash
go run ./05-composition./05-composition-and-embedding/3-bank-account/_starter
```

## Success Criteria

Your finished solution should:

- keep the shared account behavior in one reusable embedded type
- use promotion intentionally instead of copying fields or methods around
- shadow `Withdraw` only where the overdraft rule actually changes behavior
- keep the data model easy to explain without inheritance language
- pass the provided tests

## Common Failure Modes

- copying `Deposit` and `Withdraw` into every account type instead of reusing `Account`
- describing embedding as inheritance instead of promoted access to a composed type
- shadowing methods without being able to explain why the outer behavior differs
- mutating balances with value receivers so updates silently affect copies

## Next Step

After you complete this exercise, continue to [Section 07](../../../06-strings-and-text) if you
are ready to move on.



