package main

import (
	"fmt"
	"strings"
)

// Section 08: Quality & Testing - Interfaces for testability
//
// NEXT UP: TE.8

func te_7Summary(input string) string {
	return strings.TrimSpace(strings.ToLower(input))
}

func main() {
	fmt.Println("summary:", te_7Summary("  Interface Seams  "))
}
