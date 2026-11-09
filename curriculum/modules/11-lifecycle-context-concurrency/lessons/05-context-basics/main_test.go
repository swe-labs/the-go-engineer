package main

import (
	"context"
	"testing"
	"time"
)

func TestBackgroundContext(t *testing.T) {
	ctx := context.Background()
	if ctx.Err() != nil {
		t.Errorf("Background().Err() = %v, want nil", ctx.Err())
	}
	if ctx.Done() != nil {
		t.Error("Background().Done() should be nil")
	}
}

func TestWithCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	if ctx.Err() != nil {
		t.Errorf("before cancel: Err() = %v", ctx.Err())
	}
	cancel()
	if ctx.Err() != context.Canceled {
		t.Errorf("after cancel: Err() = %v, want Canceled", ctx.Err())
	}
}

func TestWithTimeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	<-ctx.Done()
	if ctx.Err() != context.DeadlineExceeded {
		t.Errorf("after timeout: Err() = %v, want DeadlineExceeded", ctx.Err())
	}
}

func TestWithTimeoutCancelFirst(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Hour)
	cancel()
	if ctx.Err() != context.Canceled {
		t.Errorf("cancel before timeout: Err() = %v, want Canceled", ctx.Err())
	}
}

func TestBackgroundIsNilErr(t *testing.T) {
	ctx := context.Background()
	if ctx.Err() != nil {
		t.Errorf("Background().Err() = %v, want nil", ctx.Err())
	}
}
