package main

import (
	"testing"
	"time"
)

func TestChannelsCompiles(t *testing.T) {
}

func TestUnbufferedChannelSendReceive(t *testing.T) {
	ch := make(chan int)
	go func() { ch <- 1 }()
	got := <-ch
	if got != 1 {
		t.Fatalf("expected 1, got %d", got)
	}
}

func TestUnbufferedChannelBlockingEh(t *testing.T) {
	ch := make(chan struct{})
	blocked := make(chan bool, 1)

	go func() {
		ch <- struct{}{}
		blocked <- true
	}()

	time.Sleep(10 * time.Millisecond)
	select {
	case <-blocked:
		t.Fatal("send should have blocked until receive")
	default:
	}

	<-ch
	<-blocked
}
