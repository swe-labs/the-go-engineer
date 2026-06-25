# Repository and service seam

## Learning objective

Abstract database access behind repository interfaces, implement a service layer that depends on interfaces, and test business logic with fake repositories using dependency injection.

## Why this matters

Direct `database/sql` calls scattered across business logic make code untestable — every test needs a real database. By defining a repository interface and injecting it into a service layer, you can test business logic with in-memory fakes, swap databases without changing business code, and isolate unit tests from integration tests.

## Mental model

```
┌──────────────────────────────────────────────────┐
│                   Service Layer                  │
│   (business logic, validation, orchestration)    │
│                  depends on                      │
│           ┌──────────────────────┐               │
│           │   Repository interface                │
│           └──────────────────────┘               │
│                     ▲                            │
│                     │ implements                  │
│           ┌──────────────────────┐               │
│           │  SQLiteRepository    │               │
│           │  (or FakeRepository) │               │
│           └──────────────────────┘               │
└──────────────────────────────────────────────────┘
```

The service never imports `database/sql`. It only knows the interface. In production, you inject `SQLiteRepository`. In tests, you inject `FakeRepository`.

## Core idea

Define the data access contract as an interface:

```go
type UserRepository interface {
    GetByID(ctx context.Context, id int) (*User, error)
    List(ctx context.Context) ([]User, error)
    Create(ctx context.Context, user *User) error
    Delete(ctx context.Context, id int) error
}
```

Implement it with SQLite:

```go
type userRepository struct {
    db *sql.DB
}

func (r *userRepository) GetByID(ctx context.Context, id int) (*User, error) {
    // ... SQL query
}
```

Implement a fake for tests:

```go
type fakeUserRepo struct {
    users map[int]User
    nextID int
}

func (f *fakeUserRepo) GetByID(_ context.Context, id int) (*User, error) {
    u, ok := f.users[id]
    if !ok { return nil, sql.ErrNoRows }
    return &u, nil
}
```

## Under the hood

The repository pattern is an application of the Dependency Inversion Principle (DIP): high-level modules (service) should not depend on low-level modules (database). Both should depend on abstractions (interface).

In Go, interfaces are satisfied implicitly. The test fake doesn't need to import the SQLite driver. It only needs to match the method signatures.

## How Go uses it

- **Clean architecture**: Ports (interfaces) and adapters (implementations) pattern.
- **Hexagonal architecture**: Repository is a "driven port" (outbound adapter).
- **Testing**: Swap real DB for fake in unit tests; use real DB in integration tests.
- **Multiple backends**: PostgreSQL for production, SQLite for local dev — swap via repository interface.

## Go example

```go
package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	_ "modernc.org/sqlite"
)

type User struct {
	ID    int
	Name  string
	Email string
}

type UserRepository interface {
	GetByID(ctx context.Context, id int) (*User, error)
	List(ctx context.Context) ([]User, error)
	Create(ctx context.Context, name, email string) (*User, error)
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetByID(ctx context.Context, id int) (*User, error) {
	u := &User{}
	err := r.db.QueryRowContext(ctx, "SELECT id, name, email FROM users WHERE id = ?", id).
		Scan(&u.ID, &u.Name, &u.Email)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (r *userRepository) List(ctx context.Context) ([]User, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, name, email FROM users ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, rows.Err()
}

func (r *userRepository) Create(ctx context.Context, name, email string) (*User, error) {
	result, err := r.db.ExecContext(ctx, "INSERT INTO users (name, email) VALUES (?, ?)", name, email)
	if err != nil {
		return nil, err
	}
	id, _ := result.LastInsertId()
	return r.GetByID(ctx, int(id))
}

type UserService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Register(ctx context.Context, name, email string) (*User, error) {
	if name == "" {
		return nil, fmt.Errorf("name cannot be empty")
	}
	if email == "" {
		return nil, fmt.Errorf("email cannot be empty")
	}
	return s.repo.Create(ctx, name, email)
}

func (s *UserService) GetUser(ctx context.Context, id int) (*User, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *UserService) ListUsers(ctx context.Context) ([]User, error) {
	return s.repo.List(ctx)
}

func main() {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	db.Exec(`CREATE TABLE users (id INTEGER PRIMARY KEY, name TEXT NOT NULL, email TEXT NOT NULL UNIQUE)`)

	repo := NewUserRepository(db)
	svc := NewUserService(repo)

	alice, _ := svc.Register(context.Background(), "Alice", "alice@example.com")
	bob, _ := svc.Register(context.Background(), "Bob", "bob@example.com")
	fmt.Printf("Registered: %s (id=%d), %s (id=%d)\n", alice.Name, alice.ID, bob.Name, bob.ID)

	users, _ := svc.ListUsers(context.Background())
	fmt.Printf("Total users: %d\n", len(users))
}
```

## Step-by-step execution

For `svc.Register(ctx, "Alice", "alice@example.com")`:

1. `Register` validates name and email are non-empty.
2. Calls `repo.Create(ctx, "Alice", "alice@example.com")`.
3. `repo.Create` executes `INSERT INTO users`.
4. Gets the `LastInsertId`.
5. Calls `repo.GetByID` to return the full user with generated ID.
6. Returns `*User` to the caller.

In a test, step 3-5 would use a fake repository that stores users in a `map[int]User` — no database needed.

## Common mistakes

- **Leaking `*sql.DB` into the service layer**: The service should never import `database/sql`. The repository interface and domain types should be in a separate package.
- **Too many methods on the repository**: Keep interfaces focused. Split into `UserReader` and `UserWriter` if needed.
- **Fake repositories that don't match real behavior**: The fake should replicate error conditions (`ErrNoRows`, duplicate key errors) to test edge cases.
- **Over-abstraction**: If the app has only a few queries, a full repository pattern may be overkill. Start simple and extract when needed.
- **Missing context propagation**: Repository methods must accept `context.Context` and pass it to `QueryContext`/`ExecContext`.

## Debugging walkthrough

**Scenario**: Tests pass with fake repository but fail in production.

```go
// Fake repo — always succeeds
func (f *fakeUserRepo) Create(name, email string) *User {
    f.nextID++
    f.users[f.nextID] = User{ID: f.nextID, Name: name, Email: email}
    return &User{ID: f.nextID, Name: name, Email: email}
}
```

**Problem**: The fake always succeeds. The real repo can fail on UNIQUE constraint, FK violation, or NOT NULL. Production errors are untested.

**Fix**: Add error injection to the fake:

```go
type fakeUserRepo struct {
    users    map[int]User
    nextID   int
    shouldFail bool
}

func (f *fakeUserRepo) Create(ctx context.Context, name, email string) (*User, error) {
    if f.shouldFail {
        return nil, fmt.Errorf("database unavailable")
    }
    // ...
}
```

## Production notes

- **Package structure**: Put interfaces in a `repository` or `store` package, implementations in `sqlite` or `postgres` subpackage.
- **Constructor injection**: Use `NewXxx(db) *Xxx` pattern. Don't use global `init()` or `sync.Once` for repositories.
- **Transaction handling**: The repository interface should expose `BeginTx(ctx) (Tx, error)` for multi-operation transactions. The transaction type should also be an interface.
- **Mock generation**: Use `mockgen` (gomock) or `moq` to generate mocks from interfaces. Or write small fakes by hand — they're simpler to understand.

## Performance implications

- The interface indirection adds zero cost at the call site (Go inlines interface calls when the concrete type is known).
- A fake repository in tests eliminates database round-trips — tests run in milliseconds instead of seconds.
- The repository pattern encourages batched queries (List returns all rows) rather than N+1 loops, improving database performance.

## Practice task

1. Define a `ProductRepository` interface with `GetByID`, `List`, `Create`, `Update`.
2. Implement it with SQLite.
3. Create a `ProductService` with `CreateProduct`, `GetProduct`, `UpdatePrice` (with validation).
4. Write a fake repository.
5. Test the service using the fake, covering success and error cases.

## Tests / verification

```bash
go run ./curriculum/modules/09-sql-postgres-persistence/lessons/24-repository-and-service-seam
go test ./curriculum/modules/09-sql-postgres-persistence/lessons/24-repository-and-service-seam
```

## Review questions

1. What is the Dependency Inversion Principle and how does the repository pattern apply it?
2. Why should the service layer not import `database/sql`?
3. What advantage does a fake repository have over mocking with a library?
4. How does the repository pattern improve test performance?
5. When might the repository pattern be unnecessary overhead?

## NEXT UP

Integration tests — testing with a real database, setup/teardown, and build tags.
