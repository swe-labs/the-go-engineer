// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

package main

import "fmt"

func main() {
	fmt.Println("Odd numbers until the stop point:")

	for i := 1; i <= 10; i++ {
		if i%2 == 0 {
			continue
		}

		if i == 7 {
			break
		}

		fmt.Println(i)
	}

	fmt.Println("\n---------------------------------------------------")
	fmt.Println("NEXT UP: CF.4 switch")
	fmt.Println("Current: CF.3 (break / continue)")
	fmt.Println("---------------------------------------------------")
}
