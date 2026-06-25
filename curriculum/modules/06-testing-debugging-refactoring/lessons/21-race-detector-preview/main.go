package main

import (
	"fmt"
	"sync"
)

type Counter struct {
	mu    sync.Mutex
	value int
}

func main() {
	var wg sync.WaitGroup
	c := &Counter{}

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			c.mu.Lock()
			c.value++
			c.mu.Unlock()
		}()
	}

	wg.Wait()
	fmt.Println("Counter:", c.value)
}

// UnsafeIncrement demonstrates a data race (used in tests).
func UnsafeIncrement(c *Counter) {
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			c.value++ // intentional race
		}()
	}
	wg.Wait()
}
