package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Job struct {
	ID    int
	Delay time.Duration
}

type Result struct {
	JobID  int
	Output string
}

func worker(ctx context.Context, id int, jobs <-chan Job, results chan<- Result) {
	for {
		select {
		case <-ctx.Done():
			return
		case job, ok := <-jobs:
			if !ok {
				return
			}
			time.Sleep(job.Delay)
			results <- Result{JobID: job.ID, Output: fmt.Sprintf("worker %d processed job %d", id, job.ID)}
		}
	}
}

func main() {
	numJobs := 10
	numWorkers := 3

	jobs := make(chan Job, numJobs)
	results := make(chan Result, numJobs)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup
	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			worker(ctx, id, jobs, results)
		}(w)
	}

	for j := 1; j <= numJobs; j++ {
		jobs <- Job{ID: j, Delay: time.Duration(j%3) * 10 * time.Millisecond}
	}
	close(jobs)

	go func() {
		wg.Wait()
		close(results)
	}()

	for r := range results {
		fmt.Println(r.Output)
	}
}
