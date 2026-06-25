package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	// Create a single directory
	err := os.Mkdir("mydir", 0755)
	if err != nil && !os.IsExist(err) {
		fmt.Fprintf(os.Stderr, "mkdir error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Created mydir/")

	// Create nested directories
	err = os.MkdirAll("a/b/c/d", 0755)
	if err != nil {
		fmt.Fprintf(os.Stderr, "mkdirall error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Created a/b/c/d/")

	// List directory contents
	entries, err := os.ReadDir(".")
	if err != nil {
		fmt.Fprintf(os.Stderr, "readdir error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("\nCurrent directory contents:")
	for _, e := range entries {
		info, _ := e.Info()
		size := info.Size()
		mode := e.Type()
		fmt.Printf("  %s %s (%d bytes)\n", mode, e.Name(), size)
	}

	// Remove recursively
	err = os.RemoveAll("a")
	if err != nil {
		fmt.Fprintf(os.Stderr, "removeall error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("\nRemoved a/")

	// Remove single directory (must be empty)
	err = os.Remove("mydir")
	if err != nil {
		fmt.Fprintf(os.Stderr, "remove error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Removed mydir/")

	// Use filepath.Walk to traverse
	fmt.Println("\nWalking root:")
	filepath.WalkDir(".", func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		fmt.Printf("  %s\n", path)
		return nil
	})
}
