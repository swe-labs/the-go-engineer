# Input validation

## Learning objective

Validate HTTP request input using manual checks (required fields, type assertions, range constraints) and the `go-playground/validator` library, and return structured validation errors.

## Why this matters

Unvalidated input is the root cause of SQL injection, XSS, mass assignment, and logic bugs in web applications. Any data arriving from an HTTP request — headers, query strings, path parameters, and bodies — must be treated as hostile until validated. Go's explicit error handling makes validation verbose but also forces you to be deliberate. Using a validation library reduces boilerplate while keeping the explicit error model intact.

## Mental model

Think of validation as a sieve. Raw input goes in the top — strings, numbers, missing fields. Each sieve layer catches a specific kind of problem: missing required fields, wrong types, values out of range, format violations (email, UUID). Only what passes through all sieves is trusted and passed to business logic. Validation errors are collected and returned as a structured response so the client can fix all issues in one round trip.

## Core idea

Validation in Go follows a layered approach:

1. **Presence check**: required fields must be non-empty / non-zero.
2. **Type check**: string fields that should be numbers, booleans, dates.
3. **Range check**: numeric min/max, string length min/max.
4. **Format check**: email, URL, UUID, regex pattern.
5. **Business check**: unique username, sufficient balance, valid state transition.

For simple cases, manual checks with `if` statements are clear and dependency-free. For complex schemas, `go-playground/validator` provides struct tags:

```go
type CreateUserRequest struct {
    Name  string `json:"name"  validate:"required,min=2,max=100"`
    Email string `json:"email" validate:"required,email"`
    Age   int    `json:"age"   validate:"gte=18,lte=150"`
}
```

## Under the hood

`go-playground/validator` uses reflection to read struct tags, then applies registered validation functions. Each tag like `required`, `min=2`, `email` maps to a function that reads the field's value and returns an error. The library caches struct information after the first validation, so subsequent validations on the same type are fast. Manual validation with `if` statements compiles to direct comparisons — no reflection, maximum performance, zero dependencies.

## How Go uses it

```go
import "github.com/go-playground/validator/v10"

var validate = validator.New()

func createUser(w http.ResponseWriter, r *http.Request) {
    var req CreateUserRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, `{"error":"invalid JSON"}`, 400)
        return
    }
    if err := validate.Struct(req); err != nil {
        // Convert validation errors to structured response
        writeValidationError(w, err)
        return
    }
    // req is safe to use
}
```

## Go example

```go
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
)

type ProductInput struct {
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Category string  `json:"category"`
	Email    string  `json:"email"`
}

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

var validCategories = map[string]bool{
	"electronics": true,
	"clothing":    true,
	"food":        true,
}

var emailRegex = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)

func validateProduct(input ProductInput) []ValidationError {
	var errs []ValidationError

	if strings.TrimSpace(input.Name) == "" {
		errs = append(errs, ValidationError{Field: "name", Message: "name is required"})
	} else if len(input.Name) < 3 {
		errs = append(errs, ValidationError{Field: "name", Message: "name must be at least 3 characters"})
	}

	if input.Price <= 0 {
		errs = append(errs, ValidationError{Field: "price", Message: "price must be greater than 0"})
	} else if input.Price > 10000 {
		errs = append(errs, ValidationError{Field: "price", Message: "price must not exceed 10000"})
	}

	if !validCategories[input.Category] {
		errs = append(errs, ValidationError{Field: "category", Message: fmt.Sprintf("category must be one of: electronics, clothing, food")})
	}

	if input.Email != "" && !emailRegex.MatchString(input.Email) {
		errs = append(errs, ValidationError{Field: "email", Message: "invalid email format"})
	}

	return errs
}

func productHandler(w http.ResponseWriter, r *http.Request) {
	var input ProductInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid JSON"})
		return
	}

	errs := validateProduct(input)
	if len(errs) > 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(map[string]interface{}{"errors": errs})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(input)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /products", productHandler)
	log.Fatal(http.ListenAndServe(":8080", mux))
}
```

## Step-by-step execution

1. Client sends `POST /products` with body `{"name":"AB","price":-5,"category":"toys"}`.
2. JSON decoded into `ProductInput` → `Name: "AB"`, `Price: -5`, `Category: "toys"`.
3. `validateProduct` runs checks:
   - `Name` is not empty but length is 2 (< 3) → error: "name must be at least 3 characters".
   - `Price` <= 0 → error: "price must be greater than 0".
   - `Category` "toys" not in valid set → error: "category must be one of: electronics, clothing, food".
4. Three errors collected → returned as JSON array with 422 Unprocessable Entity.
5. Client fixes all three issues and resends.
6. All checks pass → product stored and 201 returned.

## Common mistakes

- Returning 400 Bad Request for all validation failures. Use 422 Unprocessable Entity for semantic validation errors (wrong values) and 400 for syntactic errors (malformed JSON).
- Stopping validation at the first error. Collect all errors so the client can fix everything at once.
- Not trimming whitespace. A name of `"   "` passes the `!= ""` check. Always trim user input before validation.
- Using reflection-based validators on performance-critical hot paths without caching. `go-playground/validator` caches struct info, but manual validation is still faster.
- Forgetting that zero values are often invalid: `int` field `0` may mean "not provided" or "valid value 0". Use pointers or `*int` to distinguish.

## Debugging walkthrough

A validator always rejects valid email addresses:

```go
emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)
```

**Symptom**: `"User@Example.com"` fails validation.

**Root cause**: The regex is case-sensitive (`[a-z]`), but email `User@Example.com` has uppercase letters. While the `local` part of an email is technically case-sensitive, in practice most systems treat it as case-insensitive.

**Fix**: Lowercase the email before validation, or make the regex case-insensitive with the `(?i)` flag:

```go
emailRegex := regexp.MustCompile(`(?i)^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)
```

Or use `strings.ToLower(input.Email)` before checking.

## Production notes

- Never trust client-side validation alone. Always validate server-side.
- Use a single validation layer, not scattered checks across handlers. Centralize rules in dedicated validation functions or a validator instance.
- Log validation failures at DEBUG level for monitoring bad actors, but never log the raw input (may contain PII or secrets).
- For large teams, define validation rules in a shared schema that both Go and frontend code can consume. Otherwise, front-end and backend validation inevitably diverge.

## Performance implications

- Manual validation with `if` statements is the fastest approach: no reflection, no allocations beyond the error slice.
- `go-playground/validator` uses reflection on first use of each struct type (~1-2µs overhead), then caches. Subsequent validations are ~200-500ns.
- Regex-based validation (email, phone) is expensive. Pre-compile regexes with `regexp.MustCompile` at package level. Avoid running regex on every keystroke in a search endpoint — debounce or use simpler checks for high-throughput paths.
- Collecting errors in a slice is O(n) in the number of fields. This is negligible for typical API payloads (10-50 fields).

## Practice task

Write a function `validateUser(input map[string]interface{}) []ValidationError` that checks:
- `"username"` is required, 3-30 characters, alphanumeric only.
- `"age"` is required, must be an integer between 13 and 150.
- `"email"` is optional but must be a valid email format if provided.
- `"role"` is optional, must be one of: "admin", "editor", "viewer".

Also write an HTTP handler that accepts POST requests, calls `validateUser`, and returns 422 with the errors array. Test with valid and invalid inputs.

## Tests / verification

```bash
go run ./curriculum/modules/08-http-rest-apis/lessons/06-input-validation
go test ./curriculum/modules/08-http-rest-apis/lessons/06-input-validation
```

## Review questions

1. What HTTP status code is most appropriate for validation errors? Why not 400?
2. Why should validation errors be collected all at once rather than failing fast on the first error?
3. When would you use manual validation (if statements) instead of a library like `go-playground/validator`?
4. How can you distinguish between a missing field and a zero-value field when decoding JSON?
5. What security vulnerabilities can arise from inadequate input validation?

## NEXT UP

Response writing — writing response bodies, setting status codes, encoding JSON, setting Content-Type, and streaming responses.
