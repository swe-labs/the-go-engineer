// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// RUN: go run ./24-errgroup-and-pools/4-bounded-pipeline-exercise
package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"
)

// ============================================================================
// Section 24: errgroup & sync.Pool — Exercise: Image Resizer SOLUTION
// ============================================================================

// 1. Create the sync.Pool
var bufPool = sync.Pool{
	New: func() any {
		// Pre-allocate 2MB capacity per buffer
		return bytes.NewBuffer(make([]byte, 0, 2*1024*1024))
	},
}

func main() {
	imageIDs := []string{"img1", "img2", "img3", "img4", "img5", "img6", "img7", "img8", "img9", "img10", "img11", "imgError", "img13"}

	fmt.Println("Starting batch job...")
	start := time.Now()

	// 2. Initialize errgroup.WithContext
	g, ctx := errgroup.WithContext(context.Background())

	// 3. Limit concurrency to 4
	g.SetLimit(4)

	for _, id := range imageIDs {
		id := id // Capture loop variable

		// 4. Launch the job via g.Go
		g.Go(func() error {
			return processImage(ctx, id)
		})
	}

	// 5. Wait for the group to finish
	if err := g.Wait(); err != nil {
		fmt.Printf("❌ Batch job failed: %v\n", err)
	} else {
		fmt.Printf("✅ Batch job completed successfully in %v\n", time.Since(start))
	}
}

func processImage(ctx context.Context, id string) error {
	// Respect cancellation
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	// 6. Get a buffer from the pool and ensure it returns
	buf := bufPool.Get().(*bytes.Buffer)
	defer func() {
		buf.Reset()
		bufPool.Put(buf)
	}()

	// Simulate heavy work
	buf.WriteString("simulated image data for " + id)
	time.Sleep(100 * time.Millisecond) // Simulate processing time

	if id == "imgError" {
		return fmt.Errorf("corrupt image data for %s", id)
	}

	log.Printf("Processed %s (buffer capacity: %d)", id, buf.Cap())
	return nil
}
