# Opslane Milestone 03 Assessment — Rubric

## Assessment Dimensions

### 1. Password Security & Bcrypt Engineering (30%)
- **Excellent (5 pts)**: Passwords hashed with salted bcrypt; cost factor parameterized; verification uses `CompareHashAndPassword`; generic errors prevent username enumeration; hashes never logged or serialized.
- **Proficient (3 pts)**: Bcrypt implemented correctly, but uses hardcoded cost or distinct error messages that leak user existence.
- **Developing (1 pt)**: Uses fast hashing algorithms (MD5, SHA-256) or stores plain passwords.
- **Unacceptable (0 pts)**: No password security implemented.

### 2. Token Issuance & Algorithm Defense (30%)
- **Excellent (5 pts)**: JWT claims include tenant and user identities with registered expiration; HMAC signing algorithm strictly validated against key confusion; refresh tokens use `crypto/rand`.
- **Proficient (3 pts)**: JWT issued and verified, but algorithm check omitted or refresh tokens use `math/rand`.
- **Developing (1 pt)**: Tokens issued without expiration or signature verification fails.
- **Unacceptable (0 pts)**: Tokens are malformed or fail verification.

### 3. Authentication Architecture & Session Lifecycle (20%)
- **Excellent (5 pts)**: Separation of concerns across `PasswordManager`, `TokenManager`, and `AuthService`; token expiry returns explicit sentinel `ErrTokenExpired`; short TTL limits blast radius.
- **Proficient (3 pts)**: Workflow functions correctly, but error handling conflates expiration with invalid signatures.
- **Developing (1 pt)**: Monolithic auth function with poor separation of hashing and token logic.
- **Unacceptable (0 pts)**: Auth service does not compile or pass tests.

### 4. Conceptual & Threat Modeling Defense (20%)
- **Excellent (5 pts)**: Articulates timing attack mitigations, JWT structure, bcrypt string syntax, algorithm confusion defense, and refresh token rotation mechanics.
- **Proficient (3 pts)**: Answers conceptual questions accurately with minor gaps in timing attacks or token anatomy.
- **Developing (1 pt)**: Cannot explain why constant-time comparison is necessary or what bcrypt cost means.
- **Unacceptable (0 pts)**: No explanation provided.

## Grading Thresholds
- **12–14 pts**: Exceeds expectations (Gold standard)
- **10–11 pts**: Meets expectations (Passing)
- **7–9 pts**: Needs revision
- **0–6 pts**: Unsatisfactory — review prerequisites and retake
