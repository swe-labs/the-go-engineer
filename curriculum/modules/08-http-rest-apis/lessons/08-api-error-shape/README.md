# API error shape

## Learning objective

Design a consistent, structured JSON error response format with error codes and messages, and implement it across all API handlers. Understand and compare with RFC 7807 (Problem Details).

## Why this matters

Clients consuming your API need to handle errors programmatically. If every endpoint returns a different error shape, the client must special-case each one. A consistent error format means the client can write a single error handler: check the status code, parse the error body, display the message. This is the difference between a professional API that teams can integrate with in hours and one that requires weeks of trial-and-error debugging.

## Mental model

Think of error responses as a contract. Every error — whether validation failure, authentication rejection, or internal crash — uses the same envelope. The envelope has a status code (the HTTP method-level indicator), an error code (a machine-readable string for programmatic handling), and a message (a human-readable explanation). Optionally, it has details (field-level errors) and a trace ID for debugging. The client opens the envelope and knows exactly where to look, regardless of which endpoint produced the error.

## Core idea

A consistent JSON error response should include:

```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "The request contains invalid fields",
    "details": [
      {"field": "email", "message": "must be a valid email address"}
    ],
    "trace_id": "abc-123-def"
  }
}
```

Key principles:
- **Code**: machine-readable string like `"NOT_FOUND"`, `"VALIDATION_ERROR"`, `"UNAUTHORIZED"`. Clients use this for logic branching.
- **Message**: human-readable string suitable for display or logging.
- **Details**: optional array of per-field errors for validation failures.
- **Trace ID**: optional correlation ID for debugging across services.

RFC 7807 (Problem Details for HTTP APIs) standardizes this as:

```json
{
  "type": "https://api.example.com/errors/validation",
  "title": "Validation Error",
  "status": 422,
  "detail": "email must be a valid email address",
  "instance": "/api/users"
}
```

Go can implement either format. The key is consistency.

## Under the hood

The error response is just JSON written to `ResponseWriter`. What makes it "structured" is the Go type definition and the fact that every handler uses the same type. A shared `APIError` struct and a `writeError` helper function ensure consistency:

```go
type APIError struct {
    Code    string       `json:"code"`
    Message string       `json:"message"`
    Details []FieldError `json:"details,omitempty"`
}
```

Every handler calls `writeError(w, status, apiError)` instead of ad-hoc `json.Encode` calls. Middleware catches panics and unrecovered errors and also calls `writeError`. This guarantees every error response has the same shape.

## How Go uses it

```go
type FieldError struct {
    Field   string `json:"field"`
    Message string `json:"message"`
}

type APIError struct {
    Code    string       `json:"code"`
    Message string       `json:"message"`
    Details []FieldError `json:"details,omitempty"`
}

func writeError(w http.ResponseWriter, status int, err APIError) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(err)
}
```

## Go example

```go
package main

import (
	"encoding/json"
	"log"
	"net/http"
)

// APIError is the standard error envelope.
type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type APIError struct {
	Code    string       `json:"code"`
	Message string       `json:"message"`
	Details []FieldError `json:"details,omitempty"`
}

func writeError(w http.ResponseWriter, status int, err APIError) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(err)
}

var errNotFound = APIError{
	Code:    "NOT_FOUND",
	Message: "The requested resource was not found",
}

var errInternal = APIError{
	Code:    "INTERNAL_ERROR",
	Message: "An unexpected error occurred",
}

func getUserHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, APIError{
			Code:    "BAD_REQUEST",
			Message: "Missing user ID",
		})
		return
	}

	// Simulate not-found lookup.
	if id != "42" {
		writeError(w, http.StatusNotFound, errNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"id": "42", "name": "Alice"})
}

func createUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name  string `json:"name"`
		Email string `json:"email"`
		Age   int    `json:"age"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, APIError{
			Code:    "INVALID_JSON",
			Message: "The request body is not valid JSON",
		})
		return
	}

	var errs []FieldError
	if input.Name == "" {
		errs = append(errs, FieldError{Field: "name", Message: "Name is required"})
	}
	if input.Email == "" {
		errs = append(errs, FieldError{Field: "email", Message: "Email is required"})
	}
	if input.Age < 18 {
		errs = append(errs, FieldError{Field: "age", Message: "Must be at least 18"})
	}

	if len(errs) > 0 {
		writeError(w, http.StatusUnprocessableEntity, APIError{
			Code:    "VALIDATION_ERROR",
			Message: "The request contains invalid fields",
			Details: errs,
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(input)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /users/{id}", getUserHandler)
	mux.HandleFunc("POST /users", createUserHandler)
	log.Fatal(http.ListenAndServe(":8080", mux))
}
```

## Step-by-step execution

1. Client sends `GET /users/999`.
2. `getUserHandler` receives `id = "999"`.
3. `id != "42"` → calls `writeError(w, 404, errNotFound)`.
4. `writeError` sets `Content-Type: application/json`, calls `WriteHeader(404)`, encodes the `APIError` struct.
5. Client receives `{"code":"NOT_FOUND","message":"The requested resource was not found"}`.
6. Client sends `POST /users` with body `{"name":"","email":"","age":10}`.
7. `createUserHandler` decodes body, runs validation, collects three field errors.
8. Calls `writeError(w, 422, APIError{Code: "VALIDATION_ERROR", Details: [...]})`.
9. Client gets a 422 with a `details` array containing all three field errors.
10. Client can display each field error next to the corresponding input field.

## Common mistakes

- Returning error messages in different shapes across endpoints. One returns `{"error":"msg"}`, another returns `{"message":"msg"}`. This forces clients to check each endpoint individually.
- Returning 200 with an error in the body. If an error occurs, use the appropriate 4xx/5xx status code so the client can check the status code first.
- Exposing internal error details to the client (database errors, stack traces). Always log the full error server-side and return a sanitized message.
- Not including a trace/correlation ID. Without it, matching client-reported errors to server logs is manual and slow.
- Hardcoding error codes as strings everywhere. Define constants: `const ErrNotFound = "NOT_FOUND"`.

## Debugging walkthrough

A client reports that validation errors return HTML instead of JSON:

**Symptom**: When validation fails, the response has `Content-Type: text/plain` and contains a Go error string like `"json: cannot unmarshal string into Go value of type int"`.

**Root cause**: Somewhere in the handler chain, `http.Error(w, err.Error(), 400)` is called instead of the structured `writeError` function. The handler uses `http.Error` (which sets `Content-Type: text/plain`) as a fallback.

**Fix**: Ensure every code path that writes an error calls `writeError`, not `http.Error`. Use a middleware to catch panics and unexpected errors, routing them through the same structured error function.

## Production notes

- Define error code constants in a single package so they are consistent across handlers and can be documented.
- Use `trace_id` from the request context to correlate errors in distributed systems. Generate one in middleware if not provided by the client.
- For RFC 7807 compliance, set `Content-Type: application/problem+json` instead of `application/json`.
- Document every error code in your API reference. Clients need to know which error codes to expect from each endpoint.
- Consider a custom error type that implements `error` and carries an HTTP status code and `APIError`:

```go
type HTTPError struct {
    Status int
    APIError
}
func (e HTTPError) Error() string { return e.Message }
```

## Performance implications

- A consistent error format adds negligible overhead: a few struct fields and one JSON encode per error.
- Pre-declare common error structs (`errNotFound`, `errInternal`) at package level to avoid allocating them per request.
- Validation error details with many fields can be large (4 KB+). Consider truncating the details array if it exceeds a reasonable size (50 fields).

## Practice task

Define an `APIError` struct and `writeError` helper. Implement a handler `deleteResourceHandler` that:
- Accepts `DELETE /resources/{id}`.
- Returns `NOT_FOUND` if `id` is not `"42"`.
- Returns `BAD_REQUEST` if `id` is empty.
- Returns `INTERNAL_ERROR` if `id` is `"panic"` (simulate a panic, caught by middleware).
- On success, returns 204 No Content with no body.

Write table-driven tests for all four cases.

## Tests / verification

```bash
go run ./curriculum/modules/08-http-rest-apis/lessons/08-api-error-shape
go test ./curriculum/modules/08-http-rest-apis/lessons/08-api-error-shape
```

## Review questions

1. What are the minimum fields a structured error response should contain?
2. Why should error codes be strings like `"VALIDATION_ERROR"` rather than numeric codes?
3. What is RFC 7807 and what `Content-Type` does it specify?
4. Why should you avoid returning `err.Error()` or stack traces in production API responses?
5. What is the benefit of defining error code constants vs. using string literals?

## NEXT UP

Status codes — choosing the right HTTP status code for every response and avoiding common misclassification mistakes.
