package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
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
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	ciphertext := aesGCM.Seal(nil, nonce, plaintext, nil)
	return base64.StdEncoding.EncodeToString(append(nonce, ciphertext...)), nil
}

func decrypt(encoded string, key []byte) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonceSize := aesGCM.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	return aesGCM.Open(nil, nonce, ciphertext, nil)
}

func loadFromEnv(prefix string) map[string]string {
	secrets := make(map[string]string)
	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		if len(pair) != 2 {
			continue
		}
		if strings.HasPrefix(pair[0], prefix) {
			secrets[pair[0]] = pair[1]
		}
	}
	return secrets
}

func main() {
	fmt.Println("=== Secrets Management Demo ===")

	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		panic(err)
	}
	sm := NewSecretManager(key)

	sm.Set("DB_PASSWORD", "s3cret!p4ss")
	sm.Set("API_KEY", "sk-live-abc123def456")

	secretValue, _ := sm.Get("DB_PASSWORD")
	fmt.Println("Stored secret:", "DB_PASSWORD =", secretValue)

	encrypted, err := encrypt([]byte(`{"DB_PASSWORD":"s3cret!p4ss","API_KEY":"sk-live-abc123def456"}`), key)
	if err != nil {
		panic(err)
	}
	fmt.Println("Encrypted secrets:", encrypted[:40]+"...")

	decryptedBytes, err := decrypt(encrypted, key)
	if err != nil {
		panic(err)
	}
	var decryptedMap map[string]string
	json.Unmarshal(decryptedBytes, &decryptedMap)
	fmt.Println("Decrypted secrets:", decryptedMap)

	fmt.Println("\n=== Environment Variable Secrets ===")
	os.Setenv("APP_DB_URL", "postgres://user:pass@localhost:5432/db")
	os.Setenv("APP_REDIS_URL", "redis://localhost:6379")
	envSecrets := loadFromEnv("APP_")
	for k, v := range envSecrets {
		fmt.Printf("  %s = %s\n", k, v)
	}

	fmt.Println("\nNote: In production, use a secret store like HashiCorp Vault.")
	fmt.Println("Never hardcode secrets or commit them to version control.")
}
