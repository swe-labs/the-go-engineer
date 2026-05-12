// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

package models

import "time"

// User represents the domain model for a user in our system.
// User (Struct): represents the domain model for a user in our system.
type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"-"` // Omit from JSON for security
	CreatedAt time.Time `json:"created_at"`
	Profile   Profile   `json:"profile"`
}

// Profile represents the extended user information.
// Profile (Struct): represents the extended user information.
type Profile struct {
	UserID    int       `json:"user_id"`
	Avatar    string    `json:"avatar"`
	CreatedAt time.Time `json:"created_at"`
}
