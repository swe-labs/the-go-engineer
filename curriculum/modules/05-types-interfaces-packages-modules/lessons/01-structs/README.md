# Structs

## Learning objective

Define and use Go struct types with fields, tags, and literals; distinguish named, anonymous, and zero-value structs; and compare struct values by value.

## Why this matters

A `struct` is Go's primary mechanism for grouping related data into a single composite type. Nearly every nontrivial Go program defines at least one struct. Structs form the data layer of APIs, database models, configuration, JSON payloads, and protocol messages. Mastering structs means you can model any domain accurately and safely.

## Mental model

Think of a struct as a blueprint for a record with named slots. Each slot holds a value of a specific type. When you create a struct value, you fill those slots. Two struct values are equal if every corresponding slot holds an equal value. The struct itself is a value, not a reference — assigning one struct to another copies all slots.

## Core idea

A struct is a composite type defined with the `type` keyword followed by the struct literal `struct { ... }` containing field declarations:

```go
type Person struct {
    Name string
    Age  int
}
```

Fields are accessed with dot notation: `p.Name`. Struct fields can have **tags** — string literals attached to field declarations that carry metadata (used by JSON encoders, ORMs, validators). Tags are invisible to normal code but accessible via `reflect`.

A struct literal creates a value:

```go
p := Person{Name: "Alice", Age: 30}   // named fields
p2 := Person{"Bob", 25}               // positional (order-dependent, rare)
```

**Zero-value struct**: When you declare a struct variable without initialization, every field gets its zero value (`""`, `0`, `nil`, etc.):

```go
var p Person   // p.Name == "", p.Age == 0
```

**Anonymous structs**: Declared inline without a named type:

```go
point := struct{X, Y int}{X: 1, Y: 2}
```

**Struct comparison**: Structs are comparable if all their fields are comparable. Two structs are equal if all corresponding fields are equal.

## Under the hood

A struct value is a contiguous block of memory containing each field's bytes in declaration order. The compiler inserts padding between fields to satisfy alignment requirements (a `uint64` must be 8-byte aligned, etc.). This means `sizeof` a struct may be larger than the sum of its fields due to padding.

Field tags are stored as a `reflect.StructTag` string in the compiled binary. At runtime, `reflect` can parse and read tags — no runtime overhead during normal field access.

## How Go uses it

- **JSON serialization**: field tags like `json:"name,omitempty"` control encoding/decoding.
- **Database models**: structs map to table rows; ORM tags specify column names.
- **API request/response types**: every endpoint has a struct shape.
- **Configuration**: structs group related config values.
- **Error wrapping**: `fmt.Errorf("... %w", err)` returns a struct implementing `error`.

## Go example

```go
package main

import (
    "encoding/json"
    "fmt"
)

type Product struct {
    ID    int     `json:"id"`
    Name  string  `json:"name"`
    Price float64 `json:"price,omitempty"`
}

func main() {
    // Zero-value struct
    var p1 Product
    fmt.Printf("zero: %+v\n", p1)

    // Named struct literal
    p2 := Product{ID: 1, Name: "Laptop", Price: 999.99}
    fmt.Printf("literal: %+v\n", p2)

    // Anonymous struct
    point := struct{ X, Y int }{X: 3, Y: 4}
    fmt.Printf("anonymous: %+v\n", point)

    // Struct comparison
    a := Product{ID: 1, Name: "Mouse"}
    b := Product{ID: 1, Name: "Mouse"}
    fmt.Println("a == b:", a == b)

    // JSON tags in action
    jsonBytes, _ := json.Marshal(p2)
    fmt.Println("json:", string(jsonBytes))

    // Structs are copied on assignment
    p3 := p2
    p3.Price = 0
    fmt.Printf("original unchanged: %.2f\n", p2.Price)
}
```

## Step-by-step execution

For `p2 := Product{ID: 1, Name: "Laptop", Price: 999.99}`:

1. Compiler reserves memory for a `Product` value (3 fields + padding).
2. Field `ID` (int) is set to `1`.
3. Field `Name` (string) is set to `"Laptop"` (string header: pointer + length).
4. Field `Price` (float64) is set to `999.99`.
5. The resulting value is assigned to `p2`.
6. Later, `json.Marshal(p2)` reads the `json` tag on each field to determine the output key name.

## Common mistakes

- Mistake: Forgetting that struct assignment copies all fields, so modifying one does not affect the other.
  - Fix: Use a pointer (`&Product{...}`) when you need shared access.

- Mistake: Assuming `==` works on structs with slice or map fields.
  - Why: Slices and maps are not comparable. `==` on a struct containing them is a compile error.
  - Fix: Use `reflect.DeepEqual` or write a custom `Equal` method.

- Mistake: Expecting positional struct literals to be self-documenting or refactor-safe.
  - Fix: Always use named field literals (`Field: value`) except in trivial cases.

- Mistake: Ignoring struct padding, leading to wasted memory in large arrays.
  - Fix: Order fields from largest to smallest alignment to minimize padding.

## Debugging walkthrough

Consider this code that fails to compile:

```go
type Data struct {
    Items []string
}
func main() {
    a := Data{Items: []string{"x"}}
    b := Data{Items: []string{"x"}}
    fmt.Println(a == b) // COMPILE ERROR
}
```

**Symptom**: `invalid operation: a == b (struct containing []string cannot be compared)`.

**Root cause**: `Data` contains a slice field `Items`, making the struct incomparable with `==`.

**Fix**: Use `reflect.DeepEqual(a, b)` or write a comparison method:

```go
func (d Data) Equal(other Data) bool {
    if len(d.Items) != len(other.Items) {
        return false
    }
    for i := range d.Items {
        if d.Items[i] != other.Items[i] {
            return false
        }
    }
    return true
}
```

## Production notes

- **Field tags** are the standard place for serialization metadata. Use `json`, `yaml`, `xml`, `db`, `validate`, `form`, etc. consistently.
- **Never embed a mutex by value** — embed `*sync.Mutex` or use a pointer field.
- **Zero-value structs** should be usable. Design your struct so that `var x T` is a valid initial state (the "zero value is useful" Go philosophy).
- **JSON omitempty** with `omitempty` drops zero-value fields from output. Use `pointer` fields when you need to distinguish "not set" from "zero".
- **`string` tag option** in JSON forces the field to be quoted: `json:"id,string"`.

## Performance implications

- **Copy cost**: Assigning or passing a large struct by value copies all bytes. For structs > ~100 bytes, prefer `*T` to avoid copying.
- **Padding**: Large arrays of structs waste memory and CPU cache if fields are poorly ordered. Order fields by decreasing alignment: `*T`, `int64`, `float64`, `int32`, `int16`, `bool`, `string` (pointer+length).
- **Tag reflection** is slow. Avoid `reflect` in hot paths. Use code generation (e.g., `easyjson`) for high-performance serialization.

## Practice task

Define a `Book` struct with fields `Title string`, `Author string`, `Pages int`, and a JSON tag `json:"-"` on `Pages` to omit it from JSON output. Write a `main()` that:
1. Creates a `Book` literal.
2. Prints the struct with `%+v`.
3. JSON-marshals it and prints the JSON string.
4. Creates two identical `Book` values (differing only in `Pages`) and compares them with `==`, printing the result.

## Tests / verification

```bash
go run ./curriculum/modules/05-types-interfaces-packages-modules/lessons/01-structs
go test ./curriculum/modules/05-types-interfaces-packages-modules/lessons/01-structs
```

## Review questions

1. What is the zero value of `type Point struct { X, Y int; Z *float64 }`?
2. Can you use `==` to compare two structs that both contain a `[]byte` field? Why or why not?
3. What is the difference between `p := Person{Name: "Alice"}` and `var p Person; p.Name = "Alice"`?
4. When would you use an anonymous struct instead of a named struct type?
5. How do field tags affect memory layout?

## NEXT UP

Methods — attaching behavior to your named types with receiver functions.
