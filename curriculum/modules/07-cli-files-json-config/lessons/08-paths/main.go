package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	// Join paths portably
	base := "docs"
	sub := "2024"
	file := "report.txt"
	fullPath := filepath.Join(base, sub, file)
	fmt.Println("Joined path:", fullPath)

	// Clean a messy path
	messy := "docs/../docs//2024/./report.txt"
	clean := filepath.Clean(messy)
	fmt.Println("Cleaned path:", clean)

	// Split a path
	dir, name := filepath.Split(fullPath)
	fmt.Printf("Dir: %q, File: %q\n", dir, name)

	// Walk directory tree
	root := "."
	fmt.Println("Walking directory tree from", root)
	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		fmt.Printf("  %s (%d bytes)\n", path, info.Size())
		return nil
	})
}
