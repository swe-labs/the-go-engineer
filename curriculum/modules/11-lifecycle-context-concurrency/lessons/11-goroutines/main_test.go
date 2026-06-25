package main

import (
	"testing"
)

func TestGoroutinesCompiles(t *testing.T) {
}

func TestGoroutineBasic(t *testing.T) {
	done := make(chan struct{}, 1)
	go func() {
		done <- struct{}{}
	}()
	<-done
}
