package main

import "fmt"

// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// Section 06: Backend, APIs & Databases - REST design principles
//
// Run: go run ./06-backend-db/01-web-and-database/apis/1-rest-design-principles

func main() {
	fmt.Println("=== API.1 REST design principles ===")
	fmt.Println("Learn the resource-oriented rules that make REST APIs predictable for humans and clients.")
	fmt.Println()
	fmt.Println("- Name resources as nouns, not RPC verbs.")
	fmt.Println("- Let HTTP verbs express the action when possible.")
	fmt.Println("- Prefer predictable payload and error shapes over one-off convenience.")
	fmt.Println()
	fmt.Println("API consistency saves more engineering time than clever endpoint naming ever will.")
	fmt.Println()
	fmt.Println("---------------------------------------------------")
	fmt.Println("NEXT UP: API.2")
	fmt.Println("Current: API.1 (rest design principles)")
	fmt.Println("---------------------------------------------------")
}
