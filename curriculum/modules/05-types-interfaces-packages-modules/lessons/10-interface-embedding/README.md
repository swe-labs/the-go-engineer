# Interface embedding

## Learning objective

Embed interfaces within interfaces to compose larger contracts, follow Go naming conventions for composed interfaces, and resolve method conflicts.

## Why this matters

Interface embedding is how Go builds complex abstractions from simple ones. `io.ReadWriter`, `io.ReadWriteCloser`, and `io.ReadSeeker` are composed from the same small building blocks (`Reader`, `Writer`, `Closer`, `Seeker`). This pattern appears throughout the standard library and in well-designed Go packages. Mastering it lets you create expressive, modular interfaces.

## Mental model

Embedding an interface is like making a to-do list by copying items from smaller lists. An interface that embeds `io.Reader` and `io.Writer` requires both `Read` and `Write`. There is no special runtime structure — the compiler merges the method requirements into one set. The composed interface is a union of its parts.

## Core idea

**Embedding interfaces**:

```go
type Reader interface {
    Read(p []byte) (n int, err error)
}

type Writer interface {
    Write(p []byte) (n int, err error)
}

type Closer interface {
    Close() error
}

type ReadWriter interface {
    Reader
    Writer
}

type ReadWriteCloser interface {
    Reader
    Writer
    Closer
}
```

**Naming conventions**:

| Composition | Name |
|---|---|
| Reader + Writer | `ReadWriter` |
| Reader + Closer | `ReadCloser` |
| Writer + Closer | `WriteCloser` |
| Reader + Writer + Closer | `ReadWriteCloser` |
| Reader + Seeker | `ReadSeeker` |

The pattern is: concatenate the method names (removing duplicates), sort by convention (Read < Write < Close < Seek).

**Conflict resolution**: If two embedded interfaces both declare a method with the same name and signature, it is fine — the method requirement appears once. If they declare the same name with different signatures, it is a compile error.

```go
type A interface { F() int }
type B interface { F() string }

type C interface { A; B } // COMPILE ERROR: conflicting F
```

## Under the hood

Interface embedding is purely a compile-time mechanism. The compiler merges the method sets of all embedded interfaces into the new interface. The resulting `itab` has one entry per unique method. There is no runtime overhead from embedding.

## How Go uses it

- **io.ReadWriter** — the canonical composed interface.
- **io.ReadWriteCloser** — used by network connections (`net.Conn`).
- **io.ReadSeeker** — used by `os.File`, `bytes.Reader`.
- **http.ResponseWriter** — "embeds" `io.Writer` conceptually (actually declares `Write` directly) plus `Header()` and `WriteHeader()`.
- **sort.Interface** — named, not composed, because the three methods are fundamental to sorting.
- **hash.Hash** — embeds `io.Writer` and adds `Sum`, `Reset`, `Size`, `BlockSize`.

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

type Validator interface {
	Validate() error
}

type Formatter interface {
	Format() string
}

type ValidatedFormatter interface {
	Validator
	Formatter
}

type User struct {
	Name  string
	Email string
}

func (u User) Validate() error {
	if u.Name == "" {
		return fmt.Errorf("name is required")
	}
	if u.Email == "" {
		return fmt.Errorf("email is required")
	}
	return nil
}

func (u User) Format() string {
	return fmt.Sprintf("%s <%s>", u.Name, u.Email)
}

func main() {
	var rw ReadWriter = &strings.Builder{}
	rw.Write([]byte("hello"))
	var buf [16]byte
	n, _ := rw.Read(buf[:])
	fmt.Println("ReadWriter:", string(buf[:n]))

	var vf ValidatedFormatter = User{Name: "Alice", Email: "alice@example.com"}
	if err := vf.Validate(); err != nil {
		fmt.Println("validation failed:", err)
		return
	}
	fmt.Println("formatted:", vf.Format())
}
```

## Step-by-step execution

For `rw.Read(buf[:])` where `rw` is `ReadWriter` backed by `*strings.Builder`:

1. `ReadWriter` embeds `Reader`, which requires `Read(p []byte) (n int, err error)`.
2. At compile time, `*strings.Builder` is checked for `Read` and `Write` methods.
3. `*strings.Builder` has both. Assignment is valid.
4. At runtime, `rw.Read(buf[:])` does a method dispatch through the `itab` to `(*strings.Builder).Read`.
5. `buf[:n]` contains `"hello"`.

## Common mistakes

- Mistake: Naming a composed interface inconsistently (e.g., `WriterReader` instead of `ReadWriter`).
  - Fix: Follow the standard pattern: `ReadWriteCloser`, `ReadSeeker`, etc.

- Mistake: Embedding an interface in a struct and expecting the struct to satisfy the interface without implementing the methods.
  - Fix: Embedding an interface in a struct gives the struct access to the interface's methods only if the embedded field is set to a concrete implementation.

- Mistake: Creating deep interface embedding chains (A embeds B embeds C embeds D).
  - Fix: Keep embedding at most 2-3 levels. Deep chains are hard to read.

- Mistake: Embedding the same interface twice (e.g., `A { io.Reader; io.ReadWriter }`).
  - Fix: Redundant embedding is harmless but unnecessary. Embed only the smallest set needed.

## Debugging walkthrough

This code fails to compile:

```go
type A interface {
    Read() int
}

type B interface {
    Read() string
}

type C interface {
    A
    B
}
```

**Symptom**: `duplicate method Read` or `A and B conflict`.

**Root cause**: Both `A` and `B` declare `Read` with different signatures. The compiler cannot merge them into `C`.

**Fix**: Rename one of the methods, or do not embed both in the same interface.

## Production notes

- **Embed standard interfaces** (`io.Reader`, `io.Writer`, etc.) rather than redefining them. Consistency with the standard library makes your code familiar.
- **Small composed interfaces** are preferred over large monolithic ones. `ReadWriter` is better than a 10-method `IO` interface.
- **Document composed interfaces**: if you create `StorerLoader` that embeds `Storer` and `Loader`, document the contract.
- **Interface embedding in structs** is used for dependency injection: embed `http.ResponseWriter` in a test double.

## Performance implications

- **Zero runtime overhead**: Interface embedding is only a compile-time merging of method sets.
- **Larger interfaces** mean larger `itab` entries (more methods to look up), but the difference is negligible.
- **Type assertion cost** is the same regardless of interface size.

## Practice task

Define an interface `Saver` with method `Save() error` and an interface `Loader` with method `Load() error`. Compose them into `Storage` interface. Create a `FileStorage` struct that implements both. Write a function `Backup(s Storage)` that calls `Save()` and then `Load()`. In `main()`, demonstrate with `FileStorage`.

## Tests / verification

```bash
go run ./curriculum/modules/05-types-interfaces-packages-modules/lessons/10-interface-embedding
go test ./curriculum/modules/05-types-interfaces-packages-modules/lessons/10-interface-embedding
```

## Review questions

1. How do you compose an interface that requires both `Read` and `Close`?
2. What happens if two embedded interfaces both declare `Close() error`?
3. What is the naming convention for a composed interface that combines `Reader`, `Writer`, and `Closer`?
4. Can you embed an interface in a struct? What does that mean at runtime?
5. Why is interface embedding better than defining a single large interface with all methods?

## NEXT UP

Stringer — the `fmt.Stringer` interface and custom string formatting.
