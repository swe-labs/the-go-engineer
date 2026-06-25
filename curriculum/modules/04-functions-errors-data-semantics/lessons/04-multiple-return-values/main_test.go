package main

import "testing"

func TestParsePositive(t *testing.T) {
	tests := []struct {
		input   string
		wantVal int
		wantErr bool
	}{
		{"42", 42, false},
		{"1", 1, false},
		{"0", 0, true},
		{"-1", 0, true},
		{"abc", 0, true},
		{"999", 999, false},
	}
	for _, tc := range tests {
		gotVal, gotErr := parsePositive(tc.input)
		if tc.wantErr && gotErr == nil {
			t.Errorf("parsePositive(%q) expected error, got %d", tc.input, gotVal)
		}
		if !tc.wantErr && gotErr != nil {
			t.Errorf("parsePositive(%q) unexpected error: %v", tc.input, gotErr)
		}
		if !tc.wantErr && gotVal != tc.wantVal {
			t.Errorf("parsePositive(%q) = %d; want %d", tc.input, gotVal, tc.wantVal)
		}
	}
}
