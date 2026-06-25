package main

import "fmt"

func sumUntil(nums []int, target int) (int, int) {
	sum := 0
	count := 0
	for i := 0; i < len(nums); i++ {
		if nums[i] < 0 {
			continue
		}
		sum += nums[i]
		count++
		if sum >= target {
			break
		}
	}
	return sum, count
}

func main() {
	nums := []int{3, -1, 5, 2, -3, 4}
	target := 10
	sum, count := sumUntil(nums, target)
	fmt.Printf("nums=%v, target=%d → sum=%d, count=%d\n", nums, target, sum, count)

	fmt.Println("--- C-style ---")
	for i := 0; i < 5; i++ {
		fmt.Printf("%d ", i)
	}
	fmt.Println()

	fmt.Println("--- while-style ---")
	n := 0
	for n < 3 {
		fmt.Printf("%d ", n)
		n++
	}
	fmt.Println()

	fmt.Println("--- labeled break ---")
outer:
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if i*j > 2 {
				break outer
			}
			fmt.Printf("(%d,%d) ", i, j)
		}
	}
	fmt.Println()
}
