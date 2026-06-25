package main

import (
	"fmt"
	"sort"
	"strings"
	"unicode"
)

func WordFrequency(text string) map[string]int {
	freq := make(map[string]int)
	for _, word := range strings.Fields(text) {
		cleaned := strings.TrimFunc(strings.ToLower(word), func(r rune) bool {
			return !unicode.IsLetter(r) && !unicode.IsNumber(r)
		})
		if cleaned != "" {
			freq[cleaned]++
		}
	}
	return freq
}

type WordCount struct {
	Word  string
	Count int
}

func TopWords(freq map[string]int, n int) []WordCount {
	if len(freq) == 0 || n <= 0 {
		return nil
	}
	wcs := make([]WordCount, 0, len(freq))
	for word, count := range freq {
		wcs = append(wcs, WordCount{Word: word, Count: count})
	}
	sort.Slice(wcs, func(i, j int) bool {
		if wcs[i].Count != wcs[j].Count {
			return wcs[i].Count > wcs[j].Count
		}
		return wcs[i].Word < wcs[j].Word
	})
	if n > len(wcs) {
		n = len(wcs)
	}
	return wcs[:n]
}

func main() {
	text := `The quick brown fox jumps over the lazy dog. The dog barks, and the fox runs away. Quick as a flash, the fox disappears into the woods. The dog sniffs the ground but the fox is gone. What a chase! The fox wins again.`

	freq := WordFrequency(text)
	fmt.Println("Total unique words:", len(freq))

	for w, c := range freq {
		if c > 1 {
			fmt.Printf("  %s: %d\n", w, c)
		}
	}

	fmt.Println("\nTop 5 words:")
	for _, wc := range TopWords(freq, 5) {
		fmt.Printf("  %s: %d\n", wc.Word, wc.Count)
	}

	empty := WordFrequency("")
	fmt.Println("\nEmpty input length:", len(empty))
	fmt.Println("TopWords on empty:", TopWords(empty, 5))
}
