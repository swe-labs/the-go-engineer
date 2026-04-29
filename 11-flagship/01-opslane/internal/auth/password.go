// Copyright (c) 2026 Rasel Hossen
// See LICENSE for usage terms.

package auth

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"

	"golang.org/x/crypto/bcrypt"
)

var (
	// ErrInvalidCredentials hides whether the email, password, or token detail was wrong.
	ErrInvalidCredentials = errors.New("invalid credentials")
	// ErrWeakPassword tells callers that the proposed password failed the local policy.
	ErrWeakPassword = errors.New("password does not meet policy")
)

const minimumPasswordLength = 12

// ValidatePassword enforces the minimum local password policy before hashing.
func ValidatePassword(password string) error {
	if utf8.RuneCountInString(password) < minimumPasswordLength {
		return fmt.Errorf("%w: must be at least %d characters", ErrWeakPassword, minimumPasswordLength)
	}

	var hasLetter bool
	var hasDigit bool
	for _, r := range password {
		if unicode.IsLetter(r) {
			hasLetter = true
		}
		if unicode.IsDigit(r) {
			hasDigit = true
		}
	}

	if !hasLetter || !hasDigit {
		return fmt.Errorf("%w: must include at least one letter and one number", ErrWeakPassword)
	}

	return nil
}

// HashPassword validates and stores a password using bcrypt.
func HashPassword(password string) (string, error) {
	if err := ValidatePassword(password); err != nil {
		return "", err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("hash password: %w", err)
	}

	return string(hash), nil
}

// VerifyPassword compares a bcrypt hash with a candidate password.
func VerifyPassword(hash, password string) error {
	if strings.TrimSpace(hash) == "" || password == "" {
		return ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return ErrInvalidCredentials
	}

	return nil
}
