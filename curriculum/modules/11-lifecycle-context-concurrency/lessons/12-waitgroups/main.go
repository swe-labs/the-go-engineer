package main

import (
	"fmt"
	"sync"
	"time"
)

func fetchUser(id int) {
	time.Sleep(time.Duration(id) * 50 * time.Millisecond)
	fmt.Printf("User %d fetched\n", id)
}

type WorkItem struct {
	ID   int
	Data string
}

type Result struct {
	Item WorkItem
	Err  error
}

func processItems(items []WorkItem) []Result {
	var wg sync.WaitGroup
	results := make([]Result, len(items))

	for i, item := range items {
		wg.Add(1)
		go func(idx int, wi WorkItem) {
			defer wg.Done()
			// Simulate processing
			time.Sleep(20 * time.Millisecond)
			results[idx] = Result{Item: wi}
		}(i, item)
	}

	wg.Wait()
	return results
}

func main() {
	fmt.Println("=== Basic WaitGroup pattern ===")
	var wg sync.WaitGroup

	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			fetchUser(id)
		}(i)
	}

	wg.Wait()
	fmt.Println("All fetches complete")

	fmt.Println("\n=== Typed WaitGroup pattern ===")
	items := []WorkItem{
		{ID: 1, Data: "alpha"},
		{ID: 2, Data: "beta"},
		{ID: 3, Data: "gamma"},
	}

	results := processItems(items)
	for _, r := range results {
		fmt.Printf("Result: %+v\n", r.Item)
	}

	fmt.Println("\n=== Multiple batches ===")
	for batch := 0; batch < 3; batch++ {
		var batchWG sync.WaitGroup
		for j := 0; j < 3; j++ {
			batchWG.Add(1)
			go func(b, j int) {
				defer batchWG.Done()
				fmt.Printf("Batch %d, job %d\n", b, j)
			}(batch, j)
		}
		batchWG.Wait()
		fmt.Printf("Batch %d complete\n", batch)
	}
}
