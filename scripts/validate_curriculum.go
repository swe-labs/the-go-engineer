//go:build ignore

package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	args := []string{"run", "./internal/tools/curriculum"}
	if len(os.Args) > 1 {
		args = append(args, os.Args[1:]...)
	} else {
		args = append(args, "validate-all")
	}

	cmd := exec.Command("go", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Curriculum validation failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Curriculum validation passed successfully.")
}
