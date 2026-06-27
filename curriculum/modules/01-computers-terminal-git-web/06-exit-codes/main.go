package main

import (
	"fmt"
	"os"
)

// ExitCodeByFileStatus returns 0 for success (file exists) or 1 for failure.
func ExitCodeByFileStatus(fileExists bool) int {
	if fileExists {
		return 0
	}
	return 1
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: exit-codes <filename>")
		os.Exit(2)
	}
	_, err := os.Stat(os.Args[1])
	exists := err == nil
	code := ExitCodeByFileStatus(exists)
	fmt.Printf("exit code: %d\n", code)
	os.Exit(code)
}
