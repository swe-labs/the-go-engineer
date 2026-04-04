// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"time"

	"golang.org/x/sync/errgroup"
)

// ============================================================================
// Section 24: errgroup & sync.Pool — Exercise: Image Resizer
// Level: Intermediate/Advanced
// ============================================================================
//
// THE TASK:
// You are building an image processing pipeline. You receive a list of image
// IDs that need to be fetched, resized, and saved.
//
// Because image buffers are large (e.g., 2MB each), allocating a new buffer
// for every image will crash the server under load or cause massive GC spikes.
// Because the server has limited CPU cores, you cannot process all 20 images
// simultaneously.
//
// REQUIREMENTS:
// 1. Process all image IDs concurrently using errgroup.
// 2. Limit the concurrency so NO MORE than 4 images are processed at the same time.
//    (HINT: Use g.SetLimit).
// 3. Stop early if ANY processing fails.
//    (HINT: Use errgroup.WithContext, and trigger a context cancellation).
// 4. Use a `sync.Pool` to recycle the `*bytes.Buffer` used for the payload.
//
// RUN: go run ./24-errgroup-and-pools/4-bounded-pipeline-exercise/_starter
// ============================================================================

// TODO 1: Create a sync.Pool for *bytes.Buffer. Ensure you allocate enough capacity (e.g. 2MB).
// var bufPool = sync.Pool{ ... }

func main() {
	imageIDs := []string{"img1", "img2", "img3", "img4", "img5", "img6", "img7", "img8", "img9", "img10", "img11", "imgError", "img13"}

	fmt.Println("Starting batch job...")
	start := time.Now()

	// TODO 2: Initialize an errgroup.WithContext
	// g, ctx := ...
	var g errgroup.Group

	// TODO 3: Limit the concurrency to 4
	// g.SetLimit(4)

	for _, id := range imageIDs {
		id := id // Capture loop variable

		// TODO 4: Launch the processing inside the errgroup
		// g.Go(func() error { ... return processImage(ctx, id) })
		_ = id
	}

	// TODO 5: Wait for the group to finish and check for errors
	if err := g.Wait(); err != nil {
		fmt.Printf("❌ Batch job failed: %v\n", err)
	} else {
		fmt.Printf("✅ Batch job completed successfully in %v\n", time.Since(start))
	}
}

// processImage simulates processing an image.
func processImage(ctx context.Context, id string) error {
	// First check if the context is cancelled before starting heavy work
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	// TODO 6: Get a buffer from the pool, use it, and defer putting it back
	// buf := bufPool.Get().(*bytes.Buffer)
	// defer ...
	// buf.Reset()
	var buf *bytes.Buffer = bytes.NewBuffer(make([]byte, 0, 2*1024*1024)) // REMOVE THIS LINE and use the pool

	// Simulate heavy work
	buf.WriteString("simulated image data for " + id)
	time.Sleep(100 * time.Millisecond) // Simulate processing time

	if id == "imgError" {
		return fmt.Errorf("corrupt image data for %s", id)
	}

	log.Printf("Processed %s (buffer capacity: %d)", id, buf.Cap())
	return nil
}
