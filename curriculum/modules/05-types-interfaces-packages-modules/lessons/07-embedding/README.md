# Embedding

## Learning objective

Distinguish embedded fields from named fields, describe method promotion and interface embedding, resolve naming conflicts, and explain why embedding is not subclassing.

## Why this matters

Embedding is one of Go's most distinctive features. It enables code reuse without inheritance, keeps composition explicit, and powers the standard library's interface design (e.g., `io.ReadWriter`, `sort.Interface`). Understanding embedding is essential for reading idiomatic Go and designing clean abstractions.

## Mental model

Embedding is automatic delegation. When you embed `A` in `B`, every method and field of `A` is automatically available on `B`. It is as if you wrote forwarding methods by hand — but the compiler does it for you. B does not "inherit" from A; B "has" an A and conveniently exposes A's interface.

## Core idea

**Embedded vs named field**:

```go
type Named struct {
    Logger // embedded — methods promoted
    Log    *Logger // named — no promotion
}
```

**Method promotion**: If `Logger` has a `Print` method, `Named` also has a `Print` method (forwarding to `Logger.Print`).

**Interface embedding**: An interface can embed other interfaces:

```go
type Reader interface {
    Read(p []byte) (n int, err error)
}

type Closer interface {
    Close() error
}

type ReadCloser interface {
    Reader
    Closer
}
```

A type satisfies `ReadCloser` by implementing both `Read` and `Close`.

**Conflict resolution**: If `B` embeds `A` and `C`, and both `A` and `C` have a method `F`, then accessing `b.F()` is ambiguous — compile error. Resolve by specifying `b.A.F()` or `b.C.F()`.

**Embedding is NOT subclassing**:
- No virtual methods. No dynamic dispatch on embedded methods.
- The embedded type cannot access the outer type's fields.
- You cannot override methods. If `B` embeds `A` and defines `F`, `B.F` shadows `A.F` — but `A.F` still exists at `B.A.F()`.

## Under the hood

For struct embedding, the compiler creates a field with the type name and generates forwarding methods for each promoted method. These are plain functions, not virtual — the receiver of the forwarded method is the embedded field, not the outer struct. Interface embedding is purely a compile-time merging of method requirements.

## How Go uses it

- **io package**: `io.ReadWriter = io.Reader + io.Writer`.
- **http package**: `http.ResponseWriter` embeds `io.Writer` and adds `Header()` and `WriteHeader()`.
- **error wrapping**: `fmt.Errorf("%w", err)` creates an error that embeds the original.
- **sort package**: `sort.Reverse` embeds `sort.Interface` and overrides `Less`.
- **Database models**: embed a `Model` struct with `ID`, `CreatedAt`, `UpdatedAt`, `DeletedAt`.

## Go example

```go
package main

import (
	"fmt"
	"io"
	"strings"
)

type Reader interface {
	Read(p []byte) (n int, err error)
}

type Writer interface {
	Write(p []byte) (n int, err error)
}

type ReadWriter interface {
	Reader
	Writer
}

type Logger struct {
	prefix string
}

func (l Logger) Log(msg string) {
	fmt.Printf("[%s] %s\n", l.prefix, msg)
}

type Server struct {
	Logger // embedded
	Host   string
}

func main() {
	s := Server{
		Logger: Logger{prefix: "api"},
		Host:   "localhost:8080",
	}
	s.Log("server started") // promoted method

	// Interface embedding demonstration
	var rw ReadWriter = &strings.Builder{}
	rw.Write([]byte("hello"))
	var buf [16]byte
	n, _ := rw.Read(buf[:])
	fmt.Println("read:", string(buf[:n]))
}
```

## Step-by-step execution

For `s.Log("server started")`:

1. Compiler looks for `Log` method on `Server` — not a direct method.
2. Checks embedded fields of `Server`: `Logger` has `Log`.
3. Compiler generates forwarding: `Server.Log(msg)` calls `s.Logger.Log(msg)`.
4. At runtime, `Logger.Log` is called with receiver `s.Logger`.
5. Inside `Log`, `l.prefix` is `"api"`, so output is `[api] server started`.

For interface embedding with `ReadWriter`:

1. `ReadWriter` requires `Read(p []byte) (n int, err error)` and `Write(p []byte) (n int, err error)`.
2. `strings.Builder` has both methods (value receiver on `Write`, pointer receiver on `Read`).
3. `var rw ReadWriter = &strings.Builder{}` compiles because `*strings.Builder` satisfies `ReadWriter`.
4. At runtime, `rw.Write(...)` calls `(*strings.Builder).Write`.

## Common mistakes

- Mistake: Thinking embedded methods can be overridden like virtual methods.
  - Fix: If `Server` defines `Log(msg string)`, it shadows `Logger.Log`. But `s.Logger.Log()` still works. No polymorphism.

- Mistake: Embedding the same interface twice (e.g., `io.ReadWriter` already embeds `io.Reader`; also embedding `io.Reader` directly causes no harm but is redundant).
  - Fix: Embed only the composed interface.

- Mistake: Assuming conflict resolution is automatic.
  - Fix: When two embedded types have the same method, you must disambiguate with the explicit type name.

- Mistake: Embedding a type expecting the outer type's methods to be accessible from the inner type.
  - Fix: The embedded type knows nothing about the outer type. If needed, pass a function or interface reference.

## Debugging walkthrough

This code fails:

```go
type A struct{}
func (A) F() {}

type B struct{}
func (B) F() {}

type C struct {
    A
    B
}

func main() {
    var c C
    c.F() // COMPILE ERROR: ambiguous
}
```

**Symptom**: `ambiguous selector c.F`.

**Root cause**: Both `A.F` and `B.F` are promoted to `C`. The compiler cannot determine which one you mean.

**Fix**: Specify the exact path:

```go
c.A.F() // OK
c.B.F() // OK
```

Or define `C.F` to delegate explicitly:

```go
func (c C) F() { c.A.F() }
```

## Production notes

- **Use embedding for delegation**, not for modeling taxonomy. If you find yourself asking "is this a kind of that?", use an interface instead.
- **Interface embedding** keeps interfaces composable and small. Prefer composing small interfaces over defining large ones.
- **Shadowing with intent**: When you shadow an embedded method, document that you are intentionally overriding behavior.
- **Deep embedding** (> 2 levels) hurts readability. Keep hierarchies flat.

## Performance implications

- **Zero overhead**: Promoted method calls are resolved at compile time and compiled to direct function calls.
- **Interface embedding** has no runtime cost — it is just a compile-time union of method sets.
- **Shadowed methods** add no overhead; the compiler selects the correct target at compile time.

## Practice task

Define an interface `Pinger` with `Ping() error` and an interface `HealthChecker` that embeds `Pinger` and adds `Status() map[string]string`. Implement `HealthChecker` with a struct `Monitor` that embeds a `Pinger` implementation. In `main()`, create a `Monitor` and call both `Ping()` and `Status()`. Verify with tests.

## Tests / verification

```bash
go run ./curriculum/modules/05-types-interfaces-packages-modules/lessons/07-embedding
go test ./curriculum/modules/05-types-interfaces-packages-modules/lessons/07-embedding
```

## Review questions

1. How does an embedded field differ from a named field in a struct?
2. What happens when two embedded types both have a method named `Close`?
3. Can an embedded type access the fields of the outer type?
4. Does embedding a type mean the outer type satisfies the embedded type's interfaces?
5. How does Go resolve a promoted field name that conflicts with a direct field of the outer type?

## NEXT UP

Interfaces — Go's implicit satisfaction model, interface values, and the empty interface.
