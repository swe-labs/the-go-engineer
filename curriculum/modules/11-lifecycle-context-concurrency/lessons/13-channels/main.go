package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("=== Unbuffered channel basics ===")
	ch := make(chan int)

	go func() {
		ch <- 42
	}()

	val := <-ch
	fmt.Println("Received:", val)

	fmt.Println("\n=== Channel as synchronization ===")
	done := make(chan struct{})

	go func() {
		fmt.Println("Working...")
		time.Sleep(50 * time.Millisecond)
		close(done)
	}()

	<-done
	fmt.Println("Done signal received")

	fmt.Println("\n=== Two-way communication ===")
	req := make(chan string)
	resp := make(chan string)

	go func() {
		msg := <-req
		resp <- "echo: " + msg
	}()

	req <- "hello"
	reply := <-resp
	fmt.Println(reply)

	fmt.Println("\n=== Multiple sends ===")
	ch2 := make(chan int)

	go func() {
		for i := 0; i < 5; i++ {
			ch2 <- i * 10
		}
		close(ch2)
	}()

	for n := range ch2 {
		fmt.Println("Got:", n)
	}
}
