// Copyright (c) 2026 Rasel Hossen
// See LICENSE for usage terms.

package main

import "fmt"

func main() {
	cart := []string{"TSHIRT", "MUG", "HAT", "BOOK", "KEYBOARD"}

	var subtotal float64

	fmt.Println("Processing checkout:")

	for _, item := range cart {
		var price float64

		switch item {
		case "TSHIRT":
			price = 20.00
		case "MUG":
			price = 12.50
		case "HAT":
			price = 18.00
		case "BOOK":
			price = 25.99
		}

		if price == 0 {
			fmt.Printf("skip %s: unknown item\n", item)
			continue
		}

		if item == "BOOK" {
			originalPrice := price
			price = price * 0.90
			fmt.Printf("%s promo: %.2f -> %.2f\n", item, originalPrice, price)
		} else {
			fmt.Printf("%s: %.2f\n", item, price)
		}

		subtotal += price
	}

	fmt.Printf("subtotal: %.2f\n", subtotal)
}
