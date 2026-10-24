# Mockery

## When you need this

Use automated interface mocking (via tools like Mockery or test doubles) when unit testing business logic in services that interact with external databases, third-party payment APIs, or cloud storage interfaces that cannot or should not be invoked during fast local test suites.

## When you do not need this

Do not create mocks for simple data structures, utility functions, or when an in-memory fake (like an in-memory map repository) is trivial to construct. Never mock internal private implementation details of the package under test.

## Prerequisites

- core-16-15

## Mental model

A mock is a test double that implements an interface, records which methods were called with what parameters, and returns predetermined responses configured by the test. Mockery automates the generation of these mocks directly from Go interface definitions.

## Visual model

```mermaid
flowchart LR
    subgraph TestExecution["Unit Test (order_service_test.go)"]
        TEST["TestCreateOrder()"]
        MOCK["MockPaymentGateway"]
        SVC["OrderService"]
    end

    subgraph Expectations["Configured Behavior"]
        EXPECT["mock.On('Charge', 100).Return(true, nil)"]
        ASSERT["mock.AssertExpectations(t)"]
    end

    TEST --> EXPECT
    TEST -->|Calls ProcessOrder()| SVC
    SVC -->|Invokes Charge(100)| MOCK
    TEST --> ASSERT
```

## How it works

Mockery uses Go's go/packages AST parser to discover interfaces across your codebase. For each interface method, it generates a struct embedding testify/mock.Mock. The generated method calls m.Called(args...) and returns the mapped types, enabling assertion of call counts and parameters in tests.

## Go example

```go
package main

import (
	"errors"
	"fmt"
)

// PaymentGateway defines the contract for external transaction processing.
type PaymentGateway interface {
	Charge(userID string, amount int) error
}

// MockPaymentGateway simulates a Mockery-generated test double.
type MockPaymentGateway struct {
	Calls        map[string]int
	ReturnErrors map[string]error
}

func NewMockPaymentGateway() *MockPaymentGateway {
	return &MockPaymentGateway{
		Calls:        make(map[string]int),
		ReturnErrors: make(map[string]error),
	}
}

func (m *MockPaymentGateway) Charge(userID string, amount int) error {
	m.Calls[userID]++
	return m.ReturnErrors[userID]
}

// CheckoutService processes purchases.
type CheckoutService struct {
	gateway PaymentGateway
}

func (s *CheckoutService) Process(userID string, amount int) error {
	if amount <= 0 {
		return errors.New("invalid amount")
	}
	return s.gateway.Charge(userID, amount)
}

// Demo demonstrates Mockery and satisfies test requirements.
func Demo() string {
	mock := NewMockPaymentGateway()
	service := CheckoutService{gateway: mock}
	_ = service.Process("user_101", 500)
	return "Mockery verified"
}

func main() {
	fmt.Println(Demo())
}
```

## Step-by-step execution

1. Developer defines an interface in domain code (e.g. PaymentGateway).
2. Mockery parses the interface and generates MockPaymentGateway implementing all methods.
3. In a test, the developer configures expectations: mock.On("Charge", ...).Return(...).
4. The service is injected with the mock and executed.
5. The test calls mock.AssertExpectations(t) to verify that all expected calls occurred with exact parameters.

## Tradeoffs

- Advantage: True Unit Isolation: Tests execute in milliseconds without Docker containers or real network dependencies.
- Advantage: Precise Failure Simulation: Easily simulate rare edge cases like payment gateway timeouts or socket resets.
- Disadvantage: Refactoring Friction: Over-mocking binds tests tightly to internal implementation details rather than observable behavior.

## Common mistakes

| Mistake | Why it happens | Fix |
|---|---|---|
| Mocking everything instead of using real domain logic | Over-zealous isolation | Only mock external I/O boundaries (databases, HTTP, third-party APIs) |
| Forgetting to assert expectations at the end of the test | Tests pass even if methods were never called | Always call mock.AssertExpectations(t) |

## Debugging walkthrough

### Scenario: Test fails with "Unexpected call to Charge()"

**Cause:** The service invoked the mock with arguments that did not match the configured expectation (e.g., passing "user_102" instead of "user_101").

**Fix:** Inspect the diff in the test error output, verify whether the service's argument preparation logic is correct, and adjust either the production logic or test expectation accordingly.

## Practice task

Define an EmailSender interface with a Send(to, body string) error method. Write a mock implementation with call counting, and write a test for an OnboardingService that asserts the welcome email is sent exactly once.

## Production notes

Mockery is the de-facto mocking tool in Go enterprise microservices, generating mocks automatically via //go:generate mockery --name=Repository to support comprehensive test suites.

## Related core lessons

- Module 15.15 — Demo preparation
- Module 17.15 — Demo preparation
