# Opslane Flagship — Production Reference Implementation

## Mission

Opslane is a production-grade multi-tenant B2B operations and order fulfillment engine designed to illustrate full-stack Go systems engineering under strict real-world reliability, transactional, and concurrency constraints.

This unified flagship capstone brings together all 15 checkpoints across:
1. Environment configuration and strict invariant schema validation
2. Relational persistence with SQLite / PostgreSQL compatibility
3. Cryptographic token authentication and Role-Based Access Control (RBAC)
4. Strict multi-tenant data boundary isolation
5. HTTP REST routing with standard library Go 1.22 mux patterns
6. Transactional order lifecycle state machine with rollback guarantees
7. Idempotent payment processing with replay deduplication
8. Asynchronous event auditing and backpressured worker pools
9. Singleflight caching layers preventing stampede thundering herds
10. Full observability with structured `slog` logging and Prometheus metrics
11. Phased graceful shutdown with zero in-flight request drops

---

## Architectural Pillars

```
                                  +-----------------------+
                                  |   HTTP Client / Web   |
                                  +-----------+-----------+
                                              |
                          (Bearer Token / Correlation ID / Tenant ID)
                                              v
               +-------------------------------------------------------------+
               |                       Opslane Server                        |
               |                                                             |
               |  +--------------------+             +--------------------+  |
               |  |  Auth & RBAC Mux   |             | Structured Logging |  |
               |  |  (Go 1.22 routing) |             |  (log/slog JSON)   |  |
               |  +---------+----------+             +--------------------+  |
               |            |                                                |
               |            v                                                |
               |  +--------------------+             +--------------------+  |
               |  |  Order & Payment   |<----------->| Singleflight Cache |  |
               |  |   Domain Engine    |             |  (Stampede-Proof)  |  |
               |  +----+---------+-----+             +--------------------+  |
               |       |         |                                           |
               |       |         v                                           |
               |       |  +--------------------+     +--------------------+  |
               |       |  |  Serializable Tx   |     | Bounded Worker     |  |
               |       |  |  (Idempotency Key) |     | Pool Backpressure  |  |
               |       |  +---------+----------+     +---------+----------+  |
               |       |            |                          |             |
               +-------|------------|--------------------------|-------------+
                       |            v                          v
                       |  +--------------------+     +--------------------+
                       +->|  Relational Store  |<----+ Background Event   |
                          | (Multi-Tenant DB)  |     |   Audit Queue      |
                          +--------------------+     +--------------------+
```

---

## Deliverables

1. **Production Reference Server (`_solution/main.go`)**: Fully implemented runnable server with routing, transactional engine, worker pools, singleflight caching, and telemetry.
2. **Student Starter Scaffold (`_starter/main.go`)**: Instructional scaffolding with architectural type models and guided `TODO` milestones.
3. **Automated Verification Suite (`main_test.go`)**:
   - `TestHealthAndReadinessEndpoints`: Probes and liveness validation.
   - `TestMultiTenantIsolation`: Proves Tenant A cannot query or mutate Tenant B records.
   - `TestIdempotentPaymentWorkflow`: Guarantees replay idempotency and single charge.
   - `TestCrossTenantPaymentPrevention`: Verifies defense against unauthorized cross-tenant payments.
   - `TestHTTPAPI_RBACAndFullLifecycle`: Verifies RBAC privileges and Prometheus metrics.

---

## Verification Commands

To execute the test suite against the reference solution:
```bash
go test -v -race ./curriculum/modules/18-flagship-opslane/projects/opslane-flagship/_solution
```

To run the test suite against your student starter:
```bash
go test -v -race ./curriculum/modules/18-flagship-opslane/projects/opslane-flagship/_starter
```

---

## API Endpoints

| Method | Path | Auth / Role | Description |
|---|---|---|---|
| `GET` | `/healthz` | Public | Liveness probe returning service UP |
| `GET` | `/readyz` | Public | Readiness probe checking database connectivity |
| `GET` | `/metrics` | Public | Operations and request telemetry |
| `POST` | `/api/v1/orders` | Operator / Admin | Creates an order scoped to `X-Tenant-ID` |
| `GET` | `/api/v1/orders` | Viewer / Operator / Admin | Lists orders strictly within `X-Tenant-ID` |
| `GET` | `/api/v1/orders/{id}` | Viewer / Operator / Admin | Retrieves single order within tenant boundary |
| `POST` | `/api/v1/orders/{id}/pay` | Operator / Admin | Executes idempotent payment workflow |

---

## Evaluation Rubric

- **Correctness (35%)**: Order state machine adheres to strict progression rules; duplicate idempotency keys return existing transaction without re-execution.
- **Security & Multi-Tenancy (25%)**: Complete data isolation across tenants; viewer roles denied write operations.
- **Concurrency & Performance (25%)**: Background jobs handled via bounded worker pool with backpressure; singleflight caching prevents database stampedes.
- **Reliability & Telemetry (15%)**: All HTTP routes instrumented with correlation IDs, structured logs, atomic metrics, and clean graceful shutdown.
