package main

import (
	"fmt"
	"os"
	"os/exec"
)

// GetPID returns the current process ID.
func GetPID() int {
	return os.Getpid()
}

// GetPPID returns the parent process ID.
func GetPPID() int {
	return os.Getppid()
}

// RunCommand runs a command and returns its exit code.
// Returns -1 if the command could not be started.
func RunCommand(name string, args ...string) (int, error) {
	cmd := exec.Command(name, args...)
	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return exitErr.ExitCode(), nil
		}
		return -1, err
	}
	return 0, nil
}

func main() {
	fmt.Printf("PID: %d, PPID: %d\n", GetPID(), GetPPID())

	code, err := RunCommand("go", "version")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to run command: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("child exit code: %d\n", code)
}
