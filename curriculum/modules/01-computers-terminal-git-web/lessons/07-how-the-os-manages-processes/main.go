package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

// ProcessMetadata encapsulates basic identity of an active process.
type ProcessMetadata struct {
	PID  int
	PPID int
	OS   string
	Arch string
}

func main() {
	meta := ProcessMetadata{
		PID:  os.Getpid(),
		PPID: os.Getppid(),
		OS:   runtime.GOOS,
		Arch: runtime.GOARCH,
	}

	fmt.Println("=== Active Process Metadata ===")
	fmt.Printf("Process ID (PID):         %d\n", meta.PID)
	fmt.Printf("Parent Process ID (PPID): %d (Your Terminal Shell)\n", meta.PPID)
	fmt.Printf("Operating System / Arch:  %s / %s\n", meta.OS, meta.Arch)

	fmt.Println("\n=== Spawning and Managing a Child Process ===")
	var childCmd *exec.Cmd

	if runtime.GOOS == "windows" {
		childCmd = exec.Command("cmd.exe", "/c", "echo Hello from Child Process && exit 0")
	} else {
		childCmd = exec.Command("sh", "-c", "echo 'Hello from Child Process' && exit 0")
	}

	// Capture output
	output, err := childCmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Child process failed: %v\n", err)
		return
	}

	fmt.Printf("Child Output: %s", string(output))
	if childCmd.ProcessState != nil {
		fmt.Printf("Child Exit Code: %d (Success: %t)\n",
			childCmd.ProcessState.ExitCode(),
			childCmd.ProcessState.Success(),
		)
	}

	fmt.Println("\nProcess lifecycle demonstrated: Parent spawned Child -> Child executed -> Parent reaped exit code.")
}
