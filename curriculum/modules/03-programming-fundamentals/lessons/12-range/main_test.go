package main

import "testing"

func TestWordCount(t *testing.T) {
	tests := []struct {
		input string
		want  map[string]int
	}{
		{"hello world hello", map[string]int{"hello": 2, "world": 1}},
		{"", map[string]int{}},
		{"hello", map[string]int{"hello": 1}},
		{"a b c a b a", map[string]int{"a": 3, "b": 2, "c": 1}},
		{"  leading and trailing  ", map[string]int{"leading": 1, "and": 1, "trailing": 1}},
		{"Hello, 世界! Hello.", map[string]int{"Hello": 2, "世界": 1}},
	}

	for _, tc := range tests {
		got := wordCount(tc.input)
		if len(got) != len(tc.want) {
			t.Errorf("wordCount(%q) = %v; want %v", tc.input, got, tc.want)
			continue
		}
		for k, v := range tc.want {
			if got[k] != v {
				t.Errorf("wordCount(%q)[%q] = %d; want %d", tc.input, k, got[k], v)
			}
		}
	}
}

func TestMergeCounts(t *testing.T) {
	tests := []struct {
		maps []map[string]int
		want map[string]int
	}{
		{
			[]map[string]int{{"a": 1}, {"b": 2}},
			map[string]int{"a": 1, "b": 2},
		},
		{
			[]map[string]int{{"a": 1, "b": 2}, {"a": 3, "c": 4}},
			map[string]int{"a": 4, "b": 2, "c": 4},
		},
		{
			[]map[string]int{{}, {}},
			map[string]int{},
		},
		{
			[]map[string]int{{"x": 5}},
			map[string]int{"x": 5},
		},
		{
			[]map[string]int{{"a": 1}, {"a": 2}, {"a": 3}},
			map[string]int{"a": 6},
		},
	}

	for _, tc := range tests {
		got := mergeCounts(tc.maps...)
		if len(got) != len(tc.want) {
			t.Errorf("mergeCounts(%v) = %v; want %v", tc.maps, got, tc.want)
			continue
		}
		for k, v := range tc.want {
			if got[k] != v {
				t.Errorf("mergeCounts(%v)[%q] = %d; want %d", tc.maps, k, got[k], v)
			}
		}
	}
}

func TestRangeSliceCopySemantics(t *testing.T) {
	nums := []int{1, 2, 3}
	for i, v := range nums {
		// v is a copy of nums[i] at this iteration — modifying v does NOT affect nums
		if i == 0 {
			nums[1] = 999 // subsequent iterations re-read nums[i] — v will be 999
		}
		if i == 0 && v != 1 {
			t.Errorf("expected v=1, got v=%d", v)
		}
	}
}

func TestRangeArrayCopySemantics(t *testing.T) {
	arr := [3]int{1, 2, 3}
	for i, v := range arr {
		if i == 0 {
			arr[1] = 999 // modify original — should not affect range copy
		}
		if i == 1 && v != 2 {
			t.Errorf("range over array copied value should be 2, got %d", v)
		}
	}
}

func TestRangeMapIteration(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	seen := make(map[string]bool)
	for k := range m {
		seen[k] = true
	}
	if len(seen) != 3 {
		t.Errorf("expected 3 keys, got %d", len(seen))
	}
	for k := range m {
		if !seen[k] {
			t.Errorf("key %q not visited", k)
		}
	}
}

func TestRangeStringRunes(t *testing.T) {
	s := "世"
	count := 0
	for i, r := range s {
		if i != 0 {
			t.Errorf("expected byte offset 0, got %d", i)
		}
		if r != 0x4E16 {
			t.Errorf("expected rune U+4E16, got U+%04X", r)
		}
		count++
	}
	if count != 1 {
		t.Errorf("expected 1 rune, got %d", count)
	}
}
