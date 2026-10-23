//go:build ignore

package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	cmd := exec.Command("go", "run", "./internal/tools/curriculum", "validate-all")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Curriculum validation failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Curriculum validation passed successfully.")
}
