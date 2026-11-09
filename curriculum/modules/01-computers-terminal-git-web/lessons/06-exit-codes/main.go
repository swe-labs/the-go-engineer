package main

import (
	"fmt"
	"os"
)

// CheckFileStatus validates if a file exists and returns a standardized exit code.
func CheckFileStatus(path string) int {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "Error: file '%s' does not exist.\n", path)
			return 1 // General error: file missing
		}
		fmt.Fprintf(os.Stderr, "Error accessing file: %v\n", err)
		return 1
	}

	if info.IsDir() {
		fmt.Fprintf(os.Stderr, "Error: '%s' is a directory, not a file.\n", path)
		return 1
	}

	fmt.Printf("Success: file '%s' found (%d bytes).\n", path, info.Size())
	return 0 // Success
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: exit-codes <filename>")
		fmt.Println("Exiting with status code 2 (Invalid Arguments)")
		os.Exit(2)
	}

	target := os.Args[1]

	// Handle simulated panic demo
	if target == "--simulate-panic" {
		fmt.Println("Triggering runtime panic...")
		panic("critical unexpected failure") // Will exit with status 2
	}

	exitCode := CheckFileStatus(target)
	fmt.Printf("Process terminating with exit code: %d\n", exitCode)
	os.Exit(exitCode)
}
