// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// RUN: go run ./14-application-architecture/grpc/2-streaming/server
package main

import (
	"fmt"
	"io"
	"log/slog"
	"net"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/rasel9t6/the-go-engineer/14-application-architecture/grpc/gen"
)

// ============================================================================
// Section 26: gRPC — Streaming (Server & Bidirectional)
// Level: Expert
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Server-side streaming: one request → many responses
//   - Bidirectional streaming: many requests ↔ many responses
//   - Handling stream closures and EOF (End Of File)
//   - Concurrency patterns with gRPC streams
//
// ENGINEERING DEPTH:
//   Unary RPC is familiar because it looks like HTTP/1.1. Streaming is where
//   gRPC's usage of HTTP/2 shines. A single TCP connection is multiplexed
//   to handle multiple concurrent streams, each capable of pushing data
//   independently.
//
//   SERVER STREAMING (WatchOrderStatus):
//     Useful for "push" notifications. Instead of the client polling every 5s,
//     the server holds the stream open and sends a message only when state
//     changes. This reduces latency and CPU usage on both ends.
//
//   BIDIRECTIONAL STREAMING (ProcessOrderStream):
//     The most powerful RPC type. The client sends a stream of items, and the
//     server responds to each. This is perfect for high-throughput batch
//     processing where you want to know the result of item 1 while still
//     uploading item 500.
//
// RUN:
//   Terminal 1: go run ./14-application-architecture/grpc/2-streaming/server
//   Terminal 2: go run ./14-application-architecture/grpc/2-streaming/client
// ============================================================================

type StreamingOrderServer struct {
	pb.UnimplementedOrderServiceServer
	logger *slog.Logger
}

// WatchOrderStatus implements server-side streaming.
// The client sends one request, and we send back multiple status updates.
func (s *StreamingOrderServer) WatchOrderStatus(req *pb.WatchOrderStatusRequest, stream pb.OrderService_WatchOrderStatusServer) error {
	if req.OrderId == "" {
		return status.Error(codes.InvalidArgument, "order_id is required")
	}

	statuses := []string{"pending", "processing", "packaging", "shipped", "delivered"}

	s.logger.Info("started watching order", "order_id", req.OrderId)

	for _, st := range statuses {
		// Prepare the update
		update := &pb.OrderStatusUpdate{
			OrderId:   req.OrderId,
			Status:    st,
			Message:   fmt.Sprintf("Order status changed to %s", st),
			Timestamp: time.Now().Unix(),
		}

		// Push the update over the stream
		if err := stream.Send(update); err != nil {
			s.logger.Error("failed to send status update", "error", err)
			return err
		}

		s.logger.Debug("pushed status update", "order_id", req.OrderId, "status", st)

		// Simulate some work between status changes
		time.Sleep(2 * time.Second)
	}

	s.logger.Info("finished watching order", "order_id", req.OrderId)
	return nil
}

// ProcessOrderStream implements bidirectional streaming.
// The client sends a stream of OrderItems; the server sends back a stream of ProcessResults.
func (s *StreamingOrderServer) ProcessOrderStream(stream pb.OrderService_ProcessOrderStreamServer) error {
	s.logger.Info("started processing bidirectional order stream")

	for {
		// Read the next item from the client's stream
		item, err := stream.Recv()
		if err == io.EOF {
			// Client is done sending. We return nil to close our side too.
			s.logger.Info("client finished stream")
			return nil
		}
		if err != nil {
			s.logger.Error("stream receive error", "error", err)
			return err
		}

		s.logger.Info("received item to process", "product_id", item.ProductId, "qty", item.Quantity)

		// Simulate business logic (accepting even quantities, rejecting odd)
		result := &pb.ProcessResult{
			ProductId: item.ProductId,
			Accepted:  true,
			Reason:    "Order validated and processed",
		}
		if item.Quantity%2 != 0 {
			result.Accepted = false
			result.Reason = "Odd quantities require manual approval"
		}

		// Send the result back to the client immediately
		if err := stream.Send(result); err != nil {
			s.logger.Error("stream send error", "error", err)
			return err
		}
	}
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		logger.Error("failed to listen", "error", err)
		os.Exit(1)
	}

	server := grpc.NewServer()
	pb.RegisterOrderServiceServer(server, &StreamingOrderServer{logger: logger})

	logger.Info("Streaming gRPC server listening", "addr", ":50051")
	if err := server.Serve(lis); err != nil {
		logger.Error("server failed", "error", err)
		os.Exit(1)
	}
}
