package main

import (
	"fmt"
)

func main() {
	fmt.Println("=== Close and receive ===")
	ch := make(chan int, 3)
	ch <- 10
	ch <- 20
	ch <- 30
	close(ch)

	for i := 0; i < 5; i++ {
		v, ok := <-ch
		fmt.Printf("v=%d ok=%t\n", v, ok)
	}

	fmt.Println("\n=== Range over closed channel ===")
	ch2 := make(chan string, 2)
	ch2 <- "alpha"
	ch2 <- "beta"
	close(ch2)

	for s := range ch2 {
		fmt.Println(s)
	}

	fmt.Println("\n=== Close signals completion ===")
	done := make(chan struct{})
	go func() {
		fmt.Println("work complete")
		close(done)
	}()

	<-done
	fmt.Println("main notified")

	fmt.Println("\n=== Multiple receivers on close ===")
	ch3 := make(chan int, 5)
	for i := 1; i <= 5; i++ {
		ch3 <- i * 100
	}
	close(ch3)

	total := 0
	for v := range ch3 {
		total += v
	}
	fmt.Println("Sum:", total)
}
