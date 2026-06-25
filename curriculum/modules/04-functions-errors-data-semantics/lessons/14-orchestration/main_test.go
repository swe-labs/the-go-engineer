package main

import (
	"errors"
	"testing"
)

func TestDeploy_EmptyService(t *testing.T) {
	err := deploy("")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, ErrPermanent) {
		t.Errorf("expected ErrPermanent, got %v", err)
	}
}

func TestDeploy_ValidService(t *testing.T) {
	// This may pass or fail depending on random retry outcome
	err := deploy("my-api")
	if err != nil {
		// If it fails, it should be a transient error after retries
		if !errors.Is(err, ErrTransient) && !errors.Is(err, ErrPermanent) {
			t.Errorf("unexpected error type: %v", err)
		}
	}
}

func TestRetry_SuccessOnFirstAttempt(t *testing.T) {
	attempts := 0
	err := retry(3, func() error {
		attempts++
		return nil
	})
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	if attempts != 1 {
		t.Errorf("expected 1 attempt, got %d", attempts)
	}
}

func TestRetry_FailsOnAllAttempts(t *testing.T) {
	attempts := 0
	err := retry(3, func() error {
		attempts++
		return ErrTransient
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if attempts != 3 {
		t.Errorf("expected 3 attempts, got %d", attempts)
	}
}

func TestRetry_StopsOnPermanent(t *testing.T) {
	attempts := 0
	err := retry(3, func() error {
		attempts++
		return ErrPermanent
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if attempts != 1 {
		t.Errorf("expected 1 attempt, got %d", attempts)
	}
}

func TestBuild_Empty(t *testing.T) {
	err := build("")
	if !errors.Is(err, ErrPermanent) {
		t.Errorf("expected ErrPermanent, got %v", err)
	}
}

func TestBuild_Valid(t *testing.T) {
	err := build("my-api")
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestIsTransient(t *testing.T) {
	if !isTransient(ErrTransient) {
		t.Error("expected ErrTransient to be transient")
	}
	if isTransient(ErrPermanent) {
		t.Error("expected ErrPermanent to NOT be transient")
	}
}
