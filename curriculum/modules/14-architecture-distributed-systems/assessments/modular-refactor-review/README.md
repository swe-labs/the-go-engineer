# Project Assessment — Modular Refactor Review

## Purpose

Validate mastery for the target curriculum items.

## Targets

- `project-modular-refactor-review`

## Evidence

Submit answers, code evidence, test output, and explanation where required.

---

## 1. Assessment Overview & Industrial Context

In fast-growing software organizations, greenfield prototypes frequently evolve into monolithic, tightly coupled systems where HTTP handlers directly access global database connections, perform complex domain calculations, make synchronous third-party API calls, and emit notifications inline. 

Such architectures create significant technical debt:
- **Global Mutable State**: Shared global maps or database connections make isolated unit testing impossible without cross-test race conditions and state leakage.
- **Untestable Business Logic**: Inability to test domain invariants (e.g. discount rules, validation) without spinning up full HTTP servers and mocking external network endpoints.
- **Cascading Failures**: Infrastructure errors (payment provider timeouts) directly crash or hang the business application layer.

In this assessment, you will refactor a legacy e-commerce checkout monolith (`main.go`) into a clean, modular Hexagonal (Ports and Adapters) Architecture. You will establish domain-driven boundaries, invert infrastructure dependencies through Go interfaces, implement constructor-based dependency injection, and prove architectural decoupling using isolated unit tests with mock adapters.

---

## 2. Core Competencies Tested

- **Hexagonal Architecture (Ports & Adapters)**: Separating pure business domain models from inbound drivers (HTTP handlers) and outbound secondary adapters (databases, payment processors, notification services).
- **Dependency Inversion Principle (DIP)**: Defining port interfaces that the core domain owns, forcing concrete infrastructure adapters to depend on domain interfaces rather than vice versa.
- **Constructor Injection & Composition Root**: Assembling the application's dependency graph explicitly in `main.go` without global singletons or service locators.
- **Isolated Unit Testing with Mocks**: Verifying complex business logic using lightweight in-memory and mock adapters without network calls or database dependencies.
- **Domain Invariant Enforcement**: Encapsulating validation, subtotal calculation, and business rules (e.g., 10% discount for orders >= $100) within pure domain entities.

---

## 3. Zero-Magic Architectural Breakdown

```
                       ┌───────────────────────────────────────┐
                       │             HTTP Handler              │ (Inbound Driver)
                       │         (handler.go / REST)           │
                       └──────────────────┬────────────────────┘
                                          │ Invokes
                                          ▼
                       ┌───────────────────────────────────────┐
                       │             OrderService              │ (Application Core)
                       │           (service.go)                │
                       │   - Coordinates checkout flow         │
                       │   - Enforces transactional steps      │
                       └───────┬──────────┬──────────┬─────────┘
                               │ Uses     │ Uses     │ Uses
                               ▼          ▼          ▼
                       ┌──────────────┐ ┌─────────┐ ┌───────────────┐
                       │  OrderRepo   │ │ Payment │ │ Notification  │ (Ports / Interfaces)
                       │  (ports.go)  │ │ Gateway │ │    Service    │
                       └──────▲───────┘ └────▲────┘ └───────▲───────┘
                              │ Implements   │ Implements   │ Implements
                       ┌──────┴───────┐ ┌────┴────┐ ┌───────┴───────┐
                       │  In-Memory   │ │  Mock   │ │ Mock Email/   │ (Outbound Adapters)
                       │  Repository  │ │ Payment │ │   Notifier    │ (adapters.go)
                       └──────────────┘ └─────────┘ └───────────────┘
```

### Why Go Interfaces Belongs at the Point of Consumption
In Go, interfaces are satisfied implicitly without an `implements` keyword. Idiomatic Hexagonal Architecture defines ports (interfaces) in the consumer package (`ports.go` or alongside the domain service). The domain service specifies *what it requires* to fulfill the use case; concrete adapters in outer layers implement those contracts.

### The Problem with Global Singletons
In legacy monoliths, global variables (`var globalDB = ...`) bypass function parameter signatures, hiding dependencies from callers and preventing parallel test execution (`t.Parallel()`). Explicit constructor injection (`NewOrderService(repo, payment, notifier)`) makes dependencies transparent and allows substituting test doubles effortlessly.

---

## 4. Implementation Requirements & Task Breakdown

### Task 1: Domain Entities & Business Invariants
In `domain.go`:
- Define `OrderItem` (`SKU`, `Quantity`, `PriceCents`).
- Define `OrderStatus` constants (`StatusPending`, `StatusPaid`, `StatusCancelled`).
- Define `Order` struct with standard JSON annotations.
- Implement `(o *Order) CalculateTotals() error`:
  - Validate customer name, email, and at least one item present (`ErrInvalidOrder`).
  - Validate item SKU non-empty, quantity > 0, price > 0 (`ErrInvalidItem`).
  - Calculate subtotal: $\sum (\text{PriceCents} \times \text{Quantity})$.
  - Apply 10% discount if subtotal >= 10,000 cents ($100).
  - Set `Total = Subtotal - Discount`.

### Task 2: Port Interfaces
In `ports.go`:
- `OrderRepository`:
  - `Save(ctx context.Context, order Order) error`
  - `FindByID(ctx context.Context, id string) (Order, error)`
  - `List(ctx context.Context) ([]Order, error)`
- `PaymentGateway`:
  - `Charge(ctx context.Context, amountCents int, paymentToken string) (string, error)`
- `NotificationService`:
  - `SendOrderConfirmation(ctx context.Context, order Order) error`

### Task 3: Application Service Layer
In `service.go`:
- Define `CreateOrderDTO` struct.
- Define `OrderService` struct storing `repo`, `payment`, `notifier`.
- Implement constructor `NewOrderService(...)`.
- Implement `Checkout(ctx context.Context, dto CreateOrderDTO) (Order, error)`:
  1. Instantiate new `Order` with `StatusPending` and generated ID.
  2. Invoke `order.CalculateTotals()`.
  3. Charge customer via `payment.Charge(ctx, order.Total, dto.PaymentToken)`.
  4. On payment error, abort without saving to repo and return wrapped error.
  5. Update status to `StatusPaid` and assign `PaymentRef`.
  6. Persist order via `repo.Save(ctx, order)`.
  7. Dispatch notification via `notifier.SendOrderConfirmation(ctx, order)`.

### Task 4: Infrastructure Adapters & HTTP Handlers
In `adapters.go`:
- Implement `InMemoryOrderRepository` with `sync.RWMutex` thread safety.
- Implement `MockPaymentGateway` recording calls and rejecting token `"fail_token"`.
- Implement `MockNotificationService` recording sent recipient emails.
In `handler.go`:
- Implement `NewHTTPHandler(service *OrderService) *HTTPHandler`.
- Provide `Routes() http.Handler` mounting `POST /orders`, `GET /orders/{id}`, `GET /orders`.

---

## 5. Verification & Acceptance Criteria

Run all unit tests in the project directory:

```bash
go test -v -race ./curriculum/modules/14-architecture-distributed-systems/projects/modular-refactor-review/_solution
```

### Acceptance Checklist
- [ ] `TestOrderService_SuccessfulCheckoutAndDiscount`: verifies full checkout pipeline, 10% discount calculation, payment recording, notification dispatch, and persistence.
- [ ] `TestOrderService_PaymentFailureRollback`: verifies that failed payments abort without saving orders or sending notifications.
- [ ] `TestHTTPHandler_EndToEnd`: validates HTTP request parsing, response status codes (201 Created, 200 OK, 400 Bad Request), and JSON serialization.
- [ ] Zero reliance on global mutable state or package-level locks.
- [ ] All tests pass cleanly under `-race`.

---

## 6. Common Pitfalls & Anti-Patterns

1. **Leaking Infrastructure Details into Domain Models**:
   *Anti-pattern*: Adding SQL tags or HTTP headers to domain structs in `domain.go`.
   *Correction*: Keep domain entities pure; map HTTP DTOs and database models at the boundary layer.

2. **Inverted Dependency Direction**:
   *Anti-pattern*: Importing concrete adapters (`adapters.go`) into the domain service (`service.go`).
   *Correction*: The service must only reference interfaces defined in `ports.go`.

3. **Incomplete Transaction Rollback on Failure**:
   *Anti-pattern*: Saving an order to the repository *before* verifying payment authorization.
   *Correction*: Always verify external preconditions (charge authorization) before committing state to persistence.

---

## 7. 4-Tier Assessment Rubric

| Level | Criteria & Observable Evidence |
|---|---|
| **Level 1: Novice** | Code fails to compile; business logic remains embedded directly in HTTP handlers; relies on global variables; tests cannot run without external servers. |
| **Level 2: Advanced Beginner** | Domain structs created, but service imports concrete database implementations directly rather than port interfaces. Rollback logic on payment failure is missing. |
| **Level 3: Competent** | Hexagonal boundaries strictly maintained. Port interfaces decoupled from infrastructure. Constructor dependency injection used throughout. Mocks enable 100% test coverage. |
| **Level 4: Proficient** | Flawless domain-driven design. Context propagation maintained across all ports. Error handling wraps sentinel domain errors cleanly. Implements thread-safe adapters with zero race conditions. |

---

## 8. Submission & Reflection

1. Submit the terminal test output from `go test -v -race ./...`.
2. Provide a diagram or dependency graph showing package import directions.
3. Answer: *Why does Go's structural typing (implicit interface fulfillment) make hexagonal architecture easier to maintain compared to nominally typed languages like Java or C#?*
