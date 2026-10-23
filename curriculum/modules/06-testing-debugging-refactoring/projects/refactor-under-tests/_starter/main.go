package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var (
	ErrInvalidInvoice  = errors.New("invalid invoice payload")
	ErrInvalidQuantity = errors.New("quantity must be positive")
	ErrInvalidPrice    = errors.New("price must be positive")
)

// Legacy un-refactored monolithic function.
// TODO: Refactor this into clean domain types (Invoice, LineItem, CustomerTier) with separate
// parsing, calculation, and formatting methods, while ensuring ProcessInvoice passes all characterization tests.

func ProcessInvoice(raw string) (string, int, error) {
	lines := strings.Split(strings.TrimSpace(raw), "\n")
	if len(lines) == 0 || lines[0] == "" {
		return "", 0, ErrInvalidInvoice
	}

	headerParts := strings.Split(strings.TrimSpace(lines[0]), ",")
	if len(headerParts) < 3 {
		return "", 0, fmt.Errorf("%w: missing header fields", ErrInvalidInvoice)
	}

	orderID := strings.TrimSpace(headerParts[0])
	state := strings.ToUpper(strings.TrimSpace(headerParts[1]))
	tier := strings.ToUpper(strings.TrimSpace(headerParts[2]))

	subtotal := 0
	itemCount := 0

	for i := 1; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if line == "" {
			continue
		}

		parts := strings.Split(line, ",")
		if len(parts) < 3 {
			return "", 0, fmt.Errorf("%w: invalid line format", ErrInvalidInvoice)
		}

		qty, err := strconv.Atoi(strings.TrimSpace(parts[1]))
		if err != nil || qty <= 0 {
			return "", 0, fmt.Errorf("%w on line %d", ErrInvalidQuantity, i+1)
		}

		price, err := strconv.Atoi(strings.TrimSpace(parts[2]))
		if err != nil || price <= 0 {
			return "", 0, fmt.Errorf("%w on line %d", ErrInvalidPrice, i+1)
		}

		subtotal += qty * price
		itemCount++
	}

	if itemCount == 0 {
		return "", 0, fmt.Errorf("%w: no items found", ErrInvalidInvoice)
	}

	discountRate := 0.0
	if tier == "GOLD" {
		discountRate = 0.15
	} else if tier == "SILVER" {
		discountRate = 0.10
	} else if tier == "BRONZE" {
		discountRate = 0.05
	}

	discountCents := int(float64(subtotal) * discountRate)
	taxable := subtotal - discountCents

	taxRate := 0.05
	if state == "CA" {
		taxRate = 0.09
	} else if state == "NY" {
		taxRate = 0.08
	} else if state == "TX" {
		taxRate = 0.0625
	}

	taxCents := int(float64(taxable) * taxRate)
	total := taxable + taxCents

	summary := fmt.Sprintf("Order %s: Subtotal $%0.2f, Discount -$%0.2f, Tax $%0.2f, Total $%0.2f",
		orderID,
		float64(subtotal)/100,
		float64(discountCents)/100,
		float64(taxCents)/100,
		float64(total)/100,
	)

	return summary, total, nil
}

func main() {
	sampleCSV := `ORD-987,CA,GOLD
Widget Pro,2,5000
Gadget Basic,1,3000`

	summary, totalCents, err := ProcessInvoice(sampleCSV)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Println(summary)
	fmt.Printf("Calculated total: %d cents\n", totalCents)
}
