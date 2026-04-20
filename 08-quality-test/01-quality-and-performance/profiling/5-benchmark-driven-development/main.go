package main

import (
	"fmt"
	"strings"
)

// Section 08: Quality & Testing - Benchmark-driven development
//
// NEXT UP: PR.6

func pr_5Summary(input string) string {
	return strings.TrimSpace(strings.ToLower(input))
}

func main() {
	fmt.Println("summary:", pr_5Summary("  Benchmark Driven Development  "))
}
