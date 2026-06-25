# Package boundaries

## Learning objective

Design Go packages with clear boundaries, use internal packages to hide implementation details, enforce dependency direction, avoid cyclic dependencies, and organize packages by domain rather than by layer.

## Why this matters

The Go compiler enforces exactly two visibility levels within a module: exported (accessible to all importers) and unexported (accessible only within the package). There is no `protected`, `friend`, or subpackage visibility. This simplicity means package boundaries are the _only_ architectural boundary you get for free. Every exported name is a commitment. Every import creates a dependency. Getting package boundaries right is the single highest-leverage architectural skill in Go.

## Mental model

A package is a walled garden. Inside the walls, everything is shared -- functions, types, variables -- whether exported or not. Outside the walls, only the officially gated exports are visible. The wall has a cost: every time you cross it, you use the gate (the exported API). The benefit is that the garden inside can be replanted, restructured, or replaced as long as the gates stay the same.

A good package boundary is one where the cost of crossing is justified by the freedom it gives the code inside. A bad boundary is one where you constantly need to cross it (too many packages) or where you never benefit from the isolation (too few packages).

## Core idea

A package boundary serves three purposes:

1. **Compilation unit**: Go compiles packages independently and caches the results. Changing one package only recompiles its dependents.
2. **Visibility boundary**: Exported names form the contract. Unexported names can change without notice.
3. **Dependency node**: Packages form a directed acyclic graph. The compiler rejects cycles at build time.

The `internal/` package convention is a module-level visibility boundary. Packages rooted at `internal/` can only be imported by code that shares a common ancestor directory. This lets you create multi-package subsystems where some packages are internal to the subsystem but exported within the module.

Dependency direction should flow from high-level policy (what the application should do) toward low-level detail (how it does it). Infrastructure packages should be leaves in the import graph, never roots.

Cyclic dependencies are a signal that the boundary is drawn incorrectly. The fix is either to merge the two packages or to extract shared types into a third package that both depend on.

Domain-driven package layout organizes packages by business concept (`pkg/customer/`, `pkg/order/`, `pkg/inventory/`) rather than by technical layer (`pkg/models/`, `pkg/services/`, `pkg/repositories/`). Layer-organized packages create high coupling because every model is in one place and every service is in another. Domain-organized packages keep related types together and reduce the blast radius of changes.

## Under the hood

Go's import resolution works as follows:

1. The compiler resolves each import path to a directory on disk.
2. It reads all `.go` files in that directory that share the same package declaration.
3. It builds a symbol table of exported names.
4. Import cycles are detected during dependency graph construction: the compiler walks the import graph depth-first, marking each package as visiting. If it encounters a package already marked as visiting, it reports an import cycle and stops. This is O(V+E) where V is packages and E is imports.
5. The `internal/` convention is enforced by `go/build`. During import resolution, the tool checks whether the importing package's directory tree contains the `internal/` directory as an ancestor. If not, the import is rejected.

## How Go uses it

The Go standard library is the best example of package boundary design. `net/http` defines `Handler` and `ResponseWriter` as interfaces. It does not import `json`, `xml`, or `template` -- those are separate packages. `io` defines `Reader` and `Writer`. `os.File` implements them. `bytes.Buffer` implements them. The interfaces are in the consuming package, enabling any implementation to be plugged in.

Production Go projects follow the same pattern:
- `github.com/kubernetes/kubernetes/pkg/api/` owns the API types.
- `github.com/kubernetes/kubernetes/pkg/kubelet/` owns the kubelet logic.
- `github.com/kubernetes/kubernetes/pkg/scheduler/` owns the scheduler.
Each is a self-contained package with a small exported surface.

## Go example

```go
package main

import (
	"fmt"
)

type ExportedService struct {
	repo Repository
}

func NewExportedService(repo Repository) *ExportedService {
	return &ExportedService{repo: repo}
}

type Repository interface {
	Save(item string) error
	Find(id string) (string, error)
}

type inMemoryStore struct {
	data map[string]string
}

func newInMemoryStore() *inMemoryStore {
	return &inMemoryStore{data: make(map[string]string)}
}

func (s *inMemoryStore) Save(item string) error {
	s.data[item] = item
	return nil
}

func (s *inMemoryStore) Find(id string) (string, error) {
	val, ok := s.data[id]
	if !ok {
		return "", fmt.Errorf("not found: %s", id)
	}
	return val, nil
}

func main() {
	repo := newInMemoryStore()
	svc := NewExportedService(repo)

	_ = svc.repo.Save("user-1")
	result, _ := svc.repo.Find("user-1")
	fmt.Println("Found:", result)

	_, err := svc.repo.Find("user-2")
	fmt.Println("Error:", err)
}
```

This demonstrates the key package-boundary concepts even within a single package. `Repository` is an exported interface -- the public contract. `inMemoryStore` is an unexported implementation -- it can be replaced entirely without affecting consumers. `NewExportedService` is the constructor that accepts the interface. The main function only interacts with the exported surface.

## Step-by-step execution

1. `newInMemoryStore()` creates an `*inMemoryStore` with an empty map. This function is unexported -- only code within the same package can call it.
2. `NewExportedService(repo)` stores the `Repository` interface value. The compiler checks that `*inMemoryStore` satisfies `Repository` at the assignment site (it does, because `Save` and `Find` match).
3. `svc.repo.Save("user-1")` calls through the interface. The runtime dispatches to `inMemoryStore.Save`.
4. `svc.repo.Find("user-1")` retrieves the value. The method performs a map lookup.
5. `svc.repo.Find("user-2")` fails because the key does not exist.

The critical design property: callers do not know about `inMemoryStore`. If you switch to a database-backed store, you write a new type that satisfies `Repository` and pass it to `NewExportedService`. No code in `ExportedService` changes.

## Common mistakes

- **Putting every type in one package**: a single `models/` or `types/` package with 200 exported types has no internal cohesion. Every file imports `models`, creating coupling across the entire codebase. Any change to any type requires recompiling every dependent package.
- **Making packages too small**: a package with one 5-line type and one 3-line function adds import overhead without any benefit. The cost of understanding a package boundary must be justified by the cohesion within.
- **Circular imports between related packages**: `users.go` imports `orders.go` and `orders.go` imports `users.go`. The Go compiler rejects this, forcing you to either merge the packages or extract shared types into a third package. The circular import is a signal that the boundary is wrong.
- **Exporting everything**: every exported name is a commitment. Unexport by default. Only export what external consumers need.

## Debugging walkthrough

Consider three packages with a circular dependency:

```
pkg/user/     imports pkg/order/
pkg/order/    imports pkg/user/
```

**Symptom**: `go build` fails with `import cycle not allowed`.

**Investigation**: Run `go list -f '{{.ImportPath}} imports: {{.Imports}}' ./pkg/...` to see the full dependency graph. Identify the cycle.

**Root cause**: Both packages need a shared type (e.g., `OrderSummary` containing user info). They could not decide where to put it, so each tried to define it and import from the other.

**Fix**: Create a `pkg/shared/` or better, a `pkg/order/summary.go` that defines only the types needed by `pkg/user/`. Or merge the two packages into a `pkg/sales/` domain package. The domain-driven approach (merge) is preferred because it also fixes the cohesion problem.

## Production notes

- **`internal/` is your friend**: put implementation details in `internal/`. This prevents external consumers from depending on them and gives you freedom to refactor.
- **Package naming**: the package name should describe what it provides, not what it contains. `package user` is good. `package models` is not. Avoid generic names like `utils`, `helpers`, `common`.
- **Minimal exported surface**: start with unexported types and only export when a consumer outside the package needs them. Exporting a type is easy; unexporting it can break consumers.
- **Dependency inversion**: define interfaces in the consumer package, not the implementor. `package user` defines `UserRepository`. `package postgres` implements it. This prevents `package user` from importing `package postgres`.
- **Test packages**: use `_test` suffix packages to test your exported API from an external perspective. This catches accidental coupling to unexported details.

## Performance implications

- **Package count**: the Go compiler compiles each package independently. More packages mean more compilation overhead. However, the build cache means only changed packages and their dependents are recompiled.
- **Interface call overhead**: crossing a package boundary through an interface method costs a dynamic dispatch. In non-hot paths, this is immeasurable. In hot loops, consider whether the abstraction is necessary.
- **Inlining**: the compiler can inline across package boundaries for direct calls. Interface calls prevent inlining.

## Practice task

Split the following code into three packages in separate directories under a parent module:

1. `store/` -- an exported `Store` interface with `Get(key string) (string, error)` and `Set(key, value string) error`.
2. `memory/` -- an unexported `memoryStore` implementation of `Store`, plus a `New()` constructor.
3. `main/` imports `store/` and `memory/`, creates a store, and uses it.

Then verify that the `memory/` package cannot be imported by a hypothetical external module (simulate by trying to import it from outside the module). Observe that `memory` is importable within the module but not outside.

## Tests / verification

```bash
go run ./curriculum/modules/14-architecture-distributed-systems/lessons/02-package-boundaries
go test ./curriculum/modules/14-architecture-distributed-systems/lessons/02-package-boundaries
```

The tests verify that the in-memory store saves and finds items, returns an error for missing keys, and that the exported service correctly uses a repository via the interface. After completing the practice task, add tests for your `Store` interface implementations.

## Review questions

1. What happens if you create a circular import between two packages? How does Go detect it?
2. What is the difference between putting a type in `internal/` vs leaving it unexported within a package?
3. Given `pkg/user/` (depends on `pkg/db/`) and `pkg/order/` (depends on `pkg/db/`), what change would you make if `pkg/db/` needs to call a function in `pkg/user/`?
4. Why does Go not have `protected` or `friend` visibility? How does the `internal/` convention fill this gap?
5. Name two signals that a package boundary is drawn incorrectly.

## NEXT UP

Service layer -- isolating business logic behind a dedicated layer with dependency injection via constructor.
