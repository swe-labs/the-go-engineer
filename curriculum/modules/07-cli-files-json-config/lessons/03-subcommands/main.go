package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: mytool <command> [flags]")
		fmt.Println("Commands:")
		fmt.Println("  greet  Print a greeting")
		fmt.Println("  add    Add two numbers")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "greet":
		greetCmd := flag.NewFlagSet("greet", flag.ExitOnError)
		name := greetCmd.String("name", "World", "name to greet")
		greetCmd.Parse(os.Args[2:])
		fmt.Printf("Hello, %s!\n", *name)

	case "add":
		addCmd := flag.NewFlagSet("add", flag.ExitOnError)
		a := addCmd.Int("a", 0, "first number")
		b := addCmd.Int("b", 0, "second number")
		addCmd.Parse(os.Args[2:])
		fmt.Printf("%d + %d = %d\n", *a, *b, *a+*b)

	default:
		fmt.Printf("Unknown command: %q\n", os.Args[1])
		os.Exit(1)
	}
}
