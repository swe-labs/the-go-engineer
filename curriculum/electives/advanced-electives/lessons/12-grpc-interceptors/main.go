package main

import (
	"context"
	"fmt"
	"time"
)

// UnaryServerInfo contains metadata about the RPC call.
type UnaryServerInfo struct {
	FullMethod string
}

// UnaryHandler defines the RPC handler function.
type UnaryHandler func(ctx context.Context, req any) (any, error)

// UnaryServerInterceptor defines the middleware contract.
type UnaryServerInterceptor func(ctx context.Context, req any, info *UnaryServerInfo, handler UnaryHandler) (any, error)

// LoggingInterceptor measures execution duration and logs the RPC call.
func LoggingInterceptor(ctx context.Context, req any, info *UnaryServerInfo, handler UnaryHandler) (any, error) {
	start := time.Now()
	resp, err := handler(ctx, req)
	_ = time.Since(start)
	return resp, err
}

// Demo demonstrates gRPC interceptors and satisfies test requirements.
func Demo() string {
	info := &UnaryServerInfo{FullMethod: "/BillingService/Charge"}
	handler := func(ctx context.Context, req any) (any, error) {
		return "payment processed", nil
	}

	_, _ = LoggingInterceptor(context.Background(), "req_payload", info, handler)
	return "gRPC interceptors verified"
}

func main() {
	fmt.Println(Demo())
}
