# Contact Directory

## Learning objective

Implement a contact directory manager using Go data structures: slices, maps, and custom types.

## Why this matters

CRUD (Create, Read, Update, Delete) operations are the foundation of nearly every application. This project consolidates maps, slices, strings, and control flow into a practical contact management tool.

## Prerequisites

Complete all Module 03 lessons (01–23). Complete the [Pricing Checkout](../pricing-checkout/README.md) project first.

## Available concepts

This project draws on values, expressions, variables, types, zero values, type conversions, constants, boolean logic, if/else, switch, loops, range, arrays, slices, maps, comma-ok idiom, strings, and pointers.

## Project overview

The `_starter/` directory contains a Go program with missing function bodies. Your job is to implement each function correctly so that:

- The program compiles and runs without errors.
- All tests pass.
- `gofmt -l .` produces no output.
- `go vet .` produces no output.

You may check `_solution/` only after you have attempted all fixes yourself.

## Tasks

### Task 1 — Implement directory operations

Complete the following functions in `_starter/main.go`:

| Function | Behaviour |
|----------|-----------|
| `NewContact(name, email, phone string) Contact` | Creates a `Contact` with trimmed fields. If any field is empty after trimming, return a contact with empty fields and an error. |
| `AddContact(dir Directory, c Contact) error` | Adds a contact keyed by phone. Returns error if phone already exists or if contact has an empty phone. |
| `GetContact(dir Directory, phone string) (Contact, error)` | Returns the contact for the given phone. Returns error if not found. |
| `UpdateContact(dir Directory, phone string, c Contact) error` | Updates the contact at the given phone. Returns error if not found. |
| `DeleteContact(dir Directory, phone string) error` | Deletes the contact at the given phone. Returns error if not found. |
| `ListContacts(dir Directory) []Contact` | Returns all contacts sorted by name (case-insensitive). Returns nil if empty. |
| `SearchContacts(dir Directory, query string) []Contact` | Returns contacts whose name or email contains the query (case-insensitive). Returns nil if no matches. |

### Task 2 — Add verification

The starter `main()` has a placeholder — replace it with a complete directory workflow: add several contacts, list them, search, update, and delete.

### Task 3 — Prove correctness

Make all tests in `_starter/main_test.go` pass.

## Verification

After completing all tasks, run from `_starter/`:

```bash
go run .
```

Expected output (approximately):
```
Added Alice
Added Bob
Added Charlie
--- All contacts ---
Alice (alice@example.com) - 555-0101
Bob (bob@example.com) - 555-0102
Charlie (charlie@example.com) - 555-0103
--- Search for 'alice' ---
Alice (alice@example.com) - 555-0101
--- After deleting Bob ---
Alice (alice@example.com) - 555-0101
Charlie (charlie@example.com) - 555-0103
```

```bash
go test -v .
```

Expected output: all tests pass.

```bash
gofmt -l .
```

Expected output: (nothing — all files formatted)

```bash
go vet .
```

Expected output: (nothing — no issues)

## Rubric

| Criterion | Excellent (5) | Good (3) | Needs Work (0) |
|-----------|--------------|----------|----------------|
| **Correctness** | All CRUD operations work correctly; all tests pass | 3–4 operations work; most tests pass | Fewer than 3 operations work |
| **Data handling** | Proper error handling for missing, duplicate, and empty inputs | Basic error handling with gaps | Missing or incorrect error handling |
| **Proof** | Tests cover normal cases and edge cases (not found, duplicate, empty, case sensitivity) | Tests cover normal cases only | No meaningful tests |
| **Code quality** | Idiomatic Go, clean map/slice usage, proper error returns | Mostly readable with minor issues | Hard to read or unidiomatic |

## NEXT UP

Module 03 assessment — [Checkpoint](../../assessments/checkpoint/README.md)
