package main

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
)

//go:embed static/*
var embeddedFiles embed.FS

func main() {
	// Read from embedded FS
	data, err := fs.ReadFile(embeddedFiles, "static/hello.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "read embedded file error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Embedded file: %s", data)

	// Walk embedded FS
	fmt.Println("\nEmbedded files:")
	fs.WalkDir(embeddedFiles, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		info, _ := d.Info()
		fmt.Printf("  %s (%d bytes)\n", path, info.Size())
		return nil
	})

	// Open a directory on disk as an FS
	diskFS := os.DirFS(".")
	data, err = fs.ReadFile(diskFS, "main.go")
	if err != nil {
		fmt.Fprintf(os.Stderr, "read disk file error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("\nmain.go on disk: %d bytes\n", len(data))
}
