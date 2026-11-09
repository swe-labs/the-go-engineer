package main

import (
	"errors"
	"fmt"
	"sync"
)

type Product struct {
	ID    string
	Name  string
	Price float64
}

type ProductRepository interface {
	Save(product Product) error
	FindByID(id string) (Product, error)
	FindAll() ([]Product, error)
	Delete(id string) error
}

type InMemoryProductRepo struct {
	mu    sync.RWMutex
	store map[string]Product
}

func NewInMemoryProductRepo() *InMemoryProductRepo {
	return &InMemoryProductRepo{store: make(map[string]Product)}
}

func (r *InMemoryProductRepo) Save(product Product) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.store[product.ID] = product
	return nil
}

func (r *InMemoryProductRepo) FindByID(id string) (Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.store[id]
	if !ok {
		return Product{}, errors.New("product not found")
	}
	return p, nil
}

func (r *InMemoryProductRepo) FindAll() ([]Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]Product, 0, len(r.store))
	for _, p := range r.store {
		result = append(result, p)
	}
	return result, nil
}

func (r *InMemoryProductRepo) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.store, id)
	return nil
}

type ProductService struct {
	repo ProductRepository
}

func NewProductService(repo ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) AddProduct(id, name string, price float64) error {
	if price <= 0 {
		return errors.New("price must be positive")
	}
	return s.repo.Save(Product{ID: id, Name: name, Price: price})
}

func (s *ProductService) ListProducts() ([]Product, error) {
	return s.repo.FindAll()
}

func main() {
	repo := NewInMemoryProductRepo()
	svc := NewProductService(repo)

	_ = svc.AddProduct("p1", "Widget", 9.99)
	_ = svc.AddProduct("p2", "Gadget", 24.99)

	products, _ := svc.ListProducts()
	for _, p := range products {
		fmt.Printf("%s: %s ($%.2f)\n", p.ID, p.Name, p.Price)
	}
}
