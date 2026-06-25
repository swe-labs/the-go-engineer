package main

import (
	"errors"
	"fmt"
	"os"
)

func writeToFile(filename, content string) (err error) {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}
	defer func() {
		if cerr := file.Close(); cerr != nil && err == nil {
			err = fmt.Errorf("close file: %w", cerr)
		}
	}()
	_, err = file.WriteString(content)
	if err != nil {
		return fmt.Errorf("write file: %w", err)
	}
	return nil
}

func safeProcess(items []string) error {
	var errs []error
	for _, item := range items {
		filename := item + ".txt"
		content := "Content for " + item
		if err := writeToFileToLoop(filename, content); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

// writeToFileToLoop opens and closes a file within a single iteration (no defer in loop).
func writeToFileToLoop(filename, content string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("create: %w", err)
	}
	_, err = file.WriteString(content)
	closeErr := file.Close()
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	if closeErr != nil {
		return fmt.Errorf("close: %w", closeErr)
	}
	return nil
}

func main() {
	if err := writeToFile("hello.txt", "Hello, World!"); err != nil {
		fmt.Println("writeToFile error:", err)
	} else {
		fmt.Println("writeToFile: success")
	}

	if err := safeProcess([]string{"a", "b"}); err != nil {
		fmt.Println("safeProcess error:", err)
	} else {
		fmt.Println("safeProcess: success")
	}

	// Cleanup test files
	os.Remove("hello.txt")
	os.Remove("a.txt")
	os.Remove("b.txt")
}
