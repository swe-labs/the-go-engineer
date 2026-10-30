package main

import (
	"fmt"
	"runtime/debug"
	"strings"
)

// SafeExecute runs an arbitrary function and intercepts any runtime panic,
// returning the recovered value and the top frames of the runtime stack trace.
func SafeExecute(fn func()) (recoveredVal any, stack string) {
	defer func() {
		if r := recover(); r != nil {
			recoveredVal = r
			rawStack := string(debug.Stack())
			lines := strings.Split(rawStack, "\n")
			// Extract top 6 lines of stack trace for concise pedagogical inspection
			if len(lines) > 6 {
				stack = strings.Join(lines[:6], "\n")
			} else {
				stack = rawStack
			}
		}
	}()

	fn()
	return nil, ""
}

// AccessElement attempts an index access that may trigger an out-of-bounds runtime panic.
func AccessElement(items []int, idx int) int {
	return items[idx]
}

func main() {
	fmt.Println("=== Deconstructing Go Runtime Panics & Stack Traces ===")

	numbers := []int{10, 20, 30}
	fmt.Printf("Slice initialized: %v (Length: %d)\n", numbers, len(numbers))

	// Safely trigger and capture an index out-of-bounds runtime panic
	fmt.Println("\nAttempting out-of-bounds access at index 5...")
	recoveredVal, stack := SafeExecute(func() {
		_ = AccessElement(numbers, 5)
	})

	if recoveredVal != nil {
		fmt.Printf("✓ Successfully intercepted runtime panic: %v\n\n", recoveredVal)
		fmt.Println("--- Anatomy of the Runtime Stack Trace ---")
		fmt.Println(stack)
		fmt.Println("\nInsight: The stack trace reads bottom-to-top (caller -> callee -> crash point).")
		fmt.Println("The top frame points directly to the function and line number where the CPU halted.")
	} else {
		fmt.Println("No panic occurred.")
	}
}
