package auth

import (
	"errors"
	"testing"
)

func TestHashPasswordVerifiesOriginalPassword(t *testing.T) {
	t.Parallel()

	hash, err := HashPassword("CorrectHorse7Battery")
	if err != nil {
		t.Fatalf("HashPassword returned error: %v", err)
	}

	if hash == "CorrectHorse7Battery" {
		t.Fatal("password hash must not store the raw password")
	}

	if err := VerifyPassword(hash, "CorrectHorse7Battery"); err != nil {
		t.Fatalf("VerifyPassword returned error: %v", err)
	}
}

func TestVerifyPasswordRejectsWrongPassword(t *testing.T) {
	t.Parallel()

	hash, err := HashPassword("CorrectHorse7Battery")
	if err != nil {
		t.Fatalf("HashPassword returned error: %v", err)
	}

	err = VerifyPassword(hash, "WrongHorse7Battery")
	if !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("VerifyPassword error = %v, want ErrInvalidCredentials", err)
	}
}

func TestValidatePasswordRejectsWeakPassword(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		password string
	}{
		{name: "too short", password: "short7"},
		{name: "missing number", password: "correcthorsebattery"},
		{name: "missing letter", password: "123456789012"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePassword(tt.password)
			if !errors.Is(err, ErrWeakPassword) {
				t.Fatalf("ValidatePassword error = %v, want ErrWeakPassword", err)
			}
		})
	}
}
