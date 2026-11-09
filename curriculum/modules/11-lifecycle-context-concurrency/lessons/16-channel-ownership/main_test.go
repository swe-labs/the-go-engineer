package main

import "testing"

func TestChannelOwnershipCompiles(t *testing.T) {
}

func TestOwnerConsumerPattern(t *testing.T) {
	ch := owner()
	results := consumer(ch)
	if len(results) != 5 {
		t.Fatalf("expected 5 results, got %d", len(results))
	}
}

func TestMultipleSendersCoordinator(t *testing.T) {
	results := splitOwnerWithMultipleSenders()
	if len(results) != 9 {
		t.Fatalf("expected 9 results, got %d", len(results))
	}
}

func TestChannelPassing(t *testing.T) {
	start := make(chan chan int)
	go func() {
		inner := <-start
		inner <- 42
	}()

	myCh := make(chan int)
	start <- myCh
	val := <-myCh
	if val != 42 {
		t.Fatalf("expected 42, got %d", val)
	}
}
