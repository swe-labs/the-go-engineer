// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

package main

// ============================================================================
// Stage 07: Concurrency — Concurrent Downloader (Exercise)
// Level: Advanced
// ============================================================================
//
// RUN: go run ./07-concurrency/01-concurrency/goroutines/7-downloader
// ============================================================================

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"
)

func DownloadFile(url, destDir string) error {
	filename := filepath.Base(url)
	filePath := filepath.Join(destDir, filename)

	out, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer out.Close()

	fmt.Println("Downloading", url)
	start := time.Now()

	resp, err := http.Get(url)
	if err != nil {
		_ = os.Remove(filePath)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		_ = os.Remove(filePath)
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	fmt.Printf("Download %s took %s\n", filename, time.Since(start))
	return nil
}

func SequentialDownloader(urls []string, destDir string) error {
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return err
	}

	start := time.Now()
	for _, url := range urls {
		if err := DownloadFile(url, destDir); err != nil {
			fmt.Println("Error downloading", url, err)
			continue
		}
	}

	fmt.Printf("All sequential downloads (%d files) completed in %s\n", len(urls), time.Since(start))
	return nil
}

type Result struct {
	URL      string
	Filename string
	Size     int64
	Duration time.Duration
	Error    error
}

func ConcurrentDownloader(urls []string, destDir string, maxConcurrent int) error {
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return err
	}

	// 1. Result Aggregator Channel
	// Unbuffered channel where all workers will send their final Success/Error structs.
	results := make(chan Result)

	// 2. The WaitGroup
	// Acts as an atomic counter to track exactly how many Goroutines are currently running.
	var wg sync.WaitGroup

	// 3. The "Semaphore" Pattern
	// A buffered channel initialized with `maxConcurrent` capacity.
	// If the buffer is full (e.g., length is 3), the 4th Goroutine will BLOCK
	// when trying to write to it, essentially throttling the OS file decriptions!
	limiter := make(chan struct{}, maxConcurrent)

	for _, url := range urls {
		wg.Add(1) // Increment active worker count

		// 4. Closure Variable Capture
		// We MUST pass `url` as a function argument. If we just read the `url`
		// variable from the outer loop, the inner Goroutines would race and often
		// end up downloading the final URL multiple times.
		go func(url string) {
			defer wg.Done() // Decrement counter when function exits

			limiter <- struct{}{}        // Acquire Semaphore token (blocks if full)
			defer func() { <-limiter }() // Release Semaphore token back to the buffer

			start := time.Now()
			filename := filepath.Base(url)
			filePath := filepath.Join(destDir, filename)

			out, err := os.Create(filePath)
			if err != nil {
				results <- Result{URL: url, Error: err}
				return
			}
			defer out.Close()

			resp, err := http.Get(url)
			if err != nil {
				results <- Result{URL: url, Error: err}
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				results <- Result{URL: url, Error: fmt.Errorf("bad status: %s", resp.Status)}
				return
			}
			size, err := io.Copy(out, resp.Body)
			if err != nil {
				results <- Result{URL: url, Error: err}
				return
			}
			timeSince := time.Since(start)
			results <- Result{URL: url, Filename: filename, Size: size, Duration: timeSince, Error: nil}
		}(url)
	}

	// 5. The Background Waiter
	// We MUST run wg.Wait() in a separate Goroutine!
	// If we ran `wg.Wait()` directly here in the main function, the main thread
	// would block forever, and wouldn't be able to reach the `for result := range results`
	// loop below to consume the data. The workers would then get deadlocked trying
	// to send to the unbuffered `results` channel!
	go func() {
		wg.Wait()
		close(results) // Close channel, signaling the range loop to terminate
	}()

	var totalSize int64
	var errors []error
	start := time.Now()

	for result := range results {
		if result.Error != nil {
			fmt.Printf("Error downloading %s: %s\n", result.URL, result.Error.Error())
			errors = append(errors, result.Error)
		} else {
			totalSize += result.Size
			fmt.Printf("Downloaded %s (%d bytes) in %s\n", result.Filename, result.Size, result.Duration)
		}
	}

	startedSince := time.Since(start)
	fmt.Printf("All downloads completed in %s, Total: %d bytes\n", startedSince, totalSize)
	if len(errors) > 0 {
		return fmt.Errorf("errors downloading: %+v", errors)
	}

	return nil
}
func main() {

	urls := []string{"https://file-examples.com/storage/fe062d525c682cad199a167/2017/10/file_example_JPG_1MB.jpg",
		"https://go.dev/images/go-logo-white.svg"}

	err := ConcurrentDownloader(urls, "./", 3)
	if err != nil {
		fmt.Println(err)
		return
	}

	log.Println("Done")

}
