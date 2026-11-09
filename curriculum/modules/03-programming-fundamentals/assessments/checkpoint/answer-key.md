# Answer Key

## Question 1

`x` is `0` — the zero value for `int`.

- `var x int` declares without initializing; x gets zero value.
- `var y int = 10` declares and explicitly initializes.
- `z := 20` is short declaration — type is inferred from the literal.

Use `var x int` for package-level declarations or when zero value is desired. Use `var y int = 10` when you need an explicit type but want to declare at package level. Use `z := 20` inside functions for concise initialization with inferred type.

## Question 2

`s[0]` prints `72` (the byte value of `'H'`). `len(s)` returns `5` (the byte length).

A beginner might expect `s[0]` to print `"H"` because other languages treat strings as character arrays. In Go, `s[0]` indexes bytes, not characters. For ASCII text the byte and character index align, but for multi-byte UTF-8 (e.g., `"é"` is 2 bytes) they diverge. Use `for i, r := range s` to iterate runes.

## Question 3

An array has a fixed length known at compile time (`[3]int`). A slice is a view into an underlying array (`[]int`).

When appending beyond capacity, Go allocates a new backing array (typically 2x the old capacity), copies existing elements, and returns a new slice header pointing to the new array.

Slice header (on 64-bit): pointer to backing array (8 bytes), length (8 bytes), capacity (8 bytes) = 24 bytes total. Multiple slices can share the same backing array, which is the source of aliasing bugs.

## Question 4

It prints `0`. Accessing a missing key in a Go map returns the zero value for the value type (0 for int). It doesn't crash.

To distinguish missing vs. zero, use the comma-ok idiom: `v, ok := m["b"]`. If `ok` is `false`, the key does not exist.

## Question 5

It prints `[10 20 40]`.

The aliasing bug: `items[:2]` and `items[3:]` both reference the same backing array. `append(items[:2], items[3:]...)` writes `40` into position 2 of the original backing array, overwriting `30`. The original `items` slice becomes `[10 20 40 40]`. This is a common source of subtle bugs — use `copy` or create a new slice to avoid mutating the source.

## Question 6

The program panics at runtime with "invalid memory address or nil pointer dereference".

`var p *int` declares p as a nil pointer (zero value for pointer types). Dereferencing `*p` attempts to read from address 0, which the OS prevents.

Fix: initialize the pointer before dereferencing — either point to an existing variable (`x := 42; p := &x`) or use `new` (`p := new(int); *p = 42`). Always check `if p != nil` before dereferencing.

## Question 7

Type conversion in Go explicitly changes a value from one type to another: `int(3.14)`. Go requires explicit conversions between different numeric types — there is no implicit numeric promotion.

`var x int = 3.14` fails because 3.14 is an untyped constant with a fractional part — Go's constant resolution allows assignment to float64 but not to int. `var x float64 = 3` succeeds because 3 is an untyped integer constant that can be represented exactly as float64. The rule: untyped constants are implicitly converted only when the conversion is lossless or when assigned to a compatible type.

## Question 8

`A = 0`, `B = 1`, `C = 2`.

`iota` is a predeclared identifier that increments within a `const` block. It resets to 0 for each new const block. Each line in a const block increments iota by 1.

Practical use: enumerating flags, bitmask values, or sequential constants like days of the week, HTTP status code categories, or log levels.

## Question 9

```go
func Reverse(s string) string {
    runes := []rune(s)
    for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
        runes[i], runes[j] = runes[j], runes[i]
    }
    return string(runes)
}
```

Byte-by-byte swap fails for non-ASCII strings because a byte is not a character. Multi-byte UTF-8 runes like `"é"` (2 bytes) or `"世"` (3 bytes) would be split and produce invalid UTF-8. Converting to `[]rune` first ensures each element is a full code point, so swapping is safe.

## Question 10

A nil pointer dereference occurs when a program attempts to read or write through a pointer that is `nil`.

**Scenario 1 — Declaration without initialization:**
```go
var p *int
*p = 42 // panic
```
Fix: `p := new(int)` or `var x int; p := &x`

**Scenario 2 — Map access through nil map:**
```go
var m map[string]int
m["key"] = 1 // panic — nil map assignment
```
Fix: `m := make(map[string]int)` — reading from nil map returns zero value, but writing panics.

**Scenario 3 — Nil pointer receiver:**
```go
type T struct { V int }
func (t *T) Method() { fmt.Println(t.V) }
var t *T
t.Method() // panic — t is nil, but Method accesses t.V
```
Fix: check `if t == nil` at the start of pointer receiver methods, or ensure the pointer is initialized before calling.
