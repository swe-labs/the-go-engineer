//go:build integration

package main

import (
	"fmt"
	"testing"
	"time"
)

// This test will ONLY run if you execute:
// go test -v -tags=integration ./05-packages-io/01-modules-and-packages/4-build-tags
func TestSlowDatabaseIntegration(t *testing.T) {
	fmt.Println("Running a slow database interaction test...")
	time.Sleep(2 * time.Second)
	fmt.Println("Done!")
}
