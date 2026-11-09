# Questions

Answer each question with concrete reasoning and evidence.

## Question 1

```go
var x int
var y int = 10
z := 20
```

What is the value of `x`? Explain the difference between the three declaration forms and when you would use each.

## Question 2

```go
s := "Hello"
fmt.Println(s[0])
fmt.Println(len(s))
```

What does `s[0]` print? What does `len(s)` return? Why might these results surprise someone who expects `s[0]` to print `"H"` and `len(s)` to return `5`?

## Question 3

Explain the difference between a slice and an array in Go. What happens when you append to a slice beyond its capacity? Draw the memory layout of a slice header.

## Question 4

```go
m := map[string]int{"a": 1}
v := m["b"]
fmt.Println(v)
```

What does this program print? Why doesn't it crash even though `"b"` is not in the map? How would you distinguish between a key that exists with value `0` and a key that does not exist?

## Question 5

```go
items := []int{10, 20, 30, 40}
fmt.Println(append(items[:2], items[3:]...))
```

What does this program print? Explain the aliasing bug that can occur with overlapping slice operations.

## Question 6

```go
func main() {
    var p *int
    fmt.Println(*p)
}
```

What happens when you run this program? Why? How would you fix it?

## Question 7

Explain what a type conversion is in Go and when it is required. Why does `var x int = 3.14` fail to compile but `var x float64 = 3` succeed?

## Question 8

```go
const (
    A = iota
    B
    C
)
```

What are the values of `A`, `B`, and `C`? How does `iota` work and what is a practical use for it?

## Question 9

Write a function `Reverse(s string) string` that reverses a string. Explain why a byte-by-byte swap would fail for non-ASCII strings, and how your implementation handles that correctly.

## Question 10

What is a nil pointer dereference? Give three concrete scenarios in Go where a nil pointer panic can occur. For each scenario, show the code that would cause it and the fix.
