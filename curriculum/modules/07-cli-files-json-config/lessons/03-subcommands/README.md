# Subcommands

## Learning objective

Build multi-command CLIs using the subcommand pattern with independent `flag.FlagSet` per command, command routing via `os.Args`, and proper usage messages.

## Why this matters

Tools like `git`, `docker`, `kubectl`, and `go` itself are multi-command CLIs: `git commit`, `git push`, `git log`. Each subcommand has its own flags and behavior. The subcommand pattern is the standard architecture for every non-trivial CLI tool in Go.

## Mental model

A subcommand router is a switch statement over `os.Args[1]`. Each case creates its own `flag.FlagSet` with flags specific to that subcommand. The flag sets are independent — `-a` in `git add` means something different from `-a` in `git commit`.

```
os.Args: mytool greet -name=Alice
         ^^^^^^ ^^^^^ ^^^^^^^^^^^
         prog   cmd   cmd-specific flags
```

## Core idea

The subcommand pattern has three parts:

1. **Router** — inspect `os.Args[1]` to identify the subcommand name.
2. **Per-command FlagSet** — create a `flag.NewFlagSet("name", ...)` with flags for that subcommand.
3. **Parsing** — call `fs.Parse(os.Args[2:])` to parse only the args after the subcommand name.

A `flag.FlagSet` is an independent flag parser. Unlike the global `flag.CommandLine`, it does not pollute the global namespace and does not call `os.Exit(2)` on errors unless you use `flag.ExitOnError`.

## Under the hood

`flag.NewFlagSet(name, errorHandling)` creates a new `FlagSet` with its own internal map of flag definitions. The `name` is used in usage messages. The `errorHandling` parameter is one of:

- `flag.ContinueOnError` — return error from `Parse()`, do not exit.
- `flag.ExitOnError` — print error and `os.Exit(2)`.
- `flag.PanicOnError` — panic on parse error (used in tests).

The route is a simple string comparison. There is no registry or reflection — just a switch statement or a map of `string -> func`.

## How Go uses it

- Create a `flag.FlagSet` per subcommand with `flag.NewFlagSet("cmdname", flag.ExitOnError)`.
- Define flags on the set with `fs.String()`, `fs.Int()`, etc.
- Parse with `fs.Parse(os.Args[2:])`.
- Access flag values by dereferencing pointers.
- Use `fs.Args()` for positional args after the subcommand's flags.
- Handle the "no subcommand" case with a help/usage message.

## Go example

```go
package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: mytool <command> [flags]")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "greet":
		greetCmd := flag.NewFlagSet("greet", flag.ExitOnError)
		name := greetCmd.String("name", "World", "name to greet")
		greetCmd.Parse(os.Args[2:])
		fmt.Printf("Hello, %s!\n", *name)

	case "add":
		addCmd := flag.NewFlagSet("add", flag.ExitOnError)
		a := addCmd.Int("a", 0, "first number")
		b := addCmd.Int("b", 0, "second number")
		addCmd.Parse(os.Args[2:])
		fmt.Printf("%d + %d = %d\n", *a, *b, *a+*b)

	default:
		fmt.Printf("Unknown command: %q\n", os.Args[1])
		os.Exit(1)
	}
}
```

## Step-by-step execution

Run `go run main.go greet -name=Alice`:

1. `len(os.Args)` is 3, so the no-command check passes.
2. `os.Args[1]` = `"greet"`, matches the `case "greet"` branch.
3. `flag.NewFlagSet("greet", flag.ExitOnError)` creates a fresh FlagSet.
4. `greetCmd.String("name", "World", ...)` registers the `-name` flag.
5. `greetCmd.Parse(os.Args[2:])` parses `["-name=Alice"]`.
6. `*name` is `"Alice"`, prints `"Hello, Alice!"`.

Run `go run main.go add -a=10 -b=3`:

1. `os.Args[1]` = `"add"`, matches the `case "add"` branch.
2. `flag.NewFlagSet("add", ...)` creates a new FlagSet.
3. Two flags `-a` and `-b` are registered.
4. `Parse` picks up `-a=10` and `-b=3`.
5. Prints `"10 + 3 = 13"`.

## Common mistakes

- **Parsing `os.Args[1:]` instead of `os.Args[2:]`**: The subcommand name itself is not a flag; skip it.
- **Using the global `flag.Parse()`**: Global flags would consume the subcommand as a flag value.
- **Missing subcommand check**: Panics on `os.Args[1]` when no args are provided.
- **Reusing the same FlagSet**: Each subcommand must have its own FlagSet; reuse mixes flags.

## Debugging walkthrough

Buggy code:

```go
flag.Parse()  // Global parse — consumes "greet" as an unknown flag
switch os.Args[1] {
case "greet":
	// ...
}
```

**Symptom**: Running `go run main.go greet -name=Alice` prints "flag provided but not defined: -name" or ignores the subcommand.

**Investigation**: Comment out `flag.Parse()` and print `os.Args` before the switch.

**Fix**: Remove global `flag.Parse()`. Each subcommand parses its own segment with `fs.Parse(os.Args[2:])`.

## Production notes

- Use `flag.ContinueOnError` if you want to validate flags per subcommand without aborting the program (e.g., in tests).
- For large CLIs (10+ subcommands), consider using a library like `spf13/cobra` which provides subcommand routing, help generation, and completion scripts.
- Keep the router function small — delegate each subcommand to its own function: `runGreet(greetCmd)`.
- Use a consistent exit code convention: 0 for success, 1 for user error, 2 for flag parsing errors.

## Performance implications

- Each `flag.FlagSet` allocates a map of flag definitions. Creating a FlagSet per invocation is cheap (microseconds).
- The switch statement is O(1) per subcommand. For large numbers of subcommands, prefer a `map[string]func([]string)` lookup.

## Practice task

Extend the example with a third subcommand `"calc"` that accepts `-op` (string, one of "add", "sub", "mul", "div"), `-x` (float64), and `-y` (float64), and prints the result. If division by zero, print an error. Save as `main.go`.

## Tests / verification

```bash
go run ./curriculum/modules/07-cli-files-json-config/lessons/03-subcommands
go test ./curriculum/modules/07-cli-files-json-config/lessons/03-subcommands
```

## Review questions

1. Why must each subcommand have its own `flag.FlagSet` instead of using the global `flag` package?
2. What is the value of `os.Args[2]` in the command `mytool greet -name=Alice`?
3. What happens if `flag.ExitOnError` is used and the user provides an invalid flag value?
4. How would you implement a default subcommand (e.g., `mytool` alone runs `mytool help`)?
5. What are the three error-handling behaviors provided by `flag.NewFlagSet`?

## NEXT UP

Standard input and standard output — reading and writing the fundamental I/O streams.
