package main

import (
	"fmt"
	"os"
	"strings"
)

type SecretStore interface {
	Get(key string) (string, error)
	Set(key, value string) error
}

type EnvSecretStore struct{}

func (e *EnvSecretStore) Get(key string) (string, error) {
	val := os.Getenv(key)
	if val == "" {
		return "", fmt.Errorf("secret %s not found", key)
	}
	return val, nil
}

func (e *EnvSecretStore) Set(key, value string) error {
	return os.Setenv(key, value)
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

func (sm *SecretsManager) Set(key, value string) error {
	return sm.store.Set(key, value)
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

type DotEnvEntry struct {
	Key   string
	Value string
}

func ParseDotEnv(content string) ([]DotEnvEntry, error) {
	var entries []DotEnvEntry
	for _, line := range strings.Split(content, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid .env line: %s", line)
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		value = strings.Trim(value, `"'`)
		entries = append(entries, DotEnvEntry{Key: key, Value: value})
	}
	return entries, nil
}

func Sanitize(value string) string {
	if len(value) <= 8 {
		return strings.Repeat("*", len(value))
	}
	return value[:2] + strings.Repeat("*", len(value)-4) + value[len(value)-2:]
}

type SecretValidator struct{}

func (v *SecretValidator) ValidateLength(key, value string, minLen int) error {
	if len(value) < minLen {
		return fmt.Errorf("secret %s is too short (%d < %d)", key, len(value), minLen)
	}
	return nil
}

func (v *SecretValidator) ValidateNotEmpty(key, value string) error {
	if value == "" {
		return fmt.Errorf("secret %s is empty", key)
	}
	return nil
}

func main() {
	store := NewInMemoryStore()
	sm := NewSecretsManager(store)

	sm.Set("DB_PASSWORD", "s3cur3P@ss!")
	sm.Set("API_KEY", "sk-abc123def456ghi789")
	sm.Set("JWT_SECRET", "my-jwt-secret-key-2026")

	dsnTemplate := "postgres://app:${DB_PASSWORD}@localhost:5432/mydb"
	dsn, err := sm.ResolveTemplate(dsnTemplate)
	if err != nil {
		fmt.Fprintf(os.Stderr, "template error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Resolved DSN:", dsn)
	fmt.Println("Sanitized DSN:", Sanitize(dsn))
	fmt.Println("Sanitized API_KEY:", Sanitize("sk-abc123def456ghi789"))

	dotenv := `# Database config
DB_HOST=localhost
DB_PORT=5432
DB_PASSWORD=supersecret
`
	entries, err := ParseDotEnv(dotenv)
	if err != nil {
		fmt.Fprintf(os.Stderr, "parse error: %v\n", err)
		os.Exit(1)
	}
	for _, e := range entries {
		fmt.Printf("Loaded: %s=%s\n", e.Key, Sanitize(e.Value))
	}
}
