// Copyright (c) 2026 Rasel Hossen
// See LICENSE for usage terms.

package db

import (
	"context"
	"database/sql"
	"fmt"
)

var schemaStatements = []string{
	`CREATE TABLE IF NOT EXISTS tenants (
		id BIGSERIAL PRIMARY KEY,
		name TEXT NOT NULL,
		slug TEXT NOT NULL UNIQUE,
		created_at TIMESTAMPTZ NOT NULL
	);`,
	`CREATE TABLE IF NOT EXISTS users (
		id BIGSERIAL PRIMARY KEY,
		tenant_id BIGINT NOT NULL,
		email TEXT NOT NULL,
		display_name TEXT NOT NULL,
		password_hash TEXT NOT NULL,
		role TEXT NOT NULL,
		created_at TIMESTAMPTZ NOT NULL,
		UNIQUE (tenant_id, id),
		UNIQUE (tenant_id, email),
		FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE
	);`,
	`CREATE TABLE IF NOT EXISTS orders (
		id BIGSERIAL PRIMARY KEY,
		tenant_id BIGINT NOT NULL,
		user_id BIGINT NOT NULL,
		status TEXT NOT NULL,
		total_cents BIGINT NOT NULL,
		currency TEXT NOT NULL,
		idempotency_key TEXT NOT NULL,
		created_at TIMESTAMPTZ NOT NULL,
		updated_at TIMESTAMPTZ NOT NULL,
		UNIQUE (tenant_id, id),
		UNIQUE (tenant_id, idempotency_key),
		FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
		FOREIGN KEY (tenant_id, user_id) REFERENCES users(tenant_id, id) ON DELETE NO ACTION
	);`,
	`CREATE TABLE IF NOT EXISTS payments (
		id BIGSERIAL PRIMARY KEY,
		tenant_id BIGINT NOT NULL,
		order_id BIGINT NOT NULL,
		status TEXT NOT NULL,
		provider_reference TEXT NOT NULL,
		amount_cents BIGINT NOT NULL,
		failure_reason TEXT NOT NULL DEFAULT '',
		created_at TIMESTAMPTZ NOT NULL,
		updated_at TIMESTAMPTZ NOT NULL,
		UNIQUE (tenant_id, provider_reference),
		FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
		FOREIGN KEY (tenant_id, order_id) REFERENCES orders(tenant_id, id) ON DELETE CASCADE
	);`,
	`CREATE INDEX IF NOT EXISTS idx_orders_tenant_status ON orders(tenant_id, status);`,
	`CREATE INDEX IF NOT EXISTS idx_orders_tenant_user ON orders(tenant_id, user_id);`,
	`CREATE INDEX IF NOT EXISTS idx_orders_tenant_created_at ON orders(tenant_id, created_at DESC);`,
	`CREATE INDEX IF NOT EXISTS idx_payments_tenant_order ON payments(tenant_id, order_id);`,
	`CREATE INDEX IF NOT EXISTS idx_payments_tenant_status ON payments(tenant_id, status);`,
}

func Migrate(ctx context.Context, database *sql.DB) error {
	for _, statement := range schemaStatements {
		if _, err := database.ExecContext(ctx, statement); err != nil {
			return fmt.Errorf("apply migration: %w", err)
		}
	}

	return nil
}
