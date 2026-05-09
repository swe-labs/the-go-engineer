// Copyright (c) 2026 Rasel Hossen
// See LICENSE for usage terms.

package models

import "time"

// Package models provides the core domain entity types for the Opslane backend.
// All entities are tenant-scoped to ensure multi-tenant isolation.

// UserRole (Type): defines the RBAC (Role-Based Access Control) level for a user within a tenant.
// It controls what resources a user can access and what operations they can perform.
type UserRole string

const (
	// UserRoleAdmin grants full read/write access to all tenant resources.
	UserRoleAdmin UserRole = "admin"
	// UserRoleMember grants read/write access to standard operational resources (e.g., Orders).
	UserRoleMember UserRole = "member"
	// UserRoleBilling grants read-only access to operational resources, but full access to Payments.
	UserRoleBilling UserRole = "billing"
)

// User (Struct): represents a human identity that has been granted access to a specific Tenant.
// Users are strictly tenant-scoped - they cannot access resources outside their TenantID.
//
// Invariant: Every User must have a valid TenantID; the database enforces this via foreign key.
type User struct {
	ID           int64     `json:"id"`
	TenantID     int64     `json:"tenant_id"`
	Email        string    `json:"email"`
	DisplayName  string    `json:"display_name"`
	PasswordHash string    `json:"-"`
	Role         UserRole  `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
}
