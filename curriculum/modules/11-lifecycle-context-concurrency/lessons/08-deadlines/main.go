package main

import (
	"context"
	"errors"
	"fmt"
	"time"
)

func operationWithDeadline(ctx context.Context) error {
	deadline, ok := ctx.Deadline()
	if ok {
		fmt.Println("Deadline set:", deadline.Format(time.RFC3339))
		fmt.Println("Time until deadline:", time.Until(deadline).Round(time.Millisecond))
	} else {
		fmt.Println("No deadline set")
	}

	select {
	case <-time.After(1 * time.Second):
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func main() {
	deadline := time.Now().Add(200 * time.Millisecond)
	ctx, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel()

	err := operationWithDeadline(ctx)
	if errors.Is(err, context.DeadlineExceeded) {
		fmt.Println("Deadline exceeded")
	} else if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Completed before deadline")
	}
}
