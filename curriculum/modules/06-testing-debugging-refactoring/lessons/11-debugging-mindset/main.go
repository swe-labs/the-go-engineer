package main

import "fmt"

func SumUntil(n int) int {
	sum := 0
	for i := 1; i <= n; i++ {
		sum += i
	}
	return sum
}

func Average(nums []int) float64 {
	if len(nums) == 0 {
		return 0
	}
	sum := 0
	for _, n := range nums {
		sum += n
	}
	return float64(sum) / float64(len(nums))
}

func SumEvens(nums []int) int {
	sum := 0
	for _, n := range nums {
		if n%2 == 0 {
			sum += n
		}
	}
	return sum
}

func main() {
	fmt.Println("SumUntil(5) =", SumUntil(5))
	fmt.Println("Average([2,4,6]) =", Average([]int{2, 4, 6}))
	fmt.Println("Average([]) =", Average([]int{}))
	fmt.Println("SumEvens([1,2,3,4,5]) =", SumEvens([]int{1, 2, 3, 4, 5}))
}
