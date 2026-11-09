package mathutil

import "testing"

func TestAdd(t *testing.T) {
	if got := Add(2, 3); got != 5 {
		t.Errorf("Add(2,3) = %d; want 5", got)
	}
	if got := Add(-1, 1); got != 0 {
		t.Errorf("Add(-1,1) = %d; want 0", got)
	}
}

func TestMul(t *testing.T) {
	if got := Mul(3, 4); got != 12 {
		t.Errorf("Mul(3,4) = %d; want 12", got)
	}
	if got := Mul(0, 100); got != 0 {
		t.Errorf("Mul(0,100) = %d; want 0", got)
	}
}
