package main

import (
	"context"
	"database/sql"
	"errors"
	"testing"
)

func setupTestStore(t *testing.T) *IsolatedStore {
	t.Helper()
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}

	store, err := NewIsolatedStore(db)
	if err != nil {
		t.Fatalf("failed to init isolated store: %v", err)
	}
	return store
}

func TestTenantIsolation_MissingContext(t *testing.T) {
	store := setupTestStore(t)
	rawCtx := context.Background()

	_, err := store.Create(rawCtx, "ord-1", "Orphaned Item", 50.0)
	if !errors.Is(err, ErrMissingTenantContext) {
		t.Fatalf("expected ErrMissingTenantContext, got: %v", err)
	}

	_, err = store.List(rawCtx)
	if !errors.Is(err, ErrMissingTenantContext) {
		t.Fatalf("expected ErrMissingTenantContext on List, got: %v", err)
	}
}

func TestTenantIsolation_CrossTenantQuery(t *testing.T) {
	// Core practice task requirement:
	// Verify resources created for one tenant are completely invisible to queries scoped to another tenant.
	store := setupTestStore(t)

	ctxTenantA := ContextWithTenant(context.Background(), "tenant-alpha")
	ctxTenantB := ContextWithTenant(context.Background(), "tenant-beta")

	// 1. Tenant Alpha creates an order
	_, err := store.Create(ctxTenantA, "ord-alpha-100", "Confidential Alpha Blueprint", 15000.0)
	if err != nil {
		t.Fatalf("failed to create order for tenant A: %v", err)
	}

	// 2. Tenant Alpha can read it
	orderA, err := store.GetByID(ctxTenantA, "ord-alpha-100")
	if err != nil || orderA == nil {
		t.Fatalf("expected tenant A to read own order, got order: %v, err: %v", orderA, err)
	}

	// 3. Tenant Beta queries the same order ID -> MUST RETURN NIL (Data boundary)
	orderB, err := store.GetByID(ctxTenantB, "ord-alpha-100")
	if err != nil {
		t.Fatalf("expected nil error on query, got: %v", err)
	}
	if orderB != nil {
		t.Fatalf("CRITICAL SECURITY LEAK: Tenant Beta read Tenant Alpha's data: %+v", orderB)
	}
}

func TestTenantIsolation_CrossTenantList(t *testing.T) {
	store := setupTestStore(t)

	ctxA := ContextWithTenant(context.Background(), "tenant-a")
	ctxB := ContextWithTenant(context.Background(), "tenant-b")

	// Seed 2 orders for Tenant A and 1 for Tenant B
	_, _ = store.Create(ctxA, "a-1", "Item A1", 10.0)
	_, _ = store.Create(ctxA, "a-2", "Item A2", 20.0)
	_, _ = store.Create(ctxB, "b-1", "Item B1", 99.0)

	listA, err := store.List(ctxA)
	if err != nil {
		t.Fatalf("failed to list orders for tenant A: %v", err)
	}
	if len(listA) != 2 {
		t.Errorf("expected tenant A to see 2 orders, got %d", len(listA))
	}
	for _, o := range listA {
		if o.TenantID != "tenant-a" {
			t.Errorf("tenant A list contains foreign tenant ID: %s", o.TenantID)
		}
	}

	listB, err := store.List(ctxB)
	if err != nil {
		t.Fatalf("failed to list orders for tenant B: %v", err)
	}
	if len(listB) != 1 || listB[0].ID != "b-1" {
		t.Errorf("expected tenant B to see only order b-1, got: %v", listB)
	}
}

func TestTenantIsolation_CrossTenantMutationAttempt(t *testing.T) {
	store := setupTestStore(t)

	ctxA := ContextWithTenant(context.Background(), "tenant-a")
	ctxB := ContextWithTenant(context.Background(), "tenant-b")

	// Tenant A creates an item
	_, _ = store.Create(ctxA, "target-order", "Original Title", 100.0)

	// Tenant B attempts to UPDATE Tenant A's item -> Must affect 0 rows
	updated, err := store.Update(ctxB, "target-order", "Malicious Injected Title", 0.0)
	if err != nil {
		t.Fatalf("unexpected update error: %v", err)
	}
	if updated {
		t.Fatalf("SECURITY VIOLATION: Tenant B successfully updated Tenant A's record!")
	}

	// Verify Tenant A's item remains unmodified
	original, _ := store.GetByID(ctxA, "target-order")
	if original.Title != "Original Title" || original.Amount != 100.0 {
		t.Errorf("expected original data to be intact, got: %+v", original)
	}

	// Tenant B attempts to DELETE Tenant A's item -> Must affect 0 rows
	deleted, err := store.Delete(ctxB, "target-order")
	if err != nil {
		t.Fatalf("unexpected delete error: %v", err)
	}
	if deleted {
		t.Fatalf("SECURITY VIOLATION: Tenant B successfully deleted Tenant A's record!")
	}
}
