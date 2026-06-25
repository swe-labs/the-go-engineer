package main

import "testing"

func TestCanAccess(t *testing.T) {
	tests := []struct {
		name     string
		role     string
		isAdmin  bool
		isActive bool
		hour     int
		want     bool
	}{
		{"admin active any hour", "admin", true, true, 3, true},
		{"admin inactive", "admin", true, false, 3, false},
		{"active non-admin business hours", "user", false, true, 10, true},
		{"active non-admin after hours", "user", false, true, 20, false},
		{"inactive non-admin business hours", "user", false, false, 10, false},
		{"banned active business hours", "banned", false, true, 10, false},
		{"empty role", "", false, true, 10, false},
		{"business hour boundary 9am", "guest", false, true, 9, true},
		{"business hour boundary 5pm", "guest", false, true, 17, false},
		{"before business hours 8am", "user", false, true, 8, false},
		{"midnight", "user", false, true, 0, false},
	}

	for _, tc := range tests {
		got := canAccess(tc.role, tc.isAdmin, tc.isActive, tc.hour)
		if got != tc.want {
			t.Errorf("%s: canAccess(%q, %t, %t, %d) = %t; want %t",
				tc.name, tc.role, tc.isAdmin, tc.isActive, tc.hour, got, tc.want)
		}
	}
}

func TestShortCircuit(t *testing.T) {
	called := false
	mocked := func() bool {
		called = true
		return true
	}

	// false && ... should not evaluate the right side
	result := false && mocked()
	if result {
		t.Error("false && mocked() should be false")
	}
	if called {
		t.Error("mocked() was called on false && — short-circuit violated")
	}

	// true || ... should not evaluate the right side
	called = false
	result = true || mocked()
	if !result {
		t.Error("true || mocked() should be true")
	}
	if called {
		t.Error("mocked() was called on true || — short-circuit violated")
	}
}

func TestDeMorgan(t *testing.T) {
	for _, a := range []bool{true, false} {
		for _, b := range []bool{true, false} {
			// !(A && B) == !A || !B
			if !(a && b) != (!a || !b) {
				t.Errorf("De Morgan && failed for a=%t b=%t", a, b)
			}
			// !(A || B) == !A && !B
			if !(a || b) != (!a && !b) {
				t.Errorf("De Morgan || failed for a=%t b=%t", a, b)
			}
		}
	}
}

func TestOperatorPrecedence(t *testing.T) {
	// ! has higher precedence than &&, which has higher than ||
	if !true && false != ((!true) && false) {
		t.Error("! precedence violation")
	}
	if true && false || true != ((true && false) || true) {
		t.Error("&& vs || precedence violation")
	}
}

func TestLeapYearExpression(t *testing.T) {
	isLeap := func(year int) bool {
		return year%400 == 0 || (year%4 == 0 && year%100 != 0)
	}
	tests := []struct {
		year int
		want bool
	}{
		{2000, true},
		{1900, false},
		{2024, true},
		{2023, false},
		{2400, true},
		{2100, false},
	}
	for _, tc := range tests {
		got := isLeap(tc.year)
		if got != tc.want {
			t.Errorf("isLeap(%d) = %t; want %t", tc.year, got, tc.want)
		}
	}
}
