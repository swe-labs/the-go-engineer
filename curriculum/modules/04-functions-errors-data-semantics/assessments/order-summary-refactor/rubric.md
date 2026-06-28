# Order Summary Refactor — Project Rubric

Grade each criterion as Excellent (E), Good (G), or Needs Work (N). A passing project requires at least 3/4 at Good or better.

## 1. Custom Error Types

| Level | Description |
|---|---|
| **Excellent** | Defines `ValidationError` and `OrderError` that implement `error`, `Unwrap`, and `Is` with meaningful fields. Uses them consistently. |
| **Good** | Defines at least one custom error type but missing `Unwrap` or `Is`, or fields are not well-chosen. |
| **Needs Work** | No custom error types; uses only `errors.New` or `fmt.Errorf` with `%v`. |

## 2. Error Wrapping and Matching

| Level | Description |
|---|---|
| **Excellent** | Uses `fmt.Errorf` with `%w` at every call site. Both `errors.Is` and `errors.As` are used correctly for sentinel and custom-type matching. |
| **Good** | Wrapping is present but inconsistent; either `errors.Is` or `errors.As` is used but not both. |
| **Needs Work** | Errors are returned as flat strings; no structured matching in tests or code. |

## 3. Defer and Cleanup

| Level | Description |
|---|---|
| **Excellent** | Uses `defer` for simulated resource cleanup. The `defer` runs on both success and panic paths. Cleanup is not duplicated. |
| **Good** | Uses `defer` but cleanup simulation is trivial or the cleanup logic is duplicated. |
| **Needs Work** | No `defer` usage; resources are not released. |

## 4. Test Quality

| Level | Description |
|---|---|
| **Excellent** | Table-driven tests cover success, each validation rule, sentinel matching via `errors.Is`, custom type extraction via `errors.As`, and multi-error scenarios. Subtests are used. |
| **Good** | Tests cover success and basic failure but miss sentinel matching or multi-error cases. |
| **Needs Work** | Tests are missing, incomplete, or use hard-coded values without table-driven patterns. |
