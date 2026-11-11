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

	// Deliberate data race: 100 goroutines mutate SharedCounter.value with no
	// lock. Without -race this may still print 1000 (races are intermittent);
	// with `go run -race .` the detector reports the read-write race and the
	// program exits with status 66. That report is the demonstration.
	racyWG := sync.WaitGroup{}
	racyWG.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			defer racyWG.Done()
			racyAdd(safe, 1)
		}()
	}
	racyWG.Wait()
	fmt.Printf("Racy counter (racyAdd, no lock): %d\n", safe.Value())

	fmt.Println("\nRun without -race: the racy counter may silently print a wrong value.")
	fmt.Println("Run with: go run -race . - the detector prints the race report and exits with status 66.")
	fmt.Println("Run tests with: go test -race . (the tests themselves are race-free).")
}
