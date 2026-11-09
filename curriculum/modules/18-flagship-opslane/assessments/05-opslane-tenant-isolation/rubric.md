# Opslane Milestone 05 Assessment — Rubric

## Assessment Dimensions

### 1. Context-Bound Tenant Propagation (30%)
- **Excellent (5 pts)**: Mandatory tenant extraction from context; missing or invalid context immediately returns `ErrMissingTenantContext`; zero reliance on unverified caller parameters.
- **Proficient (3 pts)**: Tenant extracted from context, but fallback defaults allow unverified access in error cases.
- **Developing (1 pt)**: Tenant passed as an unverified function argument without context check.
- **Unacceptable (0 pts)**: No tenant scoping implemented.

### 2. Multi-Tenant Schema & Primary Key Design (25%)
- **Excellent (5 pts)**: Composite primary key `(tenant_id, id)` defined; indexes support tenant query patterns; foreign keys enforce referential integrity.
- **Proficient (3 pts)**: Schema includes `tenant_id`, but lacks composite primary key or index.
- **Developing (1 pt)**: Global single ID primary key without tenant indexing.
- **Unacceptable (0 pts)**: Schema does not partition tenant data.

### 3. Isolated Query & Mutation Enforcement (25%)
- **Excellent (5 pts)**: All CRUD operations (`Create`, `GetByID`, `List`, `Update`, `Delete`) enforce `tenant_id = ?`; mutating operations verify `RowsAffected() > 0`; not found returns clean `nil, nil`.
- **Proficient (3 pts)**: Queries filter by `tenant_id`, but `Update`/`Delete` omit `RowsAffected` validation.
- **Developing (1 pt)**: Some queries omit `tenant_id` filter, allowing cross-tenant leaks.
- **Unacceptable (0 pts)**: Operations leak or mutate cross-tenant data.

### 4. Conceptual & Defense-in-Depth Defense (20%)
- **Excellent (5 pts)**: Accurately explains IDOR attack mitigation, index clustering mechanics, Row-Level Security (RLS) integration, and context-bound design rationale.
- **Proficient (3 pts)**: Answers conceptual questions accurately with minor gaps in RLS or index clustering.
- **Developing (1 pt)**: Cannot articulate why IDOR occurs or how context prevents confused deputy attacks.
- **Unacceptable (0 pts)**: No explanation provided.

## Grading Thresholds
- **12–14 pts**: Exceeds expectations (Gold standard)
- **10–11 pts**: Meets expectations (Passing)
- **7–9 pts**: Needs revision
- **0–6 pts**: Unsatisfactory — review prerequisites and retake
