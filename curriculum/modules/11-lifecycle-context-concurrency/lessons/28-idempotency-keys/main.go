package main

import (
	"fmt"
	"sync"
)

type IdempotencyStore struct {
	mu    sync.Mutex
	store map[string]string
}

func NewIdempotencyStore() *IdempotencyStore {
	return &IdempotencyStore{store: make(map[string]string)}
}

func (s *IdempotencyStore) Get(key string) (string, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	val, ok := s.store[key]
	return val, ok
}

func (s *IdempotencyStore) Set(key, value string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.store[key]; exists {
		return false
	}
	s.store[key] = value
	return true
}

type PaymentRequest struct {
	ID       string
	Amount   int
	Currency string
}

func processPayment(store *IdempotencyStore, req PaymentRequest) (string, bool) {
	if result, ok := store.Get(req.ID); ok {
		return result, true
	}
	result := fmt.Sprintf("charged %d %s", req.Amount, req.Currency)
	if store.Set(req.ID, result) {
		return result, false
	}
	if result, ok := store.Get(req.ID); ok {
		return result, true
	}
	return "", false
}

func main() {
	store := NewIdempotencyStore()
	req := PaymentRequest{ID: "pay_123", Amount: 5000, Currency: "USD"}

	result1, dedup1 := processPayment(store, req)
	fmt.Printf("First call: %s (dedup=%v)\n", result1, dedup1)

	result2, dedup2 := processPayment(store, req)
	fmt.Printf("Second call: %s (dedup=%v)\n", result2, dedup2)
}
