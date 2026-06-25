# Inspecting variables

## Learning objective

Use Delve's variable inspection commands — `print`, `locals`, `args`, `whatis` — to examine values, types, and structural relationships at any point during execution.

## Why this matters

Setting breakpoints and stepping through code lets you navigate the program. But the real payoff is examining variable state at each stop. A single `print` can confirm or refute a hypothesis about a bug. Without variable inspection, you are debugging blind — guessing what values might be based on log output. With it, you see the exact state of every variable, struct field, pointer, and slice at the moment execution pauses.

## Mental model

Each breakpoint hit freezes the program in a specific stack frame. You can query this frame like a snapshot: "what is the value of `x`?", "what are all the local variables?", "what type is this interface value?". The debugger reads the process memory directly, so the values are authoritative — they are not logs from a previous run or approximations.

## Core idea

Delve variable inspection commands:

| Command | Short | Action |
|---|---|---|
| `print <expr>` | `p` | Print the value of an expression |
| `locals` | — | Print all local variables in the current frame |
| `args` | — | Print all function arguments in the current frame |
| `whatis <expr>` | — | Print the type of an expression |
| `set <var> = <val>` | — | Change a variable's value (use with caution) |

Expressions can be:
- Simple variables: `p x`, `p name`
- Struct fields: `p user.Name`
- Pointer dereferences: `p *ptr`, `p ptr.Field` (Delve auto-dereferences)
- Slice elements: `p arr[0]`
- Composite expressions: `p len(arr)`, `p user.Name + "!"`

## Under the hood

When you type `print x`, Delve:
1. Looks up `x` in the symbol table of the current stack frame.
2. Determines its type from DWARF debug info (embedded in the binary by the compiler).
3. Reads the raw bytes of `x` from the process memory at the frame's stack offset (or heap address for pointers).
4. Formats the bytes according to the type and prints them.

For compound types (structs, slices, maps), Delve reads all constituent fields. For slices, it reads the underlying array pointer, length, and capacity. For interfaces, it reads the type descriptor and the data pointer.

## How Go uses it

- **Validate assumptions**: Before debugging, you might assume `x` is positive. `print x` might reveal it's `-1`. Hypothesis refuted.
- **Inspect complex types**: `print` works with structs, maps, slices, and channels. Delve shows all fields and elements.
- **Compare values across frames**: Use `frame <n>` to switch to a different stack frame, then `print` variables in that frame's scope.
- **Watch expressions**: Use `on <bp> p <expr>` to print an expression each time a breakpoint is hit without stopping.
- **Modify variables for testing**: `set x = 42` changes a variable mid-execution. Useful for testing error paths without restarting.

## Go example

```go
package main

import "fmt"

type User struct {
	ID    int
	Name  string
	Score float64
	Tags  []string
}

func main() {
	u := User{
		ID:    1,
		Name:  "Alice",
		Score: 95.5,
		Tags:  []string{"go", "debugging", "delve"},
	}
	processUser(u)
}

func processUser(u User) {
	bonus := computeBonus(u.Score)
	fmt.Printf("User %d (%s): final score %.1f\n", u.ID, u.Name, u.Score+bonus)
}

func computeBonus(score float64) float64 {
	if score > 90 {
		return 5.0
	}
	return 1.0
}
```

Debugging session:
```bash
dlv debug
(dlv) break main.processUser
(dlv) continue
(dlv) print u
(dlv) print u.Name
(dlv) print u.Tags
(dlv) print u.Tags[1]
(dlv) locals
(dlv) args
(dlv) whatis u
(dlv) whatis u.Score
(dlv) next
(dlv) print bonus
(dlv) continue
```

## Step-by-step execution

Break at `processUser`:

1. `print u` shows the full struct:
```
main.User {ID: 1, Name: "Alice", Score: 95.5, Tags: []string {"go", "debugging", "delve"}}
```
2. `print u.Tags[1]` → `"debugging"`.
3. `locals` shows all locals: `u` and `bonus` (bonus not yet assigned, shows zero value).
4. `whatis u.Score` → `float64`.
5. `next` executes the `computeBonus` call. `print bonus` → `5`.
6. `next` prints the output line. `continue` finishes.

## Common mistakes

- **Printing before variable is in scope**: If you break at `main.main`, `u` does not exist yet. Set the breakpoint after the variable is declared.
- **Assuming Delve knows your current frame**: `print x` shows `x` in the *current* frame. If you're in a library function (e.g. `fmt.Printf`), your local variables are not in scope. Use `frame 1` to switch to your code's frame.
- **Not distinguishing value vs pointer**: `print p` on a pointer shows the address. `print *p` dereferences it. Delve auto-dereferences for field access (`p.Field` works), but for assignment you need `*p`.
- **Modifying variables with `set` unintentionally**: `set` changes process memory. If you modify a loop counter, the loop may behave differently. Only use `set` when you intend to alter execution.
- **Large data structures**: `print` with a very large slice or map can produce thousands of lines. Delve truncates long output — use slice bounds like `arr[0:10]` to limit.

## Debugging walkthrough

Suppose `computeBonus` returns `1.0` instead of `5.0` for a score of `95.5`. We step through:

```
(dlv) break main.computeBonus
(dlv) continue
> main.computeBonus() ./main.go:29
(dlv) print score
95.5
(dlv) next              # evaluate if score > 90
(dlv) next              # should go into the if branch
(dlv) next              # return 5.0
(dlv) print score + 5   # confirm
100.5
(dlv) continue
```

The bug was not here — `computeBonus` worked correctly. The bug was *caller-side*: maybe `u.Score` was not what we thought. Go back and inspect the input.

## Production notes

- **Read-only debugging**: In production-adjacent debugging (core dumps, crash traces), you cannot use `set`. Core dumps are read-only snapshots. Use `print` and `locals` exclusively.
- **Remote debugging**: `print` over a remote Delve connection has the same capabilities but adds network latency. For large structs, the round trip may be noticeable.
- **Sensitive data**: `print` shows variable values including passwords, tokens, and PII. Be mindful when sharing debug output or attaching to a process that handles sensitive data.
- **DWARF size**: Debug symbols increase binary size by 2–5x. In CI/deployment, strip them with `-ldflags="-s -w"`. Keep them in development.

## Performance implications

- `print` reads process memory directly via ptrace / debug events. Reading a single int is ~1 µs. Reading a large struct (10KB+) may take ~10–100 µs.
- `locals` iterates over all variables in DWARF for the current frame. On a function with 100+ locals, this might take a few milliseconds.
- None of these commands affect the running program's performance — they only execute while the process is stopped.

## Practice task

Start `dlv debug` on this lesson's program. Set a breakpoint on `main.processUser`. Use `print` to inspect the `User` struct. Print individual fields. Use `locals` to see both the argument and local variables. Use `whatis` to check the type of `u.Tags`. Step into `computeBonus` and print `score` before and after the `if` condition. Use `set` to change `score` to `80.0` and step through to see the different return value.

## Tests / verification

```bash
go run ./curriculum/modules/06-testing-debugging-refactoring/lessons/17-inspecting-variables
go test ./curriculum/modules/06-testing-debugging-refactoring/lessons/17-inspecting-variables
```

## Review questions

1. How do you print the value of a local variable named `counter` in Delve?
2. What is the difference between `p x` when `x` is a pointer vs `p *x`?
3. What does the `locals` command show?
4. How do you determine the type of an expression at runtime in Delve?
5. Why might `print x` return "could not find symbol" even though `x` is declared in the source?

## NEXT UP

Debugging tests — using `dlv test` to step through test execution and debug test failures.
