package main

import "fmt"

func merge(a, b []int) []int {
	result := make([]int, 0, len(a)+len(b))
	result = append(result, a...)
	result = append(result, b...)
	return result
}

func dedup(sorted []int) []int {
	if len(sorted) == 0 {
		return []int{}
	}
	result := make([]int, 0, len(sorted))
	result = append(result, sorted[0])
	for i := 1; i < len(sorted); i++ {
		if sorted[i] != sorted[i-1] {
			result = append(result, sorted[i])
		}
	}
	return result
}

func main() {
	a := []int{1, 2, 3}
	b := []int{4, 5, 6}
	fmt.Println("merge:", merge(a, b))

	dups := []int{1, 1, 2, 3, 3, 3, 4, 5, 5}
	fmt.Println("dedup:", dedup(dups))
	fmt.Println("dedup empty:", dedup([]int{}))
	fmt.Println("dedup single:", dedup([]int{7}))
}
