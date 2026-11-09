package main

import (
	"fmt"
	"strings"
)

func concatPlus(parts []string) string {
	result := ""
	for _, p := range parts {
		result += p
	}
	return result
}

func concatBuilder(parts []string) string {
	var b strings.Builder
	for _, p := range parts {
		b.WriteString(p)
	}
	return b.String()
}

func main() {
	parts := []string{"hello", " ", "world", " ", "from", " ", "Go"}
	fmt.Println("concatPlus:", concatPlus(parts))
	fmt.Println("concatBuilder:", concatBuilder(parts))
}
