# Receiver sets

## Learning objective

Determine the method set of `T` and `*T`, explain why the value receiver set excludes pointer methods, and apply interface satisfaction rules correctly.

## Why this matters

Receiver sets are the bridge between methods and interfaces. When you assign a value of type `T` to an interface variable, Go checks the method set of `T` — not `*T`. Understanding which methods are available on `T` vs `*T` prevents confusing compile errors and clarifies why some types satisfy interfaces only through pointers.

## Mental model

Every type `T` has a method set — the set of methods callable on a value of that type. `*T` also has a method set, which is a superset of `T`'s. The rule is simple: **value methods are in both sets; pointer methods are only in `*T`'s set**. When you have a value, you cannot call a pointer method on it directly; when you have a pointer, you can call everything.

## Core idea

**Method set of `T`** (value type `T`):
- All methods declared with receiver `T`.

**Method set of `*T`** (pointer to `T`):
- All methods declared with receiver `*T`.
- All methods declared with receiver `T` (Go automatically dereferences).

**Interface satisfaction**: A type `T` satisfies an interface `I` if `T`'s method set contains all methods of `I`. If `I` requires a method with a pointer receiver, only `*T` satisfies `I` — not `T`.

```go
type Stringer interface {
    String() string
}

type Counter struct {
    value int
}

func (c Counter) String() string {  // value receiver
    return fmt.Sprintf("Count: %d", c.value)
}

// Counter satisfies Stringer (method is in value set)
var s Stringer = Counter{}  // OK

func (c *Counter) Increment() {  // pointer receiver
    c.value++
}

// Counter does NOT have Increment in its value receiver set
// var s2 interface{ Increment() } = Counter{}  // COMPILE ERROR
var s2 interface{ Increment() } = &Counter{}  // OK
```

## Under the hood

The compiler builds method sets during type-checking. For each named type `T`, it scans all methods in the package with receiver `T` or `*T`. The method set of `T` includes only `T`-receiver methods. The method set of `*T` includes both. This is a compile-time check — no runtime dispatch overhead for method set resolution.

## How Go uses it

- **Interface satisfaction checking**: `io.Reader` requires `Read(p []byte) (n int, err error)`. If declared with a pointer receiver, only `*os.File` satisfies it.
- **JSON marshaling**: `json.Marshal` checks for `json.Marshaler` on the value passed. If `MarshalJSON` uses a pointer receiver, you must pass a pointer.
- **Error wrapping**: `Unwrap() error` is typically on `*fmt.wrapError`.
- **Sorting**: `sort.Interface` methods (`Len`, `Less`, `Swap`) are pointer-receiver methods on `*sort.IntSlice`.

## Go example

```go
package main

import "fmt"

type Greeter interface {
	Greet() string
}

type Person struct {
	Name string
}

func (p Person) Greet() string { // value receiver
	return "Hello, I'm " + p.Name
}

type Robot struct {
	Model string
}

func (r *Robot) Greet() string { // pointer receiver
	return "Beep boop, model " + r.Model
}

func main() {
	var g Greeter

	g = Person{Name: "Alice"} // OK: value satisfies Greeter
	fmt.Println(g.Greet())

	g = &Person{Name: "Bob"} // OK: *T also satisfies
	fmt.Println(g.Greet())

	// g = Robot{Model: "R2"} // COMPILE ERROR: Robot does not satisfy Greeter
	g = &Robot{Model: "R2"} // OK: *Robot satisfies Greeter
	fmt.Println(g.Greet())
}
```

## Step-by-step execution

For `g = &Robot{Model: "R2"}`:

1. Compiler checks if `*Robot`'s method set contains `Greet() string`.
2. `*Robot`'s method set includes all pointer-receiver methods of `Robot`, including `(*Robot).Greet`.
3. Assignment is valid.
4. At runtime, `g.Greet()` calls `(*Robot).Greet` with receiver `&Robot{Model: "R2"}`.
5. Inside `Greet`, receiver is dereferenced and `Model` is read.

For `g = Person{Name: "Alice"}`:

1. Compiler checks `Person`'s method set for `Greet() string`.
2. `Person`'s method set includes `Person.Greet` (value receiver).
3. Assignment is valid.

## Common mistakes

- Mistake: Trying to assign a value to an interface when the interface method uses a pointer receiver.
  - Fix: Assign `&T{}` instead of `T{}`.

- Mistake: Assuming a type satisfies an interface because the method set of `*T` has all methods.
  - Fix: Check the method set of `T` specifically. Only `*T` includes pointer-receiver methods.

- Mistake: Declaring an interface with methods that all use value receivers, then switching one to a pointer receiver — this breaks value assignments.
  - Fix: Be deliberate about receiver choice on interface methods.

- Mistake: Forgetting that a nil pointer of type `*T` can satisfy an interface, leading to non-nil interface values holding nil pointers.
  - Fix: Check the concrete value inside the interface, not just whether the interface is nil.

## Debugging walkthrough

This code does not compile:

```go
type ByteCounter int

func (c *ByteCounter) Write(p []byte) (int, error) {
    *c += ByteCounter(len(p))
    return len(p), nil
}

func main() {
    var w io.Writer
    w = ByteCounter(0) // COMPILE ERROR
}
```

**Symptom**: `ByteCounter does not satisfy io.Writer (Write method has pointer receiver)`.

**Root cause**: `Write` has a pointer receiver. `ByteCounter`'s method set does not include `Write`. Only `*ByteCounter` satisfies `io.Writer`.

**Fix**: Use `w = new(ByteCounter)` or `w = &ByteCounter{}`.

## Production notes

- **Interface hygiene**: If a type is meant to satisfy an interface as a value, ensure all required methods use value receivers.
- **Consistency**: If one method needs a pointer receiver (e.g., for mutation), all methods on that type typically use pointer receivers. This avoids confusion about interface satisfaction.
- **Documentation**: Clearly document whether a type satisfies an interface as `T` or `*T`.
- **Code generation**: Tools like `impl` can verify interface satisfaction at compile time.

## Performance implications

- **Method set size** has no runtime cost. All checks are compile-time.
- **Indirect calls**: Calling a method through an interface is slower than a direct call (one extra dereference). But this is the cost of polymorphism, not of receiver sets.
- **Heap allocation**: Assigning a value to an interface variable wraps it in an `interface{}` value, which may cause heap allocation if the value does not fit in two machine words.

## Practice task

Define an interface `Speaker` with method `Speak() string`. Create two types: `Dog` (value receiver `Speak`) and `Cat` (pointer receiver `Speak`). In `main()`, show which assignments work: `var s Speaker = Dog{}`, `var s Speaker = &Dog{}`, `var s Speaker = Cat{}`, `var s Speaker = &Cat{}`. Print results. Verify with tests.

## Tests / verification

```bash
go run ./curriculum/modules/05-types-interfaces-packages-modules/lessons/05-receiver-sets
go test ./curriculum/modules/05-types-interfaces-packages-modules/lessons/05-receiver-sets
```

## Review questions

1. What methods are in the method set of `T`?
2. What methods are in the method set of `*T`?
3. If an interface requires a method with a pointer receiver, can you assign a value of type `T` to that interface variable?
4. Does Go automatically generate a value-receiver method for every pointer-receiver method?
5. Why does assigning a value to an interface require computing the method set at compile time?

## NEXT UP

Composition — Go's answer to inheritance, using struct embedding to reuse fields and methods.
