# Secrets management

## Learning objective

Manage secrets in Go applications using environment variables, encrypted storage, secret stores, and proper `.gitignore` patterns to prevent accidental exposure of credentials and cryptographic keys.

## Why this matters

Hardcoded secrets in source code are the most common cause of credential leaks. The 2023 GitHub secret scanning report found over 1 million secret leaks per year, including API keys, database passwords, and cloud credentials. A leaked secret can lead to data breaches, cloud account takeover, and massive financial liability. Go engineers must follow secret management best practices to ensure credentials are never committed to version control and are stored securely in production.

## Mental model

Secrets are the keys to your digital kingdom. They unlock databases, cloud accounts, third-party APIs, and encryption systems. Treat secrets like physical keys: you would not leave your house key taped to the front door (hardcoded in source code) or pinned to a public bulletin board (committed to a public repo).

A secret manager is a secure vault. You authenticate with the vault using a short-lived token, and the vault gives you the keys you need. If a key is compromised, you rotate it in the vault without redeploying the application.

## Core idea

Secret management principles:

1. **Never hardcode secrets** in source code. Use environment variables, secret files, or a secret store.
2. **Never commit secrets** to version control. Add `.env`, `*.key`, `secrets.*` to `.gitignore`.
3. **Encrypt secrets at rest** when stored outside a secret manager.
4. **Rotate secrets regularly** and especially after a suspected leak.
5. **Use the principle of least privilege**: each service should have access only to the secrets it needs.
6. **Audit secret access**: know who accessed which secret and when.

| Storage method | Security | Complexity | Use case |
|---|---|---|---|
| Environment variables | Low | None | Local dev, non-sensitive config |
| `.env` files (gitignored) | Low | Minimal | Multi-service local dev |
| Encrypted config files | Medium | Low | CI/CD, small deployments |
| HashiCorp Vault | High | High | Enterprise production |
| AWS Secrets Manager | High | Medium | AWS-native deployments |
| GCP Secret Manager | High | Medium | GCP-native deployments |
| Kubernetes Secrets | Medium | Medium | Kubernetes deployments |

## Under the hood

Environment variables are inherited from the parent process. When you run `export DB_PASSWORD=s3cret`, that value is accessible to the current shell and its child processes. In Go, `os.Getenv("DB_PASSWORD")` reads from the process's environment block. Environment variables are visible in `/proc/pid/environ` on Linux and can be leaked through error pages, crash dumps, and child processes.

Secret managers (Vault, AWS Secrets Manager) store secrets encrypted at rest and in transit. Applications authenticate to the secret manager using a short-lived token (e.g., Vault's Kubernetes auth via a service account token, or AWS IAM roles via instance metadata). Secrets are fetched at application startup (or on-demand) and cached in memory.

Encryption at rest uses symmetric-key algorithms like AES-256-GCM. The encryption key is itself a secret that must be managed securely (often derived from a master key or stored in a hardware security module).

## How Go uses it

Go applications typically load secrets at startup:

```go
dbPassword := os.Getenv("DB_PASSWORD")
if dbPassword == "" {
    log.Fatal("DB_PASSWORD not set")
}
```

For production, Go services use SDKs for secret stores:

- HashiCorp Vault: `github.com/hashicorp/vault/api`
- AWS Secrets Manager: `github.com/aws/aws-sdk-go-v2/service/secretsmanager`
- GCP Secret Manager: `cloud.google.com/go/secretmanager`
- Kubernetes: mounted as files in `/etc/secrets/`

For local development, tools like `direnv` load `.envrc` files, and projects use `godotenv` to load `.env` files.

## Go example

```go
package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"os"
)

type SecretManager struct {
	encryptionKey []byte
	secrets       map[string]string
}

func NewSecretManager(key []byte) *SecretManager {
	return &SecretManager{
		encryptionKey: key,
		secrets:       make(map[string]string),
	}
}

func (sm *SecretManager) Set(name, value string) {
	sm.secrets[name] = value
}

func (sm *SecretManager) Get(name string) (string, bool) {
	v, ok := sm.secrets[name]
	return v, ok
}

func encrypt(plaintext []byte, key []byte) (string, error) {
	block, _ := aes.NewCipher(key)
	aesGCM, _ := cipher.NewGCM(block)
	nonce := make([]byte, aesGCM.NonceSize())
	io.ReadFull(rand.Reader, nonce)
	ciphertext := aesGCM.Seal(nil, nonce, plaintext, nil)
	return base64.StdEncoding.EncodeToString(append(nonce, ciphertext...)), nil
}

func decrypt(encoded string, key []byte) ([]byte, error) {
	data, _ := base64.StdEncoding.DecodeString(encoded)
	block, _ := aes.NewCipher(key)
	aesGCM, _ := cipher.NewGCM(block)
	nonceSize := aesGCM.NonceSize()
	return aesGCM.Open(nil, data[:nonceSize], data[nonceSize:], nil)
}

func main() {
	key := make([]byte, 32)
	rand.Read(key)
	sm := NewSecretManager(key)
	sm.Set("DB_PASSWORD", "s3cret!p4ss")
	val, _ := sm.Get("DB_PASSWORD")
	fmt.Println("Secret:", val)
}
```

## Step-by-step execution

For a Go service deploying to production with secrets:

1. Developer sets `DB_PASSWORD` in the production environment (via Kubernetes Secret, Vault, or cloud console).
2. At startup, the Go application calls `os.Getenv("DB_PASSWORD")`.
3. If the variable is empty, the application exits with an error (fail-fast).
4. The application uses the value to open a database connection.
5. The value remains in memory for the lifetime of the process.
6. When the process exits, the memory is reclaimed by the OS.

For rotation:
1. Operator updates the secret in the vault/secret store.
2. Operator restarts the application (or sends SIGUSR1 to trigger reload).
3. The new process picks up the updated secret.
4. The old process exits, and its memory (with the old secret) is freed.

## Common mistakes

- Mistake: Hardcoding secrets in source code that gets committed to version control.
  - Why it happens: Developers hardcode credentials for convenience during development.
  - Fix: Use environment variables or `.env` files from the start. Add `*.env` and `secrets.*` to `.gitignore`.

- Mistake: Storing secrets in environment variables that are logged or exposed in error pages.
  - Why it happens: A structured logger captures all environment variables, or a panic handler dumps the environment.
  - Fix: Never log environment variables. Redact known secret names in log output.

- Mistake: Using a single secret for multiple purposes (JWT signing, DB encryption, API keys).
  - Why it happens: It is easier to manage one key than many.
  - Fix: Each secret should have a single purpose. Compromise of one key should not affect other systems.

- Mistake: Not rotating secrets regularly or after a suspected leak.
  - Why it happens: Rotation is manual and disruptive.
  - Fix: Automate rotation. Use a secret store that supports versioning and automatic rotation schedules.

- Mistake: Storing encryption keys next to encrypted data in the same repository.
  - Why it happens: Developers encrypt config files but commit the encryption key in the same repo.
  - Fix: The encryption key must be stored separately from the encrypted data (e.g., in a CI/CD secret variable).

## Debugging walkthrough

Consider this broken setup:

```go
func main() {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("DB connection failed: %v", err)
	}
	defer db.Close()
}
```

Symptom: The application connects to a staging database in production.

Investigation: Check the actual environment variable value:

```go
fmt.Println("DATABASE_URL:", os.Getenv("DATABASE_URL"))
```

Root cause: The `DATABASE_URL` environment variable was set to the staging database URL in the CI/CD configuration, and no one updated it for the production deployment.

Fix: Use separate environment configurations with validation:

```go
func mustGetenv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("required environment variable %s not set", key)
	}
	return val
}

func main() {
	env := mustGetenv("APP_ENV") // "production", "staging", "development"
	dbURL := mustGetenv("DATABASE_URL")
	if env == "production" && !strings.Contains(dbURL, "prod") {
		log.Fatal("DATABASE_URL does not look like a production URL")
	}
	db, err := sql.Open("postgres", dbURL)
	// ...
}
```

## Production notes

- Always use a secret store or cloud secrets manager for production secrets. Environment variables are acceptable for non-sensitive config but not for credentials.
- Encrypt secrets at rest. AES-256-GCM is the recommended algorithm.
- Rotate secrets at least every 90 days, and immediately after any suspected compromise.
- Use short-lived credentials where possible: database tokens with TTLs, IAM roles instead of long-lived keys.
- Audit secret access: know which service accessed which secret and when.
- In Go, use `runtime.SetFinalizer` or `memguard` to zero out secrets in memory after use, preventing them from being dumped in core files.
- `.gitignore` must include: `*.env`, `.env.*`, `*.key`, `secrets.*`, `credentials.*`, `*.pem` (for private keys).

## Performance implications

- Reading `os.Getenv` is O(n) in the number of environment variables and takes a few microseconds. It is called once at startup.
- Secret store SDK calls add network latency (typically 10-100ms). Cache secrets in memory with a TTL to avoid per-request calls.
- Encryption/decryption with AES-GCM is fast (hardware-accelerated on modern CPUs, ~1 GB/s throughput).
- The performance bottleneck in secret management is rarely the crypto or the API call; it is the operational complexity of rotating and auditing secrets.

## Practice task

Write a function `readEncryptedConfig(path string, key []byte) (map[string]string, error)` that reads an encrypted JSON file, decrypts it with AES-256-GCM, and returns the config map.

Then write a function `writeEncryptedConfig(path string, config map[string]string, key []byte) error` that encrypts and writes the config.

Write a `main()` that:
1. Generates a random 32-byte encryption key.
2. Writes an encrypted config file with keys `DB_PASSWORD`, `API_KEY`, `JWT_SECRET`.
3. Reads back the encrypted file and prints the decrypted values.
4. Confirms that reading with a wrong key returns an error.

## Tests / verification

```bash
go run ./curriculum/modules/10-auth-security/lessons/20-secrets-management
go test ./curriculum/modules/10-auth-security/lessons/20-secrets-management
```

The existing tests verify secret storage and retrieval, encryption/decryption roundtrips, decryption failure with wrong keys, and loading secrets from environment variables.

## Review questions

1. Why are environment variables insufficient for production secret management?
2. How does the principle of least privilege apply to secret management in a microservices architecture?
3. What encryption algorithm and mode are recommended for encrypting secrets at rest?
4. How would you implement secret rotation without downtime in a Go API service?
5. Why should you never log environment variables, and how can you prevent accidental exposure?

## NEXT UP

Dependency security -- managing supply chain risks, verifying dependencies, minimizing attack surface, and understanding the software supply chain in Go projects.
