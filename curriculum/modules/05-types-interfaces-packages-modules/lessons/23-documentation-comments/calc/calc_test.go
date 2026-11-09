package calc

import (
	"fmt"
	"testing"
)

func ExampleSum() {
	total := Sum(1, 2, 3, 4)
	fmt.Println(total)
	// Output: 10
}

func ExampleCalculator() {
	c := NewCalculator()
	c.Add(5)
	c.Add(10)
	fmt.Println(c.Result())
	// Output: 15
}

func TestSum(t *testing.T) {
	if got := Sum(1, 2); got != 3 {
		t.Errorf("Sum(1,2) = %d; want 3", got)
	}
	if got := Sum(); got != 0 {
		t.Errorf("Sum() = %d; want 0", got)
	}
}

func TestCalculator(t *testing.T) {
	c := NewCalculator()
	if got := c.Result(); got != 0 {
		t.Errorf("initial result = %d; want 0", got)
	}
	c.Add(10)
	c.Add(20)
	if got := c.Result(); got != 30 {
		t.Errorf("after add = %d; want 30", got)
	}
}
