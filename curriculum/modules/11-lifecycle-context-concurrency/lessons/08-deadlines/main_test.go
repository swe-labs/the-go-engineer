package main

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestDeadlineExceeded(t *testing.T) {
	deadline := time.Now().Add(10 * time.Millisecond)
	ctx, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel()

	err := operationWithDeadline(ctx)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("expected DeadlineExceeded, got %v", err)
	}
}

func TestDeadlineNotExceeded(t *testing.T) {
	deadline := time.Now().Add(5 * time.Second)
	ctx, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel()

	fast := make(chan error, 1)
	go func() {
		fast <- operationWithDeadline(ctx)
	}()

	select {
	case err := <-fast:
		if err != nil {
			t.Errorf("expected nil, got %v", err)
		}
	case <-time.After(50 * time.Millisecond):
	}
}

func TestDeadlineWithTimeoutEquivalence(t *testing.T) {
	t1, cancel1 := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel1()

	t2, cancel2 := context.WithDeadline(context.Background(), time.Now().Add(50*time.Millisecond))
	defer cancel2()

	<-t1.Done()
	<-t2.Done()

	if t1.Err() != t2.Err() {
		t.Errorf("WithTimeout err = %v, WithDeadline err = %v", t1.Err(), t2.Err())
	}
}

func TestDeadlinePast(t *testing.T) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(-1*time.Second))
	defer cancel()

	<-ctx.Done()
	if !errors.Is(ctx.Err(), context.DeadlineExceeded) {
		t.Errorf("Err() = %v, want DeadlineExceeded", ctx.Err())
	}
}

func TestDeadlineReturnsCorrectTime(t *testing.T) {
	deadline := time.Now().Add(1 * time.Hour)
	ctx, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel()

	got, ok := ctx.Deadline()
	if !ok {
		t.Fatal("Deadline() returned ok=false")
	}
	if !got.Equal(deadline) {
		t.Errorf("Deadline() = %v, want %v", got, deadline)
	}
}
