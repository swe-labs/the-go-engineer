# Composition

## Learning objective

Compose types using struct embedding, access promoted fields and methods, and explain why Go prefers composition over inheritance.

## Why this matters

Go rejects classical inheritance. Instead of subclassing, Go uses **composition** — embedding one struct inside another. This produces flatter, more flexible designs that are easier to refactor and test. Every Go engineer must understand composition to model relationships between types without fighting the language.

## Mental model

Think of embedding as copy-paste of fields and methods. When you embed `type A struct { X int }` inside `B`, it is as if `B` got an `X` field of its own, plus any methods `A` has. But the "copy" is logical — at runtime, the embedded field is a named sub-value, accessed through a synthetic field named after the embedded type.

## Core idea

**Struct embedding** uses the type name without a field name:

```go
type Address struct {
    City, State string
}

type Person struct {
    Name    string
    Address // embedded
}
```

**Promoted fields**: Fields of `Address` are accessible directly on `Person`:

```go
p := Person{Name: "Alice", Address: Address{City: "NYC", State: "NY"}}
fmt.Println(p.City)     // promoted: p.City == p.Address.City
fmt.Println(p.Address.City) // also works via full path
```

**Promoted methods**: Methods of `Address` are also promoted to `Person`.

**Composition over inheritance**: Instead of "a `Person` IS-A `Address`", Go says "a `Person` HAS-A `Address`" — but provides syntactic sugar for delegation. This avoids the fragile base class problem.

**Flattening promoted names**: If both `Person` and `Address` have a field `ID`, the outer field shadows the inner one. Access `p.Address.ID` to disambiguate.

## Under the hood

Embedding is purely a compile-time mechanism. The embedded field is stored as a regular field named after the type. Promotion creates implicit forwarding methods/field accessors in the compiler's symbol table. At runtime, there is no difference between a promoted field access and a direct field access — the offset is computed at compile time.

## How Go uses it

- **HTTP handlers**: `type MyServer struct { *http.Server }` — embed to add methods.
- **Logging**: `type Logger struct { *log.Logger }` — add level filtering.
- **Database models**: embed a `BaseModel` with `ID`, `CreatedAt`, `UpdatedAt`.
- **Testing**: embed `testing.T` in a custom test helper to forward `Error`, `Fatal`, etc.
- **Error wrapping**: `fmt.Errorf("... %w", err)` creates a wrapped error via embedding.

## Go example

```go
package main

import "fmt"

type Address struct {
	City, State string
}

func (a Address) Full() string {
	return a.City + ", " + a.State
}

type Person struct {
	Name    string
	Address // embedded
}

type Employee struct {
	Person   // embedded
	Position string
}

func main() {
	p := Person{
		Name:    "Alice",
		Address: Address{City: "Portland", State: "OR"},
	}
	fmt.Println(p.Name)
	fmt.Println(p.City)         // promoted
	fmt.Println(p.Address.City) // explicit
	fmt.Println(p.Full())       // promoted method

	e := Employee{
		Person:   Person{Name: "Bob", Address: Address{City: "Seattle", State: "WA"}},
		Position: "Engineer",
	}
	fmt.Println(e.Name, e.Position)
	fmt.Println(e.Full()) // double promotion: Employee → Person → Address
}
```

## Step-by-step execution

For `p.City` where `p = Person{Name: "Alice", Address: Address{City: "Portland"}}`:

1. Compiler sees `p.City` on type `Person`.
2. `City` is not a direct field of `Person`, so compiler checks promoted fields.
3. `Person` embeds `Address`. `Address` has field `City`.
4. Compiler rewrites `p.City` to `p.Address.City`.
5. At runtime, the offset of `Address` within `Person` is computed, then offset of `City` within `Address`.
6. Value `"Portland"` is loaded.

For `e.Full()` where `e` is `Employee`:

1. Compiler looks for `Full` method on `Employee` — not found.
2. Checks promoted methods from `Employee`'s embedded fields (`Person`).
3. `Person` does not have `Full`, but embeds `Address`.
4. `Address` has `Full` method.
5. Compiler generates an implicit forwarding: `Employee.Full()` → `e.Person.Full()` → `e.Person.Address.Full()`.

## Common mistakes

- Mistake: Treating embedding as inheritance and expecting polymorphism through embedding.
  - Fix: Embedding promotes fields/methods but does not create an IS-A relationship. Use interfaces for polymorphism.

- Mistake: Name collision between a promoted field and the outer type's field.
  - Fix: The outer field always shadows the inner one. Use the full path to access the inner field.

- Mistake: Embedding a type multiple times in the same struct (e.g., embedding `Address` twice).
  - Fix: A type cannot be embedded more than once in the same struct. Use named fields instead.

- Mistake: Assuming embedding zero-values a pointer field.
  - Fix: `type T struct { *log.Logger }` — the embedded pointer is nil. Initialize it explicitly.

## Debugging walkthrough

This code fails to compile:

```go
type A struct {
    Value int
}
type B struct {
    A
    A // COMPILE ERROR: duplicate field A
}
```

**Symptom**: `duplicate field A`.

**Root cause**: You cannot embed the same type twice in one struct. The promoted field name would clash.

**Fix**: Use a named field: `type B struct { First A; Second A }`.

## Production notes

- **Embed for delegation, not for inheritance**. Embedding a type says "my type has a `Logger`", not "my type IS-A `Logger`".
- **Pointer embedding** (`*T`) is common for optional components: `type Server struct { *log.Logger }`. Nil-safe methods on `*log.Logger` make this work even when nil.
- **Test helpers**: Embed `*testing.T` in a custom struct to add assertion helpers: `type Asserter struct { *testing.T }`.
- **Avoid deep embedding chains**. Two levels is usually enough. Deep promotion chains make debugging harder.

## Performance implications

- **Zero overhead**: Promoted field access compiles to the same machine code as direct field access. No indirection.
- **Method promotion** creates an implicit wrapper function. The compiler usually inlines it, making the cost zero.
- **Struct size**: Embedding adds the size of the embedded type as a field. Padding rules apply as usual.

## Practice task

Define a `Contact` struct with embedded `Address` and an `Email string` field. Add a method `Info() string` to `Contact` that returns the name, email, and full address (using the promoted `Full()` method). In `main()`, create a `Contact` value and print `Info()`.

## Tests / verification

```bash
go run ./curriculum/modules/05-types-interfaces-packages-modules/lessons/06-composition
go test ./curriculum/modules/05-types-interfaces-packages-modules/lessons/06-composition
```

## Review questions

1. How does struct embedding differ from having a named field of the embedded type?
2. What happens if the outer type and the embedded type both have a field named `ID`?
3. Does embedding create an IS-A relationship? Why or why not?
4. How does Go resolve a promoted method call on a multi-level embedding chain?
5. Can you embed a pointer to a type? What is the zero value of such an embedded field?

## NEXT UP

Embedding — deeper dive into promoted methods, interface embedding, and conflict resolution.
