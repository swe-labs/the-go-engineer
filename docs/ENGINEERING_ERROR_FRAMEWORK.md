# Engineering Error Framework

## The Gap

Current curriculum teaches:

```go
func doSomething() error {
    return fmt.Errorf("something failed")
}
```

But does NOT teach:

- When to return vs wrap vs panic
- User error vs system error vs fatal error
- Error as control flow

This is a critical gap for "expert engineer" thinking.

---

## The Framework: Three Error Categories

### 1. User Errors (Return)

**Definition**: Input validation failures, business logic violations  
**Handle**: Return immediately, no logging needed (caller decides)

```go
type UserError struct {
    Code    string
    Message string
}

func (e *UserError) Error() string {
    return e.Message
}

// Examples:
// - Invalid email format
// - Insufficient balance
// - Item not found
// - Permission denied
```

**When to use**:

- Input validation
- Business rule violations
- "Not found" scenarios

**Pattern**:

```go
if !isValid(email) {
    return nil, &UserError{
        Code:    "INVALID_EMAIL",
        Message: "email is invalid",
    }
}
```

---

### 2. System Errors (Wrap + Propagate)

**Definition**: Infrastructure failures, external service failures  
**Handle**: Wrap with context, propagate to caller

```go
type SystemError struct {
    Code    string
    Message string
    Cause   error
}

func (e *SystemError) Error() string {
    return fmt.Sprintf("%s: %v", e.Message, e.Cause)
}

func (e *SystemError) Unwrap() error {
    return e.Cause
}

// Examples:
// - Database connection failed
// - HTTP request timeout
// - File read error
// - External API failure
```

**When to use**:

- Database operations
- Network calls
- File I/O
- External service calls

**Pattern**:

```go
func getUser(id string) (*User, error) {
    row := db.QueryRow("SELECT * FROM users WHERE id = ?", id)

    var user User
    err := row.Scan(&user.ID, &user.Name, &user.Email)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, &UserError{
                Code:    "USER_NOT_FOUND",
                Message: fmt.Sprintf("user %s not found", id),
            }
        }

        // Wrap system error, preserve original
        return nil, &SystemError{
            Code:    "DB_QUERY_FAILED",
            Message: "failed to query user",
            Cause:   err,
        }
    }
    return &user, nil
}
```

---

### 3. Fatal Errors (Log + Exit)

**Definition**: Program bugs, unrecoverable states  
**Handle**: Log and exit gracefully

```go
type FatalError struct {
    Message string
    Cause   error
}

func (e *FatalError) Error() string {
    return fmt.Sprintf("FATAL: %s: %v", e.Message, e.Cause)
}

func HandleFatal(err error) {
    log.Error(err)
    // Cleanup if needed
    os.Exit(1)
}

// Examples:
// - Database schema migration failure
// - Configuration file missing
// - Required environment variable not set
```

**When to use**:

- Startup configuration errors
- Database connection failure at startup
- Missing required resources
- Bug detection (should not happen)

**Pattern**:

```go
func main() {
    cfg, err := loadConfig()
    if err != nil {
        if _, ok := err.(*ConfigError); ok {
            // Config errors are fatal - can't start
            log.Fatal("Failed to load configuration: ", err)
        }
    }

    db, err := connectDB(cfg)
    if err != nil {
        log.Fatal("Failed to connect to database: ", err)
    }
}
```

---

## Error Hierarchy

```
Error Interface
    │
    ├── UserError (validation, business logic)
    │   └── Caller handles: show user-friendly message
    │
    ├── SystemError (infrastructure, external)
    │   └── Caller handles: wrap + propagate + log
    │
    └── FatalError (bugs, unrecoverable)
        └── Caller handles: log + exit
```

---

## Decision Tree

```
Is this error caused by user input or business logic?
    │
    ├── YES → UserError → Return to caller
    │
    └── NO → Is this a bug or unexpected failure?
              │
              ├── BUG (should not happen) → FatalError → Log + Exit
              │
              └── UNEXPECTED (can happen) → SystemError → Wrap + Propagate
```

---

## In Practice

### Handler Layer

```go
func (h *Handler) HandleError(w http.ResponseWriter, err error) {
    var userErr *UserError
    var sysErr *SystemError

    switch {
    case errors.As(err, &userErr):
        // User error → 400 Bad Request
        http.Error(w, userErr.Message, http.StatusBadRequest)

    case errors.Is(err, context.DeadlineExceeded):
        // System error → 504 Gateway Timeout
        http.Error(w, "request timeout", http.StatusGatewayTimeout)

    case errors.Is(err, sql.ErrNoRows):
        // Not found → 404
        http.Error(w, "resource not found", http.StatusNotFound)

    default:
        // Unknown error → 500 Internal Server Error
        // Log details for debugging, show generic to user
        log.Error(err)
        http.Error(w, "internal server error", http.StatusInternalServerError)
    }
}
```

### Service Layer

```go
func (s *OrderService) CreateOrder(req CreateOrderRequest) (*Order, error) {
    // Validate input → UserError
    if err := validateOrderRequest(req); err != nil {
        return nil, err
    }

    // Check inventory → might return UserError
    if !s.hasInventory(req.Items) {
        return nil, &UserError{
            Code:    "INSUFFICIENT_INVENTORY",
            Message: "not enough inventory",
        }
    }

    // Save to DB → might return SystemError
    order, err := s.repo.CreateOrder(req)
    if err != nil {
        return nil, &SystemError{
            Code:    "DB_CREATE_FAILED",
            Message: "failed to create order",
            Cause:   err,
        }
    }

    return order, nil
}
```

### Repository Layer

```go
func (r *OrderRepo) CreateOrder(req CreateOrderRequest) (*Order, error) {
    // Query DB
    // If error: return wrapped SystemError
    // NEVER return raw error to upper layer
}
```

---

## The Senior Engineer Mindset

A junior engineer writes:

```go
if err != nil {
    return err
}
```

A senior engineer writes:

```go
if err != nil {
    if errors.Is(err, sql.ErrNoRows) {
        return nil, fmt.Errorf("user not found: %w", err)
    }
    return nil, fmt.Errorf("database error: %w", err)
}
```

But a **Google-level engineer** writes:

```go
if err != nil {
    switch {
    case errors.Is(err, sql.ErrNoRows):
        return nil, &UserError{Code: "NOT_FOUND", Message: "resource not found"}
    case errors.Is(err, context.DeadlineExceeded):
        return nil, &SystemError{Code: "TIMEOUT", Message: "request timeout", Cause: err}
    default:
        return nil, &SystemError{Code: "INTERNAL", Message: "unexpected error", Cause: err}
    }
}
```

---

## Production Error Patterns

### 1. Error Codes (Not Just Messages)

```go
// Instead of just messages:
errors.New("failed to connect to database")

// Use error codes:
&SystemError{
    Code:    "DB_CONNECTION_FAILED",
    Message: "failed to connect to database",
    Cause:   err,
}
```

**Why**:

- Programmatic handling
- Client can understand error type
- Translation to user-facing messages

---

### 2. Error Wrapping

```go
// WRONG - loses original error
return errors.New("query failed")

// RIGHT - preserves original error
return fmt.Errorf("query user by id: %w", err)

// BEST - custom error with cause
return &SystemError{
    Code:    "DB_QUERY_FAILED",
    Message: "query user by id",
    Cause:   err,
}
```

---

### 3. Error Handling in HTTP

```go
func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    err := h.handle(w, r)
    if err != nil {
        h.logError(r, err)
        h.writeError(w, err)
    }
}

func (h Handler) writeError(w http.ResponseWriter, err error) {
    var userErr *UserError
    if errors.As(err, &userErr) {
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(map[string]string{
            "error":   userErr.Code,
            "message": userErr.Message,
        })
        return
    }

    // Default: 500
    w.WriteHeader(http.StatusInternalServerError)
    // Don't leak internal error details!
    json.NewEncoder(w).Encode(map[string]string{
        "error":   "INTERNAL_ERROR",
        "message": "something went wrong",
    })
}
```

---

## Anti-Patterns to Avoid

### 1. Never Ignore Errors

```go
// BAD
_ = someFunction()

// ACCEPTABLE in specific cases
func main() {
    // At startup, configuration errors should fatal
    if err := initConfig(); err != nil {
        log.Fatal(err)
    }
}

// NEVER in production code
func process() {
    result, _ := doSomething()  // Don't know if it failed!
}
```

### 2. Never Return Raw External Errors

```go
// BAD - leaks implementation
func GetUser(id string) (*User, error) {
    user, err := db.QueryRow(...)
    return user, err  // Exposes DB details!
}

// GOOD - wraps in system error
func GetUser(id string) (*User, error) {
    user, err := db.QueryRow(...)
    if err != nil {
        return nil, &SystemError{
            Code: "DB_QUERY_FAILED",
            Cause: err,
        }
    }
    return user, nil
}
```

### 3. Never Use Panic for Control Flow

```go
// BAD - panic is for fatal errors
func findUser(id string) *User {
    user, err := db.Find(id)
    if err != nil {
        panic(err)  // Crashes entire program!
    }
    return user
}

// GOOD - return error
func findUser(id string) (*User, error) {
    user, err := db.Find(id)
    if err != nil {
        return nil, err
    }
    return user, nil
}
```

---

## Summary

| Error Type      | When                             | How              | HTTP Status   |
| --------------- | -------------------------------- | ---------------- | ------------- |
| **UserError**   | Input validation, business rules | Return           | 400, 404, 403 |
| **SystemError** | Infrastructure, external         | Wrap + Propagate | 500, 502, 503 |
| **FatalError**  | Bugs, unrecoverable              | Log + Exit       | N/A (exit)    |

---

_This framework should be taught after basic error handling and before production engineering._
