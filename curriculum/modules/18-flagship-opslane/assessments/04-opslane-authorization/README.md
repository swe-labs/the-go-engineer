# Opslane authorization

## Mission

Implement fine-grained Role-Based Access Control (RBAC) and context-driven permission enforcement middleware for Opslane's multi-tenant APIs.

## Prerequisites

- `opslane-03`

## Mental Model

Whereas authentication determines *who* a user is, authorization enforces *what* that user is permitted to do. In Opslane, permissions are decoupled from users via roles. Roles (`owner`, `admin`, `member`, `viewer`) map to explicit granular capabilities (`order:read`, `order:create`, `order:cancel`, `tenant:manage`, `audit:view`). When an HTTP request arrives, the authentication middleware extracts the user's claims and injects an immutable `UserContext` into `r.Context()`. Downstream route decorators (`RequirePermission`) inspect this context to authorize actions before handler execution, returning standard HTTP 403 Forbidden on violations.

## Visual Model

```text
Inbound HTTP Request
  │ (Header: "Authorization: Bearer <token>")
  ▼
┌─────────────────────────────────────────────────────────────┐
│ AuthMiddleware.Authenticate                                 │
│   1. Validates Bearer token                                 │
│   2. Injects UserContext(UserID, TenantID, Role) into ctx   │
└──────────────────────────────┬──────────────────────────────┘
                               │ r.WithContext(ctx)
                               ▼
┌─────────────────────────────────────────────────────────────┐
│ RequirePermission(PermOrderCancel) Guard                    │
│   1. UserFromContext(ctx)                                   │
│   2. user.HasPermission("order:cancel")?                   │
└──────────────┬───────────────────────────────┬──────────────┘
               │ Allowed (Owner, Admin)        │ Denied (Member, Viewer)
               ▼                               ▼
┌──────────────────────────────┐ ┌────────────────────────────┐
│ Target Handler               │ │ HTTP 403 Forbidden         │
│ (e.g. POST /orders/cancel)   │ │ {"error": "forbidden", ...}│
└──────────────────────────────┘ └────────────────────────────┘
```

## Machine View

Under the hood:
1. `contextKey` is defined as a private unexported type `type contextKey string` to prevent collisions with other packages storing values in the same context tree.
2. `UserContext` provides `HasPermission`, which executes an $O(k)$ lookup across the role's permission slice.
3. `RequirePermission` acts as an HTTP handler decorator. If the user is missing from context, it emits HTTP 401 Unauthorized. If the user's role lacks the requested permission, it sets `w.WriteHeader(http.StatusForbidden)` and renders an explanatory JSON payload without leaking internal authorization logic.
4. If authorized, execution passes seamlessly to `next(w, r)` without overhead.

## Run Instructions

```bash
# Run authorization test suite
go test -v ./curriculum/modules/18-flagship-opslane/assessments/04-opslane-authorization

# Run standalone router
go run ./curriculum/modules/18-flagship-opslane/assessments/04-opslane-authorization
```

## Code Walkthrough

1. **Context Key & Roles (lines 11–24)**: Defines unexported `userContextKey` and role constants (`RoleOwner`, `RoleAdmin`, `RoleMember`, `RoleViewer`).
2. **Permissions & RBAC Matrix (lines 26–58)**: Defines granular permissions and the authoritative `RolePermissions` mapping table.
3. **`UserContext` & Helpers (lines 61–90)**:
   - `HasPermission(perm)`: Checks if the user's role contains the required permission.
   - `UserFromContext(ctx)`: Type-safe context retrieval with nil guard.
   - `ContextWithUser(ctx, u)`: Injects user identity into context.
4. **`AuthMiddleware.Authenticate` (lines 93–133)**: Extracts Bearer authorization header, resolves user identity, and updates request context.
5. **`RequirePermission` (lines 136–157)**: Enforces RBAC permissions, halting execution with 401 or 403 when checks fail.
6. **`BuildRouter` (lines 160–185)**: Declares protected routes with appropriate permission guards.

## Try It

1. Run tests with `token-viewer`: verify that `GET /api/v1/orders` succeeds (200), but `POST /api/v1/orders` fails with 403 Forbidden.
2. Run tests with `token-member`: verify that creating an order succeeds (201), but cancelling returns 403 Forbidden.
3. Run tests with `token-owner`: verify that all operations, including tenant settings management, succeed.

Common mistakes:
- Using exported built-in types (e.g. `string("user")`) as context keys, risking key collisions with third-party middleware.
- Conflating 401 Unauthorized (missing or invalid credentials) with 403 Forbidden (authenticated caller lacks permission).
- Hardcoding role checks (`if role == "admin"`) in business logic instead of checking granular permissions (`HasPermission(perm)`).

## In Production

Scalable SaaS applications avoid hardcoding role names in handlers. By gating endpoints on granular permissions (`order:cancel`) rather than roles (`admin`), systems can dynamically modify role definitions or support custom tenant-defined roles without code refactoring.

## Thinking Questions

1. Why should route handlers be guarded by *permissions* rather than *roles*?
2. What is the difference between an unexported custom type context key versus a string literal context key?
3. How does Attribute-Based Access Control (ABAC) differ from Role-Based Access Control (RBAC)?

## Next Step

Proceed to [Milestone 05 — Tenant Isolation](../05-opslane-tenant-isolation/README.md) to enforce strict tenant scoping in data access and repository layers.
