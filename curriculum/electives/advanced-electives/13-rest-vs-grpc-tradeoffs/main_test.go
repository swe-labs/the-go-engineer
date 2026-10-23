package main

import "testing"

func TestDemo(t *testing.T) {
	got := Demo()
	want := "REST vs gRPC tradeoffs verified"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
