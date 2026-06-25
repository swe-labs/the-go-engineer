package main

import "testing"

func TestSumUntil(t *testing.T) {
	got := SumUntil(5)
	want := 15
	if got != want {
		t.Errorf("SumUntil(5) = %d; want %d", got, want)
	}
}

func TestAverage(t *testing.T) {
	got := Average([]int{2, 4, 6})
	want := 4.0
	if got != want {
		t.Errorf("Average([2,4,6]) = %f; want %f", got, want)
	}
}

func TestAverageEmpty(t *testing.T) {
	got := Average([]int{})
	want := 0.0
	if got != want {
		t.Errorf("Average([]) = %f; want %f", got, want)
	}
}

func TestSumEvens(t *testing.T) {
	got := SumEvens([]int{1, 2, 3, 4, 5})
	want := 6
	if got != want {
		t.Errorf("SumEvens([1,2,3,4,5]) = %d; want %d", got, want)
	}
}

func TestSumEvensEmpty(t *testing.T) {
	got := SumEvens([]int{})
	want := 0
	if got != want {
		t.Errorf("SumEvens([]) = %d; want %d", got, want)
	}
}

func TestSumEvensAllEvens(t *testing.T) {
	got := SumEvens([]int{2, 4, 6, 8})
	want := 20
	if got != want {
		t.Errorf("SumEvens([2,4,6,8]) = %d; want %d", got, want)
	}
}

func TestSumEvensAllOdds(t *testing.T) {
	got := SumEvens([]int{1, 3, 5})
	want := 0
	if got != want {
		t.Errorf("SumEvens([1,3,5]) = %d; want %d", got, want)
	}
}
