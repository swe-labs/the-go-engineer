# Boolean logic

## Learning objective

Write and evaluate boolean expressions using `&&`, `||`, and `!`, apply short-circuit evaluation and De Morgan's laws, and use boolean logic correctly in control flow, nil checks, and boundary conditions.

## Why this matters

Boolean logic is the foundation of every conditional: `if`, `for`, `switch`, and `continue` all depend on boolean expressions. Bugs in boolean logic are notoriously hard to spot because they often pass initial testing and surface only in edge cases. Misplaced negation, wrong operator precedence, or misunderstanding short-circuit evaluation can lead to nil pointer dereferences, infinite loops, or silent data corruption. Mastering boolean logic is not theoretical — it directly determines whether your code handles real-world inputs correctly.

## Mental model

Think of boolean values as signals on a wire: `true` = current flowing, `false` = no current. Operators are logic gates:

- `&&` (AND) connects two signals in series: both must be on for current to flow.
- `||` (OR) connects them in parallel: at least one must be on.
- `!` (NOT) is an inverter: turns on to off, off to on.

Short-circuit evaluation is a physical property of these gates: in a series (`&&`), if the first switch is off, the circuit is already broken — the second switch is never checked. In parallel (`||`), if the first path conducts, the second is irrelevant.

The result of any boolean expression is always a `bool` — there is no truthy/falsy in Go. You cannot use an integer where a `bool` is expected, and vice versa.

## Core idea

```go
var a, b bool = true, false
fmt.Println(a && b)  // false
fmt.Println(a || b)  // true
fmt.Println(!a)      // false
```

Truth tables:

| A | B | A && B | A \|\| B | !A |
|---|---|--------|---------|----|
| T | T | T | T | F |
| T | F | F | T | F |
| F | T | F | T | T |
| F | F | F | F | T |

Operator precedence (high to low): `!` > `&&` > `||`. Like arithmetic, use parentheses to override.

Short-circuit evaluation rules:

- `expr1 && expr2`: if `expr1` is `false`, `expr2` is **not evaluated**.
- `expr1 || expr2`: if `expr1` is `true`, `expr2` is **not evaluated**.

This is not a performance optimization — it is a correctness guarantee relied upon by idioms like `p != nil && p.Name != ""`.

## Under the hood

The `bool` type is a single byte in memory, where `0` represents `false` and any non-zero value (conventionally `1`) represents `true`. The `&&` operator compiles to a conditional branch: if the first operand is `false`, jump past the second operand's evaluation. The `||` operator compiles similarly: if the first operand is `true`, jump past the second operand.

At the machine level, the compiler emits `CMP` (compare) and `Jcc` (conditional jump) instructions. Branch prediction in modern CPUs makes these jumps nearly free when the branch is predictable (e.g., nil checks on pointers that are almost always non-nil).

De Morgan's laws describe how to distribute negation over conjunction and disjunction:

- `!(A && B)` is equivalent to `!A || !B`
- `!(A || B)` is equivalent to `!A && !B`

These are logical identities, not Go-specific, but they help simplify complex negated conditions.

## How Go uses it

- **Nil safety**: `if p != nil && p.Valid()` — short-circuit prevents nil dereference.
- **Boundary checks**: `if i >= 0 && i < len(slice)` — safe index guard.
- **Input validation**: `if name == "" || age < 0` — reject invalid input.
- **Loop termination**: `for i < n && !done { ... }` — compound loop condition.
- **Configuration defaults**: `val := provided || default` — but Go does not have `||` for non-bool; use if/else.
- **De Morgan simplification**: `if !(status != "ok" && retries > 0)` rewritten as `if status == "ok" || retries == 0`.

## Go example

```go
package main

import "fmt"

func main() {
	// Truth table
	fmt.Println("AND truth table:")
	for _, a := range []bool{true, false} {
		for _, b := range []bool{true, false} {
			fmt.Printf("  %5t && %5t = %5t\n", a, b, a&&b)
		}
	}

	fmt.Println("\nOR truth table:")
	for _, a := range []bool{true, false} {
		for _, b := range []bool{true, false} {
			fmt.Printf("  %5t || %5t = %5t\n", a, b, a||b)
		}
	}

	// Short-circuit demonstration
	fmt.Println("\nShort-circuit:")
	fmt.Printf("  false && sideEffect() = %t (sideEffect not called)\n", false && sideEffect())
	fmt.Printf("  true || sideEffect()  = %t (sideEffect not called)\n", true || sideEffect())

	// De Morgan's laws
	a, b := true, false
	fmt.Println("\nDe Morgan's laws:")
	fmt.Printf("  !(a && b) == !a || !b: %t\n", !(a && b) == (!a || !b))
	fmt.Printf("  !(a || b) == !a && !b: %t\n", !(a || b) == (!a && !b))
}

func sideEffect() bool {
	fmt.Println("  sideEffect() WAS called!")
	return true
}
```

## Step-by-step execution

Trace `result := i >= 0 && i < len(slice) && slice[i] > 0` with `i = -1`, `slice = [5]int{1,2,3,4,5}`:

1. Evaluate `i >= 0`: `-1 >= 0` → `false`.
2. Short-circuit: `false && ...` evaluates to `false` without checking the next two conditions.
3. `result = false`. `slice[i]` is never accessed — no index-out-of-bounds panic.

Trace `ok := p == nil || p.Valid()` with `p = nil`:

1. Evaluate `p == nil`: `true`.
2. Short-circuit: `true || ...` evaluates to `true` without calling `p.Valid()`.
3. `ok = true`. Safe — `p.Valid()` is never called on nil.

Trace De Morgan transformation. Original:

```go
if !(status == "active" || retries < 3) { ... }
```

Applying `!(A || B) = !A && !B`:

```go
if status != "active" && retries >= 3 { ... }
```

Both are logically equivalent. The second form is often easier to read because the negation is distributed to each sub-condition.

## Common mistakes

- **Using `&&` instead of `||` (or vice versa)**: `if age < 0 && age > 150` — always false because no number is both < 0 and > 150. Should be `||`.

- **Relying on short-circuit for side effects**: `if count > 0 && incrementCounter()` — if `count` is 0, `incrementCounter` is never called. Use explicit statements for side effects.

- **Writing `if !x == true`**: `!x` is already a bool. Just write `if !x`.

- **Forgetting parentheses with `!` and `&&`/`||`**: `!a && b` means `(!a) && b`, not `!(a && b)`. The `!` operator binds tighter than `&&` and `||`.

- **Using bitwise `&` or `|` instead of `&&` or `||`**: `&` and `|` work on integers (bitwise) and do not short-circuit. Using them on booleans compiles but evaluates both sides unconditionally.

- **Negating a complex condition incorrectly**: `if !(a > 0 && b > 0)` becomes `if a <= 0 || b <= 0`, not `if a <= 0 && b <= 0`.

## Debugging walkthrough

```go
package main

import "fmt"

func main() {
	items := []string{"a", "b", "c"}
	i := 3
	if i >= 0 && i < len(items) && items[i] == "c" {
		fmt.Println("found c at index", i)
	}
}
```

**Symptom**: No output. Index 3 is out of bounds for a 3-element slice (valid indices 0-2). But there is no panic — short-circuit saves us.

**Investigation**: Trace with `i = 3`. `i >= 0` is `true`. `i < len(items)` is `3 < 3` = `false`. Short-circuit: `false && ...` → `false`. The `items[i]` access never executes.

**Fix**: Change `i` to 2, or adjust the condition to use the correct upper bound.

Now consider a common nil-check bug:

```go
type Config struct {
	Timeout int
}

func getTimeout(c *Config) int {
	if c != nil || c.Timeout > 0 {
		return c.Timeout
	}
	return 30
}
```

**Symptom**: Panic: nil pointer dereference on `c.Timeout`.

**Root cause**: `||` short-circuits on `true`. When `c` is nil, `c != nil` is `false`, so the right side `c.Timeout > 0` **is evaluated** — dereferencing nil.

**Fix**: Use `&&` instead of `||`:

```go
if c != nil && c.Timeout > 0 {
	return c.Timeout
}
```

## Production notes

- **Always check nil before dereference**: The Go proverb "Don't check for nil, provide a sensible zero value" applies only to your own APIs. When receiving a pointer from external code, guard with `p != nil && ...`.
- **Prefer positive conditions**: `if ok { ... }` is easier to read than `if !failed { ... }`. When a condition is complex, extract it into a well-named function: `func isValid(input string) bool`.
- **Use De Morgan to simplify**: Complex negations should be rewritten to avoid double negatives. `if !(x == 0)` → `if x != 0`.
- **Avoid deep nesting**: Combine conditions at the top level and return early. Prefer `if !condition { return }` over wrapping the entire function body in `if condition { ... }`.
- **Document non-obvious short-circuit dependencies**: If the evaluation order matters for correctness beyond the standard idioms, add a comment.

## Performance implications

- `&&` and `||` compile to conditional branches. Well-predicted branches cost ~1 cycle; mispredicted branches cost 10-20 cycles.
- Short-circuit evaluation can be a net performance win because it avoids evaluating expensive expressions (function calls, map lookups) when the result is already determined.
- There is no memory allocation in boolean expressions — `bool` values are stack-allocated or register-resident.
- The compiler may optimize `!` by inverting a condition flag rather than explicitly computing a boolean value.

## Practice task

Write a function `canAccess(role string, isAdmin, isActive bool, hour int) bool` that returns `true` if access is permitted:
- `isActive` must be `true`.
- `isAdmin` grants access at any hour.
- Non-admins can only access during business hours (9 <= hour < 17).
- `role` must be non-empty (`""` is rejected).
- If `role` is `"banned"`, access is denied regardless of other flags.

In `main()`, test at least 8 combinations covering denied-banned, active-admin, inactive-nonadmin-in-hours, and active-nonadmin-out-of-hours.

## Tests / verification

```bash
go run ./curriculum/modules/03-programming-fundamentals/lessons/08-boolean-logic
go test ./curriculum/modules/03-programming-fundamentals/lessons/08-boolean-logic
```

## Review questions

1. What is the result of `true && false || true` in Go? Show the step-by-step evaluation with precedence.
2. In `if p != nil && p.Name != ""`, why is it safe to use `&&` instead of `||`?
3. Apply De Morgan's law to `!(age >= 18 && country == "US")` and write the equivalent expression.
4. What is the difference between `&&` and `&` when used with `bool` operands?
5. Write a boolean expression that is `true` when a year is a leap year (divisible by 400, or divisible by 4 but not by 100).

## NEXT UP

If and else — how Go's conditional branching works with boolean expressions.
