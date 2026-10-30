package main

import (
	"fmt"
	"runtime"
	"strings"
)

// EnvironmentInfo captures essential runtime details about the Go installation.
type EnvironmentInfo struct {
	Version  string
	GOOS     string
	GOARCH   string
	Compiler string
	NumCPU   int
	GOROOT   string
}

// GetEnvironmentInfo extracts active Go runtime and system metadata.
func GetEnvironmentInfo() EnvironmentInfo {
	return EnvironmentInfo{
		Version:  runtime.Version(),
		GOOS:     runtime.GOOS,
		GOARCH:   runtime.GOARCH,
		Compiler: runtime.Compiler,
		NumCPU:   runtime.NumCPU(),
		GOROOT:   runtime.GOROOT(),
	}
}

// VerifyInstallation checks whether the environment is properly configured.
func VerifyInstallation(info EnvironmentInfo) (bool, []string) {
	var issues []string

	if !strings.HasPrefix(info.Version, "go") && !strings.HasPrefix(info.Version, "devel") {
		issues = append(issues, "invalid Go version prefix detected")
	}
	if info.GOOS == "" {
		issues = append(issues, "target operating system (GOOS) is undefined")
	}
	if info.GOARCH == "" {
		issues = append(issues, "target architecture (GOARCH) is undefined")
	}
	if info.NumCPU <= 0 {
		issues = append(issues, "logical CPU count must be greater than zero")
	}

	return len(issues) == 0, issues
}

func main() {
	fmt.Println("=== Go Installation & Environment Diagnostics ===")

	info := GetEnvironmentInfo()
	fmt.Printf("Go Version:         %s\n", info.Version)
	fmt.Printf("Target OS (GOOS):   %s\n", info.GOOS)
	fmt.Printf("Target Arch (GOARCH): %s\n", info.GOARCH)
	fmt.Printf("Compiler:           %s\n", info.Compiler)
	fmt.Printf("Logical CPUs:       %d\n", info.NumCPU)
	if info.GOROOT != "" {
		fmt.Printf("GOROOT Directory:   %s\n", info.GOROOT)
	} else {
		fmt.Printf("GOROOT Directory:   (standard embedded default)\n")
	}

	ok, issues := VerifyInstallation(info)
	if !ok {
		fmt.Println("\nConfiguration Issues Detected:")
		for _, issue := range issues {
			fmt.Printf("  - %s\n", issue)
		}
		return
	}

	fmt.Println("\nStatus: Go toolchain installation is verified and ready for Module 02.")
}
