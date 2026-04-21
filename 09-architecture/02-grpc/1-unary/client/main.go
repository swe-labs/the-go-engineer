// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// RUN: go run ./09-architecture/02-grpc/1-unary/client
package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	pb "github.com/rasel9t6/the-go-engineer/09-architecture/02-grpc/gen"
)

// ============================================================================
// Section 26: gRPC — Unary Client
// Level: Beginner
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Creating a gRPC connection and client stub
//   - Using context timeouts with gRPC (essential — no default timeout)
//   - Sending metadata (headers) with a request
//   - Inspecting gRPC errors: extracting status codes
//   - The generated client stub: one method per RPC
//
// ENGINEERING DEPTH:
//   Unlike net/http which has configurable timeouts on the Client struct,
//   gRPC has NO default deadline. A call to CreateOrder() without a context
//   timeout will hang forever if the server is unresponsive. This is why
//   EVERY gRPC call must use context.WithTimeout() or context.WithDeadline().
//
//   The client stub is type-safe: if you change the proto and regenerate,
//   any call site using the old API will fail to compile. This is the key
//   advantage over REST with hand-written clients that only fail at runtime.
//
// RUN:
//   Terminal 1: go run ./09-architecture/02-grpc/1-unary/server
//   Terminal 2: go run ./09-architecture/02-grpc/1-unary/client
// ============================================================================

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// =========================================================================
	// 1. Establish a gRPC connection (channel)
	// =========================================================================
	// grpc.Dial creates a connection to the server.
	// PRODUCTION: Replace insecure.NewCredentials() with TLS credentials.
	//   creds, err := credentials.NewClientTLSFromFile("cert.pem", "")
	conn, err := grpc.NewClient("localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		logger.Error("failed to connect", "error", err)
		os.Exit(1)
	}
	defer conn.Close()

	// =========================================================================
	// 2. Create the typed client stub (generated code)
	// =========================================================================
	// The stub has one method for each RPC defined in the .proto file.
	// It wraps the raw gRPC connection with type-safe call methods.
	client := pb.NewOrderServiceClient(conn)

	// =========================================================================
	// 3. Attach metadata (equivalent to HTTP request headers)
	// =========================================================================
	// metadata.AppendToOutgoingContext adds k/v pairs to the gRPC request context.
	// The server reads them with metadata.FromIncomingContext(ctx).
	ctx := metadata.AppendToOutgoingContext(context.Background(),
		"x-request-id", "req_client_001",
		"x-user-agent", "go-engineer-client/1.0",
	)

	// =========================================================================
	// 4. Unary call: CreateOrder
	// =========================================================================
	// CRITICAL: Always set a deadline. gRPC has no default timeout.
	callCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	logger.Info("creating order...")
	resp, err := client.CreateOrder(callCtx, &pb.CreateOrderRequest{
		CustomerId: "cust_42",
		ProductId:  "prod_laptop_x1",
		Quantity:   2,
	})
	if err != nil {
		handleGRPCError(logger, "CreateOrder", err)
		os.Exit(1)
	}

	logger.Info("order created",
		"order_id", resp.OrderId,
		"status", resp.Status,
		"total_cost", fmt.Sprintf("$%.2f", resp.TotalCost),
	)

	// =========================================================================
	// 5. Unary call: GetOrder
	// =========================================================================
	getCtx, cancel2 := context.WithTimeout(ctx, 5*time.Second)
	defer cancel2()

	order, err := client.GetOrder(getCtx, &pb.GetOrderRequest{
		OrderId: resp.OrderId,
	})
	if err != nil {
		handleGRPCError(logger, "GetOrder", err)
		os.Exit(1)
	}

	logger.Info("order retrieved",
		"order_id", order.OrderId,
		"customer_id", order.CustomerId,
		"quantity", order.Quantity,
		"status", order.Status,
	)

	// =========================================================================
	// 6. Demonstrate error handling: not found
	// =========================================================================
	_, err = client.GetOrder(getCtx, &pb.GetOrderRequest{OrderId: "ord_nonexistent"})
	handleGRPCError(logger, "GetOrder (not found)", err)

	// =========================================================================
	// 7. Demonstrate error handling: invalid argument
	// =========================================================================
	_, err = client.CreateOrder(callCtx, &pb.CreateOrderRequest{
		CustomerId: "", // Missing required field
		Quantity:   0,
	})
	handleGRPCError(logger, "CreateOrder (invalid)", err)
	fmt.Println("\n---------------------------------------------------")
	fmt.Println("🚀 NEXT UP: GS.1 signal.NotifyContext")
	fmt.Println("   Current: GR.3 (unary client)")
	fmt.Println("---------------------------------------------------")
}

// handleGRPCError extracts the gRPC status code and message from an error.
// Use status.FromError() — never switch on the error string.
func handleGRPCError(logger *slog.Logger, method string, err error) {
	if err == nil {
		return
	}

	s, ok := status.FromError(err)
	if !ok {
		// Not a gRPC error — could be a network failure
		logger.Error("non-gRPC error", "method", method, "error", err)
		return
	}

	switch s.Code() {
	case codes.NotFound:
		logger.Warn("resource not found", "method", method, "message", s.Message())
	case codes.InvalidArgument:
		logger.Warn("invalid request", "method", method, "message", s.Message())
	case codes.DeadlineExceeded:
		logger.Error("request timed out", "method", method)
	case codes.Unavailable:
		logger.Error("service unavailable — retry later", "method", method)
	default:
		logger.Error("rpc failed", "method", method, "code", s.Code(), "message", s.Message())
	}

	// KEY TAKEAWAY:
	// - grpc.NewClient() creates the connection; pb.NewXxxServiceClient(conn) creates the stub
	// - ALWAYS use context.WithTimeout() — gRPC has no default deadline
	// - metadata.AppendToOutgoingContext adds headers to the outgoing request
	// - status.FromError(err) extracts the gRPC status code safely
	// - Switch on status.Code(), never on the error message string
}
