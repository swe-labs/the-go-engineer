# FE.7 Order Summary

## Mission

Build one small order-summary program that proves the learner can combine functions, validation,
multiple return values, and explicit errors in one readable flow.

This is the foundations milestone for the section.

## Prerequisites

Complete these first:

- `FE.1` functions basics
- `FE.2` parameters and returns
- `FE.3` multiple return values
- `FE.4` errors as values
- `FE.5` validation
- `FE.6` orchestration

## What You Will Build

Implement a small order-summary program that:

1. validates an order name
2. validates a slice of item prices
3. validates a shipping value
4. sums the prices with one helper
5. orchestrates the whole flow with `processOrder`
6. returns either a readable summary or an explicit error

## Visual Model

```text
order name ---+
prices -------+--> processOrder(...)
shipping -----+
```

```text
processOrder
   |
   +--> validateOrderName
   +--> validatePrices
   +--> validateShipping
   +--> sumPrices
   +--> buildSummary
   +--> return (summary, nil)
```

```text
failure path:
one validation fails
processOrder returns ("", error)
```

## Machine View

This milestone is about function coordination more than low-level memory detail.

The important machine truth is:

- `main()` calls one top-level function
- `processOrder` calls helpers in order
- each helper either returns success or returns an error
- the first error stops the flow
- only the success path reaches the summary builder

That is how a readable application flow grows out of smaller functions.

## Files

- [main.go](./main.go): complete runnable solution
- [_starter/main.go](./_starter/main.go): starter scaffold with TODOs
- [main_test.go](./main_test.go): automated proof surface for the milestone

## Run Instructions

Run the completed solution:

```bash
go run ./01-foundations/05-functions-and-errors/7-order-summary
```

Run the starter scaffold:

```bash
go run ./01-foundations/05-functions-and-errors/7-order-summary/_starter
```

Run the automated verification surface:

```bash
go test ./01-foundations/05-functions-and-errors/7-order-summary
```

## Recommended Learning Flow

1. Read this README first.
2. Open [_starter/main.go](./_starter/main.go) and list the functions you need.
3. Build the validators first.
4. Add `sumPrices`, `buildSummary`, and `processOrder` after that.
5. Compare with [main.go](./main.go) only after your own attempt.

## Code Walkthrough

### `func validateOrderName(name string) error {`

This first validator checks whether the order name is meaningful enough to continue.

### `func validatePrices(prices []int) error {`

This validator rejects:

- an empty slice
- any negative price

That keeps the later math honest.

### `func validateShipping(shipping int) error {`

This validator is intentionally small.
It teaches that not every helper needs to be complex to be useful.

### `func sumPrices(prices []int) int {`

This helper does one job only:
add all prices together.

### `func buildSummary(name string, subtotal int, shipping int) string {`

This helper keeps formatting separate from validation and calculation.

### `func processOrder(name string, prices []int, shipping int) (string, error) {`

This is the milestone's most important function.
It owns the sequence of the work:

1. validate the name
2. validate the prices
3. validate shipping
4. calculate subtotal
5. build the final summary

### `return "", err`

Whenever validation fails, the function returns immediately with:

- an empty summary
- the real error

That is the foundations version of honest failure handling.

### `return buildSummary(name, subtotal, shipping), nil`

This success line makes the final contract clear:

- summary on success
- `nil` error on success

## Try It

1. Add one more price to the success case and observe the total.
2. Change shipping to a negative number and follow the failure path.
3. Replace the order name with only spaces and confirm the name validation still fails.

## Success Criteria

Your finished solution should:

- reject empty or invalid input before calculating
- keep validation, calculation, and formatting in separate helpers
- return explicit errors instead of hiding failure
- keep `main()` readable by delegating to `processOrder`
- pass the supplied tests

## Common Failure Modes

- trying to calculate totals before validating the inputs
- returning a misleading summary when an error should stop the flow
- mixing formatting and validation into one giant function
- forgetting to check `err` before using the returned summary

## Verification Surface

Use these three proof surfaces together:

1. `go run ./01-foundations/05-functions-and-errors/7-order-summary`
2. `go run ./01-foundations/05-functions-and-errors/7-order-summary/_starter`
3. `go test ./01-foundations/05-functions-and-errors/7-order-summary`

## Next Step

After this milestone, move to
[2 Types and Design](../../../docs/stages/02-types-and-design.md).
