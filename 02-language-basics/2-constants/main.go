// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

package main

import "fmt"

const (
	Host = "127.0.0.1"
	Port = ":8080"
	User = "root"
)

var (
	isRunning bool
)

func main() {
	AppName := "Go"
	fmt.Println(AppName)

	const pi float64 = 3.1415926
	fmt.Println(pi)

	const rate float32 = 5.2
	fmt.Println(rate)

	fmt.Printf("Server: %s%s (User: %s)\n", Host, Port, User)
	fmt.Printf("Running: %t\n", isRunning)

	fmt.Println("---------------------------------------------------")
	fmt.Println("NEXT UP: LB.3 enums")
	fmt.Println("Current: LB.2 (constants)")
	fmt.Println("---------------------------------------------------")
}
