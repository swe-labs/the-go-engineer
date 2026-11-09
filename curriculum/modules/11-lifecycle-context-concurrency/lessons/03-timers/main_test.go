package main

import (
	"testing"
	"time"
)

func TestTimerFires(t *testing.T) {
	timer := time.NewTimer(10 * time.Millisecond)
	select {
	case <-timer.C:
	case <-time.After(100 * time.Millisecond):
		t.Fatal("timer did not fire within 100ms")
	}
}

func TestTimerStop(t *testing.T) {
	timer := time.NewTimer(1 * time.Hour)
	if !timer.Stop() {
		t.Error("Stop returned false on active timer")
	}
	select {
	case <-timer.C:
		t.Error("timer fired after being stopped")
	default:
	}
}

func TestTimerReset(t *testing.T) {
	timer := time.NewTimer(10 * time.Millisecond)
	<-timer.C
	timer.Reset(10 * time.Millisecond)
	select {
	case <-timer.C:
	case <-time.After(100 * time.Millisecond):
		t.Fatal("reset timer did not fire")
	}
}

func TestTimerAfter(t *testing.T) {
	start := time.Now()
	<-time.After(20 * time.Millisecond)
	elapsed := time.Since(start)
	if elapsed < 15*time.Millisecond {
		t.Errorf("time.After fired too soon: %v", elapsed)
	}
}
