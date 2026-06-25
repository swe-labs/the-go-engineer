package main

import (
	"context"
	"testing"
)

func TestCreateDocument_MissingTenant(t *testing.T) {
	store := NewTenantStore()
	_, err := store.AddDocument(context.Background(), "title", "content")
	if err == nil {
		t.Fatal("expected error when tenant not in context")
	}
}

func TestGetDocuments_TenantIsolation(t *testing.T) {
	store := NewTenantStore()
	ctxA := context.WithValue(context.Background(), tenantKey, "tenant-a")
	ctxB := context.WithValue(context.Background(), tenantKey, "tenant-b")

	store.AddDocument(ctxA, "Doc A1", "content")
	store.AddDocument(ctxA, "Doc A2", "content")
	store.AddDocument(ctxB, "Doc B1", "content")

	docsA, _ := store.GetDocuments(ctxA)
	if len(docsA) != 2 {
		t.Errorf("expected 2 docs for tenant-a, got %d", len(docsA))
	}

	docsB, _ := store.GetDocuments(ctxB)
	if len(docsB) != 1 {
		t.Errorf("expected 1 doc for tenant-b, got %d", len(docsB))
	}

	// Verify tenant A can't see B's docs.
	for _, d := range docsA {
		if d.TenantID != "tenant-a" {
			t.Errorf("tenant-a got doc from tenant %s", d.TenantID)
		}
	}
}

func TestGetDocumentsNoFilter_CrossTenantLeak(t *testing.T) {
	store := NewTenantStore()
	ctxA := context.WithValue(context.Background(), tenantKey, "tenant-a")
	ctxB := context.WithValue(context.Background(), tenantKey, "tenant-b")

	store.AddDocument(ctxA, "Doc A1", "content")
	store.AddDocument(ctxB, "Doc B1", "content")

	all, _ := store.GetDocumentsNoFilter(context.Background())
	if len(all) != 2 {
		t.Errorf("expected 2 docs without filter, got %d", len(all))
	}
}

func TestGetDocument_TenantScoped(t *testing.T) {
	store := NewTenantStore()
	ctxA := context.WithValue(context.Background(), tenantKey, "tenant-a")
	ctxB := context.WithValue(context.Background(), tenantKey, "tenant-b")

	docA, _ := store.AddDocument(ctxA, "Secret A", "content")

	// Tenant A can get their own doc.
	got, err := store.GetDocument(ctxA, docA.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.Title != "Secret A" {
		t.Errorf("expected 'Secret A', got %s", got.Title)
	}

	// Tenant B cannot get tenant A's doc.
	_, err = store.GetDocument(ctxB, docA.ID)
	if err == nil {
		t.Error("expected tenant B to be denied access to tenant A's document")
	}
}

func TestCreateDocument_Valid(t *testing.T) {
	store := NewTenantStore()
	ctx := context.WithValue(context.Background(), tenantKey, "test-tenant")
	doc, err := store.AddDocument(ctx, "Test Doc", "test content")
	if err != nil {
		t.Fatal(err)
	}
	if doc.ID == "" {
		t.Error("expected non-empty ID")
	}
	if doc.Title != "Test Doc" {
		t.Errorf("expected 'Test Doc', got %s", doc.Title)
	}
}
