# Methods

## Learning objective

Write methods on any named Go type using a receiver parameter; distinguish methods from functions; and use method expressions to convert methods to function values.

## Why this matters

Methods attach behavior to data. In Go, you can define methods on any named type — not just structs. This lets you implement interfaces, encapsulate logic with your types, and write code that reads naturally (`user.Name()`, `user.Save()`). Understanding methods is the prerequisite for interfaces, the heart of Go's type system.

## Mental model

A method is a function with a **receiver** — an extra parameter that appears before the function name. When you call `obj.Method()`, Go passes `obj` as the receiver argument. The receiver binds the function to the type, just like `self` in Python or `this` in JavaScript, but explicit and visible in the signature.

## Core idea

Method declaration syntax:

```go
func (receiver Type) MethodName(params) results {
    // body
}
```

The receiver appears as a parameter list with one element before the function name. By convention the receiver name is one or two letters (the first letter of the type).

**Value receiver** — the method receives a copy of the value:

```go
type Circle struct {
    Radius float64
}

func (c Circle) Area() float64 {
    return 3.14159 * c.Radius * c.Radius
}
```

Call: `c := Circle{Radius: 5}; area := c.Area()`.

**Method vs function**: The same logic can be written as `func Area(c Circle) float64`. The method form enables interface satisfaction and chaining; the function form is more explicit and can be used as a value without binding.

**Method on any named type**: Methods are not limited to structs. You can define methods on any type declared with `type`:

```go
type Celsius float64
func (c Celsius) F() float64 { return float64(c)*1.8 + 32 }
```

**Method expression**: You can reference a method as a function value by specifying the type: `Circle.Area` produces `func(Circle) float64`. Call it with the receiver as the first argument: `Circle.Area(c)`.

## Under the hood

The compiler desugars methods into regular functions with the receiver as the first parameter. `func (c Circle) Area() float64` becomes `func Circle_Area(c Circle) float64` internally. Method calls are direct function calls — no dynamic dispatch unless the method is on an interface. Value receiver methods receive a copy of the struct in the caller's stack frame.

## How Go uses it

- **Interface satisfaction**: A type satisfies an interface by implementing its methods.
- **Error types**: `Error() string` method makes any type an `error`.
- **String formatting**: `String() string` method makes any type a `Stringer`.
- **Getters**: `func (p *Person) Name() string { return p.name }`
- **Setters**: `func (p *Person) SetName(n string) { p.name = n }`

## Go example

```go
package main

import "fmt"

type Celsius float64

func (c Celsius) F() float64 {
	return float64(c)*1.8 + 32
}

type Rectangle struct {
	Width, Height float64
}

func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

func (r Rectangle) Perimeter() float64 {
	return 2 * (r.Width + r.Height)
}

func main() {
	c := Celsius(100.0)
	fmt.Printf("%.1f°C = %.1f°F\n", c, c.F())

	rect := Rectangle{Width: 3, Height: 4}
	fmt.Printf("area: %.1f, perimeter: %.1f\n", rect.Area(), rect.Perimeter())

	methodExpr := Rectangle.Area
	fmt.Printf("via method expr: %.1f\n", methodExpr(rect))
}
```

## Step-by-step execution

For `rect.Area()` with `rect = Rectangle{Width: 3, Height: 4}`:

1. Go evaluates the receiver `rect`: a `Rectangle` value `{3, 4}`.
2. The compiler resolves `Area` on type `Rectangle`.
3. The receiver value is copied onto the call stack (new stack frame for `Area`).
4. Inside `Area`, `r.Width` reads the copy's `Width` field: `3`.
5. `r.Height` reads the copy's `Height` field: `4`.
6. Returns `3 * 4 = 12`.
7. Caller receives `12`.

For `methodExpr(rect)` where `methodExpr := Rectangle.Area`:

1. `Rectangle.Area` is a method expression, typed as `func(Rectangle) float64`.
2. Calling `methodExpr(rect)` passes `rect` as the first (receiver) argument.
3. Same as above from step 2.

## Common mistakes

- Mistake: Defining a method with the same name on the same type in the same package.
  - Fix: Go does not allow overloading. Use unique names.

- Mistake: Expecting a method to modify its receiver when using a value receiver.
  - Fix: Use a pointer receiver (`*T`) to mutate the original.

- Mistake: Confusing method expression `T.Method` with method value `x.Method`.
  - `T.Method` is `func(T)`: the receiver is explicit. `x.Method` is `func()`: the receiver is already bound to `x`.

- Mistake: Defining methods on built-in types directly (e.g., `func (x int) Double()`).
  - Fix: You can only define methods on named types in the same package. Create a `type MyInt int` alias.

## Debugging walkthrough

This code does not compile:

```go
package main

type MyInt int

func (x int) Double() int {
    return x * 2
}

func main() {
    var x int = 5
    println(x.Double())
}
```

**Symptom**: `cannot define new methods on non-local type int`.

**Root cause**: Methods can only be defined on types declared in the same package. `int` is a built-in type, not defined in this package.

**Fix**: Define the method on the named type `MyInt`:

```go
type MyInt int

func (x MyInt) Double() MyInt {
    return x * 2
}
```

Or call the method on `MyInt` instead of `int`.

## Production notes

- **Receiver naming**: Use short names (one or two letters) consistently. `c` for `Celsius`, `r` for `Rectangle`.
- **Getters don't use `Get` prefix** in Go: `user.Name()` not `user.GetName()`.
- **Setters** may use `Set` prefix: `user.SetName("Alice")`.
- **Method on `error`**: The `Error()` method is the only requirement for the `error` interface.
- **Methods on basic types**: Creating named types from primitives and adding methods is common for domain modeling.

## Performance implications

- **Value receiver copy cost**: For small types (ints, bools, small structs), copying is negligible. For large structs, prefer pointer receivers.
- **Inlining**: The compiler can inline small value-receiver methods, eliminating call overhead.
- **Escape analysis**: Value receivers keep values on the stack when possible; pointer receivers can cause heap allocation if the pointer escapes.

## Practice task

Define a named type `Miles float64`. Add a method `ToKilometers()` that returns `Miles * 1.60934`. Also define a struct `Trip` with fields `Distance Miles` and `DurationMinutes float64`. Add a method `Speed()` that returns miles per hour. In `main()`, create a `Trip` value and print its distance in kilometers and its speed. Verify with the test.

## Tests / verification

```bash
go run ./curriculum/modules/05-types-interfaces-packages-modules/lessons/02-methods
go test ./curriculum/modules/05-types-interfaces-packages-modules/lessons/02-methods
```

## Review questions

1. What is the syntax difference between a method and a function in Go?
2. Can you define a method on `type MySlice []string`? Why or why not?
3. What is a method expression, and how do you call it?
4. What does `T.Method` represent, and how does it differ from `x.Method`?
5. If you call a value-receiver method on a pointer, does Go dereference it automatically?

## NEXT UP

Pointer receivers — mutating the receiver and understanding when to use `*T`.
