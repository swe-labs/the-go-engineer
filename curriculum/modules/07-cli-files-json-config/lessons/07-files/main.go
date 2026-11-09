package main

import (
	"fmt"
	"os"
)

func main() {
	// Write file
	data := []byte("Hello, file!\nLine 2\nLine 3\n")
	err := os.WriteFile("test_output.txt", data, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "write error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Written test_output.txt")

	// Read file
	content, err := os.ReadFile("test_output.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "read error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Read %d bytes:\n%s", len(content), content)

	// Open and read with io.Copy to stdout
	f, err := os.Open("test_output.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "open error: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()

	info, _ := f.Stat()
	fmt.Printf("File size via Stat: %d bytes\n", info.Size())

	// Clean up
	os.Remove("test_output.txt")
	fmt.Println("Cleaned up test_output.txt")
}
