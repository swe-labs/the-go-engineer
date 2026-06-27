package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func CreateFile(path, content string) (string, int64, bool, string, error) {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", 0, false, "", err
	}
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return "", 0, false, "", err
	}
	return ReadFile(path)
}

func ReadFile(path string) (string, int64, bool, string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", 0, false, "", err
	}
	info, err := os.Stat(path)
	if err != nil {
		return "", 0, false, "", err
	}
	return path, info.Size(), info.IsDir(), string(data), nil
}

func main() {
	fmt.Println("=== Files, Bytes, Directories, and Paths ===")
	fmt.Println()

	baseDir, _ := os.Getwd()
	filePath := filepath.Join(baseDir, "output", "hello.txt")

	path, size, isDir, content, err := CreateFile(filePath, "Hello, filesystem!")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Created file:  %s\n", path)
	fmt.Printf("Absolute path: %s\n", filePath)
	fmt.Printf("Size:          %d bytes\n", size)
	fmt.Printf("Is directory?  %v\n", isDir)
	fmt.Printf("Content:       %q\n", content)
	fmt.Println()

	abs, _ := filepath.Abs(filePath)
	fmt.Printf("Absolute: %s\n", abs)
	fmt.Printf("Base:     %s\n", filepath.Base(filePath))
	fmt.Printf("Dir:      %s\n", filepath.Dir(filePath))
	fmt.Printf("Ext:      %s\n", filepath.Ext(filePath))
	fmt.Println()

	os.RemoveAll(filepath.Join(baseDir, "output"))
	fmt.Println("Cleaned up temporary files.")
}
