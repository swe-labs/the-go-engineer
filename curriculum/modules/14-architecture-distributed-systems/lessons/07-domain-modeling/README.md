# Domain modeling

## Learning objective

Model business domains in Go using entities, value objects, and aggregates. Capture domain logic within struct types and raise domain events for cross-context communication.

## Why this matters

Anemic domain models are widespread: structs with only fields and getters/setters, while business logic lives in services that check and manipulate these structs from the outside. This separates the data from the rules that govern it. When the rules change, you must find and update every service that touches the data. A rich domain model encapsulates data _and_ behavior together, making the code self-documenting and resistant to inconsistencies.

## Mental model

Think of a traffic intersection. The intersection has rules (stop signs, traffic lights, right-of-way). If the rules are written on a sign outside the intersection, drivers must read the sign, remember the rule, and apply it. If a driver forgets, they cause an accident. If the rules are embedded in the intersection itself (a traffic light with sensors), the intersection enforces the rules automatically.

A rich domain model is the traffic light. The data (which cars are waiting) and the rules (when to switch lights) are in the same place. An anemic model is the sign: the data (cars at intersection) and the rules (who goes first) are separate, and someone outside must coordinate them.

## Core idea

Domain-Driven Design (DDD) provides three building blocks:

**Entities** are objects with a distinct identity that persists over time and across changes. Two entities with the same field values but different IDs are different. Example: an `Order` entity has an `OrderID` that identifies it uniquely. You can change the order status, add items, and change the shipping address -- it is still the same order.

```go
type Order struct {
    ID        string
    Customer  string
    Status    OrderStatus
    CreatedAt time.Time
}
```

**Value objects** are objects defined only by their attributes. Two value objects with the same field values are interchangeable. Value objects are immutable. Example: an `Address` value object. If you change your address, you do not modify the old address -- you replace it with a new one.

```go
type Address struct {
    Street  string
    City    string
    ZipCode string
}
```

Value objects should enforce their own invariants at construction time. An `Address` should reject an empty street in its constructor.

**Aggregates** are clusters of entities and value objects treated as a single unit. An `Order` aggregate contains the order entity, its items (entities), the shipping address (value object), and the payment info (value object). External code references the aggregate only through its root entity (the `Order`). This ensures consistency: if you need to add an item, you go through `order.AddItem(...)` not through a separate `orderItems` repository.

**Domain events** record something significant that happened in the domain. The aggregate appends events that other parts of the system can react to. An `Order.Ship()` method might raise an `OrderShipped` event that triggers a notification service.

## Under the hood

Go structs with methods are the implementation mechanism. There is no framework, no annotation, no ORM base class. A domain model is plain Go code:

```go
func (o *Order) Ship() error {
    if o.Status != OrderPending {
        return errors.New("only pending orders can be shipped")
    }
    o.Status = OrderShipped
    o.Events = append(o.Events, OrderEvent{Type: "order.shipped"})
    return nil
}
```

The method is the business rule. It cannot be bypassed because it is attached to the struct. If someone tries to set `order.Status = "shipped"` directly, they are fighting the type system -- and the type system wins at compile time if the Status field is unexported.

Immutability for value objects is achieved by exporting fields but not providing mutating methods. Copy-on-write is the pattern:

```go
func (a Address) WithCity(city string) Address {
    return Address{Street: a.Street, City: city, ZipCode: a.ZipCode}
}
```

## How Go uses it

DDD in Go is pragmatic. Most production Go projects use a subset of DDD concepts:

- **Entities with methods**: `type Order struct { ... }` with methods like `Ship()`, `Cancel()`, `AddItem()`. The methods enforce business rules.
- **Value objects for common concepts**: `Money`, `Address`, `Email`, `PhoneNumber`. These are often newtype wrappers with validation in the constructor.
- **Aggregate roots with repository interfaces**: repositories operate at the aggregate level. `OrderRepository` saves and loads entire `Order` aggregates.
- **Domain events for cross-context communication**: events are recorded on the aggregate and published after the aggregate is persisted.

Not every project needs full DDD. If the domain is CRUD (create-read-update-delete with no business rules), anemic models are acceptable. DDD pays off when the domain has complex, evolving rules.

## Go example

```go
package main

import (
	"errors"
	"fmt"
	"time"
)

type Address struct {
	Street  string
	City    string
	ZipCode string
}

func NewAddress(street, city, zipCode string) (Address, error) {
	if street == "" || city == "" || zipCode == "" {
		return Address{}, errors.New("all address fields are required")
	}
	return Address{Street: street, City: city, ZipCode: zipCode}, nil
}

type OrderItem struct {
	ProductID string
	Quantity  int
	UnitPrice float64
}

func NewOrderItem(productID string, quantity int, unitPrice float64) (OrderItem, error) {
	if productID == "" {
		return OrderItem{}, errors.New("product ID is required")
	}
	if quantity <= 0 {
		return OrderItem{}, errors.New("quantity must be positive")
	}
	if unitPrice <= 0 {
		return OrderItem{}, errors.New("unit price must be positive")
	}
	return OrderItem{ProductID: productID, Quantity: quantity, UnitPrice: unitPrice}, nil
}

func (i OrderItem) Total() float64 {
	return float64(i.Quantity) * i.UnitPrice
}

type OrderStatus string

const (
	OrderPending   OrderStatus = "pending"
	OrderShipped   OrderStatus = "shipped"
	OrderDelivered OrderStatus = "delivered"
)

type OrderEvent struct {
	OrderID   string
	EventType string
	Timestamp time.Time
	Data      interface{}
}

type Order struct {
	ID        string
	Customer  string
	Items     []OrderItem
	Shipping  Address
	Status    OrderStatus
	Events    []OrderEvent
	CreatedAt time.Time
}

func NewOrder(id, customer string, items []OrderItem, shipping Address) (*Order, error) {
	if id == "" {
		return nil, errors.New("order ID is required")
	}
	if customer == "" {
		return nil, errors.New("customer is required")
	}
	if len(items) == 0 {
		return nil, errors.New("order must have at least one item")
	}
	order := &Order{
		ID: id, Customer: customer, Items: items,
		Shipping: shipping, Status: OrderPending, CreatedAt: time.Now(),
	}
	order.raiseEvent("order.created", nil)
	return order, nil
}

func (o *Order) Total() float64 {
	var total float64
	for _, item := range o.Items {
		total += item.Total()
	}
	return total
}

func (o *Order) Ship() error {
	if o.Status != OrderPending {
		return errors.New("only pending orders can be shipped")
	}
	o.Status = OrderShipped
	o.raiseEvent("order.shipped", nil)
	return nil
}

func (o *Order) raiseEvent(eventType string, data interface{}) {
	o.Events = append(o.Events, OrderEvent{
		OrderID: o.ID, EventType: eventType,
		Timestamp: time.Now(), Data: data,
	})
}

func main() {
	addr, _ := NewAddress("123 Main St", "Springfield", "12345")
	item, _ := NewOrderItem("prod-1", 2, 19.99)
	order, _ := NewOrder("ORD-001", "Alice", []OrderItem{item}, addr)

	fmt.Printf("Order %s: $%.2f (status: %s)\n", order.ID, order.Total(), order.Status)
	order.Ship()
	fmt.Printf("Order shipped: %s\n", order.Status)
	fmt.Printf("Events: %d raised\n", len(order.Events))
}
```

## Step-by-step execution

For `NewOrder("ORD-001", "Alice", items, addr)`:

1. The constructor validates ID, customer, and items. Empty ID would return an error.
2. An `Order` struct is allocated with initial values and status `OrderPending`.
3. The `order.created` domain event is appended to `Events`.
4. The `*Order` pointer is returned.

For `order.Ship()`:

1. `Ship` checks that `o.Status == OrderPending`. It is, so the method proceeds.
2. `o.Status` is set to `OrderShipped`.
3. An `order.shipped` event is appended to `Events`.
4. Nil is returned. If the order had already been shipped, an error would be returned and the status would not change.

The value objects (`Address`, `OrderItem`) are valid at construction. The entity (`Order`) enforces state transitions through methods. The aggregate holds domain events that downstream consumers can process.

## Common mistakes

- **Anemic domain model**: the `Order` struct has only fields and getters. Business logic like `Ship()` lives in a service. This means any service can change order status without going through the rule check. Put methods on the struct.
- **Mutable value objects**: `Address` has a `SetStreet` method that modifies the existing instance. This breaks value semantics. Create a new `Address` with the updated value instead.
- **Giant aggregate**: an `Order` aggregate that contains `Customer`, `PaymentHistory`, `ShipmentTracking`, and `ReturnRequests`. The aggregate should be focused. If you rarely need all that data together, split the aggregate.
- **Ignoring domain events**: events are appended but never published or consumed. Events are only useful if something processes them. The service layer should flush events after saving the aggregate.

## Debugging walkthrough

Consider an order that was shipped before payment was confirmed:

```go
order := loadOrder("ORD-001")
order.Status = "shipped" // bypasses Ship() method
```

**Symptom**: The order is shipped but payment was never collected. The customer receives the product and the business loses money.

**Root cause**: The `Status` field is exported and settable from outside the struct. The `Ship()` method's guard (`if o.Status != OrderPending`) can be bypassed.

**Fix**: Make `Status` unexported and expose it through a read-only method:

```go
type Order struct {
    status OrderStatus
}
func (o *Order) Status() OrderStatus { return o.status }
func (o *Order) Ship() error {
    if o.status != OrderPending { return errors.New("...") }
    o.status = OrderShipped
    return nil
}
```

Now the only way to change status is through a method that enforces the business rule.

## Production notes

- **Repository persistence**: the repository loads and saves aggregates. It calls the aggregate constructor to build the type from database rows. The aggregate's events are flushed after save.
- **Event publishing**: after the repository saves the aggregate, the service layer publishes domain events to a message bus. This ensures at-most-once semantics within the transaction boundary.
- **Value object storage**: value objects like `Address` are often stored as columns in the same table as the aggregate root. JSON columns work well for complex value objects.
- **Bounded context alignment**: entities in different bounded contexts can have different representations. The `Order` in ordering context has `Status` as `pending/shipped/delivered`. The `Order` in billing context has `Status` as `unpaid/paid/refunded`. They are different types.

## Performance implications

- **Validation in constructors**: value object validation happens on every construction. For batch operations (importing 100K orders), this can add up. Consider bulk constructors that defer validation.
- **Event slice allocation**: appending domain events to a slice allocates. For aggregates that generate many events, pre-allocate the slice.
- **Aggregate size**: loading a large aggregate with hundreds of items from the database loads everything into memory. Consider lazy-loading or paginating child collections in extreme cases.

## Practice task

Add a `Cancel()` method to the `Order` aggregate. The rules:

1. Only pending orders can be cancelled.
2. Cancelling sets status to a new `OrderCancelled` status.
3. Cancelling raises an `order.cancelled` domain event.
4. If any items have been shipped individually (add a `ShippedQuantity` field to `OrderItem`), cancellation of those items is denied.

Write table-driven tests for all paths: successful cancellation, cancelling a shipped order, and partial cancellation.

## Tests / verification

```bash
go run ./curriculum/modules/14-architecture-distributed-systems/lessons/07-domain-modeling
go test ./curriculum/modules/14-architecture-distributed-systems/lessons/07-domain-modeling
```

The existing tests verify address creation, order item total calculation, order creation with validation, and shipping transitions. After completing the practice task, add tests for the cancellation logic.

## Review questions

1. What is the difference between an entity and a value object in DDD? Give a Go example of each.
2. Why should an aggregate be loaded and saved as a single unit?
3. How do domain events enable communication between bounded contexts?
4. What is the risk of making all struct fields exported (public)?
5. When would you choose an anemic domain model over a rich domain model?

## NEXT UP

Invariants -- enforcing business rules at system boundaries to prevent invalid state.
