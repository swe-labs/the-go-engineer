package main

import "testing"

func TestEvaluate(t *testing.T) {
	got := evaluate(10, 5, '+')
	if got != 15 {
		t.Errorf("evaluate(10, 5, '+') = %d; want 15", got)
	}
	got = evaluate(10, 5, '-')
	if got != 5 {
		t.Errorf("evaluate(10, 5, '-') = %d; want 5", got)
	}
	got = evaluate(10, 5, '*')
	if got != 50 {
		t.Errorf("evaluate(10, 5, '*') = %d; want 50", got)
	}
	got = evaluate(10, 5, '/')
	if got != 2 {
		t.Errorf("evaluate(10, 5, '/') = %d; want 2", got)
	}
	got = evaluate(10, 5, '%')
	if got != 0 {
		t.Errorf("evaluate(10, 5, '%%') = %d; want 0", got)
	}
	got = evaluate(7, 3, '+')
	if got != 10 {
		t.Errorf("evaluate(7, 3, '+') = %d; want 10", got)
	}
	got = evaluate(7, 3, '-')
	if got != 4 {
		t.Errorf("evaluate(7, 3, '-') = %d; want 4", got)
	}
	got = evaluate(7, 3, '*')
	if got != 21 {
		t.Errorf("evaluate(7, 3, '*') = %d; want 21", got)
	}
	got = evaluate(7, 3, '/')
	if got != 2 {
		t.Errorf("evaluate(7, 3, '/') = %d; want 2", got)
	}
	got = evaluate(7, 3, '%')
	if got != 1 {
		t.Errorf("evaluate(7, 3, '%%') = %d; want 1", got)
	}
}
