# Engineering Error Framework

This document defines the three-tier error model used in backend, production, and flagship lessons.

## The Gap

Basic Go lessons teach:

```go
func doSomething() error {
	return fmt.Errorf("something failed")
}
```

Production-shaped Go needs more judgment:

- when to return vs wrap vs panic
- user error vs system error vs fatal error
- what to log
- what to show the caller
- what can be retried

## The Framework

### 1. User Errors

**Definition:** input validation failures and business rule violations.

**Handle:** return to the caller. Do not log as infrastructure failure.

Examples:

- invalid email format
- insufficient balance
- item not found
- permission denied

```go
type UserError struct {
	Code    string
	Message string
}

func (e *UserError) Error() string {
	return e.Message
}
```

Use stable uppercase machine codes for public or production-shaped errors:

```go
return nil, &UserError{
	Code:    "INVALID_EMAIL",
	Message: "email format is incorrect",
}
```

### 2. System Errors

**Definition:** infrastructure failures, external service failures, and unexpected operational failures.

**Handle:** wrap with context and propagate.

Examples:

- database query failed
- HTTP request timeout
- file read error
- external API failure

```go
type SystemError struct {
	Code    string
	Message string
	Cause   error
}

func (e *SystemError) Error() string {
	return fmt.Sprintf("%s: %v", e.Message, e.Cause)
}

func (e *SystemError) Unwrap() error {
	return e.Cause
}
```

Example:

```go
return nil, &SystemError{
	Code:    "DB_QUERY_FAILED",
	Message: "query user by id",
	Cause:   err,
}
```

### 3. Fatal Errors

**Definition:** startup failures, unrecoverable states, and invariant violations.

**Handle:** log and exit from the top-level coordinator, usually `main`.

Examples:

- invalid required configuration at startup
- database cannot connect during startup
- migration failure
- missing required resource

Do not call `log.Fatal` from worker goroutines.

## Decision Tree

```text
Is this caused by user input or business rules?
  yes -> UserError
  no  -> Is this unrecoverable at startup or an invariant violation?
          yes -> FatalError
          no  -> SystemError
```

## Handler Layer Pattern

```go
func (h *Handler) HandleError(w http.ResponseWriter, err error) {
	var userErr *UserError

	switch {
	case errors.As(err, &userErr):
		http.Error(w, userErr.Message, http.StatusBadRequest)

	case errors.Is(err, context.DeadlineExceeded):
		http.Error(w, "request timeout", http.StatusGatewayTimeout)

	default:
		h.logger.Error("request failed", "error", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
	}
}
```

## Repository Layer Rules

Repository code should:

- accept `context.Context`
- use parameterized queries
- close rows
- check `rows.Err()`
- wrap system failures
- translate not-found cases intentionally
- never leak raw database details into public responses

## Anti-Patterns

### Ignoring errors

```go
result, _ := doImportantWork()
```

### Returning raw external errors from boundaries

```go
return user, err
```

### Panic for control flow

```go
panic(err)
```

### Logging secrets or PII

```go
logger.Info("login", "password", password)
```

## Summary

| Error Type | When | Handling | Typical HTTP |
| --- | --- | --- | --- |
| UserError | validation, business rules | return | 400, 403, 404 |
| SystemError | infrastructure, external failures | wrap and propagate | 500, 502, 503, 504 |
| FatalError | startup or unrecoverable failure | log and exit at top level | not applicable |
