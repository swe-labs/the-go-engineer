package main

import (
	"bytes"
	"fmt"
)

func wordCount(text []byte) map[string]int {
	counts := make(map[string]int)
	words := bytes.Fields(text)
	for _, w := range words {
		counts[string(w)]++
	}
	return counts
}

func toUpperASCII(b []byte) []byte {
	result := make([]byte, len(b))
	for i, c := range b {
		if c >= 'a' && c <= 'z' {
			result[i] = c - 32
		} else {
			result[i] = c
		}
	}
	return result
}

func main() {
	text := []byte("hello world hello Go")
	fmt.Println("wordCount:", wordCount(text))

	upper := toUpperASCII(text)
	fmt.Println("toUpperASCII:", string(upper))

	var buf bytes.Buffer
	buf.WriteString("bytes.Buffer example: ")
	buf.Write(upper)
	fmt.Println(buf.String())
}
