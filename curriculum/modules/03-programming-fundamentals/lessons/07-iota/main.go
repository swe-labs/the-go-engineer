package main

import "fmt"

type Flag uint8

const (
	FlagA Flag = 1 << iota
	FlagB
	FlagC
	FlagD
)

func set(f, flag Flag) Flag {
	return f | flag
}

func clear(f, flag Flag) Flag {
	return f &^ flag
}

func has(f, flag Flag) bool {
	return f&flag != 0
}

func main() {
	var f Flag = FlagA | FlagC
	fmt.Printf("initial: %08b\n", f)

	f = set(f, FlagB)
	fmt.Printf("set(B):  %08b\n", f)

	f = clear(f, FlagC)
	fmt.Printf("clr(C):  %08b\n", f)

	fmt.Printf("has A: %t\n", has(f, FlagA))
	fmt.Printf("has B: %t\n", has(f, FlagB))
	fmt.Printf("has C: %t\n", has(f, FlagC))
	fmt.Printf("has D: %t\n", has(f, FlagD))
}
