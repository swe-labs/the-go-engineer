# Value receivers

## Learning objective

Write value-receiver methods with predictable copy semantics, explain when value receivers are preferred over pointer receivers, and analyze the performance trade-offs of copying large structs.

## Why this matters

Value receivers are the default in Go for a reason: they make immutability guarantees explicit and keep the call stack clean. Understanding when a value receiver is correct — and when it is wasteful — separates competent Go developers from great ones. Every Go engineer makes this choice daily.

## Mental model

A value receiver is a copy. Inside the method, `c` is an independent copy of the caller's value. Any writes affect only the copy. When the method returns, the copy is discarded. From the caller's perspective, the original is unchanged. This is like passing a document to a reviewer: they can mark it up, but you keep the original.

## Core idea

Syntax:

```go
type Point struct {
    X, Y float64
}

func (p Point) DistanceToOrigin() float64 {
    return math.Sqrt(p.X*p.X + p.Y*p.Y)
}
```

**Copy semantics**: `p` inside the method is a fresh copy. The caller's `Point` is unaffected.

**Immutability guarantee**: Value receivers cannot mutate the caller's data. This is useful for accessors, computations, and types that should be immutable (e.g., `time.Time`).

**When to use value receiver**:

1. Method does not mutate the receiver.
2. Receiver is small (all fields fit in a few machine words).
3. Type should be immutable by design (e.g., `Color`, `Money`).
4. No pointer methods exist on the type (consistency).

**Large struct copy cost**: Copying a large struct (hundreds of bytes) on every method call adds CPU and memory pressure. For hot-path code, measure and consider switching to a pointer receiver.

## Under the hood

A value receiver method is compiled to a function whose first parameter is the value itself: `func Point_DistanceToOrigin(p Point) float64`. The caller copies the entire struct onto the stack (or into registers for small structs). In Go's calling convention, small structs (up to ~2 machine words) may be passed in registers on modern architectures, making them nearly free.

For a struct `type Large struct { buf [1024]byte }`, calling a value receiver method copies 1024 bytes onto the call stack. This is fast for a single call but adds measurable overhead in tight loops.

## How Go uses it

- **time.Time** methods use value receivers: `t.Add(d)` returns a new `Time`, never mutating `t`.
- **math/big** types use pointer receivers because they are large (multiple `Word` slices).
- **Color types** like `image/color.RGBA` use value receivers.
- **Accessor methods**: `func (p Person) Age() int`.
- **String() and Error()** methods typically use value receivers.

## Go example

```go
package main

import (
	"fmt"
	"math"
)

type Point struct {
	X, Y float64
}

func (p Point) DistanceToOrigin() float64 {
	return math.Sqrt(p.X*p.X + p.Y*p.Y)
}

func (p Point) Translated(dx, dy float64) Point {
	return Point{X: p.X + dx, Y: p.Y + dy}
}

func main() {
	p := Point{X: 3, Y: 4}
	fmt.Printf("distance: %.2f\n", p.DistanceToOrigin())

	p2 := p.Translated(1, 2)
	fmt.Printf("original: (%.1f, %.1f)\n", p.X, p.Y)
	fmt.Printf("translated: (%.1f, %.1f)\n", p2.X, p2.Y)

	// Value receiver on a pointer automatically dereferences
	pp := &Point{X: 5, Y: 12}
	fmt.Printf("pointer distance: %.2f\n", pp.DistanceToOrigin())
}
```

## Step-by-step execution

For `p.Translated(1, 2)` where `p = Point{X: 3, Y: 4}`:

1. Caller copies `p` onto the stack (or into registers): `{X: 3, Y: 4}`.
2. Inside `Translated`, the copy's `X` is `3`, `Y` is `4`.
3. `p.X + dx` = `3 + 1` = `4`.
4. `p.Y + dy` = `4 + 2` = `6`.
5. New `Point{X: 4, Y: 6}` is returned.
6. Caller's `p` remains `{X: 3, Y: 4}`.
7. `p2` receives the returned value `{X: 4, Y: 6}`.

For `pp.DistanceToOrigin()` where `pp = &Point{X: 5, Y: 12}`:

1. Go sees method with value receiver on `*Point`.
2. Go dereferences `pp`: copies `Point{X: 5, Y: 12}` onto the stack.
3. Method works on the copy.
4. Original `pp` is unchanged.

## Common mistakes

- Mistake: Expecting a value receiver to modify the original struct.
  - Fix: Use a pointer receiver if mutation is needed.

- Mistake: Using value receivers on large structs in performance-sensitive code.
  - Fix: Benchmark with both approaches. For structs > 4 machine words, pointer receivers are usually faster.

- Mistake: Mixing value and pointer receivers on the same type without a clear rule.
  - Fix: If any method needs a pointer receiver, make all methods use pointer receivers.

- Mistake: Copying a mutex by using a value receiver on a struct containing `sync.Mutex`.
  - Fix: Always use a pointer receiver for types that contain a mutex.

## Debugging walkthrough

This code silently fails to update the struct:

```go
type Config struct {
    Timeout int
}

func (c Config) SetTimeout(t int) {
    c.Timeout = t
}

func main() {
    cfg := Config{}
    cfg.SetTimeout(30)
    fmt.Println(cfg.Timeout) // prints 0!
}
```

**Symptom**: `SetTimeout` appears to do nothing — `cfg.Timeout` stays `0`.

**Root cause**: `SetTimeout` has a value receiver. `c` inside the method is a copy. Setting `c.Timeout = t` modifies only the copy, which is discarded.

**Fix**: Change to pointer receiver: `func (c *Config) SetTimeout(t int)`.

## Production notes

- **Immutability by default** is a design advantage. Prefer value receivers for types that represent fixed values (points, colors, amounts).
- **Consistency within a package**: If a type has any pointer method, all methods conventionally use pointer receivers.
- **`go vet` warns** about unusual receiver patterns but is not exhaustive.
- **Method on small structs** is the cheapest calling convention in Go. Do not prematurely optimize with pointers.

## Performance implications

- **Small struct copy** (<= 32 bytes): Value receiver is fast, often passed in registers. No heap allocation.
- **Medium struct copy** (32-128 bytes): Copy cost is measurable. Pointer receiver avoids the copy but adds a dereference.
- **Large struct copy** (> 128 bytes): Pointer receiver is significantly cheaper. The copy also stresses the CPU cache.
- **Escape analysis**: A value receiver rarely causes heap escape. A pointer receiver can if the pointer is stored in a heap-allocated location.

## Practice task

Define a `Temperature` struct with a `Celsius float64` field. Add value receiver methods `ToFahrenheit() float64`, `ToKelvin() float64`, and `IsFreezing() bool` (true if Celsius <= 0). Do not add any pointer methods. In `main()`, create several `Temperature` values and print their conversions. Verify that calling `IsFreezing()` on a value does not let you change the original.

## Tests / verification

```bash
go run ./curriculum/modules/05-types-interfaces-packages-modules/lessons/04-value-receivers
go test ./curriculum/modules/05-types-interfaces-packages-modules/lessons/04-value-receivers
```

## Review questions

1. What is the key semantic difference between value and pointer receivers?
2. When is it safe (and idiomatic) to use a value receiver on a large struct?
3. Does Go automatically dereference a pointer when calling a value-receiver method? If so, how?
4. Why is mixing value and pointer receivers on the same type generally discouraged?
5. How does value receiver copying affect escape analysis and heap allocation?

## NEXT UP

Receiver sets — how Go determines whether a type satisfies an interface based on its method sets.
