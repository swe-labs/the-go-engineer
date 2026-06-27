package main

import (
	"fmt"
	"os"
	"os/exec"
)

func GetPID() int {
	return os.Getpid()
}

func GetPPID() int {
	return os.Getppid()
}

func RunCommand(name string, args ...string) (int, error) {
	cmd := exec.Command(name, args...)
	err := cmd.Run()
	if err != nil {
		if cmd.ProcessState != nil {
			return cmd.ProcessState.ExitCode(), nil
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
