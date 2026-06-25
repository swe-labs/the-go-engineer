package main

import (
	"encoding/json"
	"fmt"
)

type Product struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price,omitempty"`
}

type Book struct {
	Title  string `json:"title"`
	Author string `json:"author"`
	Pages  int    `json:"-"`
}

func main() {
	var p1 Product
	fmt.Printf("zero: %+v\n", p1)

	p2 := Product{ID: 1, Name: "Laptop", Price: 999.99}
	fmt.Printf("literal: %+v\n", p2)

	point := struct{ X, Y int }{X: 3, Y: 4}
	fmt.Printf("anonymous: %+v\n", point)

	a := Product{ID: 1, Name: "Mouse"}
	b := Product{ID: 1, Name: "Mouse"}
	fmt.Println("a == b:", a == b)

	jsonBytes, _ := json.Marshal(p2)
	fmt.Println("json:", string(jsonBytes))

	p3 := p2
	p3.Price = 0
	fmt.Printf("original unchanged: %.2f\n", p2.Price)

	b1 := Book{Title: "The Go Programming Language", Author: "Donovan & Kernighan", Pages: 380}
	fmt.Printf("book: %+v\n", b1)
	bookJSON, _ := json.Marshal(b1)
	fmt.Println("book json:", string(bookJSON))

	b2 := Book{Title: "The Go Programming Language", Author: "Donovan & Kernighan", Pages: 400}
	fmt.Println("b1 == b2:", b1 == b2)
}
