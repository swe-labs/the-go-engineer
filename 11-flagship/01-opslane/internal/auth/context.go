// Copyright (c) 2026 Rasel Hossen
// See LICENSE for usage terms.

package auth

import (
	"context"
	"errors"
)

// ErrMissingIdentity means protected code ran before auth middleware attached identity.
var ErrMissingIdentity = errors.New("missing authenticated identity")

type identityContextKey struct{}

// WithIdentity stores trusted auth identity on a request context.
func WithIdentity(ctx context.Context, identity Identity) context.Context {
	return context.WithValue(ctx, identityContextKey{}, identity)
}

// IdentityFromContext returns the trusted auth identity if middleware attached one.
func IdentityFromContext(ctx context.Context) (Identity, bool) {
	identity, ok := ctx.Value(identityContextKey{}).(Identity)
	return identity, ok
}

// RequireIdentity returns the trusted identity or an explicit error for protected paths.
func RequireIdentity(ctx context.Context) (Identity, error) {
	identity, ok := IdentityFromContext(ctx)
	if !ok {
		return Identity{}, ErrMissingIdentity
	}

	return identity, nil
}
