package main

import "fmt"

// Demo demonstrates go generate.
func Demo() string {
	return "go generate verified"
}

func main() {
	fmt.Println(Demo())
}
