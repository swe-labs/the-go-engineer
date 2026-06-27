package main

import "fmt"

func Stage(name string, working, staging []string) ([]string, error) {
	found := false
	for _, f := range working {
		if f == name {
			found = true
		}
	}
	if !found {
		return staging, fmt.Errorf("file not found: %s", name)
	}
	for _, f := range staging {
		if f == name {
			return staging, nil
		}
	}
	return append(staging, name), nil
}

func Unstage(name string, staging []string) []string {
	result := []string{}
	for _, f := range staging {
		if f != name {
			result = append(result, f)
		}
	}
	return result
}

func Commit(msg string, staging []string, history []string) ([]string, error) {
	if len(staging) == 0 {
		return history, fmt.Errorf("nothing staged")
	}
	return append(history, msg), nil
}

func main() {
	working := []string{"hello.txt"}
	staging := []string{}

	fmt.Println("Staging: (empty)")

	staging, err := Stage("hello.txt", working, staging)
	if err != nil {
		fmt.Println("Error:", err)
	}
	fmt.Println("After stage:", staging)

	history := []string{}
	history, err = Commit("feat: add hello.txt", staging, history)
	if err != nil {
		fmt.Println("Error:", err)
	}
	fmt.Printf("History: %d commits\n", len(history))

	working = []string{"hello.txt", "main.go"}
	staging, _ = Stage("main.go", working, staging)
	history, _ = Commit("feat: add main.go", staging, history)
	fmt.Printf("History: %d commits\n", len(history))
}
