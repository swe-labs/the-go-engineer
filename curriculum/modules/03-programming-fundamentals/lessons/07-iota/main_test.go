package main

import "testing"

func TestFlagValues(t *testing.T) {
	if FlagA != 1 {
		t.Errorf("FlagA = %d; want 1", FlagA)
	}
	if FlagB != 2 {
		t.Errorf("FlagB = %d; want 2", FlagB)
	}
	if FlagC != 4 {
		t.Errorf("FlagC = %d; want 4", FlagC)
	}
	if FlagD != 8 {
		t.Errorf("FlagD = %d; want 8", FlagD)
	}
}

func TestSet(t *testing.T) {
	var f Flag
	f = set(f, FlagA)
	if f != FlagA {
		t.Errorf("set(0, FlagA) = %d; want %d", f, FlagA)
	}
	f = set(f, FlagC)
	if f != FlagA|FlagC {
		t.Errorf("set(FlagA, FlagC) = %d; want %d", f, FlagA|FlagC)
	}
}

func TestClear(t *testing.T) {
	f := FlagA | FlagB | FlagC
	f = clear(f, FlagB)
	if f != FlagA|FlagC {
		t.Errorf("clear(FlagA|FlagB|FlagC, FlagB) = %d; want %d", f, FlagA|FlagC)
	}
	f = clear(f, FlagA)
	if f != FlagC {
		t.Errorf("clear(FlagA|FlagC, FlagA) = %d; want %d", f, FlagC)
	}
}

func TestHas(t *testing.T) {
	f := FlagB | FlagD
	tests := []struct {
		flag Flag
		want bool
	}{
		{FlagA, false},
		{FlagB, true},
		{FlagC, false},
		{FlagD, true},
	}
	for _, tc := range tests {
		got := has(f, tc.flag)
		if got != tc.want {
			t.Errorf("has(%d, %d) = %t; want %t", f, tc.flag, got, tc.want)
		}
	}
}

func TestZeroFlag(t *testing.T) {
	var f Flag
	if has(f, FlagA) {
		t.Error("zero value should not have any flag set")
	}
	if has(f, FlagB) {
		t.Error("zero value should not have any flag set")
	}
}
