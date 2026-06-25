package main

import (
	"fmt"
	"unicode"
)

func reverseRunes(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func countLetters(s string) map[rune]int {
	counts := make(map[rune]int)
	for _, r := range s {
		if unicode.IsLetter(r) {
			counts[r]++
		}
	}
	return counts
}

func main() {
	s1 := "Hello, 世界 🚀"
	fmt.Println("original:", s1)
	fmt.Println("reversed:", reverseRunes(s1))
	fmt.Println("letter counts:", countLetters(s1))

	s2 := "café"
	fmt.Println("original:", s2)
	fmt.Println("reversed:", reverseRunes(s2))
	fmt.Println("letter counts:", countLetters(s2))
}
