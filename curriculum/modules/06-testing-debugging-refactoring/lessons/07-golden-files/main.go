package main

import (
	"fmt"
	"strings"
)

func FormatPerson(name string, age int) string {
	return fmt.Sprintf("Name: %s\nAge: %d\n", name, age)
}

func GenerateReport(items []string) string {
	if len(items) == 0 {
		return "No items\n"
	}
	var b strings.Builder
	for i, item := range items {
		b.WriteString(fmt.Sprintf("%d. %s\n", i+1, item))
	}
	return b.String()
}

func main() {
	fmt.Print(FormatPerson("Alice", 30))
	fmt.Print(GenerateReport([]string{"apples", "bananas", "cherries"}))
}
