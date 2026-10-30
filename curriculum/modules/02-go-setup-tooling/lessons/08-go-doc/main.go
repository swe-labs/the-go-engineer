// Package main demonstrates the conventions and mechanics of Go doc comments.
//
// Go's documentation tool ('go doc') extracts comments directly from the Abstract
// Syntax Tree without requiring external annotations, HTML tags, or separate markdown files.
// Every exported symbol should have a comment beginning with its own name.
package main

import (
	"fmt"
	"strings"
)

// Calculator provides basic arithmetic operations with clear documentation contracts.
type Calculator struct {
	Base int
}

// Add computes and returns the sum of the base value and the given integer parameter n.
// It does not mutate the receiver.
func (c Calculator) Add(n int) int {
	return c.Base + n
}

// CleanDocString normalizes raw comment lines for terminal display.
func CleanDocString(comment string) string {
	lines := strings.Split(comment, "\n")
	var cleaned []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(strings.TrimPrefix(line, "//"))
		if trimmed != "" {
			cleaned = append(cleaned, trimmed)
		}
	}
	return strings.Join(cleaned, " ")
}

func main() {
	fmt.Println("=== Deconstructing 'go doc' Documentation Engine ===")

	calc := Calculator{Base: 10}
	result := calc.Add(5)
	fmt.Printf("Calculation: 10 + 5 = %d\n\n", result)

	sampleComment := `// Add computes and returns the sum of the base value.
// It does not mutate the receiver.`

	fmt.Println("--- Doc Comment Parsing Demonstration ---")
	fmt.Println("Raw Doc Comment in Source:")
	fmt.Println(sampleComment)
	fmt.Println("\nCleaned Doc Summary (as shown by 'go doc'):")
	fmt.Printf("  %s\n", CleanDocString(sampleComment))

	fmt.Println("\nInsight: In Go, documentation is code. Run 'go doc' on any package or symbol")
	fmt.Println("to read instantaneous, offline documentation formatted directly from source.")
}
