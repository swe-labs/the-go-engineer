# Interfaces

## Learning objective

Define and satisfy Go interfaces implicitly, understand interface values as type+value pairs, distinguish nil interfaces from nil concrete values, and use the empty interface.

## Why this matters

Interfaces are Go's primary tool for abstraction and polymorphism. Unlike languages with explicit `implements` declarations, Go uses **structural typing** — a type satisfies an interface automatically if it has the required methods. This makes interfaces flexible, testable, and decoupled. Mastering interfaces is the gateway to writing generic, composable Go code.

## Mental model

An interface is a contract: "if you have these methods, you are my type." Go checks the contract at compile time, but the binding happens when you assign a value to an interface variable. Think of an interface variable as a wrapper that stores two things: the concrete type and its value. When you call a method on the interface, Go unwraps the concrete value and calls its method.

## Core idea

**Interface type definition**:

```go
type Stringer interface {
    String() string
}
```

**Implicit satisfaction (duck typing)**: Any type that has a `String() string` method automatically satisfies `Stringer`. No `implements` keyword needed.

**Interface values (type + value pair)**: At runtime, an interface value is represented as a two-word structure: a pointer to the type information (the **concrete type**) and a pointer to the data (the **concrete value**). This is sometimes called the "iface" structure.

```go
var s Stringer
s = MyType{value: 42}
// s = (type: MyType, value: &{42})
```

**Nil interface**: An interface variable that has not been assigned anything has both type and value nil. Calling a method on it panics.

**Non-nil interface holding nil pointer**: If you assign a typed nil pointer (e.g., `var p *MyType; s = p`), the interface is non-nil (it has a type) even though the value inside is nil. This is a common gotcha.

**Empty interface**: `interface{}` (or `any` in Go 1.18+) has no methods. Every type satisfies it. It is Go's universal type, used for values of unknown type.

## Under the hood

The runtime representation of an interface is:

```go
type iface struct {
    tab  *itab   // type info + method table
    data unsafe.Pointer // pointer to concrete value
}
```

The `itab` (interface table) is a cache generated at runtime, mapping the concrete type to the interface's method set. When you assign a value to an interface, Go checks at compile time whether the type satisfies the interface, then at runtime the `itab` is created (once) and cached.

## How Go uses it

- **io.Reader** and **io.Writer** are the most ubiquitous interfaces in the standard library.
- **error** is a built-in interface with one method: `Error() string`.
- **fmt.Stringer** provides `String()` for custom string formatting.
- **sort.Interface** requires `Len`, `Less`, `Swap`.
- **http.Handler** requires `ServeHTTP(ResponseWriter, *Request)`.
- **json.Marshaler** requires `MarshalJSON() ([]byte, error)`.

## Go example

```go
package main

import "fmt"

type Stringer interface {
	String() string
}

type Book struct {
	Title  string
	Author string
}

func (b Book) String() string {
	return b.Title + " by " + b.Author
}

func printAny(v interface{}) {
	fmt.Printf("value=%v type=%T\n", v, v)
}

func main() {
	var s Stringer
	fmt.Printf("nil? %v, type=%T\n", s == nil, s)

	s = Book{Title: "1984", Author: "Orwell"}
	fmt.Println(s.String())
	fmt.Printf("type=%T\n", s)

	var p *Book
	s = p                // assigning typed nil
	fmt.Printf("nil? %v, type=%T\n", s == nil, s) // false!
	fmt.Printf("value nil? %v\n", s == nil)       // true

	var empty interface{}
	empty = 42
	printAny(empty)
	empty = "hello"
	printAny(empty)
	empty = Book{Title: "Brave New World", Author: "Huxley"}
	printAny(empty)
}
```

## Step-by-step execution

For `var s Stringer; s = Book{Title: "1984", Author: "Orwell"}`:

1. Compile-time: check that `Book` has `String() string` method — yes, value receiver.
2. Runtime: allocate interface value `s`: `{tab: &itab{Book, Stringer}, data: &Book{...}}`.
3. The concrete `Book` value is heap-allocated (or stack-allocated if the compiler proves it doesn't escape).
4. `s.String()` dereferences the `itab` to find the `String` method, then calls it with the `data` pointer as receiver.

For `var p *Book; s = p; s == nil`:

1. `p` is uninitialized `*Book`, value is `nil`.
2. `s = p` stores `{tab: &itab{*Book, Stringer}, data: nil}`.
3. `s == nil` compares the interface to `nil`. `s.tab != nil`, so `s != nil` — even though `data` is nil.

## Common mistakes

- Mistake: Checking `if err == nil` after calling a function that returns a typed nil pointer.
  - Fix: Always return `nil` explicitly (nil interface, not typed nil) or use a helper: `if err != nil { return err }` in the caller.

- Mistake: Defining an interface with too many methods (interface bloat).
  - Fix: Prefer small interfaces (1-2 methods). Compose larger interfaces from small ones.

- Mistake: Using `interface{}` when a concrete type would work.
  - Fix: Only use `interface{}` when you genuinely don't know or don't care about the type.

- Mistake: Expecting type assertions to work on nil interfaces.
  - Fix: Check `s != nil` before type-asserting.

## Debugging walkthrough

This code has a subtle nil bug:

```go
type MyError struct {
    Msg string
}

func (e *MyError) Error() string {
    if e == nil {
        return ""
    }
    return e.Msg
}

func getError() *MyError {
    return nil
}

func main() {
    var err error = getError() // *MyError is assigned to error
    if err != nil {            // TRUE — interface is non-nil!
        fmt.Println("unexpected error:", err)
    }
}
```

**Symptom**: "unexpected error: " is printed even though `getError` returns nil.

**Root cause**: `getError()` returns a nil `*MyError`, which is a typed nil. Assigning it to `error` creates a non-nil interface (type = `*MyError`, value = nil).

**Fix**: Change `getError` to return `error` directly: `func getError() error { return nil }`.

## Production notes

- **Accept interfaces, return structs** is a common Go idiom. Inputs should be abstract; outputs should be concrete.
- **Interface pollution** — defining interfaces before they are needed — is a common mistake. Define interfaces where they are used, not where they are implemented.
- **`any`** is the Go 1.18+ alias for `interface{}`. Prefer `any` in new code.
- **Nil checks** on interfaces are tricky. Use `reflect.ValueOf(err).IsNil()` with care, or avoid typed nils entirely.

## Performance implications

- **Interface calls** are indirect through the `itab`. This is ~1-2 ns slower than a direct call on modern hardware.
- **Heap allocation**: Assigning a non-pointer value to an interface typically causes a heap allocation (escape to heap). `*T` assignments may also escape if `T` is small.
- **Method table lookup** is cached. The first call to an interface method on a given (type, interface) pair builds the `itab`; subsequent calls reuse it.

## Practice task

Define an interface `Notifier` with method `Notify(message string) error`. Create two implementations: `EmailNotifier` (with field `Address string`) and `SMSNotifier` (with field `Phone string`). Write a function `SendAlerts(notifiers []Notifier, message string)` that calls `Notify` on each. In `main()`, create a slice of both types and send alerts. Verify with tests.

## Tests / verification

```bash
go run ./curriculum/modules/05-types-interfaces-packages-modules/lessons/08-interfaces
go test ./curriculum/modules/05-types-interfaces-packages-modules/lessons/08-interfaces
```

## Review questions

1. How does Go determine whether a type satisfies an interface?
2. What is the runtime representation of an interface value?
3. Why can an interface variable be non-nil even when it holds a nil pointer?
4. What types satisfy `interface{}` / `any`?
5. When should you define an interface in your package?

## NEXT UP

Small interfaces — the power of 1-2 method interfaces and the interface segregation principle.
