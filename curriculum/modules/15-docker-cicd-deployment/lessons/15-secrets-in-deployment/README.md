# Secrets in deployment

## Learning objective

Implement a secrets management abstraction in Go that supports multiple backends (env vars, in-memory, secret stores), parse .env files safely, sanitize secrets in logs, and resolve templates with secret substitution.

## Why this matters

Secrets are the most sensitive configuration a service handles. Database passwords, API keys, JWT signing keys, and TLS private keys grant access to critical infrastructure. Mishandling secrets is the leading cause of data breaches. A Go engineer must know how to load, store, and access secrets without leaking them into logs, error messages, or version control. Using a secrets abstraction (rather than ad-hoc `os.Getenv` calls) makes it possible to switch between local development (env vars or .env files) and production (Vault or AWS Secrets Manager) without code changes.

## Mental model

Secrets are config values that must be kept confidential. They differ from regular config in three ways:

1. **Access control**: Only authorized services and people should read them.
2. **Encryption at rest and in transit**: Secrets should never be stored in plaintext.
3. **Audit trail**: Every access to a secret should be logged.

Think of secrets like physical keys to a building. Config is the address (public). Secrets are the key (private). You would not tape the key next to the address on the front door, so do not put secrets in source code, config files committed to git, or environment dumps in logs.

A secret store is a lockbox. You authenticate to the lockbox, request a specific key, and receive the value. The lockbox logs every access. If someone tries to open the lockbox without authorization, an alarm sounds.

## Core idea

The secrets management hierarchy for Go services:

| Layer | Tool | Use case |
|---|---|---|
| Local dev | `.env` file (never committed) | Developer workstation |
| CI/CD | Encrypted env vars in CI (GitHub Actions secrets) | Pipeline needs API keys |
| Production | Secret store (HashiCorp Vault, AWS Secrets Manager, GCP Secret Manager) | Runtime secret access |
| Orchestration | Kubernetes Secrets (encrypted at rest) | Container environment |

The `SecretStore` interface abstracts the backend. Code depends on the interface, not a concrete implementation. This lets you swap backends without changing application logic.

## Under the hood

A production secret store like HashiCorp Vault works as follows:

1. The application authenticates to Vault using a token, Kubernetes service account, or AWS IAM role.
2. Vault verifies the identity and checks ACL policies to determine which secrets the caller can read.
3. The application requests a secret at a path (e.g., `secret/data/db-password`).
4. Vault returns the encrypted secret over TLS. The secret is decrypted in Vault's memory and sent to the application.
5. The application caches the secret in memory for a configurable TTL (default: typically minutes).

AWS Secrets Manager similarly provides encrypted storage with automatic rotation. Secrets are encrypted with KMS and accessed via the AWS SDK.

In contrast, `.env` files are plaintext files with `KEY=VALUE` lines. They are loaded into the process environment at startup. The `godotenv` library is popular for Go, but the standard library can parse them trivially. `.env` files must never be committed to version control.

## How Go uses it

Go applications typically load secrets at startup and cache them in memory for the process lifetime. The pattern:

1. At startup, authenticate to the secret store (Vault, AWS, or just read env vars).
2. Load all required secrets into a `Secrets` struct.
3. Use the secrets throughout the service without additional store calls.
4. Never log secret values. Use a sanitization function to show only the first/last few characters for debugging.

Many Go projects define a `SecretStore` interface:

```go
type SecretStore interface {
    Get(key string) (string, error)
    Set(key, value string) error
}
```

This allows testing with an in-memory store and running in production with Vault.

## Go example

```go
package main

import (
	"fmt"
	"os"
	"strings"
)

type SecretStore interface {
	Get(key string) (string, error)
}

type InMemoryStore struct {
	data map[string]string
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{data: make(map[string]string)}
}

func (s *InMemoryStore) Get(key string) (string, error) {
	val, ok := s.data[key]
	if !ok {
		return "", fmt.Errorf("secret %s not found", key)
	}
	return val, nil
}

func (s *InMemoryStore) Set(key, value string) error {
	s.data[key] = value
	return nil
}

type SecretsManager struct {
	store SecretStore
}

func NewSecretsManager(store SecretStore) *SecretsManager {
	return &SecretsManager{store: store}
}

func (sm *SecretsManager) Get(key string) (string, error) {
	return sm.store.Get(key)
}

func (sm *SecretsManager) ResolveTemplate(tmpl string) (string, error) {
	var missing []string
	result := os.Expand(tmpl, func(key string) string {
		val, err := sm.store.Get(key)
		if err != nil {
			missing = append(missing, key)
			return ""
		}
		return val
	})
	if len(missing) > 0 {
		return "", fmt.Errorf("missing secrets: %s", strings.Join(missing, ", "))
	}
	return result, nil
}

func Sanitize(value string) string {
	if len(value) <= 8 {
		return strings.Repeat("*", len(value))
	}
	return value[:2] + strings.Repeat("*", len(value)-4) + value[len(value)-2:]
}

func main() {
	store := NewInMemoryStore()
	store.Set("DB_PASSWORD", "s3cur3P@ss!")
	store.Set("API_KEY", "sk-abc123def456ghi789")

	sm := NewSecretsManager(store)

	dsn, err := sm.ResolveTemplate("postgres://app:${DB_PASSWORD}@localhost/mydb")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("DSN:", dsn)
	fmt.Println("Sanitized API key:", Sanitize("sk-abc123def456ghi789"))
}
```

## Step-by-step execution

For `ResolveTemplate("postgres://app:${DB_PASSWORD}@localhost/mydb")`:

1. `os.Expand` scans the template string and finds `${DB_PASSWORD}`.
2. The callback function is invoked with key `"DB_PASSWORD"`.
3. `sm.store.Get("DB_PASSWORD")` returns `"s3cur3P@ss!"`.
4. `os.Expand` substitutes `${DB_PASSWORD}` with `"s3cur3P@ss!"`.
5. Result: `"postgres://app:s3cur3P@ss!@localhost/mydb"`.

For `Sanitize("sk-abc123def456ghi789")`:

1. Length = 19 characters, which is > 8.
2. First 2 chars: `"sk"`.
3. Last 2 chars: `"89"`.
4. Middle: 15 asterisks (`"***************"`).
5. Result: `"sk***************89"`.

## Common mistakes

- **Committing .env files to git**: Add `.env` to `.gitignore` immediately. A committed secret is compromised and must be rotated.
- **Logging secrets in error messages**: If a database connection fails and the error includes the DSN, the password is logged. Catch the error and log a sanitized version.
- **Hardcoding fallback secrets**: `os.Getenv("DB_PASSWORD")` should not have a fallback default secret. If the env var is missing, fail fast.
- **Long-lived secret caches**: If a secret is rotated while the process is running, the old value is still used. Use a TTL or a watch mechanism for critical secrets.
- **Using the same secret across environments**: A development database password should differ from production. Never share secrets between environments.

## Debugging walkthrough

Consider a service that fails to authenticate to the database in production:

```
Error: pq: password authentication failed for user "app"
```

**Symptom**: The database rejects the connection. The service works locally.

**Investigation**:

1. Check which secret store is configured. If using Vault, verify the Vault token is valid and the policy allows reading the secret path.
2. Check if the secret exists in the store:
   ```go
   store := vault.NewClient(token)
   secret, err := store.Get("production/db-password")
   ```
3. Compare the secret value with the actual database password set by the DBA.

**Root cause**: The database password was rotated by the DBA, but the secret in Vault was not updated. The service is using the old password.

**Fix**: Update the secret in Vault: `vault kv put secret/db-password value=newpass`. The service needs a restart or TTL-based refresh to pick up the new value.

Another scenario:

```
Error: secret DB_PASSWORD not found
```

**Root cause**: The `InMemoryStore` does not have the secret set. In production, the Vault client is not configured, so the `EnvSecretStore` is used, but `DB_PASSWORD` is not set in the environment either.

**Fix**: Ensure the secret is set in the environment or the secret store is properly configured and authenticated.

## Production notes

In production Go services:

- **Use a dedicated secrets library**: `github.com/hashicorp/vault/api` for Vault, `aws-sdk-go-v2/service/secretsmanager` for AWS, or `cloud.google.com/go/secretmanager` for GCP.
- **Short-lived credentials**: Prefer dynamic secrets (Vault's database secret engine) over static passwords. Dynamic secrets are valid for a limited time and auto-expire.
- **Secret rotation**: Implement a rotation handler that refreshes secrets before they expire. For Kubernetes, tools like `external-secrets` or `SealedSecrets` manage this.
- **Audit logging**: Log secret access events (which secret, what time, by which service) to a secure, immutable audit log.
- **Minimum privilege**: Each service should only have access to the secrets it needs. A web service does not need the database root password.

## Performance implications

- **Secret store calls add latency**: Each call to Vault adds 10-50ms. Cache secrets in memory and only refresh on TTL expiry.
- **In-memory cache is zero-overhead**: Reading a secret from a Go map is sub-microsecond.
- **Vault token renewal**: Vault tokens have a TTL. Automatic renewal adds periodic load. A service with 1000 instances each renewing every hour creates negligible load.
- **Encrypted env vars in CI**: GitHub Actions secrets are decrypted at workflow start and injected as environment variables. Zero runtime overhead.

## Practice task

Write a function `LoadSecretsFromDotEnv(path string) (*SecretsManager, error)` that:

1. Opens a file at the given path.
2. Reads lines in `KEY=VALUE` format (skip comments and blank lines).
3. Populates an `InMemoryStore` with the parsed entries.
4. Returns a `SecretsManager` backed by the populated store.
5. If the file does not exist, return an error.

Then write a `main()` that creates a temporary `.env` file with at least three secrets, loads them, resolves a template URL, and prints the result with sanitized values.

## Tests / verification

```bash
go run ./curriculum/modules/15-docker-cicd-deployment/lessons/15-secrets-in-deployment
go test ./curriculum/modules/15-docker-cicd-deployment/lessons/15-secrets-in-deployment
```

The existing tests verify the in-memory store, env secret store, template resolution, dot env parsing (with comments, quotes, invalid lines), sanitization of short and long values, and secret validation. After completing the practice task, add tests for `LoadSecretsFromDotEnv` covering missing files, valid files, and files with comments.

## Review questions

1. What is the difference between config and secrets? Give two examples of each.
2. Why should secrets never be hardcoded or committed to version control? What is the remediation if a secret is committed?
3. How does the `SecretStore` interface make testing easier compared to using `os.Getenv` directly throughout the codebase?
4. What are the advantages of dynamic secrets (e.g., Vault database secrets engine) over static database passwords?
5. A developer reports a secret value appears in log output. What are two ways this could happen, and how would you fix each?

## NEXT UP

Vulnerability scanning: detecting CVEs in dependencies and Docker images before they reach production.
