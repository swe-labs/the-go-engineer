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
	mu              sync.RWMutex
	tokens          map[string]*RefreshTokenInfo
	revokedFamilies map[string]bool
	hmacKey         []byte
}

type RefreshTokenInfo struct {
	TokenHash string
	UserID    string
	FamilyID  string
	ExpiresAt time.Time
	Used      bool
	Revoked   bool
}

func NewRefreshTokenStore() *RefreshTokenStore {
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		panic(err) // crypto/rand failure is unrecoverable for a token store
	}
	return &RefreshTokenStore{
		tokens:          make(map[string]*RefreshTokenInfo),
		revokedFamilies: make(map[string]bool),
		hmacKey:         key,
	}
}

func generateRandomString(length int) string {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		panic(err) // crypto/rand failure is unrecoverable for token generation
	}
	return hex.EncodeToString(bytes)
}

func (s *RefreshTokenStore) hashToken(token string) string {
	h := hmac.New(sha256.New, s.hmacKey)
	h.Write([]byte(token))
	return hex.EncodeToString(h.Sum(nil))
}

func (s *RefreshTokenStore) issueInFamilyLocked(userID, familyID string, accessTTL, refreshTTL time.Duration) *TokenPair {
	accessToken := generateRandomString(32)
	refreshToken := generateRandomString(64)

	refreshHash := s.hashToken(refreshToken)

	s.tokens[refreshHash] = &RefreshTokenInfo{
		TokenHash: refreshHash,
		UserID:    userID,
		FamilyID:  familyID,
		ExpiresAt: time.Now().Add(refreshTTL),
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int(accessTTL.Seconds()),
	}
}

func (s *RefreshTokenStore) revokeFamilyLocked(familyID string) {
	s.revokedFamilies[familyID] = true
	for _, info := range s.tokens {
		if info.FamilyID == familyID {
			info.Revoked = true
		}
	}
}

func (s *RefreshTokenStore) Issue(userID string, accessTTL, refreshTTL time.Duration) *TokenPair {
	familyID := generateRandomString(16)
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.issueInFamilyLocked(userID, familyID, accessTTL, refreshTTL)
}

func (s *RefreshTokenStore) Rotate(refreshToken string, accessTTL, refreshTTL time.Duration) (*TokenPair, error) {
	hash := s.hashToken(refreshToken)

	s.mu.Lock()
	defer s.mu.Unlock()

	info, ok := s.tokens[hash]
	if !ok {
		return nil, fmt.Errorf("refresh token not found")
	}

	if info.Revoked || s.revokedFamilies[info.FamilyID] {
		return nil, fmt.Errorf("refresh token revoked")
	}

	if info.Used {
		// Replay of a consumed token indicates theft: revoke the entire family.
		s.revokeFamilyLocked(info.FamilyID)
		return nil, fmt.Errorf("refresh token reuse detected: family revoked")
	}

	if time.Now().After(info.ExpiresAt) {
		return nil, fmt.Errorf("refresh token expired")
	}

	info.Used = true
	return s.issueInFamilyLocked(info.UserID, info.FamilyID, accessTTL, refreshTTL), nil
}

func (s *RefreshTokenStore) Revoke(refreshToken string) bool {
	hash := s.hashToken(refreshToken)
	s.mu.Lock()
	defer s.mu.Unlock()
	info, ok := s.tokens[hash]
	if !ok {
		return false
	}
	s.revokeFamilyLocked(info.FamilyID)
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
	if info.Revoked || s.revokedFamilies[info.FamilyID] {
		return "", fmt.Errorf("token revoked")
	}
	if info.Used {
		return "", fmt.Errorf("token already used")
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
