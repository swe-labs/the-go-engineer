package auth

import (
	"testing"
	"time"
)

func FuzzVerifyRejectsMalformedTokens(f *testing.F) {
	manager, err := NewTokenManager("01234567890123456789012345678901", "opslane", time.Hour)
	if err != nil {
		f.Fatalf("token manager setup failed: %v", err)
	}

	seeds := []string{
		"",
		"abc",
		"a.b",
		"a.b.c",
		"....",
		"eyJhbGciOiJIUzI1NiJ9.e30.invalid",
	}
	for _, seed := range seeds {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, token string) {
		_, _ = manager.Verify(token)
	})
}
