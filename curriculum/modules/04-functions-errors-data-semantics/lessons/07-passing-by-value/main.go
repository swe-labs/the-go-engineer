package main

import "fmt"

func main() {
	original := []int{1, 2, 3, 4, 5}
	rotated := rotate(original, 2)
	fmt.Println("original:", original)
	fmt.Println("rotated:", rotated)
}

func rotate(nums []int, n int) []int {
	if len(nums) == 0 {
		return []int{}
	}
	n = n % len(nums)
	result := make([]int, len(nums))
	for i, v := range nums {
		newPos := (i - n + len(nums)) % len(nums)
		result[newPos] = v
	}
	return result
}
