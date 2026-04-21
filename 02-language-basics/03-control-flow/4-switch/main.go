// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

package main

import "fmt"

func main() {
	day := "Monday"

	switch day {
	case "Saturday", "Sunday":
		fmt.Println("Weekend mode.")
	case "Monday":
		fmt.Println("Start-of-week mode.")
	default:
		fmt.Println("Regular workday mode.")
	}

	score := 82

	switch {
	case score >= 90:
		fmt.Println("Excellent result.")
	case score >= 80:
		fmt.Println("Strong result.")
	case score >= 70:
		fmt.Println("Passing result.")
	default:
		fmt.Println("Needs more work.")
	}

	fmt.Println("\n---------------------------------------------------")
	fmt.Println("NEXT UP: CF.5 defer-basics")
	fmt.Println("Current: CF.4 (switch)")
	fmt.Println("---------------------------------------------------")
}
