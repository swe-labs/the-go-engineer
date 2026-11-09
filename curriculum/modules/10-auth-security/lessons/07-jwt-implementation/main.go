package main

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type CustomClaims struct {
	jwt.RegisteredClaims
	Role     string `json:"role"`
	TenantID string `json:"tenant_id"`
}

var (
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
)

func init() {
	var err error
	privateKey, err = rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatalf("Failed to generate key: %v", err)
	}
	publicKey = &privateKey.PublicKey
}

func IssueToken(userID, role, tenantID string, ttl time.Duration) (string, error) {
	now := time.Now()
	claims := CustomClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "the-go-engineer",
			Subject:   userID,
			Audience:  jwt.ClaimStrings{"api.thegoengineer.com"},
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(now),
			ID:        fmt.Sprintf("%s-%d", userID, now.UnixNano()),
		},
		Role:     role,
		TenantID: tenantID,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(privateKey)
}

func ValidateToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return publicKey, nil
	})
	if err != nil {
		return nil, fmt.Errorf("token validation failed: %w", err)
	}
	claims, ok := token.Claims.(*CustomClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}
	return claims, nil
}

func CreateJWT(userID, role string, secret []byte, ttl time.Duration) (string, error) {
	now := time.Now()
	claims := jwt.RegisteredClaims{
		Subject:   userID,
		ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
		IssuedAt:  jwt.NewNumericDate(now),
		ID:        fmt.Sprintf("%s-%d", userID, now.UnixNano()),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token.Header["role"] = role
	return token.SignedString(secret)
}

func VerifyJWT(tokenStr string, secret []byte) (string, string, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secret, nil
	})
	if err != nil {
		return "", "", err
	}
	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok || !token.Valid {
		return "", "", fmt.Errorf("invalid token")
	}
	role, _ := token.Header["role"].(string)
	return claims.Subject, role, nil
}

func main() {
	// Standalone usage
	token, _ := IssueToken("alice", "admin", "acme-corp", 1*time.Hour)
	fmt.Println("Token:", token[:80], "...")

	claims, err := ValidateToken(token)
	if err != nil {
		log.Fatal("Validation failed:", err)
	}
	fmt.Printf("Validated: user=%s role=%s tenant=%s\n", claims.Subject, claims.Role, claims.TenantID)

	// HMAC example
	hmacToken, _ := CreateJWT("bob", "editor", []byte("secret-key"), 30*time.Minute)
	sub, role, err := VerifyJWT(hmacToken, []byte("secret-key"))
	if err != nil {
		log.Fatal("HMAC validation failed:", err)
	}
	fmt.Printf("HMAC: user=%s role=%s\n", sub, role)

	// HTTP server
	http.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var req struct {
			UserID   string `json:"user_id"`
			Role     string `json:"role"`
			TenantID string `json:"tenant_id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}
		token, err := IssueToken(req.UserID, req.Role, req.TenantID, 1*time.Hour)
		if err != nil {
			http.Error(w, "token generation failed", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"token": token})
	})

	http.Handle("/protected", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
			http.Error(w, `{"error":"missing token"}`, http.StatusUnauthorized)
			return
		}
		tokenStr := strings.TrimPrefix(auth, "Bearer ")
		claims, err := ValidateToken(tokenStr)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusUnauthorized)
			return
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "authenticated", "user_id": claims.Subject,
			"role": claims.Role, "tenant_id": claims.TenantID,
		})
	}))

	fmt.Println("JWT server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
