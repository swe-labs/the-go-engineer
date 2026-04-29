// Copyright (c) 2026 Rasel Hossen
// See LICENSE for usage terms.

package models

import "time"

// Tenant is the root isolation boundary in the SaaS platform.
// Every other domain entity (Users, Orders, Payments) MUST belong to exactly one Tenant.
type Tenant struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	CreatedAt time.Time `json:"created_at"`
}
