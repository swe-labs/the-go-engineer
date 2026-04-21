// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

package models

import "time"

// ============================================================================
// Internal Package: Models
// Level: Advanced
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Defining cross-boundary struct definitions
//
// ENGINEERING DEPTH:
//   Models belong in their own package (`models`) so they can be imported by
//   the `handlers` and `repository` packages WITHOUT causing a circular dependency!
//   This is a master-level requirement for Go Package Design.
// ============================================================================

type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"-"` // DO NOT serialize password hashes over JSON!
	CreatedAt time.Time `json:"created_at"`
}

type Post struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	AuthorID  int       `json:"author_id"`
	CreatedAt time.Time `json:"created_at"`
}
