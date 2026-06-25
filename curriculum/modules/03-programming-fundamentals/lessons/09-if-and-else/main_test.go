package main

import "testing"

func TestClassifyTemperature(t *testing.T) {
	tests := []struct {
		temp int
		want string
	}{
		{45, "dangerously hot"},
		{41, "dangerously hot"},
		{40, "hot"},
		{35, "hot"},
		{31, "hot"},
		{30, "warm"},
		{25, "warm"},
		{21, "warm"},
		{20, "mild"},
		{15, "mild"},
		{11, "mild"},
		{10, "cool"},
		{5, "cool"},
		{0, "cool"},
		{-1, "freezing"},
		{-100, "freezing"},
	}
	for _, tc := range tests {
		got := classifyTemperature(tc.temp)
		if got != tc.want {
			t.Errorf("classifyTemperature(%d) = %q; want %q", tc.temp, got, tc.want)
		}
	}
}

func TestDivide(t *testing.T) {
	if _, err := divide(5, 0); err == nil {
		t.Error("divide(5, 0) expected error, got nil")
	}
	if result, err := divide(10, 3); err != nil {
		t.Errorf("divide(10, 3) unexpected error: %v", err)
	} else if result != 10.0/3.0 {
		t.Errorf("divide(10, 3) = %f; want %f", result, 10.0/3.0)
	}
}

func TestMax(t *testing.T) {
	tests := []struct {
		a, b, want int
	}{
		{3, 5, 5},
		{5, 3, 5},
		{4, 4, 4},
		{-1, 0, 0},
		{-5, -3, -3},
	}
	for _, tc := range tests {
		got := max(tc.a, tc.b)
		if got != tc.want {
			t.Errorf("max(%d, %d) = %d; want %d", tc.a, tc.b, got, tc.want)
		}
	}
}

func TestValidateAge(t *testing.T) {
	tests := []struct {
		s       string
		wantErr bool
	}{
		{"25", false},
		{"18", false},
		{"17", true},
		{"0", true},
		{"abc", true},
		{"200", true},
		{"151", true},
		{"-5", true},
	}
	for _, tc := range tests {
		err := validateAge(tc.s)
		if tc.wantErr && err == nil {
			t.Errorf("validateAge(%q) expected error, got nil", tc.s)
		}
		if !tc.wantErr && err != nil {
			t.Errorf("validateAge(%q) unexpected error: %v", tc.s, err)
		}
	}
}

func TestGradeClassification(t *testing.T) {
	tests := []struct {
		score int
		want  string
	}{
		{-5, "invalid"},
		{0, "F"},
		{45, "F"},
		{59, "F"},
		{60, "D"},
		{69, "D"},
		{70, "C"},
		{79, "C"},
		{80, "B"},
		{89, "B"},
		{90, "A"},
		{95, "A"},
		{100, "A"},
		{101, "invalid"},
	}
	for _, tc := range tests {
		got := gradeClassification(tc.score)
		if got != tc.want {
			t.Errorf("gradeClassification(%d) = %q; want %q", tc.score, got, tc.want)
		}
	}
}
