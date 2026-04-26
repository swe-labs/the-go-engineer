package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/rasel9t6/the-go-engineer/11-flagship/01-opslane/internal/db"
	"github.com/rasel9t6/the-go-engineer/11-flagship/01-opslane/internal/models"
)

func TestOrderServiceCreateOrderCreatesPendingOrderAndReservesInventory(t *testing.T) {
	t.Parallel()

	repo := newStubOrderRepository()
	inventory := &stubInventoryCoordinator{}
	service := NewOrderService(repo, inventory)

	result, err := service.CreateOrder(context.Background(), CreateOrderInput{
		TenantID:       7,
		UserID:         42,
		TotalCents:     2500,
		Currency:       "usd",
		IdempotencyKey: " checkout-123 ",
	})
	if err != nil {
		t.Fatalf("CreateOrder returned error: %v", err)
	}

	if !result.Created {
		t.Fatal("Created = false, want true")
	}
	if result.Order.Status != models.OrderStatusPending {
		t.Fatalf("status = %q, want pending", result.Order.Status)
	}
	if result.Order.Currency != "USD" {
		t.Fatalf("currency = %q, want USD", result.Order.Currency)
	}
	if len(inventory.reserveCalls) != 1 {
		t.Fatalf("reserve calls = %d, want 1", len(inventory.reserveCalls))
	}
	if repo.createCalls != 1 {
		t.Fatalf("create calls = %d, want 1", repo.createCalls)
	}
}

func TestOrderServiceCreateOrderReturnsExistingOrderForIdempotentRetry(t *testing.T) {
	t.Parallel()

	repo := newStubOrderRepository()
	repo.ordersByKey[repo.makeKey(7, "checkout-123")] = models.Order{
		ID:             101,
		TenantID:       7,
		UserID:         42,
		Status:         models.OrderStatusPending,
		TotalCents:     2500,
		Currency:       "USD",
		IdempotencyKey: "checkout-123",
		CreatedAt:      time.Now().UTC(),
		UpdatedAt:      time.Now().UTC(),
	}

	inventory := &stubInventoryCoordinator{}
	service := NewOrderService(repo, inventory)

	result, err := service.CreateOrder(context.Background(), CreateOrderInput{
		TenantID:       7,
		UserID:         42,
		TotalCents:     2500,
		Currency:       "USD",
		IdempotencyKey: "checkout-123",
	})
	if err != nil {
		t.Fatalf("CreateOrder returned error: %v", err)
	}

	if result.Created {
		t.Fatal("Created = true, want false")
	}
	if result.Order.ID != 101 {
		t.Fatalf("order id = %d, want 101", result.Order.ID)
	}
	if len(inventory.reserveCalls) != 0 {
		t.Fatalf("reserve calls = %d, want 0", len(inventory.reserveCalls))
	}
	if len(inventory.releaseCalls) != 1 {
		t.Fatalf("release calls = %d, want 1", len(inventory.releaseCalls))
	}
	if repo.createCalls != 0 {
		t.Fatalf("create calls = %d, want 0", repo.createCalls)
	}
}

func TestOrderServiceCreateOrderSurfacesReleaseFailureOnIdempotentRetry(t *testing.T) {
	t.Parallel()

	repo := newStubOrderRepository()
	repo.ordersByKey[repo.makeKey(7, "checkout-123")] = models.Order{
		ID:             101,
		TenantID:       7,
		UserID:         42,
		Status:         models.OrderStatusPending,
		TotalCents:     2500,
		Currency:       "USD",
		IdempotencyKey: "checkout-123",
		CreatedAt:      time.Now().UTC(),
		UpdatedAt:      time.Now().UTC(),
	}

	inventory := &stubInventoryCoordinator{
		releaseErr: errors.New("release failed"),
	}
	service := NewOrderService(repo, inventory)

	_, err := service.CreateOrder(context.Background(), CreateOrderInput{
		TenantID:       7,
		UserID:         42,
		TotalCents:     2500,
		Currency:       "USD",
		IdempotencyKey: "checkout-123",
	})
	if err == nil {
		t.Fatal("CreateOrder error = nil, want failure")
	}
	if !errors.Is(err, ErrInventoryUnavailable) {
		t.Fatalf("CreateOrder error = %v, want ErrInventoryUnavailable", err)
	}
	if len(inventory.reserveCalls) != 0 {
		t.Fatalf("reserve calls = %d, want 0", len(inventory.reserveCalls))
	}
	if len(inventory.releaseCalls) != 1 {
		t.Fatalf("release calls = %d, want 1", len(inventory.releaseCalls))
	}
}

func TestOrderServiceCreateOrderReleasesInventoryAfterDuplicateInsertRace(t *testing.T) {
	t.Parallel()

	repo := newStubOrderRepository()
	repo.createErr = db.ErrDuplicateValue
	repo.createHook = func(order *models.Order) {
		repo.seedOrder(models.Order{
			ID:             101,
			TenantID:       order.TenantID,
			UserID:         order.UserID,
			Status:         models.OrderStatusPending,
			TotalCents:     order.TotalCents,
			Currency:       order.Currency,
			IdempotencyKey: order.IdempotencyKey,
			CreatedAt:      time.Now().UTC(),
			UpdatedAt:      time.Now().UTC(),
		})
	}

	inventory := &stubInventoryCoordinator{}
	service := NewOrderService(repo, inventory)

	result, err := service.CreateOrder(context.Background(), CreateOrderInput{
		TenantID:       7,
		UserID:         42,
		TotalCents:     2500,
		Currency:       "USD",
		IdempotencyKey: "checkout-123",
	})
	if err != nil {
		t.Fatalf("CreateOrder returned error: %v", err)
	}

	if result.Created {
		t.Fatal("Created = true, want false")
	}
	if len(inventory.reserveCalls) != 1 {
		t.Fatalf("reserve calls = %d, want 1", len(inventory.reserveCalls))
	}
	if len(inventory.releaseCalls) != 1 {
		t.Fatalf("release calls = %d, want 1", len(inventory.releaseCalls))
	}
}

func TestOrderServiceCreateOrderSurfacesRollbackReleaseFailure(t *testing.T) {
	t.Parallel()

	createErr := errors.New("insert failed")

	repo := newStubOrderRepository()
	repo.createErr = createErr

	inventory := &stubInventoryCoordinator{
		releaseErr: errors.New("release failed"),
	}
	service := NewOrderService(repo, inventory)

	_, err := service.CreateOrder(context.Background(), CreateOrderInput{
		TenantID:       7,
		UserID:         42,
		TotalCents:     2500,
		Currency:       "USD",
		IdempotencyKey: "checkout-123",
	})
	if err == nil {
		t.Fatal("CreateOrder error = nil, want failure")
	}
	if !errors.Is(err, createErr) {
		t.Fatalf("CreateOrder error = %v, want wrapped create error", err)
	}
	if !errors.Is(err, ErrInventoryUnavailable) {
		t.Fatalf("CreateOrder error = %v, want ErrInventoryUnavailable", err)
	}
	if len(inventory.releaseCalls) != 1 {
		t.Fatalf("release calls = %d, want 1", len(inventory.releaseCalls))
	}
}

func TestOrderServiceCreateOrderRejectsInvalidInput(t *testing.T) {
	t.Parallel()

	service := NewOrderService(newStubOrderRepository(), &stubInventoryCoordinator{})

	_, err := service.CreateOrder(context.Background(), CreateOrderInput{
		TenantID:       7,
		UserID:         42,
		TotalCents:     0,
		Currency:       "US",
		IdempotencyKey: "",
	})
	if !errors.Is(err, ErrInvalidOrder) {
		t.Fatalf("CreateOrder error = %v, want ErrInvalidOrder", err)
	}
}

func TestOrderServiceTransitionOrderRejectsInvalidStatusMove(t *testing.T) {
	t.Parallel()

	repo := newStubOrderRepository()
	repo.seedOrder(models.Order{
		ID:             101,
		TenantID:       7,
		UserID:         42,
		Status:         models.OrderStatusPending,
		TotalCents:     2500,
		Currency:       "USD",
		IdempotencyKey: "checkout-123",
	})

	service := NewOrderService(repo, &stubInventoryCoordinator{})

	_, err := service.TransitionOrder(context.Background(), TransitionOrderRequest{
		TenantID: 7,
		OrderID:  101,
		Status:   models.OrderStatusPaid,
	})
	if !errors.Is(err, ErrInvalidStatusTransition) {
		t.Fatalf("TransitionOrder error = %v, want ErrInvalidStatusTransition", err)
	}
}

func TestOrderServiceTransitionOrderReleasesInventoryOnFailure(t *testing.T) {
	t.Parallel()

	repo := newStubOrderRepository()
	repo.seedOrder(models.Order{
		ID:             101,
		TenantID:       7,
		UserID:         42,
		Status:         models.OrderStatusProcessing,
		TotalCents:     2500,
		Currency:       "USD",
		IdempotencyKey: "checkout-123",
	})

	inventory := &stubInventoryCoordinator{}
	service := NewOrderService(repo, inventory)

	order, err := service.TransitionOrder(context.Background(), TransitionOrderRequest{
		TenantID: 7,
		OrderID:  101,
		Status:   models.OrderStatusFailed,
	})
	if err != nil {
		t.Fatalf("TransitionOrder returned error: %v", err)
	}

	if order.Status != models.OrderStatusFailed {
		t.Fatalf("status = %q, want failed", order.Status)
	}
	if len(inventory.releaseCalls) != 1 {
		t.Fatalf("release calls = %d, want 1", len(inventory.releaseCalls))
	}
}

func TestOrderServiceTransitionOrderReservesInventoryForRetry(t *testing.T) {
	t.Parallel()

	repo := newStubOrderRepository()
	repo.seedOrder(models.Order{
		ID:             101,
		TenantID:       7,
		UserID:         42,
		Status:         models.OrderStatusFailed,
		TotalCents:     2500,
		Currency:       "USD",
		IdempotencyKey: "checkout-123",
	})

	inventory := &stubInventoryCoordinator{}
	service := NewOrderService(repo, inventory)

	order, err := service.TransitionOrder(context.Background(), TransitionOrderRequest{
		TenantID: 7,
		OrderID:  101,
		Status:   models.OrderStatusProcessing,
	})
	if err != nil {
		t.Fatalf("TransitionOrder returned error: %v", err)
	}

	if order.Status != models.OrderStatusProcessing {
		t.Fatalf("status = %q, want processing", order.Status)
	}
	if len(inventory.reserveCalls) != 1 {
		t.Fatalf("reserve calls = %d, want 1", len(inventory.reserveCalls))
	}
}

func TestOrderServiceTransitionOrderSurfacesRollbackReleaseFailure(t *testing.T) {
	t.Parallel()

	updateErr := errors.New("update failed")

	repo := newStubOrderRepository()
	repo.updateErr = updateErr
	repo.seedOrder(models.Order{
		ID:             101,
		TenantID:       7,
		UserID:         42,
		Status:         models.OrderStatusFailed,
		TotalCents:     2500,
		Currency:       "USD",
		IdempotencyKey: "checkout-123",
	})

	inventory := &stubInventoryCoordinator{
		releaseErr: errors.New("release failed"),
	}
	service := NewOrderService(repo, inventory)

	_, err := service.TransitionOrder(context.Background(), TransitionOrderRequest{
		TenantID: 7,
		OrderID:  101,
		Status:   models.OrderStatusProcessing,
	})
	if err == nil {
		t.Fatal("TransitionOrder error = nil, want failure")
	}
	if !errors.Is(err, updateErr) {
		t.Fatalf("TransitionOrder error = %v, want wrapped update error", err)
	}
	if !errors.Is(err, ErrInventoryUnavailable) {
		t.Fatalf("TransitionOrder error = %v, want ErrInventoryUnavailable", err)
	}
	if len(inventory.reserveCalls) != 1 {
		t.Fatalf("reserve calls = %d, want 1", len(inventory.reserveCalls))
	}
	if len(inventory.releaseCalls) != 1 {
		t.Fatalf("release calls = %d, want 1", len(inventory.releaseCalls))
	}
}

func TestOrderServiceTransitionOrderRetriesReleaseForIdempotentTerminalState(t *testing.T) {
	t.Parallel()

	repo := newStubOrderRepository()
	repo.seedOrder(models.Order{
		ID:             101,
		TenantID:       7,
		UserID:         42,
		Status:         models.OrderStatusFailed,
		TotalCents:     2500,
		Currency:       "USD",
		IdempotencyKey: "checkout-123",
	})

	inventory := &stubInventoryCoordinator{}
	service := NewOrderService(repo, inventory)

	order, err := service.TransitionOrder(context.Background(), TransitionOrderRequest{
		TenantID: 7,
		OrderID:  101,
		Status:   models.OrderStatusFailed,
	})
	if err != nil {
		t.Fatalf("TransitionOrder returned error: %v", err)
	}
	if order.Status != models.OrderStatusFailed {
		t.Fatalf("status = %q, want failed", order.Status)
	}
	if len(inventory.releaseCalls) != 1 {
		t.Fatalf("release calls = %d, want 1", len(inventory.releaseCalls))
	}
}

type stubOrderRepository struct {
	ordersByID  map[int64]models.Order
	ordersByKey map[string]models.Order
	createErr   error
	createHook  func(order *models.Order)
	updateErr   error
	createCalls int
	nextOrderID int64
}

func newStubOrderRepository() *stubOrderRepository {
	return &stubOrderRepository{
		ordersByID:  make(map[int64]models.Order),
		ordersByKey: make(map[string]models.Order),
		nextOrderID: 101,
	}
}

func (r *stubOrderRepository) CreateOrder(_ context.Context, order *models.Order) error {
	r.createCalls++
	if r.createHook != nil {
		r.createHook(order)
	}
	if r.createErr != nil {
		return r.createErr
	}

	order.ID = r.nextOrderID
	r.nextOrderID++
	order.CreatedAt = time.Now().UTC()
	order.UpdatedAt = order.CreatedAt
	r.seedOrder(*order)
	return nil
}

func (r *stubOrderRepository) GetOrderByID(_ context.Context, tenantID, orderID int64) (models.Order, error) {
	order, ok := r.ordersByID[orderID]
	if !ok || order.TenantID != tenantID {
		return models.Order{}, sql.ErrNoRows
	}

	return order, nil
}

func (r *stubOrderRepository) GetOrderByIdempotencyKey(_ context.Context, tenantID int64, idempotencyKey string) (models.Order, error) {
	order, ok := r.ordersByKey[r.makeKey(tenantID, idempotencyKey)]
	if !ok {
		return models.Order{}, sql.ErrNoRows
	}

	return order, nil
}

func (r *stubOrderRepository) UpdateOrderStatus(_ context.Context, tenantID, orderID int64, status models.OrderStatus) (models.Order, error) {
	if r.updateErr != nil {
		return models.Order{}, r.updateErr
	}

	order, ok := r.ordersByID[orderID]
	if !ok || order.TenantID != tenantID {
		return models.Order{}, sql.ErrNoRows
	}

	order.Status = status
	order.UpdatedAt = time.Now().UTC()
	r.seedOrder(order)
	return order, nil
}

func (r *stubOrderRepository) seedOrder(order models.Order) {
	r.ordersByID[order.ID] = order
	r.ordersByKey[r.makeKey(order.TenantID, order.IdempotencyKey)] = order
}

func (r *stubOrderRepository) makeKey(tenantID int64, idempotencyKey string) string {
	return fmt.Sprintf("%d::%s", tenantID, idempotencyKey)
}

type stubInventoryCoordinator struct {
	reserveCalls []InventoryReservation
	releaseCalls []InventoryReservation
	reserveErr   error
	releaseErr   error
}

func (s *stubInventoryCoordinator) Reserve(_ context.Context, reservation InventoryReservation) error {
	s.reserveCalls = append(s.reserveCalls, reservation)
	return s.reserveErr
}

func (s *stubInventoryCoordinator) Release(_ context.Context, reservation InventoryReservation) error {
	s.releaseCalls = append(s.releaseCalls, reservation)
	return s.releaseErr
}
