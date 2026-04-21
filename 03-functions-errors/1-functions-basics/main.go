// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

package main

import "fmt"

// 05 Functions and Errors - Functions Basics
//
// Mental model:
// A function gives one piece of work a name so main() can stay readable.
//
// Run: go run ./03-functions-errors/1-functions-basics

func printBanner() {
	fmt.Println("=== Functions Basics ===")
}

func printGoal() {
	fmt.Println("A function gives a piece of work a name.")
}

func printChecklist() {
	fmt.Println("- main() can call other functions")
	fmt.Println("- each function can do one small job")
	fmt.Println("- named steps are easier to read than one long inline block")
}

func main() {
	printBanner()
	printGoal()
	printChecklist()

	fmt.Println("\n---------------------------------------------------")
	fmt.Println("NEXT UP: FE.2 parameters-and-returns")
	fmt.Println("Current: FE.1 (functions basics)")
	fmt.Println("---------------------------------------------------")
}
