package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

func printNumbers(wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 1; i <= 5; i++ {
		time.Sleep(10 * time.Millisecond)
		fmt.Printf("goroutine: %d\n", i)
	}
}

func sayHello() {
	fmt.Println("Hello from a goroutine")
}

func main() {
	fmt.Println("=== Goroutine basics ===")
	go sayHello()
	time.Sleep(50 * time.Millisecond)

	fmt.Println("\n=== Goroutine scheduling ===")
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 3; i++ {
			runtime.Gosched()
			fmt.Println("lambda goroutine")
		}
	}()

	go printNumbers(&wg)
	wg.Wait()

	fmt.Println("\n=== GOMAXPROCS ===")
	fmt.Printf("Logical CPUs: %d\n", runtime.NumCPU())
	fmt.Printf("GOMAXPROCS: %d\n", runtime.GOMAXPROCS(0))

	fmt.Println("\n=== Goroutine stack growth ===")
	stackTester(0)
}

func stackTester(depth int) {
	if depth >= 1000 {
		return
	}
	if depth%100 == 0 {
		fmt.Printf("stack depth: %d\n", depth)
	}
	stackTester(depth + 1)
}
