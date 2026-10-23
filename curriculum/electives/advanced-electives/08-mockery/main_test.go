package main

import "testing"

func TestDemo(t *testing.T) {
    got := Demo()
    want := "Mockery verified"
    if got != want {
        t.Errorf("got %q, want %q", got, want)
    }
}
