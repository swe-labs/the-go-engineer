package main

import (
	"fmt"
	"strings"
)

// Section 08: Quality & Testing - Integration tests
//
// NEXT UP: TE.10

func te_9Summary(input string) string {
	return strings.TrimSpace(strings.ToLower(input))
}

func main() {
	fmt.Println("summary:", te_9Summary("  Integration Boundaries  "))
}
