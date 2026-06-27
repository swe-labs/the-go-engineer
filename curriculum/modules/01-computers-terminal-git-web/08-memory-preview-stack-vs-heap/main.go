package main

import "fmt"

func StackValue() int {
	x := 42
	return x
}

func AddNumbers(a, b int) int {
	sum := a + b
	return sum
}

func SumArray(nums [5]int) int {
	sum := 0
	for _, v := range nums {
		sum += v
	}
	return sum
}

func main() {
	s := StackValue()
	fmt.Printf("Stack value: %d\n", s)

	sum := AddNumbers(10, 20)
	fmt.Printf("Sum: %d\n", sum)

	arr := [5]int{0, 1, 2, 3, 4}
	total := SumArray(arr)
	fmt.Printf("Array sum: %d\n", total)
}
