package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

var (
	ErrMaxRetriesExceeded = errors.New("max retries exceeded")
	ErrDownloadFailed     = errors.New("download failed")
)

// DownloadJob defines a target URL and buffer/writer destination.
type DownloadJob struct {
	ID          string
	URL         string
	Destination io.Writer
}

// DownloadResult summarizes the outcome of a download task.
type DownloadResult struct {
	ID           string
	URL          string
	BytesWritten int64
	Duration     time.Duration
	Retries      int
	Error        error
}

// DownloaderConfig holds concurrency and retry parameters.
type DownloaderConfig struct {
	Workers        int
	MaxRetries     int
	InitialBackoff time.Duration
	Client         *http.Client
}

// ConcurrentDownloader manages parallel downloads with retries.
type ConcurrentDownloader struct {
	workers        int
	maxRetries     int
	initialBackoff time.Duration
	client         *http.Client
	totalBytes     atomic.Int64
}

// NewConcurrentDownloader creates a new Downloader.
func NewConcurrentDownloader(cfg DownloaderConfig) *ConcurrentDownloader {
	if cfg.Workers <= 0 {
		cfg.Workers = 4
	}
	if cfg.MaxRetries < 0 {
		cfg.MaxRetries = 2
	}
	if cfg.InitialBackoff <= 0 {
		cfg.InitialBackoff = 10 * time.Millisecond
	}
	if cfg.Client == nil {
		cfg.Client = &http.Client{Timeout: 10 * time.Second}
	}

	return &ConcurrentDownloader{
		workers:        cfg.Workers,
		maxRetries:     cfg.MaxRetries,
		initialBackoff: cfg.InitialBackoff,
		client:         cfg.Client,
	}
}

// DownloadAll executes a batch of download jobs concurrently across bounded workers.
func (d *ConcurrentDownloader) DownloadAll(ctx context.Context, jobs []DownloadJob) []DownloadResult {
	if len(jobs) == 0 {
		return nil
	}

	results := make([]DownloadResult, len(jobs))
	jobChan := make(chan struct {
		index int
		job   DownloadJob
	}, len(jobs))

	for idx, j := range jobs {
		jobChan <- struct {
			index int
			job   DownloadJob
		}{index: idx, job: j}
	}
	close(jobChan)

	var wg sync.WaitGroup
	workers := d.workers
	if workers > len(jobs) {
		workers = len(jobs)
	}

	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for item := range jobChan {
				select {
				case <-ctx.Done():
					results[item.index] = DownloadResult{
						ID:    item.job.ID,
						URL:   item.job.URL,
						Error: ctx.Err(),
					}
				default:
					start := time.Now()
					written, retries, err := d.downloadWithRetry(ctx, item.job)
					results[item.index] = DownloadResult{
						ID:           item.job.ID,
						URL:          item.job.URL,
						BytesWritten: written,
						Duration:     time.Since(start),
						Retries:      retries,
						Error:        err,
					}
				}
			}
		}()
	}

	wg.Wait()
	return results
}

// downloadWithRetry executes a single download with exponential backoff.
func (d *ConcurrentDownloader) downloadWithRetry(ctx context.Context, job DownloadJob) (int64, int, error) {
	backoff := d.initialBackoff
	retries := 0

	for {
		written, err := d.downloadSingle(ctx, job)
		if err == nil {
			return written, retries, nil
		}

		// Non-retryable errors
		if errors.Is(ctx.Err(), context.Canceled) || errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return 0, retries, ctx.Err()
		}

		if retries >= d.maxRetries {
			return 0, retries, fmt.Errorf("%w: %v", ErrMaxRetriesExceeded, err)
		}

		retries++
		select {
		case <-ctx.Done():
			return 0, retries, ctx.Err()
		case <-time.After(backoff):
			backoff *= 2
		}
	}
}

func (d *ConcurrentDownloader) downloadSingle(ctx context.Context, job DownloadJob) (int64, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, job.URL, nil)
	if err != nil {
		return 0, err
	}

	resp, err := d.client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return 0, fmt.Errorf("%w: unexpected HTTP status %d", ErrDownloadFailed, resp.StatusCode)
	}

	written, err := io.Copy(job.Destination, resp.Body)
	if err != nil {
		return 0, err
	}

	d.totalBytes.Add(written)
	return written, nil
}

// TotalBytesDownloaded returns aggregate downloaded bytes across all jobs.
func (d *ConcurrentDownloader) TotalBytesDownloaded() int64 {
	return d.totalBytes.Load()
}

func main() {
	downloader := NewConcurrentDownloader(DownloaderConfig{
		Workers:        2,
		MaxRetries:     2,
		InitialBackoff: 50 * time.Millisecond,
	})

	var buf1, buf2 bytes.Buffer
	jobs := []DownloadJob{
		{ID: "job-1", URL: "https://go.dev/robots.txt", Destination: &buf1},
		{ID: "job-2", URL: "https://golang.org/robots.txt", Destination: &buf2},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	results := downloader.DownloadAll(ctx, jobs)
	for _, res := range results {
		slog.Info("download completed",
			"id", res.ID,
			"bytes", res.BytesWritten,
			"retries", res.Retries,
			"error", res.Error,
		)
	}
}
