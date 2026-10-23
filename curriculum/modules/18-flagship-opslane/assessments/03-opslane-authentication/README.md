# Opslane authentication

## Mission

Implement enterprise-grade authentication for Opslane, integrating salted cryptographic password hashing with `bcrypt`, stateless JSON Web Token (JWT) issuance and validation with `golang-jwt/jwt/v5`, cryptographically random refresh token generation, and tenant-scoped identity claims.

## Prerequisites

- `opslane-02`

## Mental Model

Authentication verifies *who* the caller is before allowing access to tenant resources. In Opslane, authentication combines slow, adaptive one-way hashing for credentials at rest (`bcrypt` with configurable work factor) and fast, stateless digital signatures for API requests (`HMAC-SHA256` signed JWTs). When a user presents credentials to `AuthService.Login()`, the server verifies the password against the stored hash in constant time, issues a short-lived access token carrying tenant and user identity claims, and generates a cryptographically random refresh token for session lifecycle management.

## Visual Model

```text
┌─────────────────────────────────────────────────────────────┐
│ Client Login Request (email, password)                      │
└──────────────────────────────┬──────────────────────────────┘
                               │ POST /api/v1/auth/login
                               ▼
┌─────────────────────────────────────────────────────────────┐
│ AuthService.Login                                           │
│   1. Retrieve user by email                                 │
│   2. bcrypt.CompareHashAndPassword(hash, password)          │
└──────────────────────────────┬──────────────────────────────┘
                               │ If password matches
                               ▼
┌─────────────────────────────────────────────────────────────┐
│ TokenManager Issuance                                       │
│   ├─► Claims: TenantID, UserID, Role, Email, ExpiresAt(15m) │
│   ├─► Signs Access Token via HMAC-SHA256 (HS256)            │
│   └─► Generates 32-byte Refresh Token via crypto/rand       │
└──────────────────────────────┬──────────────────────────────┘
                               │
                               ▼
┌─────────────────────────────────────────────────────────────┐
│ AuthResponse: AccessToken, RefreshToken, ExpiresAt          │
└─────────────────────────────────────────────────────────────┘
```

## Machine View

Under the hood:
1. `PasswordManager.Hash` calls `bcrypt.GenerateFromPassword`. Bcrypt incorporates an automatic 128-bit salt and an exponential cost parameter ($2^{\text{cost}}$ iterations), rendering rainbow table and GPU-accelerated dictionary attacks computationally infeasible.
2. `TokenManager.GenerateAccessToken` constructs a JWT payload with standard registered claims (`exp`, `iat`, `iss`, `sub`) plus custom tenant claims (`tenant_id`, `role`). It signs the payload using `jwt.SigningMethodHS256`.
3. `TokenManager.ValidateAccessToken` inspects the token's header to explicitly verify that the signing algorithm is HMAC (`*jwt.SigningMethodHMAC`). This defends against critical JWT algorithm substitution vulnerabilities (e.g. attackers switching `alg` to `none` or an asymmetric public key).
4. `GenerateRefreshToken` generates 32 bytes of cryptographically secure pseudo-random entropy from the OS kernel via `crypto/rand.Read`, encoded to a 64-character hexadecimal string.

## Run Instructions

```bash
# Run authentication test suite
go test -v ./curriculum/modules/18-flagship-opslane/assessments/03-opslane-authentication

# Run standalone login demo
go run ./curriculum/modules/18-flagship-opslane/assessments/03-opslane-authentication
```

## Code Walkthrough

1. **Domain & Claims Structs (lines 21–46)**: Declares `UserAuth`, `Claims` (embedding `jwt.RegisteredClaims` with `TenantID`, `UserID`, `Role`), and `AuthResponse`.
2. **`PasswordManager` (lines 48–71)**: Encapsulates `bcrypt.GenerateFromPassword` and `bcrypt.CompareHashAndPassword`. Defaults to `bcrypt.DefaultCost` (10), while allowing lower cost in test environments for fast test execution.
3. **`TokenManager` (lines 73–142)**:
   - `GenerateAccessToken`: Builds `Claims` with 15-minute expiration and signs with symmetric secret key.
   - `GenerateRefreshToken`: Pulls 32 bytes from `crypto/rand`.
   - `ValidateAccessToken`: Parses token string, enforces HMAC signing method check, verifies signature, and checks expiration.
4. **`AuthService` (lines 144–200)**: Coordinates credential lookup, bcrypt password verification, and token issuance returning `*AuthResponse`.

## Try It

1. Run the test suite: verify that valid credentials issue a token and invalid credentials return `ErrInvalidCredentials`.
2. Test token tampering: alter a single byte in the JWT signature and verify `ValidateAccessToken` rejects it with `ErrInvalidToken`.
3. Test expiration: generate a token with a negative TTL and verify it returns `ErrTokenExpired`.

Common mistakes:
- Using fast cryptographic hashes (MD5, SHA-256) for passwords instead of slow adaptive hashes like bcrypt.
- Skipping the algorithm validation check in `jwt.ParseWithClaims`, leaving the service vulnerable to key confusion exploits.
- Using `math/rand` instead of `crypto/rand` for refresh tokens and secret generation.

## In Production

In production cloud environments, JWT signing keys should be rotated periodically or backed by asymmetric key pairs (RS256 / EdDSA) with JWKS endpoints. Short access token lifespans (10–15 minutes) limit the blast radius of stolen tokens, while refresh tokens stored securely in HttpOnly cookies allow seamless session renewal without re-prompting users for credentials.

## Thinking Questions

1. Why is MD5 or SHA-256 unsuitable for storing passwords, even when salted?
2. What is the JWT algorithm confusion attack (`alg: none` or RSA public key confusion) and how does our parser prevent it?
3. What are the security tradeoffs between stateful session IDs in Redis versus stateless JWT tokens?

## Next Step

Proceed to [Milestone 04 — Authorization](../04-opslane-authorization/README.md) to implement Role-Based Access Control (RBAC) and tenant authorization middleware.
