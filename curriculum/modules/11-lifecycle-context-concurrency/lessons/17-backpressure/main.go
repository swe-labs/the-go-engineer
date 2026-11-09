package main

import (
	"fmt"
	"sync"
	"time"
)

// Bounded channel as backpressure: producer blocks when buffer is full.
func boundedProducer(done <-chan struct{}, limit int) <-chan int {
	ch := make(chan int, limit)
	go func() {
		defer close(ch)
		for i := 1; ; i++ {
			select {
			case ch <- i:
				fmt.Printf("enqueued %d (len=%d)\n", i, len(ch))
			case <-done:
				return
			}
		}
	}()
	return ch
}

// Load shedding: drop items when buffer is full.
func loadSheddingProducer(limit int) <-chan int {
	ch := make(chan int, limit)
	go func() {
		defer close(ch)
		for i := 1; i <= 20; i++ {
			select {
			case ch <- i:
				fmt.Printf("accepted %d\n", i)
			default:
				fmt.Printf("REJECTED %d (buffer full)\n", i)
			}
			time.Sleep(5 * time.Millisecond)
		}
	}()
	return ch
}

// Monitor queue depth.
type MonitoredChannel struct {
	ch    chan int
	count int
	mu    sync.Mutex
}

func (m *MonitoredChannel) Send(v int) bool {
	select {
	case m.ch <- v:
		m.mu.Lock()
		m.count++
		m.mu.Unlock()
		return true
	default:
		fmt.Printf("queue depth: %d/%d\n", len(m.ch), cap(m.ch))
		return false
	}
}

func main() {
	fmt.Println("=== Backpressure with bounded channel ===")
	done := make(chan struct{})
	ch := boundedProducer(done, 3)

	for i := 0; i < 10; i++ {
		<-ch
		time.Sleep(20 * time.Millisecond)
	}
	close(done)
	fmt.Println("Main done consuming")

	fmt.Println("\n=== Load shedding ===")
	time.Sleep(50 * time.Millisecond)
	ch2 := loadSheddingProducer(3)

	for v := range ch2 {
		time.Sleep(10 * time.Millisecond)
		fmt.Printf("processed %d\n", v)
		if v >= 8 {
			break
		}
	}
}
