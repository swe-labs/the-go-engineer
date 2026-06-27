package main

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

type EnvVar struct {
	Key   string
	Value string
	Set   bool
}

func GetEnv(key string) EnvVar {
	val, ok := os.LookupEnv(key)
	return EnvVar{
		Key:   key,
		Value: val,
		Set:   ok,
	}
}

func main() {
	fmt.Println("=== Environment Variables ===")
	fmt.Println()

	vars := []string{"PATH", "HOME", "USER", "GOPATH", "GOROOT", "MY_CUSTOM_VAR"}
	for _, key := range vars {
		ev := GetEnv(key)
		if ev.Set {
			fmt.Printf("%s = %s\n", ev.Key, ev.Value)
		} else {
			fmt.Printf("%s = [not set]\n", ev.Key)
		}
	}

	fmt.Println()
	fmt.Println("All environment variables (sorted):")
	fmt.Println()

	all := os.Environ()
	sort.Strings(all)
	for _, pair := range all {
		parts := strings.SplitN(pair, "=", 2)
		if len(parts) == 2 {
			fmt.Printf("%s = %s\n", parts[0], parts[1])
		}
	}
}
