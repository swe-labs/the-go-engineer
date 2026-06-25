package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	info, _ := os.Stdin.Stat()
	isPipe := (info.Mode() & os.ModeCharDevice) == 0

	if isPipe {
		fmt.Fprintln(os.Stderr, "Reading piped input...")
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			line := scanner.Text()
			fmt.Println(strings.ToUpper(line))
		}
		if err := scanner.Err(); err != nil {
			fmt.Fprintf(os.Stderr, "error reading stdin: %v\n", err)
			os.Exit(1)
		}
	} else {
		fmt.Println("Interactive mode. Type lines (Ctrl+C to exit):")
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			line := scanner.Text()
			if line == "exit" || line == "quit" {
				break
			}
			fmt.Println(strings.ToUpper(line))
		}
	}
}
