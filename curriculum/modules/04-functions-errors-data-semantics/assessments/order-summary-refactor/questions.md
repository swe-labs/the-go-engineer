# Order Summary Refactor — Project Questions

Answer each question with concrete reasoning and evidence from your implementation.

## Question 1

Your `ValidationError` type includes a `Field` and `Rule`. Why is carrying structured fields better than returning a formatted string like `"field Name: required"`? What does this enable for the caller that a string cannot?

## Question 2

Your `OrderError` wraps an underlying error via `Unwrap()`. When would you choose a wrapping custom type (`OrderError` with `Err error`) over a sentinel error (`var ErrBadOrder = errors.New(...)`)? Give a concrete scenario where each is the better choice.

## Question 3

The refactored `processOrder` uses `defer res.Close()` for a simulated resource. What happens if `processOrder` panics partway through processing items — does the deferred `Close` still run? Explain using the semantics of `defer`.

## Question 4

In your tests, you use both `errors.Is` (to match sentinels) and `errors.As` (to extract custom types). If you could only use one of these two functions for testing errors, which would you choose and why? Consider both the `ValidationError` and `OrderError` types in your answer.

## Question 5

The starter code panics for all validation failures. Your refactored code returns `(result, error)` instead. What concrete advantages does returning errors provide over panicking? Name at least two, referencing how the caller can now handle failures.
