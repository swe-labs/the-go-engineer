# TI.8 Custom Error Types

## Mission

- Define custom struct types that satisfy the built-in `error` interface.
- Enrich error instances with structured metadata for better troubleshooting.
- Utilize `errors.As` and `errors.Is` to safely inspect and handle specific error types.

## Prerequisites

- `TI.2` Methods
- `TI.3` Interfaces

## Mental Model

In Go, an **Error** is any type that implements the single-method `error` interface: `Error() string`. While simple string-based errors are useful for basic reporting, production systems often require structured data to differentiate between failure modes (e.g., a "Validation Error" vs. a "Connection Timeout"). By creating custom error structs, we can carry machine-readable context that allows calling code to make informed decisions about recovery or reporting.

## Visual Model

```mermaid
graph TD
    A[error Interface] --> B[Custom Error Struct]
    B --> C[Error() string]
```

## Machine View

When a custom struct is returned as an `error`, the runtime packs it into the standard 2-word interface structure:

1.  **ITab**: Points to the concrete type's method set (including the `Error()` method).
2.  **Data**: Points to the instance of the custom struct containing its metadata fields.

Functions like `errors.As(err, &target)` use reflection to verify if the dynamic type stored in the interface can be assigned to the target pointer, allowing you to "unwrap" the structured data safely.

## Run Instructions

```bash
go run ./04-types-design/8-custom-errors
```

## Code Walkthrough

- **Field Selection**: Add fields that provide actionable context (e.g., `StatusCode`, `RetryAfter`, `FieldID`).
- **Standard Satisfaction**: Always implement the `Error() string` method on the custom struct (or a pointer to it).
- **Type Checking**: Use `errors.As` instead of raw type assertions (`err.(MyError)`) to support wrapped errors in modern Go.

## Try It

1. In `main.go`, add a `TimeoutError` struct that includes a `Duration` field.
2. Implement the `Error()` method for `TimeoutError`.
3. Create a function that randomly returns a `TimeoutError` or a `ValidationError`.
4. Update the `handleError` function to detect and log the specific duration of a `TimeoutError` using `errors.As`.

## In Production

- **API Gateways**: Returning specific status codes and user-friendly messages for different failure modes.
- **Form Validation**: Carrying a map of field names to specific validation messages.
- **Database Drivers**: Distinguishing between constraint violations (which are client errors) and connection losses (which are system errors).

## Thinking Questions

1. Why is it generally safer to use `errors.As` than a direct type assertion on an error value?
2. How does structured error data improve the observability of a system compared to parsing error strings?
3. Should custom error types use value receivers or pointer receivers for their `Error()` method? Does it matter?

## Next Step

Next: `TI.9` -> [`04-types-design/9-generics`](../9-generics/README.md)
