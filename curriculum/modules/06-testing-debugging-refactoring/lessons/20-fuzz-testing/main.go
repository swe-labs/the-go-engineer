package main

import (
	"errors"
	"fmt"
	"unicode/utf8"
)

func main() {
	inputs := []string{"hello", "world", "héllo", ""}
	for _, s := range inputs {
		result, err := Reverse(s)
		if err != nil {
			fmt.Printf("Reverse(%q) error: %v\n", s, err)
		} else {
			fmt.Printf("Reverse(%q) = %q\n", s, result)
		}
	}
}

func Reverse(s string) (string, error) {
	if !utf8.ValidString(s) {
		return "", errors.New("invalid utf-8")
	}
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes), nil
}
