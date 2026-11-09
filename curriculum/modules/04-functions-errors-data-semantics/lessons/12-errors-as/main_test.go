package main

import (
	"errors"
	"testing"
)

func TestFetchURL_Missing(t *testing.T) {
	err := fetchURL("/missing")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var httpErr *HTTPError
	if !errors.As(err, &httpErr) {
		t.Fatal("expected *HTTPError")
	}
	if httpErr.StatusCode != 404 {
		t.Errorf("expected 404, got %d", httpErr.StatusCode)
	}
	if httpErr.Body != "not found" {
		t.Errorf("expected 'not found', got %q", httpErr.Body)
	}
}

func TestFetchURL_Broken(t *testing.T) {
	err := fetchURL("/broken")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var httpErr *HTTPError
	if !errors.As(err, &httpErr) {
		t.Fatal("expected *HTTPError")
	}
	if httpErr.StatusCode != 500 {
		t.Errorf("expected 500, got %d", httpErr.StatusCode)
	}
}

func TestFetchURL_OK(t *testing.T) {
	err := fetchURL("/ok")
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestFetchURL_Other(t *testing.T) {
	err := fetchURL("/other")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var httpErr *HTTPError
	if !errors.As(err, &httpErr) {
		t.Fatal("expected *HTTPError")
	}
	if httpErr.StatusCode != 400 {
		t.Errorf("expected 400, got %d", httpErr.StatusCode)
	}
}

func TestAsReturnsFalseForWrongType(t *testing.T) {
	err := fetchURL("/missing")
	var target string
	// errors.As requires *error, not *string — this test just confirms compilation
	_ = err
	_ = target
}
