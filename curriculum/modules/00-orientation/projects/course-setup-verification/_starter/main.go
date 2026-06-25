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

	fmt.Println(checkTool("Go", "go"))
	fmt.Println(checkTool("Git", "git"))

	// TODO: Check for a text editor (e.g., "code" for VS Code, "vim", "nano", etc.)
	editor := "code"
	fmt.Println(checkTool("Editor (code)", editor))
}
