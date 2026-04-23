//go:build integration

package db

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/rasel9t6/the-go-engineer/11-flagship/01-opslane/internal/config"
	"github.com/rasel9t6/the-go-engineer/11-flagship/01-opslane/internal/models"
	"github.com/testcontainers/testcontainers-go"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestStoreSupportsTenantScopedRecords(t *testing.T) {
	t.Parallel()

	defer func() {
		if r := recover(); r != nil {
			message := fmt.Sprint(r)
			if strings.Contains(strings.ToLower(message), "docker") || strings.Contains(strings.ToLower(message), "rootless") {
				t.Skipf("skip postgres integration test: %v", r)
			}
			panic(r)
		}
	}()

	ctx := context.Background()
	database := openPostgresTestDatabase(t, ctx)
	store := NewStore(database)

	tenantA := &models.Tenant{Name: "Alpha Inc", Slug: "alpha"}
	if err := store.CreateTenant(ctx, tenantA); err != nil {
		t.Fatalf("CreateTenant(alpha) failed: %v", err)
	}

	tenantB := &models.Tenant{Name: "Beta Inc", Slug: "beta"}
	if err := store.CreateTenant(ctx, tenantB); err != nil {
		t.Fatalf("CreateTenant(beta) failed: %v", err)
	}

	storedTenant, err := store.GetTenantBySlug(ctx, tenantA.Slug)
	if err != nil {
		t.Fatalf("GetTenantBySlug failed: %v", err)
	}
	if storedTenant.ID != tenantA.ID {
		t.Fatalf("tenant id = %d, want %d", storedTenant.ID, tenantA.ID)
	}

	userA := &models.User{
		TenantID:     tenantA.ID,
		Email:        "owner@shared.test",
		DisplayName:  "Owner A",
		PasswordHash: "hash-a",
		Role:         models.UserRoleAdmin,
	}
	if err := store.CreateUser(ctx, userA); err != nil {
		t.Fatalf("CreateUser(userA) failed: %v", err)
	}

	userB := &models.User{
		TenantID:     tenantB.ID,
		Email:        "owner@shared.test",
		DisplayName:  "Owner B",
		PasswordHash: "hash-b",
		Role:         models.UserRoleAdmin,
	}
	if err := store.CreateUser(ctx, userB); err != nil {
		t.Fatalf("CreateUser(userB) failed: %v", err)
	}

	storedUser, err := store.GetUserByEmail(ctx, tenantA.ID, userA.Email)
	if err != nil {
		t.Fatalf("GetUserByEmail failed: %v", err)
	}
	if storedUser.TenantID != tenantA.ID {
		t.Fatalf("stored user tenant = %d, want %d", storedUser.TenantID, tenantA.ID)
	}

	order := &models.Order{
		TenantID:       tenantA.ID,
		UserID:         userA.ID,
		Status:         models.OrderStatusPending,
		TotalCents:     159900,
		Currency:       "USD",
		IdempotencyKey: "order-alpha-1",
	}
	if err := store.CreateOrder(ctx, order); err != nil {
		t.Fatalf("CreateOrder failed: %v", err)
	}

	crossTenantOrder := &models.Order{
		TenantID:       tenantB.ID,
		UserID:         userA.ID,
		Status:         models.OrderStatusPending,
		TotalCents:     209900,
		Currency:       "USD",
		IdempotencyKey: "order-beta-invalid-user",
	}
	if err := store.CreateOrder(ctx, crossTenantOrder); err == nil {
		t.Fatal("CreateOrder allowed cross-tenant user reference")
	}

	payment := &models.Payment{
		TenantID:          tenantA.ID,
		OrderID:           order.ID,
		Status:            models.PaymentStatusPending,
		ProviderReference: "pay-alpha-1",
		AmountCents:       order.TotalCents,
	}
	if err := store.CreatePayment(ctx, payment); err != nil {
		t.Fatalf("CreatePayment failed: %v", err)
	}

	crossTenantPayment := &models.Payment{
		TenantID:          tenantB.ID,
		OrderID:           order.ID,
		Status:            models.PaymentStatusPending,
		ProviderReference: "pay-beta-invalid-order",
		AmountCents:       order.TotalCents,
	}
	if err := store.CreatePayment(ctx, crossTenantPayment); err == nil {
		t.Fatal("CreatePayment allowed cross-tenant order reference")
	}

	orders, err := store.ListOrdersByTenant(ctx, tenantA.ID)
	if err != nil {
		t.Fatalf("ListOrdersByTenant failed: %v", err)
	}
	if len(orders) != 1 {
		t.Fatalf("orders len = %d, want 1", len(orders))
	}

	payments, err := store.ListPaymentsByOrder(ctx, tenantA.ID, order.ID)
	if err != nil {
		t.Fatalf("ListPaymentsByOrder failed: %v", err)
	}
	if len(payments) != 1 {
		t.Fatalf("payments len = %d, want 1", len(payments))
	}
}

func openPostgresTestDatabase(t *testing.T, ctx context.Context) *sql.DB {
	t.Helper()

	container, err := tcpostgres.Run(ctx,
		"postgres:16-alpine",
		tcpostgres.WithDatabase("opslane_test"),
		tcpostgres.WithUsername("opslane"),
		tcpostgres.WithPassword("secretpassword"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").WithOccurrence(2).WithStartupTimeout(10*time.Second),
		),
	)
	if err != nil {
		t.Skipf("skip postgres integration test: %v", err)
	}

	t.Cleanup(func() {
		_ = container.Terminate(context.Background())
	})

	dsn, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatalf("connection string: %v", err)
	}

	cfg := config.DatabaseConfig{
		DSN:             dsn,
		MaxOpenConns:    4,
		MaxIdleConns:    2,
		ConnMaxIdleTime: time.Minute,
		ConnMaxLifetime: 5 * time.Minute,
	}

	database, err := Open(ctx, cfg)
	if err != nil {
		t.Fatalf("open postgres db: %v", err)
	}
	t.Cleanup(func() {
		_ = database.Close()
	})

	if err := Migrate(ctx, database); err != nil {
		t.Fatalf("migrate postgres db: %v", err)
	}

	return database
}
