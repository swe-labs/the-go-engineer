package main

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestWorkerPool(t *testing.T) {
	tests := []struct {
		name       string
		numJobs    int
		numWorkers int
		wantOutput int
	}{
		{"3 workers 5 jobs", 5, 3, 5},
		{"1 worker 10 jobs", 10, 1, 10},
		{"5 workers 0 jobs", 0, 5, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jobs := make(chan Job, tt.numJobs)
			results := make(chan Result, tt.numJobs)

			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()

			var wg sync.WaitGroup
			for w := 1; w <= tt.numWorkers; w++ {
				wg.Add(1)
				go func(id int) {
					defer wg.Done()
					worker(ctx, id, jobs, results)
				}(w)
			}

			for j := 1; j <= tt.numJobs; j++ {
				jobs <- Job{ID: j, Delay: time.Millisecond}
			}
			close(jobs)

			go func() {
				wg.Wait()
				close(results)
			}()

			count := 0
			for range results {
				count++
			}
			if count != tt.wantOutput {
				t.Errorf("got %d results, want %d", count, tt.wantOutput)
			}
		})
	}
}
