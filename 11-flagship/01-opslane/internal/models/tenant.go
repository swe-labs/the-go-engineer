// Copyright (c) 2026 Rasel Hossen
// See LICENSE for usage terms.

package models

import "time"

// Tenant (Struct): is the root isolation boundary in the SaaS platform.
// Every other domain entity (Users, Orders, Payments) MUST belong to exactly one Tenant.
//
// Boundary: Tenant is the top-level parent entity; all other entities reference it via TenantID.
// Invariant: Slug must be unique across all tenants to support subdomain-based routing.
type Tenant struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	CreatedAt time.Time `json:"created_at"`
}
