package main

import (
	"fmt"
	"sync"
	"time"
)

func deadlockExample() {
	var mu1, mu2 sync.Mutex
	var wg sync.WaitGroup

	wg.Add(2)
	go func() {
		defer wg.Done()
		mu1.Lock()
		time.Sleep(10 * time.Millisecond)
		mu2.Lock()
		mu2.Unlock()
		mu1.Unlock()
	}()

	go func() {
		defer wg.Done()
		mu2.Lock()
		time.Sleep(10 * time.Millisecond)
		mu1.Lock()
		mu1.Unlock()
		mu2.Unlock()
	}()
	wg.Wait()
}

func safeOrderedLock() {
	var mu1, mu2 sync.Mutex
	var wg sync.WaitGroup

	wg.Add(2)
	go func() {
		defer wg.Done()
		mu1.Lock()
		mu2.Lock()
		mu2.Unlock()
		mu1.Unlock()
	}()

	go func() {
		defer wg.Done()
		mu1.Lock()
		mu2.Lock()
		mu2.Unlock()
		mu1.Unlock()
	}()
	wg.Wait()
}

func main() {
	fmt.Println("Deadlocks: Coffman conditions")
	fmt.Println("Run tests to verify safe locking order")
}
