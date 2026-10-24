package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

func main() {
	fmt.Println("=== Terminal Subprocess & Stream Redirection Demo ===")

	var cmd *exec.Cmd
	var inputData string

	if runtime.GOOS == "windows" {
		// Use PowerShell on Windows
		inputData = "apple\nbanana\ncherry\nblueberry\napple pie\n"
		cmd = exec.Command("powershell", "-NoProfile", "-Command", "$input | Select-String 'apple'")
	} else {
		// Use standard POSIX grep on Linux/macOS
		inputData = "apple\nbanana\ncherry\nblueberry\napple pie\n"
		cmd = exec.Command("grep", "apple")
	}

	// 1. Wire stdin pipe with simulated input data
	cmd.Stdin = strings.NewReader(inputData)

	// 2. Attach separate memory buffers for stdout (fd 1) and stderr (fd 2)
	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf

	fmt.Println("Feeding stdin to child process:")
	fmt.Print(inputData)

	// 3. Execute process and wait for termination
	err := cmd.Run()

	fmt.Println("\n=== Process Output Results ===")
	fmt.Printf("Captured stdout (fd 1):\n%s", stdoutBuf.String())

	if stderrBuf.Len() > 0 {
		fmt.Printf("Captured stderr (fd 2):\n%s\n", stderrBuf.String())
	} else {
		fmt.Println("Captured stderr (fd 2): (clean, 0 error bytes)")
	}

	// 4. Inspect exit code
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			fmt.Printf("Process exited with non-zero code: %d\n", exitErr.ExitCode())
		} else {
			fmt.Printf("Execution failed: %v\n", err)
		}
	} else {
		fmt.Println("Process exited successfully with Exit Code 0.")
	}
}
