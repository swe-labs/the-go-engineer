package main

import (
	"fmt"
	"os"
	"strings"
)

func simulateShell(input string) string {
	parts := strings.Fields(input)
	if len(parts) == 0 {
		return ""
	}

	switch parts[0] {
	case "echo":
		return strings.Join(parts[1:], " ")
	case "hello":
		return "Hello from the Go shell!"
	case "exit":
		os.Exit(0)
		return ""
	default:
		return "command not found: " + parts[0]
	}
}

func main() {
	fmt.Println("=== Terminal Basics ===")
	fmt.Println("Type a command (echo, hello, exit):")
	fmt.Println()

	var input string
	for {
		fmt.Print("> ")
		_, err := fmt.Scanln(&input)
		if err != nil {
			break
		}
		output := simulateShell(input)
		if output != "" {
			fmt.Println(output)
		}
	}

	fmt.Println("Terminal simulation ended.")
}
