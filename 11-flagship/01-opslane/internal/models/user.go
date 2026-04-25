// Copyright (c) 2026 Rasel Hossen
// See LICENSE for usage terms.

package models

import "time"

type UserRole string

const (
	UserRoleAdmin   UserRole = "admin"
	UserRoleMember  UserRole = "member"
	UserRoleBilling UserRole = "billing"
)

type User struct {
	ID           int64     `json:"id"`
	TenantID     int64     `json:"tenant_id"`
	Email        string    `json:"email"`
	DisplayName  string    `json:"display_name"`
	PasswordHash string    `json:"-"`
	Role         UserRole  `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
}
