// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// RUN: go run ./15-code-generation/2-mockery
package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/rasel9t6/the-go-engineer/15-code-generation/2-mockery/internal/storage"
)

// ============================================================================
// Section 15: Code Generation — Mockery
// Level: Expert
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - How to use mockery for automatic mock generation
//   - Integrating mockery with //go:generate
//   - The benefits of "mocking at the consumer's request"
//   - Improving developer experience (DX) through automated stubs
//
// ENGINEERING DEPTH:
//   Manual mocking (seen in Section 13) is great for small projects, but scales
//   poorly. Every time an interface changes, you must manually update every mock.
//   Mockery solves this by parsing your Go code and generating type-safe mock
//   structs that integrate with testify/mock.
//
//   THE PATTERN:
//     1. Define an interface: `type Storer interface { ... }`
//     2. Add the directive: `//go:generate mockery --name=Storer`
//     3. Run the generator: `go generate ./...`
//     4. Use the mock: `m := new(mocks.Storer)`
//
// RUN: go run ./15-code-generation/2-mockery
// ============================================================================

// UserManager uses a Storer to manage users.
type UserManager struct {
	store storage.Storer
}

func (m *UserManager) WelcomeUser(ctx context.Context, id string) (string, error) {
	u, err := m.store.GetUser(ctx, id)
	if err != nil {
		return "", fmt.Errorf("failed to get user: %w", err)
	}

	return fmt.Sprintf("Welcome, %s!", u.Email), nil
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	logger.Info("Mockery Lesson")
	fmt.Println("\nTo generate a mock for the Storer interface, run:")
	fmt.Println("  go generate ./15-code-generation/2-mockery/...")

	fmt.Println("\nThis command triggers the //go:generate mockery directive")
	fmt.Println("found in internal/storage/storage.go.")

	fmt.Println("\nOnce generated, you can use the mock in your tests like this:")
	fmt.Println(`
    func TestWelcomeUser(t *testing.T) {
        m := new(mocks.Storer)
        m.On("GetUser", mock.Anything, "42").Return(&storage.User{Email: "test@example.com"}, nil)

        mgr := &UserManager{store: m}
        msg, _ := mgr.WelcomeUser(context.Background(), "42")

        assert.Equal(t, "Welcome, test@example.com!", msg)
        m.AssertExpectations(t)
    }`)

	fmt.Println("\n---------------------------------------------------")
	fmt.Println("🚀 NEXT UP: CG.3 sqlc (restored expert content)")
	fmt.Println("   Current: CG.2 (mockery)")
	fmt.Println("---------------------------------------------------")
}
