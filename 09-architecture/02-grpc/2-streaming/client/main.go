// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// RUN: go run ./09-architecture/02-grpc/2-streaming/client
package main

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/rasel9t6/the-go-engineer/09-architecture/02-grpc/gen"
)

// ============================================================================
// Section 26: gRPC — Streaming Client
// Level: Expert
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Calling a server-side streaming method: recv() in a loop
//   - Calling a bidirectional streaming method: send() and recv() concurrently
//   - Using grpc.NewClient with insecure credentials for local dev
//   - Proper stream termination (CloseSend)
//
// RUN:
//   Terminal 1: go run ./09-architecture/02-grpc/2-streaming/server
//   Terminal 2: go run ./09-architecture/02-grpc/2-streaming/client
// ============================================================================

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// Connect to the gRPC server
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Error("failed to connect", "error", err)
		os.Exit(1)
	}
	defer conn.Close()

	client := pb.NewOrderServiceClient(conn)

	// =========================================================================
	// Demo 1: Server-side Streaming (WatchOrderStatus)
	// =========================================================================
	fmt.Println("\n--- [Demo 1: Server-side Streaming] ---")
	watchOrder(client, logger, "ord_123")

	// =========================================================================
	// Demo 2: Bidirectional Streaming (ProcessOrderStream)
	// =========================================================================
	fmt.Println("\n--- [Demo 2: Bidirectional Streaming] ---")
	processOrders(client, logger)
}

func watchOrder(client pb.OrderServiceClient, logger *slog.Logger, orderID string) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Call the server-streaming method
	stream, err := client.WatchOrderStatus(ctx, &pb.WatchOrderStatusRequest{OrderId: orderID})
	if err != nil {
		logger.Error("failed to watch order", "error", err)
		return
	}

	logger.Info("watching order status...", "order_id", orderID)

	for {
		// Receive the next update from the server
		update, err := stream.Recv()
		if err == io.EOF {
			// Server closed the stream
			logger.Info("server closed the watch stream")
			break
		}
		if err != nil {
			logger.Error("stream receive error", "error", err)
			return
		}

		logger.Info("status update received",
			"status", update.Status,
			"msg", update.Message,
			"timestamp", time.Unix(update.Timestamp, 0).Format(time.RFC3339),
		)
	}
}

func processOrders(client pb.OrderServiceClient, logger *slog.Logger) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Call the bidirectional streaming method
	stream, err := client.ProcessOrderStream(ctx)
	if err != nil {
		logger.Error("failed to start process stream", "error", err)
		return
	}

	// We'll use a wait group or a done channel to coordinate our send/recv
	done := make(chan struct{})

	// 1. Goroutine for receiving results from the server
	go func() {
		logger.Info("receiver: started")
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				logger.Info("receiver: server closed stream")
				close(done)
				return
			}
			if err != nil {
				logger.Error("receiver: stream error", "error", err)
				close(done)
				return
			}
			logger.Info("receiver: result received", "product_id", res.ProductId, "accepted", res.Accepted, "reason", res.Reason)
		}
	}()

	// 2. Main thread: send multiple items to the server
	logger.Info("sender: sending items...")
	items := []*pb.OrderItem{
		{ProductId: "prod_apple", Quantity: 2},
		{ProductId: "prod_banana", Quantity: 3}, // This should get rejected (odd quantity)
		{ProductId: "prod_cherry", Quantity: 10},
	}

	for _, item := range items {
		if err := stream.Send(item); err != nil {
			logger.Error("sender: failed to send item", "error", err)
			break
		}
		logger.Info("sender: item sent", "product_id", item.ProductId)
		time.Sleep(1 * time.Second)
	}

	// Important: close the sending side of the stream once we're done
	if err := stream.CloseSend(); err != nil {
		logger.Error("sender: failed to close send stream", "error", err)
	}
	logger.Info("sender: finished sending")

	// Wait for the receiver to finish
	<-done
}
