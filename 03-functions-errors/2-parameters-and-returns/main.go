// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

package main

import "fmt"

// 05 Functions and Errors - Parameters and Returns
//
// Mental model:
// Parameters bring values into a function, and return values send results back.
//
// Run: go run ./03-functions-errors/2-parameters-and-returns

func announceCart(name string) {
	fmt.Printf("Checking %s\n", name)
}

func sumPrices(prices []int) int {
	total := 0

	for _, price := range prices {
		total += price
	}

	return total
}

func labelPrice(name string, total int) string {
	return fmt.Sprintf("%s total: %d", name, total)
}

func main() {
	prices := []int{12, 18, 25}

	announceCart("starter cart")

	total := sumPrices(prices)
	summary := labelPrice("starter cart", total)

	fmt.Println(summary)

	fmt.Println("\n---------------------------------------------------")
	fmt.Println("NEXT UP: FE.3 multiple-return-values")
	fmt.Println("Current: FE.2 (parameters and returns)")
	fmt.Println("---------------------------------------------------")
}
