# Documentation comments

## Learning objective

Write godoc-formatted documentation comments for packages, types, functions, and methods; link identifiers across documentation; mark deprecated items; use the `go doc` command; and write example functions that double as tests.

## Why this matters

Documentation is how other developers — including your future self — understand what your code does and how to use it. Go has a built-in documentation culture: `godoc` (now `go doc`) reads comments and generates documentation automatically. Well-documented Go code is a hallmark of professional engineering. Poorly documented code forces readers to reverse-engineer intent from implementation, wasting time and causing bugs.

## Mental model

A documentation comment is a conversation with the reader. Imagine you are explaining your code to a colleague who knows Go but has never seen your package.

- **Package comment**: "Here is what this package is for and when you would use it."
- **Type comment**: "Here is what this type represents."
- **Function/ Method comment**: "Here is what this function does, what it expects, what it returns."

The `go doc` tool is the reader's interface: it extracts these conversations and presents them as a beautifully formatted manual.

## Core idea

### Doc comment format

Doc comments appear immediately before a top-level declaration, with no blank line between the comment and the declaration:

```go
// Package mathutil provides basic arithmetic operations.
package mathutil

// Add returns the sum of a and b.
func Add(a, b int) int {
    return a + b
}

// A Point represents a 2D coordinate.
type Point struct {
    X, Y int
}
```

Rules:
- Start with the name of the declared item (e.g., `Add returns...`, `Point represents...`).
- Use complete sentences. Period at the end.
- The first sentence becomes the short summary shown in `go doc` listings.
- Use blank lines in the comment to separate paragraphs.
- Indent code examples with a tab.

### Linking identifiers

In doc comments, bracket syntax `[Identifier]` creates a hyperlink to that identifier's documentation:

```go
// [Add] computes the sum of two integers.
// See also [Subtract] and [Multiply].
func Add(a, b int) int
```

This works for:
- Symbols in the same package: `[Add]`
- Symbols in other packages: `[math.Sqrt]`, `[io.Reader]`
- Type members: `[Point.X]`

### Deprecation comments

To mark an identifier as deprecated, add `// Deprecated:` followed by a message:

```go
// Deprecated: Use Add instead.
func OldAdd(a, b int) int {
    return a + b
}
```

The `go doc` tool and linters (like `staticcheck`) highlight deprecated symbols.

### `go doc` command

```bash
go doc <package>          # show package documentation
go doc <package>.Func     # show function documentation
go doc <package>.Type     # show type documentation
go doc <package>.Type.Method  # show method documentation
```

Examples:
```bash
go doc fmt
go doc fmt.Println
go doc http.Server
go doc http.Server.ListenAndServe
```

### Example functions

Example functions are runnable code that appears in `go doc` output. They are placed in `_test.go` files:

```go
func ExampleAdd() {
    sum := Add(2, 3)
    fmt.Println(sum)
    // Output: 5
}
```

Key rules:
- Function name starts with `Example`.
- If it is for a type: `ExampleType_Method`.
- The `// Output:` comment at the end is the expected output. The test runner verifies it matches.
- Examples appear in `go doc` output in order of declaration.

## Under the hood

The `go/doc` package parses Go source files, extracts doc comments, and builds a structured documentation tree. The `go/doc/comment` package (Go 1.19+) parses the new doc comment format with support for headings, lists, links, and code blocks.

The `// Deprecated:` marker is detected by `go/doc` and by analysis tools like `staticcheck` and `golangci-lint`. The Go compiler does not warn about deprecated usage — that is left to linters and IDE integration.

Example functions are compiled and run as part of `go test`. The `// Output:` comment uses the testing framework's golden output mechanism: the test captures stdout and compares it to the expected output.

## How Go uses it

- **Standard library**: Every package in the standard library follows these conventions. Run `go doc fmt` to see.
- **`go doc` command**: The primary way developers read documentation without leaving the terminal.
- **`pkgsite`**: The web-based Go documentation browser (at `pkg.go.dev`) renders godoc comments.
- **Linters**: `staticcheck` flags missing doc comments on exported names (configurable). `golint` used to require them.
- **IDEs**: VSCode and GoLand render doc comments as hover tooltips.
- **`gopls`**: The Go language server provides documentation on hover using the doc comment infrastructure.

## Go example

```go
// Package greeter provides friendly greeting functions.
//
// Use [Hello] to create a greeting string, or [Hi] to print
// a casual greeting directly.
//
// Example:
//
//	msg := greeter.Hello("World")
//	fmt.Println(msg)
package greeter

import "fmt"

// Hello returns a greeting for the given name.
//
// The greeting includes the name and an exclamation mark.
//
// Example:
//
//	greeting := greeter.Hello("Alice")
//	fmt.Println(greeting) // Hello, Alice!
func Hello(name string) string {
	return fmt.Sprintf("Hello, %s!", name)
}

// Hi prints a casual greeting to stdout.
//
// Deprecated: Use [Hello] and print the result yourself.
func Hi(name string) {
	fmt.Printf("Hi, %s!\n", name)
}
```

```go
// Example function — lives in a _test.go file.

func ExampleHello() {
	greeting := Hello("World")
	fmt.Println(greeting)
	// Output: Hello, World!
}
```

## Step-by-step execution

Using `go doc` on the above package:

1. `go doc ./greeter` — shows the package comment and a summary of exported names:
   ```
   Package greeter provides friendly greeting functions.
   ...
   func Hello(name string) string
   func Hi(name string)  // Deprecated
   ```

2. `go doc ./greeter Hello` — shows the full doc for `Hello`:
   ```
   func Hello(name string) string
       Hello returns a greeting for the given name.
       ...
   ```

3. Running tests with examples: `go test -run ExampleHello` — executes the function, captures output, compares to `// Output:`.

## Common mistakes

- **Blank line between comment and declaration**: This breaks the association — the comment becomes a general comment, not a doc comment. No blank line allowed.
- **Starting a function comment with "This function"**: Redundant. Start with the function name directly: "FunctionName does X."
- **Not documenting `error` returns**: Always document what errors a function can return and under what conditions.
- **Forgetting `// Output:` in example functions**: Without it, the example is compiled and run but no output is verified. Include `// Output:` to make it a tested example.
- **Over-documenting internal details**: Doc comments explain what and why, not how. Implementation details belong in implementation, not documentation.
- **Deprecating without suggesting a replacement**: Always say what to use instead: `// Deprecated: Use NewFunc instead.`
- **Not using link syntax `[Identifier]`**: Links make documentation navigable. Use them for cross-references.

## Debugging walkthrough

Symptom: `go doc` shows no documentation for an exported function.

**Root cause**: The doc comment has a blank line between the comment and `func`:

```go
// MyFunc does something.

func MyFunc() {}  // Wrong! Blank line breaks the doc comment.
```

**Fix**: Remove the blank line:

```go
// MyFunc does something.
func MyFunc() {}
```

Symptom: An example function compiles but does not appear in `go doc` output.

**Root cause**: The example function is in a regular `.go` file, not a `_test.go` file. Example functions must be in test files.

**Fix**: Move the example to `*_test.go`.

Symptom: `go test` reports a failure in an example function:

```
got:
Hello, World!
want:
Hello world!
```

**Root cause**: The `// Output:` comment does not match the actual output (including exact punctuation and spacing).

**Fix**: Update the `// Output:` comment to match the actual output exactly.

## Production notes

- **Every exported name gets a doc comment**: This is the Go community standard. Linters can enforce it. In code review, undocumented exported symbols should be flagged.
- **Package comments are important**: Every package should have a package-level doc comment explaining its purpose. For packages with many files, the comment typically goes in `doc.go`.
- **`doc.go` convention**: A file named `doc.go` containing only the package comment and `package` declaration is a common pattern for packages with complex documentation.
- **Deprecation policy**: Mark deprecated items with `// Deprecated:` and provide a replacement. After two major versions, consider removing the deprecated item entirely.
- **CI integration**: Run `go vet ./...` to catch common documentation issues. Use `staticcheck` to enforce doc comment presence.
- **`//go:embed` and documentation**: When embedding files, document the purpose and format of embedded data.

## Performance implications

- Doc comments are parsed at compile time into AST nodes. The memory overhead is negligible (a few KB per package).
- Example functions are compiled as test binaries, which adds to `go test` compilation time. Keep examples concise.
- The `go doc` command is fast — it reads source files, not compiled binaries.
- There is no runtime cost for doc comments. They are stripped from compiled binaries.
- Linking identifiers in doc comments is resolved at documentation-generation time, not at compile time.

## Practice task

Create a package `calc` in subdirectory `calc/` with:

1. A package comment documenting the package.
2. An exported function `Sum(nums ...int) int` with a doc comment.
3. An exported type `Calculator` with a doc comment and an exported method `Result() int`.
4. A deprecated function `OldSum` that delegates to `Sum`.
5. An example function `ExampleSum` in `calc_test.go` with `// Output:`.
6. Link identifiers: In the `Calculator` doc comment, link to `Sum`.

Then use `go doc ./calc` and `go doc ./calc Sum` to verify the documentation renders correctly.

## Tests / verification

```bash
go run ./curriculum/modules/05-types-interfaces-packages-modules/lessons/23-documentation-comments
go test ./curriculum/modules/05-types-interfaces-packages-modules/lessons/23-documentation-comments
```

## Review questions

1. What is the rule for the first sentence of a Go doc comment?
2. How do you create a hyperlink to another documented identifier in a godoc comment?
3. What is the syntax for marking a function as deprecated?
4. Why must example functions live in `_test.go` files?
5. What happens if there is a blank line between a doc comment and the declaration?

## NEXT UP

Congratulations on completing Module 05! You now understand Go's type system with structs, methods, interfaces, and the module system. Next up: Module 06 — Testing, Debugging, and Refactoring.
