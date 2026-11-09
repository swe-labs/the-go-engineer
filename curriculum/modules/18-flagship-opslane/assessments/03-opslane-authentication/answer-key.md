# Opslane Milestone 03 Assessment — Answer Key

## Multiple Choice (1 point each)

**1. b) Bcrypt is an adaptive, CPU- and memory-intensive key derivation function whose work factor can be increased to thwart GPU and ASIC hardware brute-force cracking.**  
*Explanation:* SHA-256 is designed to be as fast as possible for integrity checks (billions of hashes/sec on GPUs). Bcrypt is deliberately slow and memory-hard. Its work factor ($2^{\text{cost}}$ iterations) can be adjusted upward as hardware accelerates, making brute-force dictionary attacks infeasible.

**2. b) To prevent JWT algorithm confusion attacks, where an attacker changes the `alg` header to `none` or an asymmetric algorithm to bypass signature verification.**  
*Explanation:* In algorithm confusion attacks, an attacker modifies the token header to `{"alg": "none"}` or changes RS256 to HS256 (signing with the server's public key as an HMAC secret). Explicitly checking `t.Method.(*jwt.SigningMethodHMAC)` ensures that only expected symmetric HMAC algorithms are accepted.

**3. b) `math/rand` uses a deterministic pseudo-random generator with predictable seeds, allowing attackers to predict future refresh tokens; `crypto/rand` draws cryptographically secure entropy from the OS kernel.**  
*Explanation:* `math/rand` is a linear congruential generator designed for simulation speed. If an attacker observes a few outputs, they can reconstruct the internal state and predict subsequent tokens. `crypto/rand` reads from `/dev/urandom` or Windows CryptGenRandom, providing true cryptographic unpredictability.

**4. b) To enable stateless authorization in downstream services without querying the database on every HTTP request.**  
*Explanation:* Storing `tenant_id` and `role` inside the signed JWT claims allows API gateways and downstream microservices to verify the token signature and enforce tenant boundaries immediately, eliminating database lookup latency on every request.

**5. b) A generic `ErrInvalidCredentials` for both cases, preventing username enumeration attacks.**  
*Explanation:* If an API returns "user not found" when an email does not exist and "invalid password" when the email exists, an attacker can enumerate valid corporate email addresses by analyzing server error responses. Returning a generic error masks user existence.

**6. b) To limit the window of vulnerability if an access token is stolen or intercepted in transit, relying on refresh tokens for re-authentication.**  
*Explanation:* Stateless JWTs cannot be easily revoked before expiration without maintaining a revocation list. A short lifespan (10–15 min) ensures that if an access token leaks, the attacker's window of opportunity expires quickly.

---

## Short Answer (2 points each)

**7. Constant-time comparison vs. standard string equality:**
- **Standard `==`:** Compares strings byte-by-byte and terminates early on the first mismatched byte. An attacker can measure nanosecond differences in response times to guess valid secrets character-by-character (Timing Attack).
- **Constant-Time Comparison:** Evaluates all bytes regardless of where differences occur, taking constant execution time and preventing timing side-channel attacks.
*Grading: 1 pt for explaining early termination vs timing side-channels; 1 pt for explaining constant-time execution.*

**8. Anatomy of a JSON Web Token (JWT):**
1. **Header (Base64URL):** Specifies token metadata, including token type (`typ: "JWT"`) and signing algorithm (`alg: "HS256"`).
2. **Payload (Base64URL):** Contains registered claims (`exp`, `iat`, `iss`, `sub`) and custom business claims (`tenant_id`, `role`, `email`).
3. **Signature:** Computed as `HMAC_SHA256(Base64URL(Header) + "." + Base64URL(Payload), SecretKey)`. Verifies that the header and payload have not been tampered with.
*Grading: 1 pt for defining the three components; 1 pt for explaining how the signature is calculated.*

**9. Bcrypt hash format:**
A bcrypt hash string (e.g. `$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy`) contains:
- `$2a$`: The bcrypt specification algorithm identifier.
- `$10$`: The logarithmic cost factor ($2^{10} = 1024$ key derivation rounds).
- Next 22 characters: The 128-bit salt encoded in base-64.
- Remaining 31 characters: The resulting key/hash value.
*Grading: 1 pt for identifying algorithm and cost factor ($2^{\text{cost}}$); 1 pt for explaining the embedded salt and hash.*

**10. Refresh token rotation and theft detection:**
- **Mechanism:** Whenever a client exchanges a refresh token for a new access token, the server invalidates the used refresh token and issues a brand new refresh token in response.
- **Theft Detection:** If an attacker and a legitimate client both attempt to use the same refresh token, the server detects that an already-used (revoked) refresh token was presented. The server treats this as a breach and immediately revokes all active refresh tokens for that user session family, forcing all parties to re-authenticate.
*Grading: 1 pt for explaining one-time use token exchange; 1 pt for explaining breach detection and session family revocation.*
