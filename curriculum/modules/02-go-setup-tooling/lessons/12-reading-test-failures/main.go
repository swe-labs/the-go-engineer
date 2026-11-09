package main

import (
	"fmt"
)

// Add returns the sum of two integers.
func Add(a, b int) int {
	return a + b
}

// FormatFailureMessage formats an idiomatic Go test failure string adhering to
// the standard "got X; want Y" diagnostic pattern.
func FormatFailureMessage(callDescription string, got, want any) string {
	return fmt.Sprintf("%s = %v; want %v", callDescription, got, want)
}

func main() {
	fmt.Println("=== Deconstructing Go Test Failure Reports ===")

	a, b := 2, 3
	actual := Add(a, b)
	fmt.Printf("Add(%d, %d) returned: %d\n\n", a, b, actual)

	// Demonstrating the canonical Go failure report format
	simulatedInput := "Add(2, 2)"
	simulatedGot := 5
	simulatedWant := 4
	failureReport := FormatFailureMessage(simulatedInput, simulatedGot, simulatedWant)

	fmt.Println("--- Anatomy of an Idiomatic Test Failure Log ---")
	fmt.Println("--- FAIL: TestAdd (0.00s)")
	fmt.Printf("    calc_test.go:14: %s\n\n", failureReport)

	fmt.Println("Insight: Go testing discourages opaque assertion macros (like assert.Equal).")
	fmt.Println("Standard 'got != want' messages clearly display what the code did vs what was expected.")
}
