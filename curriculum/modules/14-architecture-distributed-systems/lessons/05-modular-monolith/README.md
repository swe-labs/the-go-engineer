# Modular monolith

## Learning objective

Design a modular monolith with bounded contexts communicating through interfaces, evolve the architecture from monolith to services when needed, and recognize when a monolith is the right choice.

## Why this matters

Every distributed system tutorial starts with microservices. But for 90% of projects, a well-structured monolith outperforms microservices in velocity, reliability, and operational cost. The problem is not the monolith itself -- it is the _unstructured_ monolith where every package imports every other package. A modular monolith applies the same bounded-context discipline as microservices but keeps everything in one process. When you later need to extract a service, the interfaces are already in place.

## Mental model

A modular monolith is a city with distinct neighborhoods. The residential zone (customer management) does not import from the industrial zone (inventory management). Each neighborhood has a clear border (the `package` boundary) and communicates with other neighborhoods through defined channels (interfaces). The city government (the main function) wires the neighborhoods together at startup.

If a neighborhood grows too large, it can be split off into its own city (microservice) without redesigning the communication channels. The interfaces between neighborhoods become the service contracts.

## Core idea

A modular monolith has three properties:

1. **Single process**: everything compiles into one binary and runs in one OS process.
2. **Bounded contexts**: the codebase is divided into domains that do not share internal types. Ordering does not import from Billing's internal packages.
3. **Interface-based communication**: bounded contexts communicate through interfaces defined in shared packages or in the consuming context.

The bounded context is the key concept from Domain-Driven Design. Each context has its own ubiquitous language, its own data model, and its own business rules. A `User` in the ordering context may have different fields than a `User` in the marketing context. They are separate types, possibly with different names or different fields, living in different packages.

Internal communication between bounded contexts happens through interfaces. The ordering context defines an `InvoiceService` interface. The billing context implements it. The main function wires them together:

```go
func main() {
    billing := billing.NewService(billing.NewRepo(db))
    ordering := ordering.NewService(ordering.NewRepo(db), billing)
}
```

This monolith-first approach means you get the discipline of service boundaries without the overhead of network calls, serialization, deployment coordination, and observability. When team boundaries or scalability requirements demand a split, the interfaces already define the contract.

## Under the hood

The modular monolith leverages two Go features:

- **Package visibility**: each bounded context is a Go package (or a tree of packages under `internal/`). Unexported types cannot leak across context boundaries.
- **Interface satisfaction**: contexts communicate through interfaces. The compiler checks satisfaction at the wiring point in `main`. If the billing context changes its API, the ordering context does not compile until the interface is satisfied.

Because everything is in one process, communication is synchronous function calls. There is no serialization, no network latency, no partial failure. This makes the modular monolith significantly simpler than a distributed system.

A modular monolith does not prevent you from scaling: you can run multiple instances behind a load balancer. The database is the bottleneck, not the application process.

## How Go uses it

Many successful Go projects started as modular monoliths:

- **GitHub's monolith**: the monolith was organized into service-like components internally. Each component had a clear API. GitHub extracted services incrementally over years, not months.
- **Segment's monolith**: Segment's Go backend was a modular monolith where each integration was a package implementing a common interface. New integrations did not require modifying the core.
- **Basecamp's monolith**: Basecamp famously runs a monolith serving millions of users. Their codebase is organized into modules with clear boundaries.

Go's standard library is itself a modular monolith. `net/http` is a bounded context that does not import `json` or `template`. Communication happens through interfaces (`Handler`, `ResponseWriter`).

## Go example

```go
package main

import (
	"errors"
	"fmt"
	"time"
)

type Order struct {
	ID        string
	Total     float64
	Status    string
	CreatedAt time.Time
}

type Invoice struct {
	OrderID string
	Amount  float64
	Paid    bool
}

type OrderStore interface {
	Save(order Order) error
	FindByID(id string) (Order, error)
}

type InvoiceStore interface {
	Save(invoice Invoice) error
	FindByOrderID(orderID string) (Invoice, error)
}

// BillingContext - handles invoicing and payments.
type BillingService struct {
	invoices InvoiceStore
}

func NewBillingService(invoices InvoiceStore) *BillingService {
	return &BillingService{invoices: invoices}
}

func (s *BillingService) CreateInvoice(order Order) (Invoice, error) {
	inv := Invoice{OrderID: order.ID, Amount: order.Total, Paid: false}
	if err := s.invoices.Save(inv); err != nil {
		return Invoice{}, err
	}
	return inv, nil
}

// OrderingContext - handles order placement and orchestration.
type OrderingBoundedContext struct {
	orders  OrderStore
	billing *BillingService
}

func NewOrderingBoundedContext(orders OrderStore, billing *BillingService) *OrderingBoundedContext {
	return &OrderingBoundedContext{orders: orders, billing: billing}
}

func (ctx *OrderingBoundedContext) PlaceOrder(id string, total float64) (Order, error) {
	if total <= 0 {
		return Order{}, errors.New("total must be positive")
	}
	order := Order{ID: id, Total: total, Status: "placed", CreatedAt: time.Now()}
	if err := ctx.orders.Save(order); err != nil {
		return Order{}, err
	}
	inv, err := ctx.billing.CreateInvoice(order)
	if err != nil {
		return Order{}, fmt.Errorf("billing failed: %w", err)
	}
	fmt.Printf("Invoice %s created for order %s ($%.2f)\n", inv.OrderID, order.ID, inv.Amount)
	return order, nil
}

func main() {
	ctx := NewOrderingBoundedContext(newInMemoryOrderStore(), NewBillingService(newInMemoryInvoiceStore()))
	order, _ := ctx.PlaceOrder("ORD-001", 129.99)
	fmt.Printf("Order placed: %s (status: %s)\n", order.ID, order.Status)
}
```

## Step-by-step execution

For `ctx.PlaceOrder("ORD-001", 129.99)`:

1. `PlaceOrder` validates the total is positive. 129.99 passes.
2. An `Order` struct is constructed with a "placed" status.
3. The order is saved to `OrderStore` (in-memory map).
4. `ctx.billing.CreateInvoice(order)` is called. This crosses from the ordering context into the billing context.
5. `BillingService.CreateInvoice` constructs an `Invoice` and saves it to `InvoiceStore`.
6. The invoice is returned to the ordering context, which prints a confirmation.
7. The order is returned.

The key property: the ordering context does not know about `InvoiceStore` or how invoicing works. It only knows `BillingService.CreateInvoice`. If the invoicing logic changes inside the billing context, the ordering context is unaffected.

## Common mistakes

- **Sharing database tables between contexts**: the ordering context writes to `orders` table and the billing context reads from it directly. This couples the two contexts at the database level. Each context should own its data. Communication goes through interfaces, not shared tables.
- **Sharing domain types between contexts**: both contexts use the same `User` struct from a shared `models` package. When ordering needs a different User field than billing, the shared struct must change for both. Each context should define its own types for the concepts it needs.
- **Synchronous orchestration where async would be better**: ordering calls billing synchronously to create an invoice. If billing is slow, ordering is slow. Consider publishing an event and having billing consume it asynchronously.
- **Leaking internal packages**: putting `internal/` packages in a shared location where other contexts import them. Each context should have its own `internal/` directory.

## Debugging walkthrough

Consider a modular monolith where the ordering context imports directly from billing's internal package:

```go
// ordering/service.go
import "project/billing/internal/tax"
```

**Symptom**: A change in billing's tax calculation breaks ordering's build.

**Investigation**: Run `go list -f '{{.Imports}}' ./...` and look for imports of `internal/` packages across context boundaries.

**Root cause**: The ordering context bypassed the `BillingService` interface and imported an internal implementation detail.

**Fix**: Move the tax calculation behind a `BillingService.TaxFor(order) float64` method. The ordering context calls the interface method. The billing context controls the tax logic. The `internal/tax` package stays internal.

## Production notes

- **Monolith-first**: start as a monolith. Split only when there is a proven need (team boundaries, scaling bottlenecks, deployment frequency conflicts). Premature splitting adds overhead without benefit.
- **Wiring in main**: all dependency injection happens in `main()` or a dedicated `wire.go`. Contexts do not import each other's constructors. This makes the dependency graph visible in one place.
- **Testing across contexts**: test each context's service layer in isolation with mock dependencies of other contexts. Integration tests wire all contexts together with in-memory stores.
- **Extraction readiness**: when you split a context into a service, the interface becomes the RPC contract (gRPC or HTTP). The interface methods should already accept and return serializable types. The adapter code (HTTP client/server) is the only new code.

## Performance implications

- **In-process call overhead**: calling between contexts is a function call with no serialization. This is 100-1000x faster than an RPC call (microseconds vs milliseconds).
- **No network overhead**: no connection pooling, no retries, no timeouts for inter-context calls. This simplifies both code and operations.
- **Single-process memory**: all contexts share the same heap. A memory leak in one context affects the entire process. Use `pprof` to identify allocation sources.
- **Startup time**: one binary is faster to start than N services. This matters for rapid deployment cycles and developer feedback loops.

## Practice task

Add a `ShippingContext` as a third bounded context. Define a `ShippingService` interface in a shared location, implement it in the shipping context, and wire it into the ordering context so that `PlaceOrder` also creates a shipping label after the invoice. Write tests for each context in isolation using mock dependencies.

## Tests / verification

```bash
go run ./curriculum/modules/14-architecture-distributed-systems/lessons/05-modular-monolith
go test ./curriculum/modules/14-architecture-distributed-systems/lessons/05-modular-monolith
```

The existing tests verify that the ordering context places orders and creates invoices through the billing context. Tests use in-memory stores for isolation. After completing the practice task, add tests for the shipping context.

## Review questions

1. What is the difference between a modular monolith and a traditional monolith?
2. How do bounded contexts communicate in a modular monolith? What Go features enable this?
3. Under what circumstances would you extract a bounded context into a separate microservice?
4. Why should each bounded context own its data rather than sharing a database?
5. What is the advantage of a monolith-first approach over starting with microservices?

## NEXT UP

Hexagonal architecture -- ports and adapters for building testable, infrastructure-independent core domain logic.
