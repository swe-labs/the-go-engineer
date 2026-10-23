package main

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"
)

func TestConcurrentDownloader_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("payload for " + r.URL.Path))
	}))
	defer server.Close()

	downloader := NewConcurrentDownloader(DownloaderConfig{
		Workers: 2,
		Client:  server.Client(),
	})

	var buf1, buf2 bytes.Buffer
	jobs := []DownloadJob{
		{ID: "j1", URL: server.URL + "/file1", Destination: &buf1},
		{ID: "j2", URL: server.URL + "/file2", Destination: &buf2},
	}

	results := downloader.DownloadAll(t.Context(), jobs)

	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}

	for _, res := range results {
		if res.Error != nil {
			t.Fatalf("job %s failed: %v", res.ID, res.Error)
		}
		if res.BytesWritten <= 0 {
			t.Fatalf("job %s wrote 0 bytes", res.ID)
		}
	}

	if buf1.String() != "payload for /file1" {
		t.Fatalf("unexpected content for buf1: %q", buf1.String())
	}
	if buf2.String() != "payload for /file2" {
		t.Fatalf("unexpected content for buf2: %q", buf2.String())
	}
}

func TestConcurrentDownloader_RetryOnTransientFailure(t *testing.T) {
	var attempts atomic.Int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		count := attempts.Add(1)
		if count == 1 {
			// First attempt fails with 503
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		// Subsequent attempt succeeds
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("recovered data"))
	}))
	defer server.Close()

	downloader := NewConcurrentDownloader(DownloaderConfig{
		Workers:        1,
		MaxRetries:     2,
		InitialBackoff: 10 * time.Millisecond,
		Client:         server.Client(),
	})

	var buf bytes.Buffer
	jobs := []DownloadJob{
		{ID: "retry-job", URL: server.URL, Destination: &buf},
	}

	results := downloader.DownloadAll(t.Context(), jobs)

	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Error != nil {
		t.Fatalf("expected successful recovery on retry, got %v", results[0].Error)
	}
	if results[0].Retries != 1 {
		t.Fatalf("expected exactly 1 retry, got %d", results[0].Retries)
	}
	if buf.String() != "recovered data" {
		t.Fatalf("expected recovered data, got %q", buf.String())
	}
}

func TestConcurrentDownloader_MaxRetriesExceeded(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	downloader := NewConcurrentDownloader(DownloaderConfig{
		Workers:        1,
		MaxRetries:     2,
		InitialBackoff: 5 * time.Millisecond,
		Client:         server.Client(),
	})

	var buf bytes.Buffer
	jobs := []DownloadJob{
		{ID: "failing-job", URL: server.URL, Destination: &buf},
	}

	results := downloader.DownloadAll(t.Context(), jobs)

	if len(results) != 1 || results[0].Error == nil {
		t.Fatal("expected error when max retries exceeded, got nil")
	}
	if !errors.Is(results[0].Error, ErrMaxRetriesExceeded) {
		t.Fatalf("expected ErrMaxRetriesExceeded, got %v", results[0].Error)
	}
	if results[0].Retries != 2 {
		t.Fatalf("expected 2 retries attempted, got %d", results[0].Retries)
	}
}

func TestConcurrentDownloader_ContextCancellation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	downloader := NewConcurrentDownloader(DownloaderConfig{
		Workers: 1,
		Client:  server.Client(),
	})

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()

	var buf bytes.Buffer
	jobs := []DownloadJob{
		{ID: "cancel-job", URL: server.URL, Destination: &buf},
	}

	results := downloader.DownloadAll(ctx, jobs)

	if len(results) != 1 || results[0].Error == nil {
		t.Fatal("expected error due to context cancellation, got nil")
	}
}
