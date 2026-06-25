# Comma-ok idiom

## Learning objective

Use the comma-ok (comma-ok) two-value form to safely access map entries, perform type assertions, and receive from channels — and explain why omitting the `ok` check is a common source of bugs.

## Why this matters

Go does not have exceptions, and it avoids sentinel values (like `-1` or `NULL`) wherever possible. Instead, it returns a second boolean value to indicate success. This pattern appears in three fundamental operations: map access, type assertions, and channel receives. Mastering the comma-ok idiom means you handle missing keys gracefully, convert interface types safely, and drain channels without panic — skills required in every Go codebase of non-trivial size.

## Mental model

Think of the comma-ok idiom as a **guarded access** pattern. The first value is the thing you want (or its zero value). The second value (`ok`) is a boolean that answers "Did the operation succeed?" You should always check `ok` unless you are certain the operation cannot fail, and even then, checking is cheap insurance.

```go
// Without ok: you get the zero value, but is it real or a default?
v := m["key"]

// With ok: you know for sure
v, ok := m["key"]   // ok == true  → key exists
                     // ok == false → key not present
```

## Core idea

The comma-ok idiom is the two-value return form used in three distinct Go constructs:

### 1. Map access

```go
val, ok := myMap[key]
```

- If `key` exists: `val` is the stored value, `ok` is `true`.
- If `key` does not exist: `val` is the zero value for `V`, `ok` is `false`.

### 2. Type assertion

```go
val, ok := someInterface.(ConcreteType)
```

- If the interface holds a value of `ConcreteType`: `val` is the concrete value, `ok` is `true`.
- If not: `val` is the zero value of `ConcreteType`, `ok` is `false`.
- Without `ok`, a failed type assertion **panics**.

### 3. Channel receive

```go
val, ok := <-ch
```

- If the channel is open and a value is received: `val` is the value, `ok` is `true`.
- If the channel is closed and empty: `val` is the zero value, `ok` is `false`.
- Without `ok`, you cannot distinguish a zero value from a closed channel.

## Under the hood

The comma-ok pattern is a **language-level feature**, not a library convention. The compiler recognizes these three contexts and generates different code for the single-value vs two-value form:

- **Map access**: The single-value form calls `runtime.mapaccess1`; the two-value form calls `runtime.mapaccess2`. The difference is internal: `mapaccess2` returns an additional boolean indicating whether the key's hash bucket contained a matching key.

- **Type assertion**: The single-value form (`x.(T)`) generates a `runtime.assertI2I` call that panics on failure. The two-value form generates `runtime.assertI2I2` which returns a boolean instead of panicking.

- **Channel receive**: The single-value form `<-ch` blocks forever on a nil channel and returns zero values on a closed channel indistinguishably. The two-value form returns `ok=false` when the channel is closed and drained.

The `ok` variable is a regular `bool`. There is nothing special about the name `ok` — it is a convention. You could name it `found`, `present`, `exists`, or any valid identifier.

## How Go uses it

The comma-ok idiom appears throughout idiomatic Go:

- **Safely reading caches**: `if val, ok := cache[key]; ok { ... }`
- **Type switches**: the `switch v := x.(type)` form implicitly uses the comma-ok pattern for each case.
- **Draining a channel**: `for v := range ch` is syntactic sugar over `for v, ok := <-ch; ok; v, ok = <-ch { ... }`
- **Asserting to an interface**: checking if a value implements another interface: `if w, ok := r.Writer.(http.Flusher); ok { ... }`

The standard library itself uses this pattern. For example, reading from a `net.Conn` that also implements `http.Hijacker`:

```go
if hj, ok := w.(http.Hijacker); ok {
    hj.Hijack()
}
```

## Go example

```go
package main

import "fmt"

func describe(i interface{}) {
	fmt.Printf("(%v, %T) -> ", i, i)

	v, ok := i.(string)
	if ok {
		fmt.Println("is a string:", v)
		return
	}

	w, ok := i.(int)
	if ok {
		fmt.Println("is an int:", w)
		return
	}

	fmt.Println("unknown type")
}

func main() {
	scores := map[string]int{
		"alice": 92,
		"bob":   85,
	}

	if score, ok := scores["alice"]; ok {
		fmt.Println("Alice score:", score)
	}

	if _, ok := scores["eve"]; !ok {
		fmt.Println("Eve not found")
	}

	describe("hello")
	describe(42)
	describe(3.14)

	ch := make(chan int, 2)
	ch <- 10
	ch <- 20
	close(ch)

	for v, ok := <-ch; ok; v, ok = <-ch {
		fmt.Println("Received:", v)
	}

	v, ok := <-ch
	fmt.Printf("After close: v=%d, ok=%v\n", v, ok)
}
```

## Step-by-step execution

For the map lookup `if score, ok := scores["alice"]; ok { ... }`:

1. Compiler sees the two-value assignment form.
2. It calls `runtime.mapaccess2(scores, "alice")`.
3. `mapaccess2` hashes `"alice"`, finds bucket, compares keys.
4. Key found: returns `(92, true)`.
5. `ok` is `true`, so the `if` body executes.

For the type assertion `v, ok := i.(string)` where `i` is `42` (an `int`):

1. Compiler generates `runtime.assertI2I2` call.
2. Runtime checks the dynamic type of `i` against `string`.
3. Types don't match: returns `("", false)`.
4. `ok` is `false`, body does not execute. No panic.

For the channel receive `v, ok := <-ch` on a closed channel:

1. Go's channel runtime checks the `closed` flag and the buffer.
2. If buffer is empty: returns `(0, false)`.
3. The `ok: false` signals the channel is drained.

## Common mistakes

- **Ignoring `ok` on type assertions**: `str := i.(string)` panics if `i` is not a string. Always use `str, ok := i.(string)` and handle the `false` case.
- **Not distinguishing zero from missing**: `if val := m["key"]; val == 0` conflates "key not present" with "value is zero". Use `val, ok := m["key"]`.
- **Using single-value channel receive in a loop**: `for { v := <-ch }` cannot detect when the channel is closed. Use `for v := range ch` or the two-value form.
- **Shadowing `ok` in nested scopes**: each `if _, ok := ...` creates a new `ok` variable. Be careful if you intend to use `ok` after the `if` block.
- **Checking `ok` on a closed channel but continuing to read**: After `ok=false`, the channel will continue to return `(zero, false)` forever — it will not block. You must break or return.

## Debugging walkthrough

Consider this code:

```go
package main

import "fmt"

func main() {
	var data interface{} = "hello world"
	length := data.(int) // panic: interface conversion: string is not int
	fmt.Println(length)
}
```

**Symptom**: Runtime panic.

**Root cause**: Blind type assertion without checking compatibility.

**Fix with comma-ok**:

```go
if length, ok := data.(int); ok {
    fmt.Println("Length:", length)
} else {
    fmt.Println("data is not an int")
}
```

**Another example**:

```go
func getCached(key string) string {
    cache := map[string]string{"a": "alpha"}
    return cache[key] // If key not found, returns ""
}
```

**Symptom**: Returns empty string for missing keys — indistinguishable from a cached empty string.

**Fix**:

```go
func getCached(key string) (string, bool) {
    cache := map[string]string{"a": "alpha"}
    val, ok := cache[key]
    return val, ok
}
```

## Production notes

- Always use the comma-ok form for type assertions unless you are certain the type matches and want a panic on failure (which is almost never — prefer explicit handling).
- The cost of the two-value form vs single-value form is negligible — a single boolean comparison. There is no performance reason to skip `ok`.
- In hot paths, you can structure your maps and type assertions so that `ok` is `true` in the common case, keeping the fast path linear.
- For channels, `for v := range ch` is cleaner than a manual comma-ok loop, but the comma-ok form is essential when you need to detect closure explicitly (e.g., in select statements).
- Some linters (like `errcheck`-style rules) can flag ignored `ok` values.

## Performance implications

| Operation | Single-value | Two-value | Overhead |
|---|---|---|---|
| Map lookup | `mapaccess1` | `mapaccess2` | ~1 extra bool compare |
| Type assertion | `assertI2I` (panics) | `assertI2I2` | Branch instead of panic path |
| Channel receive | `chanrecv1` | `chanrecv2` | ~1 extra bool compare |

The performance difference between single-value and two-value forms is **negligible** — a single boolean test-and-branch. Never sacrifice correctness for this micro-optimization.

## Practice task

Write a function `LookupSafe(data map[string]interface{}, key string, expectedType string) (interface{}, bool, error)` that:

1. Checks if `key` exists in `data`. If not, returns `(nil, false, nil)`.
2. Attempts a type assertion on the value to the type named by `expectedType` (which can be `"string"`, `"int"`, `"float64"`, `"bool"`).
3. If the assertion succeeds, returns `(value, true, nil)`.
4. If the assertion fails, returns `(nil, false, fmt.Errorf("expected %s, got %T", expectedType, val))`.
5. Write a `main()` that exercises all paths — key exists and type matches, key exists and type mismatches, key missing.

## Tests / verification

```bash
go run ./curriculum/modules/03-programming-fundamentals/lessons/18-comma-ok-idiom
go test ./curriculum/modules/03-programming-fundamentals/lessons/18-comma-ok-idiom
```

## Review questions

1. What are the three Go constructs that support the comma-ok form?
2. What happens if you do a single-value type assertion and the assertion fails?
3. Why does `for v := range ch` not need an explicit `ok` check?
4. How can you distinguish a map entry whose value is `0` from a missing key?
5. What does `val, ok := <-ch` return after a closed channel has been fully drained?

## NEXT UP

Strings (core-03-19): Explore Go's string type — immutable bytes under the hood, UTF-8 support, and the standard library for text manipulation.
