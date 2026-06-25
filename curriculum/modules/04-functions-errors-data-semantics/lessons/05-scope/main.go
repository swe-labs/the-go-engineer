package main

import "fmt"

var name = "global"

func main() {
	fmt.Println("demoScope result:", demoScope())
}

func demoScope() string {
	globalCopy := name
	name := "local"
	localCopy := name

	if true {
		name := "block"
		blockCopy := name

		// Return in order: block, local, global
		return blockCopy + " " + localCopy + " " + globalCopy
	}

	// unreachable but keeps return type satisfied
	return ""
}
