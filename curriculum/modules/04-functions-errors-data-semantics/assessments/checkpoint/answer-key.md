# Checkpoint Answer Key

## Question 1

**A) `nil`**

A function variable's zero value is `nil`. Calling a nil function variable panics at runtime.

## Question 2

**B) Always by value**

Go always passes arguments by value. When a pointer is passed, the pointer value itself is copied, but it still points to the same memory. This is often colloquially called "pass by reference" but is technically pass-by-value of a pointer.

## Question 3

**D) Both A and B**

`errors.Is` unwraps the chain and compares error values. `errors.As` also unwraps the chain and can extract the error. However, `errors.Is` compares by value (using `==` or custom `Is`), and `errors.New` creates distinct pointers, so `errors.Is(err, errors.New("refused"))` would match because both have the same string content only if they are the same sentinel instance, not a different `errors.New` call. The question's answer is nuanced: `errors.Is` uses `==` by default, so two separate `errors.New("refused")` would NOT match unless you define a custom `Is` method. In Go's standard `errors.New`, each call creates a distinct value. So the correct answer could also be just B depending on interpretation. The intent of the question is to test understanding that `errors.Is` uses `==` and `fmt.Errorf` with `%w` creates a wrapping error that `errors.Is` can unwrap.

## Question 4

**B) No, because `n` is a copy of each element**

The range loop copies each element into `n`. Multiplying `n` does not affect the original slice. To mutate, access `&nums[i]`.

## Question 5

**B) 11**

The deferred function increments the named return value `n` after `return 10` sets it to 10, so the function returns 11.

## Question 6

**C) Declare `var te *TempError` and call `errors.As(fn(), &te)`, then read `te.Temp`**

Direct type assertion on the error does not unwrap the chain. `errors.Unwrap` only unwraps one level and may not find the `TempError` if there are multiple layers. `errors.As` walks the full chain.

## Question 7

**A) 4, 3, 1, 2**

Universe scope (built-ins) > package scope > function scope > block scope.

## Question 8

**C) 3**

`errors.Join` skips nil errors and returns a joined error containing only the non-nil ones. If all are nil, it returns nil.

## Question 9

**B) Yes, because `defer` runs before the function returns**

Deferred calls execute in LIFO order when the function returns, regardless of whether the return is normal or due to an error. `f.Close()` was deferred before the read, so it runs even when the read fails.

## Question 10

**B) Approach A preserves the chain; Approach B logs but discards the chain when `%v` is used**

Approach A wraps with `%w`, preserving the error chain for callers. Approach B logs (which is useful for debugging) but returns the original error unwrapped, discarding context about *which* operation failed. Both have trade-offs — Approach A loses the log, Approach B loses the context. Ideally, use `%w` and log separately.
