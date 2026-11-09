package main

import "testing"

func TestWordFrequency(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected map[string]int
	}{
		{
			name:     "empty string",
			input:    "",
			expected: map[string]int{},
		},
		{
			name:     "single word",
			input:    "hello",
			expected: map[string]int{"hello": 1},
		},
		{
			name:     "multiple words",
			input:    "hello world hello",
			expected: map[string]int{"hello": 2, "world": 1},
		},
		{
			name:     "case insensitivity",
			input:    "Hello HELLO hello",
			expected: map[string]int{"hello": 3},
		},
		{
			name:     "punctuation stripping",
			input:    "hello! world. hello?",
			expected: map[string]int{"hello": 2, "world": 1},
		},
		{
			name:     "numbers preserved",
			input:    "go1 go2 go1",
			expected: map[string]int{"go1": 2, "go2": 1},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := WordFrequency(tc.input)
			if len(got) != len(tc.expected) {
				t.Errorf("len = %d; want %d; got %v", len(got), len(tc.expected), got)
			}
			for k, v := range tc.expected {
				if got[k] != v {
					t.Errorf("freq[%q] = %d; want %d", k, got[k], v)
				}
			}
		})
	}
}

func TestTopWords(t *testing.T) {
	tests := []struct {
		name     string
		freq     map[string]int
		n        int
		expected []WordCount
	}{
		{
			name:     "nil map",
			freq:     nil,
			n:        3,
			expected: nil,
		},
		{
			name:     "empty map",
			freq:     map[string]int{},
			n:        3,
			expected: nil,
		},
		{
			name:     "n=0 returns nil",
			freq:     map[string]int{"a": 1},
			n:        0,
			expected: nil,
		},
		{
			name:     "n larger than entries",
			freq:     map[string]int{"a": 1, "b": 2},
			n:        10,
			expected: []WordCount{{"b", 2}, {"a", 1}},
		},
		{
			name:     "top 2 of 4",
			freq:     map[string]int{"a": 1, "b": 5, "c": 3, "d": 2},
			n:        2,
			expected: []WordCount{{"b", 5}, {"c", 3}},
		},
		{
			name:     "tie broken alphabetically",
			freq:     map[string]int{"zebra": 5, "apple": 5},
			n:        2,
			expected: []WordCount{{"apple", 5}, {"zebra", 5}},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := TopWords(tc.freq, tc.n)
			if len(got) != len(tc.expected) {
				t.Fatalf("len = %d; want %d; got %v", len(got), len(tc.expected), got)
			}
			for i := range got {
				if got[i] != tc.expected[i] {
					t.Errorf("result[%d] = %+v; want %+v", i, got[i], tc.expected[i])
				}
			}
		})
	}
}

func TestNilMapBehavior(t *testing.T) {
	var m map[string]int

	if l := len(m); l != 0 {
		t.Errorf("len(nil map) = %d; want 0", l)
	}

	if v := m["anything"]; v != 0 {
		t.Errorf("read from nil map = %d; want 0", v)
	}

	delete(m, "anything")
}
