package main

import (
	"errors"
	"fmt"
)

type Item struct {
	Name       string
	PriceCents int
}

func NewItem(name string, priceCents int) Item {
	return Item{Name: name, PriceCents: priceCents}
}

func TaxRate(country string) float64 {
	switch country {
	case "US":
		return 0.07
	case "GB":
		return 0.20
	case "DE":
		return 0.19
	case "JP":
		return 0.10
	default:
		return 0.0
	}
}

func DiscountPercent(quantity int) float64 {
	switch {
	case quantity >= 100:
		return 0.15
	case quantity >= 50:
		return 0.10
	case quantity >= 10:
		return 0.05
	default:
		return 0.0
	}
}

func FormatPrice(cents int) string {
	if cents < 0 {
		return fmt.Sprintf("-$%d.%02d", -cents/100, -cents%100)
	}
	return fmt.Sprintf("$%d.%02d", cents/100, cents%100)
}

func CalculateTotal(items []Item, country string) (int, error) {
	if len(items) == 0 {
		return 0, errors.New("no items to calculate")
	}

	sum := 0
	for _, item := range items {
		sum += item.PriceCents
	}

	quantity := len(items)
	discount := DiscountPercent(quantity)
	if discount > 0 {
		sum -= int(float64(sum) * discount)
	}

	tax := TaxRate(country)
	if tax > 0 {
		sum += int(float64(sum) * tax)
	}

	return sum, nil
}

func main() {
	items := []Item{
		NewItem("Laptop", 99999),
		NewItem("Mouse", 2499),
		NewItem("Keyboard", 8999),
	}

	subtotal := 0
	for _, item := range items {
		subtotal += item.PriceCents
		fmt.Printf("Item: %s - %s\n", item.Name, FormatPrice(item.PriceCents))
	}

	fmt.Printf("Subtotal: %s\n", FormatPrice(subtotal))

	qty := len(items)
	discount := DiscountPercent(qty)
	discountAmount := int(float64(subtotal) * discount)
	fmt.Printf("Discount (%d%%): %s\n", int(discount*100), FormatPrice(discountAmount))

	afterDiscount := subtotal - discountAmount
	tax := TaxRate("US")
	taxAmount := int(float64(afterDiscount) * tax)
	fmt.Printf("Tax (US %d%%): %s\n", int(tax*100), FormatPrice(taxAmount))

	total, err := CalculateTotal(items, "US")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Printf("Total: %s\n", FormatPrice(total))
}
