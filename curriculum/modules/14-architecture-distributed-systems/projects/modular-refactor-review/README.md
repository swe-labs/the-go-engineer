# Modular Refactor Review

## Learning objective

Refactor a tightly coupled, monolithic Go application into a modular, testable hexagonal (ports and adapters) architecture. Master domain-driven separation of concerns, dependency inversion via interface abstraction, constructor-based dependency injection, and characterization test verification without mutating global state.

## Why this matters

As codebases grow, initial prototypes frequently degenerate into sprawling monoliths where HTTP handlers directly perform database operations, business calculations, third-party payment requests, and notification dispatches. This creates untestable codebases prone to regressions and cascade failures. Senior Go engineers must know how to decouple business rules from infrastructure mechanisms, allowing business logic to be tested in isolation and infrastructure adapters to be swapped seamlessly.

## Prerequisites

Complete Module 14 lessons on modular design, hexagonal architecture, ports and adapters pattern, interface segregation, dependency injection, and clean boundaries.

## Available concepts

This project draws on domain modeling, Go interfaces (`ports`), dependency inversion principle, constructor injection, adapter implementations (in-memory repositories, mock payment gateways, mock notification services), standard library `net/http` routing, and unit/integration testing with `httptest`.

## Project overview

The `_starter/` directory contains an untestable legacy monolith (`main.go`). It handles e-commerce checkout, but exhibits severe architectural flaws:
- Global mutable database map (`globalDB`) guarded by a global lock.
- Business rules (discount calculation, order total validation) tangled inside the HTTP handler.
- Concrete third-party payment calls and email logging hardcoded directly into the request flow.
- Zero unit testability without mutating global variables or spawning external servers.

Your mission is to refactor this monolith into clean layered packages following Hexagonal Architecture:
1. `domain.go`: Pure domain entities (`Order`, `OrderItem`, `OrderStatus`) and business invariants (`CalculateTotals`).
2. `ports.go`: Inbound/outbound interfaces (`OrderRepository`, `PaymentGateway`, `NotificationService`).
3. `service.go`: Application use case orchestration (`OrderService`) with constructor dependency injection.
4. `adapters.go`: Swappable infrastructure adapters (`InMemoryOrderRepository`, `MockPaymentGateway`, `MockNotificationService`).
5. `handler.go`: HTTP REST controller converting HTTP requests/responses into domain service calls.
6. `main.go`: Clean composition root assembling the dependency graph.

## Architectural Blueprint

```
                      ┌───────────────────────────────────────┐
                      │             HTTP Handler              │ (Inbound Adapter)
                      └──────────────────┬────────────────────┘
                                         │ Calls
                                         ▼
                      ┌───────────────────────────────────────┐
                      │             OrderService              │ (Application Core)
                      │   - Orchestrates checkout flow        │
                      │   - Enforces business transactions    │
                      └───────┬──────────┬──────────┬─────────┘
                              │ Uses     │ Uses     │ Uses
                              ▼          ▼          ▼
                      ┌──────────────┐ ┌─────────┐ ┌───────────────┐
                      │  OrderRepo   │ │ Payment │ │ Notification  │ (Ports / Interfaces)
                      │  Interface   │ │ Gateway │ │    Service    │
                      └──────▲───────┘ └────▲────┘ └───────▲───────┘
                             │ Implements   │ Implements   │ Implements
                      ┌──────┴───────┐ ┌────┴────┐ ┌───────┴───────┐
                      │  In-Memory   │ │  Mock   │ │ Mock Email/   │ (Outbound Adapters)
                      │  Repository  │ │ Payment │ │   Notifier    │
                      └──────────────┘ └─────────┘ └───────────────┘
```

## Tasks

### Task 1 — Define Domain Models and Business Invariants

In `domain.go`:
- Define `OrderItem` and `Order` structs with clean JSON annotations.
- Define `OrderStatus` (`PENDING`, `PAID`, `CANCELLED`).
- Implement `(o *Order) CalculateTotals() error`:
  - Validate customer name, email, item quantities, and item prices.
  - Compute subtotal across all items.
  - Apply the 10% discount rule if the subtotal is $100 (10,000 cents) or greater.
  - Set `Total = Subtotal - Discount`.

### Task 2 — Specify Port Interfaces

In `ports.go`:
- Define `OrderRepository` interface: `Save(ctx, Order) error`, `FindByID(ctx, string) (Order, error)`, `List(ctx) ([]Order, error)`.
- Define `PaymentGateway` interface: `Charge(ctx context.Context, amountCents int, paymentToken string) (string, error)`.
- Define `NotificationService` interface: `SendOrderConfirmation(ctx context.Context, order Order) error`.

### Task 3 — Implement Use Case Service

In `service.go`:
- Define `OrderService` holding private fields for the three port interfaces.
- Implement constructor `NewOrderService(repo OrderRepository, payment PaymentGateway, notifier NotificationService) *OrderService`.
- Implement `Checkout(ctx context.Context, dto CreateOrderDTO) (Order, error)`:
  - Construct order and call `order.CalculateTotals()`.
  - Execute payment via `s.payment.Charge(...)`. If it fails, abort and do not persist.
  - Update status to `StatusPaid` and assign payment reference.
  - Persist order via `s.repo.Save(ctx, order)`.
  - Dispatch notification via `s.notifier.SendOrderConfirmation(ctx, order)`.
  - Return created order.

### Task 4 — Implement Adapters & Composition Root

- In `adapters.go`: Implement `InMemoryOrderRepository`, `MockPaymentGateway`, and `MockNotificationService`.
- In `handler.go`: Implement `HTTPHandler` with endpoints `POST /orders`, `GET /orders/{id}`, `GET /orders`.
- In `main.go`: Wire all dependencies together in `main()` and launch the HTTP server with graceful shutdown.

## Verification

Run the automated test suite:

```bash
go test -v ./curriculum/modules/14-architecture-distributed-systems/projects/modular-refactor-review/_solution
```

Expected output:
```
=== RUN   TestOrderService_SuccessfulCheckoutAndDiscount
--- PASS: TestOrderService_SuccessfulCheckoutAndDiscount (0.00s)
=== RUN   TestOrderService_PaymentFailureRollback
--- PASS: TestOrderService_PaymentFailureRollback (0.00s)
=== RUN   TestHTTPHandler_EndToEnd
--- PASS: TestHTTPHandler_EndToEnd (0.00s)
PASS
ok  	github.com/rasel9t6/the-go-engineer/curriculum/modules/14-architecture-distributed-systems/projects/modular-refactor-review/_solution	0.818s
```

## Rubric

| Criterion | Excellent (5) | Good (3) | Needs Work (0) |
|---|---|---|---|
| **Architectural Separation** | Clean separation across domain, ports, use cases, adapters, and composition root; zero cross-layer bleeding | Extracted interfaces but domain logic still relies on HTTP or DB types | Still relies on global variables or monolithic handlers |
| **Dependency Injection** | All external dependencies injected via constructors; 100% testable without global mutations | Injected some dependencies but left some hardcoded globals | Direct instantiation of concrete infrastructure |
| **Domain Logic Integrity** | Invariants, pricing, and discount rules fully encapsulated and verified by unit tests | Invariants mostly work but discount calculation lives in HTTP layer | Business rules missing or broken |
| **Test Quality** | Unit and integration tests verify happy paths, rollback on payment failures, and HTTP status codes | Tests exist but only cover happy path | No automated tests |

## NEXT UP

[Dockerized Service with CI/CD](../../15-docker-cicd-deployment/projects/dockerized-service-with-ci-cd/README.md)
