//go:build !race

package main

import "testing"

func TestUnsafeIncrementRaces(t *testing.T) {
	c := &Counter{}
	UnsafeIncrement(c)
	t.Logf("unsafe counter value: %d (non-deterministic)", c.value)
}
