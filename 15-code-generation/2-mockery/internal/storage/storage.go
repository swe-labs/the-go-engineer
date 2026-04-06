// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

package storage

import "context"

// User represents a system user.
type User struct {
	ID    string
	Email string
}

// Storer is the interface that Mockery will generate a mock for.
// In Go, we mock at the consumer's request. By using //go:generate
// we ensure the mock is always up to date with the interface.
//
//go:generate mockery --name=Storer --output=../../mocks --case=underscore
type Storer interface {
	GetUser(ctx context.Context, id string) (*User, error)
	SaveUser(ctx context.Context, u *User) error
}
