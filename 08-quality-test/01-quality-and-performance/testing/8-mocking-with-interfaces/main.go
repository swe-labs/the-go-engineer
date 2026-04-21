package main

import (
	"fmt"
	"strings"
)

// Section 08: Quality & Testing - Mocking with interfaces
//
// NEXT UP: TE.9

func te_8Summary(input string) string {
	return strings.TrimSpace(strings.ToLower(input))
}

func main() {
	fmt.Println("summary:", te_8Summary("  Mocking Collaborators  "))
}
