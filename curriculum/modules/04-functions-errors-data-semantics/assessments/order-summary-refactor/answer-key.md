# Order Summary Refactor — Answer Key

## Question 1

Structured fields enable the caller to inspect *which* field failed and *why* programmatically, rather than parsing a string. For example, the caller can use `errors.As` to extract a `ValidationError` and then check `ve.Field` to decide whether to retry with different input or show a field-level error message in a UI. A formatted string would require fragile string parsing that breaks if the error message format changes.

## Question 2

Use a sentinel error (`ErrBadOrder`) when the caller only needs to know *whether* a specific failure occurred, regardless of context. For example, "order not found" is always the same condition.

Use a wrapping custom type (`OrderError`) when the error carries dynamic information (a code, a wrapped cause) that the caller needs to inspect. For example, `OrderError{Code: ErrCodeProcessing, Err: someNetworkErr}` lets the caller check the code and unwrap to find the root cause. The custom type also allows custom `Is()` matching on code alone, ignoring the wrapped message.

## Question 3

Yes, the deferred `res.Close()` still runs. In Go, `defer` functions execute when the enclosing function returns, *including* when it panics. The runtime unwinds the stack and runs all deferred calls in LIFO order before the panic propagates. So even if `processOrder` panics after processing three of five items, the simulated resource is still released.

## Question 4

`errors.As` is the more powerful choice. With `errors.As`, you can extract a `*ValidationError` and inspect its fields (`Field`, `Rule`), or extract an `*OrderError` and check its `Code`. You can also match sentinel errors indirectly by extracting a wrapper type that contains the sentinel and then comparing. `errors.Is` can only answer a yes/no question about sentinel membership and cannot extract structured data. However, in practice you should use both: `errors.Is` for simple sentinel checks and `errors.As` when you need to read fields from a custom type.

## Question 5

1. **Caller choice**: When errors are returned, the caller can decide how to handle each failure — retry, log, degrade, or abort. A panic forces immediate termination, giving the caller no control.

2. **Composability**: Returned errors can be wrapped with `fmt.Errorf` and `%w`, adding context at each layer. Tests can use `errors.Is` and `errors.As` to verify specific error conditions. Panics cannot be inspected or wrapped; recovering from a panic loses the structured information about what went wrong.
