// Copyright (c) 2026 Rasel Hossen
// See LICENSE for usage terms.

package auth

import (
	"context"
	"fmt"
	"strings"

	"github.com/rasel9t6/the-go-engineer/11-flagship/01-opslane/internal/models"
)

// UserLookup is the minimum repository behavior auth needs to verify credentials.
type UserLookup interface {
	GetUserByEmail(ctx context.Context, tenantID int64, email string) (models.User, error)
}

// LoginRequest carries the tenant-scoped credentials required to issue a token.
type LoginRequest struct {
	TenantID int64
	Email    string
	Password string
}

// LoginResult is the safe auth response returned after successful authentication.
type LoginResult struct {
	Token    string
	Identity Identity
}

// Service verifies credentials and issues tenant-scoped access tokens.
type Service struct {
	users  UserLookup
	tokens *TokenManager
}

// NewService wires the auth service to its persistence and token dependencies.
func NewService(users UserLookup, tokens *TokenManager) *Service {
	return &Service{
		users:  users,
		tokens: tokens,
	}
}

// Login verifies a tenant-scoped user password and returns a signed access token.
func (s *Service) Login(ctx context.Context, req LoginRequest) (LoginResult, error) {
	if s == nil || s.users == nil || s.tokens == nil {
		return LoginResult{}, fmt.Errorf("auth service is not configured")
	}

	if req.TenantID <= 0 || strings.TrimSpace(req.Email) == "" || req.Password == "" {
		return LoginResult{}, ErrInvalidCredentials
	}

	user, err := s.users.GetUserByEmail(ctx, req.TenantID, req.Email)
	if err != nil {
		return LoginResult{}, ErrInvalidCredentials
	}

	if err := VerifyPassword(user.PasswordHash, req.Password); err != nil {
		return LoginResult{}, ErrInvalidCredentials
	}

	if user.TenantID != req.TenantID {
		return LoginResult{}, ErrInvalidCredentials
	}

	identity := Identity{
		UserID:   user.ID,
		TenantID: user.TenantID,
		Email:    user.Email,
		Role:     user.Role,
	}

	token, err := s.tokens.Issue(identity)
	if err != nil {
		return LoginResult{}, fmt.Errorf("issue auth token: %w", err)
	}

	verifiedIdentity, err := s.tokens.Verify(token)
	if err != nil {
		return LoginResult{}, fmt.Errorf("verify issued auth token: %w", err)
	}

	return LoginResult{
		Token:    token,
		Identity: verifiedIdentity,
	}, nil
}
