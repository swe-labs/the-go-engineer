# Module 10 Assessment — Authentication, Authorization, and Security

## Purpose

Verify that the learner can independently apply and explain JWT, bcrypt hashing, OAuth2 flows, API keys, TLS configuration, refresh token rotation, RBAC, and security headers.

## Targets

- `core-10-01`
- `core-10-02`
- `core-10-03`
- `core-10-04`
- `core-10-05`
- `core-10-06`
- `core-10-07`
- `core-10-08`
- `core-10-09`
- `core-10-10`
- `core-10-11`
- `core-10-12`
- `core-10-13`
- `core-10-14`
- `core-10-15`
- `core-10-16`
- `core-10-17`
- `core-10-18`
- `core-10-19`
- `core-10-20`
- `core-10-21`
- `core-10-22`
- `core-10-23`
- `core-10-24`
- `core-10-25`

## Scope

This checkpoint covers:

- Cryptographic foundations: one-way hashing vs symmetric/asymmetric encryption, salt generation, and slow hashing algorithms (`bcrypt`, `argon2id`)
- Password management: constant-time comparison (`subtle.ConstantTimeCompare`), preventing timing attacks, and credential leak prevention in serialization
- Token-based authentication: JWT header, payload claims, cryptographic signature verification, expiration enforcement, and secret rotation
- Authorization models: Role-Based Access Control (RBAC), Attribute-Based Access Control (ABAC), and multi-tenant isolation
- Security middleware: extracting Bearer credentials, injecting claims into request context, rejecting unauthenticated and unauthorized requests
- Web security defenses: Cross-Origin Resource Sharing (`CORS`), Content Security Policy (`CSP`), `X-Frame-Options`, `X-Content-Type-Options`, and Strict-Transport-Security (`HSTS`)
- Transport Layer Security: TLS 1.3 configuration, certificates, and mutual TLS (mTLS) fundamentals
- OAuth2 and OpenID Connect flows: authorization code grant with PKCE, client credentials, and token revocation
- OWASP Top 10 mitigation: SQL injection, broken object-level authorization (IDOR), security misconfiguration, and credential stuffing defenses

## Evidence required

Submit:

1. **Authentication middleware**: Composable middleware verifying Bearer tokens and injecting validated claims into request contexts
2. **Cryptographic hashing & token issuance**: Implementation demonstrating bcrypt password hashing and HMAC-SHA256 JWT generation
3. **Security-hardened endpoint suite**: Endpoints enforcing RBAC and multi-tenant isolation accompanied by automated negative tests
4. **Written explanations**: Answers to the five review questions below

## Time expectation

Plan for 2 to 3 hours.

## Rubric

| Criterion | Weight | Excellent (5) | Good (3) | Needs Work (0) |
|-----------|--------|--------------|----------|----------------|
| **Conceptual clarity** | 30% | Clearly explains difference between authentication and authorization, cryptographic hashing mechanics, JWT verification, and RBAC schemas | Explains basics of auth, tokens, and roles with minor conceptual gaps | Cannot distinguish authentication from authorization or explain token validation |
| **Implementation** | 35% | Complete, clean implementation with bcrypt hashing, JWT validation, RBAC enforcement, and defense-in-depth security headers | Working implementation with minor role or header configuration flaws | Plaintext passwords or unverified tokens accepted |
| **Proof** | 25% | Strong negative test suite verifying rejection of tampered tokens, expired credentials, unauthorized roles, and cross-tenant attempts | Basic tests covering happy path only | No test evidence or verification |
| **Communication** | 10% | Professional explanation connecting security patterns to OWASP Top 10 mitigation, token rotation, and audit compliance | Basic explanation with some clarity | Unable to articulate security trade-offs or threat models |

**Passing score**: 80%

## Review questions

1. Explain the difference between authentication and authorization. How does JWT prove identity without a database lookup?
2. Walk me through RBAC. How would you model roles, permissions, and resource access in a SQL schema?
3. Show me a SQL injection attack and explain exactly why parameterized queries prevent it at the driver level.
4. What is a timing attack, and how does Go's `crypto/subtle` package prevent timing side-channels during credential verification?
5. What is the security advantage of short-lived access tokens paired with rotating refresh tokens stored in HTTP-only secure cookies?

## Passing standard

Score of 80% or higher across the four rubric criteria, with all required evidence submitted.

## Retake policy

Learner may retry after fixing rubric gaps and writing a short correction note.

## NEXT UP

[Module 11: Lifecycle, Context, and Concurrency](../../11-lifecycle-context-concurrency/README.md)
