package main

import "testing"

func TestAdd(t *testing.T) {
	tests := []struct {
		name string
		a    int
		b    int
		want int
	}{
		{"positive integers", 2, 3, 5},
		{"negative integers", -1, -2, -3},
		{"mixed signs", -5, 10, 5},
		{"zeros", 0, 0, 0},
		{"identity with zero", 42, 0, 42},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Add(tt.a, tt.b)
			if got != tt.want {
				t.Errorf("Add(%d, %d) = %d; want %d", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func BenchmarkAdd(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = Add(i, 42)
	}
}

func TestMainExecution(t *testing.T) {
	main()
}
