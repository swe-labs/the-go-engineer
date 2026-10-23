# Opslane Milestone 04 Assessment — Rubric

## Assessment Dimensions

### 1. RBAC Matrix & Capability Mapping (30%)
- **Excellent (5 pts)**: Fine-grained permissions separated from roles; authoritative mapping matrix implemented; `HasPermission` evaluates membership correctly; Least Privilege enforced across all tiers.
- **Proficient (3 pts)**: Permissions defined, but roles check strings directly rather than using typed constants.
- **Developing (1 pt)**: Roles check broad permissions without granular separation.
- **Unacceptable (0 pts)**: No permission matrix defined.

### 2. Context Safety & Key Isolation (25%)
- **Excellent (5 pts)**: Unexported custom type used for context key to prevent collision; `UserFromContext` performs safe type assertion with nil guard; context is immutable.
- **Proficient (3 pts)**: Context used correctly, but uses exported or built-in string key.
- **Developing (1 pt)**: Context type assertions made without comma-ok check, risking runtime panics.
- **Unacceptable (0 pts)**: Context handling crashes on missing values.

### 3. Middleware Pipeline & Status Code Semantics (25%)
- **Excellent (5 pts)**: Strict distinction between 401 Unauthorized (missing/bad token) and 403 Forbidden (insufficient permissions); JSON error payloads descriptive; clean handler wrapping.
- **Proficient (3 pts)**: Middleware works, but returns 401 for permission denials or 403 for missing tokens.
- **Developing (1 pt)**: Middleware fails to pass context to `next.ServeHTTP`.
- **Unacceptable (0 pts)**: Middleware does not intercept requests.

### 4. Conceptual & Access Control Defense (20%)
- **Excellent (5 pts)**: Accurately explains Least Privilege, type-safe context retrieval, RBAC vs ABAC trade-offs, and confused deputy attack mitigation.
- **Proficient (3 pts)**: Answers conceptual questions accurately with minor gaps in ABAC or context internals.
- **Developing (1 pt)**: Confuses 401 and 403 or cannot explain context collision risks.
- **Unacceptable (0 pts)**: No explanation provided.

## Grading Thresholds
- **12–14 pts**: Exceeds expectations (Gold standard)
- **10–11 pts**: Meets expectations (Passing)
- **7–9 pts**: Needs revision
- **0–6 pts**: Unsatisfactory — review prerequisites and retake
