# API keys

## Learning objective

Generate, hash, validate, and rotate API keys in Go using cryptographically secure randomness, store only hashed keys in the database, and scope keys to specific permissions.

## Why this matters

API keys are the primary authentication mechanism for developer-facing APIs. Stripe, Twilio, GitHub, and OpenAI all use API keys. A leaked API key can lead to data breaches, account takeover, and financial fraud. The engineering challenge is balancing usability (a single key string that never changes) with security (keys that can be revoked and rotated). In Go services, a well-designed API key system includes cryptographically random generation, server-side hashing, permission scoping, and rotation support.

## Mental model

An API key is a shared secret between the service and the client. Think of it as a long, random password that a program uses instead of a human. The service gives the client a key once (at creation time) and stores only a hash of it. If the database is breached, the attacker gets hashes, not usable keys. If a key leaks, the service rotates the key: the old hash is replaced with a new one. Keys can be scoped to specific permissions: a read-only key cannot write data, even if it is valid.

## Core idea

An API key system has five parts:

1. **Generation**: produce a cryptographically random string with a prefix for identification (e.g., `api_key_abc123...`).
2. **Storage**: store only a SHA-256 hash of the key in the database. The plaintext key is returned once at creation.
3. **Validation**: on each request, hash the presented key and look up the hash in the database.
4. **Scoping**: attach permissions (read, write, admin) to each key.
5. **Rotation**: allow replacing a key without changing the scoping or metadata.

Key format best practice: `prefix + random_bytes(base62)`. Example: `api_key_6HzmLMNPqA9rBx8Cw4fDgE2hJk5nV7` (prefix = "api_key_", random = 22 alphanumeric chars).

| Component | Size | Example |
|---|---|---|
| Prefix (optional) | 5-15 chars | `api_key_`, `svc_key_` |
| Random payload | 32 bytes -> ~44 base64 chars | `6HzmLMNPqA9rBx8Cw4fDgE2hJk5nV7` |
| Total | ~50-60 chars | `api_key_6HzmLMNPqA9rBx8Cw4fDgE2hJk5nV7` |

## Under the hood

Key generation uses `crypto/rand` to produce 32+ bytes of entropy, then encodes them with base62 or base64url. The prefix serves two purposes: it identifies the key type (live/test, service, environment) and makes keys recognizable in logs (without exposing the full key).

Storage uses SHA-256 hashing (not bcrypt — API keys are high-entropy random strings that do not benefit from slow hashing). SHA-256 is appropriate because:

- API keys have 128+ bits of entropy (cannot be brute-forced).
- Fast hashing means low latency on validation.
- Each key is unique (no salt needed, though HMAC with a pepper is a defense-in-depth option).

Validation flow:

```
Request -> Extract key from header -> SHA-256(key) -> DB lookup by hash -> Load scopes -> Authorize
```

## How Go uses it

The Go standard library provides everything needed for API key systems:

- `crypto/rand`: cryptographic random bytes.
- `crypto/sha256`: fast hashing for storage.
- `encoding/base64`: encoding random bytes to strings.

Third-party packages like `github.com/gorilla/mux` or `chi` provide middleware patterns for key validation.

## Go example

```go
package main

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

// APIKey represents a stored API key record.
type APIKey struct {
	ID          string    `json:"id"`
	KeyPrefix   string    `json:"key_prefix"`
	KeyHash     string    `json:"key_hash"`  // SHA-256 of the full key
	Description string    `json:"description"`
	Scopes      []string  `json:"scopes"`
	CreatedAt   time.Time `json:"created_at"`
	LastUsedAt  time.Time `json:"last_used_at,omitempty"`
	ExpiresAt   time.Time `json:"expires_at,omitempty"`
}

// KeyStore is a thread-safe API key store (simulates a database).
type KeyStore struct {
	mu       sync.RWMutex
	keysByID map[string]*APIKey
	keysByHash map[string]*APIKey
}

func NewKeyStore() *KeyStore {
	return &KeyStore{
		keysByID:   make(map[string]*APIKey),
		keysByHash: make(map[string]*APIKey),
	}
}

func (ks *KeyStore) Add(key *APIKey) {
	ks.mu.Lock()
	defer ks.mu.Unlock()
	ks.keysByID[key.ID] = key
	ks.keysByHash[key.KeyHash] = key
}

func (ks *KeyStore) GetByHash(hash string) (*APIKey, bool) {
	ks.mu.RLock()
	defer ks.mu.RUnlock()
	k, ok := ks.keysByHash[hash]
	return k, ok
}

func (ks *KeyStore) Delete(id string) {
	ks.mu.Lock()
	defer ks.mu.Unlock()
	if k, ok := ks.keysByID[id]; ok {
		delete(ks.keysByID, id)
		delete(ks.keysByHash, k.KeyHash)
	}
}

// GenerateAPIKey creates a new API key with a prefix and scopes.
// Returns the full plaintext key (shown once) and the stored record.
func GenerateAPIKey(prefix string, scopes []string, ttl time.Duration) (string, *APIKey, error) {
	if prefix == "" {
		prefix = "sk"
	}
	// Generate 32 random bytes.
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", nil, fmt.Errorf("random: %w", err)
	}
	// Encode as base62-like using base64url (no padding).
	randomPart := base64.RawURLEncoding.EncodeToString(b)

	// Full key with prefix. Example: "api_key_abc123..."
	fullKey := prefix + "_" + randomPart

	// Hash the full key for storage.
	hash := HashKey(fullKey)

	now := time.Now()
	record := &APIKey{
		ID:          generateID(),
		KeyPrefix:   prefix + "_",
		KeyHash:     hash,
		Description: "Auto-generated key",
		Scopes:      scopes,
		CreatedAt:   now,
		ExpiresAt:   now.Add(ttl),
	}
	return fullKey, record, nil
}

// HashKey returns the SHA-256 hex hash of a key.
func HashKey(key string) string {
	h := sha256.Sum256([]byte(key))
	return hex.EncodeToString(h[:])
}

// ValidateKey checks a presented key against the store.
func ValidateKey(store *KeyStore, presentedKey string, requiredScope string) (*APIKey, error) {
	hash := HashKey(presentedKey)
	record, ok := store.GetByHash(hash)
	if !ok {
		return nil, errors.New("invalid key")
	}
	if !record.ExpiresAt.IsZero() && time.Now().After(record.ExpiresAt) {
		return nil, errors.New("key expired")
	}
	if requiredScope != "" && !hasScope(record.Scopes, requiredScope) {
		return nil, fmt.Errorf("insufficient scope: need %s", requiredScope)
	}
	return record, nil
}

func hasScope(scopes []string, required string) bool {
	for _, s := range scopes {
		if s == required || s == "admin" {
			return true
		}
		if strings.HasSuffix(s, ":*") {
			prefix := strings.TrimSuffix(s, ":*")
			if strings.HasPrefix(required, prefix+":") {
				return true
			}
		}
	}
	return false
}

func generateID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

// RotateAPIKey creates a new key with the same metadata as an existing one.
func RotateAPIKey(store *KeyStore, existingID, prefix string) (string, *APIKey, error) {
	store.mu.RLock()
	existing, ok := store.keysByID[existingID]
	store.mu.RUnlock()
	if !ok {
		return "", nil, errors.New("key not found")
	}
	// Generate new key with same scopes.
	fullKey, record, err := GenerateAPIKey(prefix, existing.Scopes, time.Until(existing.ExpiresAt))
	if err != nil {
		return "", nil, err
	}
	// Delete old key.
	store.Delete(existingID)
	// Add new key.
	store.Add(record)
	return fullKey, record, nil
}

// APIKeyMiddleware validates API keys on protected routes.
func APIKeyMiddleware(store *KeyStore, requiredScope string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key := r.Header.Get("X-API-Key")
			if key == "" {
				key = r.URL.Query().Get("api_key") // fallback (less secure)
			}
			if key == "" {
				http.Error(w, `{"error":"missing API key"}`, http.StatusUnauthorized)
				return
			}
			_, err := ValidateKey(store, key, requiredScope)
			if err != nil {
				http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func main() {
	store := NewKeyStore()

	// Create keys with different scopes.
	key1, record1, _ := GenerateAPIKey("api_key", []string{"documents:read", "documents:write"}, 365*24*time.Hour)
	store.Add(record1)
	fmt.Printf("Created key 1: %s\n  Scopes: %v\n  Hash (stored): %s...\n\n", key1, record1.Scopes, record1.KeyHash[:16])

	key2, record2, _ := GenerateAPIKey("api_key", []string{"documents:read"}, 365*24*time.Hour)
	store.Add(record2)
	fmt.Printf("Created key 2 (read-only): %s\n  Scopes: %v\n\n", key2, record2.Scopes)

	// Simulate reading key 1 from logs for testing.
	_ = key1

	// Set up HTTP routes.
	mux := http.NewServeMux()
	mux.Handle("/api/documents", APIKeyMiddleware(store, "documents:read")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "read access granted"})
	})))
	mux.Handle("/api/documents/write", APIKeyMiddleware(store, "documents:write")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "write access granted"})
	})))

	log.Println("Server on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

// Ensure the hmac import is used.
var _ = hmac.Equal
```

## Step-by-step execution

For key creation with `GenerateAPIKey("api_key", ["documents:read", "documents:write"], 365*24*h)`:

1. `crypto/rand` generates 32 random bytes (256 bits of entropy).
2. Base64url encoding produces a 43-character string.
3. Full key: `api_key_<43-char-random>`.
4. `SHA-256(full_key)` produces a 64-character hex hash.
5. Record stored with hash, scopes, timestamps.
6. Plaintext key returned to user exactly once.
7. User stores the key securely.

For a request with `X-API-Key: api_key_<...>`:

1. Middleware extracts the key from the header.
2. `SHA-256(presented_key)` produces the hash.
3. Database lookup by hash returns the record (or nil).
4. Check: record exists? Expired? Has required scope?
5. If all checks pass: request proceeds. If any fails: 403 Forbidden.

## Common mistakes

- Using API keys as the sole authentication mechanism for human users — API keys lack MFA, granular revocation, and rotation ergonomics.
- Logging API keys in request logs or error messages — keys in logs are accessible to anyone with log access.
- Storing API keys in plaintext in the database — a database breach exposes all keys.
- Generating predictable API keys — sequential or time-based keys are guessable. Always use `crypto/rand`.
- Using bcrypt to hash API keys — bcrypt is for low-entropy passwords. API keys have 256 bits of entropy; SHA-256 is sufficient and faster.
- Not supporting key rotation — once a key is exposed, the only mitigation is a new key with the same permissions.

## Debugging walkthrough

Consider this API key validation code:

```go
func validateAPIKey(r *http.Request, db *sql.DB) (*APIKey, error) {
    key := r.Header.Get("X-API-Key")
    row := db.QueryRow("SELECT id, scopes FROM api_keys WHERE key_hash = $1", key)
    // BUG: storing the raw key, not the hash!
    // ...
}
```

**Symptom**: API keys work when first created, but stop working after a few hours.

**Investigation**: The database stores the raw key in the `key_hash` column. The query looks up the raw key directly. This works until the key is rotated or the database is refreshed from a backup.

**Root cause**: The developer confused the raw API key with its hash. The API key is never stored — only its hash should be.

**Fix**:

```go
func validateAPIKey(r *http.Request, db *sql.DB) (*APIKey, error) {
    key := r.Header.Get("X-API-Key")
    hash := sha256Hex(key)
    row := db.QueryRow("SELECT id, scopes FROM api_keys WHERE key_hash = $1", hash)
    // ...
}
```

## Production notes

- Keys should never be returned after creation. Show them once and provide a "regenerate" option.
- Use HMAC-SHA256 with a server-side pepper for additional protection: `hash = HMAC-SHA256(pepper, key)`. This prevents key hashes from being salted offline if the database and code are both stolen.
- Implement key listing with masking: `api_key_...aB3x` shows only the last 4 characters.
- Rate-limit key validation attempts to prevent brute-force of leaked key hashes.
- Support multiple active keys per user for zero-downtime rotation: the old key and new key are both valid during the rotation window.

## Performance implications

- SHA-256 hashing of API keys: ~0.001ms per hash. Negligible.
- Database lookup by hash: O(log n) with an index on `key_hash`. Sub-millisecond for millions of keys.
- Scope check: O(s) where s is the number of scopes per key (typically < 20).
- Total validation overhead: < 2ms per request, dominated by database query latency.

## Practice task

Write Go functions `CreateAPIKey(prefix string, scopes []string) (fullKey string, hash string, err error)` and `ValidateAPIKey(store map[string]*APIKey, presentedKey, requiredScope string) error` where:
- `CreateAPIKey` generates a key with `crypto/rand`, SHA-256 hashes it, stores the hash, and returns the full key.
- `ValidateAPIKey` hashes the presented key, looks up the hash, checks scope and expiry.
- Then write a `main()` that creates 3 keys (read-only, write, admin), validates each against the appropriate scope, and tests a revoked key.

## Tests / verification

```bash
go run ./curriculum/modules/10-auth-security/lessons/10-api-keys
go test ./curriculum/modules/10-auth-security/lessons/10-api-keys
```

The test file `main_test.go` contains table-driven tests that verify:
- `CreateAPIKey` returns a key matching the prefix pattern.
- `CreateAPIKey` returns a non-empty hash.
- `ValidateAPIKey` returns nil for valid key and scope.
- `ValidateAPIKey` returns error for invalid key.
- `ValidateAPIKey` returns error for wrong scope.

## Review questions

1. Why should API keys be hashed with SHA-256 rather than stored in plaintext?
2. Why is bcrypt not the right choice for hashing API keys (unlike passwords)?
3. What is the purpose of a key prefix like `api_key_`?
4. How would you implement key rotation without downtime?
5. An API key validation middleware runs before every request. What could go wrong if the database lookup is slow?

## NEXT UP

RBAC — role-based access control with roles, permissions, and middleware enforcement in Go.
