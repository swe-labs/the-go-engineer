package main

import (
	"fmt"
	"unicode"
)

func wordCount(s string) map[string]int {
	counts := make(map[string]int)
	start := -1
	for i, r := range s {
		if unicode.IsLetter(r) {
			if start == -1 {
				start = i
			}
		} else {
			if start != -1 {
				counts[s[start:i]]++
				start = -1
			}
		}
	}
	if start != -1 {
		counts[s[start:]]++
	}
	return counts
}

func mergeCounts(maps ...map[string]int) map[string]int {
	merged := make(map[string]int)
	for _, m := range maps {
		for k, v := range m {
			merged[k] += v
		}
	}
	return merged
}

func main() {
	// Slice range demo
	nums := []int{10, 20, 30}
	for i, v := range nums {
		fmt.Printf("slice[%d] = %d\n", i, v)
	}

	// String range demo (runes)
	s := "Hi, 世界"
	for i, r := range s {
		fmt.Printf("byte=%d rune=%c\n", i, r)
	}

	// Map range demo
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	for k, v := range m {
		fmt.Printf("map[%q] = %d\n", k, v)
	}

	// Practice task
	wc := wordCount("hello world hello")
	fmt.Println("wordCount:", wc)

	m1 := map[string]int{"hello": 2, "world": 1}
	m2 := map[string]int{"hello": 1, "go": 3}
	merged := mergeCounts(m1, m2)
	fmt.Println("merged:", merged)
}
