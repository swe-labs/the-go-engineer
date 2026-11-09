# Checkpoint Rubric

Grade each criterion as Excellent (E), Good (G), or Needs Work (N). A passing checkpoint requires at least 3/5 at Good or better.

## 1. Function Design

| Level | Description |
|---|---|
| **Excellent** | Understands function signatures, parameters, return values, multiple returns, and scope. Can reason about the call stack and pass-by-value behavior. |
| **Good** | Understands basic function syntax and returns but struggles with scope rules or pass-by-value implications. |
| **Needs Work** | Cannot write a correct function signature or confuses return with print. |

## 2. Error Handling Correctness

| Level | Description |
|---|---|
| **Excellent** | Uses errors-as-values idiom consistently. Chooses `fmt.Errorf` with `%w` for wrapping and selects `errors.Is` / `errors.As` appropriately. |
| **Good** | Returns errors from functions but does not wrap consistently or uses `%v` instead of `%w`. |
| **Needs Work** | Panics for expected failures or ignores returned errors. |

## 3. Error Type Implementation

| Level | Description |
|---|---|
| **Excellent** | Can design a custom error type that implements `error`, `Unwrap`, and `Is`. Understands how `errors.As` extracts values from the chain. |
| **Good** | Implements basic `Error() string` but omits `Unwrap` or `Is` methods. |
| **Needs Work** | Does not know how to define a custom error type or implement an interface. |

## 4. Testing Completeness

| Level | Description |
|---|---|
| **Excellent** | Writes table-driven tests covering success, failure, error matching, and edge cases. Uses subtests and error-checking helpers. |
| **Good** | Tests basic success/failure paths but does not test error matching or edge cases. |
| **Needs Work** | No tests, or tests only cover the happy path. |

## 5. Code Quality

| Level | Description |
|---|---|
| **Excellent** | Uses `defer` appropriately, avoids panics for expected errors, organizes code cleanly, follows Go idioms. |
| **Good** | Uses `defer` in basic cases but may panic for expected errors or have muddled organization. |
| **Needs Work** | No `defer` usage, panics for all failures, or monolithic functions with no separation of concerns. |
