// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// RUN: go run ./07-concurrency/02-concurrency-patterns/4-bounded-pipeline-exercise
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
// Stage 07: Concurrency Patterns - Exercise: Image Resizer Solution
// ============================================================================

var bufPool = sync.Pool{
	New: func() any {
		return bytes.NewBuffer(make([]byte, 0, 2*1024*1024))
	},
}

func main() {
	imageIDs := []string{"img1", "img2", "img3", "img4", "img5", "img6", "img7", "img8", "img9", "img10", "img11", "imgError", "img13"}

	fmt.Println("Starting batch job...")
	start := time.Now()

	g, ctx := errgroup.WithContext(context.Background())
	g.SetLimit(4)

	for _, id := range imageIDs {
		id := id
		g.Go(func() error {
			return processImage(ctx, id)
		})
	}

	if err := g.Wait(); err != nil {
		fmt.Printf("[FAIL] Batch job failed: %v\n", err)
	} else {
		fmt.Printf("[OK] Batch job completed successfully in %v\n", time.Since(start))
	}
}

func processImage(ctx context.Context, id string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	buf := bufPool.Get().(*bytes.Buffer)
	defer func() {
		buf.Reset()
		bufPool.Put(buf)
	}()

	buf.WriteString("simulated image data for " + id)
	time.Sleep(100 * time.Millisecond)

	if id == "imgError" {
		return fmt.Errorf("corrupt image data for %s", id)
	}

	log.Printf("Processed %s (buffer capacity: %d)", id, buf.Cap())
	return nil
}
