package main

import (
	"testing"
)

func TestInMemoryProductRepo_SaveAndFind(t *testing.T) {
	repo := NewInMemoryProductRepo()
	p := Product{ID: "p1", Name: "Widget", Price: 9.99}
	err := repo.Save(p)
	if err != nil {
		t.Fatalf("Save failed: %v", err)
	}
	found, err := repo.FindByID("p1")
	if err != nil {
		t.Fatalf("FindByID failed: %v", err)
	}
	if found.Name != "Widget" {
		t.Errorf("expected Widget, got %s", found.Name)
	}
}

func TestInMemoryProductRepo_FindByID_NotFound(t *testing.T) {
	repo := NewInMemoryProductRepo()
	_, err := repo.FindByID("nonexistent")
	if err == nil {
		t.Fatal("expected error for nonexistent product")
	}
}

func TestInMemoryProductRepo_FindAll(t *testing.T) {
	repo := NewInMemoryProductRepo()
	repo.Save(Product{ID: "a", Name: "A", Price: 1})
	repo.Save(Product{ID: "b", Name: "B", Price: 2})
	products, err := repo.FindAll()
	if err != nil {
		t.Fatalf("FindAll failed: %v", err)
	}
	if len(products) != 2 {
		t.Errorf("expected 2 products, got %d", len(products))
	}
}

func TestInMemoryProductRepo_Delete(t *testing.T) {
	repo := NewInMemoryProductRepo()
	repo.Save(Product{ID: "p1", Name: "Widget", Price: 9.99})
	err := repo.Delete("p1")
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
	_, err = repo.FindByID("p1")
	if err == nil {
		t.Fatal("expected error after delete")
	}
}

func TestProductService_AddProduct_InvalidPrice(t *testing.T) {
	svc := NewProductService(NewInMemoryProductRepo())
	err := svc.AddProduct("p1", "Widget", 0)
	if err == nil {
		t.Fatal("expected error for zero price")
	}
}

func TestProductService_ListProducts(t *testing.T) {
	svc := NewProductService(NewInMemoryProductRepo())
	svc.AddProduct("p1", "Widget", 9.99)
	svc.AddProduct("p2", "Gadget", 24.99)
	products, err := svc.ListProducts()
	if err != nil {
		t.Fatalf("ListProducts failed: %v", err)
	}
	if len(products) != 2 {
		t.Errorf("expected 2 products, got %d", len(products))
	}
}
