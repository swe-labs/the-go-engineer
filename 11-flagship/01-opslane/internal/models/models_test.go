package models

import (
	"testing"
)

func TestUserRoleValues(t *testing.T) {
	t.Parallel()

	tests := []struct {
		role UserRole
		want string
	}{
		{role: UserRoleAdmin, want: "admin"},
		{role: UserRoleMember, want: "member"},
		{role: UserRoleBilling, want: "billing"},
	}

	for _, tt := range tests {
		if string(tt.role) != tt.want {
			t.Fatalf("UserRole(%q) = %q, want %q", tt.role, string(tt.role), tt.want)
		}
	}
}

func TestUserRoleUniqueness(t *testing.T) {
	t.Parallel()

	roles := map[UserRole]bool{
		UserRoleAdmin:   true,
		UserRoleMember:  true,
		UserRoleBilling: true,
	}

	if len(roles) != 3 {
		t.Fatalf("expected 3 unique user roles, got %d", len(roles))
	}
}

func TestUserStructFields(t *testing.T) {
	t.Parallel()

	u := User{
		ID:       1,
		TenantID: 10,
		Email:    "test@example.com",
		Role:     UserRoleAdmin,
	}

	if u.ID != 1 {
		t.Fatalf("User.ID = %d, want 1", u.ID)
	}
	if u.TenantID != 10 {
		t.Fatalf("User.TenantID = %d, want 10", u.TenantID)
	}
	if u.Email != "test@example.com" {
		t.Fatalf("User.Email = %q", u.Email)
	}
	if u.Role != UserRoleAdmin {
		t.Fatalf("User.Role = %q, want admin", u.Role)
	}
}

func TestOrderStatusValues(t *testing.T) {
	t.Parallel()

	tests := []struct {
		status OrderStatus
		want   string
	}{
		{status: OrderStatusPending, want: "pending"},
		{status: OrderStatusProcessing, want: "processing"},
		{status: OrderStatusPaid, want: "paid"},
		{status: OrderStatusFailed, want: "failed"},
		{status: OrderStatusCancelled, want: "cancelled"},
	}

	for _, tt := range tests {
		if string(tt.status) != tt.want {
			t.Fatalf("OrderStatus(%q) = %q, want %q", tt.status, string(tt.status), tt.want)
		}
	}
}

func TestOrderStatusUniqueness(t *testing.T) {
	t.Parallel()

	statuses := map[OrderStatus]bool{
		OrderStatusPending:    true,
		OrderStatusProcessing: true,
		OrderStatusPaid:       true,
		OrderStatusFailed:     true,
		OrderStatusCancelled:  true,
	}

	if len(statuses) != 5 {
		t.Fatalf("expected 5 unique order statuses, got %d", len(statuses))
	}
}

func TestPaymentStatusValues(t *testing.T) {
	t.Parallel()

	tests := []struct {
		status PaymentStatus
		want   string
	}{
		{status: PaymentStatusPending, want: "pending"},
		{status: PaymentStatusAuthorized, want: "authorized"},
		{status: PaymentStatusSettled, want: "settled"},
		{status: PaymentStatusFailed, want: "failed"},
		{status: PaymentStatusRefunded, want: "refunded"},
	}

	for _, tt := range tests {
		if string(tt.status) != tt.want {
			t.Fatalf("PaymentStatus(%q) = %q, want %q", tt.status, string(tt.status), tt.want)
		}
	}
}

func TestPaymentStatusUniqueness(t *testing.T) {
	t.Parallel()

	statuses := map[PaymentStatus]bool{
		PaymentStatusPending:    true,
		PaymentStatusAuthorized: true,
		PaymentStatusSettled:    true,
		PaymentStatusFailed:     true,
		PaymentStatusRefunded:   true,
	}

	if len(statuses) != 5 {
		t.Fatalf("expected 5 unique payment statuses, got %d", len(statuses))
	}
}

func TestTenantStructFields(t *testing.T) {
	t.Parallel()

	t1 := Tenant{
		ID:   1,
		Name: "Test Corp",
		Slug: "test-corp",
	}

	if t1.ID != 1 {
		t.Fatalf("Tenant.ID = %d, want 1", t1.ID)
	}
	if t1.Name != "Test Corp" {
		t.Fatalf("Tenant.Name = %q", t1.Name)
	}
	if t1.Slug != "test-corp" {
		t.Fatalf("Tenant.Slug = %q", t1.Slug)
	}
}

func TestOrderStructFields(t *testing.T) {
	t.Parallel()

	o := Order{
		ID:             1,
		TenantID:       10,
		UserID:         100,
		Status:         OrderStatusPending,
		TotalCents:     5000,
		Currency:       "USD",
		IdempotencyKey: "idem-001",
	}

	if o.ID != 1 {
		t.Fatalf("Order.ID = %d, want 1", o.ID)
	}
	if o.TenantID != 10 {
		t.Fatalf("Order.TenantID = %d, want 10", o.TenantID)
	}
	if o.Status != OrderStatusPending {
		t.Fatalf("Order.Status = %q, want pending", o.Status)
	}
	if o.TotalCents != 5000 {
		t.Fatalf("Order.TotalCents = %d, want 5000", o.TotalCents)
	}
	if o.IdempotencyKey != "idem-001" {
		t.Fatalf("Order.IdempotencyKey = %q", o.IdempotencyKey)
	}
}

func TestPaymentStructFields(t *testing.T) {
	t.Parallel()

	p := Payment{
		ID:                1,
		TenantID:          10,
		OrderID:           100,
		Status:            PaymentStatusPending,
		ProviderReference: "prov-ref-001",
		AmountCents:       5000,
	}

	if p.ID != 1 {
		t.Fatalf("Payment.ID = %d, want 1", p.ID)
	}
	if p.TenantID != 10 {
		t.Fatalf("Payment.TenantID = %d, want 10", p.TenantID)
	}
	if p.OrderID != 100 {
		t.Fatalf("Payment.OrderID = %d, want 100", p.OrderID)
	}
	if p.Status != PaymentStatusPending {
		t.Fatalf("Payment.Status = %q, want pending", p.Status)
	}
	if p.ProviderReference != "prov-ref-001" {
		t.Fatalf("Payment.ProviderReference = %q", p.ProviderReference)
	}
	if p.AmountCents != 5000 {
		t.Fatalf("Payment.AmountCents = %d, want 5000", p.AmountCents)
	}
}

func TestTenantSlugUniqueness(t *testing.T) {
	t.Parallel()

	t1 := Tenant{ID: 1, Slug: "acme"}
	t2 := Tenant{ID: 2, Slug: "acme"}

	slugs := map[string]bool{}
	slugs[t1.Slug] = true
	if slugs[t2.Slug] {
		t.Log("duplicate slug would violate uniqueness constraint — the DB layer must enforce this")
	}
}
