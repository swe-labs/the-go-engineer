package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var (
	ErrInvalidInvoice = errors.New("invalid invoice payload")
	ErrInvalidQuantity = errors.New("quantity must be positive")
	ErrInvalidPrice    = errors.New("price must be positive")
)

type CustomerTier string

const (
	TierGold     CustomerTier = "GOLD"
	TierSilver   CustomerTier = "SILVER"
	TierBronze   CustomerTier = "BRONZE"
	TierStandard CustomerTier = "STANDARD"
)

func (t CustomerTier) DiscountRate() float64 {
	switch t {
	case TierGold:
		return 0.15
	case TierSilver:
		return 0.10
	case TierBronze:
		return 0.05
	default:
		return 0.0
	}
}

type LineItem struct {
	Name       string
	Quantity   int
	PriceCents int
}

func (l LineItem) Subtotal() int {
	return l.Quantity * l.PriceCents
}

type Invoice struct {
	OrderID    string
	State      string
	Tier       CustomerTier
	Items      []LineItem
}

func (inv *Invoice) Subtotal() int {
	sum := 0
	for _, item := range inv.Items {
		sum += item.Subtotal()
	}
	return sum
}

func (inv *Invoice) DiscountCents() int {
	return int(float64(inv.Subtotal()) * inv.Tier.DiscountRate())
}

func (inv *Invoice) TaxRate() float64 {
	switch strings.ToUpper(inv.State) {
	case "CA":
		return 0.09
	case "NY":
		return 0.08
	case "TX":
		return 0.0625
	default:
		return 0.05
	}
}

func (inv *Invoice) TaxCents() int {
	taxable := inv.Subtotal() - inv.DiscountCents()
	return int(float64(taxable) * inv.TaxRate())
}

func (inv *Invoice) TotalCents() int {
	return inv.Subtotal() - inv.DiscountCents() + inv.TaxCents()
}

func ParseInvoice(raw string) (*Invoice, error) {
	lines := strings.Split(strings.TrimSpace(raw), "\n")
	if len(lines) == 0 || lines[0] == "" {
		return nil, ErrInvalidInvoice
	}

	headerParts := strings.Split(strings.TrimSpace(lines[0]), ",")
	if len(headerParts) < 3 {
		return nil, fmt.Errorf("%w: missing header fields (expected order_id,state,tier)", ErrInvalidInvoice)
	}

	orderID := strings.TrimSpace(headerParts[0])
	state := strings.TrimSpace(headerParts[1])
	tier := CustomerTier(strings.ToUpper(strings.TrimSpace(headerParts[2])))

	invoice := &Invoice{
		OrderID: orderID,
		State:   state,
		Tier:    tier,
		Items:   make([]LineItem, 0, len(lines)-1),
	}

	for i := 1; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if line == "" {
			continue
		}

		parts := strings.Split(line, ",")
		if len(parts) < 3 {
			return nil, fmt.Errorf("%w: invalid item format on line %d", ErrInvalidInvoice, i+1)
		}

		name := strings.TrimSpace(parts[0])
		qty, err := strconv.Atoi(strings.TrimSpace(parts[1]))
		if err != nil || qty <= 0 {
			return nil, fmt.Errorf("%w on line %d", ErrInvalidQuantity, i+1)
		}

		price, err := strconv.Atoi(strings.TrimSpace(parts[2]))
		if err != nil || price <= 0 {
			return nil, fmt.Errorf("%w on line %d", ErrInvalidPrice, i+1)
		}

		invoice.Items = append(invoice.Items, LineItem{
			Name:       name,
			Quantity:   qty,
			PriceCents: price,
		})
	}

	if len(invoice.Items) == 0 {
		return nil, fmt.Errorf("%w: invoice must have at least one line item", ErrInvalidInvoice)
	}

	return invoice, nil
}

// ProcessInvoice maintains backwards-compatibility with the legacy API while delegating to clean domain models.
func ProcessInvoice(raw string) (string, int, error) {
	invoice, err := ParseInvoice(raw)
	if err != nil {
		return "", 0, err
	}

	total := invoice.TotalCents()
	summary := fmt.Sprintf("Order %s: Subtotal $%0.2f, Discount -$%0.2f, Tax $%0.2f, Total $%0.2f",
		invoice.OrderID,
		float64(invoice.Subtotal())/100,
		float64(invoice.DiscountCents())/100,
		float64(invoice.TaxCents())/100,
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
