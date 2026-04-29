package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestConcurrentDownloader(t *testing.T) {
	// Create a mock HTTP server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("dummy content"))
	}))
	defer ts.Close()

	// Create temp dir for downloads
	tmpDir, err := os.MkdirTemp("", "downloader-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Run downloader against mock server
	urls := []string{ts.URL + "/file1.txt", ts.URL + "/file2.txt"}
	err = ConcurrentDownloader(urls, tmpDir, 2)
	if err != nil {
		t.Fatalf("ConcurrentDownloader returned error: %v", err)
	}

	// Verify files were downloaded
	for _, f := range []string{"file1.txt", "file2.txt"} {
		info, err := os.Stat(filepath.Join(tmpDir, f))
		if err != nil {
			t.Errorf("Expected file %s to be downloaded, got error: %v", f, err)
		} else if info.Size() == 0 {
			t.Errorf("Expected file %s to have content, but size is 0", f)
		}
	}
}
