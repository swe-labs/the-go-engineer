// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

package main

import (
	"fmt"
	"math"
	"strings"
)

func main() {
	greeting := "hello, go developer!"
	upper := strings.ToUpper(greeting)
	hasGo := strings.Contains(greeting, "go")
	parts := strings.Split("one,two,three", ",")

	fmt.Println("Original:   ", greeting)
	fmt.Println("Uppercase:  ", upper)
	fmt.Println("Contains go:", hasGo)
	fmt.Println("Split parts:", parts)
	fmt.Printf("Pi:         %.2f\n", math.Pi)
	fmt.Printf("Sqrt(144):  %.0f\n", math.Sqrt(144))
	fmt.Printf("2^10:       %.0f\n", math.Pow(2, 10))
	fmt.Println()
	fmt.Println("---------------------------------------------------")
	fmt.Println("NEXT UP: GT.4 dev-environment")
	fmt.Println("Current: GT.3 (how-go-works)")
	fmt.Println("---------------------------------------------------")
}
