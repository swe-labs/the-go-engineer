# Project Assessment — Secure Authenticated API

## Purpose

Evaluate your ability to build a zero-compromise, security-hardened REST API in Go. This assessment verifies that you can implement cryptographic password hashing (`bcrypt`), signed JSON Web Token (JWT) lifecycle management, Role-Based Access Control (RBAC), multi-tenant boundary isolation, defensive HTTP security headers, and automated negative authorization tests.

## Targets

- `project-secure-authenticated-api`

## Scope

This assessment covers:

- Cryptographic password hashing using `golang.org/x/crypto/bcrypt` with `bcrypt.DefaultCost`
- JSON credential sanitization: preventing credential leakage using `json:"-"` struct tags
- Stateless token authentication: generating and validating HMAC-SHA256 signed JWTs (`github.com/golang-jwt/jwt/v5`)
- Authentication middleware: validating `Authorization: Bearer <token>` headers and injecting claims into `context.Context`
- Role-Based Access Control (RBAC) middleware: enforcing role hierarchy (`ADMIN` > `MEMBER` > `VIEWER`) and returning `403 Forbidden`
- Multi-tenant boundary isolation: verifying that authenticated users cannot view or delete cross-tenant resources
- Defense-in-depth HTTP security headers: `X-Content-Type-Options: nosniff`, `X-Frame-Options: DENY`, and `Content-Security-Policy`
- Comprehensive negative testing: asserting rejection of forged tokens, expired signatures, unauthenticated access, and cross-tenant tampering

## Evidence required

Submit:

1. Your completed `_starter/main.go` with password hashing, JWT generation/parsing, RBAC middleware, and tenant enforcement implemented
2. Output of `go test -v .` showing all unit and negative security tests passing
3. A written explanation answering: Walk me through your negative authorization test suite: how do your automated tests prove that Tenant A cannot access Tenant B's data even with a valid JWT?

## Time expectation

Plan for 90 to 120 minutes.

## Rubric

| Criterion | Weight | Excellent (5) | Good (3) | Needs Work (0) |
|-----------|--------|--------------|----------|----------------|
| **Correctness** | 35% | Password hashing, token issuance, authentication middleware, and RBAC hierarchy work per specification; all automated tests pass | Most endpoints work; minor role hierarchy or token parsing edge case | Authentication non-functional or tokens not verified |
| **Multi-tenant isolation** | 25% | Strict boundary verification ensures Tenant A cannot view or mutate Tenant B documents under any circumstances (403 Forbidden) | Tenant checks present but incomplete across delete or update paths | Cross-tenant data leakage occurs |
| **Security hardening** | 25% | Bcrypt cost factor applied; password hash never serialized to JSON; signature algorithm explicitly verified; security headers active on all responses | Security headers or JSON tags missing | Passwords stored in plaintext or fast hashing (MD5/SHA) used |
| **Explanation & Reasoning** | 15% | Learner articulates cryptographic differences between hashing and signatures, and explains negative test suite assertions | Explains auth with minor conceptual gaps | Unable to explain token validation or tenant isolation mechanics |

**Passing score**: 80%

## Review questions

1. Walk me through your negative authorization test suite: how do your automated tests prove that Tenant A cannot access Tenant B's data even with a valid JWT?
2. Explain the cryptographic difference between token authentication (JWT signature verification) and password storage (bcrypt work factor hashing).
3. What security headers (`CORS`, `CSP`, `X-Frame-Options`) and token expiration/refresh mechanics did you implement to mitigate token theft and replay attacks?
4. Why must passwords never be hashed with fast algorithms like SHA-256 or MD5, and why is `bcrypt`'s configurable cost factor essential?
5. What is the danger of storing sensitive claims (such as raw permissions or passwords) inside a JWT payload?

## Passing standard

You pass when all test suites pass, passwords are encrypted with bcrypt, multi-tenant boundaries are strictly enforced, and you can explain token security and RBAC mechanics.

## Retake policy

Learner may revise code and resubmit after addressing all rubric gaps.

## NEXT UP

[Module 10 Checkpoint](../checkpoint/README.md)
