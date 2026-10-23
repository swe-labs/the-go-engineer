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

// Domain Models

type Tenant struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Tier      string    `json:"tier"` // "starter", "pro", "enterprise"
	CreatedAt time.Time `json:"created_at"`
}

type User struct {
	ID           string    `json:"id"`
	TenantID     string    `json:"tenant_id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Role         string    `json:"role"` // "owner", "admin", "member", "viewer"
	CreatedAt    time.Time `json:"created_at"`
}

type Order struct {
	ID          string    `json:"id"`
	TenantID    string    `json:"tenant_id"`
	UserID      string    `json:"user_id"`
	TotalAmount float64   `json:"total_amount"`
	Currency    string    `json:"currency"`
	Status      string    `json:"status"` // "PENDING", "PROCESSING", "COMPLETED", "FAILED", "CANCELLED"
	Notes       string    `json:"notes,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Payment struct {
	ID             string    `json:"id"`
	TenantID       string    `json:"tenant_id"`
	OrderID        string    `json:"order_id"`
	IdempotencyKey string    `json:"idempotency_key"`
	Amount         float64   `json:"amount"`
	Status         string    `json:"status"` // "PENDING", "SUCCESS", "FAILED"
	CreatedAt      time.Time `json:"created_at"`
}

type AuditLog struct {
	ID        string    `json:"id"`
	TenantID  string    `json:"tenant_id"`
	ActorID   string    `json:"actor_id"`
	Action    string    `json:"action"`
	EntityID  string    `json:"entity_id"`
	Details   string    `json:"details"`
	CreatedAt time.Time `json:"created_at"`
}

// Migration definition
type Migration struct {
	Version int
	Name    string
	Up      string
	Down    string
}

// SchemaMigrations defines the foundational database schema history
var SchemaMigrations = []Migration{
	{
		Version: 1,
		Name:    "create_tenants_and_users",
		Up: `
		CREATE TABLE IF NOT EXISTS tenants (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			tier TEXT NOT NULL DEFAULT 'starter',
			created_at TIMESTAMP NOT NULL
		);

		CREATE TABLE IF NOT EXISTS users (
			id TEXT PRIMARY KEY,
			tenant_id TEXT NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
			email TEXT NOT NULL UNIQUE,
			password_hash TEXT NOT NULL,
			role TEXT NOT NULL DEFAULT 'member',
			created_at TIMESTAMP NOT NULL
		);
		CREATE INDEX IF NOT EXISTS idx_users_tenant ON users(tenant_id);
		`,
		Down: `
		DROP TABLE IF EXISTS users;
		DROP TABLE IF EXISTS tenants;
		`,
	},
	{
		Version: 2,
		Name:    "create_orders_and_payments",
		Up: `
		CREATE TABLE IF NOT EXISTS orders (
			id TEXT PRIMARY KEY,
			tenant_id TEXT NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
			user_id TEXT NOT NULL REFERENCES users(id),
			total_amount REAL NOT NULL,
			currency TEXT NOT NULL DEFAULT 'USD',
			status TEXT NOT NULL DEFAULT 'PENDING',
			notes TEXT,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL
		);
		CREATE INDEX IF NOT EXISTS idx_orders_tenant_status ON orders(tenant_id, status);

		CREATE TABLE IF NOT EXISTS payments (
			id TEXT PRIMARY KEY,
			tenant_id TEXT NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
			order_id TEXT NOT NULL REFERENCES orders(id),
			idempotency_key TEXT NOT NULL UNIQUE,
			amount REAL NOT NULL,
			status TEXT NOT NULL,
			created_at TIMESTAMP NOT NULL
		);
		CREATE INDEX IF NOT EXISTS idx_payments_idempotency ON payments(idempotency_key);
		`,
		Down: `
		DROP TABLE IF EXISTS payments;
		DROP TABLE IF EXISTS orders;
		`,
	},
	{
		Version: 3,
		Name:    "create_audit_logs",
		Up: `
		CREATE TABLE IF NOT EXISTS audit_logs (
			id TEXT PRIMARY KEY,
			tenant_id TEXT NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
			actor_id TEXT NOT NULL,
			action TEXT NOT NULL,
			entity_id TEXT NOT NULL,
			details TEXT NOT NULL,
			created_at TIMESTAMP NOT NULL
		);
		CREATE INDEX IF NOT EXISTS idx_audit_tenant_time ON audit_logs(tenant_id, created_at);
		`,
		Down: `
		DROP TABLE IF EXISTS audit_logs;
		`,
	},
}

// Database Connection & Pooling

func ConnectDatabase(dsn string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Production connection pool tuning
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(15 * time.Minute)
	db.SetConnMaxIdleTime(5 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("database ping failed: %w", err)
	}

	return db, nil
}

// RunMigrations executes pending schema migrations sequentially in a transaction.
func RunMigrations(db *sql.DB, migrations []Migration) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version INTEGER PRIMARY KEY,
			name TEXT NOT NULL,
			applied_at TIMESTAMP NOT NULL
		);
	`)
	if err != nil {
		return fmt.Errorf("failed to ensure schema_migrations table: %w", err)
	}

	for _, m := range migrations {
		var exists int
		err := db.QueryRow("SELECT COUNT(1) FROM schema_migrations WHERE version = ?", m.Version).Scan(&exists)
		if err != nil {
			return fmt.Errorf("failed to check migration version %d: %w", m.Version, err)
		}
		if exists > 0 {
			continue // Already applied
		}

		tx, err := db.Begin()
		if err != nil {
			return fmt.Errorf("failed to begin migration transaction %d: %w", m.Version, err)
		}

		if _, err := tx.Exec(m.Up); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("migration %d (%s) failed: %w", m.Version, m.Name, err)
		}

		now := time.Now().UTC()
		if _, err := tx.Exec("INSERT INTO schema_migrations (version, name, applied_at) VALUES (?, ?, ?)", m.Version, m.Name, now); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("failed to record migration %d: %w", m.Version, err)
		}

		if err := tx.Commit(); err != nil {
			return fmt.Errorf("failed to commit migration %d: %w", m.Version, err)
		}
	}
	return nil
}

// OrderRepository provides data access operations for Orders.
type OrderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) Create(ctx context.Context, o *Order) error {
	query := `
		INSERT INTO orders (id, tenant_id, user_id, total_amount, currency, status, notes, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err := r.db.ExecContext(ctx, query,
		o.ID, o.TenantID, o.UserID, o.TotalAmount, o.Currency, o.Status, o.Notes, o.CreatedAt, o.UpdatedAt,
	)
	return err
}

func (r *OrderRepository) GetByID(ctx context.Context, tenantID, orderID string) (*Order, error) {
	query := `
		SELECT id, tenant_id, user_id, total_amount, currency, status, notes, created_at, updated_at
		FROM orders
		WHERE tenant_id = ? AND id = ?
	`
	var o Order
	var notes sql.NullString
	err := r.db.QueryRowContext(ctx, query, tenantID, orderID).Scan(
		&o.ID, &o.TenantID, &o.UserID, &o.TotalAmount, &o.Currency, &o.Status, &notes, &o.CreatedAt, &o.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if notes.Valid {
		o.Notes = notes.String
	}
	return &o, nil
}

func main() {
	db, err := ConnectDatabase(":memory:")
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	defer db.Close()

	if err := RunMigrations(db, SchemaMigrations); err != nil {
		log.Fatalf("Migrations failed: %v", err)
	}

	fmt.Println("Opslane database and models initialized successfully with 3 schema migrations.")
}
