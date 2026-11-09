package main

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sync"
	"time"
)

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

type RefreshTokenStore struct {
	mu      sync.RWMutex
	tokens  map[string]*RefreshTokenInfo
	hmacKey []byte
}

type RefreshTokenInfo struct {
	TokenHash string
	UserID    string
	ExpiresAt time.Time
	Revoked   bool
}

func NewRefreshTokenStore() *RefreshTokenStore {
	key := make([]byte, 32)
	rand.Read(key)
	return &RefreshTokenStore{
		tokens:  make(map[string]*RefreshTokenInfo),
		hmacKey: key,
	}
}

func generateRandomString(length int) string {
	bytes := make([]byte, length)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

func (s *RefreshTokenStore) hashToken(token string) string {
	h := hmac.New(sha256.New, s.hmacKey)
	h.Write([]byte(token))
	return hex.EncodeToString(h.Sum(nil))
}

func (s *RefreshTokenStore) Issue(userID string, accessTTL, refreshTTL time.Duration) *TokenPair {
	accessToken := generateRandomString(32)
	refreshToken := generateRandomString(64)

	refreshHash := s.hashToken(refreshToken)

	s.mu.Lock()
	s.tokens[refreshHash] = &RefreshTokenInfo{
		TokenHash: refreshHash,
		UserID:    userID,
		ExpiresAt: time.Now().Add(refreshTTL),
	}
	s.mu.Unlock()

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int(accessTTL.Seconds()),
	}
}

func (s *RefreshTokenStore) Rotate(refreshToken string, accessTTL, refreshTTL time.Duration) (*TokenPair, error) {
	hash := s.hashToken(refreshToken)

	s.mu.Lock()
	info, ok := s.tokens[hash]
	if !ok {
		s.mu.Unlock()
		return nil, fmt.Errorf("refresh token not found")
	}

	if info.Revoked {
		s.mu.Unlock()
		return nil, fmt.Errorf("refresh token revoked")
	}

	if time.Now().After(info.ExpiresAt) {
		s.mu.Unlock()
		return nil, fmt.Errorf("refresh token expired")
	}

	delete(s.tokens, hash)
	s.mu.Unlock()

	return s.Issue(info.UserID, accessTTL, refreshTTL), nil
}

func (s *RefreshTokenStore) Revoke(refreshToken string) bool {
	hash := s.hashToken(refreshToken)
	s.mu.Lock()
	defer s.mu.Unlock()
	info, ok := s.tokens[hash]
	if !ok {
		return false
	}
	info.Revoked = true
	return true
}

func (s *RefreshTokenStore) Validate(refreshToken string) (string, error) {
	hash := s.hashToken(refreshToken)
	s.mu.RLock()
	defer s.mu.RUnlock()
	info, ok := s.tokens[hash]
	if !ok {
		return "", fmt.Errorf("token not found")
	}
	if info.Revoked {
		return "", fmt.Errorf("token revoked")
	}
	if time.Now().After(info.ExpiresAt) {
		return "", fmt.Errorf("token expired")
	}
	return info.UserID, nil
}

func main() {
	fmt.Println("=== Refresh Token Rotation Demo ===")

	store := NewRefreshTokenStore()
	accessTTL := 15 * time.Minute
	refreshTTL := 30 * 24 * time.Hour

	pair := store.Issue("user-abc", accessTTL, refreshTTL)
	fmt.Printf("Issued tokens:\n")
	fmt.Printf("  Access Token:  %s...\n", pair.AccessToken[:16])
	fmt.Printf("  Refresh Token: %s...\n", pair.RefreshToken[:16])
	fmt.Printf("  Expires in:    %d seconds\n", pair.ExpiresIn)

	userID, err := store.Validate(pair.RefreshToken)
	if err != nil {
		fmt.Printf("Validation failed: %v\n", err)
	} else {
		fmt.Printf("Validation successful for user: %s\n", userID)
	}

	newPair, err := store.Rotate(pair.RefreshToken, accessTTL, refreshTTL)
	if err != nil {
		fmt.Printf("Rotation failed: %v\n", err)
	} else {
		fmt.Printf("Rotated tokens:\n")
		fmt.Printf("  New Access Token:  %s...\n", newPair.AccessToken[:16])
		fmt.Printf("  New Refresh Token: %s...\n", newPair.RefreshToken[:16])
	}

	_, err = store.Validate(pair.RefreshToken)
	if err != nil {
		fmt.Printf("Old refresh token rejected after rotation: %v\n", err)
	}

	store.Revoke(newPair.RefreshToken)
	_, err = store.Validate(newPair.RefreshToken)
	if err != nil {
		fmt.Printf("Revoked token rejected: %v\n", err)
	}

	fmt.Println("\n=== Summary ===")
	fmt.Println("Access tokens: short-lived (15 min)")
	fmt.Println("Refresh tokens: long-lived (30 days), rotated on use")
	fmt.Println("Revocation: immediate invalidation")
}
