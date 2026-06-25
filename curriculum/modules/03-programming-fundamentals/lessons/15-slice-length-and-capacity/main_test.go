package main

import "testing"

func TestGrowStepsLength(t *testing.T) {
	result := growSteps(2, 10)
	if len(result) != 10 {
		t.Errorf("growSteps(2, 10) len = %d; want 10", len(result))
	}
}

func TestGrowStepsValues(t *testing.T) {
	result := growSteps(5, 5)
	for i, v := range result {
		if v != i {
			t.Errorf("growSteps(5, 5)[%d] = %d; want %d", i, v, i)
		}
	}
}

func TestGrowStepsZeroInitial(t *testing.T) {
	result := growSteps(0, 5)
	if len(result) != 5 {
		t.Errorf("growSteps(0, 5) len = %d; want 5", len(result))
	}
	for i, v := range result {
		if v != i {
			t.Errorf("growSteps(0, 5)[%d] = %d; want %d", i, v, i)
		}
	}
}

func TestGrowStepsCapacityGrows(t *testing.T) {
	// With initial cap 2 and 20 appends, capacity must have grown
	result := growSteps(2, 20)
	if cap(result) < 20 {
		t.Errorf("growSteps(2, 20) final cap = %d; expected at least 20", cap(result))
	}
}

func TestGrowStepsNoNegative(t *testing.T) {
	s := make([]int, 0, 3)
	if len(s) != 0 || cap(s) != 3 {
		t.Fatal("setup failed")
	}
	s = append(s, 1, 2, 3)
	if len(s) != 3 || cap(s) != 3 {
		t.Fatal("expected full capacity")
	}
	s = append(s, 4)
	if cap(s) <= 3 {
		t.Error("expected capacity to grow after exceeding initial cap")
	}
}
