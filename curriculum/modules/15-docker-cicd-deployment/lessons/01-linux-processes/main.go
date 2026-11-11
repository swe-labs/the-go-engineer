package main

import (
	"fmt"
	"os/exec"
)

func runCommand(name string, args ...string) (string, int, error) {
	cmd := exec.Command(name, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return string(out), exitErr.ExitCode(), exitErr
		}
		return "", -1, err
	}
	return string(out), cmd.ProcessState.ExitCode(), nil
}

func main() {
	out, code, err := runCommand("go", "version")
	if err != nil {
		fmt.Printf("Error: %v (exit %d)\n", err, code)
		return
	}
	fmt.Printf("Output: %s", out)
	fmt.Printf("Exit code: %d\n", code)

	// An unknown go subcommand exits with status 2 on every platform.
	_, code, err = runCommand("go", "this-subcommand-does-not-exist")
	if err != nil {
		fmt.Printf("Expected failure, exit code: %d\n", code)
	}
}
