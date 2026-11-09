package main

import (
	"fmt"
	"runtime/debug"
)

func main() {
	fmt.Println("=== Debugging panics ===")

	// Nil pointer dereference demo with recovery.
	fmt.Println("\n1. Nil pointer dereference:")
	SafeConfigRead()

	// Index out of range demo.
	fmt.Println("\n2. Index out of range:")
	SafeSliceAccess([]int{1, 2, 3}, 5)

	// Division by zero with recovery.
	fmt.Println("\n3. Division by zero:")
	SafeDivideDemo()

	// Panic in goroutine — cannot recover from outside.
	fmt.Println("\n4. Panic in goroutine (will crash):")
	// We skip this in demo to avoid crashing.
	fmt.Println("   (skipped — would crash the program)")
}

type Server struct {
	DB *Database
}

type Database struct {
	DSN string
}

func SafeConfigRead() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from nil pointer:", r)
			fmt.Printf("Stack:\n%s\n", debug.Stack())
		}
	}()

	s := &Server{DB: nil}
	_ = s.DB.DSN // panic: nil pointer dereference
}

func SafeSliceAccess(nums []int, idx int) (int, bool) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Recovered from index panic: %v\n", r)
		}
	}()
	if idx < 0 {
		return 0, false
	}
	return nums[idx], true
}

func SafeDivideDemo() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from divide by zero:", r)
		}
	}()
	a := 10
	b := 0
	_ = a / b // panic: integer divide by zero
}

func SafeDivide(a, b []int) []int {
	result := make([]int, 0, len(a))
	for i := 0; i < len(a); i++ {
		func(idx int) {
			defer func() {
				if r := recover(); r != nil {
					result = append(result, 0)
				}
			}()
			divisor := 1
			if idx < len(b) {
				divisor = b[idx]
			}
			result = append(result, a[idx]/divisor)
		}(i)
	}
	return result
}
