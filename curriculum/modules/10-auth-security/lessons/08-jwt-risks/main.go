package main

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	rsaPrivateKey *rsa.PrivateKey
	rsaPublicKey  *rsa.PublicKey
	hmacSecret    = []byte("my-super-secret-key-change-in-production")
)

func init() {
	var err error
	rsaPrivateKey, err = rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatalf("key generation: %v", err)
	}
	rsaPublicKey = &rsaPrivateKey.PublicKey
}

type DenyList struct {
	mu sync.RWMutex
	m  map[string]time.Time
}

func NewDenyList() *DenyList {
	return &DenyList{m: make(map[string]time.Time)}
}

func (dl *DenyList) Revoke(jti string, expiry time.Time) {
	dl.mu.Lock()
	dl.m[jti] = expiry
	dl.mu.Unlock()
}

func (dl *DenyList) IsRevoked(jti string) bool {
	dl.mu.RLock()
	exp, ok := dl.m[jti]
	dl.mu.RUnlock()
	if !ok {
		return false
	}
	if time.Now().After(exp) {
		dl.mu.Lock()
		delete(dl.m, jti)
		dl.mu.Unlock()
		return false
	}
	return true
}

type SafeClaims struct {
	jwt.RegisteredClaims
	Role  string `json:"role"`
	Scope string `json:"scope"`
}

func IssueToken(userID, role string) (string, error) {
	now := time.Now()
	claims := SafeClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			ExpiresAt: jwt.NewNumericDate(now.Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(now),
			ID:        fmt.Sprintf("%s-%d", userID, now.UnixNano()),
		},
		Role:  role,
		Scope: "api:read",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(rsaPrivateKey)
}

func VerifyTokenSafe(tokenStr string, denyList *DenyList) (*SafeClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &SafeClaims{}, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodRS256 {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return rsaPublicKey, nil
	})
	if err != nil {
		return nil, fmt.Errorf("parse failed: %w", err)
	}
	claims, ok := token.Claims.(*SafeClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}
	if denyList.IsRevoked(claims.ID) {
		return nil, errors.New("token revoked")
	}
	return claims, nil
}

// Vulnerable verifier for demonstration.
func VerifyTokenVulnerable(tokenStr string) (string, error) {
	token, err := jwt.Parse(tokenStr, nil)
	if err != nil {
		return "", err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", errors.New("invalid token")
	}
	sub, _ := claims["sub"].(string)
	return sub, nil
}

func VerifyTokenWithAlgCheck(tokenStr string, key interface{}) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodRS256 {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return key, nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}

func DemonstrateForgery(publicKey *rsa.PublicKey) string {
	forged := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": "admin",
		"exp": time.Now().Add(1 * time.Hour).Unix(),
		"iat": time.Now().Unix(),
	})
	// In the real attack, the HMAC secret would be the public key bytes.
	keyBytes := []byte(fmt.Sprintf("%v", publicKey))
	forgedStr, _ := forged.SignedString(keyBytes)
	return forgedStr
}

func main() {
	fmt.Println("=== JWT Risk Demonstrations ===")
	fmt.Println()

	token, _ := IssueToken("alice", "user")
	fmt.Println("Legitimate token:", token[:60], "...")
	fmt.Println()

	denyList := NewDenyList()
	claims, err := VerifyTokenSafe(token, denyList)
	if err != nil {
		log.Fatal("Unexpected:", err)
	}
	fmt.Printf("Verified: user=%s, role=%s, jti=%s\n\n", claims.Subject, claims.Role, claims.ID)

	denyList.Revoke(claims.ID, claims.ExpiresAt.Time)
	_, err = VerifyTokenSafe(token, denyList)
	if err != nil && strings.Contains(err.Error(), "revoked") {
		fmt.Println("Revoked token correctly rejected")
	}

	// Algorithm confusion demo.
	forged := DemonstrateForgery(rsaPublicKey)
	fmt.Println("Forged token (HS256 with public key):", forged[:60], "...")

	sub, err := VerifyTokenVulnerable(forged)
	if err != nil {
		fmt.Println("Safe: Vulnerable verifier rejected:", err)
	} else {
		fmt.Println("WARNING: Vulnerable verifier accepted forged token as:", sub)
	}

	_, err = VerifyTokenWithAlgCheck(forged, rsaPublicKey)
	if err != nil {
		fmt.Println("Safe verifier rejected forged token:", err)
	}

	// None-alg demo.
	header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"none","typ":"JWT"}`))
	payload := base64.RawURLEncoding.EncodeToString([]byte(`{"sub":"admin","role":"admin","exp":9999999999}`))
	noneToken := fmt.Sprintf("%s.%s.", header, payload)

	sub2, err := VerifyTokenVulnerable(noneToken)
	if err != nil {
		fmt.Println("Safe: None-alg rejected:", err)
	} else {
		fmt.Println("WARNING: None-alg token accepted as:", sub2)
	}
}

var _ = hmac.Equal
