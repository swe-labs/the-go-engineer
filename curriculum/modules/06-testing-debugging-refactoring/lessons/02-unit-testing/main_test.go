package main

import "testing"

func TestSplitHostPortValid(t *testing.T) {
	host, port, err := SplitHostPort("localhost:8080")
	if err != nil {
		t.Fatalf("SplitHostPort(\"localhost:8080\") unexpected error: %v", err)
	}
	if host != "localhost" {
		t.Errorf("host = %q; want %q", host, "localhost")
	}
	if port != 8080 {
		t.Errorf("port = %d; want %d", port, 8080)
	}
}

func TestSplitHostPortMissingPort(t *testing.T) {
	_, _, err := SplitHostPort("localhost")
	if err == nil {
		t.Fatal("expected error for missing port")
	}
}

func TestSplitHostPortInvalidPort(t *testing.T) {
	_, _, err := SplitHostPort("localhost:abc")
	if err == nil {
		t.Fatal("expected error for invalid port")
	}
}

func TestSplitHostPortPortZero(t *testing.T) {
	_, _, err := SplitHostPort("localhost:0")
	if err == nil {
		t.Fatal("expected error for port 0")
	}
}

func TestSplitHostPortPortTooHigh(t *testing.T) {
	_, _, err := SplitHostPort("localhost:70000")
	if err == nil {
		t.Fatal("expected error for port too high")
	}
}
