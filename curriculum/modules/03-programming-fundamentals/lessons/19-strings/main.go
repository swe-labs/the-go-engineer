package main

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
)

func CleanAndSplit(text string) ([]string, error) {
	cleaned := strings.TrimSpace(text)
	cleaned = strings.ToLower(cleaned)

	var filtered []rune
	for _, r := range cleaned {
		if unicode.IsLetter(r) || unicode.IsNumber(r) || unicode.IsSpace(r) {
			filtered = append(filtered, r)
		} else {
			filtered = append(filtered, ' ')
		}
	}

	words := strings.Fields(string(filtered))

	var result []string
	for _, w := range words {
		if len(w) >= 3 {
			result = append(result, w)
		}
	}

	if len(result) == 0 {
		return nil, errors.New("empty result")
	}
	return result, nil
}

func WordCount(words []string) map[string]int {
	freq := make(map[string]int)
	for _, w := range words {
		freq[w]++
	}
	return freq
}

func main() {
	input := "  Hello, World! Go is   amazing.  "
	words, err := CleanAndSplit(input)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Words:", words)

	freq := WordCount(words)
	fmt.Println("Frequency:", freq)

	if _, err := CleanAndSplit(""); err != nil {
		fmt.Println("Empty input error:", err)
	}

	if _, err := CleanAndSplit("a an"); err != nil {
		fmt.Println("Short words error:", err)
	}

	fmt.Println("\n--- CleanAndSplit demo ---")
	for i, w := range words {
		fmt.Printf("%d: %s\n", i, w)
	}
}
