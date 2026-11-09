package main

import (
	"context"
	"log/slog"
	"os"
)

type cidKey struct{}

// WithCorrelationID stores a correlation ID in context.
func WithCorrelationID(ctx context.Context, cid string) context.Context {
	return context.WithValue(ctx, cidKey{}, cid)
}

// GetCorrelationID retrieves the correlation ID from context.
func GetCorrelationID(ctx context.Context) string {
	cid, _ := ctx.Value(cidKey{}).(string)
	return cid
}

// SimulateServiceCall logs with correlation ID context.
func SimulateServiceCall(ctx context.Context, service, action string) {
	slog.Info("service call",
		"service", service,
		"action", action,
		"correlation_id", GetCorrelationID(ctx),
	)
}

// CallServices simulates propagating a correlation ID through a chain of services.
func CallServices(ctx context.Context, services []string) []string {
	var chain []string
	cid := GetCorrelationID(ctx)
	for _, svc := range services {
		chain = append(chain, cid)
		slog.Info("entering service", "service", svc, "correlation_id", cid)
	}
	return chain
}

func main() {
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	slog.SetDefault(slog.New(handler))

	ctx := WithCorrelationID(context.Background(), "corr-abc-123")
	SimulateServiceCall(ctx, "api-gateway", "authenticate")
	SimulateServiceCall(ctx, "user-service", "get_profile")
	SimulateServiceCall(ctx, "payment-service", "charge")
	CallServices(ctx, []string{"auth", "orders", "inventory", "notifications"})
}
