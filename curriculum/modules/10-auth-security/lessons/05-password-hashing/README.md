# Password hashing

## Learning objective

Hash and verify passwords using bcrypt from `golang.org/x/crypto`, configure cost factors appropriately, and explain why SHA-256, AES, and custom hash functions are unsafe for password storage.

## Why this matters

Every week, another database breach leaks millions of passwords. If you store passwords in plaintext or with a fast hash, every leaked password is instantly cracked. In 2022 alone, over 24 billion credentials were exposed in data breaches. As a Go engineer, you are responsible for user credentials. Bcrypt, scrypt, and Argon2 are the only acceptable choices for password hashing. There is no excuse for storing a plaintext password or using `sha256(password)` — these are real vulnerabilities that lead to real account takeovers.

## Mental model

Password hashing is a meat grinder: you put in a password, and what comes out is unrecognizable ground meat. You cannot reconstruct the original cut of meat from the ground meat. The grinder has a "difficulty dial" (cost factor): turning it up makes the grinding slower. A fast grinder (SHA-256) produces ground meat in microseconds — an attacker with a GPU can try billions of passwords per second. A slow grinder (bcrypt with cost 12) takes 250ms per password — the attacker can only try 4 passwords per second. The grinder also adds a random pinch of salt so that grinding the same password twice produces different results.

## Core idea

Three properties of a secure password hash:

1. **One-way**: given the hash, it is computationally infeasible to recover the original password.
2. **Slow**: the hash function is deliberately expensive to compute, making brute-force attacks impractical.
3. **Salted**: a random value (salt) is added to each password before hashing, so identical passwords produce different hashes.

| Algorithm | Salt | Memory-hard | Recommended cost | GPU-resistant |
|---|---|---|---|---|
| bcrypt | 16 bytes, embedded in hash | No | Cost 10-14 | Moderate (sequential algorithm) |
| scrypt | 32 bytes, embedded in hash | Yes | N=2^15, r=8, p=1 | Strong (requires memory) |
| Argon2id | 16 bytes, embedded in hash | Yes | time=1, mem=64MB, threads=4 | Strongest (winner of PHC) |

bcrypt is the most widely supported and is the default for most Go applications. Argon2id is the modern choice for new systems.

## Under the hood

bcrypt is based on the Blowfish cipher. The algorithm:

1. Generates a 16-byte random salt.
2. Sets up the Blowfish key schedule with the salt.
3. Iterates the key schedule for `2^cost` rounds. Each round depends on the previous one, making it inherently sequential (cannot be parallelized on GPU).
4. Outputs a 60-byte string: `$2a$<cost>$<22-char-salt><31-char-hash>`.

Go's bcrypt implementation in `golang.org/x/crypto/bcrypt` handles salt generation, cost factor iteration, and hash comparison. The most common mistake is not checking the error from `bcrypt.GenerateFromPassword` — it can fail if the cost is too high for the current machine.

## How Go uses it

The standard library in Go does not include bcrypt. The official package is `golang.org/x/crypto/bcrypt`:

```go
import "golang.org/x/crypto/bcrypt"
```

Two functions cover 99% of use cases:

- `bcrypt.GenerateFromPassword(password []byte, cost int) ([]byte, error)`: hashes a password with a random salt.
- `bcrypt.CompareHashAndPassword(hash, password []byte) error`: compares a password against a hash. Returns nil on match.

For scrypt: `golang.org/x/crypto/scrypt`.
For Argon2: `golang.org/x/crypto/argon2`.

## Go example

```go
package main

import (
	"errors"
	"fmt"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// userDB simulates a database of hashed passwords.
type userDB struct {
	hashes map[string][]byte
}

func newUserDB() *userDB {
	return &userDB{hashes: make(map[string][]byte)}
}

// CreateUser hashes the password and stores the hash.
func (db *userDB) CreateUser(username, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hash failed: %w", err)
	}
	db.hashes[username] = hash
	return nil
}

// Authenticate checks the password against the stored hash.
func (db *userDB) Authenticate(username, password string) error {
	hash, ok := db.hashes[username]
	if !ok {
		return errors.New("user not found")
	}
	if err := bcrypt.CompareHashAndPassword(hash, []byte(password)); err != nil {
		return errors.New("invalid password")
	}
	return nil
}

// CostBenchmark measures hashing time at different cost factors.
func CostBenchmark(password string, costs []int) {
	for _, cost := range costs {
		start := time.Now()
		hash, err := bcrypt.GenerateFromPassword([]byte(password), cost)
		duration := time.Since(start)
		if err != nil {
			log.Printf("Cost %d: error - %v", cost, err)
			continue
		}
		fmt.Printf("Cost %2d: %v (hash length: %d bytes)\n", cost, duration.Round(time.Millisecond), len(hash))
	}
}

// HashPassword is a standalone function for testing.
func HashPassword(password string, cost int) ([]byte, error) {
	if len(password) < 8 {
		return nil, errors.New("password too short: minimum 8 characters")
	}
	if cost < 4 || cost > 31 {
		return nil, errors.New("cost must be between 4 and 31")
	}
	return bcrypt.GenerateFromPassword([]byte(password), cost)
}

// VerifyPassword compares password against a bcrypt hash.
func VerifyPassword(hash []byte, password string) bool {
	return bcrypt.CompareHashAndPassword(hash, []byte(password)) == nil
}

func main() {
	db := newUserDB()

	// Create users with hashed passwords.
	db.CreateUser("alice", "correct-horse-battery-staple")
	db.CreateUser("bob", "p@ssw0rd-rocks!")

	// Successful authentication.
	if err := db.Authenticate("alice", "correct-horse-battery-staple"); err != nil {
		log.Fatal("Unexpected auth failure:", err)
	}
	fmt.Println("alice: authenticated successfully")

	// Failed authentication (wrong password).
	if err := db.Authenticate("alice", "wrong-password"); err != nil {
		fmt.Println("alice with wrong password: rejected -", err)
	}

	// Failed authentication (nonexistent user).
	if err := db.Authenticate("eve", "password"); err != nil {
		fmt.Println("eve: rejected -", err)
	}

	// Benchmark cost factors.
	fmt.Println("\n=== Cost benchmark ===")
	CostBenchmark("test-password-123", []int{4, 6, 8, 10, 12, 14})
}
```

## Step-by-step execution

For `db.CreateUser("alice", "correct-horse-battery-staple")` with default cost (10):

1. `bcrypt.GenerateFromPassword` generates a random 16-byte salt.
2. Sets up Blowfish key schedule with the password and salt.
3. Iterates the key schedule 2^10 = 1024 times (~100ms on modern hardware).
4. Returns a 60-byte encoded hash string.
5. Stored in `db.hashes["alice"]`.

For `db.Authenticate("alice", "correct-horse-battery-staple")`:

1. Retrieves hash from `db.hashes["alice"]`.
2. Parses the hash to extract salt, cost, and expected hash.
3. Runs the bcrypt key schedule with the provided password and the same salt.
4. Compares the computed hash against the stored hash using `subtle.ConstantTimeCompare`.
5. Passwords match: returns nil.

For a wrong password: step 4 produces a different hash, `CompareHashAndPassword` returns an error.

## Common mistakes

- Using fast hash functions (SHA-256, MD5) for passwords — fast hashes can be brute-forced at billions of attempts/second on GPU.
- Implementing custom password hashing logic — use bcrypt, scrypt, or Argon2 from the Go ecosystem. Do not invent your own.
- Storing bcrypt hashes in a `VARCHAR(60)` field without testing with higher cost — bcrypt output can be up to 60 bytes for cost 31, but the hash string length varies with cost.
- Setting the cost too high for production — cost 14 takes ~1 second on a typical CPU, making login a terrible user experience.
- Not checking the error from `GenerateFromPassword` — it can return `bcrypt.ErrHashTooLong` if the password exceeds 72 bytes.
- Truncating passwords before hashing — bcrypt silently truncates at 72 bytes. Pre-hash with SHA-256 first if you need longer passwords.

## Debugging walkthrough

Consider this authentication code:

```go
func login(username, password string) bool {
    hash, _ := db.getPasswordHash(username)
    return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}
```

**Symptom**: Some users cannot log in even with correct passwords. The error is silently swallowed.

**Investigation**: Add logging to see which step fails:

```go
func login(username, password string) bool {
    hash, err := db.getPasswordHash(username)
    if err != nil {
        log.Printf("Login failed: user %s not found: %v", username, err)
        return false
    }
    err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    if err != nil {
        log.Printf("Login failed: user %s password mismatch: %v", username, err)
        return false
    }
    return true
}
```

**Root cause**: The original code ignores the error from `getPasswordHash`. If the database lookup fails (e.g., wrong case in username), it passes an empty hash to `CompareHashAndPassword`, which returns an error, and the login silently fails.

**Fix**: Check database errors explicitly before comparing. Also, use constant-time comparison via bcrypt (which does it internally) and never reveal whether the username or the password was wrong (to prevent user enumeration).

## Production notes

- Use `bcrypt.DefaultCost` (10) in development, `bcrypt.MinCost` (4) for unit tests (to keep tests fast), and cost 12-14 in production depending on your hardware.
- Hash passwords at registration and store the hash. Never log, return, or expose the hash in API responses.
- When users log in, bcrypt verification timing is independent of password length (no timing leak).
- For password rotation: store the previous hash to prevent password reuse, but hash the new password before comparing with the old hash.
- In high-throughput login systems (100+ logins/second), bcrypt at cost 12 can saturate CPU. Consider Argon2id with reduced memory parameters for better throughput.

## Performance implications

- bcrypt cost 10: ~100ms per hash on modern CPU. At cost 12: ~250ms. At cost 14: ~1s.
- bcrypt is not memory-hard, so a GPU with 4096 cores can parallelize many bcrypt attempts. scrypt and Argon2 are memory-hard, requiring megabytes of memory per attempt, reducing GPU advantage.
- Login endpoints with bcrypt are susceptible to DoS: an attacker can send many login requests to overload the CPU. Implement rate limiting on login endpoints.
- User registration is less performance-critical (happens once per user) but still use cost 10-12.

## Practice task

Write Go functions `HashPassword(password string, cost int) (string, error)` and `VerifyPassword(hash, password string) bool` that:
- `HashPassword` returns an error if password is empty or cost is outside 4-31.
- `HashPassword` uses `golang.org/x/crypto/bcrypt`.
- `VerifyPassword` returns true if password matches the hash, false otherwise.
- Then write a `main()` that hashes three passwords (valid, empty, and one with cost 4), verifies each, and benchmarks cost 4 vs cost 10.

## Tests / verification

```bash
go run ./curriculum/modules/10-auth-security/lessons/05-password-hashing
go test ./curriculum/modules/10-auth-security/lessons/05-password-hashing
```

The test file `main_test.go` contains table-driven tests that verify:
- `HashPassword` returns a hash starting with `$2a$`.
- `HashPassword` returns an error for empty passwords.
- `HashPassword` returns an error for cost values outside 4-31.
- `VerifyPassword` returns true for matching password and hash.
- `VerifyPassword` returns false for non-matching passwords.

## Review questions

1. Why is SHA-256 considered unsafe for password hashing despite being a cryptographically secure hash function?
2. What does the bcrypt cost factor control, and how does it affect both security and user experience?
3. Why is a random salt important even if your password hashing algorithm is slow?
4. bcrypt truncates passwords at 72 bytes. How would you handle passwords longer than 72 bytes without reducing security?
5. If a database of bcrypt hashes is leaked, what is the attacker's best strategy to recover passwords, and what makes it expensive?

## NEXT UP

Sessions and cookies — how server-managed session state combines with secure cookie flags for web authentication.
