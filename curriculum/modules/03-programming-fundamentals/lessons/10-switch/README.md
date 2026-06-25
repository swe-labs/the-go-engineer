# Switch

## Learning objective

Write switch statements with expression switches, case lists, default clauses, fallthrough, type switches, and expressionless switches to replace long if/else chains.

## Why this matters

A switch statement is Go's clean alternative to long if/else-if chains. It improves readability, makes the compiler's job easier (jump-table optimization), and provides features unavailable in if chains: multiple case values, type dispatch, and fallthrough. Type switches are especially valuable in Go interfaces, where the concrete type behind an interface value determines behavior. Professional Go code uses switch statements for state machines, command dispatch, type-based serialization, and error categorization.

## Mental model

Think of an expression switch as a multi-way junction: the expression at the top selects one path. Each `case` label is a possible value. The switch enters the first matching case, executes its body, and then _automatically exits_ — no `break` needed. This is the opposite of C, where you must `break` explicitly.

A type switch is similar, but instead of matching values, it matches the dynamic type of an interface value.

An expressionless switch (`switch { ... }`) is a sequence of boolean conditions evaluated top-to-bottom — a cleaner if/else-if chain.

## Core idea

```go
switch expr {
case val1, val2:
	// execute if expr == val1 or expr == val2
case val3:
	// execute if expr == val3
default:
	// execute if no case matches
}
```

Key properties:

- Cases are evaluated top to bottom. First match wins.
- `case` lists can contain multiple values separated by commas.
- `default` matches anything and goes last by convention (but can be placed anywhere).
- No implicit fallthrough — execution stops at the end of each case.
- `fallthrough` explicitly continues to the next case (rarely needed).
- The switch expression may be omitted (`switch { ... }`), making each `case` a boolean expression.
- A type switch (`switch v := x.(type)`) matches the dynamic type of an interface.

## Under the hood

The Go compiler analyzes switch statements for optimization:

- **Jump tables**: When cases are consecutive integers with few gaps, the compiler emits a jump table — O(1) dispatch. This is common for `iota`-based constants.
- **Binary search**: For sparse or non-integer cases (strings), the compiler emits a binary search tree optimized for the specific case set.
- **If/else chain**: For expressionless switches or complex case expressions, the compiler emits sequential comparisons.
- **Type switch**: Implemented as a type assertion chain. Each case checks if the interface's dynamic type matches using the runtime's type descriptor comparison.

The absence of automatic fallthrough is a deliberate design choice. Go's creators observed that fallthrough in C was the default and caused more bugs than convenience. In Go, fallthrough is explicit and visible.

## How Go uses it

- **Command dispatch**: `switch cmd { case "start": ... case "stop": ... }` in CLI tools.
- **Error type checking**: `switch e := err.(type)` to handle different error implementations.
- **State machines**: `switch state { case idle: ... case running: ... }`.
- **Type-based serialization**: `switch v := val.(type) { case int: ... case string: ... }` in JSON encoders.
- **HTTP route matching**: `switch r.Method { case "GET": ... case "POST": ... }`.
- **Multi-error handling**: `switch { case errors.Is(err, io.EOF): ... case errors.Is(err, os.ErrNotExist): ... }`.
- **Zero-value guards**: `switch x := someFunc.(type) { case nil: ... default: ... }`.

## Go example

```go
package main

import "fmt"

func describe(v interface{}) string {
	switch v := v.(type) {
	case nil:
		return "nil"
	case int:
		return fmt.Sprintf("int %d", v)
	case float64:
		return fmt.Sprintf("float64 %f", v)
	case string:
		return fmt.Sprintf("string %q (len %d)", v, len(v))
	case bool:
		if v {
			return "true"
		}
		return "false"
	case []int:
		return fmt.Sprintf("[]int len=%d cap=%d", len(v), cap(v))
	default:
		return fmt.Sprintf("unknown type %T", v)
	}
}

func main() {
	fmt.Println(describe(nil))
	fmt.Println(describe(42))
	fmt.Println(describe(3.14))
	fmt.Println(describe("hello"))
	fmt.Println(describe(true))
	fmt.Println(describe([]int{1, 2, 3}))
	fmt.Println(describe(struct{}{}))

	// Expressionless switch as if/else chain
	score := 85
	switch {
	case score >= 90:
		fmt.Println("A")
	case score >= 80:
		fmt.Println("B")
	case score >= 70:
		fmt.Println("C")
	case score >= 60:
		fmt.Println("D")
	default:
		fmt.Println("F")
	}

	// Fallthrough example
	i := 2
	switch i {
	case 1:
		fmt.Println("one")
	case 2:
		fmt.Println("two")
		fallthrough
	case 3:
		fmt.Println("three")
	default:
		fmt.Println("other")
	}
}
```

## Step-by-step execution

Trace `describe(42)`:

1. `v.(type)` extracts the dynamic type of `v`: `int`.
2. Compiler generates a type switch: compare `int` against each case type in order.
3. Case `nil`: no match (`int != nil`).
4. Case `int`: match! Execute body.
5. `v` is now typed as `int` with value `42`.
6. `fmt.Sprintf("int %d", v)` → `"int 42"`.

Trace expressionless switch with `score = 85`:

1. `case score >= 90`: `85 >= 90` → `false`. Next case.
2. `case score >= 80`: `85 >= 80` → `true`. Execute body: print `"B"`.
3. Switch exits — remaining cases are skipped.

Trace `fallthrough`:

1. `i = 2`. `case 1`: no match.
2. `case 2`: match. Print `"two"`. Then `fallthrough` forces execution into case 3.
3. `case 3`: print `"three"`. No more case bodies. Switch exits.

## Common mistakes

- **Expecting C-style automatic fallthrough**: Go stops at the end of each case. If you want fallthrough, write `fallthrough` explicitly.

- **Using `fallthrough` on non-consecutive cases**: `fallthrough` skips the next case's condition check but does not re-evaluate it. `case 2: fallthrough; case 5:` prints case 5's body after case 2 matches, which is usually not intended.

- **Forgetting `default` for unexpected inputs**: A switch without `default` silently does nothing on unmatched input. Often you want to log or handle the unexpected case.

- **Using `break` unnecessarily**: Go switch cases already break automatically. `break` is only needed to exit an enclosing loop or to break out of a `select`.

- **Confusing type switch with expression switch**: `switch v := x.(type)` requires the `.(type)` syntax. Forgetting it turns the switch into an expression switch on `x`'s type, which compiles but does something different.

- **Shadowing `v` in type switch**: `switch v := x.(type)` declares `v` with the matched type in each case. The original variable is not modified.

## Debugging walkthrough

```go
package main

import "fmt"

func main() {
	day := 3
	switch day {
	case 1:
		fmt.Println("Monday")
	case 2:
		fmt.Println("Tuesday")
	case 3:
		fmt.Println("Wednesday")
	case 4:
		fmt.Println("Thursday")
	case 5:
		fmt.Println("Friday")
	case 6:
		fmt.Println("Saturday")
	case 7:
		fmt.Println("Sunday")
	}
}
```

**Symptom**: No output for `day = 3` — it should print "Wednesday".

**Investigation**: Check the values. The switch works correctly; the output is "Wednesday". If output is missing, the value might be outside 1-7.

**Bug version**:

```go
day := 3
switch day:
case 1:
	fmt.Println("Monday")
```

**Symptom**: Compile error — `switch day:` should be `switch day {` (curly brace, not colon).

Now consider:

```go
func classify(x interface{}) string {
	switch x.(type) {
	case int:
		return "int"
	case string:
		return "string"
	}
	return "unknown"
}
```

**Symptom**: Works, but the matched value is not accessible inside the case. To use the typed value, use `switch v := x.(type)`.

## Production notes

- **Use expressionless switch for long if/else-if chains**: `switch { case x < 0: ... case x == 0: ... default: ... }` is more readable than `if x < 0 {} else if x == 0 {} else {}`.
- **Keep cases short**: Each case body should be 1-5 lines. Extract complex logic into named functions.
- **Type switch with error wrapping**: In Go 1.13+, use `errors.As` inside a type switch for wrapped errors: `switch { case errors.As(err, &target): ... }`.
- **Avoid `fallthrough` in most code**: It often indicates confusion. If you need fallthrough logic, consider restructuring.
- **Order cases from most specific to most general**: Place `default` last, narrowest cases first in type switches to avoid shadowing by broader types.
- **Prefer switch over if/else for equality comparisons**: It is faster, clearer, and the compiler can optimize it.

## Performance implications

- Integer switch on dense, consecutive values: jump table O(1). Same as a computed goto.
- String switch on small sets: binary search O(log n).
- Expressionless switch: linear scan O(n) — each case evaluated sequentially.
- Type switch: linear scan over the case types O(n). The compiler could optimize using a hash of type pointers but currently does not.
- Switch with many cases (50+) may be slower than a map lookup. Profile before micro-optimizing.

## Practice task

Write a function `calc(a, b float64, op string) (float64, error)` that uses a switch statement on `op`:
- `"add"`, `"sum"`, `"+"` → return `a + b`.
- `"sub"`, `"diff"`, `"-"` → return `a - b`.
- `"mul"`, `"prod"`, `"*"` → return `a * b`.
- `"div"`, `"/"` → return `a / b`, error if `b == 0`.
- Anything else → error `"unknown operator: <op>"`.

Then write a `describe` function using a type switch that accepts any value and returns:
- `"int: <value>"` for int types.
- `"float: <value>"` for float64 types.
- `"string: <value>"` for string types.
- `"bool: <value>"` for bool types.
- `"unknown type: <T>"` for everything else.

## Tests / verification

```bash
go run ./curriculum/modules/03-programming-fundamentals/lessons/10-switch
go test ./curriculum/modules/03-programming-fundamentals/lessons/10-switch
```

## Review questions

1. Does Go's `switch` require a `break` at the end of each case to prevent fallthrough? Why or why not?
2. What is the difference between `switch x.(type)` and `switch v := x.(type)`?
3. What happens if no case matches and there is no `default` clause?
4. Write a switch statement equivalent to `if x > 0 { "positive" } else if x < 0 { "negative" } else { "zero" }`.
5. What is the output of `switch 2 { case 1: fmt.Print("a"); fallthrough; case 2: fmt.Print("b"); fallthrough; case 3: fmt.Print("c") }`?

## NEXT UP

Loops — Go's `for` keyword as the only looping construct, covering all iteration patterns.
