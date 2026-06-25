# Why architecture exists

## Learning objective

Identify coupling and cohesion in Go code, explain Conway's law and its effect on software structure, evaluate architectural drivers (team size, change velocity, domain complexity), and apply a tradeoff-first mindset when choosing between design alternatives.

## Why this matters

Every Go project starts with a single `main.go`. Without deliberate architecture, the codebase congeals into a mass of intertwined dependencies that resist change. When the business asks for a new feature, what should be a one-file change becomes a multi-day spelunking expedition. Architecture is what keeps a codebase malleable as it scales from 1K to 100K lines. Professional engineers are judged not by how fast they write code, but by how long the codebase stays productive.

## Mental model

Architecture is a set of constraints you choose to impose on yourself. Think of a city: without zoning laws, anyone can build anything anywhere, but the result is chaos where one building's plumbing runs through another's foundation. Zoning (architecture) restricts what goes where and how things connect. The restrictions are annoying, but they prevent the chaos that makes the city unlivable at scale.

In code, architecture means deciding which packages exist, what each package exports, and how packages depend on each other. The earlier you impose gentle structure, the later you hit the wall where the codebase becomes unproductive.

## Core idea

_Architecture_ is the set of design decisions that are expensive to change. Good architecture maximizes the parts of the system that can change independently. Two primary metrics measure this:

- **Coupling**: how much one module depends on another. Low coupling means changing module A does not require changing module B.
- **Cohesion**: how closely the elements within a module belong together. High cohesion means the module has a single, well-defined responsibility.

Conway's law states that organizations design systems that mirror their communication structure. If two teams must coordinate on every change, the code will have high coupling between their modules. If a team is organized around a business domain, the code will naturally align to that domain.

Architectural drivers are the forces that push you toward one structure vs another:

| Driver | Effect |
|---|---|
| Team size | More teams need more boundaries |
| Change velocity | Fast-changing domains need isolation |
| Domain complexity | Complex domains need rich models |
| Operational constraints | Uptime, latency, compliance |

## Under the hood

Go's compiler enforces a directed acyclic import graph. This single constraint -- no circular imports -- forces developers to think about dependency direction at compile time rather than debugging it at runtime. The `go tool` resolves imports by walking the dependency graph; if it detects a cycle, it reports it as a compile error and stops.

The package-level visibility (`exported` vs `unexported`) is enforced by the compiler. An exported name is accessible to any importer; an unexported name is only visible within its package. This binary visibility system is deliberately simple -- there is no `protected` or `friend` -- and it prevents the subtle coupling patterns found in languages with richer visibility controls.

## How Go uses it

Go's package system is the primary architectural primitive. Production Go projects use packages to:

- **Isolate domains**: `pkg/customer/`, `pkg/order/`, `pkg/payment/`
- **Enforce boundaries**: `internal/` packages cannot be imported by external code
- **Define contracts**: interfaces in consumer packages, not producer packages
- **Control dependency direction**: high-level policy packages import low-level detail packages, never the reverse

The standard library itself is organized this way. `net/http` defines the `Handler` interface. Concrete implementations live in other packages. The `http` package does not import `json`, `xml`, or any particular router.

## Go example

```go
package main

import (
	"errors"
	"fmt"
)

type Order struct {
	ID     string
	Amount float64
	Status string
}

type BadOrderService struct{}

func (BadOrderService) Process(orderID string, amount float64) error {
	if amount <= 0 {
		return errors.New("invalid amount")
	}
	fmt.Printf("BadOrderService: processing order %s for $%.2f\n", orderID, amount)
	return nil
}

type Logger interface {
	Log(message string)
}

type PaymentGateway interface {
	Charge(amount float64) error
}

type OrderValidator interface {
	Validate(amount float64) error
}

type OrderProcessor struct {
	logger    Logger
	payment   PaymentGateway
	validator OrderValidator
}

func NewOrderProcessor(logger Logger, payment PaymentGateway, validator OrderValidator) *OrderProcessor {
	return &OrderProcessor{logger: logger, payment: payment, validator: validator}
}

func (p *OrderProcessor) Process(orderID string, amount float64) error {
	p.logger.Log("processing order " + orderID)
	if err := p.validator.Validate(amount); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}
	if err := p.payment.Charge(amount); err != nil {
		return fmt.Errorf("payment failed: %w", err)
	}
	p.logger.Log("order " + orderID + " processed successfully")
	return nil
}

type ConsoleLogger struct{}
func (ConsoleLogger) Log(message string) { fmt.Println("LOG:", message) }

type StripeGateway struct{}
func (StripeGateway) Charge(amount float64) error { fmt.Printf("Stripe: charging $%.2f\n", amount); return nil }

type AmountValidator struct{}
func (AmountValidator) Validate(amount float64) error {
	if amount <= 0 { return errors.New("amount must be positive") }
	return nil
}

func main() {
	tight := BadOrderService{}
	tight.Process("ORD-001", 49.99)

	loose := NewOrderProcessor(ConsoleLogger{}, StripeGateway{}, AmountValidator{})
	loose.Process("ORD-002", 99.99)
}
```

## Step-by-step execution

For `NewOrderProcessor(...)` -> `Process("ORD-002", 99.99)`:

1. `NewOrderProcessor` receives three concrete dependencies as interface values. Go checks at compile time that each concrete type satisfies its interface; no runtime type assertion is needed.
2. `Process` calls `logger.Log("processing order ORD-002")`. The `ConsoleLogger.Log` method runs, printing the log line.
3. `Process` calls `validator.Validate(99.99)`. `AmountValidator.Validate` checks `99.99 > 0`, which is true, so it returns nil.
4. `Process` calls `payment.Charge(99.99)`. `StripeGateway.Charge` prints the charging message and returns nil.
5. `Process` calls `logger.Log("order ORD-002 processed successfully")`. A final log line is printed.

Compare with `BadOrderService.Process`: it hard-codes validation logic, output formatting, and payment processing into one method. To change the payment provider, you must edit the struct. To change the log destination, you must edit the struct. The `OrderProcessor` version lets you swap any dependency by passing a different implementation at construction time.

## Common mistakes

- **Jumping to microservices**: every new Go project starts with a gRPC service in its own repo, separate deployment, and a Kubernetes manifest. The overhead of service boundaries (network calls, serialization, deployment coordination, observability) overwhelms the team before they have product-market fit. Start with a modular monolith; extract services only when there is a proven scalability or team-boundary need.

- **Architecture astronaut pattern**: implementing hexagonal architecture with seven layers of interfaces before writing a single business function. The abstractions solve problems that do not exist yet and make the code impossible to navigate. Architecture should be refactored toward, not predicted upfront.

- **Letting accidental architecture become permanent**: the directory structure chosen in week 1 (`pkg/`, `internal/`, `cmd/`) becomes the team's mental model. When the service needs to split, the team is reluctant to restructure because that is how it has always been done. Architecture must be treated as malleable.

## Debugging walkthrough

Consider a service where a bug fix requires touching three unrelated packages:

```go
// package payment
func ProcessSale(amount float64) error { /* writes to payment DB */ }

// package email
func SendReceipt(amount float64) error { /* reads order DB */ }

// package order
func CreateOrder(item string) error { /* writes to order DB */ }
```

**Symptom**: A change to `ProcessSale` to add tax calculation breaks `SendReceipt` because both read from the same database table with different assumptions about currency formatting.

**Root cause**: High coupling. `ProcessSale` and `SendReceipt` share a database schema but live in different packages with no coordination. A change in one package has side effects in the other.

**Investigation**: Run `go list -f '{{.Imports}}' ./pkg/...` to see the dependency graph. Use `go mod graph` to visualize module-level dependencies. Look for packages that import many others -- they are coupling hubs.

**Fix**: Extract a shared `billing` package with a clear interface that both `payment` and `email` depend on. The `billing` package owns the database schema, and both `payment` and `email` call `billing` methods. Now changes to the schema are isolated in one place.

## Production notes

- **Coupling budget**: every import has a cost. When reviewing a PR, look at the import block. If a package imports 20+ others, it likely has too many responsibilities.
- **Cohesion signal**: if you cannot describe a package's purpose in one sentence, it is not cohesive.
- **Team alignment**: match package boundaries to team boundaries. If two teams own the same package, Conway's law guarantees coordination pain.
- **Refactoring tolerance**: structure the codebase so that a package can be extracted into a separate service without rewriting everything. Use interfaces at boundary points from day one.
- **Governance**: some teams use "architectural fitness functions" -- automated checks that enforce dependency rules (e.g., `pkg/order/` may not import `pkg/infrastructure/`). Tools like `go vet` with custom analyzers can enforce these rules in CI.

## Performance implications

Architecture decisions have a performance profile that is often invisible until it matters:

- **Package boundary crossing costs**: calling a function through an interface has a small overhead (one indirect function call vs a direct call). In hot loops (millions of iterations), this can add up. Profile before optimizing.
- **Import graph size**: more packages mean more compilation units. Go's compiler is fast, but a project with hundreds of packages will have slower incremental builds than one with dozens. The `go build` cache mitigates this for unchanged packages.
- **Inlining inhibition**: interface calls prevent inlining. If a hot path goes through an interface method, the compiler cannot inline it. For latency-sensitive code, consider concrete types in the inner loop.
- **Over-abstraction cost**: every indirection layer adds allocations and CPU overhead. A service that wraps a repository that wraps a cache that wraps a database has four allocation sites per query. Measure, then abstract.

## Practice task

Take this tightly coupled code:

```go
type ReportGenerator struct{}

func (r ReportGenerator) GenerateCSV() string {
    // reads from database, formats CSV, returns string
    return "col1,col2\n1,2\n"
}
```

Refactor it into three cohesive components with low coupling:

1. A `DataFetcher` interface (abstracts the data source)
2. A `CSVFormatter` (converts data to CSV format)
3. A `ReportService` that composes the two via constructor injection

Write the interfaces, implementations, and a `main()` that uses the refactored design.

## Tests / verification

```bash
go run ./curriculum/modules/14-architecture-distributed-systems/lessons/01-why-architecture-exists
go test ./curriculum/modules/14-architecture-distributed-systems/lessons/01-why-architecture-exists
```

The existing tests verify that `OrderProcessor` validates, logs, and processes payments correctly using mock dependencies. They also verify that `BadOrderService` rejects invalid amounts. After completing the practice task, add tests for your refactored `ReportService`.

## Review questions

1. What is the difference between coupling and cohesion? Give a Go example of high coupling and one of low coupling.
2. How does Conway's law affect the decision to split a monolith into microservices?
3. You have two packages: `pkg/user` (10 exported functions, one responsibility) and `pkg/utils` (50 exported functions, unrelated utilities). Which has higher cohesion? Which is likely more maintainable?
4. A new engineer on your team wants to add a `DeletedAt` timestamp to the User struct in `pkg/order/`. The `pkg/user/` package also reads User data. What architecture concern does this raise?
5. Under what circumstances would you choose a modular monolith over microservices?

## NEXT UP

Package boundaries -- how to design, name, and organize Go packages for maximum maintainability.
