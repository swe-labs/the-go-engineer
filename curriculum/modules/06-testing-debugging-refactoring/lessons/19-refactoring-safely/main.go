package main

import (
	"fmt"
	"strings"
)

func main() {
	names := []string{"alice", "bob", "charlie"}
	result := process(names)
	fmt.Println(result)
}

func process(names []string) string {
	result := ""
	for _, n := range names {
		result += strings.ToUpper(n[:1]) + strings.ToLower(n[1:]) + ","
	}
	return strings.TrimSuffix(result, ",")
}
