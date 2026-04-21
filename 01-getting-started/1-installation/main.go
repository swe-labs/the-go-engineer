// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

package main

import (
	"fmt"
	"runtime"
)

func main() {
	fmt.Println("Go installation looks healthy.")
	fmt.Printf("Go version:   %s\n", runtime.Version())
	fmt.Printf("OS/Arch:      %s/%s\n", runtime.GOOS, runtime.GOARCH)
	fmt.Printf("Logical CPUs: %d\n", runtime.NumCPU())
	fmt.Println()
	fmt.Println("---------------------------------------------------")
	fmt.Println("NEXT UP: GT.2 hello-world")
	fmt.Println("Current: GT.1 (installation)")
	fmt.Println("---------------------------------------------------")
}
