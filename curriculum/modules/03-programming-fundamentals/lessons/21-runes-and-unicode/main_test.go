package main

import "testing"

func TestReverseRunes(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"hello", "olleh"},
		{"café", "éfac"},
		{"世界", "界世"},
		{"", ""},
		{"a", "a"},
		{"🚀", "🚀"},
		{"Hello, 世界 🚀", "🚀 界世 ,olleH"},
	}

	for _, tc := range tests {
		got := reverseRunes(tc.input)
		if got != tc.want {
			t.Errorf("reverseRunes(%q) = %q; want %q", tc.input, got, tc.want)
		}
	}
}

func TestCountLetters(t *testing.T) {
	tests := []struct {
		input string
		want  map[rune]int
	}{
		{"hello", map[rune]int{'h': 1, 'e': 1, 'l': 2, 'o': 1}},
		{"Hello, 世界 🚀", map[rune]int{'H': 1, 'e': 1, 'l': 2, 'o': 1, '世': 1, '界': 1}},
		{"", map[rune]int{}},
		{"123!@#", map[rune]int{}},
		{"café", map[rune]int{'c': 1, 'a': 1, 'f': 1, 'é': 1}},
	}

	for _, tc := range tests {
		got := countLetters(tc.input)
		if len(got) != len(tc.want) {
			t.Errorf("countLetters(%q) len=%d; want len=%d", tc.input, len(got), len(tc.want))
			continue
		}
		for k, v := range tc.want {
			if got[k] != v {
				t.Errorf("countLetters(%q)[%c] = %d; want %d", tc.input, k, got[k], v)
			}
		}
	}
}

func TestReverseRunesMultiByte(t *testing.T) {
	s := "éfac"
	got := reverseRunes(s)
	if got != "café" {
		t.Errorf("reverseRunes(%q) = %q; want %q", s, got, "café")
	}
}
