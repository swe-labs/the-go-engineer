# Section 4: Functions & Errors

## Beginner → Expert Mapping

| Topic | Level | Importance | Engineering Concept |
|-------|-------|------------|---------------------|
| Functions | Beginner | High | First-class citizens, multi-return |
| Variadics | Beginner | Medium | `...` spreading arguments |
| Closures | Intermediate | High | State retention, anonymous functions |
| Errors | Intermediate | **Critical** | Sentinel errors, idiomatic explicit checking |
| Defer & Panic | Advanced | High | LIFO execution, stack unwinding |

## Engineering Depth
Go handles errors as pure values rather than throwing exceptions. This forces the engineer to trace the exact control flow. You will also learn the mechanics of `defer`, which executes via LIFO (Last-In-First-Out) stack popping, making it ideal for resource cleanup (mutex unlocks, file closures, DB tx rolls).

## References
1. **[Official Blog]** [Error handling and Go](https://go.dev/blog/error-handling-and-go)
2. **[Official Blog]** [Defer, Panic, and Recover](https://go.dev/blog/defer-panic-and-recover)

---

## 🏗 Exercise: Custom Error Handling App (`8-error-handling`)

This project will force you to implement custom error structs and use `errors.Is`/`errors.As`.

### Step-by-Step Instructions & Hints
1. **Define the Error State:** Create a struct `MathError` containing details about a failed operation (Inputs, Time, Message).
2. **Implement the `error` Interface:** Add the `.Error() string` method to your custom struct.
   - *Hint:* Use a pointer receiver `*MathError` if the struct is large, to avoid copying it across the stack.
3. **Create the Failing Function:** Write a `Divide(a, b int) (int, error)` function.
   - *Hint:* If `b == 0`, return your custom `MathError`.
4. **Inspect the Error:** In `main()`, check the returned error.
   - *Hint:* **Do not** use `if err != nil { string matching }`. This is an ANTI-PATTERN.
   - *Hint:* Use `errors.As(err, &myMathErr)` to extract the struct fields safely and idiomatically.


## Learning Path

| ID | Lesson | Concept | Requires |
| --- | --- | --- | --- |
| FE.1 | [functions](./1-function) | Parameters · return types · pass-by-value | 🟢 entry |
| FE.2 | [closures &amp; recursion](./2-function-2) | Captured variables · anonymous functions · stack frames | FE.1 |
| FE.3 | [variadic functions](./3-variadic-func) | ...T syntax · slice internally · spread operator | FE.1, FE.2 |
| FE.4 | [multiple return values](./4-function-multi-values) | (value, error) convention · named returns · naked return | FE.1, FE.3 |
| FE.5 | [custom errors](./5-custom-error) | error interface · sentinel errors · errors.Is | FE.4 |
| FE.6 | [error wrapping](./5b-error-wrapping) | fmt.Errorf %w · errors.As · errors.Join | FE.5 |
| FE.7 | [defer](./6-defer) | LIFO execution · argument evaluation timing · resource cleanup | FE.1, FE.5 |
| FE.8 | [panic &amp; recover](./7-panic-recover) | Stack unwinding · recover inside defer only | FE.7 |
