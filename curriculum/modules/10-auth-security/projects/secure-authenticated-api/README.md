# Secure Authenticated API

## Learning objective

Build a zero-compromise, security-hardened REST API with cryptographic password hashing (`bcrypt`), signed JSON Web Token (JWT) lifecycle management, Role-Based Access Control (RBAC), and strict multi-tenant boundary isolation.

## Why this matters

Authentication and authorization flaws consistently rank among the most critical OWASP vulnerabilities (Broken Object Level Authorization / IDOR, Sensitive Data Exposure, and Broken Access Control). Writing production Go services requires constant-time password comparisons, tamper-evident stateless tokens, defensive context propagation, and verified cross-tenant leak prevention.

## Prerequisites

Complete all Module 10 lessons (01–22) on security fundamentals, cryptographic hashing, JWT anatomy, authorization middleware, OWASP Top 10 defenses, and tenant isolation.

## Available concepts

This project draws on `golang.org/x/crypto/bcrypt`, `github.com/golang-jwt/jwt/v5`, HMAC-SHA256 digital signatures, HTTP middleware chains, `context.Context` claims propagation, security response headers (`X-Content-Type-Options`, `X-Frame-Options`), and negative test assertion suites.

## Project overview

The `_starter/` directory provides route layouts, in-memory models, and handler definitions. Your job is to implement cryptographic hashing, token issuance/parsing, RBAC enforcement, and tenant isolation so that:

- Passwords are encrypted with bcrypt and never serialized to JSON (`json:"-"`).
- Authenticated endpoints require a valid, unexpired Bearer JWT signed with the server secret.
- Role hierarchy (`ADMIN` > `MEMBER` > `VIEWER`) is strictly enforced:
  - `VIEWER` can read documents.
  - `MEMBER` can create and read documents.
  - `ADMIN` can delete documents.
- Multi-tenant boundaries are strictly upheld: Users from Tenant A are strictly prohibited from viewing or deleting documents belonging to Tenant B (returning 403 Forbidden).
- Security headers are applied to every response.

## RBAC Matrix

| Role | Hierarchy Level | Allowed Actions |
|---|---|---|
| `ADMIN` | 3 | Full access: create, read, delete documents |
| `MEMBER` | 2 | Create and read documents |
| `VIEWER` | 1 | Read-only access to documents |

## API Endpoints

| Method | Path | Auth Required | Required Role | Description |
|---|---|---|---|---|
| `POST` | `/auth/register` | No | None | Register new user with hashed password |
| `POST` | `/auth/login` | No | None | Authenticate with email/password; returns JWT |
| `GET` | `/auth/me` | Yes (Bearer) | Any | Returns authenticated caller claims |
| `POST` | `/documents` | Yes (Bearer) | `MEMBER` | Create a tenant-scoped document |
| `GET` | `/documents/{id}` | Yes (Bearer) | `VIEWER` | Read document (strictly scoped to caller's tenant) |
| `DELETE` | `/documents/{id}` | Yes (Bearer) | `ADMIN` | Delete document (strictly scoped to caller's tenant) |

## Tasks

### Task 1 — Cryptographic password hashing

In `_starter/main.go`:
- Implement `HashPassword(password string) (string, error)` using `bcrypt.GenerateFromPassword` with `bcrypt.DefaultCost`.
- Implement `CheckPassword(hash, password string) bool` using `bcrypt.CompareHashAndPassword`.
- Verify that `User.PasswordHash` has the `json:"-"` struct tag.

### Task 2 — JWT lifecycle & Claims extraction

- Implement `GenerateJWT(user User, secret []byte, ttl time.Duration) (string, error)`:
  - Populate `UserID`, `TenantID`, `Role`, `ExpiresAt`, and `IssuedAt`.
  - Sign using `jwt.SigningMethodHS256`.
- Implement `ParseJWT(tokenStr string, secret []byte) (*Claims, error)`:
  - Validate HMAC signing method and signature.
  - Verify expiration and subject claims.

### Task 3 — Auth & RBAC Middleware

- Implement `AuthMiddleware(secret []byte)`:
  - Extract `Authorization: Bearer <token>`.
  - Reject missing, expired, or tampered tokens with `401 Unauthorized`.
  - Inject validated `*Claims` into `r.Context()`.
- Implement `RequireRole(minRole Role)`:
  - Compare caller's role level against `minRole`.
  - Reject insufficient privilege with `403 Forbidden`.

### Task 4 — Cross-tenant isolation enforcement

In document handlers:
- Enforce that `doc.TenantID == claims.TenantID`.
- If an authenticated user from Tenant B attempts to read or delete a document owned by Tenant A, return `403 Forbidden` with zero data leakage.

## Verification

Run all automated unit and security tests:

```bash
go test -v ./curriculum/modules/10-auth-security/projects/secure-authenticated-api/_starter
```

Expected output:
```
=== RUN   TestPasswordHashingAndJsonExclusion
--- PASS: TestPasswordHashingAndJsonExclusion (0.10s)
=== RUN   TestLoginAuthentication
--- PASS: TestLoginAuthentication (0.15s)
=== RUN   TestNegativeAuthSuite
--- PASS: TestNegativeAuthSuite (0.05s)
=== RUN   TestRBACHierarchy
--- PASS: TestRBACHierarchy (0.14s)
=== RUN   TestCrossTenantIsolation
--- PASS: TestCrossTenantIsolation (0.10s)
=== RUN   TestSecurityHeaders
--- PASS: TestSecurityHeaders (0.00s)
PASS
```

## Rubric

| Criterion | Excellent (5) | Good (3) | Needs Work (0) |
|---|---|---|---|
| **Cryptographic Hygiene** | Passwords hashed with bcrypt; constant-time verification; hashes never leaked in JSON | Passwords hashed but cost < 10 or partial JSON leakage | Plaintext storage or broken hash comparison |
| **Authentication & Tokens** | JWTs properly signed, expired tokens rejected, tampered signatures caught deterministically | JWTs valid but missing expiry validation | Unverified tokens accepted |
| **RBAC Authorization** | Role hierarchy strictly enforced; insufficient role returns 403 Forbidden | Roles checked but hardcoded without hierarchy | Unauthenticated or unauthorized access allowed |
| **Tenant Boundary Defense** | Cross-tenant access strictly blocked with 403 Forbidden; zero metadata leakage | Checks present but query allows cross-tenant leaks | IDOR vulnerability allowing access across tenants |

## NEXT UP

[Bounded Worker Pool](../11-lifecycle-context-concurrency/projects/bounded-worker-pool/README.md)
