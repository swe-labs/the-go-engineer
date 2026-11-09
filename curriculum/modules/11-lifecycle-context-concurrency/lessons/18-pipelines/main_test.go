package main

import (
	"context"
	"testing"
)

func TestPipelinesCompiles(t *testing.T) {
}

func TestBasicPipeline(t *testing.T) {
	ctx := context.Background()
	gen := generator(ctx, 1, 2, 3)
	mul := multiply(ctx, gen, 3)
	add2 := add(ctx, mul, 1)

	var results []int
	for v := range add2 {
		results = append(results, v)
	}
	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}
	// 1*3+1=4, 2*3+1=7, 3*3+1=10
	expected := []int{4, 7, 10}
	for i, v := range results {
		if v != expected[i] {
			t.Errorf("index %d: expected %d, got %d", i, expected[i], v)
		}
	}
}

func TestFanOutFanIn(t *testing.T) {
	ctx := context.Background()
	gen := generator(ctx, 1, 2)
	fanned := fanOut(ctx, gen, 2)
	merged := fanIn(ctx, fanned...)

	count := 0
	for range merged {
		count++
	}
	if count != 2 {
		t.Fatalf("expected 2 values, got %d", count)
	}
}

func TestPipelineCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	gen := generator(ctx, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10)

	// read a few then cancel
	val := <-gen
	if val != 1 {
		t.Fatalf("expected 1, got %d", val)
	}
	cancel()

	// should stop quickly
	done := make(chan struct{})
	go func() {
		for range gen {
		}
		close(done)
	}()
	<-done
}
