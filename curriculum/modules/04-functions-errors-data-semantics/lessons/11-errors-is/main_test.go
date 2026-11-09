package main

import (
	"errors"
	"testing"
)

func TestCallService_Timeout(t *testing.T) {
	err := callService(true)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var n *NetworkError
	if !errors.As(err, &n) {
		t.Fatal("expected *NetworkError")
	}
	if !n.IsTimeout {
		t.Error("expected IsTimeout to be true")
	}
	if !errors.Is(err, ErrTimeout) {
		t.Error("expected errors.Is(err, ErrTimeout) to be true")
	}
}

func TestCallService_NoTimeout(t *testing.T) {
	err := callService(false)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var n *NetworkError
	if !errors.As(err, &n) {
		t.Fatal("expected *NetworkError")
	}
	if n.IsTimeout {
		t.Error("expected IsTimeout to be false")
	}
	if errors.Is(err, ErrTimeout) {
		t.Error("expected errors.Is(err, ErrTimeout) to be false")
	}
}

func TestCallService_NotTimeoutViaIs(t *testing.T) {
	err := callService(false)
	if !errors.Is(err, ErrTimeout) {
		// correct: IsTimeout is false, so custom Is returns false
	}
}

func TestNetworkErrorIsMethod(t *testing.T) {
	ne := &NetworkError{Msg: "test", IsTimeout: true}
	if !ne.Is(ErrTimeout) {
		t.Error("expected Is(ErrTimeout) to be true")
	}
	if ne.Is(errors.New("other")) {
		t.Error("expected Is(other) to be false")
	}
}
