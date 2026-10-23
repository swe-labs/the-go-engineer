# Opslane Milestone 02 Assessment — Rubric

## Assessment Dimensions

### 1. Multi-Tenant Data Modeling & Query Isolation (30%)
- **Excellent (5 pts)**: All tables and domain entities include mandatory `tenant_id` foreign keys; repository queries strictly filter on `tenant_id`; sensitive fields (`PasswordHash`) protected from serialization.
- **Proficient (3 pts)**: Models defined and queries filtered on `tenant_id`, but minor omissions in composite indices or null handling.
- **Developing (1 pt)**: Repository methods accept `orderID` alone without `tenantID` check, enabling cross-tenant data leakage.
- **Unacceptable (0 pts)**: Models lack tenant scoping.

### 2. Migration Engineering & Atomicity (25%)
- **Excellent (5 pts)**: Schema evolution versioned with clear up/down logic; migrations executed within atomic transactions; tracking table records execution timestamp.
- **Proficient (3 pts)**: Migrations execute sequentially, but rollback on DDL failure is omitted.
- **Developing (1 pt)**: Schema created via raw untracked scripts without versioning.
- **Unacceptable (0 pts)**: Migrations fail to run or leave schema in inconsistent state.

### 3. Connection Pooling & Context Propagation (25%)
- **Excellent (5 pts)**: Connection pool parameters (`MaxOpenConns`, `MaxIdleConns`, `ConnMaxLifetime`) tuned defensively; ping verified with context timeout; all queries pass `context.Context`.
- **Proficient (3 pts)**: Database connected and pinged, but pool bounds use defaults and lack timeouts.
- **Developing (1 pt)**: Omits context propagation in repository queries (`Exec` instead of `ExecContext`).
- **Unacceptable (0 pts)**: Connection leaks or fails to open database.

### 4. Conceptual & Distributed Systems Defense (20%)
- **Excellent (5 pts)**: Accurately explains connection pool lifecycle, B-tree composite index ordering, idempotent payment key mechanics, and migration race condition mitigations.
- **Proficient (3 pts)**: Answers conceptual questions accurately with minor gaps in index mechanics or concurrency locks.
- **Developing (1 pt)**: Cannot articulate why composite index order matters or how null scanning works.
- **Unacceptable (0 pts)**: No explanation provided.

## Grading Thresholds
- **12–14 pts**: Exceeds expectations (Gold standard)
- **10–11 pts**: Meets expectations (Passing)
- **7–9 pts**: Needs revision
- **0–6 pts**: Unsatisfactory — review prerequisites and retake
