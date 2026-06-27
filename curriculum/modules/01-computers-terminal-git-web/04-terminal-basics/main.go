package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type CommandResult struct {
	Input    string
	Output   string
}

func simulateShell(input string) CommandResult {
	parts := strings.Fields(input)
	if len(parts) == 0 {
		return CommandResult{Input: input, Output: ""}
	}

	switch parts[0] {
	case "echo":
		return CommandResult{Input: input, Output: strings.Join(parts[1:], " ")}
	case "hello":
		return CommandResult{Input: input, Output: "Hello from the Go shell!"}
	case "exit":
		os.Exit(0)
		return CommandResult{}
	default:
		return CommandResult{Input: input, Output: "command not found: " + parts[0]}
	}
}

func main() {
	fmt.Println("=== Terminal Basics ===")
	fmt.Println("Type a command (echo, hello, exit):")
	fmt.Println()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		result := simulateShell(line)
		if result.Output != "" {
			fmt.Println(result.Output)
		}
		fmt.Print("> ")
	}

	fmt.Println("Terminal simulation ended.")
}
