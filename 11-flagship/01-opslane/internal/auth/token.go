// Copyright (c) 2026 Rasel Hossen
// See LICENSE for usage terms.

package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/swe-labs/the-go-engineer/11-flagship/01-opslane/internal/models"
)

var (
	// ErrInvalidToken (Error): signals a malformed, incorrectly signed, or invalid-claims token
	ErrInvalidToken = errors.New("invalid token")
	// ErrExpiredToken (Error): indicates the token is beyond its allowed lifetime
	ErrExpiredToken = errors.New("token expired")
)

// Identity (Struct): holds the trusted tenant-scoped user data carried after authentication
type Identity struct {
	UserID    int64           `json:"user_id"`
	TenantID  int64           `json:"tenant_id"`
	Email     string          `json:"email"`
	Role      models.UserRole `json:"role"`
	IssuedAt  time.Time       `json:"issued_at"`
	ExpiresAt time.Time       `json:"expires_at"`
}

// TokenManager (Struct): issues and verifies HMAC-signed JWT-compatible access tokens
//
// NOTE: This is a teaching JWT implementation to demonstrate cryptographic
// signatures and identity extraction without external dependencies.
// Production systems should usually use mature libraries or managed identity
// infrastructure.
type TokenManager struct {
	secret []byte
	issuer string
	ttl    time.Duration
	now    func() time.Time
}

// tokenHeader (Struct): internal JWT header with algorithm and type fields
type tokenHeader struct {
	Algorithm string `json:"alg"`
	Type      string `json:"typ"`
}

// tokenClaims (Struct): internal JWT payload with identity and expiry fields
type tokenClaims struct {
	Issuer    string          `json:"iss"`
	Subject   string          `json:"sub"`
	UserID    int64           `json:"uid"`
	TenantID  int64           `json:"tid"`
	Email     string          `json:"email"`
	Role      models.UserRole `json:"role"`
	IssuedAt  int64           `json:"iat"`
	ExpiresAt int64           `json:"exp"`
}

// NewTokenManager (Constructor): builds a token manager from validated runtime auth settings
func NewTokenManager(secret, issuer string, ttl time.Duration) (*TokenManager, error) {
	if len(secret) < 32 {
		return nil, fmt.Errorf("token secret must be at least 32 characters")
	}

	if strings.TrimSpace(issuer) == "" {
		return nil, fmt.Errorf("token issuer must not be empty")
	}

	if ttl <= 0 {
		return nil, fmt.Errorf("token ttl must be positive")
	}

	return &TokenManager{
		secret: []byte(secret),
		issuer: issuer,
		ttl:    ttl,
		now:    func() time.Time { return time.Now().UTC() },
	}, nil
}

// Issue (Method): creates a signed JWT token for a validated tenant-scoped identity
func (m *TokenManager) Issue(identity Identity) (string, error) {
	if identity.UserID <= 0 {
		return "", fmt.Errorf("identity user id must be positive")
	}

	if identity.TenantID <= 0 {
		return "", fmt.Errorf("identity tenant id must be positive")
	}

	if strings.TrimSpace(identity.Email) == "" {
		return "", fmt.Errorf("identity email must not be empty")
	}

	if strings.TrimSpace(string(identity.Role)) == "" {
		return "", fmt.Errorf("identity role must not be empty")
	}

	issuedAt := m.now()
	expiresAt := issuedAt.Add(m.ttl)
	claims := tokenClaims{
		Issuer:    m.issuer,
		Subject:   strconv.FormatInt(identity.UserID, 10),
		UserID:    identity.UserID,
		TenantID:  identity.TenantID,
		Email:     identity.Email,
		Role:      identity.Role,
		IssuedAt:  issuedAt.Unix(),
		ExpiresAt: expiresAt.Unix(),
	}

	headerSegment, err := encodeSegment(tokenHeader{Algorithm: "HS256", Type: "JWT"})
	if err != nil {
		return "", fmt.Errorf("encode token header: %w", err)
	}

	payloadSegment, err := encodeSegment(claims)
	if err != nil {
		return "", fmt.Errorf("encode token claims: %w", err)
	}

	signingInput := headerSegment + "." + payloadSegment
	signature := sign(signingInput, m.secret)

	return signingInput + "." + base64.RawURLEncoding.EncodeToString(signature), nil
}

// Verify (Method): checks token signature, issuer, expiry, and required tenant identity claims
func (m *TokenManager) Verify(token string) (Identity, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return Identity{}, ErrInvalidToken
	}

	signingInput := parts[0] + "." + parts[1]
	wantSignature := sign(signingInput, m.secret)

	gotSignature, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil {
		return Identity{}, ErrInvalidToken
	}

	if subtle.ConstantTimeCompare(gotSignature, wantSignature) != 1 {
		return Identity{}, ErrInvalidToken
	}

	var header tokenHeader
	if err := decodeSegment(parts[0], &header); err != nil {
		return Identity{}, ErrInvalidToken
	}

	if header.Algorithm != "HS256" || header.Type != "JWT" {
		return Identity{}, ErrInvalidToken
	}

	var claims tokenClaims
	if err := decodeSegment(parts[1], &claims); err != nil {
		return Identity{}, ErrInvalidToken
	}

	if claims.Issuer != m.issuer {
		return Identity{}, ErrInvalidToken
	}

	if claims.UserID <= 0 || claims.TenantID <= 0 || claims.Subject == "" {
		return Identity{}, ErrInvalidToken
	}

	if claims.Subject != strconv.FormatInt(claims.UserID, 10) {
		return Identity{}, ErrInvalidToken
	}

	if strings.TrimSpace(claims.Email) == "" || strings.TrimSpace(string(claims.Role)) == "" {
		return Identity{}, ErrInvalidToken
	}

	expiresAt := time.Unix(claims.ExpiresAt, 0).UTC()
	if !m.now().Before(expiresAt) {
		return Identity{}, ErrExpiredToken
	}

	return Identity{
		UserID:    claims.UserID,
		TenantID:  claims.TenantID,
		Email:     claims.Email,
		Role:      claims.Role,
		IssuedAt:  time.Unix(claims.IssuedAt, 0).UTC(),
		ExpiresAt: expiresAt,
	}, nil
}

// encodeSegment (Function): base64url-encodes a JSON-marshalled value for JWT segment
func encodeSegment(value any) (string, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(data), nil
}

// decodeSegment (Function): base64url-decodes and JSON-unmarshals a JWT segment into value
func decodeSegment(segment string, value any) error {
	data, err := base64.RawURLEncoding.DecodeString(segment)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, value)
}

// sign (Function): computes an HMAC-SHA256 signature for JWT signing input
func sign(input string, secret []byte) []byte {
	mac := hmac.New(sha256.New, secret)
	mac.Write([]byte(input))
	return mac.Sum(nil)
}
