# Mocking tradeoffs

## Learning objective

Evaluate when mocking is worth the complexity, mock external services using interfaces, use `testify/mock` for call verification, and understand mock fragility.

## Why this matters

Mocks let you verify that a function was called with the right arguments, in the right order, the right number of times. This is essential when testing code that sends emails, charges credit cards, or publishes events — side effects that you cannot assert on by checking state. But mocks introduce tight coupling to implementation details. Over-mocking makes tests brittle: a refactor that changes call patterns breaks tests even if behavior is preserved. Knowing when to mock and when to use a fake is a critical skill.

## Mental model

A mock is a recording spy. It sits between your code and the dependency, recording every call. The test then asks: "Was `SendEmail` called with the arguments `("alice@example.com", "Welcome!")` exactly once?" If yes, the mock passes. If the call happened with different arguments, or not at all, the mock fails.

```
Code → mock.SendEmail(to, subject) → [recorded]
Test → mock.AssertCalled(t, "SendEmail", "alice@...", "Welcome!")
```

The tradeoff: mocks verify *process* (how code was called), while fakes verify *state* (what the result was). Process verification is more precise but more brittle.

## Core idea

Mocking is appropriate when:

| When to mock | Example |
|---|---|
| **External service** | HTTP API, email, SMS, payment gateway |
| **Side effects** | Logging, metrics, audit trail |
| **Hard to observe** | Event publishing, message queue |
| **Non-deterministic** | Random number generator, current time |

Mocking is not appropriate when:

| Don't mock | Alternative |
|---|---|
| **Data storage** | Use an in-memory fake |
| **Simple computation** | Test the result directly |
| **Internal helpers** | Test through the public API |

In Go, mocking is done via interfaces. Your code depends on an interface; the test provides a mock implementation. The most popular library is `testify/mock`:

```go
type EmailSender interface {
    Send(to, subject, body string) error
}

type MockEmailSender struct {
    mock.Mock
}

func (m *MockEmailSender) Send(to, subject, body string) error {
    args := m.Called(to, subject, body)
    return args.Error(0)
}
```

## Under the hood

`testify/mock.Mock` stores expectations in an internal slice. When a method is called, `m.Called(...)` looks up the matching expectation by method name and argument values, records the call, and returns the configured return values.

If no expectation matches, `Called` panics (which `t.Fatal` catches). If a call was expected but never happens, `m.AssertExpectations(t)` reports it.

This is fundamentally different from a fake: a fake has real logic (store and retrieve from a map), while a mock has configured behavior (return error when called with "bad@example.com").

## How Go uses it

The Go standard library does not use `testify/mock` internally, but the ecosystem does. The `net/http/httptest` package provides real servers rather than mocks for HTTP. However, for external API clients, mocks are common:

```go
// client.go
type APIClient interface {
    GetUser(id string) (User, error)
}

// handler.go
func HandleGetUser(w http.ResponseWriter, r *http.Request, client APIClient) {
    // ...
}
```

In tests, `MockAPIClient` verifies that the handler calls `GetUser` with the right ID.

## Go example

```go
package main

import (
	"fmt"
)

// EmailSender sends emails.
type EmailSender interface {
	Send(to, subject, body string) error
}

// Notifier sends notifications via email.
type Notifier struct {
	sender EmailSender
	from   string
}

func NewNotifier(sender EmailSender, from string) *Notifier {
	return &Notifier{sender: sender, from: from}
}

func (n *Notifier) NotifyUser(email, message string) error {
	return n.sender.Send(email, "Notification", message)
}

func main() {
	fmt.Println("See tests for mocking example")
}
```

```go
// main_test.go
package main

import (
	"testing"
	"github.com/stretchr/testify/mock"
)

// MockEmailSender implements EmailSender for testing.
type MockEmailSender struct {
	mock.Mock
}

func (m *MockEmailSender) Send(to, subject, body string) error {
	args := m.Called(to, subject, body)
	return args.Error(0)
}

func TestNotifyUser(t *testing.T) {
	mockSender := new(MockEmailSender)
	notifier := NewNotifier(mockSender, "noreply@example.com")

	mockSender.On("Send", "alice@example.com", "Notification", "Welcome!").Return(nil)

	err := notifier.NotifyUser("alice@example.com", "Welcome!")

	if err != nil {
		t.Errorf("NotifyUser returned error: %v", err)
	}
	mockSender.AssertExpectations(t)
}

func TestNotifyUserFails(t *testing.T) {
	mockSender := new(MockEmailSender)
	notifier := NewNotifier(mockSender, "noreply@example.com")

	mockSender.On("Send", "bad@example.com", "Notification", "Test").Return(fmt.Errorf("send failed"))

	err := notifier.NotifyUser("bad@example.com", "Test")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	mockSender.AssertExpectations(t)
}
```

## Step-by-step execution

Running `TestNotifyUser`:

1. `MockEmailSender` is created.
2. `mockSender.On("Send", "alice@example.com", "Notification", "Welcome!").Return(nil)` registers an expectation: expect `Send` to be called with these exact arguments, return `nil`.
3. `notifier.NotifyUser("alice@example.com", "Welcome!")` calls `mockSender.Send(...)` internally.
4. Inside `MockEmailSender.Send`, `m.Called(to, subject, body)` matches the expectation and returns `nil`.
5. `AssertExpectations` checks that the expected call happened. It did. Pass.

If `Notifier.NotifyUser` had sent to the wrong address, the mock would fail: either no expectation matches (panic) or `AssertExpectations` reports an unmatched call.

## Common mistakes

- **Over-mocking.** Mocking every interface, even for simple data transformations. This leads to brittle tests that break on any refactor.

- **Mocking concrete types.** Go interfaces are implicit. If your code depends on a concrete type, you cannot mock it. Always depend on interfaces.

- **Incorrect argument matchers.** `mock.Anything` is useful but overused. Prefer exact matchers for critical arguments (user ID, email).

- **Not calling `AssertExpectations`.** Without it, a mock that was never called still passes. Always verify.

- **Mocking the same interface multiple times.** Share mock types; do not redefine them per test file.

- **Testing implementation, not behavior.** A mock verifies *how* code was called, not *what* happened. If a refactor changes the call pattern but preserves behavior, the mock test fails unnecessarily.

## Debugging walkthrough

A mock test fails with "assert: mock: I don't know what to return":

```go
mockSender.On("Send", "alice@example.com", "Notification", "Welcome!").Return(nil)
err := notifier.NotifyUser("alice@example.com", "Welcome!")
```

**Symptom**: Panic: `mock: I don't know what to return because the method doesn't return anything`.

**Investigation**: The `Send` method in `EmailSender` has different parameters than expected. The real `Send` signature is `Send(to, subject string) error` but the mock expects three arguments.

**Root cause**: The interface changed but the mock expectation was not updated.

**Fix**: Update the expectation to match the current interface:
```go
mockSender.On("Send", "alice@example.com", "Welcome!").Return(nil)
```

## Production notes

- **Mock external boundaries only.** Mock at the edge of your system: HTTP clients, database drivers, message queues. Do not mock internal collaborators.
- **Prefer integration tests for critical paths.** A mock can say "the call succeeded", but only an integration test with a real API proves the system actually works end-to-end.
- **Fakes + integration tests > mocks.** Use in-memory fakes for most tests, plus a small number of integration tests against real systems. Mocks fill the gap when neither is practical.
- **Mock expectations should be specific.** Use `mock.MatchedBy` for complex matching instead of `mock.Anything`.
- **testify/suite.** For test suites with setup/teardown, `testify/suite` provides struct-based test organization compatible with mocks.

## Performance implications

- Mocks add a method call and a map lookup per mocked call. Overhead is negligible.
- `testify/mock` uses reflection to match arguments. For hot loops, this is measurable but rarely matters in tests.
- The real cost of mocks is maintenance, not CPU time.

## Practice task

Define an interface `PaymentProcessor` with a method `Charge(amount int, currency string) (string, error)`. Write a function `Purchase(processor PaymentProcessor, item string, price int) (string, error)` that charges the card and returns the transaction ID. Mock the `PaymentProcessor` in tests to verify:
- `Charge` is called with the correct amount and currency.
- The transaction ID is returned when charge succeeds.
- An error is returned when the charge fails.

## Tests / verification

```bash
go run ./curriculum/modules/06-testing-debugging-refactoring/lessons/10-mocking-tradeoffs
go test -v ./curriculum/modules/06-testing-debugging-refactoring/lessons/10-mocking-tradeoffs
```

## Review questions

1. What distinguishes a mock from a fake?
2. When would you choose a mock over a fake?
3. What is "mock fragility" and how does it manifest?
4. Why must the dependency be expressed as an interface for mocking to work?
5. What does `mock.AssertExpectations` do and why is it important?

## NEXT UP

Debugging mindset — systematic approaches to diagnosing and fixing bugs.
