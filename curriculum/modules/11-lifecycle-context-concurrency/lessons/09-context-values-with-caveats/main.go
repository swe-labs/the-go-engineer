package main

import (
	"context"
	"fmt"
)

type contextKey string

const (
	keyUserID    contextKey = "user_id"
	keyRequestID contextKey = "request_id"
)

func handler(ctx context.Context) {
	userID, _ := ctx.Value(keyUserID).(string)
	reqID, _ := ctx.Value(keyRequestID).(string)
	fmt.Printf("Handling request %s for user %s\n", reqID, userID)
}

func main() {
	ctx := context.Background()
	ctx = context.WithValue(ctx, keyUserID, "user-42")
	ctx = context.WithValue(ctx, keyRequestID, "req-abc-123")
	handler(ctx)
}
