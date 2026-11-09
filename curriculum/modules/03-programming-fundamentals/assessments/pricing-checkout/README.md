# Project Assessment — Pricing Checkout

## Purpose

Evaluate your ability to build a complete pricing calculator using Go fundamentals: variables, types, control flow, slices, maps, and formatted output. This assessment verifies that you can implement real arithmetic logic with edge cases and prove correctness through automated tests.

## Targets

- `project-pricing-checkout`

## Scope

This assessment covers:

- Variable declaration and initialization
- Type conversions between integers and floats
- Constants and iota-based enumerations
- Control flow with if/else and switch
- Slice and map operations
- String formatting with `fmt.Sprintf`
- Integer cent arithmetic (avoiding floating-point currency bugs)
- Error handling for edge cases (empty slices, negative values)

## Evidence required

Submit:

1. Your completed `_starter/main.go` with all function bodies implemented
2. Output of `go test -v .` from the `_starter/` directory showing all tests passing
3. Output of `go run .` showing a correctly formatted checkout summary
4. Output of `gofmt -l .` (must be empty) and `go vet .` (must be clean)

## Time expectation

Plan for 60 to 90 minutes.

## Rubric

| Criterion | Weight | Excellent (5) | Good (3) | Needs Work (0) |
|-----------|--------|--------------|----------|----------------|
| **Correctness** | 30% | All 5 functions implemented correctly; all tests pass including edge cases (empty slice, zero quantity, unknown country) | 3–4 functions work; most tests pass | Fewer than 3 functions work |
| **Pricing logic** | 25% | Tax rates, quantity discounts, and cent formatting all produce exact expected values; no floating-point currency errors | Minor rounding drift in one area | Multiple pricing errors or float currency bugs |
| **Proof** | 25% | Tests cover normal cases and edge cases (empty items, negative cents, zero tax, boundary quantities at 10/50/100) | Tests cover normal happy-path cases only | No meaningful test coverage |
| **Code quality** | 20% | Idiomatic Go: proper error returns, clean switch/if logic, `fmt.Sprintf` formatting, no unused variables | Mostly readable with minor style issues | Unformatted, unidiomatic, or hard to follow |

**Passing score**: 70% (weighted total of 3.5/5 average across dimensions)

## Review questions

Be prepared to answer these questions during a review session:

1. Why does the project use integer cents instead of `float64` for currency? What bug would `float64` introduce for a value like `$0.10 + $0.20`?
2. Walk through your `CalculateTotal` function: what happens when the items slice is empty? What error do you return and why?
3. Your `DiscountPercent` function uses quantity thresholds. What happens at the boundary — does quantity 10 get the 5% discount or not? How did you verify this?
4. Explain the difference between `fmt.Sprintf("$%d.%02d", cents/100, cents%100)` and using `float64(cents)/100.0` for formatting. Which is safer and why?
5. What does `go vet` check for that the compiler does not? Give one example from your code.

## Passing standard

You pass when all tests pass, your checkout output matches the expected format, and you can explain your pricing logic without guessing about boundary cases or type conversion behavior.

## Retake policy

If you score below 70%, review the lessons corresponding to your rubric gaps (types, control flow, or error handling), fix the issues, and resubmit with a short correction note explaining what you changed and why.

## NEXT UP

[Project Assessment: Contact Directory](../contact-directory/README.md)
