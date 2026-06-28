# Checkpoint Questions — Functions, Errors, and Data Semantics

Answer each question with concrete reasoning.

## Question 1 (Recall)

What is the zero value of a function variable declared as `var fn func(int) int`?

A) `nil`  
B) A no-op function that returns 0  
C) `func(int) int {}`  
D) Compilation error — function variables must be initialized

## Question 2 (Recall)

In Go, are arguments passed by value or by reference?

A) Always by reference  
B) Always by value  
C) By value for primitives, by reference for structs  
D) It depends on whether the parameter type is a pointer

## Question 3 (Application)

Given `err := fmt.Errorf("connect: %w", errors.New("refused"))`, which call returns `true`?

A) `errors.Is(err, errors.New("refused"))`  
B) `errors.As(err, &target)`  
C) `errors.Unwrap(err) == errors.New("refused")`  
D) Both A and B

## Question 4 (Application)

```go
func doubleAll(nums []int) {
    for i, n := range nums {
        n *= 2
    }
}
```

After calling `doubleAll(slice)`, does `slice[0]` change? Why or why not?

A) Yes, because slices are references  
B) No, because `n` is a copy of each element  
C) Yes, because the loop variable modifies the original  
D) No, because the function does not return the slice

## Question 5 (Application)

```go
func getCount() (n int) {
    defer func() { n++ }()
    return 10
}
```

What does this function return?

A) 10  
B) 11  
C) 0  
D) It panics

## Question 6 (Application)

```go
type TempError struct{ Temp int }
func (e *TempError) Error() string { return fmt.Sprintf("temp %d", e.Temp) }

fn := func() error {
    return fmt.Errorf("sensor: %w", &TempError{Temp: 99})
}
```

How do you extract the `Temp` field from the returned error?

A) `fn().(*TempError).Temp`  
B) `errors.Unwrap(fn()).(*TempError).Temp`  
C) Declare `var te *TempError` and call `errors.As(fn(), &te)`, then read `te.Temp`  
D) Neither A, B, nor C works

## Question 7 (Recall)

Order these from outermost to innermost scope in a Go program:

1. Function scope  
2. Block scope  
3. Package scope  
4. Universe scope

A) 4, 3, 1, 2  
B) 3, 4, 2, 1  
C) 1, 2, 3, 4  
D) 4, 2, 3, 1

## Question 8 (Application)

A validation function calls `errors.Join` with three non-nil errors and one nil error. How many errors does `errors.Join` return?

A) 0  
B) 1  
C) 3  
D) 4

## Question 9 (Analysis)

```go
func readFile(path string) (data []byte, err error) {
    f, err := os.Open(path)
    if err != nil {
        return nil, fmt.Errorf("open: %w", err)
    }
    defer f.Close()
    data, err = io.ReadAll(f)
    if err != nil {
        return nil, fmt.Errorf("read: %w", err)
    }
    return data, nil
}
```

If `io.ReadAll` returns an error, will `f.Close()` still execute?

A) No, because the function returns early  
B) Yes, because `defer` runs before the function returns  
C) Only if `f.Close()` is also deferred  
D) No, because `data` is nil

## Question 10 (Analysis)

Consider two error-handling approaches:

Approach A:
```go
if err != nil {
    return fmt.Errorf("process order %s: %w", orderID, err)
}
```

Approach B:
```go
if err != nil {
    log.Printf("process order %s: %v", orderID, err)
    return err
}
```

Which approach preserves the full error chain for callers, and what trade-off does the other approach make?

A) Approach A preserves the chain; Approach B logs and re-returns the same error  
B) Approach A preserves the chain; Approach B logs but discards the chain when `%v` is used  
C) Approach B preserves the chain because it never wraps; Approach A breaks the chain  
D) Both preserve the chain equally
