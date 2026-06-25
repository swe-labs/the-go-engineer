package main

import (
	"fmt"
	"sync"
)

type SharedCounter struct {
	mu    sync.Mutex
	value int
}

func (c *SharedCounter) Add(n int) {
	c.mu.Lock()
	c.value += n
	c.mu.Unlock()
}

func (c *SharedCounter) Value() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.value
}

func racyAdd(c *SharedCounter, n int) {
	// Direct write without lock simulates a race
	c.value += n
}

func main() {
	var wg sync.WaitGroup
	safe := &SharedCounter{}
	n := 1000

	wg.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			safe.Add(1)
		}()
	}
	wg.Wait()

	fmt.Printf("Safe counter: %d (expected %d)\n", safe.Value(), n)

	fmt.Println("\nRun this file with: go run -race .")
	fmt.Println("Run tests with: go test -race .")
}
