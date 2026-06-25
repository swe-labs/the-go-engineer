# Validation

## Learning objective

Implement input validation patterns in Go, accumulate multiple validation errors using `errors.Join`, and design validation that returns actionable error information.

## Why this matters

A production API receives thousands of requests per minute. Every field must be checked before it reaches a database query, a file write, or a business logic operation. Invalid input that reaches the database causes SQL errors. Invalid input that reaches business logic causes panics or data corruption. Validation is the first and most important defensive layer in any Go service. Go's explicit error handling makes validation patterns clear and testable.

## Mental model

Validation is a gatekeeper at the system boundary. Raw input arrives unchecked; the gatekeeper inspects each field against a set of rules. If all rules pass, the input is allowed through. If any rule fails, the gatekeeper returns a report of what failed and why. The report is a structured error (or multiple errors) that the caller can inspect to tell the user exactly what to fix.

## Core idea

Validation in Go follows a simple contract: a function takes input and returns an error. The error is `nil` when the input is valid. The error is non-nil when the input is invalid, and the error describes what is wrong.

Key patterns:

1. **Single-field validation**: A function checks one field and returns a descriptive error on failure.
2. **Struct-level validation**: A function validates all fields of a struct and collects errors.
3. **Accumulating errors**: Multiple validation failures are collected (not short-circuited) using `errors.Join`.
4. **Validation structs**: Custom types carry structured failure data (which field, what rule, what value).

## Under the hood

`errors.Join` (added in Go 1.20) returns an error that wraps all provided errors:

```go
func Join(errs ...error) error
```

It returns nil if all errs are nil. The returned error implements `Unwrap() []error` which returns the slice. This means `errors.Is` and `errors.As` search across all joined errors, not just a linear chain.

Before `errors.Join`, developers used the `hashicorp/go-multierror` package or custom error slices. `errors.Join` is now the standard library solution.

## How Go uses it

Validation appears at every input boundary:

```go
func validateUser(u User) error {
    var errs []error
    if u.Name == "" {
        errs = append(errs, fmt.Errorf("name is required"))
    }
    if u.Age < 0 || u.Age > 150 {
        errs = append(errs, fmt.Errorf("age must be between 0 and 150, got %d", u.Age))
    }
    if u.Email != "" && !strings.Contains(u.Email, "@") {
        errs = append(errs, fmt.Errorf("email is invalid: %q", u.Email))
    }
    return errors.Join(errs...)
}
```

Production validation is often more structured, using a validation struct with fields:

```go
type FieldError struct {
    Field string
    Value interface{}
    Rule  string
}

func (f *FieldError) Error() string {
    return fmt.Sprintf("%s: rule %q violated (value=%v)", f.Field, f.Rule, f.Value)
}
```

This lets callers use `errors.As` to extract `*FieldError` and render field-level error messages in a UI or API response.

## Go example

```go
package main

import (
	"errors"
	"fmt"
	"strings"
)

type FieldError struct {
	Field string
	Value interface{}
	Rule  string
}

func (f *FieldError) Error() string {
	return fmt.Sprintf("%s: rule %q violated (value=%v)", f.Field, f.Rule, f.Value)
}

type User struct {
	Name  string
	Age   int
	Email string
}

func validateUser(u User) error {
	var errs []error
	if strings.TrimSpace(u.Name) == "" {
		errs = append(errs, &FieldError{Field: "Name", Value: u.Name, Rule: "required"})
	}
	if u.Age < 0 || u.Age > 150 {
		errs = append(errs, &FieldError{Field: "Age", Value: u.Age, Rule: "range:0-150"})
	}
	if u.Email == "" {
		errs = append(errs, &FieldError{Field: "Email", Value: u.Email, Rule: "required"})
	} else if !strings.Contains(u.Email, "@") {
		errs = append(errs, &FieldError{Field: "Email", Value: u.Email, Rule: "valid-email"})
	}
	return errors.Join(errs...)
}

func main() {
	users := []User{
		{Name: "", Age: -1, Email: ""},
		{Name: "Alice", Age: 30, Email: "alice@example.com"},
		{Name: "Bob", Age: 200, Email: "bob"},
	}
	for _, u := range users {
		err := validateUser(u)
		if err == nil {
			fmt.Printf("User %q: valid\n", u.Name)
			continue
		}
		fmt.Printf("User %q: %v\n", u.Name, err)
		var fe *FieldError
		for i := 0; errors.As(err, &fe); i++ {
			if i > 0 {
				fmt.Print(", ")
			}
			fmt.Printf("  [%d] field=%q rule=%q", i, fe.Field, fe.Rule)
			// Reset fe for next iteration
			fe = nil
		}
		fmt.Println()
	}
}
```

## Step-by-step execution

For `validateUser(User{Name: "", Age: -1, Email: ""})`:

1. `Name` is empty → append `&FieldError{Field:"Name", Rule:"required"}`.
2. `Age` is -1 (< 0) → append `&FieldError{Field:"Age", Rule:"range:0-150"}`.
3. `Email` is empty → append `&FieldError{Field:"Email", Rule:"required"}`.
4. `errors.Join(errs...)` creates a joined error containing all three.
5. `err.Error()` prints all three errors joined by `"\n"`.
6. `errors.As(err, &fe)` walks the joined errors via `Unwrap() []error` and finds the first `*FieldError`.

## Common mistakes

- **Short-circuiting instead of accumulating**: Returning on the first validation error (`return err` inside each if block) hides other problems. The user must fix one field, resubmit, hit the next error, fix it, resubmit, etc. Always accumulate when possible.

- **Returning generic errors**: `return errors.New("invalid input")` tells the user nothing actionable. Include which field failed, what value was provided, and what rule was violated.

- **Not trimming whitespace**: `" "` is not empty but is visually indistinguishable from empty in a form. Trim before validating.

- **Performing side effects before validation**: If a function mutates state and then validates, a validation failure leaves state partially mutated. Validate first, then act.

## Debugging walkthrough

```go
package main

import "fmt"

func processOrder(price float64, qty int) error {
	if price < 0 {
		return fmt.Errorf("negative price: %.2f", price)
	}
	if qty <= 0 {
		return fmt.Errorf("non-positive quantity: %d", qty)
	}
	// ... process order
	return nil
}

func main() {
	err := processOrder(-5.0, -1)
	fmt.Println(err)
}
```

**Symptom**: Only prints `"negative price: -5.00"`. The user fixes the price, resubmits, and then gets `"non-positive quantity: -1"`.

**Root cause**: Short-circuit validation: the quantity error is never reported because the function returns on the first error.

**Fix**: Collect all errors and return them together:

```go
func processOrder(price float64, qty int) error {
    var errs []error
    if price < 0 {
        errs = append(errs, fmt.Errorf("negative price: %.2f", price))
    }
    if qty <= 0 {
        errs = append(errs, fmt.Errorf("non-positive quantity: %d", qty))
    }
    return errors.Join(errs...)
}
```

## Production notes

- **Validate at the boundary**: Validate input as soon as it enters your system (HTTP handler, CLI flag parser, message queue consumer). Never let invalid data propagate inward.

- **Use structured validation errors**: In APIs that return JSON, map validation errors to a structured response like `{"errors": [{"field": "name", "message": "required"}]}`.

- **Validate in the domain layer too**: Do not rely solely on HTTP-layer validation. Business logic often has rules that HTTP validation doesn't know about (e.g., "order total cannot exceed credit limit").

- **Third-party libraries**: `go-playground/validator` uses struct tags for declarative validation. It is widely used but adds a dependency. For many projects, manual validation functions are simpler and more explicit.

- **errors.Join is shallow**: It only joins one level. If you nest `errors.Join` inside another `errors.Join`, the structure flattens. This is usually fine but be aware when building complex error trees.

## Performance implications

- `errors.Join` allocates a new error struct. For hot validation paths (thousands per second), pre-allocate a slice with capacity.
- String formatting in validation errors allocates. In cold paths this doesn't matter. In hot paths, consider using error types that construct the message lazily.
- Validating fields is CPU-bound and fast. The bottleneck is usually I/O (reading the request body), not validation logic.

## Practice task

Write a function `validateOrder(o Order) error` where `Order` has fields `ProductID string`, `Quantity int`, and `PriceCents int64`. Rules:

- `ProductID` must be non-empty and at most 32 characters.
- `Quantity` must be between 1 and 100.
- `PriceCents` must be positive.

Use `errors.Join` to accumulate all errors. Define a `ValidationError` struct with `Field`, `Value`, and `Rule` fields. In `main`, test with a fully invalid order and print all errors.

## Tests / verification

```bash
go run ./curriculum/modules/04-functions-errors-data-semantics/lessons/13-validation
go test ./curriculum/modules/04-functions-errors-data-semantics/lessons/13-validation
```

## Review questions

1. Why should validation accumulate errors instead of returning on the first failure?
2. How does `errors.Join` differ from `fmt.Errorf` with `%w` for combining multiple errors?
3. What does `errors.Unwrap` return for an error created by `errors.Join`?
4. Describe a scenario where you would use a custom validation struct over a simple `error` string.
5. Where should validation happen in a layered Go application?

## NEXT UP

Orchestration
