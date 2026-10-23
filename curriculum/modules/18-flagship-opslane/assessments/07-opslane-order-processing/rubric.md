# Opslane Milestone 07 Assessment — Rubric

## Assessment Dimensions

### 1. State Machine Rigor & Transition Validation (30%)
- **Excellent (5 pts)**: Complete FSM state model; transitions strictly validated via lookup table; terminal states (`COMPLETED`, `CANCELLED`) protected from modification; invalid transitions return explicit error.
- **Proficient (3 pts)**: States defined, but validation uses ad-hoc conditionals rather than centralized matrix.
- **Developing (1 pt)**: Terminal states can be overwritten or transitions unverified.
- **Unacceptable (0 pts)**: No state validation implemented.

### 2. Transactional Atomicity & ACID Semantics (25%)
- **Excellent (5 pts)**: Workflows execute within `BeginTx` with serializable isolation; `defer tx.Rollback()` ensures safe abort; status update and audit log committed atomically; payment failure rolls back cleanly.
- **Proficient (3 pts)**: Transaction used, but `defer tx.Rollback()` missing or audit log written in separate uncommitted query.
- **Developing (1 pt)**: Operations execute via un-transactional queries (`db.Exec`), risking partial state corruption.
- **Unacceptable (0 pts)**: Transactions fail or crash.

### 3. External Integration & Dependency Inversion (25%)
- **Excellent (5 pts)**: `PaymentGateway` decoupled via interface; mock implementation enables deterministic testing of success and failure branches; context propagated to external calls.
- **Proficient (3 pts)**: Gateway interface used, but mock lacks failure simulation.
- **Developing (1 pt)**: Payment processor concrete dependency hardcoded directly.
- **Unacceptable (0 pts)**: Cannot test without live payment service.

### 4. Conceptual & Reliability Architecture Defense (20%)
- **Excellent (5 pts)**: Accurately explains external call risks in transactions, serializable isolation semantics, idempotency mechanics, and state machine design benefits.
- **Proficient (3 pts)**: Answers conceptual questions accurately with minor gaps in isolation levels or lock exhaustion.
- **Developing (1 pt)**: Cannot articulate why transactions must roll back on payment failure.
- **Unacceptable (0 pts)**: No explanation provided.

## Grading Thresholds
- **12–14 pts**: Exceeds expectations (Gold standard)
- **10–11 pts**: Meets expectations (Passing)
- **7–9 pts**: Needs revision
- **0–6 pts**: Unsatisfactory — review prerequisites and retake
