package payment

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/rasel9t6/the-go-engineer/11-flagship/01-opslane/internal/models"
)

func TestSimulatedGatewaySettlesByDefault(t *testing.T) {
	t.Parallel()

	gateway := NewSimulatedGateway()
	result, err := gateway.Charge(context.Background(), GatewayRequest{
		TenantID:          7,
		OrderID:           101,
		ProviderReference: "pay_123",
		AmountCents:       2500,
	})
	if err != nil {
		t.Fatalf("Charge returned error: %v", err)
	}

	if result.Status != models.PaymentStatusSettled {
		t.Fatalf("status = %q, want settled", result.Status)
	}
}

func TestSimulatedGatewayHonorsContextTimeout(t *testing.T) {
	t.Parallel()

	gateway := SimulatedGateway{Delay: 50 * time.Millisecond}
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
	defer cancel()

	_, err := gateway.Charge(ctx, GatewayRequest{
		TenantID:          7,
		OrderID:           101,
		ProviderReference: "pay_timeout",
		AmountCents:       2500,
	})
	if !errors.Is(err, ErrGatewayTimeout) {
		t.Fatalf("Charge error = %v, want ErrGatewayTimeout", err)
	}
}
