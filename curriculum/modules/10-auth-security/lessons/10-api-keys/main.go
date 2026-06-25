package main

import (
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

type APIKey struct {
	ID          string    `json:"id"`
	KeyPrefix   string    `json:"key_prefix"`
	KeyHash     string    `json:"key_hash"`
	Description string    `json:"description"`
	Scopes      []string  `json:"scopes"`
	CreatedAt   time.Time `json:"created_at"`
	LastUsedAt  time.Time `json:"last_used_at,omitempty"`
	ExpiresAt   time.Time `json:"expires_at,omitempty"`
}

type KeyStore struct {
	mu         sync.RWMutex
	keysByID   map[string]*APIKey
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

func GenerateAPIKey(prefix string, scopes []string, ttl time.Duration) (string, *APIKey, error) {
	if prefix == "" {
		prefix = "sk"
	}
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", nil, fmt.Errorf("random: %w", err)
	}
	randomPart := base64.RawURLEncoding.EncodeToString(b)
	fullKey := prefix + "_" + randomPart
	hash := HashKey(fullKey)
	now := time.Now()
	record := &APIKey{
		ID: generateID(), KeyPrefix: prefix + "_", KeyHash: hash,
		Description: "Auto-generated key", Scopes: scopes,
		CreatedAt: now, ExpiresAt: now.Add(ttl),
	}
	return fullKey, record, nil
}

func HashKey(key string) string {
	h := sha256.Sum256([]byte(key))
	return hex.EncodeToString(h[:])
}

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

func CreateAPIKey(prefix string, scopes []string) (fullKey string, hash string, err error) {
	fk, record, e := GenerateAPIKey(prefix, scopes, 365*24*time.Hour)
	if e != nil {
		return "", "", e
	}
	return fk, record.KeyHash, nil
}

func ValidateAPIKeyFunc(store map[string]*APIKey, presentedKey, requiredScope string) error {
	hash := HashKey(presentedKey)
	record, ok := store[hash]
	if !ok {
		return errors.New("invalid key")
	}
	if requiredScope != "" && !hasScope(record.Scopes, requiredScope) {
		return fmt.Errorf("insufficient scope: need %s", requiredScope)
	}
	return nil
}

func main() {
	store := NewKeyStore()

	key1, record1, _ := GenerateAPIKey("api_key", []string{"documents:read", "documents:write"}, 365*24*time.Hour)
	store.Add(record1)
	fmt.Printf("Key 1: %s...  hash: %s...  scopes: %v\n", key1[:20], record1.KeyHash[:16], record1.Scopes)

	key2, record2, _ := GenerateAPIKey("api_key", []string{"documents:read"}, 365*24*time.Hour)
	store.Add(record2)
	fmt.Printf("Key 2: %s...  scopes: %v\n", key2[:20], record2.Scopes)

	// Validate
	_, err := ValidateKey(store, key1, "documents:write")
	if err != nil {
		log.Fatal("Unexpected:", err)
	}
	fmt.Println("Key 1 validated for documents:write")

	_, err = ValidateKey(store, key1, "documents:delete")
	if err != nil {
		fmt.Println("Key 1 correctly denied for documents:delete:", err)
	}

	_, err = ValidateKey(store, "invalid-key", "documents:read")
	if err != nil {
		fmt.Println("Invalid key correctly rejected:", err)
	}

	// HTTP
	http.Handle("/api/documents", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := r.Header.Get("X-API-Key")
		if key == "" {
			http.Error(w, `{"error":"missing API key"}`, http.StatusUnauthorized)
			return
		}
		_, err := ValidateKey(store, key, "documents:read")
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusForbidden)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "read access granted"})
	}))

	fmt.Println("\nAPI key server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
