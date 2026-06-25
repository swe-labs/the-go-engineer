package main

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type Queue struct {
	mu     sync.Mutex
	items  []int
	closed bool
}

func NewQueue() *Queue {
	return &Queue{items: make([]int, 0)}
}

func (q *Queue) Enqueue(item int) error {
	q.mu.Lock()
	defer q.mu.Unlock()
	if q.closed {
		return errors.New("queue closed")
	}
	q.items = append(q.items, item)
	return nil
}

func (q *Queue) Dequeue() (int, error) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.items) == 0 {
		return 0, errors.New("queue empty")
	}
	item := q.items[0]
	q.items = q.items[1:]
	return item, nil
}

func (q *Queue) Len() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return len(q.items)
}

func (q *Queue) Close() {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.closed = true
}

type WorkerPool struct {
	queue      *Queue
	numWorkers int
	processed  atomic.Int64
	wg         sync.WaitGroup
}

func NewWorkerPool(queue *Queue, numWorkers int) *WorkerPool {
	return &WorkerPool{queue: queue, numWorkers: numWorkers}
}

func (wp *WorkerPool) Start(handler func(int)) {
	for i := range wp.numWorkers {
		wp.wg.Add(1)
		go func(id int) {
			defer wp.wg.Done()
			for {
				item, err := wp.queue.Dequeue()
				if err != nil {
					return
				}
				handler(item)
				wp.processed.Add(1)
				_ = id
			}
		}(i)
	}
}

func (wp *WorkerPool) Wait() {
	wp.wg.Wait()
}

func (wp *WorkerPool) Processed() int64 {
	return wp.processed.Load()
}

func produceInts(q *Queue, count int) {
	for i := range count {
		q.Enqueue(i)
	}
}

func main() {
	q := NewQueue()
	produceInts(q, 20)

	pool := NewWorkerPool(q, 4)
	double := func(n int) { fmt.Printf("processed: %d\n", n*2) }
	pool.Start(double)
	time.Sleep(100 * time.Millisecond)
	q.Close()
	pool.Wait()
	fmt.Printf("Total processed: %d\n", pool.Processed())
}
