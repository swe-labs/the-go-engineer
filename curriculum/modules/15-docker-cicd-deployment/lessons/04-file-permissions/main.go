package main

import (
	"fmt"
	"os"
)

func writeSecureFile(path string, data []byte) error {
	return os.WriteFile(path, data, 0600)
}

func readPermissions(path string) (os.FileMode, error) {
	info, err := os.Stat(path)
	if err != nil {
		return 0, err
	}
	return info.Mode().Perm(), nil
}

func main() {
	tmpFile, err := os.CreateTemp("", "secure-*.txt")
	if err != nil {
		panic(err)
	}
	path := tmpFile.Name()
	tmpFile.Close()
	defer os.Remove(path)

	if err := writeSecureFile(path, []byte("secret data")); err != nil {
		panic(err)
	}

	mode, err := readPermissions(path)
	if err != nil {
		panic(err)
	}
	fmt.Printf("File mode: %#o\n", mode)
	fmt.Printf("Is secure: %v\n", mode == 0600)
}
