package main

import (
	"fmt"
	"go/format"
	"os"
)

// FormatSourceCode formats raw Go code bytes according to canonical gofmt rules.
func FormatSourceCode(src []byte) ([]byte, error) {
	return format.Source(src)
}

func main() {
	fmt.Println("=== Deconstructing 'gofmt' Engine (go/format) ===")

	// Deliberately poorly formatted Go source snippet
	unformattedSnippet := []byte(`package    main
import    "fmt"
func    Calculate(   x int,y int  )   int {
    sum:=x+y
        return sum
}`)

	fmt.Println("--- 1. Unformatted Source Input ---")
	fmt.Println(string(unformattedSnippet))

	formatted, err := FormatSourceCode(unformattedSnippet)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Formatting error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\n--- 2. Canonical gofmt Output ---")
	fmt.Println(string(formatted))

	fmt.Println("Insight: gofmt parses Go source into an Abstract Syntax Tree (AST)")
	fmt.Println("and re-emits it according to the language's official grammar specification.")
	fmt.Println("Result: Zero time wasted debating code formatting in code reviews.")
}
