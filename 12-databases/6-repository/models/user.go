// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

package models

import "time"

// User represent a customer in the database with a backing table of users
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	// 1. Reflection Ignore Tags
	// The `json:"-"` struct tag explicitly tells the `encoding/json` package to
	// NEVER serialize this field. This is critical for security to guarantee we
	// don't accidentally leak password hashes in API responses.
	HashedPassword string    `json:"-"`
	CreatedAt      time.Time `json:"created_at"`
	Profile        Profile   `json:"profile"`
}

// Profile belongs to a user
type Profile struct {
	UserID  int       `json:"user_id"`
	Avatar  string    `json:"avatar"`
	Created time.Time `json:"created"`
}
