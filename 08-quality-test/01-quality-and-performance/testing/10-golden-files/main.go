package main

import (
	"fmt"
	"strings"
)

// Section 08: Quality & Testing - Golden files
//
// NEXT UP:

func te_10Summary(input string) string {
	return strings.TrimSpace(strings.ToLower(input))
}

func main() {
	fmt.Println("summary:", te_10Summary("  Golden File Expectations  "))
}
