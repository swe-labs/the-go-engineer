// Copyright (c) 2026 Rasel Hossen
// See LICENSE for usage terms.

package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/lib/pq"
	"github.com/swe-labs/the-go-engineer/11-flagship/01-opslane/internal/config"
	"github.com/swe-labs/the-go-engineer/11-flagship/01-opslane/internal/models"
)

// ErrInvalidReference (Error): means a tenant-scoped row points at a parent row that does not exist
var ErrInvalidReference = errors.New("invalid tenant-scoped reference")

// ErrDuplicateValue (Error): indicates that an insert or update violated a unique constraint
var ErrDuplicateValue = errors.New("duplicate value")

const postgresForeignKeyViolation = "23503"
const postgresUniqueViolation = "23505"

// TenantRepository (Interface): defines the data access contract for tenant records
type TenantRepository interface {
	CreateTenant(ctx context.Context, tenant *models.Tenant) error
	GetTenantBySlug(ctx context.Context, slug string) (models.Tenant, error)
}

// UserRepository (Interface): manages user identity records scoped to specific tenants
type UserRepository interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetUserByEmail(ctx context.Context, tenantID int64, email string) (models.User, error)
}

// OrderRepository (Interface): handles the lifecycle and retrieval of customer orders; strictly tenant-scoped
type OrderRepository interface {
	CreateOrder(ctx context.Context, order *models.Order) error
	GetOrderByID(ctx context.Context, tenantID, orderID int64) (models.Order, error)
	GetOrderByIdempotencyKey(ctx context.Context, tenantID int64, idempotencyKey string) (models.Order, error)
	UpdateOrderStatus(ctx context.Context, tenantID, orderID int64, status models.OrderStatus) (models.Order, error)
	ListOrdersByTenant(ctx context.Context, tenantID int64) ([]models.Order, error)
}

// PaymentRepository (Interface): manages financial transactions tied to specific orders
type PaymentRepository interface {
	CreatePayment(ctx context.Context, payment *models.Payment) error
	GetPaymentByProviderReference(ctx context.Context, tenantID int64, providerReference string) (models.Payment, error)
	UpdatePaymentStatus(ctx context.Context, tenantID int64, providerReference string, status models.PaymentStatus, failureReason string) (models.Payment, error)
	ListPaymentsByOrder(ctx context.Context, tenantID, orderID int64) ([]models.Payment, error)
}

// queryRunner (Interface): abstracts sql.DB and sql.Tx for transactional queries
type queryRunner interface {
	ExecContext(context.Context, string, ...any) (sql.Result, error)
	QueryContext(context.Context, string, ...any) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...any) *sql.Row
}

// Store (Struct): concrete PostgreSQL-backed implementation of all database repositories with transaction support
type Store struct {
	db *sql.DB
	q  queryRunner
}

// Open (Function): establishes and verifies a PostgreSQL connection pool from config
func Open(ctx context.Context, cfg config.DatabaseConfig) (*sql.DB, error) {
	database, err := sql.Open("postgres", cfg.DSN)
	if err != nil {
		return nil, fmt.Errorf("open postgres database: %w", err)
	}

	database.SetMaxOpenConns(cfg.MaxOpenConns)
	database.SetMaxIdleConns(cfg.MaxIdleConns)
	database.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)
	database.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	if err := database.PingContext(ctx); err != nil {
		database.Close()
		return nil, fmt.Errorf("ping postgres database: %w", err)
	}

	return database, nil
}

// NewStore (Constructor): creates a new PostgreSQL-backed Store from a verified connection
func NewStore(database *sql.DB) *Store {
	return &Store{
		db: database,
		q:  database,
	}
}

// Ping (Method): verifies the database connection is alive and healthy
func (s *Store) Ping(ctx context.Context) error {
	return s.db.PingContext(ctx)
}

// WithTx (Method): executes a function within a database transaction with automatic rollback on error or panic
func (s *Store) WithTx(ctx context.Context, fn func(*Store) error) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	committed := false
	defer func() {
		if !committed {
			_ = tx.Rollback()
		}
	}()

	defer func() {
		if recovered := recover(); recovered != nil {
			panic(recovered)
		}
	}()

	txStore := &Store{
		db: s.db,
		q:  tx,
	}

	if err := fn(txStore); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}
	committed = true

	return nil
}

// CreateTenant (Method): inserts a new tenant, defaulting CreatedAt to UTC; returns ErrDuplicateValue on slug conflict
func (s *Store) CreateTenant(ctx context.Context, tenant *models.Tenant) error {
	if tenant.CreatedAt.IsZero() {
		tenant.CreatedAt = time.Now().UTC()
	}

	err := s.q.QueryRowContext(
		ctx,
		`INSERT INTO tenants (name, slug, created_at) VALUES ($1, $2, $3) RETURNING id`,
		tenant.Name,
		tenant.Slug,
		tenant.CreatedAt,
	).Scan(&tenant.ID)
	if err != nil {
		if isUniqueViolation(err) {
			return fmt.Errorf("insert tenant: %w", ErrDuplicateValue)
		}
		return fmt.Errorf("insert tenant: %w", err)
	}

	return nil
}

// GetTenantBySlug (Method): retrieves a tenant by its unique URL-friendly slug
func (s *Store) GetTenantBySlug(ctx context.Context, slug string) (models.Tenant, error) {
	var tenant models.Tenant

	err := s.q.QueryRowContext(
		ctx,
		`SELECT id, name, slug, created_at FROM tenants WHERE slug = $1`,
		slug,
	).Scan(&tenant.ID, &tenant.Name, &tenant.Slug, &tenant.CreatedAt)
	if err != nil {
		return models.Tenant{}, fmt.Errorf("get tenant by slug: %w", err)
	}

	return tenant, nil
}

// CreateUser (Method): persists a new user record; returns ErrInvalidReference or ErrDuplicateValue on constraint violations
func (s *Store) CreateUser(ctx context.Context, user *models.User) error {
	if user.CreatedAt.IsZero() {
		user.CreatedAt = time.Now().UTC()
	}

	err := s.q.QueryRowContext(
		ctx,
		`INSERT INTO users (tenant_id, email, display_name, password_hash, role, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6)
		 RETURNING id`,
		user.TenantID,
		user.Email,
		user.DisplayName,
		user.PasswordHash,
		user.Role,
		user.CreatedAt,
	).Scan(&user.ID)
	if err != nil {
		if isForeignKeyViolation(err) {
			return fmt.Errorf("insert user: %w", ErrInvalidReference)
		}
		if isUniqueViolation(err) {
			return fmt.Errorf("insert user: %w", ErrDuplicateValue)
		}
		return fmt.Errorf("insert user: %w", err)
	}

	return nil
}

// GetUserByEmail (Method): finds a user by email within a specific tenant; email is unique per tenant
func (s *Store) GetUserByEmail(ctx context.Context, tenantID int64, email string) (models.User, error) {
	var user models.User

	err := s.q.QueryRowContext(
		ctx,
		`SELECT id, tenant_id, email, display_name, password_hash, role, created_at
		 FROM users
		 WHERE tenant_id = $1 AND email = $2`,
		tenantID,
		email,
	).Scan(
		&user.ID,
		&user.TenantID,
		&user.Email,
		&user.DisplayName,
		&user.PasswordHash,
		&user.Role,
		&user.CreatedAt,
	)
	if err != nil {
		return models.User{}, fmt.Errorf("get user by email: %w", err)
	}

	return user, nil
}

// CreateOrder (Method): initiates a new customer order with idempotency key uniqueness for retry safety
func (s *Store) CreateOrder(ctx context.Context, order *models.Order) error {
	now := time.Now().UTC()
	if order.CreatedAt.IsZero() {
		order.CreatedAt = now
	}
	if order.UpdatedAt.IsZero() {
		order.UpdatedAt = now
	}

	err := s.q.QueryRowContext(
		ctx,
		`INSERT INTO orders (tenant_id, user_id, status, total_cents, currency, idempotency_key, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		 RETURNING id`,
		order.TenantID,
		order.UserID,
		order.Status,
		order.TotalCents,
		order.Currency,
		order.IdempotencyKey,
		order.CreatedAt,
		order.UpdatedAt,
	).Scan(&order.ID)
	if err != nil {
		if isUniqueViolation(err) {
			return fmt.Errorf("insert order: %w", ErrDuplicateValue)
		}
		return fmt.Errorf("insert order: %w", err)
	}

	return nil
}

// ListOrdersByTenant (Method): retrieves all orders for a tenant, sorted newest first
func (s *Store) ListOrdersByTenant(ctx context.Context, tenantID int64) ([]models.Order, error) {
	rows, err := s.q.QueryContext(
		ctx,
		`SELECT id, tenant_id, user_id, status, total_cents, currency, idempotency_key, created_at, updated_at
		 FROM orders
		 WHERE tenant_id = $1
		 ORDER BY created_at DESC`,
		tenantID,
	)
	if err != nil {
		return nil, fmt.Errorf("list orders by tenant: %w", err)
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var order models.Order
		if err := rows.Scan(
			&order.ID,
			&order.TenantID,
			&order.UserID,
			&order.Status,
			&order.TotalCents,
			&order.Currency,
			&order.IdempotencyKey,
			&order.CreatedAt,
			&order.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan order: %w", err)
		}
		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate orders: %w", err)
	}

	return orders, nil
}

// GetOrderByID (Method): fetches a single order with tenant isolation enforced at the query level
func (s *Store) GetOrderByID(ctx context.Context, tenantID, orderID int64) (models.Order, error) {
	var order models.Order

	err := s.q.QueryRowContext(
		ctx,
		`SELECT id, tenant_id, user_id, status, total_cents, currency, idempotency_key, created_at, updated_at
		 FROM orders
		 WHERE tenant_id = $1 AND id = $2`,
		tenantID,
		orderID,
	).Scan(
		&order.ID,
		&order.TenantID,
		&order.UserID,
		&order.Status,
		&order.TotalCents,
		&order.Currency,
		&order.IdempotencyKey,
		&order.CreatedAt,
		&order.UpdatedAt,
	)
	if err != nil {
		return models.Order{}, fmt.Errorf("get order by id: %w", err)
	}

	return order, nil
}

// GetOrderByIdempotencyKey (Method): fetches an existing order by idempotency key for safe request replay
func (s *Store) GetOrderByIdempotencyKey(ctx context.Context, tenantID int64, idempotencyKey string) (models.Order, error) {
	var order models.Order

	err := s.q.QueryRowContext(
		ctx,
		`SELECT id, tenant_id, user_id, status, total_cents, currency, idempotency_key, created_at, updated_at
		 FROM orders
		 WHERE tenant_id = $1 AND idempotency_key = $2`,
		tenantID,
		idempotencyKey,
	).Scan(
		&order.ID,
		&order.TenantID,
		&order.UserID,
		&order.Status,
		&order.TotalCents,
		&order.Currency,
		&order.IdempotencyKey,
		&order.CreatedAt,
		&order.UpdatedAt,
	)
	if err != nil {
		return models.Order{}, fmt.Errorf("get order by idempotency key: %w", err)
	}

	return order, nil
}

// UpdateOrderStatus (Method): transitions an order to a new status and returns the updated record
func (s *Store) UpdateOrderStatus(ctx context.Context, tenantID, orderID int64, status models.OrderStatus) (models.Order, error) {
	var order models.Order
	updatedAt := time.Now().UTC()

	err := s.q.QueryRowContext(
		ctx,
		`UPDATE orders
		 SET status = $3, updated_at = $4
		 WHERE tenant_id = $1 AND id = $2
		 RETURNING id, tenant_id, user_id, status, total_cents, currency, idempotency_key, created_at, updated_at`,
		tenantID,
		orderID,
		status,
		updatedAt,
	).Scan(
		&order.ID,
		&order.TenantID,
		&order.UserID,
		&order.Status,
		&order.TotalCents,
		&order.Currency,
		&order.IdempotencyKey,
		&order.CreatedAt,
		&order.UpdatedAt,
	)
	if err != nil {
		return models.Order{}, fmt.Errorf("update order status: %w", err)
	}

	return order, nil
}

// CreatePayment (Method): records a new payment attempt with FK and uniqueness constraint checks
func (s *Store) CreatePayment(ctx context.Context, payment *models.Payment) error {
	now := time.Now().UTC()
	if payment.CreatedAt.IsZero() {
		payment.CreatedAt = now
	}
	if payment.UpdatedAt.IsZero() {
		payment.UpdatedAt = now
	}

	err := s.q.QueryRowContext(
		ctx,
		`INSERT INTO payments (tenant_id, order_id, status, provider_reference, amount_cents, failure_reason, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		 RETURNING id`,
		payment.TenantID,
		payment.OrderID,
		payment.Status,
		payment.ProviderReference,
		payment.AmountCents,
		payment.FailureReason,
		payment.CreatedAt,
		payment.UpdatedAt,
	).Scan(&payment.ID)
	if err != nil {
		if isForeignKeyViolation(err) {
			return fmt.Errorf("insert payment: %w", ErrInvalidReference)
		}
		if isUniqueViolation(err) {
			return fmt.Errorf("insert payment: %w", ErrDuplicateValue)
		}
		return fmt.Errorf("insert payment: %w", err)
	}

	return nil
}

// GetPaymentByProviderReference (Method): finds a payment by the gateway's external reference ID
func (s *Store) GetPaymentByProviderReference(ctx context.Context, tenantID int64, providerReference string) (models.Payment, error) {
	var payment models.Payment

	err := s.q.QueryRowContext(
		ctx,
		`SELECT id, tenant_id, order_id, status, provider_reference, amount_cents, failure_reason, created_at, updated_at
		 FROM payments
		 WHERE tenant_id = $1 AND provider_reference = $2`,
		tenantID,
		providerReference,
	).Scan(
		&payment.ID,
		&payment.TenantID,
		&payment.OrderID,
		&payment.Status,
		&payment.ProviderReference,
		&payment.AmountCents,
		&payment.FailureReason,
		&payment.CreatedAt,
		&payment.UpdatedAt,
	)
	if err != nil {
		return models.Payment{}, fmt.Errorf("get payment by provider reference: %w", err)
	}

	return payment, nil
}

// UpdatePaymentStatus (Method): records the outcome of a payment attempt with optional failure reason
func (s *Store) UpdatePaymentStatus(ctx context.Context, tenantID int64, providerReference string, status models.PaymentStatus, failureReason string) (models.Payment, error) {
	var payment models.Payment
	updatedAt := time.Now().UTC()

	err := s.q.QueryRowContext(
		ctx,
		`UPDATE payments
		 SET status = $3, failure_reason = $4, updated_at = $5
		 WHERE tenant_id = $1 AND provider_reference = $2
		 RETURNING id, tenant_id, order_id, status, provider_reference, amount_cents, failure_reason, created_at, updated_at`,
		tenantID,
		providerReference,
		status,
		failureReason,
		updatedAt,
	).Scan(
		&payment.ID,
		&payment.TenantID,
		&payment.OrderID,
		&payment.Status,
		&payment.ProviderReference,
		&payment.AmountCents,
		&payment.FailureReason,
		&payment.CreatedAt,
		&payment.UpdatedAt,
	)
	if err != nil {
		return models.Payment{}, fmt.Errorf("update payment status: %w", err)
	}

	return payment, nil
}

// ListPaymentsByOrder (Method): retrieves all payment attempts for an order, oldest first
func (s *Store) ListPaymentsByOrder(ctx context.Context, tenantID, orderID int64) ([]models.Payment, error) {
	rows, err := s.q.QueryContext(
		ctx,
		`SELECT id, tenant_id, order_id, status, provider_reference, amount_cents, failure_reason, created_at, updated_at
		 FROM payments
		 WHERE tenant_id = $1 AND order_id = $2
		 ORDER BY created_at ASC`,
		tenantID,
		orderID,
	)
	if err != nil {
		return nil, fmt.Errorf("list payments by order: %w", err)
	}
	defer rows.Close()

	var payments []models.Payment
	for rows.Next() {
		var payment models.Payment
		if err := rows.Scan(
			&payment.ID,
			&payment.TenantID,
			&payment.OrderID,
			&payment.Status,
			&payment.ProviderReference,
			&payment.AmountCents,
			&payment.FailureReason,
			&payment.CreatedAt,
			&payment.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan payment: %w", err)
		}
		payments = append(payments, payment)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate payments: %w", err)
	}

	return payments, nil
}

// isForeignKeyViolation (Function): checks if a Postgres error is a foreign key violation
func isForeignKeyViolation(err error) bool {
	var pqErr *pq.Error
	return errors.As(err, &pqErr) && string(pqErr.Code) == postgresForeignKeyViolation
}

// isUniqueViolation (Function): checks if a Postgres error is a unique constraint violation
func isUniqueViolation(err error) bool {
	var pqErr *pq.Error
	return errors.As(err, &pqErr) && string(pqErr.Code) == postgresUniqueViolation
}
