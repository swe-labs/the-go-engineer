package main

import (
	"bufio"
	"fmt"
	"os"
)

// Section 00: How Computers Work

func main() {
	// 1. Write to stdout asking for input
	fmt.Print("Enter your name (data for stdin): ")

	// 2. Open the stdin pipe and listen for data
	reader := bufio.NewReader(os.Stdin)

	// Read until the user presses Enter (newline character)
	input, _ := reader.ReadString('\n')

	// 3. Write the captured data back out to stdout
	fmt.Printf("Hello, %s! Your data flowed through the pipes successfully.\n", input)

	// 4. Write to stderr (usually appears on the same screen, but is technically a separate pipe)
	fmt.Fprintln(os.Stderr, "(This message was secretly sent through the stderr pipe!)")

	fmt.Println("NEXT UP: HC.5 os-processes")
}
