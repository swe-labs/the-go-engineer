package main

import "testing"

func TestCalc(t *testing.T) {
	tests := []struct {
		a, b float64
		op   string
		want float64
		err  bool
	}{
		{10, 3, "add", 13, false},
		{10, 3, "sum", 13, false},
		{10, 3, "+", 13, false},
		{10, 3, "sub", 7, false},
		{10, 3, "diff", 7, false},
		{10, 3, "-", 7, false},
		{10, 3, "mul", 30, false},
		{10, 3, "prod", 30, false},
		{10, 3, "*", 30, false},
		{10, 3, "div", 3.3333333333333335, false},
		{10, 3, "/", 3.3333333333333335, false},
		{10, 0, "div", 0, true},
		{10, 0, "/", 0, true},
		{0, 0, "div", 0, true},
		{10, 3, "unknown", 0, true},
		{10, 3, "mod", 0, true},
	}

	for _, tc := range tests {
		got, err := calc(tc.a, tc.b, tc.op)
		if tc.err {
			if err == nil {
				t.Errorf("calc(%f, %f, %q) expected error, got %f", tc.a, tc.b, tc.op, got)
			}
			continue
		}
		if err != nil {
			t.Errorf("calc(%f, %f, %q) unexpected error: %v", tc.a, tc.b, tc.op, err)
			continue
		}
		if got != tc.want {
			t.Errorf("calc(%f, %f, %q) = %f; want %f", tc.a, tc.b, tc.op, got, tc.want)
		}
	}
}

func TestDescribe(t *testing.T) {
	tests := []struct {
		v    interface{}
		want string
	}{
		{42, "int: 42"},
		{-1, "int: -1"},
		{3.14, "float: 3.140000"},
		{0.0, "float: 0.000000"},
		{"hello", "string: hello"},
		{"", "string: "},
		{true, "bool: true"},
		{false, "bool: false"},
		{[]int{1}, "unknown type: []int"},
		{nil, "unknown type: <nil>"},
	}

	for _, tc := range tests {
		got := describe(tc.v)
		if got != tc.want {
			t.Errorf("describe(%#v) = %q; want %q", tc.v, got, tc.want)
		}
	}
}

func TestCalcDivisionByZeroError(t *testing.T) {
	_, err := calc(1, 0, "div")
	if err == nil {
		t.Error("calc(1, 0, \"div\") should return error")
	}
	_, err = calc(1, 0, "/")
	if err == nil {
		t.Error("calc(1, 0, \"/\") should return error")
	}
}

func TestFallthroughDemo(t *testing.T) {
	// Verify fallthrough behaviour: case 2 falls through to case 3
	var results []string
	i := 2
	switch i {
	case 1:
		results = append(results, "one")
	case 2:
		results = append(results, "two")
		fallthrough
	case 3:
		results = append(results, "three")
	}
	if len(results) != 2 || results[0] != "two" || results[1] != "three" {
		t.Errorf("fallthrough produced %v; want [two three]", results)
	}
}
