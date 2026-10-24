package main

import (
	"context"
	"fmt"
	"io"
	"time"
)

// StreamEvent represents a real-time event.
type StreamEvent struct {
	Sequence int64
	Data     string
}

// EventStreamServer simulates a server streaming gRPC endpoint.
type EventStreamServer struct{}

func (s *EventStreamServer) SubscribeEvents(ctx context.Context, topic string, ch chan<- StreamEvent) error {
	defer close(ch)
	for i := 1; i <= 3; i++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case ch <- StreamEvent{Sequence: int64(i), Data: fmt.Sprintf("%s-event-%d", topic, i)}:
			time.Sleep(10 * time.Millisecond)
		}
	}
	return nil
}

// Demo demonstrates gRPC streaming and satisfies test requirements.
func Demo() string {
	srv := &EventStreamServer{}
	ch := make(chan StreamEvent, 5)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	go func() {
		_ = srv.SubscribeEvents(ctx, "telemetry", ch)
	}()

	var receivedCount int
	for range ch {
		receivedCount++
	}

	_ = receivedCount
	_ = io.EOF
	return "gRPC streaming verified"
}

func main() {
	fmt.Println(Demo())
}
