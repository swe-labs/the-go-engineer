# Lesson 02: What zero magic means

## Learning objective

By the end of this lesson, you will be able to identify "magic" in code — behavior that happens implicitly without being visible to the reader — and replace it with "zero magic" — explicit, readable, verifiable code. You will learn why this curriculum bans magic and how to apply the principle in Go.

## Why this matters

Magic is the enemy of learning. When a function silently defaults to port 8080, a beginner sees "it just works" but has no idea why. When a configuration value is zero but the program behaves as though it were non-zero, the learner cannot trace cause to effect. Zero magic means every behavior is visible in the code you are looking at. This makes learning faster, debugging easier, and code reviews more productive. In production, zero-magic code is code that another engineer can understand at 2 AM during an incident.

## Mental model

Imagine a vending machine. A magic vending machine dispenses a drink when you press a button, but you cannot see the price, the selection logic, or whether the machine accepted your money. A zero-magic vending machine displays the price, shows your balance, tells you which items are in stock, and prints a receipt. Both machines work, but only one lets you understand what happened.

In code, magic is the vending machine with a blank screen. Zero magic is the machine that shows everything.

## Core idea

"Magic" in software engineering refers to behavior that is real but not visible in the immediate code. Common forms of magic include:

- **Silent defaults**: A function uses port 8080 when no port is specified, but that default is documented nowhere in the code you are reading.
- **Implicit state**: A function behaves differently depending on some global variable or environment variable that is not passed as an argument.
- **Hidden side effects**: A function writes to a database, sends an email, or logs to a file without those effects being clear from its name or signature.
- **Unexplained conventions**: "We always multiply by 1000 because the API expects milliseconds" — but the code just says `* 1000` with no comment or constant name.

Zero magic eliminates each of these by making all behavior explicit.

## Under the hood

In Go, zero magic often means using explicit configuration structs, named constants, and clear function signatures that accept all dependencies as parameters.

The `Config` struct in `main.go` holds the port, timeout, and verbosity settings. `MagicServer()` creates an empty `Config{}` and uses it directly — Go initializes all fields to their zero values (0 for int, false for bool). The server "works" but the port is 0, which would fail in real code. The reader sees nothing about defaults.

`ZeroMagicServer(cfg Config)` accepts the config as a parameter and explicitly checks each field. If `cfg.Port` is 0 (the zero value for int), it sets it to 8080. If `cfg.Timeout` is 0, it sets it to 30. Every possible state is handled in plain sight.

## How Go uses it

Go itself favors zero magic. The `net/http` package requires you to explicitly set `http.Server.Addr` — there is no hidden default for the listen address. The `flag` package requires you to register flags explicitly and call `flag.Parse()`. The `os.Getenv` function returns an empty string for missing environment variables, and it is up to you to check the second return value (whether the variable was set). Go's design philosophy is that explicit is better than implicit, which is exactly the zero-magic principle.

## Go example

```go
package main

import (
	"fmt"
)

type Config struct {
	Port    int
	Timeout int
	Verbose bool
}

func MagicServer() {
	cfg := Config{}
	fmt.Println("MagicServer started on port", cfg.Port)
	fmt.Println("Timeout:", cfg.Timeout, "seconds")
	fmt.Println("Verbose:", cfg.Verbose)
}

func ZeroMagicServer(cfg Config) {
	if cfg.Port == 0 {
		cfg.Port = 8080
	}
	if cfg.Timeout == 0 {
		cfg.Timeout = 30
	}
	fmt.Println("ZeroMagicServer started on port", cfg.Port)
	fmt.Println("Timeout:", cfg.Timeout, "seconds")
	fmt.Println("Verbose:", cfg.Verbose)
}

func main() {
	fmt.Println("=== Magic approach (hidden defaults) ===")
	MagicServer()

	fmt.Println()
	fmt.Println("=== Zero-magic approach (explicit defaults) ===")
	ZeroMagicServer(Config{Port: 9090, Timeout: 60, Verbose: true})
	ZeroMagicServer(Config{})
}
```

## Step-by-step execution

1. The program defines `Config` with three fields, all of which get zero values when not initialized.
2. `MagicServer()` creates an empty `Config`. Because `Port` is 0 and `Timeout` is 0, the program prints "MagicServer started on port 0" and "Timeout: 0 seconds". This is technically correct Go — the zero value is valid — but it is misleading. A real server on port 0 would fail to bind.
3. `ZeroMagicServer(cfg Config)` checks each field before using it. If `Port` is 0, it replaces it with 8080. If `Timeout` is 0, it replaces it with 30. The caller sees exactly what will happen.
4. `main()` calls both to demonstrate the difference. The output shows that `MagicServer` uses port 0 silently, while `ZeroMagicServer` always has valid values.
5. The tests verify that `ZeroMagicServer` applies defaults correctly for zero-value, explicit-value, and partial-config scenarios.

## Common mistakes

| Mistake | Why it happens | Fix |
|---------|---------------|-----|
| Confusing zero value with "not set" | In Go, every type has a zero value (0, "", false, nil). Beginners often forget that an uninitialized int is 0, not "absent". | Check for zero explicitly, or use a pointer (`*int`) where nil means unset. |
| Assuming a function has no side effects because it doesn't mention them | A function named `Process()` might write to a database, but you cannot tell from the name. | Name functions by what they do: `ProcessAndSave`. Accept dependencies as parameters, not globals. |
| Documenting defaults in comments instead of code | A comment says "default port is 8080" but the code uses 0. The comment and code drift apart. | Put the default in code: `if port == 0 { port = 8080 }`. The code is the truth. |
| Forgetting that bool zero value is false | `if !cfg.Verbose` is true when Verbose was never set. That might be the right behavior accidentally. | Be explicit: `if !cfg.Verbose && wasSet { ... }` or use a pointer bool. |

## Debugging walkthrough

Scenario: You call `ZeroMagicServer(Config{Port: 0})` and the server starts on port 8080. You expected it to fail.

Step 1: Read the function body. Look for the line `if cfg.Port == 0 { cfg.Port = 8080 }`. There it is — the default override.

Step 2: Ask yourself: was this the intended behavior? For this lesson, yes — it applies a default. But in your own code, you might want to distinguish "user explicitly set port 0" from "user didn't set a port."

Step 3: To distinguish these cases, change `Port` to `*int`. A nil pointer means "not set"; a pointer to 0 means "explicitly zero." This is more verbose but truly zero magic.

Step 4: If the test fails with "got port 0, want 8080", the code that applies defaults was not reached. Check for a typo: `if cfg.Port == 0` (double equals) vs `if cfg.Port = 0` (single equals, assignment).

## Production notes

- Real production code should validate config at startup, not at every function call. Move the default-application logic into a `func (c Config) WithDefaults() Config` method called once when the program starts.
- Avoid global variables for configuration. Pass config explicitly through function parameters or dependency injection. This makes testing trivial and eliminates magic.
- Consider using the "functional options" pattern for complex configuration. This keeps the zero-magic principle while making call sites cleaner.

## Performance implications

- Checking integer fields for zero and assigning defaults adds a few CPU cycles per check. For server startup (done once), this is invisible. For per-request code where config is read once and reused, the cost is also negligible.
- The zero-magic approach sometimes uses more memory (e.g., `*int` instead of `int` to distinguish "not set" from "zero"). Profile before optimizing — clarity is almost always worth the small overhead.
- Function calls that check defaults inline (branch prediction) cost almost nothing on modern CPUs. The readability gain far outweighs the microsecond.

## Practice task

Add a new field `MaxRetries int` to the `Config` struct. Update `ZeroMagicServer` to default `MaxRetries` to 3 when it is zero. Add a test case to the table-driven test that verifies `MaxRetries` defaults to 3 for an empty config and respects an explicit value of 5.

## Tests / verification

Run the tests from the repository root:

```bash
go test ./curriculum/modules/00-orientation/lessons/02-what-zero-magic-means/
```

Expected output:

```
ok      github.com/rasel9t6/the-go-engineer/curriculum/modules/00-orientation/lessons/02-what-zero-magic-means
```

## Review questions

1. What is the zero value of an `int` in Go? What is the zero value of a `bool`?
2. Why does `MagicServer()` print "port 0" while appearing to work?
3. How does `ZeroMagicServer` differ from `MagicServer` in terms of the `Config` parameter?
4. What change would you make to `Config` to distinguish "user did not set a port" from "user set port to 0"?
5. Give one example of "magic" from a Go standard library function and explain how you could make it zero-magic.

## NEXT UP

Lab 03: How lessons, exercises, projects, and checkpoints work — where you learn the learning cycle that every module follows.
