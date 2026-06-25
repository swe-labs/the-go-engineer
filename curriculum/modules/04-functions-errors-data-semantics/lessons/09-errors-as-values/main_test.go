package main

import "testing"

func TestParseAge(t *testing.T) {
	tests := []struct {
		input   string
		wantAge int
		wantErr error
	}{
		{"25", 25, nil},
		{"0", 0, nil},
		{"150", 150, nil},
		{"-1", 0, ErrNegativeAge},
		{"abc", 0, ErrInvalidAge},
		{"200", 0, ErrUnreasonableAge},
	}
	for _, tc := range tests {
		age, err := parseAge(tc.input)
		if !errorsMatch(err, tc.wantErr) {
			t.Errorf("parseAge(%q) error = %v; want %v", tc.input, err, tc.wantErr)
		}
		if tc.wantErr == nil && age != tc.wantAge {
			t.Errorf("parseAge(%q) = %d; want %d", tc.input, age, tc.wantAge)
		}
	}
}

func errorsMatch(got, want error) bool {
	if want == nil {
		return got == nil
	}
	return got == want
}
