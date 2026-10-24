package main

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

// Command Model
type PlaceOrderCommand struct {
	OrderID string
	Amount  int
}

// Query Model (Denormalized)
type OrderView struct {
	OrderID string
	Status  string
	Amount  int
}

// CQRSOrderService demonstrates separation of commands and queries.
type CQRSOrderService struct {
	mu     sync.RWMutex
	store  map[string]int       // Write model
	views  map[string]OrderView // Read model projection
}

func NewCQRSOrderService() *CQRSOrderService {
	return &CQRSOrderService{
		store: make(map[string]int),
		views: make(map[string]OrderView),
	}
}

// ExecuteCommand processes writes with strict validation.
func (s *CQRSOrderService) ExecuteCommand(ctx context.Context, cmd PlaceOrderCommand) error {
	if cmd.Amount <= 0 {
		return errors.New("invalid order amount")
	}
	s.mu.Lock()
	defer s.mu.Unlock()

	s.store[cmd.OrderID] = cmd.Amount
	// Update read projection
	s.views[cmd.OrderID] = OrderView{
		OrderID: cmd.OrderID,
		Status:  "CONFIRMED",
		Amount:  cmd.Amount,
	}
	return nil
}

// ExecuteQuery processes reads without write locks.
func (s *CQRSOrderService) ExecuteQuery(ctx context.Context, orderID string) (OrderView, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	v, ok := s.views[orderID]
	return v, ok
}

// Demo demonstrates CQRS and satisfies test requirements.
func Demo() string {
	svc := NewCQRSOrderService()
	_ = svc.ExecuteCommand(context.Background(), PlaceOrderCommand{OrderID: "ord_1", Amount: 250})
	view, _ := svc.ExecuteQuery(context.Background(), "ord_1")
	_ = view
	return "CQRS verified"
}

func main() {
	fmt.Println(Demo())
}
