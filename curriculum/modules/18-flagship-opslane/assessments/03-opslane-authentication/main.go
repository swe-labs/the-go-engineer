package main

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrTokenExpired       = errors.New("authentication token has expired")
	ErrInvalidToken       = errors.New("invalid authentication token")
)

// UserAuth models user credentials and identity.
type UserAuth struct {
	ID           string
	TenantID     string
	Email        string
	PasswordHash string
	Role         string
}

// Claims encapsulates Opslane JWT claims.
type Claims struct {
	TenantID string `json:"tenant_id"`
	UserID   string `json:"user_id"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// AuthResponse holds issued tokens.
type AuthResponse struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
	TenantID     string    `json:"tenant_id"`
	UserID       string    `json:"user_id"`
}

// PasswordManager encapsulates bcrypt hashing.
type PasswordManager struct {
	cost int
}

func NewPasswordManager(cost int) *PasswordManager {
	if cost <= 0 {
		cost = bcrypt.DefaultCost
	}
	return &PasswordManager{cost: cost}
}

func (pm *PasswordManager) Hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), pm.cost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(bytes), nil
}

func (pm *PasswordManager) Verify(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// TokenManager handles JWT issuance and verification.
type TokenManager struct {
	secretKey []byte
	issuer    string
	accessTTL time.Duration
}

func NewTokenManager(secret string, accessTTL time.Duration) *TokenManager {
	return &TokenManager{
		secretKey: []byte(secret),
		issuer:    "opslane-auth-service",
		accessTTL: accessTTL,
	}
}

func (tm *TokenManager) GenerateAccessToken(u *UserAuth) (string, time.Time, error) {
	expiresAt := time.Now().Add(tm.accessTTL)

	claims := &Claims{
		TenantID: u.TenantID,
		UserID:   u.ID,
		Email:    u.Email,
		Role:     u.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    tm.issuer,
			Subject:   u.ID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(tm.secretKey)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenStr, expiresAt, nil
}

func (tm *TokenManager) GenerateRefreshToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("failed to generate random token: %w", err)
	}
	return hex.EncodeToString(b), nil
}

func (tm *TokenManager) ValidateAccessToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return tm.secretKey, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

// AuthService coordinates authentication workflows.
type AuthService struct {
	users     map[string]*UserAuth
	passwords *PasswordManager
	tokens    *TokenManager
}

func NewAuthService(secret string) *AuthService {
	pm := NewPasswordManager(bcrypt.MinCost) // MinCost for test efficiency
	tm := NewTokenManager(secret, 15*time.Minute)

	hash, _ := pm.Hash("opslane_secret_123")
	users := map[string]*UserAuth{
		"admin@acme.corp": {
			ID:           "usr-01",
			TenantID:     "tenant-acme",
			Email:        "admin@acme.corp",
			PasswordHash: hash,
			Role:         "admin",
		},
	}

	return &AuthService{
		users:     users,
		passwords: pm,
		tokens:    tm,
	}
}

func (as *AuthService) Login(email, password string) (*AuthResponse, error) {
	user, exists := as.users[email]
	if !exists {
		return nil, ErrInvalidCredentials
	}

	if !as.passwords.Verify(password, user.PasswordHash) {
		return nil, ErrInvalidCredentials
	}

	tokenStr, expiresAt, err := as.tokens.GenerateAccessToken(user)
	if err != nil {
		return nil, err
	}

	refreshToken, err := as.tokens.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		AccessToken:  tokenStr,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
		TenantID:     user.TenantID,
		UserID:       user.ID,
	}, nil
}

func main() {
	svc := NewAuthService("opslane-super-secret-key-32-bytes!")
	res, err := svc.Login("admin@acme.corp", "opslane_secret_123")
	if err != nil {
		log.Fatalf("Login failed: %v", err)
	}

	fmt.Printf("Authenticated successfully! Issued JWT token for user %s (tenant %s)\n", res.UserID, res.TenantID)
}
