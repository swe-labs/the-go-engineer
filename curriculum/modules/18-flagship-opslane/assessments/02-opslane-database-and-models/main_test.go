package main

import (
	"context"
	"testing"
	"time"
)

func setupTestDB(t *testing.T) *OrderRepository {
	t.Helper()
	db, err := ConnectDatabase(":memory:")
	if err != nil {
		t.Fatalf("failed to open in-memory db: %v", err)
	}

	if err := RunMigrations(db, SchemaMigrations); err != nil {
		t.Fatalf("failed to run migrations: %v", err)
	}

	// Seed required parent tenant and user for foreign key constraints
	now := time.Now().UTC()
	_, err = db.Exec(`
		INSERT INTO tenants (id, name, tier, created_at) VALUES ('tenant-acme', 'Acme Corp', 'enterprise', ?);
		INSERT INTO users (id, tenant_id, email, password_hash, role, created_at) 
		VALUES ('user-alice', 'tenant-acme', 'alice@acme.corp', 'hash_stub', 'owner', ?);
	`, now, now)
	if err != nil {
		t.Fatalf("failed to seed test fixture: %v", err)
	}

	return NewOrderRepository(db)
}

func TestMigrations_RunAndIdempotency(t *testing.T) {
	db, err := ConnectDatabase(":memory:")
	if err != nil {
		t.Fatalf("failed to connect: %v", err)
	}
	defer db.Close()

	// First pass
	if err := RunMigrations(db, SchemaMigrations); err != nil {
		t.Fatalf("first migration pass failed: %v", err)
	}

	// Second pass should be a no-op and succeed without duplicate errors
	if err := RunMigrations(db, SchemaMigrations); err != nil {
		t.Fatalf("second idempotent migration pass failed: %v", err)
	}

	var count int
	err = db.QueryRow("SELECT COUNT(1) FROM schema_migrations").Scan(&count)
	if err != nil || count != len(SchemaMigrations) {
		t.Errorf("expected %d migrations recorded, got %d", len(SchemaMigrations), count)
	}
}

func TestOrderRepository_CreateAndGet(t *testing.T) {
	repo := setupTestDB(t)
	ctx := context.Background()
	now := time.Now().UTC().Truncate(time.Second)

	order := &Order{
		ID:          "ord-1001",
		TenantID:    "tenant-acme",
		UserID:      "user-alice",
		TotalAmount: 149.50,
		Currency:    "USD",
		Status:      "PENDING",
		Notes:       "Expedited production order",
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := repo.Create(ctx, order); err != nil {
		t.Fatalf("failed to create order: %v", err)
	}

	fetched, err := repo.GetByID(ctx, "tenant-acme", "ord-1001")
	if err != nil {
		t.Fatalf("failed to get order: %v", err)
	}
	if fetched == nil {
		t.Fatalf("expected order to be found, got nil")
	}

	if fetched.ID != "ord-1001" || fetched.TotalAmount != 149.50 || fetched.Notes != "Expedited production order" {
		t.Errorf("unexpected fetched order data: %+v", fetched)
	}
}

func TestOrderRepository_TenantIsolation(t *testing.T) {
	repo := setupTestDB(t)
	ctx := context.Background()
	now := time.Now().UTC()

	order := &Order{
		ID:          "ord-secret-200",
		TenantID:    "tenant-acme",
		UserID:      "user-alice",
		TotalAmount: 5000.0,
		Currency:    "USD",
		Status:      "PENDING",
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := repo.Create(ctx, order); err != nil {
		t.Fatalf("failed to create order: %v", err)
	}

	// Tenant 'tenant-evil' querying for tenant-acme's order ID must receive nil
	unauthorized, err := repo.GetByID(ctx, "tenant-evil", "ord-secret-200")
	if err != nil {
		t.Fatalf("expected nil error on missing row, got: %v", err)
	}
	if unauthorized != nil {
		t.Fatalf("SECURITY VIOLATION: Cross-tenant data leak! Expected nil, got: %+v", unauthorized)
	}
}
