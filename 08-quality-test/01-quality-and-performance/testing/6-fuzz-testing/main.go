package main

import (
	"fmt"
	"strings"
)

// Section 08: Quality & Testing - Fuzz testing
//
// NEXT UP: TE.7

func te_6Summary(input string) string {
	return strings.TrimSpace(strings.ToLower(input))
}

func main() {
	fmt.Println("summary:", te_6Summary("  Fuzz Targets  "))
}
