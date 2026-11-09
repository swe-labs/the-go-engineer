package main

import (
	"bytes"
	"testing"
)

func TestWordCount(t *testing.T) {
	tests := []struct {
		text []byte
		want map[string]int
	}{
		{[]byte("hello world hello Go"), map[string]int{"hello": 2, "world": 1, "Go": 1}},
		{[]byte(""), map[string]int{}},
		{[]byte("a b c a b a"), map[string]int{"a": 3, "b": 2, "c": 1}},
	}

	for _, tc := range tests {
		got := wordCount(tc.text)
		if len(got) != len(tc.want) {
			t.Errorf("wordCount(%q) len=%d; want len=%d", string(tc.text), len(got), len(tc.want))
		}
		for k, v := range tc.want {
			if got[k] != v {
				t.Errorf("wordCount(%q)[%q] = %d; want %d", string(tc.text), k, got[k], v)
			}
		}
	}
}

func TestToUpperASCII(t *testing.T) {
	tests := []struct {
		input []byte
		want  []byte
	}{
		{[]byte("hello"), []byte("HELLO")},
		{[]byte("Go"), []byte("GO")},
		{[]byte("123!@#"), []byte("123!@#")},
		{[]byte(""), []byte("")},
		{[]byte("abc DEF"), []byte("ABC DEF")},
	}

	for _, tc := range tests {
		got := toUpperASCII(tc.input)
		if !bytes.Equal(got, tc.want) {
			t.Errorf("toUpperASCII(%q) = %q; want %q", string(tc.input), string(got), string(tc.want))
		}
	}
}

func TestToUpperASCIIDoesNotModifyInput(t *testing.T) {
	input := []byte("hello")
	original := []byte("hello")
	_ = toUpperASCII(input)
	if !bytes.Equal(input, original) {
		t.Errorf("input was modified: %q", string(input))
	}
}
