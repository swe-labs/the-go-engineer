package main

import (
	"fmt"
	"math"
)

func main() {
	fmt.Println("--- zero value ---")
	var count int
	fmt.Println("count:", count)

	fmt.Println("--- var with init ---")
	var msg string = "hello"
	fmt.Println("msg:", msg)

	fmt.Println("--- short declaration ---")
	price := 42
	fmt.Println("price:", price)

	fmt.Println("--- multiple assignment / swap ---")
	x, y := 10, 20
	fmt.Println("before swap:", x, y)
	x, y = y, x
	fmt.Println("after swap:", x, y)

	fmt.Println("--- shadowing ---")
	shadow := "outer"
	{
		shadow := "inner"
		fmt.Println("inside block:", shadow)
	}
	fmt.Println("outside block:", shadow)

	fmt.Println("--- cylinder calculations ---")
	area, volume := calculate(3, 5)
	fmt.Printf("r=3, h=5: area=%.2f, volume=%.2f\n", area, volume)
	area, volume = calculate(1.5, 2)
	fmt.Printf("r=1.5, h=2: area=%.2f, volume=%.2f\n", area, volume)
	var r, h float64
	r, h = 10, 1
	area, volume = calculate(r, h)
	fmt.Printf("r=10, h=1: area=%.2f, volume=%.2f\n", area, volume)

	fmt.Println("--- blank identifier ---")
	total := 100
	half, _ := divide(total, 2)
	fmt.Println("half of 100:", half)
}

func calculate(r, h float64) (area, volume float64) {
	area = 2 * math.Pi * r * (r + h)
	volume = math.Pi * r * r * h
	return
}

func divide(a, b int) (int, int) {
	return a / b, a % b
}
