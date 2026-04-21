# GoScale SaaS Backend - Flagship Project Specification

## Vision

A brutal, production-grade SaaS backend that simulates real-world complexity. This is where engineering thinking is tested, broken, and rebuilt.

**This is NOT a tutorial. This is a battlefield.**

---

## Project Overview

### GoScale - Multi-tenant SaaS Platform

A multi-tenant SaaS backend with:

- User authentication & authorization
- Tenant isolation
- Order processing system
- Payment pipeline (mock)
- Event-driven architecture
- Worker pools for async processing
- Rate limiting
- Structured logging + tracing
- Full observability stack

---

## Technical Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                              GoScale Architecture                           │
└─────────────────────────────────────────────────────────────────────────────┘

                              ┌─────────────────┐
                              │   Load Balancer │
                              │   (simulated)   │
                              └────────┬────────┘
                                       │
                    ┌──────────────────┼──────────────────┐
                    │                  │                  │
                    ▼                  ▼                  ▼
            ┌──────────────┐   ┌──────────────┐   ┌──────────────┐
            │  HTTP API    │   │  HTTP API    │   │  HTTP API    │
            │  (Primary)   │   │  (Admin)     │   │  (Public)    │
            └──────┬───────┘   └──────┬───────┘   └──────┬───────┘
                   │                  │                  │
                   └──────────────────┼──────────────────┘
                                      │
                                      ▼
                            ┌─────────────────┐
                            │  Middleware     │
                            │  Stack          │
                            │  - Auth         │
                            │  - Rate Limit   │
                            │  - Logging      │
                            │  - Tracing      │
                            └────────┬────────┘
                                     │
                    ┌────────────────┼────────────────┐
                    │                │                │
                    ▼                ▼                ▼
            ┌──────────────┐   ┌──────────────┐   ┌──────────────┐
            │  User        │   │  Order       │   │  Payment     │
            │  Service     │   │  Service     │   │  Service     │
            └──────┬───────┘   └──────┬───────┘   └──────┬───────┘
                   │                  │                  │
                   └──────────────────┼──────────────────┘
                                      │
                                      ▼
                            ┌─────────────────┐
                            │  Event Bus      │
                            │  (In-memory)    │
                            └────────┬────────┘
                                     │
                    ┌────────────────┼────────────────┐
                    │                │                │
                    ▼                ▼                ▼
            ┌──────────────┐   ┌──────────────┐   ┌──────────────┐
            │  Worker      │   │  Worker      │   │  Worker      │
            │  Pool        │   │  Pool        │   │  Pool        │
            │  (Orders)    │   │  (Payments) │   │  (Notifs)    │
            └──────────────┘   └──────────────┘   └──────────────┘
                                      │
                                      ▼
                            ┌─────────────────┐
                            │  Database       │
                            │  (SQLite)       │
                            │  + Redis Cache  │
                            └─────────────────┘
```

---

## Module Breakdown

### Module 1: Foundation & Configuration

**Goal**: Set up project structure and configuration system

**Files**:

- `cmd/server/main.go` - Entry point
- `internal/config/config.go` - Configuration management
- `internal/config/environment.go` - Environment handling
- `.env.example` - Environment template

**Engineering Challenges**:

- Environment-based configuration
- Secret management
- Default values vs overrides
- Configuration validation

**Production Notes**:

```
⚠️ In production:
- Never commit .env files
- Use secrets manager (Vault, AWS Secrets Manager)
- Validate config at startup, not at runtime
- Different configs for dev/staging/prod
```

---

### Module 2: Database & Models

**Goal**: Set up database layer with proper modeling

**Files**:

- `internal/models/user.go`
- `internal/models/tenant.go`
- `internal/models/order.go`
- `internal/models/payment.go`
- `internal/db/migrations.go`
- `internal/db/repository.go`

**Engineering Challenges**:

- Schema design for multi-tenancy
- Indexing strategies
- Connection pool management
- Transaction handling
- Race conditions in DB operations

**Production Notes**:

```
⚠️ In production:
- Connection pool exhaustion causes system-wide outage
- Never execute raw SQL without parameterized queries
- Use transactions for multi-step operations
- Implement proper retry logic for transient failures
- Monitor DB connections as critical metric
```

---

### Module 3: Authentication System

**Goal**: JWT-based auth with role management

**Files**:

- `internal/auth/jwt.go`
- `internal/auth/password.go`
- `internal/auth/middleware.go`
- `internal/models/user.go` (extends)

**Engineering Challenges**:

- JWT token handling
- Password hashing (bcrypt)
- Token expiration
- Refresh token flow
- Role-based access control (RBAC)
- Tenant isolation

**Production Notes**:

```
⚠️ In production:
- JWT secrets must be rotated
- Store refresh tokens securely
- Implement account lockout after N failed attempts
- Log all authentication failures
- Token invalidation (logout) requires blacklist or short expiry
```

---

### Module 4: HTTP API Layer

**Goal**: Build REST API with middleware stack

**Files**:

- `internal/handlers/user.go`
- `internal/handlers/order.go`
- `internal/handlers/payment.go`
- `internal/middleware/auth.go`
- `internal/middleware/ratelimit.go`
- `internal/middleware/logging.go`
- `internal/middleware/cors.go`

**Engineering Challenges**:

- Proper HTTP timeouts (Read/Write/Idle)
- Request validation
- Error response standardization
- Rate limiting (per tenant/user)
- Concurrent request handling
- Resource cleanup

**Production Notes**:

```
⚠️ In production:
- Missing HTTP timeouts = memory leak + resource exhaustion
- Slow clients can hold connections forever
- Rate limiting prevents DDoS
- Always validate input - never trust client
- Return consistent error formats
- Log request IDs for tracing
```

---

### Module 5: Order Processing System

**Goal**: Core business logic with concurrency

**Files**:

- `internal/services/order.go`
- `internal/services/inventory.go`
- `internal/services/validation.go`

**Engineering Challenges**:

- Order state machine
- Concurrent order processing
- Optimistic locking
- Idempotency keys
- Race conditions in inventory
- Partial fulfillment handling

**Failure Scenarios**:

```
🔴 Test: What happens when 1000 orders arrive simultaneously?
- Inventory race condition
- Database connection pool exhaustion
- Worker pool saturation

🔴 Test: What happens when payment fails after inventory reserved?
- Distributed transaction simulation
- Compensation/rollback logic

🔴 Test: What happens when worker crashes mid-processing?
- Graceful shutdown
- State recovery
```

**Production Notes**:

```
⚠️ In production:
- Always generate idempotency keys for critical operations
- Use database-level locking for inventory
- Implement circuit breaker for external services
- Queue-based processing for reliability
- Monitor order processing latency
```

---

### Module 6: Payment Pipeline (Mock)

**Goal**: Simulate payment processing with failure scenarios

**Files**:

- `internal/services/payment.go`
- `internal/payment/gateway.go`
- `internal/payment/worker.go`

**Engineering Challenges**:

- Mock external payment gateway
- Simulate network failures
- Retry with exponential backoff
- Webhook handling
- Transaction reconciliation

**Failure Scenarios**:

```
🔴 Test: Payment gateway timeout
- How does order system handle?

🔴 Test: Duplicate payment attempts
- Idempotency check

🔴 Test: Payment succeeds but webhook fails
- Reconciliation job

🔴 Test: Race between refund and chargeback
- Proper state machine
```

---

### Module 7: Event System & Worker Pools

**Goal**: Event-driven architecture with bounded concurrency

**Files**:

- `internal/events/bus.go`
- `internal/events/types.go`
- `internal/workers/pool.go`
- `internal/workers/order_processor.go`
- `internal/workers/payment_processor.go`
- `internal/workers/notification_worker.go`

**Engineering Challenges**:

- Pub/sub event bus
- Bounded worker pools (backpressure)
- Goroutine leak prevention
- Graceful shutdown (drain pattern)
- Dead letter queue
- Event ordering

**Critical Questions**:

```
❓ What happens when worker pool is saturated?
- Backpressure mechanism
- Reject new work vs queue buildup

❓ What happens when event is published but no subscriber?
- Dead letter handling

❓ How do you debug a stuck worker?
- Instrumentation
- Tracing
```

**Production Notes**:

```
⚠️ In production:
- Never use unbounded goroutines
- Always implement backpressure
- Monitor goroutine count (leak detection)
- Implement graceful shutdown with drain
- Log all worker state transitions
- Use errgroup for coordinated shutdown
```

---

### Module 8: Caching Layer

**Goal**: Redis-style caching (using in-memory for simplicity)

**Files**:

- `internal/cache/cache.go`
- `internal/cache/redis.go` (interface)
- `internal/middleware/cache.go`

**Engineering Challenges**:

- Cache invalidation patterns
- Cache-aside pattern
- Write-through vs write-back
- TTL management
- Cache stampede prevention

---

### Module 9: Observability

**Goal**: Structured logging, metrics, tracing

**Files**:

- `internal/logging/logger.go`
- `internal/logging/middleware.go`
- `internal/metrics/metrics.go`
- `internal/tracing/tracing.go`

**Engineering Challenges**:

- Structured JSON logging
- Log levels (debug, info, warn, error)
- Request tracing (correlation IDs)
- Performance metrics
- Error rate monitoring

**Production Notes**:

```
⚠️ In production:
- Use structured logging (JSON)
- Always include correlation IDs
- Log error stack traces
- Use log levels appropriately
- Implement log sampling for high-traffic endpoints
- Integrate with logging aggregation (ELK, Loki)
```

---

### Module 10: Graceful Shutdown & Deployment

**Goal**: Production deployment considerations

**Files**:

- `cmd/server/shutdown.go`
- `Dockerfile`
- `docker-compose.yml`
- `.github/workflows/ci.yml`

**Engineering Challenges**:

- Signal handling (SIGTERM, SIGINT)
- Drain in-flight requests
- Close database connections
- Flush logs
- Health check endpoint
- Zero-downtime deployment

**Critical Scenarios**:

```
🔴 Test: SIGTERM received during order processing
- Complete current operation?
- Or abort immediately?

🔴 Test: Database connection drain
- Wait for queries to complete?

🔴 Test: Health check during shutdown
- Return unhealthy?
```

---

## Engineering Scenarios (Failure-Based Learning)

### Scenario 1: The Traffic Spike

```
Condition: 10,000 requests/second for 30 seconds

Expected Failures:
- Database connection pool exhaustion
- Worker pool saturation
- Memory exhaustion
- Rate limiter saturation

Learner Task:
1. Identify which component fails first
2. Add proper limits
3. Implement backpressure
4. Add circuit breakers
```

### Scenario 2: The Memory Leak

```
Condition: Run service for 1 hour with steady traffic

Expected Failures:
- Goroutine accumulation
- Unclosed database connections
- Growing slices
- Event bus memory growth

Learner Task:
1. Use pprof to identify leak
2. Fix root cause
3. Add monitoring
```

### Scenario 3: The Race Condition

```
Condition: Two orders for same inventory item

Expected Failures:
- Double allocation
- Negative inventory
- Overselling

Learner Task:
1. Identify race condition
2. Add proper locking
3. Use database transactions
4. Test with -race flag
```

### Scenario 4: The Cascading Failure

```
Condition: Payment gateway is down

Expected Failures:
- Orders stuck in "processing"
- Timeout accumulation
- Worker pool blocking

Learner Task:
1. Implement circuit breaker
2. Add fallback behavior
3. Configure timeouts properly
4. Add retry with backoff
```

### Scenario 5: The Slow Query

```
Condition: Database query takes 30 seconds

Expected Failures:
- Connection pool exhaustion
- Request timeouts
- Thread starvation

Learner Task:
1. Identify slow query with logs
2. Add database indexes
3. Implement query timeouts
4. Add connection pool limits
```

---

## Thinking Sections (Questions for Senior Thinking)

### Architecture

```
❓ Why use modular monolith instead of microservices?
❓ Where are service boundaries?
❓ How do you handle tenant isolation?
❓ What happens when you need to split services?
```

### Concurrency

```
❓ Why use bounded worker pools?
❓ What happens when buffer is full?
❓ How do you implement backpressure?
❓ When to use channels vs mutexes?
```

### Reliability

```
❓ What is the difference between retry and idempotency?
❓ How do you handle partial failures?
❓ When to use transactions vs compensation?
❓ What is circuit breaker pattern?
```

### Performance

```
❓ Where are the hotspots in this system?
❓ When to cache vs compute?
❓ How do you measure performance?
❓ What is the cost of each operation?
```

### Operations

```
❓ How do you debug a 500 error in production?
❓ What metrics matter?
❓ How do you do zero-downtime deploys?
❓ How do you handle a runaway process?
```

---

## Assessment Model

### Phase 1: Code Review (25%)

- Code quality
- Error handling
- Security considerations
- Performance awareness

### Phase 2: Failure Injection (25%)

- Introduce bugs
- Learner must diagnose and fix
- Time limit: 30 minutes per scenario

### Phase 3: Scale Testing (25%)

- Add load
- Identify bottlenecks
- Optimize

### Phase 4: Architecture Discussion (25%)

- Explain design decisions
- Trade-off analysis
- Senior-level thinking demonstration

---

## Success Criteria

A learner who completes GoScale can:

1. ✅ Build production-grade HTTP services
2. ✅ Implement proper error handling strategy
3. ✅ Write concurrent code without races
4. ✅ Design worker pools with backpressure
5. ✅ Handle graceful shutdown
6. ✅ Add proper HTTP timeouts
7. ✅ Implement authentication/authorization
8. ✅ Use structured logging
9. ✅ Debug production issues
10. ✅ Make architectural trade-offs

---

## Implementation Roadmap

| Week | Module | Focus                      |
| ---- | ------ | -------------------------- |
| 1    | 1-2    | Foundation & DB            |
| 2    | 3-4    | Auth & HTTP                |
| 3    | 5-6    | Business Logic             |
| 4    | 7-8    | Concurrency & Caching      |
| 5    | 9-10   | Observability & Deployment |

---

_This is The Go Engineer flagship. Not a tutorial. A battlefield._
