# Pointer receivers

## Learning objective

Write methods with pointer receivers to mutate the receiver, handle nil receivers safely, and decide when a pointer receiver is the right choice.

## Why this matters

Value receiver methods operate on a copy — they cannot modify the original. Pointer receivers let you change the receiver's state, avoid copying large structs, and handle nil receiver values gracefully. Nearly every Go codebase uses pointer receivers for methods that mutate state, making this essential for writing idiomatic Go.

## Mental model

A pointer receiver is a method whose receiver type is `*T` — a pointer to your named type. When you call `p.Method()`, Go passes a pointer to `p`. Inside the method, you read and write through that pointer, directly affecting the original value. If the receiver is `nil`, you can still call the method, but dereferencing the pointer inside will panic.

## Core idea

Syntax:

```go
type Counter struct {
    value int
}

func (c *Counter) Increment() {
    c.value++
}
```

**Mutating receiver**: Inside `Increment`, `c.value++` modifies the original `Counter` because `c` is a pointer.

**When to use pointer receiver**:

1. Method must mutate the receiver.
2. Receiver is a large struct (copy is expensive).
3. Receiver is a slice or map (they are reference types, but consistency may demand a pointer).
4. All methods on a type should use pointer receivers for consistency.

**Nil receiver handling**: A pointer method can be called on a nil pointer. The method can check for nil and handle it gracefully:

```go
func (c *Counter) Value() int {
    if c == nil {
        return 0
    }
    return c.value
}
```

**Pointer receiver with value types**: When you have a value and call a pointer-receiver method, Go automatically takes the address if the value is addressable:

```go
c := Counter{}
c.Increment() // Go does (&c).Increment()
```

## Under the hood

A pointer receiver method is compiled to a function whose first parameter is a pointer: `func Counter_Increment(c *Counter)`. The caller passes the address of the receiver. Inside the method, field access like `c.value` is compiled to an indirect load through the pointer. Mutation writes to the memory at that address, visible to the caller after the method returns.

## How Go uses it

- **Setters**: `func (p *Person) SetName(n string)`.
- **Accumulators**: counters, builders, buffers.
- **Large struct processing**: `func (db *Database) Query(...)`.
- **Nil-safe accessors**: methods on types that may be nil (linked list nodes, tree nodes).
- **Unmarshal methods**: `func (p *Point) UnmarshalJSON(data []byte) error`.

## Go example

```go
package main

import "fmt"

type Counter struct {
	value int
}

func (c *Counter) Increment() {
	c.value++
}

func (c *Counter) Add(n int) {
	c.value += n
}

func (c *Counter) Value() int {
	if c == nil {
		return 0
	}
	return c.value
}

func (c *Counter) Reset() {
	c.value = 0
}

func main() {
	c := Counter{}
	c.Increment()
	c.Increment()
	c.Add(3)
	fmt.Println("value:", c.Value()) // 5

	c.Reset()
	fmt.Println("after reset:", c.Value()) // 0

	var nilCounter *Counter
	fmt.Println("nil receiver:", nilCounter.Value()) // 0 (safe)
}
```

## Step-by-step execution

For `c.Increment()` where `c := Counter{value: 0}`:

1. `c` is a `Counter` value on the stack with `value = 0`.
2. Go sees method `Increment` has receiver `*Counter`.
3. Go takes the address of `c`: `&c`.
4. The address is pushed as the receiver argument.
5. Inside `Increment`, `c.value` dereferences the pointer and reads `0`.
6. `c.value++` writes `1` back through the pointer.
7. Method returns. `c.value` on the caller's stack is now `1`.

## Common mistakes

- Mistake: Forgetting to dereference a nil pointer receiver before accessing fields.
  - Fix: Always check `if c == nil { return }` at the start if the receiver can be nil.

- Mistake: Using a pointer receiver when the method does not mutate and the type is tiny.
  - Fix: Prefer value receivers for small, immutable types.

- Mistake: Mixing value and pointer receivers inconsistently.
  - Fix: Pick one convention per type. If any method needs a pointer receiver, make them all pointer receivers.

- Mistake: Calling a pointer-receiver method on a non-addressable value (e.g., map value, function return).
  - Fix: Store the value in a variable first, then call the method.

## Debugging walkthrough

This code fails:

```go
type Point struct {
    X, Y int
}

func (p *Point) Move(dx, dy int) {
    p.X += dx
    p.Y += dy
}

func main() {
    getPoint := func() Point { return Point{1,2} }
    getPoint().Move(1, 1) // COMPILE ERROR
}
```

**Symptom**: `cannot call pointer method on getPoint()` or `cannot take address of getPoint()`.

**Root cause**: `getPoint()` returns a value, not an addressable variable. Go cannot take the address of a function call result.

**Fix**: Assign to a variable first:

```go
p := getPoint()
p.Move(1, 1)
```

## Production notes

- **Consistency rule**: If any method on a type needs a pointer receiver, all methods on that type usually use pointer receivers.
- **Nil receivers** are idiomatic in Go. Use them to implement nil-safe types (e.g., a nil `*bytes.Buffer` behaves as an empty buffer).
- **Never store a pointer to a zero-value struct** unless you intend to share it. `var x T; f(&x)` is fine; `f(&T{})` creates a new allocation.
- **Interface values** containing a pointer are non-nil even if the pointer is nil. This is a common gotcha.

## Performance implications

- **Indirection cost**: A pointer receiver adds one level of indirection for each field access. For hot-path code, value receivers on small structs can be faster by avoiding the pointer load.
- **Escape analysis**: A pointer receiver may cause the receiver to escape to the heap if the pointer is stored or returned. This adds allocation overhead.
- **Large structs**: Pointer receivers avoid copying the entire struct, which can be a significant win for structs > 100 bytes or those with many fields.

## Practice task

Define a `BankAccount` struct with a private `balance float64` field. Add pointer receiver methods `Deposit(amount float64)`, `Withdraw(amount float64) error` (return error if insufficient funds), and `Balance() float64`. In `main()`, create an account, deposit 100, withdraw 30, print the balance, then withdraw 100 (should fail). Print the error.

## Tests / verification

```bash
go run ./curriculum/modules/05-types-interfaces-packages-modules/lessons/03-pointer-receivers
go test ./curriculum/modules/05-types-interfaces-packages-modules/lessons/03-pointer-receivers
```

## Review questions

1. When would you choose a pointer receiver over a value receiver?
2. Can you call a pointer-receiver method on a value that is not addressable? Give an example.
3. What happens if you call a method with a pointer receiver on a nil pointer?
4. Why do most Go style guides recommend consistency — either all value or all pointer receivers on a type?
5. How does the compiler handle `c.Increment()` when `c` is a value and `Increment` has a pointer receiver?

## NEXT UP

Value receivers — the other side of the receiver coin, with copy semantics and immutability guarantees.
