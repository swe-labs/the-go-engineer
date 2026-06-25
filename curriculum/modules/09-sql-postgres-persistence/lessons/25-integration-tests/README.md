# Integration tests

## Learning objective

Write integration tests that exercise real database operations with in-memory SQLite, implement proper setup/teardown, use build tags to separate unit from integration tests, and structure test helpers for maintainable database testing.

## Why this matters

Unit tests with fake repositories verify business logic but not SQL syntax, query correctness, or database constraints. Integration tests against a real database catch mistakes that fakes miss: wrong column names, type mismatches, constraint violations, and unexpected NULL behavior. In production Go services, integration tests are the safety net between schema changes and deployment.

## Mental model

Test pyramid with database tests:

```
         ┌─────┐
         │ E2E │    ← Full system (rare)
         ├─────┤
         │ INT │    ← Database + repos (some)
         ├─────┤
         │ UNIT│    ← Business logic (many)
         └─────┘
```

Unit tests: fast, no DB, fakes, run on every save.
Integration tests: slower, real DB, run in CI or on demand.
E2E tests: slowest, full stack, run before release.

## Core idea

An integration test for database code has three phases:

1. **Setup**: Open database, create schema, insert seed data.
2. **Execute**: Call the function being tested (e.g., `GetByID`).
3. **Teardown**: Close database, clean up resources.

In Go, `t.Cleanup` or `defer` handles teardown. Use build tags to control when integration tests run:

```go
//go:build integration

package main

func TestIntegrationSomething(t *testing.T) {
    db := setupTestDB(t)
    // ...
}
```

Run with: `go test -tags=integration ./...`

## Under the hood

SQLite's `:memory:` mode creates a database that exists only in RAM. Each `sql.Open("sqlite", ":memory:")` gets a fresh, empty database. This is faster than file-based databases and has no cleanup overhead — closing the connection destroys the database.

In PostgreSQL, you would use:
- Test containers (Docker container per test suite)
- Temporary databases (`CREATE DATABASE test_xxx`)
- Or a dedicated test database with rollback (`BEGIN; ... ; ROLLBACK`)

## How Go uses it

- **`//go:build integration`**: Separate integration tests from unit tests.
- **`TestMain`**: Suite-level setup (migrate schema once).
- **Table-driven integration tests**: Same pattern as unit tests, but with real DB queries.
- **Parallel tests**: Each test gets its own `:memory:` database for isolation.

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

//go:build integration

type User struct {
	ID    int
	Name  string
	Email string
}

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetByID(ctx context.Context, id int) (*User, error) {
	u := &User{}
	err := r.db.QueryRowContext(ctx, "SELECT id, name, email FROM users WHERE id = ?", id).
		Scan(&u.ID, &u.Name, &u.Email)
	return u, err
}

func (r *UserRepository) Create(ctx context.Context, name, email string) (*User, error) {
	result, err := r.db.ExecContext(ctx, "INSERT INTO users (name, email) VALUES (?, ?)", name, email)
	if err != nil {
		return nil, err
	}
	id, _ := result.LastInsertId()
	return r.GetByID(ctx, int(id))
}

// Main function (integration test uses same code)
func main() {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	db.Exec(`CREATE TABLE users (id INTEGER PRIMARY KEY, name TEXT NOT NULL, email TEXT NOT NULL UNIQUE)`)

	repo := NewUserRepository(db)
	u, _ := repo.Create(context.Background(), "Alice", "alice@example.com")
	fmt.Printf("Created user %d: %s\n", u.ID, u.Name)
}
```

Integration test file (`main_test.go` without build tag — tests that also work as unit tests, and a separate `integration_test.go` with build tag):

```go
// integration_test.go
//go:build integration

package main

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	_ "modernc.org/sqlite"
)

func setupDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })

	_, err = db.Exec(`CREATE TABLE users (
		id INTEGER PRIMARY KEY,
		name TEXT NOT NULL,
		email TEXT NOT NULL UNIQUE
	)`)
	if err != nil {
		t.Fatal(err)
	}
	return db
}

func TestIntegrationCreateUser(t *testing.T) {
	db := setupDB(t)
	repo := NewUserRepository(db)

	u, err := repo.Create(context.Background(), "Alice", "alice@example.com")
	if err != nil {
		t.Fatal(err)
	}
	if u.Name != "Alice" {
		t.Errorf("name = %q, want %q", u.Name, "Alice")
	}
	if u.ID == 0 {
		t.Error("expected non-zero ID")
	}
}

func TestIntegrationGetByID(t *testing.T) {
	db := setupDB(t)
	repo := NewUserRepository(db)

	created, _ := repo.Create(context.Background(), "Bob", "bob@example.com")
	got, err := repo.GetByID(context.Background(), created.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.Email != "bob@example.com" {
		t.Errorf("email = %q, want %q", got.Email, "bob@example.com")
	}
}

func TestIntegrationGetByIDNotFound(t *testing.T) {
	db := setupDB(t)
	repo := NewUserRepository(db)

	_, err := repo.GetByID(context.Background(), 999)
	if !errors.Is(err, sql.ErrNoRows) {
		t.Errorf("expected ErrNoRows, got %v", err)
	}
}

func TestIntegrationUniqueConstraint(t *testing.T) {
	db := setupDB(t)
	repo := NewUserRepository(db)

	repo.Create(context.Background(), "Alice", "alice@example.com")
	_, err := repo.Create(context.Background(), "Alice2", "alice@example.com")
	if err == nil {
		t.Fatal("expected unique constraint error")
	}
}
```

## Step-by-step execution

For `TestIntegrationCreateUser`:

1. `setupDB` opens an in-memory SQLite database and creates the `users` table.
2. `repo.Create(ctx, "Alice", "alice@example.com")` inserts a row.
3. `LastInsertId` returns the auto-generated ID.
4. `GetByID` retrieves the row.
5. Assert fields match expected values.
6. `t.Cleanup` closes the database when the test completes.

## Common mistakes

- **Sharing a database between tests**: Tests that share a DB have hidden dependencies. Each test should get its own `:memory:` database.
- **Not cleaning up between tests**: File-based databases accumulate state. Use `:memory:` or truncate tables.
- **Unit tests depending on database**: If a test needs a real DB, it's an integration test. Tag it accordingly.
- **Forgetting build tags**: Without `//go:build integration`, all tests run with `go test`. CI should run both: `go test ./...` and `go test -tags=integration ./...`.
- **Testing the ORM/framework instead of your code**: Don't test that SQLite works. Test that your repository methods produce correct SQL and handle results properly.

## Debugging walkthrough

**Scenario**: An integration test passes locally but fails in CI.

```go
func TestIntegrationCreateUser(t *testing.T) {
    db, err := sql.Open("sqlite", "file:test.db") // shared file!
    // ...
}
```

**Root cause**: The test uses a file-based database `test.db`. Local runs use the same file, so state leaks between test runs. CI has a fresh filesystem, so the first run works but subsequent runs fail because the table already exists.

**Fix**: Use `:memory:` for isolated, ephemeral databases:

```go
db, err := sql.Open("sqlite", ":memory:")
```

## Production notes

- **Test containers**: For PostgreSQL integration tests, use `testcontainers-go` to spin up a disposable PostgreSQL container per test suite.
- **Schema migrations in tests**: Apply migrations before running integration tests, not during each test. Use `TestMain` for suite-level setup.
- **Parallel integration tests**: SQLite `:memory:` databases are connection-scoped. Each parallel test gets its own DB — no conflicts.
- **CI pipeline**: Run integration tests on a dedicated step, not mixed with unit tests. Use `go test -tags=integration -count=1 ./...` for repeatable results.
- **Test data factories**: Write helper functions to insert seed data. Avoid hardcoding IDs — use `LastInsertId`.

## Performance implications

- In-memory SQLite tests run in milliseconds — comparable to unit tests.
- PostgreSQL integration tests with test containers take 5-30 seconds to start the container, then tests run in milliseconds.
- Parallel execution (using `t.Parallel()`) is safe with `:memory:` databases but not with file-based or shared databases.
- Table-driven integration tests work well: the setup cost is amortized across all test cases.

## Practice task

1. Create an integration test file `integration_test.go` with `//go:build integration`.
2. Write tests for a `ProductRepository` (Create, GetByID, List, Update, Delete).
3. Each test uses a fresh in-memory SQLite database.
4. Run with: `go test -tags=integration -v ./curriculum/modules/09-sql-postgres-persistence/lessons/25-integration-tests`
5. Also run without the tag and confirm the integration tests are excluded.

## Tests / verification

```bash
go run ./curriculum/modules/09-sql-postgres-persistence/lessons/25-integration-tests
go test -tags=integration ./curriculum/modules/09-sql-postgres-persistence/lessons/25-integration-tests
```

## Review questions

1. How does `:memory:` SQLite differ from a file-based database in testing?
2. What is the purpose of `//go:build integration` build tags?
3. Why should each integration test use its own database instance?
4. How would you structure integration tests for PostgreSQL using test containers?
5. What is the difference between a unit test with a fake repository and an integration test with a real database?

## NEXT UP

sqlc — type-safe SQL query generation from raw SQL.
