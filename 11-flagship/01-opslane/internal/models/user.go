// Copyright (c) 2026 Rasel Hossen
// See LICENSE for usage terms.

package models

import "time"

// UserRole defines the RBAC (Role-Based Access Control) level for a user within a tenant.
type UserRole string

const (
	// UserRoleAdmin grants full read/write access to all tenant resources.
	UserRoleAdmin UserRole = "admin"
	// UserRoleMember grants read/write access to standard operational resources (e.g., Orders).
	UserRoleMember UserRole = "member"
	// UserRoleBilling grants read-only access to operational resources, but full access to Payments.
	UserRoleBilling UserRole = "billing"
)

// User represents a human identity that has been granted access to a specific Tenant.
// Users are strictly tenant-scoped.
type User struct {
	ID           int64     `json:"id"`
	TenantID     int64     `json:"tenant_id"`
	Email        string    `json:"email"`
	DisplayName  string    `json:"display_name"`
	PasswordHash string    `json:"-"`
	Role         UserRole  `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
}
