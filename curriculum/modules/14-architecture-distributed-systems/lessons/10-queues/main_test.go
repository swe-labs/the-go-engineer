package main

import "testing"

func TestQueueFIFO(t *testing.T) {
	q := NewQueue()
	for i := 0; i < 5; i++ {
		q.Enqueue(i)
	}
	for i := 0; i < 5; i++ {
		got, err := q.Dequeue()
		if err != nil {
			t.Fatalf("unexpected error on dequeue %d: %v", i, err)
		}
		if got != i {
			t.Errorf("expected %d, got %d", i, got)
		}
	}
}

func TestQueueEmptyDequeue(t *testing.T) {
	q := NewQueue()
	_, err := q.Dequeue()
	if err == nil {
		t.Fatal("expected error on empty dequeue")
	}
}

func TestQueueLen(t *testing.T) {
	q := NewQueue()
	if q.Len() != 0 {
		t.Errorf("expected len 0, got %d", q.Len())
	}
	q.Enqueue(10)
	if q.Len() != 1 {
		t.Errorf("expected len 1, got %d", q.Len())
	}
	q.Dequeue()
	if q.Len() != 0 {
		t.Errorf("expected len 0, got %d", q.Len())
	}
}

func TestQueueCloseRejectsEnqueue(t *testing.T) {
	q := NewQueue()
	q.Close()
	err := q.Enqueue(1)
	if err == nil {
		t.Fatal("expected error on enqueue after close")
	}
}

func TestWorkerPoolProcessesAllItems(t *testing.T) {
	q := NewQueue()
	n := 100
	for i := 0; i < n; i++ {
		q.Enqueue(i)
	}

	results := make([]int, 0, n)
	var mu testing.TB = t

	pool := NewWorkerPool(q, 4)
	pool.Start(func(item int) {
		mu.Logf("item %d", item)
	})
	q.Close()
	pool.Wait()
	if pool.Processed() != int64(n) {
		t.Errorf("expected %d processed, got %d", n, pool.Processed())
	}
	_ = results
}

func TestQueueExactOrderPreserved(t *testing.T) {
	tests := []struct {
		name     string
		items    []int
		enqueues int
	}{
		{"sequential", []int{1, 2, 3, 4, 5}, 5},
		{"single_item", []int{42}, 1},
		{"ten_items", []int{0, 10, 20, 30, 40, 50, 60, 70, 80, 90}, 10},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			q := NewQueue()
			for _, v := range tc.items {
				q.Enqueue(v)
			}
			for i, expected := range tc.items {
				got, err := q.Dequeue()
				if err != nil {
					t.Fatalf("step %d: unexpected error: %v", i, err)
				}
				if got != expected {
					t.Errorf("step %d: expected %d, got %d", i, expected, got)
				}
			}
		})
	}
}
