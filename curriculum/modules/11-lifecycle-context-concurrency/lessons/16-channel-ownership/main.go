package main

import (
	"fmt"
	"sync"
)

// owner creates the channel, sends values, and closes it.
func owner() <-chan int {
	ch := make(chan int, 5)
	go func() {
		defer close(ch)
		for i := 1; i <= 5; i++ {
			ch <- i * 10
		}
	}()
	return ch
}

// consumer only receives. It never sends or closes.
func consumer(ch <-chan int) []int {
	var results []int
	for v := range ch {
		results = append(results, v)
	}
	return results
}

// splitOwnerWithMultipleSenders uses a coordinator goroutine that knows
// when all senders are done, then closes the channel.
func splitOwnerWithMultipleSenders() []int {
	ch := make(chan int)
	var wg sync.WaitGroup

	// Multiple producers
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 3; j++ {
				ch <- id*10 + j
			}
		}(i)
	}

	// Coordinator closes when all producers finish
	go func() {
		wg.Wait()
		close(ch)
	}()

	var results []int
	for v := range ch {
		results = append(results, v)
	}
	return results
}

func main() {
	fmt.Println("=== Owner creates, sends, closes ===")
	ch := owner()
	results := consumer(ch)
	fmt.Println("Results:", results)

	fmt.Println("\n=== Multiple senders with coordinator ===")
	results2 := splitOwnerWithMultipleSenders()
	fmt.Println("Multi-sender results:", results2)

	fmt.Println("\n=== Channel passing between goroutines ===")
	start := make(chan chan int)
	go func() {
		inner := <-start
		inner <- 99
	}()

	myCh := make(chan int)
	start <- myCh
	val := <-myCh
	fmt.Println("Passed channel value:", val)
}
