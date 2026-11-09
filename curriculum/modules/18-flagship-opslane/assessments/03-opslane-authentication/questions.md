# Opslane Milestone 03 Assessment — Authentication: Questions

Total: 10 questions (6 multiple choice × 1 pt, 4 short answer × 2 pt = 14 points). Passing score: 10/14 (71%).

---

## Multiple Choice (1 point each)

**1.** Why is bcrypt preferred over SHA-256 with salt for storing user passwords?
- a) Bcrypt outputs shorter strings that save database disk space.
- b) Bcrypt is an adaptive, CPU- and memory-intensive key derivation function whose work factor can be increased to thwart GPU and ASIC hardware brute-force cracking.
- c) Go standard library does not support SHA-256.
- d) Bcrypt allows passwords to be reversed by system administrators.

**2.** In `TokenManager.ValidateAccessToken`, why is `if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok` explicitly checked?
- a) To ensure the token payload is written in JSON.
- b) To prevent JWT algorithm confusion attacks, where an attacker changes the `alg` header to `none` or an asymmetric algorithm to bypass signature verification.
- c) To verify that the token has not expired.
- d) To check if the user has admin privileges.

**3.** Why must refresh tokens be generated using `crypto/rand` rather than `math/rand`?
- a) `crypto/rand` is faster than `math/rand`.
- b) `math/rand` uses a deterministic pseudo-random generator with predictable seeds, allowing attackers to predict future refresh tokens; `crypto/rand` draws cryptographically secure entropy from the OS kernel.
- c) `math/rand` cannot emit strings.
- d) `crypto/rand` automatically encrypts values in the database.

**4.** What is the purpose of including `TenantID` in the JWT `Claims` struct?
- a) To bill the user's credit card per API call.
- b) To enable stateless authorization in downstream services without querying the database on every HTTP request.
- c) To enforce TLS encryption.
- d) To determine the user's preferred language.

**5.** What error should `AuthService.Login` return when a user provides an incorrect password?
- a) `ErrUserNotFound` if user doesn't exist, and `ErrIncorrectPassword` if password fails.
- b) A generic `ErrInvalidCredentials` for both cases, preventing username enumeration attacks.
- c) HTTP 500 Internal Server Error.
- d) A list of previous valid passwords.

**6.** Why is access token TTL set to a relatively short duration (e.g. 15 minutes)?
- a) Because JWT tokens grow in file size over time.
- b) To limit the window of vulnerability if an access token is stolen or intercepted in transit, relying on refresh tokens for re-authentication.
- c) Because Go timers cannot exceed 15 minutes.
- d) To clear browser cache periodically.

---

## Short Answer (2 points each)

**7.** Explain the difference between constant-time string comparison (`crypto/subtle.ConstantTimeCompare`) and standard `==` string equality checks. Why is constant-time comparison necessary in cryptographic verification?

**8.** What are the three dot-separated components of a JSON Web Token (`header.payload.signature`)? Explain what information each component contains and how the signature is computed.

**9.** How does bcrypt incorporate salts into its hash string? What do the prefix `$2a$10$` and subsequent characters represent in a bcrypt hash string?

**10.** Explain how a refresh token rotation strategy works and how it detects token theft.
