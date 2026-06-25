package main

import (
	"errors"
	"testing"
)

func TestFetchData_EmptySource(t *testing.T) {
	err := fetchData("")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
	if unwrapped := errors.Unwrap(err); unwrapped == nil {
		t.Error("expected non-nil unwrapped error")
	}
}

func TestFetchData_BadHost(t *testing.T) {
	err := fetchData("badhost")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if errors.Is(err, ErrNotFound) {
		t.Errorf("did not expect ErrNotFound, got %v", err)
	}
}

func TestFetchData_Valid(t *testing.T) {
	err := fetchData("valid")
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestFetchData_Unknown(t *testing.T) {
	err := fetchData("other")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestWrapWithVDoesNotWrap(t *testing.T) {
	base := errors.New("base")
	wrapped := fmtErrorfWithV("context: %v", base)
	if errors.Unwrap(wrapped) != nil {
		t.Error("expected nil unwrap for percent v formatted error")
	}
}

// fmtErrorfWithV simulates what happens with %v (no wrap chain preserved).
func fmtErrorfWithV(format string, args ...interface{}) error {
	return errors.New(format) // simplified: actual %v behavior
}
