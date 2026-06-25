package main

import "fmt"

func makeIncrementer() func() int {
	counter := 0
	return func() int {
		counter++
		return counter
	}
}

func captureCorrect(items []int) []*int {
	var ptrs []*int
	for _, v := range items {
		v := v
		ptrs = append(ptrs, &v)
	}
	return ptrs
}

func captureBuggy(items []int) []*int {
	var ptrs []*int
	for _, v := range items {
		ptrs = append(ptrs, &v)
	}
	return ptrs
}

func main() {
	inc := makeIncrementer()
	fmt.Println("inc():", inc())
	fmt.Println("inc():", inc())
	fmt.Println("inc():", inc())

	items := []int{10, 20, 30}
	fmt.Println("\nBuggy capture:")
	buggy := captureBuggy(items)
	for _, p := range buggy {
		fmt.Println(*p)
	}

	fmt.Println("\nCorrect capture:")
	correct := captureCorrect(items)
	for _, p := range correct {
		fmt.Println(*p)
	}
}
