package main

import "testing"

func TestSecondsPerWeek(t *testing.T) {
	want := 7 * 24 * 60 * 60
	if SecondsPerWeek != want {
		t.Errorf("SecondsPerWeek = %d; want %d", SecondsPerWeek, want)
	}
}

func TestWeekdayConstants(t *testing.T) {
	tests := []struct {
		day  int
		name string
	}{
		{Monday, "Monday"},
		{Tuesday, "Tuesday"},
		{Wednesday, "Wednesday"},
		{Thursday, "Thursday"},
		{Friday, "Friday"},
		{Saturday, "Saturday"},
		{Sunday, "Sunday"},
	}

	for _, tc := range tests {
		got := weekdayName(tc.day)
		if got != tc.name {
			t.Errorf("weekdayName(%d) = %q; want %q", tc.day, got, tc.name)
		}
	}

	if Monday != 0 {
		t.Errorf("Monday should be 0, got %d", Monday)
	}
	if Tuesday != 1 {
		t.Errorf("Tuesday should be 1, got %d", Tuesday)
	}
	if Sunday != 6 {
		t.Errorf("Sunday should be 0, got %d", Sunday)
	}
}

func TestUntypedConstantsAssignable(t *testing.T) {
	var a int = DaysInWeek
	if a != 7 {
		t.Errorf("DaysInWeek assigned to int = %d; want 7", a)
	}
	var b float64 = DaysInWeek
	if b != 7.0 {
		t.Errorf("DaysInWeek assigned to float64 = %f; want 7.0", b)
	}
}

func TestWeekdayNameUnknown(t *testing.T) {
	if got := weekdayName(99); got != "unknown" {
		t.Errorf("weekdayName(99) = %q; want \"unknown\"", got)
	}
}
