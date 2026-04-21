// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// RUN: go run ./09-architecture/02-grpc/1-unary/server
package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	pb "github.com/rasel9t6/the-go-engineer/09-architecture/02-grpc/gen"
)

// ============================================================================
// Section 26: gRPC — Unary Server
// Level: Beginner → Intermediate
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Implementing a gRPC server from a generated interface
//   - gRPC status codes: the equivalent of HTTP status codes
//   - Interceptors: the gRPC equivalent of HTTP middleware
//   - Reading metadata from the gRPC context (headers)
//
// ENGINEERING DEPTH:
//   The generated OrderServiceServer interface mandates that every RPC method
//   you defined in the .proto file is implemented. The compiler enforces this —
//   forget a method and your build breaks. This is exactly the "interfaces as
//   contracts" pattern from Section 05, applied at the network boundary.
//
//   gRPC STATUS CODES — the most important ones:
//     codes.OK           — success
//     codes.NotFound     — resource doesn't exist
//     codes.InvalidArgument — bad request (wrong data type, missing field)
//     codes.Unauthenticated — missing/invalid credentials
//     codes.PermissionDenied — valid credentials, insufficient permissions
//     codes.Internal     — server-side error (like HTTP 500)
//     codes.Unavailable  — service is temporarily down (retry-able)
//     codes.DeadlineExceeded — context timed out before completion
//
// RUN:
//   Terminal 1: go run ./09-architecture/02-grpc/1-unary/server
//   Terminal 2: go run ./09-architecture/02-grpc/1-unary/client
// ============================================================================

// OrderServer implements the generated pb.OrderServiceServer interface.
// The server struct holds any dependencies (database, cache, logger).
type OrderServer struct {
	pb.UnimplementedOrderServiceServer // Embed this: auto-implements new methods as "unimplemented"
	orders                             map[string]*pb.GetOrderResponse
	logger                             *slog.Logger
}

// CreateOrder handles unary RPC: one request → one response.
func (s *OrderServer) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	// Read incoming metadata (equivalent to HTTP request headers)
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if vals := md.Get("x-request-id"); len(vals) > 0 {
			s.logger.Info("received request", "request_id", vals[0])
		}
	}

	// Input validation — return codes.InvalidArgument for bad inputs
	if req.CustomerId == "" {
		return nil, status.Error(codes.InvalidArgument, "customer_id is required")
	}
	if req.Quantity <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "quantity must be positive, got %d", req.Quantity)
	}

	// Simulate order creation
	orderID := fmt.Sprintf("ord_%d", time.Now().UnixNano())
	totalCost := float64(req.Quantity) * 29.99

	order := &pb.GetOrderResponse{
		OrderId:    orderID,
		CustomerId: req.CustomerId,
		ProductId:  req.ProductId,
		Quantity:   req.Quantity,
		Status:     "pending",
		CreatedAt:  time.Now().Unix(),
	}
	s.orders[orderID] = order

	s.logger.Info("order created",
		"order_id", orderID,
		"customer_id", req.CustomerId,
		"total_cost", totalCost,
	)

	return &pb.CreateOrderResponse{
		OrderId:   orderID,
		Status:    "created",
		TotalCost: totalCost,
	}, nil
}

// GetOrder retrieves an existing order by ID.
func (s *OrderServer) GetOrder(ctx context.Context, req *pb.GetOrderRequest) (*pb.GetOrderResponse, error) {
	if req.OrderId == "" {
		return nil, status.Error(codes.InvalidArgument, "order_id is required")
	}

	order, ok := s.orders[req.OrderId]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "order %q not found", req.OrderId)
	}

	return order, nil
}

// ============================================================================
// Interceptors — gRPC middleware
// ============================================================================
// Interceptors wrap every RPC call with cross-cutting concerns:
// authentication, logging, metrics, rate limiting, panic recovery.

// loggingInterceptor logs every unary RPC call with method name and latency.
func loggingInterceptor(logger *slog.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		start := time.Now()
		resp, err := handler(ctx, req) // Call the actual RPC handler
		latency := time.Since(start)

		code := codes.OK
		if err != nil {
			code = status.Code(err)
		}

		logger.Info("rpc call",
			"method", info.FullMethod,
			"code", code,
			"latency", latency,
		)
		return resp, err
	}
}

// recoveryInterceptor catches panics inside handlers and returns codes.Internal.
// Without this, a panic in any handler crashes the entire server process.
func recoveryInterceptor(logger *slog.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp any, err error) {
		defer func() {
			if r := recover(); r != nil {
				logger.Error("panic in rpc handler", "method", info.FullMethod, "panic", r)
				err = status.Errorf(codes.Internal, "internal server error")
			}
		}()
		return handler(ctx, req)
	}
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// Listen for incoming TCP connections
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		logger.Error("failed to listen", "error", err)
		os.Exit(1)
	}

	// Create a gRPC server with interceptors chained together.
	// Interceptors execute in the order they are added (like HTTP middleware).
	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			recoveryInterceptor(logger), // Outermost: catches panics everywhere
			loggingInterceptor(logger),  // Logs every call with timing
		),
	)

	// Register our OrderServer as the implementation for OrderService.
	orderSvc := &OrderServer{
		orders: make(map[string]*pb.GetOrderResponse),
		logger: logger,
	}
	pb.RegisterOrderServiceServer(server, orderSvc)

	logger.Info("gRPC server listening", "addr", ":50051")
	if err := server.Serve(lis); err != nil {
		logger.Error("server failed", "error", err)
		os.Exit(1)
	}

	// KEY TAKEAWAY:
	// - Implement the generated Server interface — one method per RPC
	// - Use status.Error(code, msg) not errors.New() for gRPC errors
	// - Embed UnimplementedOrderServiceServer — future-proofs against new RPCs
	// - Chain interceptors for logging, auth, recovery (gRPC middleware)
	// - metadata.FromIncomingContext reads gRPC headers
	fmt.Println("\n---------------------------------------------------")
	fmt.Println("🚀 NEXT UP: GR.3 unary client")
	fmt.Println("   Current: GR.2 (unary server)")
	fmt.Println("---------------------------------------------------")
}
