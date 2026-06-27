package main

import (
	"fmt"
	"os"
	"path/filepath"
)

type FileInfo struct {
	Path     string
	Size     int64
	IsDir    bool
	Content  string
}

func CreateFile(path, content string) (FileInfo, error) {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return FileInfo{}, err
	}
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return FileInfo{}, err
	}
	return ReadFile(path)
}

func ReadFile(path string) (FileInfo, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return FileInfo{}, err
	}
	info, err := os.Stat(path)
	if err != nil {
		return FileInfo{}, err
	}
	return FileInfo{
		Path:    path,
		Size:    info.Size(),
		IsDir:   info.IsDir(),
		Content: string(data),
	}, nil
}

func main() {
	fmt.Println("=== Files, Bytes, Directories, and Paths ===")
	fmt.Println()

	baseDir, _ := os.Getwd()
	filePath := filepath.Join(baseDir, "output", "hello.txt")

	info, err := CreateFile(filePath, "Hello, filesystem!")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Created file:  %s\n", info.Path)
	fmt.Printf("Absolute path: %s\n", filePath)
	fmt.Printf("Size:          %d bytes\n", info.Size)
	fmt.Printf("Is directory?  %v\n", info.IsDir)
	fmt.Printf("Content:       %q\n", info.Content)
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
