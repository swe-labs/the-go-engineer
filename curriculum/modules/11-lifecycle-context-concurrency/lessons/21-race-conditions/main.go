package main

import (
	"fmt"
	"sync"
)

type RacyCounter struct {
	value int
}

func (c *RacyCounter) Add(n int) {
	c.value += n
}

func (c *RacyCounter) Value() int {
	return c.value
}

type SafeCounter struct {
	mu    sync.Mutex
	value int
}

func (c *SafeCounter) Add(n int) {
	c.mu.Lock()
	c.value += n
	c.mu.Unlock()
}

func (c *SafeCounter) Value() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.value
}

func main() {
	var wg sync.WaitGroup
	racy := &RacyCounter{}
	safe := &SafeCounter{}
	n := 1000

	wg.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			racy.Add(1)
			safe.Add(1)
		}()
	}
	wg.Wait()

	fmt.Printf("Expected: %d\n", n)
	fmt.Printf("Racy counter: %d (may differ due to data race)\n", racy.Value())
	fmt.Printf("Safe counter: %d (correct)\n", safe.Value())
}
