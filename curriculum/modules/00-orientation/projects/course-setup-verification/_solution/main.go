package main

import (
	"fmt"
	"os/exec"
)

func checkTool(name string, path string) string {
	_, err := exec.LookPath(path)
	if err != nil {
		return fmt.Sprintf("NOT FOUND: %s is not installed or not in PATH", name)
	}
	return fmt.Sprintf("OK: %s is installed", name)
}

func main() {
	fmt.Println("Course Setup Verification")
	fmt.Println("=========================")
	fmt.Println()

	tools := []struct {
		name string
		path string
	}{
		{"Go", "go"},
		{"Git", "git"},
		{"VS Code", "code"},
		{"Vim", "vim"},
		{"Nano", "nano"},
	}

	for _, tool := range tools {
		fmt.Println(checkTool(tool.name, tool.path))
	}

	fmt.Println()
	fmt.Println("If any tool is NOT FOUND, install it before proceeding.")
	fmt.Println("At minimum, you need Go, Git, and a text editor.")
}
