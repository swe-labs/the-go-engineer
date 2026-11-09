package main

import (
	"testing"
	"time"
)

func TestSequential(t *testing.T) {
	tasks := []task{slowSquare, slowDouble}
	result := sequential(tasks, 3)
	// square: 3*3=9, double: 9*2=18
	if result != 18 {
		t.Errorf("sequential = %d, want %d", result, 18)
	}
}

func TestSequentialTiming(t *testing.T) {
	tasks := []task{slowSquare, slowDouble}
	start := time.Now()
	sequential(tasks, 3)
	elapsed := time.Since(start)
	// Each task sleeps 10ms, so sequential should be ~20ms
	if elapsed < 19*time.Millisecond {
		t.Errorf("sequential too fast: %v, expected ~20ms", elapsed)
	}
}
