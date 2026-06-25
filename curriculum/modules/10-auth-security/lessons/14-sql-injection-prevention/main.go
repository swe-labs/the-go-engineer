package main

import (
	"fmt"
	"strings"
)

func unsafeQuery(queryTemplate, userInput string) string {
	return strings.ReplaceAll(queryTemplate, "'{input}'", "'"+userInput+"'")
}

func safeQuery(queryTemplate, userInput string) string {
	return strings.ReplaceAll(queryTemplate, "?", fmt.Sprintf("'%s'", escapeSQLString(userInput)))
}

func escapeSQLString(s string) string {
	r := strings.ReplaceAll(s, "'", "''")
	return r
}

func isInjection(query string) bool {
	dangerous := []string{" OR ", " UNION ", " DROP ", " -- ", "/*", "*/", ";", "1'='1"}
	upper := strings.ToUpper(query)
	for _, d := range dangerous {
		if strings.Contains(upper, d) {
			return true
		}
	}
	return false
}

func simulateDBLookup(query string) string {
	if isInjection(query) {
		return "BREACHED: All users returned!"
	}
	if strings.Contains(query, "alice") {
		return "alice@example.com"
	}
	if strings.Contains(query, "bob") {
		return "bob@example.com"
	}
	return "no results"
}

func main() {
	baseQuery := "SELECT email FROM users WHERE username = '{input}'"

	fmt.Println("=== SQL Injection Simulation ===")
	fmt.Println("Base query:", baseQuery)

	fmt.Println("\n1. Normal input:")
	normalInput := "alice"
	q1 := unsafeQuery(baseQuery, normalInput)
	fmt.Printf("  Query: %s\n", q1)
	fmt.Printf("  Result: %s\n", simulateDBLookup(q1))

	fmt.Println("\n2. Injection attempt:")
	injectionInput := "alice' OR '1'='1"
	q2 := unsafeQuery(baseQuery, injectionInput)
	fmt.Printf("  Query: %s\n", q2)
	fmt.Printf("  Result: %s\n", simulateDBLookup(q2))

	fmt.Println("\n3. Injection attempt with UNION:")
	unionInput := "nonexistent' UNION SELECT * FROM users; --"
	q3 := unsafeQuery(baseQuery, unionInput)
	fmt.Printf("  Query: %s\n", q3)
	fmt.Printf("  Injection detected: %v\n", isInjection(q3))

	fmt.Println("\n4. Parameterized query (safe):")
	safeInput := "alice' OR '1'='1"
	q4 := safeQuery("SELECT email FROM users WHERE username = ?", safeInput)
	fmt.Printf("  Query: %s\n", q4)
	fmt.Printf("  Result: %s\n", simulateDBLookup(q4))

	fmt.Println("\n5. Escaped injection (safe):")
	escapeInput := "alice' OR '1'='1"
	escaped := escapeSQLString(escapeInput)
	q5 := unsafeQuery(baseQuery, escaped)
	fmt.Printf("  Escaped input: %s\n", escaped)
	fmt.Printf("  Query: %s\n", q5)
	fmt.Printf("  Injection detected: %v\n", isInjection(q5))
}
