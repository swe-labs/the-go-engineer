package main

import (
	"fmt"
	"io"
	"os"
)

// FormatGreeting constructs a standard Go welcome message.
func FormatGreeting(target string) string {
	if target == "" {
		target = "World"
	}
	return fmt.Sprintf("Hello, %s!", target)
}

// WriteGreeting sends formatted output to an io.Writer (e.g., os.Stdout).
func WriteGreeting(w io.Writer, target string) error {
	_, err := fmt.Fprintln(w, FormatGreeting(target))
	return err
}

func main() {
	// Demonstrating output to standard output stream (os.Stdout)
	if err := WriteGreeting(os.Stdout, "World"); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing greeting: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Program completed successfully via package main -> func main() entry point.")
}
