# Flags

## Learning objective

Define, parse, and access named command-line flags using Go's `flag` package, and distinguish between flags and positional (non-flag) arguments.

## Why this matters

Raw `os.Args` forces you to parse everything by hand. The `flag` package gives you a standards-compliant, type-safe, self-documenting way to handle options like `--output=file.json`, `-v`, and `-count=5`. Every professional Go tool uses flags — from the Go compiler itself to Docker, Kubernetes, and Hugo.

## Mental model

Think of flags as key-value pairs prefixed with `-` or `--`. You define them once (name, default, help text), call `flag.Parse()`, and then read the parsed values via pointers. Non-flag arguments (bare words after all flags) are collected in `flag.Args()`.

```
$ mytool -name=Alice -count=3 extra.txt
       |_____________|          |_________|
         flags                  flag.Args()
```

## Core idea

The `flag` package works in three phases:

1. **Define** — call `flag.String()`, `flag.Int()`, `flag.Bool()` etc. to declare flags and store their default values.
2. **Parse** — call `flag.Parse()` to scan `os.Args[1:]` and populate the flag variables.
3. **Read** — dereference the pointers returned from the `flag.Type()` functions to get the parsed values.

Flag syntax: `-flag`, `-flag=value`, `--flag`, `--flag=value`. Single-dash and double-dash are equivalent. For booleans, `-flag` sets it to `true`.

## Under the hood

`flag.Parse()` registers a `FlagSet` (defaulting to `CommandLine`) and iterates over `os.Args[1:]`. It matches each token against registered flags by name. When it finds a match, it calls the flag's `Value.Set(string)` method, which converts the string to the target type. The first unrecognized `-` token or any non-`-` token stops flag parsing — everything after is a non-flag argument.

Flag definitions are stored in a `map[string]*Flag` inside the `FlagSet`. Each `Flag` holds the name, usage string, default value, and a `Value` interface for set/get.

## How Go uses it

- `flag.String(name, default, usage)` returns a `*string` — you dereference with `*name`.
- `flag.Int`, `flag.Float64`, `flag.Bool`, `flag.Duration`, etc. follow the same pattern.
- `flag.Parse()` must be called before accessing flag values, typically at the top of `main()`.
- After parsing, `flag.Args()` returns the remaining positional args as `[]string`.
- `flag.NArg()` returns the count of non-flag args; `flag.NFlag()` returns how many flags were set.
- `flag.PrintDefaults()` prints help text for all defined flags.

## Go example

```go
package main

import (
	"flag"
	"fmt"
)

func main() {
	name := flag.String("name", "World", "a name to greet")
	count := flag.Int("count", 1, "number of times to greet")
	verbose := flag.Bool("verbose", false, "enable verbose output")

	flag.Parse()

	if *verbose {
		fmt.Printf("(verbose) name=%s count=%d\n", *name, *count)
	}

	for i := 0; i < *count; i++ {
		fmt.Printf("Hello, %s!\n", *name)
	}

	fmt.Println("Non-flag args:", flag.Args())
}
```

## Step-by-step execution

Run `go run main.go -name=Alice -count=3 extra.txt`:

1. `name` pointer is allocated with default `"World"`.
2. `count` pointer allocated with default `1`.
3. `verbose` pointer allocated with default `false`.
4. `flag.Parse()` scans args: matches `-name=Alice` → sets `*name = "Alice"`, matches `-count=3` → sets `*count = 3`.
5. `*verbose` is `false`, so verbose block is skipped.
6. Loop runs 3 times printing `"Hello, Alice!"`.
7. `flag.Args()` returns `["extra.txt"]`, printed at the end.

## Common mistakes

- **Calling `flag.Parse()` too late**: Using flag pointers before `Parse()` reads the defaults, not the user's values.
- **Forgetting to dereference**: `flag.String(...)` returns a pointer; using `name` instead of `*name` prints a memory address.
- **Mixing flags and positional args**: All flags must come before positional args for the default `FlagSet`; `-` stops parsing.
- **Using `flag.Int` for zero values**: You cannot distinguish "user passed 0" from "default 0". Use `flag.Func` or a custom type with a sentinel.
- **Panic on required flags**: The `flag` package has no built-in "required" flag. Validate manually after `Parse()`.

## Debugging walkthrough

Buggy code:

```go
name := flag.String("name", "World", "a name to greet")
fmt.Println(*name)  // Works fine: prints "World"
```

But this panics:

```go
name := flag.String("name", "World", "a name to greet")
name = flag.Arg(0)   // compile error: cannot use string as *string
```

**Symptom**: Compile error, then confusion about pointers vs values.

**Fix**: Keep the pointer from `flag.String` and double-check whether a function returns a value or a pointer.

## Production notes

- Call `flag.Parse()` early in `main`, before any logic except maybe log setup.
- Use `flag.CommandLine` for simple tools; for subcommands use custom `flag.FlagSet` (next lesson).
- Always validate required flags manually: `if *name == "" { flag.Usage(); os.Exit(1) }`.
- Environment variable overrides are common: check env, then use that as the default in `flag.String`.
- For help text, `flag.Usage` is a function variable — override it for custom formatting.

## Performance implications

- Flag parsing is O(n) over the argument list and negligible compared to program startup.
- Each `flag.Parse` call registers defaults in a map; there is no allocation per flag after definition.
- The `flag` package is not designed for repeated parse/reparse cycles — parse once at startup.

## Practice task

Write a program that accepts three flags: `-input` (string), `-lines` (int, default 10), and `-reverse` (bool). Print `"Input: <value>, Lines: <value>, Reverse: <value>"`. Then add validation: if `-input` is empty, print an error and exit with status 1. Save as `main.go`.

## Tests / verification

```bash
go run ./curriculum/modules/07-cli-files-json-config/lessons/02-flags -- -name=Go
go test ./curriculum/modules/07-cli-files-json-config/lessons/02-flags
```

## Review questions

1. What does `flag.String("port", "8080", "listen port")` return — a string or a pointer to a string?
2. What must you call before reading flag values, and what happens if you don't?
3. How does the `flag` package distinguish between a flag `-name` and a positional argument `-name`?
4. If you run `go run main.go -verbose file.txt -count=3`, what does `flag.Args()` contain?
5. How would you implement a required flag in the `flag` package?

## NEXT UP

Subcommands — building multi-command CLIs like `git commit` and `docker run`.
