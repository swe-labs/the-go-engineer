package db

import (
	"strings"
	"testing"
)

func TestSchemaStatementsCoverCoreOpslaneTables(t *testing.T) {
	t.Parallel()

	wantSnippets := []string{
		"CREATE TABLE IF NOT EXISTS tenants",
		"CREATE TABLE IF NOT EXISTS users",
		"CREATE TABLE IF NOT EXISTS orders",
		"CREATE TABLE IF NOT EXISTS payments",
		"UNIQUE (tenant_id, id)",
		"FOREIGN KEY (tenant_id, user_id) REFERENCES users(tenant_id, id) ON DELETE NO ACTION",
		"FOREIGN KEY (tenant_id, order_id) REFERENCES orders(tenant_id, id)",
		"idx_orders_tenant_status",
		"idx_orders_tenant_user",
		"idx_payments_tenant_order",
	}
	redundantSnippets := []string{
		"idx_tenants_slug",
		"idx_users_tenant_email",
		"ON DELETE RESTRICT",
	}

	joined := ""
	for _, statement := range schemaStatements {
		joined += statement + "\n"
	}

	for _, snippet := range wantSnippets {
		if !strings.Contains(joined, snippet) {
			t.Fatalf("schema statements missing %q", snippet)
		}
	}

	for _, snippet := range redundantSnippets {
		if strings.Contains(joined, snippet) {
			t.Fatalf("schema statements contain redundant index %q", snippet)
		}
	}
}
