# Pricing Checkout

## Learning objective

Build a pricing calculator and checkout system using Go fundamentals: variables, control flow, structs, slices, and maps.

## Why this matters

Real e-commerce systems process pricing with tax, discounts, and currency formatting. This project consolidates every Module 03 concept into one practical tool you could show in a portfolio.

## Prerequisites

Complete all Module 03 lessons (01–23).

## Available concepts

This project draws on values, expressions, variables, types, zero values, type conversions, constants, iota, boolean logic, if/else, switch, loops, range, arrays, slices, slice length and capacity, slice sharing and aliasing, maps, comma-ok idiom, strings, bytes, runes and unicode, pointers, and beginner pointer mistakes.

## Project overview

The `_starter/` directory contains a Go program with missing and broken function bodies. Your job is to implement each function correctly so that:

- The program compiles and runs without errors.
- All tests pass.
- `gofmt -l .` produces no output.
- `go vet .` produces no output.

You may check `_solution/` only after you have attempted all fixes yourself.

## Tasks

### Task 1 — Implement pricing logic

Complete the following functions in `_starter/main.go`:

| Function | Behaviour |
|----------|-----------|
| `NewItem(name string, priceCents int) Item` | Creates an `Item` with the given name and price in cents. |
| `TaxRate(country string) float64` | Returns tax rate: `"US"` → 0.07, `"GB"` → 0.20, `"DE"` → 0.19, `"JP"` → 0.10. Unknown country returns 0.0. |
| `DiscountPercent(quantity int) float64` | Quantity ≥ 100 → 0.15, ≥ 50 → 0.10, ≥ 10 → 0.05, otherwise 0.0. |
| `FormatPrice(cents int) string` | Converts cents to a dollar string: `1234` → `"$12.34"`. Negative cents returns `"-$12.34"`. |
| `CalculateTotal(items []Item, country string) (int, error)` | Sums item prices, applies tax, applies quantity discount, returns total in cents. Returns error for empty slice. |

### Task 2 — Add verification

The program should print formatted output for a sample checkout. The starter `main()` currently has a placeholder — replace it with a complete checkout flow.

### Task 3 — Prove correctness

The test file `_starter/main_test.go` contains tests for each function. Make all tests pass.

## Verification

After completing all tasks, run from `_starter/`:

```bash
go run .
```

Expected output (approximately):
```
Item: Laptop - $999.99
Item: Mouse - $24.99
Item: Keyboard - $89.99
Subtotal: $1114.97
Discount (5%): $55.75
Tax (US 7%): $74.15
Total: $1133.37
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
| **Correctness** | All functions correctly implemented; tests pass | 3–4 functions work correctly; most tests pass | Fewer than 3 functions work |
| **Pricing logic** | Tax rates, discounts, and formatting all correct | Minor errors in one area | Multiple pricing errors |
| **Proof** | Tests cover normal cases and edge cases (empty, negative, zero) | Tests cover normal cases only | No meaningful tests |
| **Code quality** | Idiomatic Go, proper error handling, clean formatting | Mostly readable with minor issues | Hard to read or unformatted |

## NEXT UP

[Contact Directory](../contact-directory/README.md)
