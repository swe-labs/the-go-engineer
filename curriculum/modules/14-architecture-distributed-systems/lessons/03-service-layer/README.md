# Service layer

## Learning objective

Design and implement a service layer that isolates business logic from transport and storage, use constructor-based dependency injection to make services testable, and recognize when a service layer adds value versus when it is a thin pass-through.

## Why this matters

HTTP handlers, gRPC interceptors, and CLI commands are transport mechanisms. Databases, caches, and APIs are storage mechanisms. Business logic -- the rules that make your application valuable -- should depend on neither. When business logic is scattered across handlers and repository calls, adding a new feature means modifying code in every layer. The service layer centralizes business rules into one place that can be tested independently of HTTP and databases.

## Mental model

The service layer is the brain of the application. The handler is the mouth (talks to the outside world). The repository is the hands (talks to storage). The brain decides what to do (business rules), the mouth receives input and delivers output, and the hands execute storage operations. The brain should not do the hands' job (SQL queries) or the mouth's job (HTTP parsing). The brain should be testable without the mouth or hands -- inject mock dependencies and test the business logic in isolation.

## Core idea

The service layer is a Go struct that:

- Embodies business operations as methods (`RegisterUser`, `PlaceOrder`, `ProcessPayment`).
- Accepts its dependencies via constructor injection as interface values.
- Contains zero transport or storage code.
- Returns domain types and domain errors.

Constructor injection means the service does not create its own dependencies. Instead, they are passed in:

```go
type UserService struct {
    repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
    return &UserService{repo: repo}
}
```

This makes the service testable: you pass a mock repository in tests and the real repository in production.

The service layer should justify its existence by containing meaningful business logic. A service method that just calls `repo.Insert` with no transformation or validation is a pass-through and adds no value. The service layer is where you enforce business rules, coordinate multiple repositories or external services, and translate between transport-layer DTOs and domain types.

## Under the hood

Go supports the service layer pattern through two language features:

- **Interfaces**: allow the service to define its dependencies abstractly. The concrete type satisfaction is implicit -- any type with matching method signatures satisfies the interface. This is checked at compile time when the concrete type is assigned to the interface variable.
- **Package visibility**: the service package exports only business types and interfaces. Implementation details are unexported. The `internal/` convention prevents external consumers from importing implementation packages directly.

The service layer does not import any transport package (`net/http`, `google.golang.org/grpc`) or any storage package (`database/sql`, `go.mongodb.org/mongo-driver`). It is a pure Go package that could be used in any context.

## How Go uses it

The service layer pattern is universal in production Go services:

- Stripe's API handlers call the service layer for every operation. The handler parses the request and calls a method on the appropriate service. The service performs business logic and calls repositories or external APIs.
- Kubernetes' API server has a service layer between the HTTP handler and the storage backend. The registry layer is similar to a service layer.
- Docker's engine API uses a service layer for image and container operations.

The pattern is consistent: handler calls service calls repository (or external API). The service layer is the only layer that contains business logic. Handlers and repositories are mechanical.

## Go example

```go
package main

import (
	"errors"
	"fmt"
)

type User struct {
	ID    string
	Name  string
	Email string
}

type UserRepository interface {
	Save(user User) error
	FindByID(id string) (User, error)
}

var errUserNotFound = errors.New("user not found")

type UserService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Register(name, email string) (User, error) {
	if name == "" {
		return User{}, errors.New("name cannot be empty")
	}
	if email == "" {
		return User{}, errors.New("email cannot be empty")
	}
	user := User{ID: fmt.Sprintf("usr_%s", email), Name: name, Email: email}
	if err := s.repo.Save(user); err != nil {
		return User{}, fmt.Errorf("save failed: %w", err)
	}
	return user, nil
}

func (s *UserService) GetUser(id string) (User, error) {
	return s.repo.FindByID(id)
}

type inMemoryUserRepo struct {
	users map[string]User
}

func newInMemoryUserRepo() *inMemoryUserRepo {
	return &inMemoryUserRepo{users: make(map[string]User)}
}

func (r *inMemoryUserRepo) Save(user User) error {
	r.users[user.ID] = user
	return nil
}

func (r *inMemoryUserRepo) FindByID(id string) (User, error) {
	user, ok := r.users[id]
	if !ok {
		return User{}, errUserNotFound
	}
	return user, nil
}

func main() {
	repo := newInMemoryUserRepo()
	svc := NewUserService(repo)

	user, _ := svc.Register("Alice", "alice@example.com")
	fmt.Printf("Registered: %+v\n", user)

	fetched, _ := svc.GetUser(user.ID)
	fmt.Printf("Fetched: %+v\n", fetched)
}
```

## Step-by-step execution

For `svc.Register("Alice", "alice@example.com")`:

1. `Register` checks that `name` is not empty. "Alice" passes.
2. `Register` checks that `email` is not empty. "alice@example.com" passes.
3. A `User` struct is constructed with a generated ID.
4. `s.repo.Save(user)` is called. The concrete `inMemoryUserRepo.Save` stores the user in the map.
5. The `User` struct is returned with a nil error.

For `svc.GetUser("usr_alice@example.com")`:

1. `s.repo.FindByID(id)` is called.
2. `inMemoryUserRepo.FindByID` looks up the map. Found.
3. The `User` struct is returned.

The key design property: neither `Register` nor `GetUser` know about the map, the HTTP request, or any transport. They only know about `UserRepository`. A test can replace the real repo with a mock and test validation rules without setting up a database.

## Common mistakes

- **Putting business logic in HTTP handlers**: the handler parses the request, calls the database, formats the response, and handles errors. When the same logic is needed for a gRPC endpoint or a background job, the handler code cannot be reused. Business logic must live in the service layer.
- **Making the service layer a thin pass-through**: the handler calls `service.CreateUser` which calls `repo.InsertUser` with no transformation, validation, or business logic. The service layer adds complexity without value. A service layer must justify its existence by containing business rules.
- **Letting services become god objects**: `internal/service/user.go` has 3000 lines with 50 methods covering user CRUD, billing, notifications, and reporting. The service layer should be split by bounded context, not by entity.
- **Constructing dependencies inside the service**: calling `sql.Open` or `redis.NewClient` inside `NewUserService`. This couples the service to specific infrastructure and makes testing impossible.

## Debugging walkthrough

Consider a service that directly opens a database connection:

```go
type UserService struct {}
func (s UserService) GetUser(id string) User {
    db, _ := sql.Open("postgres", connStr)
    row := db.QueryRow("SELECT * FROM users WHERE id = $1", id)
    // ...
}
```

**Symptom**: Testing `GetUser` requires a running Postgres instance. The test takes 10 seconds to start. On CI, the test flakes because the Postgres container is not ready.

**Investigation**: Look at the service constructor. There is none -- the service creates its own dependencies internally.

**Root cause**: The service is coupled to Postgres. It cannot be tested without Postgres.

**Fix**: Extract a `UserRepository` interface and inject it via constructor:

```go
func NewUserService(repo UserRepository) *UserService {
    return &UserService{repo: repo}
}
```

Now tests pass an in-memory implementation and run in milliseconds.

## Production notes

- **One service per bounded context**: an ordering service, a billing service, and a notification service are separate structs, possibly in separate packages. Do not create a single monolithic Service with 100 methods.
- **Service methods return domain errors**: define sentinel errors (`var ErrInsufficientBalance = errors.New("insufficient balance")`) that handlers can inspect and map to appropriate HTTP status codes.
- **Avoid service-to-service calls within the same process**: if service A needs data from service B, extract the shared logic into a third component or use direct repository access. Service-to-service calls within a monolith create artificial overhead.
- **Logging and metrics in services**: accept a `Logger` interface rather than logging directly. This lets you swap loggers or inject a test logger that captures output for assertions.

## Performance implications

- **Indirection cost**: calling through a service that immediately calls a repository adds one stack frame. In a typical web request (dozens of operations), this cost is negligible. In high-throughput batch processing, consider whether the service layer abstraction is necessary.
- **Allocation patterns**: returning domain types from service methods can cause heap allocations. If the service is in a hot path, consider reusing buffers or returning value types when possible.
- **Context propagation**: services should accept `context.Context` as the first parameter of every method. This enables cancellation, deadlines, and tracing propagation.

## Practice task

Extend the `UserService` with a `ChangeEmail(id, newEmail string) error` method that:

1. Validates the new email is not empty.
2. Checks that the new email is different from the current email.
3. Fetches the user and updates the email.
4. Saves the updated user.

Then add a mock-based test that covers success, empty email, and unchanged email cases.

## Tests / verification

```bash
go run ./curriculum/modules/14-architecture-distributed-systems/lessons/03-service-layer
go test ./curriculum/modules/14-architecture-distributed-systems/lessons/03-service-layer
```

The existing tests verify that `UserService.Register` validates inputs and saves the user, and that `UserService.GetUser` finds existing users and returns errors for missing ones. Tests use a mock `UserRepository`. After completing the practice task, add tests for your `ChangeEmail` method.

## Review questions

1. What is the purpose of the service layer? What problems does it solve?
2. Why should a service receive its dependencies via constructor injection rather than creating them internally?
3. How do you test a service method that calls a repository without setting up a real database?
4. What is the difference between a service layer and an anemic service (thin pass-through)?
5. When would you split a 3000-line UserService into multiple services?

## NEXT UP

Repository pattern deep dive -- abstracting data access behind interfaces for testability and flexibility.
