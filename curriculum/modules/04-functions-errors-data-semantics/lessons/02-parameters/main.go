package main

import "fmt"

func main() {
	fmt.Println("max() =", max())
	fmt.Println("max(3, 7, 2) =", max(3, 7, 2))
	fmt.Println("max(-5, -1, -10) =", max(-5, -1, -10))
	fmt.Println("concat(\", \", \"a\", \"b\", \"c\") =", concat(", ", "a", "b", "c"))
	fmt.Println("concat(\"-\", \"x\") =", concat("-", "x"))
	fmt.Println("concat(\" \") =", concat(" "))
}

func max(nums ...int) int {
	if len(nums) == 0 {
		return 0
	}
	m := nums[0]
	for _, n := range nums[1:] {
		if n > m {
			m = n
		}
	}
	return m
}

func concat(sep string, parts ...string) string {
	result := ""
	for i, p := range parts {
		if i > 0 {
			result += sep
		}
		result += p
	}
	return result
}
