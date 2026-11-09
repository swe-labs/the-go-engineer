# Goroutines & Concurrency Model

In Go, a **goroutine** is a lightweight thread managed by the Go runtime scheduler (the *M:N scheduler*).

## Memory Footprint

Unlike operating system threads which typically allocate 1MB to 8MB of stack space up front, a goroutine begins with a tiny stack of **2KB to 8KB** (growing and shrinking dynamically on the heap as needed).

```go
package main

import (
	"fmt"
	"time"
)

func worker(id int, ch chan<- string) {
	ch <- fmt.Sprintf("worker %d completed work", id)
}

func main() {
	results := make(chan string, 3)
	for i := 1; i <= 3; i++ {
		go worker(i, results)
	}
	for i := 1; i <= 3; i++ {
		fmt.Println(<-results)
	}
}
```

## The M:P:N Runtime Architecture

- **M (Machine)**: An operating system kernel thread managed by the OS scheduler.
- **P (Processor)**: A logical context / resource token representing execution capability (`GOMAXPROCS`).
- **G (Goroutine)**: The user-space thread descriptor containing stack pointer, instruction pointer, and scheduling state.

## Work Stealing Scheduler

When a logical processor `P` finishes executing all runnable goroutines in its local run queue, it inspects other processors' queues and steals half of their runnable goroutines to ensure even core utilization.
