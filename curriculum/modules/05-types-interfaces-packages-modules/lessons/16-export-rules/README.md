# Export rules

## Learning objective

Distinguish exported from unexported Go names, apply package-level visibility rules, understand `internal` package restrictions, and use factory functions to control construction of unexported types.

## Why this matters

Every Go package is a capsule. Export rules are the capsule wall — they determine what other packages can see and use. Design your capsule wall well: expose a clean, minimal API and hide implementation details. Poor visibility choices lead to fragile code where changing internal logic breaks external consumers, or to bloated APIs where every helper is public. Professional Go code relies on export rules for encapsulation, the same way classes rely on private/protected in OOP languages.

## Mental model

Imagine a package is a house with two kinds of items:
- **Exported** (capital letter): in the front yard — anyone walking by can see and use them.
- **Unexported** (lowercase letter): inside the house — only people inside (the same package) can see and use them.

The `internal/` directory is a gated community: only packages that share the same root can enter. This is a stronger restriction than unexported — it restricts across packages, not just within one.

## Core idea

**Exported** = name starts with a capital letter. Visible to any package that imports this package.

**Unexported** = name starts with a lowercase letter. Visible only within the same package.

```go
package widget

var Count int        // exported — other packages can read/write Count
var count int        // unexported — only widget package can access count

func New() *Widget   // exported — public constructor
func newHelper()     // unexported — internal helper
```

This rule applies to:
- Types (`type Widget struct{...}` vs `type widget struct{...}`)
- Functions (`func New() ...` vs `func newHelper() ...`)
- Variables and constants (`var Version = "1.0"` vs `var version = "0.5"`)
- Struct fields (`Name string` exported, `name string` unexported)
- Methods (`func (w *Widget) Render()` vs `func (w *Widget) render()`)

### `internal` package restriction

A package at `internal/` or any subdirectory of it can only be imported by packages rooted at the parent of `internal/`. For example:

```
project/
  go.mod           → module github.com/user/project
  internal/
    db/
      db.go        → only importable by github.com/user/project and its subpackages
  cmd/
    server/
      main.go      → can import internal/db ✓
  pkg/
    client/
      client.go    → can import internal/db ✓
```

Another module `github.com/other/project` cannot import `github.com/user/project/internal/db` — the compiler rejects it.

### Factory functions

When a type has unexported fields, external packages cannot construct it directly. A factory function (exported constructor) is the idiomatic solution:

```go
package user

type User struct {
    Name string   // exported
    age  int      // unexported
}

func New(name string, age int) *User {
    return &User{Name: name, age: age}
}
```

## Under the hood

Export checking happens at compile time in the type-checker phase. When package A imports package B and references `B.someName`, the compiler checks whether `someName` starts with an uppercase letter. If not, it reports:

```
cannot refer to unexported name B.someName
```

The `internal` restriction is enforced by the `go/build` package (and by the module-aware loader). It checks that the importing package's import path has the same root prefix as the `internal` package's containing module up to and including the `internal` directory.

## How Go uses it

- **Standard library encapsulation**: `os.File` has unexported fields; you create it via `os.Open` or `os.Create`, never by literal.
- **`http.ListenAndServe`**: takes an `http.Handler` interface; the concrete `http.ServeMux` is exported but its internals are hidden.
- **`context` package**: `context.Background()` and `context.WithCancel()` return interface values; the concrete types `emptyCtx` and `cancelCtx` are unexported.
- **`database/sql`**: `sql.DB` is exported but its connection pool internals are unexported. You get a `*sql.DB` via `sql.Open`.
- **Protobuf generated code**: Types are exported, but internal marshaling helpers are unexported.

## Go example

```go
package main

import (
	"fmt"

	"github.com/rasel9t6/the-go-engineer/curriculum/modules/05-types-interfaces-packages-modules/lessons/16-export-rules/person"
)

func main() {
	// Using the exported factory function.
	p := person.New("Alice", 30)
	fmt.Println("Name:", p.Name)
	// fmt.Println("Age:", p.age)  // compile error: p.age unexported

	// person.personData is unexported — cannot reference it here.
	// var pd person.personData  // compile error
}
```

```go
// person/person.go
package person

type personData struct {
	name string
	age  int
}

func New(name string, age int) *personData {
	return &personData{name: name, age: age}
}

func (p *personData) Name() string {
	return p.name
}

func (p *personData) age() int { // unexported method
	return p.age
}
```

## Step-by-step execution

For `person.New("Alice", 30)`:

1. The `person` package exports `New` — a function that returns `*personData`.
2. Inside `New`, a `personData` is created with `name: "Alice"` and `age: 30`.
3. The `personData` type is **unexported** — the caller gets a `*personData` but cannot write `person.personData{}`.
4. The caller can call `p.Name()` (exported) but not `p.age()` (unexported).
5. The `Name` field is accessed through the exported accessor method (or exported field if you design it that way).

For `internal` restriction:

1. Package `github.com/user/project/internal/db` declares exported function `Connect`.
2. Package `github.com/user/project/cmd/server` imports `"github.com/user/project/internal/db"` — OK, same root.
3. Package `github.com/other/module` tries the same import — compiler error: `use of internal package ... is not allowed`.

## Common mistakes

- **Assuming unexported fields are invisible to exported methods**: Unexported fields are accessible within the same package, including exported methods of the same type. A method `func (u *User) Age() int { return u.age }` works fine.
- **Thinking unexported means inaccessible to tests in the same package**: Tests in `package person` can access `personData` and its unexported fields. Only `_test` packages (with `package person_test`) cannot.
- **Exposing internal types in exported function signatures**: An exported function returning an unexported type forces callers to use type inference. It works but is confusing. Export the type or return an interface.
- **Forgetting that constants are also subject to export rules**: `const maxRetries = 3` is unexported; `const MaxRetries = 3` is exported.
- **`internal` is not recursive upward**: A package at `a/b/internal/c` is restricted to `a/b` and its children. But `a/b/internal/c/internal/d` is restricted to `a/b/internal/c` and its children — the `internal` at each level creates a new boundary.

## Debugging walkthrough

Symptom: `go build` fails with:

```
./main.go:10:14: cannot refer to unexported name person.personData
```

**Code**:
```go
package main
import ".../person"
func main() {
    var p person.personData  // error!
}
```

**Root cause**: `personData` starts with lowercase `p` — it is unexported. The `main` package cannot reference it by name.

**Fix**: Use the exported factory function:
```go
p := person.New("Alice", 30)
```

Or, if the type must be accessible, rename it to `PersonData` (capital P).

## Production notes

- **Minimal exported API**: A common Go proverb is "the bigger the interface, the weaker the abstraction." Apply this to packages: export as little as possible. Unexported functions and types can be changed freely without breaking consumers.
- **`internal` for monoliths**: In large monorepos, use `internal/` to enforce architectural boundaries between teams. Each service's internal packages are safe from accidental import by other services.
- **Breaking changes**: Adding a new exported name is safe. Removing or renaming an exported name is a breaking change. Changing an unexported name is safe. Make names unexported by default; export only when needed.
- **`golint` / `staticcheck`**: Linters warn about exported names without doc comments. Every exported name should have a comment explaining its purpose.

## Performance implications

- Export rules are purely a compile-time mechanism — zero runtime cost.
- The `internal` restriction is checked at compile time; no runtime checks.
- However, accessing unexported fields through exported accessor methods adds a function call overhead compared to direct field access. In hot paths, consider exporting the field or keeping direct access within the package.
- Export decisions affect binary size only indirectly: more exported API means more documentation and potentially more tests, but the compiled code is the same.

## Practice task

Create a package `bank` (in subdirectory `bank/`) with:

- An unexported type `account` with unexported fields `balance float64` and `owner string`.
- An exported factory function `NewAccount(owner string, initialBalance float64) *account`.
- Exported methods `Balance() float64`, `Deposit(amount float64)`, and `Withdraw(amount float64) error`.
- `Withdraw` returns an error if amount exceeds balance.

Then write `main()` that creates an account, deposits, withdraws, and prints the balance. Verify that creating an `account` literal directly from `main` is a compile-time error.

## Tests / verification

```bash
go run ./curriculum/modules/05-types-interfaces-packages-modules/lessons/16-export-rules
go test ./curriculum/modules/05-types-interfaces-packages-modules/lessons/16-export-rules
```

## Review questions

1. What determines whether a Go name is exported from a package?
2. Can an exported method access unexported fields of the same struct?
3. What does the `internal/` directory do that ordinary unexported names cannot?
4. Why would you use a factory function instead of exporting a struct type directly?
5. If you change an unexported function's signature, do you need to bump the major version of your module? Why or why not?

## NEXT UP

Project layout — how to organise Go projects using `cmd/`, `internal/`, `pkg/`, and flat layout conventions.
