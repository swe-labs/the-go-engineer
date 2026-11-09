package main

import (
	"fmt"
	"sync"
)

func main() {
	fmt.Println("=== Buffered channel basics ===")
	ch := make(chan int, 3)

	ch <- 1
	ch <- 2
	ch <- 3
	fmt.Println("Sent 3 values without blocking")

	fmt.Println(<-ch)
	fmt.Println(<-ch)
	fmt.Println(<-ch)

	fmt.Println("\n=== Channel as semaphore ===")
	sem := make(chan struct{}, 2)
	var wg sync.WaitGroup

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			sem <- struct{}{}
			fmt.Printf("Worker %d acquired semaphore\n", id)
			<-sem
			fmt.Printf("Worker %d released semaphore\n", id)
		}(i)
	}

	wg.Wait()

	fmt.Println("\n=== Non-blocking send with default ===")
	buf := make(chan int, 2)

	for i := 0; i < 4; i++ {
		select {
		case buf <- i:
			fmt.Printf("Sent %d\n", i)
		default:
			fmt.Printf("Buffer full, dropped %d\n", i)
		}
	}

	fmt.Println("\n=== Buffer capacity and length ===")
	c := make(chan string, 5)
	fmt.Printf("cap=%d len=%d\n", cap(c), len(c))

	c <- "a"
	c <- "b"
	fmt.Printf("cap=%d len=%d\n", cap(c), len(c))

	<-c
	fmt.Printf("cap=%d len=%d\n", cap(c), len(c))
}
