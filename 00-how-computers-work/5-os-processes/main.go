package main

import (
	"fmt"
	"os"
)

// Section 00: How Computers Work

func main() {
	// 1. Ask the Operating System for our assigned Process ID
	pid := os.Getpid()

	// 2. Print it to the screen
	fmt.Printf("Hello OS! I am running as Process ID: %d\n", pid)

	// 3. Manually tell the OS that we finished successfully (Code 0)
	// Even if we don't write this, Go does it for us at the end of main()
	fmt.Println("NEXT UP: GT.1 installation")
	os.Exit(0)
}
