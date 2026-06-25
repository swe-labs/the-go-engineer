# Repository pattern deep dive

## Learning objective

Design repository interfaces that abstract data access, implement in-memory repositories for testing, and swap storage backends without changing business logic.

## Why this matters

Direct database access in business logic creates two problems: tests require a real database, and changing the database requires rewriting scattered queries. The repository pattern collects all data access into an interface. Business logic calls the interface. Implementations handle the specific storage technology. This separation means you can test with an in-memory repository, switch from Postgres to MongoDB by writing a new implementation, and keep your business logic pure.

## Mental model

A repository is a collection-like interface for domain objects. You ask the repository for a `User` by ID, and you get one back. You tell the repository to save a `User`, and it persists. The repository hides whether the data lives in memory, in Postgres, in Redis, or in a CSV file. To the caller, it looks like a map with better error handling.

This abstraction is not about supporting multiple databases in production (you almost never do that). It is about making tests fast and reliable. An in-memory repository has no network, no disk, no connection pooling, no migration scripts. Tests that use it run in milliseconds and never flake due to database unavailability.

## Core idea

A repository interface is defined in the domain package (the consumer), not the infrastructure package (the implementor). This is dependency inversion: high-level modules (domain) do not depend on low-level modules (infrastructure). Both depend on abstractions.

```go
// Package domain (consumer)
type UserRepository interface {
    FindByID(id string) (*User, error)
    Save(user *User) error
    Delete(id string) error
}

// Package memory (implementor)
type UserRepo struct { /* ... */ }
func (r *UserRepo) FindByID(id string) (*User, error) { /* ... */ }
func (r *UserRepo) Save(user *User) error { /* ... */ }
func (r *UserRepo) Delete(id string) error { /* ... */ }
```

A repository should be focused on a single aggregate root. Do not create a single `Repository` interface with methods for every entity. Each aggregate gets its own repository. This keeps interfaces small and implementations focused.

Standard repository operations include:
- `Save`: persists a domain object (insert or update).
- `FindByID`: retrieves a single object by its unique identifier.
- `FindAll`: returns all objects (pagination is often added later).
- `Delete`: removes an object.
- Query methods: `FindByEmail`, `FindByStatus`, etc.

## Under the hood

Go interfaces are implicitly satisfied. The compiler checks interface satisfaction at the point of assignment:

```go
var repo UserRepository = &InMemoryUserRepo{}
```

If `*InMemoryUserRepo` does not have all the methods required by `UserRepository`, the compiler reports an error at this line. There is no `implements` keyword and no explicit registration. This means you can create an in-memory implementation in a test file alongside the interface without importing a separate package.

The `database/sql` package demonstrates this pattern. `sql.DB` is a concrete type that implements `sql.Queryer` and related interfaces. You can swap it with a mock by defining a type that satisfies the same interface.

## How Go uses it

The repository pattern appears in nearly every production Go service:

- Standard library: `io.Reader` and `io.Writer` are repository-like abstractions over data sources and sinks. `os.File`, `bytes.Buffer`, `strings.Reader`, `gzip.Reader` all implement these interfaces.
- ORMs like GORM and sqlx are often hidden behind repository interfaces in production code. The service layer never calls GORM directly; it calls `UserRepository`.
- CockroachDB's SQL layer uses repository-like abstractions for table access. The `sql` package defines interfaces that multiple storage engines implement.

## Go example

```go
package main

import (
	"errors"
	"fmt"
	"sync"
)

type Product struct {
	ID    string
	Name  string
	Price float64
}

type ProductRepository interface {
	Save(product Product) error
	FindByID(id string) (Product, error)
	FindAll() ([]Product, error)
	Delete(id string) error
}

type InMemoryProductRepo struct {
	mu    sync.RWMutex
	store map[string]Product
}

func NewInMemoryProductRepo() *InMemoryProductRepo {
	return &InMemoryProductRepo{store: make(map[string]Product)}
}

func (r *InMemoryProductRepo) Save(product Product) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.store[product.ID] = product
	return nil
}

func (r *InMemoryProductRepo) FindByID(id string) (Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.store[id]
	if !ok {
		return Product{}, errors.New("product not found")
	}
	return p, nil
}

func (r *InMemoryProductRepo) FindAll() ([]Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]Product, 0, len(r.store))
	for _, p := range r.store {
		result = append(result, p)
	}
	return result, nil
}

func (r *InMemoryProductRepo) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.store, id)
	return nil
}

type ProductService struct {
	repo ProductRepository
}

func NewProductService(repo ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) AddProduct(id, name string, price float64) error {
	if price <= 0 {
		return errors.New("price must be positive")
	}
	return s.repo.Save(Product{ID: id, Name: name, Price: price})
}

func (s *ProductService) ListProducts() ([]Product, error) {
	return s.repo.FindAll()
}

func main() {
	repo := NewInMemoryProductRepo()
	svc := NewProductService(repo)

	_ = svc.AddProduct("p1", "Widget", 9.99)
	_ = svc.AddProduct("p2", "Gadget", 24.99)

	products, _ := svc.ListProducts()
	for _, p := range products {
		fmt.Printf("%s: %s ($%.2f)\n", p.ID, p.Name, p.Price)
	}
}
```

## Step-by-step execution

For `svc.AddProduct("p1", "Widget", 9.99)`:

1. `AddProduct` validates that price is positive. 9.99 passes.
2. A `Product` value is constructed.
3. `s.repo.Save(product)` calls `InMemoryProductRepo.Save`.
4. `Save` acquires a write lock, stores the product in the map by ID, and releases the lock.
5. Nil is returned.

For `svc.ListProducts()`:

1. `s.repo.FindAll()` is called.
2. `InMemoryProductRepo.FindAll` acquires a read lock, iterates the map, builds a result slice, and returns it.
3. The service returns the slice to the caller.

The `sync.RWMutex` ensures concurrent access is safe. Reads do not block other readers (only writers), which matters when multiple goroutines query the repository simultaneously.

## Common mistakes

- **Repository returns database types**: the repository should return domain types, not `sql.Row` or `bson.M`. If the repository returns database-specific types, the service layer becomes coupled to the database. Convert inside the repository.
- **One repository per database table**: repositories should be per aggregate root, not per table. An `OrderRepository` handles the `orders` table and the `order_items` table. The caller does not know about `order_items`.
- **Repository has too many methods**: a repository with 30 query methods (`FindByStatusAndDate`, `FindByCustomerAndStatus`, etc.) is a sign the abstraction is wrong. Consider using a query object or specification pattern for complex queries.
- **In-memory repository is too simple**: a real database enforces constraints (unique, not null, foreign keys). If the in-memory repository does not enforce the same constraints, tests pass but production fails. Make the in-memory repository realistic.

## Debugging walkthrough

Consider a service that directly queries the database:

```go
func GetOrderTotal(id string) (float64, error) {
    row := db.QueryRow("SELECT total FROM orders WHERE id = $1", id)
    var total float64
    row.Scan(&total)
    return total, nil
}
```

**Symptom**: Adding a discount calculation requires modifying this function. Testing requires a database with seed data.

**Investigation**: The database query is inlined in the business logic. There is no abstraction. Every call site that needs order data must repeat the query.

**Root cause**: Missing repository abstraction.

**Fix**: Define an `OrderRepository` interface with `FindByID`, implement it with the SQL query, and inject it into the service. The service calls `repo.FindByID(id)` and applies discount logic on the returned domain object. Tests use an in-memory repo.

## Production notes

- **In-memory repos for tests**: every repository interface should have an in-memory implementation used in tests. The in-memory implementation should enforce the same constraints as the real database (unique IDs, required fields, non-nullable fields). This catches bugs before CI.
- **Repository layering**: some teams layer caches between the service and the repository. The service calls a `CachedRepository` that implements the same interface, checks the cache, and falls through to the backing repository on miss. The service does not know the cache exists.
- **Transaction support**: if a repository needs transaction support, add `Begin(ctx) (Tx, error)` and `Tx` interface with repository methods and `Commit/Rollback`. Do not leak `sql.Tx` to the service layer.
- **Error wrapping**: repositories should wrap errors with context: `fmt.Errorf("user %s: %w", id, err)`. This lets the caller distinguish between not-found errors and infrastructure errors.

## Performance implications

- **Interface dispatch**: calling a repository method through an interface has a small overhead. In a typical HTTP handler (1-5 repository calls), this is irrelevant. In batch processing (millions of calls), consider whether the abstraction is on the critical path.
- **Slice vs iterator**: `FindAll()` returning a slice loads all results into memory. For large datasets, return an iterator or use pagination. The interface can define `FindAll(limit, offset int)` or return a `cursor` type.
- **Write locks**: the in-memory repository uses `sync.RWMutex`. In high-concurrency scenarios, lock contention can be a bottleneck. Consider sharding the map by key hash.

## Practice task

Add a `FindByPriceRange(min, max float64) ([]Product, error)` method to the `ProductRepository` interface. Implement it in `InMemoryProductRepo`. Then add a `DiscountedProductService` that wraps `ProductRepository` and applies a percentage discount to products returned by `FindByPriceRange`. Write table-driven tests for both the repository method and the service.

## Tests / verification

```bash
go run ./curriculum/modules/14-architecture-distributed-systems/lessons/04-repository-pattern-deep-dive
go test ./curriculum/modules/14-architecture-distributed-systems/lessons/04-repository-pattern-deep-dive
```

The tests verify CRUD operations on the in-memory repository and validate that the service layer rejects invalid prices. After completing the practice task, add tests for the new price-range query and discounted service.

## Review questions

1. Why should repository interfaces be defined in the domain package rather than the infrastructure package?
2. What is the difference between a repository per table and a repository per aggregate root?
3. How does an in-memory repository make tests faster and more reliable?
4. What happens when you add a sixth method to a repository interface that has five implementations?
5. How would you add caching to a repository without changing the service layer or the repository interface?

## NEXT UP

Modular monolith -- organizing code into bounded contexts within a single process for clarity without the overhead of services.
