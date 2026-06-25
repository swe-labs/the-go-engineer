# Fakes before mocks

## Learning objective

Implement in-memory fake implementations of interfaces for testing, distinguish between stubs, fakes, and mocks, and choose the right testing double for the situation.

## Why this matters

When a function depends on an external system — a database, an HTTP API, a file system — you cannot test it directly without setting up that system. Testing doubles replace those dependencies with controllable substitutes. Beginners reach for mocking frameworks, but in Go, in-memory fakes are often simpler, faster, and more maintainable. Understanding the tradeoffs between stubs, fakes, and mocks is essential for writing tests that are both effective and resilient to change.

## Mental model

A testing double is a stunt double for a real dependency. There are three main kinds:

- **Stub**: Returns hardcoded answers. The simplest. No logic.
- **Fake**: Has a working (but simplified) implementation. An in-memory database is a fake. It behaves like the real thing but runs in process.
- **Mock**: Records how it is called and can verify expectations. "Expect Add(2,3) to be called exactly once."

```
Real system:   PostgreSQL database (slow, requires connection)
Fake:          map[string]string in memory (fast, no setup)
Mock:          testify/mock that records calls and verifies expectations
```

## Core idea

An **in-memory fake** implements the same interface as the real dependency but stores data in memory instead of an external system. It has real behavior — you can call `Get`, `Set`, `Delete` — but it runs in process, has no side effects, and is reset between tests.

```go
// Store interface that both real and fake implementations satisfy.
type Store interface {
    Get(key string) (string, error)
    Set(key, value string) error
    Delete(key string) error
}

// FakeStore is an in-memory fake.
type FakeStore struct {
    data map[string]string
}

func (f *FakeStore) Get(key string) (string, error) {
    val, ok := f.data[key]
    if !ok {
        return "", fmt.Errorf("key %q not found", key)
    }
    return val, nil
}

func (f *FakeStore) Set(key, value string) error {
    f.data[key] = value
    return nil
}

func (f *FakeStore) Delete(key string) error {
    delete(f.data, key)
    return nil
}
```

`FakeStore` is a fully functional `Store` that tests can use without setting up a database.

## Under the hood

A fake is not magic — it is a concrete type that implements an interface. Because Go interfaces are satisfied implicitly, any type with the right methods can serve as a fake. The fake lives in the test package (or a shared `fakes` package) and is maintained alongside the interface.

The key insight: a fake is a valid implementation for *any* test that uses the interface. You write it once and reuse it across hundreds of tests. A mock, by contrast, is configured per test.

## How Go uses it

The Go standard library provides fakes in `net/http/httptest`:

- `httptest.NewServer` starts a local HTTP server (real TCP, no mocking).
- `httptest.NewRecorder` implements `http.ResponseWriter` for testing handlers — this is a fake.

The standard library generally prefers fakes and real implementations over mocks. The `io` package has `io.Pipe` (a fake pipe), `io.MultiReader` (a fake combining reader), and `io.TeeReader` (a fake tee).

## Go example

```go
package main

import (
	"fmt"
)

// UserStore stores user information.
type UserStore interface {
	Save(name string, age int) error
	Find(name string) (int, error)
}

// RealUserStore connects to a database (simulated).
type RealUserStore struct {
	db map[string]int
}

func (r *RealUserStore) Save(name string, age int) error {
	r.db[name] = age
	return nil
}

func (r *RealUserStore) Find(name string) (int, error) {
	age, ok := r.db[name]
	if !ok {
		return 0, fmt.Errorf("user %q not found", name)
	}
	return age, nil
}

// Greeter greets users from a store.
type Greeter struct {
	store UserStore
}

func NewGreeter(store UserStore) *Greeter {
	return &Greeter{store: store}
}

func (g *Greeter) Greet(name string) (string, error) {
	age, err := g.store.Find(name)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("Hello, %s! You are %d years old.", name, age), nil
}

func main() {
	store := &RealUserStore{db: make(map[string]int)}
	store.Save("Alice", 30)
	g := NewGreeter(store)
	msg, _ := g.Greet("Alice")
	fmt.Println(msg)
}
```

```go
// main_test.go
package main

import "testing"

// FakeUserStore is an in-memory fake for UserStore.
type FakeUserStore struct {
	users map[string]int
}

func NewFakeUserStore() *FakeUserStore {
	return &FakeUserStore{users: make(map[string]int)}
}

func (f *FakeUserStore) Save(name string, age int) error {
	f.users[name] = age
	return nil
}

func (f *FakeUserStore) Find(name string) (int, error) {
	age, ok := f.users[name]
	if !ok {
		return 0, fmt.Errorf("user %q not found", name)
	}
	return age, nil
}

func TestGreeter(t *testing.T) {
	store := NewFakeUserStore()
	store.Save("Alice", 30)

	greeter := NewGreeter(store)
	msg, err := greeter.Greet("Alice")
	if err != nil {
		t.Fatalf("Greet failed: %v", err)
	}
	want := "Hello, Alice! You are 30 years old."
	if msg != want {
		t.Errorf("Greet = %q; want %q", msg, want)
	}
}

func TestGreeterNotFound(t *testing.T) {
	store := NewFakeUserStore()
	greeter := NewGreeter(store)
	_, err := greeter.Greet("Bob")
	if err == nil {
		t.Fatal("expected error for unknown user")
	}
}

func TestFakeStoreRoundTrip(t *testing.T) {
	store := NewFakeUserStore()

	if err := store.Save("Bob", 25); err != nil {
		t.Fatalf("Save failed: %v", err)
	}
	age, err := store.Find("Bob")
	if err != nil {
		t.Fatalf("Find failed: %v", err)
	}
	if age != 25 {
		t.Errorf("age = %d; want 25", age)
	}
}
```

## Step-by-step execution

Running `TestGreeter`:

1. `NewFakeUserStore()` creates an empty in-memory store.
2. `store.Save("Alice", 30)` adds Alice to the map.
3. `NewGreeter(store)` creates a Greeter with the fake store.
4. `greeter.Greet("Alice")` calls `store.Find("Alice")`, which returns `30`.
5. The greeting is composed: `"Hello, Alice! You are 30 years old."`.
6. Assertion passes.

The same Greeter code, when used in production, receives a `RealUserStore` connected to an actual database. The fake replaces the database without any change to the business logic.

## Common mistakes

- **Fakes with bugs.** A fake that does not correctly simulate the real system is worse than no fake. Test your fake against the same contract as the real implementation.

- **Fakes that are too complex.** If a fake needs its own configuration or state machine, it becomes a maintenance burden. Keep fakes simple.

- **Choosing mocks over fakes for state verification.** If you only need to check that data was stored, a fake's state is sufficient. Mocks are for verifying *that* a method was called with specific arguments, not for checking the result.

- **Not using interfaces.** Fakes only work when the dependency is expressed as an interface. If code depends on concrete types, you cannot substitute a fake.

- **Fakes with external dependencies.** An in-memory fake should have zero dependencies on external systems. If it imports a database driver, it is not a fake.

## Debugging walkthrough

A test using a fake passes but the production code fails:

```go
// FakeStore always returns nil error for Save
func (f *FakeStore) Save(key, value string) error {
    f.data[key] = value
    return nil // never fails
}
```

**Symptom**: Production `RealStore.Save` occasionally fails (disk full, connection timeout). The fake never reproduces this, so error handling in the caller is untested.

**Investigation**: The fake is too perfect. It does not simulate failure modes.

**Root cause**: The test never exercises the error path because the fake never returns errors.

**Fix**: Make the fake capable of simulating errors:
```go
type FakeStore struct {
    data    map[string]string
    failOn map[string]error // inject failures per key or operation
}
```

## Production notes

- **Fakes in shared test packages.** Put reusable fakes in a `fakes` subpackage or in an `internal/testutil` package.
- **Fakes are not for verification.** A fake's purpose is to replace a dependency, not to verify behavior. For verification, assert on the fake's state or use a mock.
- **Maintain fakes alongside interfaces.** When the interface changes, the fake must be updated. This is a compile-time check — the fake will not compile if it no longer satisfies the interface.
- **Fakes can be slow if poorly implemented.** An in-memory fake using a `map` is O(1). A fake that does redundant work can be slower than the real system.

## Performance implications

- In-memory fakes are fast — typically microseconds per operation.
- No network IO, no disk IO, no serialization overhead.
- Fakes can be shared across tests without performance cost (each test creates its own instance).
- Because fakes are fast, you can run thousands of tests that use them in milliseconds.

## Practice task

Define an interface `Calculator` with methods `Add(a, b int) int` and `Multiply(a, b int) int`. Write a `RealCalculator` that performs the operations. Write a `FakeCalculator` that stores the last operation and returns configured results. Write a function `Double(c Calculator, n int) int` that multiplies `n` by 2 using `c.Add(n, n)`. Test `Double` with both the real and fake calculator.

## Tests / verification

```bash
go run ./curriculum/modules/06-testing-debugging-refactoring/lessons/09-fakes-before-mocks
go test -v ./curriculum/modules/06-testing-debugging-refactoring/lessons/09-fakes-before-mocks
```

## Review questions

1. What is the difference between a stub, a fake, and a mock?
2. Why are in-memory fakes preferred over mocking frameworks in many Go tests?
3. When would a fake be insufficient and you would need a mock instead?
4. How does Go's implicit interface satisfaction make fakes easier to create?
5. What is the risk of a fake that is buggy or does not match real behavior?

## NEXT UP

Mocking tradeoffs — when mocking external services is worth the complexity.
