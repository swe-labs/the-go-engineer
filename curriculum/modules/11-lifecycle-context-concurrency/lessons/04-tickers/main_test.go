package main

import (
	"testing"
	"time"
)

func TestTickerFires(t *testing.T) {
	ticker := time.NewTicker(20 * time.Millisecond)
	defer ticker.Stop()

	select {
	case <-ticker.C:
	case <-time.After(100 * time.Millisecond):
		t.Fatal("ticker did not fire within 100ms")
	}
}

func TestTickerMultipleFires(t *testing.T) {
	ticker := time.NewTicker(10 * time.Millisecond)
	defer ticker.Stop()

	count := 0
	timeout := time.After(60 * time.Millisecond)
	for count < 3 {
		select {
		case <-ticker.C:
			count++
		case <-timeout:
			t.Fatalf("only got %d ticks before timeout", count)
		}
	}
}

func TestTickerStop(t *testing.T) {
	ticker := time.NewTicker(10 * time.Millisecond)
	ticker.Stop()
	select {
	case <-ticker.C:
		t.Error("ticker fired after Stop")
	case <-time.After(50 * time.Millisecond):
	}
}

func TestTickerZeroDuration(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for zero-duration ticker")
		}
	}()
	time.NewTicker(0)
}
