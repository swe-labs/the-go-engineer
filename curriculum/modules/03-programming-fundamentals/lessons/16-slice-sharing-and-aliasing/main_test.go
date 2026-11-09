package main

import "testing"

func TestSplitAndCopyIsolation(t *testing.T) {
	data := []int{0, 1, 2, 3, 4, 5}
	chunks := splitAndCopy(data, 2)
	data[0] = 99
	data[2] = 88
	data[4] = 77
	for i, ch := range chunks {
		for j, v := range ch {
			if v == 99 || v == 88 || v == 77 {
				t.Errorf("chunks[%d][%d] = %d; should be isolated from changes", i, j, v)
			}
		}
	}
}

func TestSplitAndCopyValues(t *testing.T) {
	tests := []struct {
		data      []int
		chunkSize int
		want      [][]int
	}{
		{[]int{0, 1, 2, 3, 4, 5}, 2, [][]int{{0, 1}, {2, 3}, {4, 5}}},
		{[]int{0, 1, 2, 3, 4}, 2, [][]int{{0, 1}, {2, 3}, {4}}},
		{[]int{1}, 1, [][]int{{1}}},
		{[]int{}, 2, [][]int{}},
	}

	for _, tc := range tests {
		got := splitAndCopy(tc.data, tc.chunkSize)
		if len(got) != len(tc.want) {
			t.Errorf("splitAndCopy(%v, %d) got %d chunks; want %d", tc.data, tc.chunkSize, len(got), len(tc.want))
			continue
		}
		for i := range got {
			for j := range got[i] {
				if got[i][j] != tc.want[i][j] {
					t.Errorf("splitAndCopy(%v, %d)[%d][%d] = %d; want %d", tc.data, tc.chunkSize, i, j, got[i][j], tc.want[i][j])
				}
			}
		}
	}
}

func TestUnsafeSplitAliasing(t *testing.T) {
	data := []int{0, 1, 2, 3, 4, 5}
	chunks := unsafeSplit(data, 2)
	data[1] = 99
	if chunks[0][1] != 99 {
		t.Errorf("unsafeSplit chunk should alias: chunks[0][1] = %d; want 99", chunks[0][1])
	}
}

func TestSplitAndCopyEdgeCases(t *testing.T) {
	data := []int{1, 2, 3}
	chunks := splitAndCopy(data, 10)
	if len(chunks) != 1 || len(chunks[0]) != 3 {
		t.Errorf("splitAndCopy(%v, 10) = %v; want [[1 2 3]]", data, chunks)
	}
}

func TestUnsafeSplitValues(t *testing.T) {
	data := []int{0, 1, 2, 3}
	chunks := unsafeSplit(data, 2)
	want := [][]int{{0, 1}, {2, 3}}
	if len(chunks) != len(want) {
		t.Fatalf("got %d chunks", len(chunks))
	}
	for i := range chunks {
		for j := range chunks[i] {
			if chunks[i][j] != want[i][j] {
				t.Errorf("[%d][%d] = %d; want %d", i, j, chunks[i][j], want[i][j])
			}
		}
	}
}
