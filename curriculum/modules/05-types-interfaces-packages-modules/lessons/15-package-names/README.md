# Package names

## Learning objective

Choose correct package names following Go conventions, understand the relationship between directory paths and package identity, and write package-level documentation comments.

## Why this matters

Package names are the first thing another developer sees when importing your code. A well-named package communicates its purpose immediately. A poorly named package causes confusion, import collisions, and awkward client code like `widget.NewWidget()`. Go's package naming is part of the language's emphasis on clarity and simplicity — getting it right is a mark of professional Go code.

## Mental model

A package is a folder of `.go` files with the same `package` declaration at the top. The package name is the **label** on that folder. When someone imports your package, they use this label as a qualifier for all exported names:

```go
import "your/project/httputil"
// usage: httputil.ServeHTTP(...)
```

The package name should be **short**, **lowercase**, and **descriptive** — a noun that describes what the package provides. It is **not** required to be unique across your entire project, but it must be unique within the same directory.

## Core idea

Every `.go` file starts with a `package` clause:

```go
package httputil
```

Rules:
- The package name must be a valid Go identifier (lowercase, no hyphens).
- All files in the same directory **must** declare the same package name (except for test files, which can use `package xxx_test`).
- The directory name and package name **should** match by convention, but the compiler only enforces that files in the same directory share the same package name.
- Special package `main` tells the compiler to produce an executable binary, not a library archive.
- Package names should be **short** — typically one word, rarely two.
- Avoid generic names like `util`, `common`, `lib`, `base`.

### Import paths

The import path is the full filesystem path relative to your module root:

```go
import "github.com/rasel9t6/the-go-engineer/curriculum/modules/05-types-interfaces-packages-modules/lessons/15-package-names"
```

The last element of the import path **should** match the package name, but Go does not enforce this. If they diverge, users import by the path but qualify by the package name — a guaranteed source of confusion.

### Package documentation

A package-level doc comment (written just before the `package` line) becomes the `godoc` page for the package:

```go
// Package httputil provides HTTP utility functions for request parsing,
// response writing, and middleware composition.
package httputil
```

## Under the hood

The Go compiler resolves package names at compile time. The import graph must be acyclic. Each package is compiled into an archive `.a` file, which is then linked into the final binary. The package name becomes a namespace prefix in the compiled code — all exported symbols are mangled with the package path to prevent collisions.

The `go/build` package uses `GOPATH` or module mode to resolve import paths to directories. In module mode, the module path plus the subdirectory relative to the module root forms the full import path.

## How Go uses it

- **Standard library**: Short, memorable names: `fmt`, `net`, `http`, `os`, `io`, `strings`, `time`, `sync`.
- **`main` package**: The entry point of any Go program. `func main()` in `package main` produces an executable.
- **Test packages**: External test packages use `package xxx_test` to test a package as an external consumer, ensuring only the exported API is tested.
- **Package resolution**: `goimports` automatically resolves import paths. The convention of directory = package name lets `goimports` add the correct import statement.
- **Blank imports**: `import _ "image/png"` runs the package's `init()` functions without using its names.

## Go example

```go
// Package greeter provides friendly greeting functions.
package greeter

import "fmt"

// Hello returns a greeting for the given name.
func Hello(name string) string {
	return fmt.Sprintf("Hello, %s!", name)
}

// Hi prints a casual greeting to stdout.
func Hi(name string) {
	fmt.Printf("Hi, %s!\n", name)
}
```

```go
// File: main.go (separate directory, package main)
package main

import (
	"fmt"

	"github.com/rasel9t6/the-go-engineer/curriculum/modules/05-types-interfaces-packages-modules/lessons/15-package-names/greeter"
)

func main() {
	msg := greeter.Hello("Alice")
	fmt.Println(msg)
	greeter.Hi("Bob")
}
```

## Step-by-step execution

For the `greeter` package:

1. The compiler reads `greeter.go` and sees `package greeter`.
2. Package name `greeter` is registered. Exported names: `Hello`, `Hi`.
3. The `main` package imports `greeter` and accesses its exported functions via `greeter.Hello()`.

For package name resolution during import:

1. `import "github.com/rasel9t6/the-go-engineer/.../greeter"` tells the compiler to find the directory `greeter` under the module root.
2. The compiler reads the first `.go` file in that directory, extracts the `package` clause, and uses it as the qualifier.
3. If the package name does not match the directory name, the compiler still compiles it — but humans will be confused.

## Common mistakes

- **Generic package names**: `util`, `common`, `misc`, `lib`. These become dumping grounds for unrelated functions. Prefer specific names like `strutil`, `httputil`, `conv`.
- **Package name != directory name**: A package in `foo/bar/` named `baz` forces importers to write `baz.Something()`, but the directory path suggests `bar.Something()`. Always match.
- **Hyphens in package names**: Go identifiers cannot contain hyphens. Directory names with hyphens (e.g., `my-package`) force a different package name. Avoid hyphens in directory names intended for Go packages; use underscores or concatenation (`mypackage`, `my_package`).
- **UpperCase package names**: Package names must be lowercase. `package MyUtil` is a compile error.
- **Overly long package names**: `package requesthandlerparser` is hard to read and type. Shorten to `reqhandler` or split into subpackages.
- **`main` package with library code**: A `main` package cannot be imported by other packages. Keep library code in non-main packages.

## Debugging walkthrough

Symptom: `go build` fails with:

```
found packages greeter (greeter.go) and main (main.go) in .../15-package-names
```

**Root cause**: Two `.go` files in the same directory declare different package names (`greeter` and `main`).

**Fix**: Move the `greeter` package into its own subdirectory:

```
15-package-names/
  main.go          -- package main
  greeter/
    greeter.go     -- package greeter
```

## Production notes

- **No cyclic imports**: Go does not allow package import cycles. Design your package hierarchy as a DAG (directed acyclic graph).
- **Package size guideline**: A package should have a single, focused responsibility. If a package has more than 10-15 exported symbols, consider splitting it.
- **Internal packages**: Use `internal/` to restrict importability. Packages inside `internal/` can only be imported by code rooted at the parent of `internal/`. This is the compiler-enforced way to keep packages private within a project.
- **Renaming imports**: Use import aliases sparingly — `import f "fmt"` works but obscures the source. Only alias to resolve collisions.
- **`init()` functions**: Each package can have `init()` functions that run on import, in dependency order. They are useful for registration (e.g., `image/png` registering itself with `image`), but avoid side effects in library packages.

## Performance implications

- Package name resolution has zero runtime cost — it is a compile-time concept.
- The number of imported packages affects compile time, not runtime performance.
- A deep package tree increases the number of import statements but has no runtime overhead.
- The linker strips unused packages from the binary (dead code elimination), so importing a package only for its `init()` effects brings in all its dependencies.

## Practice task

Create two packages in this lesson directory:

1. A package `mathutil` (in subdirectory `mathutil/`) with exported functions `Add(a, b int) int` and `Mul(a, b int) int`.
2. A package `strutil` (in subdirectory `strutil/`) with an exported function `Upper(s string) string` that returns the uppercase version.

Then write `main()` that imports both packages and calls each function at least once. Verify by running `go run .` from the lesson directory.

## Tests / verification

```bash
go run ./curriculum/modules/05-types-interfaces-packages-modules/lessons/15-package-names
go test ./curriculum/modules/05-types-interfaces-packages-modules/lessons/15-package-names
```

## Review questions

1. What is the relationship between a directory name and the package name declared in `.go` files within it?
2. Why should package names avoid generic words like `util` or `common`?
3. What is special about `package main`?
4. Can two different directories have the same package name? If so, how do you import both?
5. Where does a package-level doc comment go, and how is it displayed?

## NEXT UP

Export rules — how Go controls visibility of names across package boundaries with capitalisation.
