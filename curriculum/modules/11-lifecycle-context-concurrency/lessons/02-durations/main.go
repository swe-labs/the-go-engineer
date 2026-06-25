package main

import (
	"fmt"
	"time"
)

func main() {
	d := 2*time.Hour + 30*time.Minute + 15*time.Second
	fmt.Println("Duration:", d)
	fmt.Println("Hours:", d.Hours())
	fmt.Println("Minutes:", d.Minutes())
	fmt.Println("Seconds:", d.Seconds())
	fmt.Println("Milliseconds:", d.Milliseconds())

	a := 3 * time.Second
	b := 1500 * time.Millisecond
	fmt.Println("a + b:", a+b)
	fmt.Println("a - b:", a-b)
	fmt.Println("a > b:", a > b)
	fmt.Println("a / b:", a/b)

	const tick = 500 * time.Millisecond
	start := time.Now()
	time.Sleep(tick)
	fmt.Println("Slept for:", time.Since(start).Round(time.Millisecond))
}
