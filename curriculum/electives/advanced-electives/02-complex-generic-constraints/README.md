# Complex generic constraints

## When you need this

New Go developers use any or interface{} as a constraint because they do not understand Go 1.18+ constraint features. They lose compile-time type safety and reintroduce runtime type assertions. Complex generic constraints solve this by allowing type unions (|), type approximations (~T), and constraint embedding to precisely define which types are acceptable while catching violations at compile time.

## When you do not need this

When simpler concrete types or standard core library primitives suffice without adding abstraction overhead.

## Prerequisites

- elective-01

## Mental model

A constraint is a compile-time contract that defines which types a type parameter accepts. Simple constraints (any, comparable) accept broad categories. Complex constraints narrow the set using type unions (A | B | C), type approximations (~T accepts T and any type with underlying type T), and method requirements. The compiler enforces the constraint at every call site — if the type argument does not satisfy the constraint, the code does not compile.

## How it works

Go's constraint interface can include type elements: type unions (int | string) and type approximations (~int). These are not methods — they are compile-time restrictions on the set of acceptable type arguments. The compiler checks constraint satisfaction during type checking, before monomorphization. A type argument satisfies a constraint if: it implements all methods in the constraint AND it is a member of the constraint's type set (for type elements). Type approximation (~T) means the type argument's underlying type must be T, not the type itself.

## Go example

```go
package main

import "fmt"

func main() {
    fmt.Println("Elective: Complex generic constraints")
}
```

## Tradeoffs

- Complex constraints do not add runtime overhead — monomorphization still produces concrete code
- Compile-time type checking is more expensive with complex constraints due to constraint satisfaction verification
- Type approximation (~T) constraints do not affect runtime behavior — they only widen the set of acceptable types at compile time
- Deep constraint hierarchies are resolved at compile time and have no runtime cost

## Common mistakes

- Using any as a constraint when a more specific constraint is possible — loses compile-time type safety and allows invalid type arguments.
- Writing constraints that are too restrictive — preventing valid callers and reducing reuse, requiring duplicated functions for different types.
- Confusing interface{} (any) with a constraint interface — any accepts every type, while a constraint restricts which types are valid.
- Forgetting that ~int matches int AND any type with underlying type int (e.g., type MyInt int) — this is often unexpected.
- Using type unions with types that do not share common operations — the generic function body can only use operations common to all union members.

## Practice task

The learner must write a constraint with a type union and type approximation, write a generic function that uses the constraint, demonstrate which types satisfy and do not satisfy the constraint, and explain why constraint satisfaction is checked at compile time rather than runtime.

## Production notes

Complex constraints are used in Go libraries for numeric types (constraints.Integer, constraints.Float), ordered types (constraints.Ordered), and serialization (json.Marshaler + json.Unmarshaler). The Go standard library's sort.Slice is not generic because Go 1.18 generics do not support method-level type parameters — instead, custom sort functions use constraints.Ordered for sort.SortFunc in golang.org/x/exp/slices.

## Related core lessons

- Module 16.01 — Generics
- Module 16.15 — Demo preparation
