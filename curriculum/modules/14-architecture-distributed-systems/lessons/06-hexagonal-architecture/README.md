# Hexagonal architecture

## Learning objective

Design a hexagonal architecture with ports (interfaces) and adapters (implementations), separate core domain from infrastructure concerns, and swap adapters for testing without modifying core logic.

## Why this matters

Traditional layered architecture (handler -> service -> repository) has a dangerous property: the core domain depends on the database. If the database changes or you need to add a new input source (CLI, message queue, test harness), you must modify the core. Hexagonal architecture inverts this: the core defines ports (interfaces) and infrastructure provides adapters (implementations). The core has no imports pointing outward. This makes the core independently testable, swappable, and understandable.

## Mental model

Imagine a computer processor. The processor (core) defines ports: memory bus, I/O bus, interrupt lines. Peripherals (keyboard, monitor, disk) implement adapters that plug into these ports. The processor does not know which peripherals are connected. It only knows the port protocol. You can replace the keyboard with a different model without changing the processor.

In hexagonal architecture, the core domain logic is the processor. The ports are Go interfaces. The adapters are concrete implementations for HTTP, databases, message queues, and test mocks. The core imports only standard library and domain packages. Infrastructure imports the core.

## Core idea

Hexagonal architecture (also called ports and adapters) organizes code into three zones:

1. **Core domain**: pure business logic. Defines ports (interfaces) for inbound and outbound communication. Has no infrastructure imports.
2. **Driving adapters**: inbound adapters that trigger the core. Examples: HTTP handlers, gRPC servers, CLI commands, message consumers. They convert external input into core method calls.
3. **Driven adapters**: outbound adapters that the core calls to perform side effects. Examples: database repositories, email senders, payment gateways, HTTP clients. They implement core-defined interfaces.

The core never imports an adapter. Adapters always import and implement core interfaces. The main function wires adapters to ports at startup.

A driving port is a use case interface that the core exposes:
```go
type UserService interface {
    RegisterUser(id, email string) error
}
```

A driven port is a dependency interface that the core requires:
```go
type UserRepository interface {
    Save(user User) error
    FindByID(id string) (User, error)
}
```

## Under the hood

Go's implicit interface satisfaction is what makes hexagonal architecture practical. An adapter in package `postgres` implements an interface defined in package `domain` without importing any interface registry. The compiler checks conformance at the assembly point in `main`.

```
┌──────────────┐     ┌──────────────┐     ┌─────────────────┐
│  HTTP Handler │────>│  Core Domain  │────>│  UserRepository  │
│  (driving     │     │  (ports)      │     │  (driven port)   │
│   adapter)    │     └──────────────┘     └────────┬────────┘
└──────────────┘                                    │
                                           ┌────────▼────────┐
                                           │  PostgresAdapter │
                                           │  (driven adapter)│
                                           └─────────────────┘
```

The core depends on `UserRepository` (the port). `postgres` package implements `UserRepository` (the adapter). The core does not import `postgres`. This means:
- Tests can swap `PostgresAdapter` with `InMemoryAdapter` without touching the core.
- Adding a new database (MySQL, MongoDB) means writing a new adapter.
- The core can be compiled and unit-tested without any database driver linked.

## How Go uses it

The Go standard library follows hexagonal principles:

- `net/http` defines `Handler` (a driving port). Any type with a `ServeHTTP` method is an adapter. `http.HandlerFunc` adapts a function. `http.ServeMux` adapts routing logic.
- `io.Reader` and `io.Writer` are driven ports. `os.File`, `bytes.Buffer`, `strings.Reader` are adapters.
- `database/sql/driver` defines ports for SQL databases. `github.com/lib/pq` and `github.com/go-sql-driver/mysql` are adapters.

Production Go projects use hexagonal architecture extensively. The core service package imports only domain types and standard library. All infrastructure adapters live in separate packages.

## Go example

```go
package main

import (
	"errors"
	"fmt"
	"sync"
)

type User struct {
	ID    string
	Email string
}

// Driven port
type UserRepository interface {
	Save(user User) error
	FindByID(id string) (User, error)
}

// Driven port
type EmailService interface {
	SendWelcome(email string) error
}

// Core domain
type UserService struct {
	repo  UserRepository
	email EmailService
}

func NewUserService(repo UserRepository, email EmailService) *UserService {
	return &UserService{repo: repo, email: email}
}

func (s *UserService) RegisterUser(id, email string) error {
	if email == "" {
		return errors.New("email cannot be empty")
	}
	user := User{ID: id, Email: email}
	if err := s.email.SendWelcome(email); err != nil {
		return fmt.Errorf("welcome email failed: %w", err)
	}
	if err := s.repo.Save(user); err != nil {
		return fmt.Errorf("save failed: %w", err)
	}
	return nil
}

// Driven adapter: in-memory repository
type InMemoryUserRepoAdapter struct {
	mu    sync.RWMutex
	users map[string]User
}

func NewInMemoryUserRepoAdapter() *InMemoryUserRepoAdapter {
	return &InMemoryUserRepoAdapter{users: make(map[string]User)}
}

func (a *InMemoryUserRepoAdapter) Save(user User) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.users[user.ID] = user
	return nil
}

func (a *InMemoryUserRepoAdapter) FindByID(id string) (User, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	u, ok := a.users[id]
	if !ok {
		return User{}, errors.New("user not found")
	}
	return u, nil
}

// Driven adapter: console email sender
type ConsoleEmailAdapter struct{}

func (ConsoleEmailAdapter) SendWelcome(email string) error {
	fmt.Printf("Welcome email sent to %s\n", email)
	return nil
}

func main() {
	repo := NewInMemoryUserRepoAdapter()
	email := ConsoleEmailAdapter{}
	svc := NewUserService(repo, email)

	if err := svc.RegisterUser("u1", "alice@example.com"); err != nil {
		fmt.Println("Error:", err)
		return
	}
	user, _ := repo.FindByID("u1")
	fmt.Printf("Registered user: %s (%s)\n", user.ID, user.Email)
}
```

## Step-by-step execution

For `svc.RegisterUser("u1", "alice@example.com")`:

1. `RegisterUser` validates the email is not empty.
2. A `User` value is constructed.
3. `s.email.SendWelcome(email)` is called. The `ConsoleEmailAdapter.SendWelcome` prints a message and returns nil.
4. `s.repo.Save(user)` is called. `InMemoryUserRepoAdapter.Save` stores the user in the map.
5. Nil is returned.

The test swap: in tests, pass a `mockEmailAdapter` and the same `InMemoryUserRepoAdapter`. The core does not know which adapters are connected. To test the email-failure case, the mock returns an error. The core handles it correctly without any infrastructure setup.

## Common mistakes

- **Core imports infrastructure**: the core package imports `postgres` or `redis`. This breaks the dependency rule. The core should import only domain types and standard library. Use `go vet` or custom analyzers to enforce this.
- **Leaking adapter types into the core**: the repository interface returns `sql.Row` or the email interface takes `*smtp.Client`. The core should use only domain types in its interface definitions.
- **Too many ports**: defining a separate interface for every method creates interface pollution. One port per bounded context (e.g., `UserRepository` with 3-5 methods) is the right granularity.
- **Ports defined in adapter packages**: the `UserRepository` interface is defined in `postgres/` and the core imports it. This makes the core depend on the adapter package name. Always define ports in the core.

## Debugging walkthrough

```go
// Core imports postgres adapter (WRONG)
import "project/postgres"

type UserService struct {
    repo *postgres.UserRepo
}
```

**Symptom**: The core cannot be tested without Postgres. Adding a second database requires changing the core.

**Root cause**: The core imports a concrete adapter type instead of defining a port interface.

**Fix**: Extract `UserRepository` interface in the core. Define `postgres.UserRepo` as an adapter. Wire in main:

```go
// core
type UserRepository interface {
    Save(User) error
}

// main
repo := postgres.NewRepo()
svc := core.NewUserService(repo)
```

## Production notes

- **Adapter directory structure**: place adapters in an `adapters/` directory at the module root, with subdirectories for each adapter type: `adapters/postgres/`, `adapters/redis/`, `adapters/s3/`.
- **Testing with adapter swaps**: the core tests use in-memory adapters. Integration tests use real adapters (real Postgres, real email test server). The integration tests are in a different package or build-tagged file.
- **Graceful adapter degradation**: if an adapter (e.g., email sender) fails, the core should handle it gracefully. Define error types in the port interface that the core can inspect: `var ErrEmailSendFailed = errors.New("email send failed")`.
- **Adapter lifecycle**: adapters may need initialization (connection pools) and shutdown (graceful close). Provide `New` and `Close` methods on adapters. The main function calls `Close` on `os.Signal`.

## Performance implications

- **Interface dispatch overhead**: each port call goes through an interface method dispatch. In most applications, this cost is immeasurable. In hot paths (millions of calls per second), consider inlining or using concrete types.
- **Adapter allocation patterns**: adapters that allocate on every call (e.g., JSON serialization) can cause GC pressure. Profile adapter implementations separately from core logic.
- **Test adapter performance**: in-memory adapters are significantly faster than real infrastructure. Tests that use in-memory adapters run in milliseconds vs seconds for integration tests.

## Practice task

Add a third driven port to the core: a `PaymentGateway` interface. Implement a `StripePaymentAdapter` and a `MockPaymentAdapter` (for tests). Modify `UserService.RegisterUser` to accept a payment after sending the welcome email. Write a table-driven test that uses the mock adapter to verify success and payment failure scenarios.

## Tests / verification

```bash
go run ./curriculum/modules/14-architecture-distributed-systems/lessons/06-hexagonal-architecture
go test ./curriculum/modules/14-architecture-distributed-systems/lessons/06-hexagonal-architecture
```

The tests verify that user registration sends a welcome email and persists the user using the in-memory repository adapter. Tests use a mock email adapter to test both success and failure paths. After completing the practice task, add tests for the payment flow.

## Review questions

1. What is the difference between a driving adapter and a driven adapter?
2. Why should ports be defined in the core package rather than the adapter package?
3. How does Go's implicit interface satisfaction make hexagonal architecture practical?
4. How would you add a new database (e.g., MongoDB) to a hexagonal Go project?
5. What is the role of the main function in hexagonal architecture?

## NEXT UP

Domain modeling -- designing entities, value objects, and aggregates in Go for complex business domains.
