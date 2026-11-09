package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

type Product struct {
	SKU   string  `json:"sku"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

func WriteProducts(products []Product, w io.Writer) error {
	enc := json.NewEncoder(w)
	for _, p := range products {
		if err := enc.Encode(p); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	products := []Product{
		{"A100", "Widget", 9.99},
		{"A200", "Gadget", 24.99},
		{"A300", "Doohickey", 4.99},
	}

	var buf bytes.Buffer
	if err := WriteProducts(products, &buf); err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Print(buf.String())
}
