package main

import "testing"

func TestBufferedChannelsCompiles(t *testing.T) {
}

func TestBufferedChannelNonBlocking(t *testing.T) {
	ch := make(chan int, 2)
	ch <- 1
	ch <- 2
	if len(ch) != 2 {
		t.Fatalf("expected len 2, got %d", len(ch))
	}
}

func TestBufferedChannelBlocksWhenFull(t *testing.T) {
	ch := make(chan int, 1)
	ch <- 1
	select {
	case ch <- 2:
		t.Fatal("send should have blocked")
	default:
	}
}

func TestBufferedChannelBlocksWhenEmpty(t *testing.T) {
	ch := make(chan int, 1)
	select {
	case <-ch:
		t.Fatal("receive should have blocked")
	default:
	}
}
