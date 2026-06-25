package main

import "fmt"

func growSteps(initialCap, numAppends int) []int {
	s := make([]int, 0, initialCap)
	oldCap := cap(s)
	for i := 0; i < numAppends; i++ {
		s = append(s, i)
		if cap(s) != oldCap {
			fmt.Printf("capacity grew: %d -> %d (len=%d)\n", oldCap, cap(s), len(s))
			oldCap = cap(s)
		}
	}
	return s
}

func main() {
	result := growSteps(2, 10)
	fmt.Println("final slice:", result)
	fmt.Println("final len:", len(result), "cap:", cap(result))
}
