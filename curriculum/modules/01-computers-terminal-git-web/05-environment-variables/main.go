package main

import (
	"fmt"
	"os"
	"strings"
)

func GetEnv(key string) (string, bool) {
	return os.LookupEnv(key)
}

func main() {
	fmt.Println("=== Environment Variables ===")
	fmt.Println()

	vars := []string{"PATH", "HOME", "USER", "GOPATH", "GOROOT", "MY_CUSTOM_VAR"}
	for _, key := range vars {
		val, ok := GetEnv(key)
		if ok {
			fmt.Printf("%s = %s\n", key, val)
		} else {
			fmt.Printf("%s = [not set]\n", key)
		}
	}

	fmt.Println()
	fmt.Println("All environment variables:")
	fmt.Println()

	all := os.Environ()
	for _, pair := range all {
		parts := strings.SplitN(pair, "=", 2)
		if len(parts) == 2 {
			fmt.Printf("%s = %s\n", parts[0], parts[1])
		}
	}
}
