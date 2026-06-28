package main

import (
	"fmt"
)

// Item represents a product for sale.
type Item struct {
	Name       string
	PriceCents int
}

// NewItem creates a new Item with the given name and price in cents.
func NewItem(name string, priceCents int) Item {
	// TODO: implement
	return Item{}
}

// TaxRate returns the tax rate for a given country code.
// US -> 0.07, GB -> 0.20, DE -> 0.19, JP -> 0.10.
// Unknown country returns 0.0.
func TaxRate(country string) float64 {
	// TODO: implement
	return 0
}

// DiscountPercent returns the discount rate based on quantity.
// >= 100 -> 0.15, >= 50 -> 0.10, >= 10 -> 0.05, otherwise 0.0.
func DiscountPercent(quantity int) float64 {
	// TODO: implement
	return 0
}

// FormatPrice converts an amount in cents to a dollar string.
// 1234 -> "$12.34", -1234 -> "-$12.34".
func FormatPrice(cents int) string {
	// TODO: implement
	return ""
}

// CalculateTotal sums all item prices, applies tax for the given country,
// applies a quantity discount, and returns the total in cents.
// Returns an error if the items slice is empty.
func CalculateTotal(items []Item, country string) (int, error) {
	// TODO: implement
	return 0, nil
}

func main() {
	items := []Item{
		NewItem("Laptop", 99999),
		NewItem("Mouse", 2499),
		NewItem("Keyboard", 8999),
	}

	for _, item := range items {
		fmt.Printf("Item: %s - %s\n", item.Name, FormatPrice(item.PriceCents))
	}

	total, err := CalculateTotal(items, "US")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Printf("Total: %s\n", FormatPrice(total))
}
