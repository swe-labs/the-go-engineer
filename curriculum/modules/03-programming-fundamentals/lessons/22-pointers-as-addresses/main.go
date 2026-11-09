package main

import "fmt"

func swap(a, b *int) {
	tmp := *a
	*a = *b
	*b = tmp
}

func apply(nums []int, fn func(*int)) {
	for i := range nums {
		fn(&nums[i])
	}
}

func main() {
	x, y := 1, 2
	fmt.Println("before swap: x =", x, "y =", y)
	swap(&x, &y)
	fmt.Println("after swap: x =", x, "y =", y)

	vals := []int{10, 20, 30}
	fmt.Println("before apply:", vals)
	apply(vals, func(p *int) { *p *= 2 })
	fmt.Println("after apply:", vals)
}
