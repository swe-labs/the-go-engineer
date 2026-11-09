package main

import (
	"errors"
	"fmt"
)

type Item struct {
	Name  string
	Price int
}

var catalog = map[string]int{
	"item1": 100,
	"item2": 100,
}

func main() {
	items := []string{"item1", "item2"}
	cost, err := processOrder(items)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Total cost:", cost)
	}
}

func processOrder(items []string) (int, error) {
	if len(items) == 0 {
		return 0, errors.New("empty order")
	}
	total := 0
	for _, name := range items {
		price, ok := catalog[name]
		if !ok {
			return 0, errors.New("unknown item: " + name)
		}
		total += price
	}
	return total, nil
}
