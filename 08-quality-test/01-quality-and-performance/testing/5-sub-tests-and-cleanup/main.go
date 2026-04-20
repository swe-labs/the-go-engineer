package main

import (
	"fmt"
	"strings"
)

// Section 08: Quality & Testing - Sub-tests and t.Cleanup
//
// NEXT UP: TE.6

func te_5Summary(input string) string {
	return strings.TrimSpace(strings.ToLower(input))
}

func main() {
	fmt.Println("summary:", te_5Summary("  Subtests With Cleanup  "))
}
