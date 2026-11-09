package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func generator(ctx context.Context, nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for _, n := range nums {
			select {
			case out <- n:
			case <-ctx.Done():
				return
			}
		}
	}()
	return out
}

func multiply(ctx context.Context, in <-chan int, factor int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range in {
			select {
			case out <- n * factor:
			case <-ctx.Done():
				return
			}
		}
	}()
	return out
}

func add(ctx context.Context, in <-chan int, delta int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range in {
			select {
			case out <- n + delta:
			case <-ctx.Done():
				return
			}
		}
	}()
	return out
}

func fanOut(ctx context.Context, in <-chan int, n int) []<-chan int {
	outs := make([]<-chan int, n)
	for i := 0; i < n; i++ {
		outs[i] = multiply(ctx, in, i+1)
	}
	return outs
}

func fanIn(ctx context.Context, channels ...<-chan int) <-chan int {
	out := make(chan int)
	var wg sync.WaitGroup
	for _, ch := range channels {
		wg.Add(1)
		go func(c <-chan int) {
			defer wg.Done()
			for n := range c {
				select {
				case out <- n:
				case <-ctx.Done():
					return
				}
			}
		}(ch)
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	fmt.Println("=== Basic pipeline ===")
	gen := generator(ctx, 1, 2, 3, 4, 5)
	mul := multiply(ctx, gen, 2)
	add2 := add(ctx, mul, 10)

	for v := range add2 {
		fmt.Println(v)
	}
	fmt.Println()

	fmt.Println("=== Fan-out / Fan-in ===")
	gen2 := generator(ctx, 1, 2, 3)
	fanned := fanOut(ctx, gen2, 3)
	merged := fanIn(ctx, fanned...)

	for v := range merged {
		fmt.Printf("merged: %d\n", v)
	}
}
