package main

import (
	"context"
	"fmt"
	"path/filepath"
	"sync/atomic"
	"testing"
	"time"
)

var testDBCounter atomic.Uint64

func uniqueTestDSN(t *testing.T) string {
	return fmt.Sprintf("file:test_%d?mode=memory&cache=shared&_pragma=foreign_keys(1)", testDBCounter.Add(1))
}

func setupTestDB(t *testing.T) *OrderRepository {
	t.Helper()
	db, err := ConnectDatabase(uniqueTestDSN(t))
	if err != nil {
		t.Fatalf("failed to open in-memory db: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })

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
	db, err := ConnectDatabase(uniqueTestDSN(t))
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

func TestOrderRepository_ForeignKeyEnforcement(t *testing.T) {
	repo := setupTestDB(t)
	ctx := context.Background()
	now := time.Now().UTC()

	// Attempting to insert an order referencing a non-existent tenant/user must fail FK constraint
	orphanOrder := &Order{
		ID:          "ord-orphan-999",
		TenantID:    "nonexistent-tenant",
		UserID:      "nonexistent-user",
		TotalAmount: 100.0,
		Currency:    "USD",
		Status:      "PENDING",
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	err := repo.Create(ctx, orphanOrder)
	if err == nil {
		t.Fatalf("expected foreign key constraint violation inserting orphan order, got nil")
	}
}

func TestConnectDatabase_InMemoryIsolation(t *testing.T) {
	// Two independent :memory: connections must NOT share state: the driver
	// rewrites each to a uniquely named shared-cache URI, so tables created
	// in one are invisible to the other.
	dbA, err := ConnectDatabase(":memory:")
	if err != nil {
		t.Fatalf("open A: %v", err)
	}
	defer dbA.Close()

	dbB, err := ConnectDatabase(":memory:")
	if err != nil {
		t.Fatalf("open B: %v", err)
	}
	defer dbB.Close()

	if _, err := dbA.Exec("CREATE TABLE marker_a (id INTEGER PRIMARY KEY)"); err != nil {
		t.Fatalf("create table in A: %v", err)
	}

	var countA, countB int
	if err := dbA.QueryRow("SELECT COUNT(1) FROM sqlite_master WHERE name = 'marker_a'").Scan(&countA); err != nil {
		t.Fatalf("probe A: %v", err)
	}
	if err := dbB.QueryRow("SELECT COUNT(1) FROM sqlite_master WHERE name = 'marker_a'").Scan(&countB); err != nil {
		t.Fatalf("probe B: %v", err)
	}
	if countA != 1 || countB != 0 {
		t.Fatalf("SECURITY/ISOLATION VIOLATION: expected marker_a only in A (A=%d, B=%d)", countA, countB)
	}
}

func TestConnectDatabase_FileBackedForeignKeyEnforcement(t *testing.T) {
	// A plain filesystem-path DSN is normalized to a file: URI with the
	// foreign_keys pragma, so declared REFERENCES constraints fire there too.
	db, err := ConnectDatabase(filepath.Join(t.TempDir(), "opslane_file_fk.db"))
	if err != nil {
		t.Fatalf("open file-backed db: %v", err)
	}
	defer db.Close()

	if _, err := db.Exec("CREATE TABLE parent (id TEXT PRIMARY KEY); CREATE TABLE child (pid TEXT REFERENCES parent(id))"); err != nil {
		t.Fatalf("create schema: %v", err)
	}
	if _, err := db.Exec("INSERT INTO child (pid) VALUES ('missing-parent')"); err == nil {
		t.Fatalf("expected foreign key constraint violation on file-backed DSN, got nil")
	}
}

func TestConnectDatabase_DisablingPragmaIsOverridden(t *testing.T) {
	// ConnectDatabase guarantees foreign-key enforcement: a caller-supplied
	// _pragma=foreign_keys(0) is followed by the appended (1), and the driver
	// applies DSN pragmas in query order, so enforcement stays on.
	db, err := ConnectDatabase("file:" + filepath.ToSlash(filepath.Join(t.TempDir(), "opslane_fk0.db")) + "?_pragma=foreign_keys(0)")
	if err != nil {
		t.Fatalf("open db with disabling pragma: %v", err)
	}
	defer db.Close()

	if _, err := db.Exec("CREATE TABLE parent (id TEXT PRIMARY KEY); CREATE TABLE child (pid TEXT REFERENCES parent(id))"); err != nil {
		t.Fatalf("create schema: %v", err)
	}
	if _, err := db.Exec("INSERT INTO child (pid) VALUES ('missing-parent')"); err == nil {
		t.Fatalf("expected foreign key enforcement despite caller-supplied foreign_keys(0)")
	}
}
