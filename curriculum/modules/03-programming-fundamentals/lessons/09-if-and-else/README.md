# If and else

## Learning objective

Write and reason about Go if/else chains including short-statement form, boolean conditions, nesting, and readability patterns, and predict the execution path for any conditional structure.

## Why this matters

Conditional branching is the fundamental mechanism for writing programs that make decisions. Every Go service uses if/else for error handling (`if err != nil`), input validation, feature flags, request routing, authorization checks, and business logic. Without mastering if/else, you cannot write a single production function that handles the real-world complexity of varying inputs, edge cases, and error states. Professional Go code is dominated by if statements; learning to write them clearly and correctly is essential.

## Mental model

An if statement is a fork in a path. Execution reaches the fork, evaluates exactly one boolean condition, and proceeds down exactly one branch (or skips all branches if none match). Else attaches an alternative branch. Else if chains multiple forks in sequence: the first true condition wins, and the rest are skipped entirely.

Think of a series of locked doors: you walk down a hallway and try each door. The first door that is unlocked, you enter. If no door is unlocked, you either enter the else room or walk to the end of the hallway.

The short-statement form (`if err := doSomething(); err != nil`) is like a prep table before the door: you compute an intermediate result, inspect it, and decide which door to enter — all in one expression. The variable declared in the short statement is only visible within that if/else block, like a temporary note you crumple up and discard after leaving the room.

## Core idea

An if statement evaluates a boolean expression and conditionally executes a block:

```go
if condition {
    // executes when condition is true
}
```

The condition must be a boolean expression — no truthy/falsy coercion like in other languages. You can add an else branch:

```go
if condition {
    // true branch
} else {
    // false branch
}
```

Else if chains connect multiple conditions:

```go
if cond1 {
    // cond1 true
} else if cond2 {
    // cond1 false and cond2 true
} else {
    // both false
}
```

Go's special short-statement form scopes a variable to the conditional block:

```go
if x := compute(); x > 0 {
    // x is visible here
    fmt.Println(x)
} else {
    // x is also visible here
    fmt.Println(-x)
}
// x is NOT visible here — compile error
```

This is not a quirk; it is a deliberate design to limit variable scope and prevent accidental reuse.

## Under the hood

The Go compiler translates if/else into conditional jump instructions in the intermediate representation. At the machine level:

1. The compiler evaluates the boolean condition into a single register or flag value.
2. It emits a conditional jump (e.g., `JNZ` — jump if not zero) that skips the if-block when false.
3. For if/else, it adds an unconditional jump at the end of the if-block to skip the else-block.
4. For else-if chains, this pattern repeats: each condition is tested, the block is entered on true, and an unconditional jump skips the remaining chain.

The short-statement form is syntactic sugar: the compiler allocates the short-statement variable on the stack (or in a register for small types), initializes it, evaluates the condition, and reclaims the variable's scope at the closing brace — no additional runtime cost compared to declaring the variable before the if.

Boolean conditions are never short-circuited differently than you would expect: `&&` and `||` follow the same short-circuit rules as standalone boolean expressions, stopping evaluation as soon as the result is determined.

## How Go uses it

Go deliberately omits ternary operators (`?:`) and truthy/falsy coercion. Every condition must be an explicit boolean expression. This choice forces clarity: you write `if x > 0` rather than `if x`.

The standard library makes heavy use of if with short statement, especially in type assertions and error handling:

```go
if v, ok := m["key"]; ok {
    // use v
}
```

Idiomatic Go uses `if err != nil` as the dominant error-handling pattern — there is no try/catch. The entire standard library and most third-party code follows this pattern:

```go
if err := json.Unmarshal(data, &target); err != nil {
    return fmt.Errorf("unmarshal failed: %w", err)
}
```

The `go vet` tool detects certain if/else pitfalls like `if err != nil { return err }` followed by continuing execution (unreachable code warning).

## Go example

```go
package main

import (
	"errors"
	"fmt"
	"strconv"
)

func classifyTemperature(temp int) string {
	if temp > 40 {
		return "dangerously hot"
	} else if temp > 30 {
		return "hot"
	} else if temp > 20 {
		return "warm"
	} else if temp > 10 {
		return "mild"
	} else if temp >= 0 {
		return "cool"
	}
	return "freezing"
}

func divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, errors.New("division by zero")
	}
	return a / b, nil
}

func validateAge(s string) error {
	if age, err := strconv.Atoi(s); err != nil {
		return fmt.Errorf("invalid age: %w", err)
	} else if age < 0 || age > 150 {
		return fmt.Errorf("age %d out of range [0, 150]", age)
	} else if age < 18 {
		return fmt.Errorf("age %d is under 18", age)
	}
	return nil
}

func main() {
	for _, t := range []int{45, 32, 22, 15, 5, -5} {
		fmt.Printf("%3dC -> %s\n", t, classifyTemperature(t))
	}

	if result, err := divide(10, 3); err != nil {
		fmt.Println("divide error:", err)
	} else {
		fmt.Printf("10 / 3 = %.2f\n", result)
	}

	ages := []string{"25", "abc", "200", "16"}
	for _, s := range ages {
		if err := validateAge(s); err != nil {
			fmt.Printf("age %q invalid: %v\n", s, err)
		} else {
			fmt.Printf("age %q valid\n", s)
		}
	}
}
```

## Step-by-step execution

For `classifyTemperature(15)`:

1. Condition `temp > 40` → `15 > 40` → `false`. Skip block.
2. `else if temp > 30` → `15 > 30` → `false`. Skip block.
3. `else if temp > 20` → `15 > 20` → `false`. Skip block.
4. `else if temp > 10` → `15 > 10` → `true`. Enter block, return `"mild"`.
5. Remaining else-if and else are skipped; function returns `"mild"`.

For `classifyTemperature(-5)`:

1. `-5 > 40` → `false`.
2. `-5 > 30` → `false`.
3. `-5 > 20` → `false`.
4. `-5 > 10` → `false`.
5. `-5 >= 0` → `false`.
6. No branch matched; function falls through to return `"freezing"`.

For `validateAge("abc")` with short statement:

1. Evaluate `age, err := strconv.Atoi("abc")`. `err` is non-nil.
2. `if err != nil` → `true`. Enter block, return error `"invalid age: ..."`.
3. The `else if` and final `else` blocks are never reached.

## Common mistakes

- Mistake: Putting `else` on a new line: `if x { ... } \n else { ... }`.
  - Why: Go's automatic semicolon insertion places a semicolon after the closing `}` of the if-block, making `else` a standalone statement which is illegal.
  - Fix: Always write `} else {` on the same line.

- Mistake: Using `if x = 5` (assignment) instead of `if x == 5` (comparison).
  - Why: `=` is assignment, `==` is equality. Assignment in a condition is not allowed in Go (unlike C), so this is a compile error, not a silent bug.
  - Fix: Use `==` for comparison. If you intend to assign and test, use the short-statement form: `if x := 5; x > 0`.

- Mistake: Over-nesting when an early return would be clearer.
  - Why: Deep nesting reduces readability and makes it easy to miss the control flow.
  - Fix: Prefer early returns: `if err != nil { return err }` instead of `if err == nil { ... }`.

- Mistake: `else if` chain where a switch statement would be clearer.
  - Why: Chains of 4+ else-if branches are harder to scan than a switch.
  - Fix: Use `switch` for multi-way branches with different values of the same expression.

- Mistake: Referencing short-statement variable outside the if/else block.
  - Why: The variable scope is limited to the if/else block.
  - Fix: Declare the variable before the if statement if you need it afterward.

## Debugging walkthrough

Consider this broken code:

```go
package main

import "fmt"

func main() {
	grade := 85

	if grade >= 90 {
		fmt.Println("A")
	} else if grade >= 80 {
		fmt.Println("B")
	} else if grade >= 70 {
		fmt.Println("C")
	}
	if grade >= 60 {
		fmt.Println("D")
	} else {
		fmt.Println("F")
	}
}
```

**Symptom**: Input 85 prints `B` and `D` instead of just `B`.

**Investigation**: The developer intended a single if/else if/else chain. However, the third `} else if` block uses a separate `if` statement instead of `else if`. The `if grade >= 60` is an independent statement, not part of the chain. Since `85 >= 60`, both `B` and `D` print.

**Root cause**: Mixing `if` and `else if` in the same conceptual chain but using separate `if` statements for the last branch. Each independent `if` evaluates its condition regardless of prior matches.

**Fix**: Use a single if/else if/else chain:

```go
if grade >= 90 {
	fmt.Println("A")
} else if grade >= 80 {
	fmt.Println("B")
} else if grade >= 70 {
	fmt.Println("C")
} else if grade >= 60 {
	fmt.Println("D")
} else {
	fmt.Println("F")
}
```

Now input 85 matches `grade >= 80`, prints `B`, and skips all remaining branches.

## Production notes

In real Go services, if/else patterns dominate error handling. Key patterns:

- **Guard clauses**: Check error conditions first and return early. This flattens the happy path and separates error handling from business logic.
- **Feature flags**: `if os.Getenv("ENABLE_NEW_PIPELINE") == "true"` toggles between old and new code paths. Use this during gradual rollouts.
- **Input validation**: Chain `if` checks for each validation rule, returning a descriptive error for the first violation. This gives users actionable feedback.
- **Avoid deep nesting**: At 3+ levels of indentation, extract the inner block into a named function. Code reviews should flag functions with indentation depth > 4.
- **Short-statement for resource cleanup**: `if f, err := os.Open(file); err != nil { return err } else { defer f.Close(); ... }` limits f's scope to where it is used.

## Performance implications

- If/else chains with simple integer comparisons are essentially free — a few CPU cycles per condition.
- Long else-if chains (10+ branches) may be slower than a switch statement for dense integer values, as the compiler may optimize switch into a jump table.
- Boolean conditions with `&&` short-circuit: put the cheapest or most-likely-to-fail condition first to avoid unnecessary evaluations.
- The short-statement form has zero overhead compared to declaring the variable before the if. The compiler generates identical machine code.
- Deeply nested if blocks can cause branch mispredictions on the CPU pipeline, but this is rarely a measurable concern at the scale of application code.

## Practice task

Write a function `gradeClassification(score int) string` that:
- Returns `"A"` for scores 90-100.
- Returns `"B"` for scores 80-89.
- Returns `"C"` for scores 70-79.
- Returns `"D"` for scores 60-69.
- Returns `"F"` for scores 0-59.
- Returns `"invalid"` for scores outside 0-100.

Then write a `main()` that tests scores of `-5, 0, 45, 70, 85, 95, 100, 101` and prints each result. Run and verify every case.

## Tests / verification

```bash
go run ./curriculum/modules/03-programming-fundamentals/lessons/09-if-and-else
go test ./curriculum/modules/03-programming-fundamentals/lessons/09-if-and-else
```

## Review questions

1. Explain: What does `if x := f(); x > 0 { return x }` do differently from `x := f(); if x > 0 { return x }`?

2. Apply: Write an if/else chain that prints "positive", "negative", or "zero" for a given integer `n`, but using as few lines as possible.

3. Debug: `if a > b { max = a } else { max = b }` correctly finds the max of two numbers. What happens if you write `if a > b max = a else max = b` (no braces)? Will it compile?

4. Tradeoff: When would you choose `if x == 1 { ... } else if x == 2 { ... } else if x == 3 { ... }` instead of `switch x { case 1: ...; case 2: ...; case 3: ... }`? When would you choose the opposite?

5. Explain: Why does Go use `if err != nil { return err }` as the primary error-handling pattern instead of try/catch? What tradeoff does this design make?

## NEXT UP

Switch — Go's powerful multi-way branching construct with case lists, fallthrough, type switches, and expressionless forms.
