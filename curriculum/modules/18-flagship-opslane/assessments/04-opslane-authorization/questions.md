# Opslane Milestone 04 Assessment — Authorization: Questions

Total: 10 questions (6 multiple choice × 1 pt, 4 short answer × 2 pt = 14 points). Passing score: 10/14 (71%).

---

## Multiple Choice (1 point each)

**1.** What is the architectural difference between HTTP 401 Unauthorized and HTTP 403 Forbidden?
- a) 401 is used for GET requests; 403 is used for POST requests.
- b) 401 means the caller has not supplied valid authentication credentials; 403 means the caller is authenticated, but lacks sufficient permissions to perform the requested operation.
- c) 401 is an internal server error; 403 is a client timeout.
- d) 401 indicates database downtime; 403 indicates rate limiting.

**2.** Why is `contextKey` defined as `type contextKey string` rather than using built-in `string` directly in `context.WithValue`?
- a) Custom types serialize faster in JSON.
- b) To prevent context key collisions across different packages or third-party middleware that might also store values using string keys.
- c) Because Go compilers refuse to compile built-in strings in `context.WithValue`.
- d) To encrypt context values in memory.

**3.** Why does `RequirePermission` inspect granular permissions (`order:cancel`) rather than checking role names (`admin`) directly?
- a) Roles cannot be stored in Go structs.
- b) To achieve loose coupling: route handlers require specific capabilities, allowing role mappings to change or custom roles to be created without modifying endpoint code.
- c) Granular permissions avoid memory allocations.
- d) HTTP routers only understand permission strings.

**4.** If an incoming request contains no `Authorization` header, what status code does `AuthMiddleware.Authenticate` return?
- a) HTTP 403 Forbidden.
- b) HTTP 401 Unauthorized.
- c) HTTP 200 OK with empty body.
- d) HTTP 500 Internal Server Error.

**5.** Which role in Opslane's RBAC matrix is permitted to manage tenant settings (`tenant:manage`)?
- a) RoleAdmin.
- b) RoleOwner.
- c) RoleMember.
- d) RoleViewer.

**6.** What permissions does `RoleViewer` hold in the `RolePermissions` matrix?
- a) `order:read` only.
- b) `order:read` and `order:create`.
- c) `order:read` and `audit:view`.
- d) All permissions in read-only mode.

---

## Short Answer (2 points each)

**7.** Explain the Principle of Least Privilege in software security and how Opslane's RBAC matrix adheres to it.

**8.** How does `UserFromContext` safely handle cases where the context does not contain a `UserContext`? Why is type assertion check (`ok`) mandatory?

**9.** Contrast Role-Based Access Control (RBAC) with Attribute-Based Access Control (ABAC). In what scenarios would Opslane require ABAC?

**10.** Explain how an attacker might attempt Confused Deputy attacks in multi-tenant authorization, and how context-driven tenant verification prevents it.
