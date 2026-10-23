package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	_ "modernc.org/sqlite"
)

type contextKey string

const tenantContextKey = contextKey("tenant_id")

var ErrMissingTenantContext = errors.New("security error: missing required tenant in context")

// ContextWithTenant injects the validated tenant ID into the request context.
func ContextWithTenant(ctx context.Context, tenantID string) context.Context {
	return context.WithValue(ctx, tenantContextKey, tenantID)
}

// TenantFromContext extracts the tenant ID from context.
func TenantFromContext(ctx context.Context) (string, error) {
	val, ok := ctx.Value(tenantContextKey).(string)
	if !ok || val == "" {
		return "", ErrMissingTenantContext
	}
	return val, nil
}

// TenantOrder models an order entity bound to a tenant.
type TenantOrder struct {
	ID        string    `json:"id"`
	TenantID  string    `json:"tenant_id"`
	Title     string    `json:"title"`
	Amount    float64   `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
}

// IsolatedStore manages database access with strict multi-tenant boundary enforcement.
type IsolatedStore struct {
	db *sql.DB
}

func NewIsolatedStore(db *sql.DB) (*IsolatedStore, error) {
	schema := `
	CREATE TABLE IF NOT EXISTS tenant_orders (
		id TEXT NOT NULL,
		tenant_id TEXT NOT NULL,
		title TEXT NOT NULL,
		amount REAL NOT NULL,
		created_at TIMESTAMP NOT NULL,
		PRIMARY KEY (tenant_id, id)
	);
	CREATE INDEX IF NOT EXISTS idx_tenant_orders ON tenant_orders(tenant_id);
	`
	if _, err := db.Exec(schema); err != nil {
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}
	return &IsolatedStore{db: db}, nil
}

// Create inserts an order guaranteed to be bound to the caller's tenant.
func (s *IsolatedStore) Create(ctx context.Context, id, title string, amount float64) (*TenantOrder, error) {
	tenantID, err := TenantFromContext(ctx)
	if err != nil {
		return nil, err
	}

	order := &TenantOrder{
		ID:        id,
		TenantID:  tenantID,
		Title:     title,
		Amount:    amount,
		CreatedAt: time.Now().UTC(),
	}

	query := `INSERT INTO tenant_orders (id, tenant_id, title, amount, created_at) VALUES (?, ?, ?, ?, ?)`
	_, err = s.db.ExecContext(ctx, query, order.ID, order.TenantID, order.Title, order.Amount, order.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to insert order: %w", err)
	}

	return order, nil
}

// GetByID queries an order enforcing the tenant boundary.
func (s *IsolatedStore) GetByID(ctx context.Context, id string) (*TenantOrder, error) {
	tenantID, err := TenantFromContext(ctx)
	if err != nil {
		return nil, err
	}

	query := `SELECT id, tenant_id, title, amount, created_at FROM tenant_orders WHERE tenant_id = ? AND id = ?`
	var o TenantOrder
	err = s.db.QueryRowContext(ctx, query, tenantID, id).Scan(
		&o.ID, &o.TenantID, &o.Title, &o.Amount, &o.CreatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil // Not found within this tenant's silo
	}
	if err != nil {
		return nil, err
	}

	return &o, nil
}

// List returns all orders belonging solely to the authenticated tenant.
func (s *IsolatedStore) List(ctx context.Context) ([]TenantOrder, error) {
	tenantID, err := TenantFromContext(ctx)
	if err != nil {
		return nil, err
	}

	query := `SELECT id, tenant_id, title, amount, created_at FROM tenant_orders WHERE tenant_id = ? ORDER BY created_at DESC`
	rows, err := s.db.QueryContext(ctx, query, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []TenantOrder
	for rows.Next() {
		var o TenantOrder
		if err := rows.Scan(&o.ID, &o.TenantID, &o.Title, &o.Amount, &o.CreatedAt); err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}

	return orders, nil
}

// Update modifies an order only if it belongs to the authenticated tenant.
func (s *IsolatedStore) Update(ctx context.Context, id, newTitle string, newAmount float64) (bool, error) {
	tenantID, err := TenantFromContext(ctx)
	if err != nil {
		return false, err
	}

	query := `UPDATE tenant_orders SET title = ?, amount = ? WHERE tenant_id = ? AND id = ?`
	res, err := s.db.ExecContext(ctx, query, newTitle, newAmount, tenantID, id)
	if err != nil {
		return false, err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return false, err
	}

	return rowsAffected > 0, nil
}

// Delete removes an order strictly within the tenant boundary.
func (s *IsolatedStore) Delete(ctx context.Context, id string) (bool, error) {
	tenantID, err := TenantFromContext(ctx)
	if err != nil {
		return false, err
	}

	query := `DELETE FROM tenant_orders WHERE tenant_id = ? AND id = ?`
	res, err := s.db.ExecContext(ctx, query, tenantID, id)
	if err != nil {
		return false, err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return false, err
	}

	return rowsAffected > 0, nil
}

func main() {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		log.Fatalf("DB open failed: %v", err)
	}
	defer db.Close()

	store, err := NewIsolatedStore(db)
	if err != nil {
		log.Fatalf("Store init failed: %v", err)
	}

	ctxAcme := ContextWithTenant(context.Background(), "tenant-acme")
	_, _ = store.Create(ctxAcme, "ord-1", "Acme Production Widget", 99.0)

	fmt.Println("Opslane Tenant Isolation subsystem verified.")
}
