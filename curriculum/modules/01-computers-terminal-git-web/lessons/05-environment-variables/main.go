package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

// EnvVar holds key, value, and presence status for an environment variable.
type EnvVar struct {
	Key   string
	Value string
	IsSet bool
}

// GetEnv retrieves an environment variable distinguishing between empty and unset.
func GetEnv(key string) EnvVar {
	val, exists := os.LookupEnv(key)
	return EnvVar{
		Key:   key,
		Value: val,
		IsSet: exists,
	}
}

func main() {
	fmt.Println("=== Checking Critical Environment Variables ===")
	targetKeys := []string{"PATH", "HOME", "USER", "GOROOT", "GOPATH", "MY_CONFIG_SECRET"}
	for _, key := range targetKeys {
		v := GetEnv(key)
		if v.IsSet {
			displayVal := v.Value
			if len(displayVal) > 40 {
				displayVal = displayVal[:37] + "..."
			}
			fmt.Printf("✓ %-18s = %q\n", v.Key, displayVal)
		} else {
			fmt.Printf("✗ %-18s = (unset)\n", v.Key)
		}
	}

	fmt.Println("\n=== Deconstructing System $PATH Directories ===")
	pathVar := os.Getenv("PATH")
	delimiter := string(os.PathListSeparator) // ":" on Unix, ";" on Windows
	dirs := filepath.SplitList(pathVar)

	fmt.Printf("OS Path Delimiter: %q | Total PATH Directories: %d\n", delimiter, len(dirs))
	limit := 5
	if len(dirs) < limit {
		limit = len(dirs)
	}
	fmt.Printf("First %d directories searched by shell (in priority order):\n", limit)
	for i := 0; i < limit; i++ {
		fmt.Printf("  [%d] %s\n", i+1, dirs[i])
	}

	fmt.Println("\n=== Binary Resolution via $PATH (exec.LookPath) ===")
	lookupBinary := "go"
	if runtime.GOOS == "windows" {
		lookupBinary = "go.exe"
	}

	resolvedPath, err := exec.LookPath(lookupBinary)
	if err == nil {
		fmt.Printf("✓ Found '%s' on system PATH: %s\n", lookupBinary, resolvedPath)
	} else {
		fmt.Printf("✗ Binary '%s' was NOT found in any $PATH directory!\n", lookupBinary)
	}
}
