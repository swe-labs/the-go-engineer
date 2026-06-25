package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func WriteConfig(dir, content string) (string, error) {
	path := filepath.Join(dir, "config.txt")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return "", err
	}
	return path, nil
}

func ReadConfig(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func SetupEnv(key, value string) func() {
	orig, existed := os.LookupEnv(key)
	if err := os.Setenv(key, value); err != nil {
		panic(fmt.Sprintf("Setenv failed: %v", err))
	}
	if existed {
		return func() { os.Setenv(key, orig) }
	}
	return func() { os.Unsetenv(key) }
}

func main() {
	dir, _ := os.MkdirTemp("", "example")
	defer os.RemoveAll(dir)

	path, _ := WriteConfig(dir, "hello world")
	content, _ := ReadConfig(path)
	fmt.Println(content)

	cleanup := SetupEnv("MY_VAR", "hello")
	defer cleanup()
	fmt.Println("MY_VAR =", os.Getenv("MY_VAR"))
}
