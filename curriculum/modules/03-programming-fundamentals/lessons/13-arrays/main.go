package main

import "fmt"

func sumAndAverage(nums [8]float64) (sum float64, avg float64) {
	for _, v := range nums {
		sum += v
	}
	avg = sum / float64(len(nums))
	return
}

func reverse(arr *[6]int) {
	for i, j := 0, len(arr)-1; i < j; i, j = i+1, j-1 {
		arr[i], arr[j] = arr[j], arr[i]
	}
}

func modifyAttempt(arr [6]int) {
	arr[0] = 999
}

func main() {
	nums := [8]float64{1.5, 2.5, 3.0, 4.0, 5.5, 6.0, 7.5, 8.0}
	sum, avg := sumAndAverage(nums)
	fmt.Printf("sum=%.1f avg=%.2f\n", sum, avg)

	orig := [6]int{1, 2, 3, 4, 5, 6}
	fmt.Println("before reverse:", orig)
	reverse(&orig)
	fmt.Println("after reverse:", orig)

	cp := orig
	modifyAttempt(cp)
	fmt.Println("original after modifyAttempt:", orig)
	fmt.Println("copy inside modifyAttempt was a copy:", cp)
}
