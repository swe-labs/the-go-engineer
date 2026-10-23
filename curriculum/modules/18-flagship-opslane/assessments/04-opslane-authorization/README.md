# Opslane authorization

## Mission

Understand and apply Opslane authorization in the context of professional Go software engineering.

## Prerequisites

- opslane-03

## Mental Model

Authorization is the application's security guard — authentication checks your ID, but authorization checks whether you have the right pass to enter a specific room.

## Visual Model

```text
input -> concept boundary -> behavior -> observable result
```

## Machine View

Under the hood, authorization can use Casbin (external library), a custom middleware that checks against a permission map, or OPA (Open Policy Agent). Opslane implements a lightweight middleware-based approach with a permission cache for performance.

## Run Instructions

```bash
Read the lesson and complete the practice task.
go test ./curriculum/modules/18-flagship-opslane/assessments/04-opslane-authorization
```

## Code Walkthrough

Read the code from top to bottom. Connect each line to the mental model above. The important question is not only what the line does, but why the program needs that line.

## Try It

1. Run the command exactly as shown.
2. Change one input.
3. Predict the output before running again.
4. Explain the result in your own words.

Common mistakes:

- Checking authorization only in the frontend or client-side code.
- Using hard-coded role checks instead of a configurable permission system.
- Forgetting to check authorization on every protected endpoint.

## In Production

Authorization in Opslane follows the same patterns as Kubernetes RBAC, AWS IAM, and GitHub's permission model — roles, permissions, and resource scopes.

## Concrete Proof: Negative Authorization Test Suite

To prove authorization robustness, your test suite must implement four mandatory negative test cases:
1. **Unauthenticated Request:** Calling protected endpoints without a Bearer token returns `401 Unauthorized` with structured JSON error `{ "error": "unauthorized" }`.
2. **Role Mismatch (RBAC):** Authenticated user with role `viewer` attempting `POST /api/v1/orders` returns `403 Forbidden` with `{ "error": "insufficient_permissions" }`.
3. **Cross-Tenant Boundary Rejection:** User belonging to Tenant A querying `/api/v1/tenants/tenant-B/resources` returns `404 Not Found` (preferred to avoid tenant enumeration) or `403 Forbidden`.
4. **Token Tampering / Invalid Signature:** Request with altered JWT payload returns immediate `401 Unauthorized`.

## Thinking Questions

1. What problem does this concept solve?
2. What would break if this concept were removed?
3. How would you test the behavior?

## Next Step

Continue to the next item listed in metadata: `opslane-05`.
