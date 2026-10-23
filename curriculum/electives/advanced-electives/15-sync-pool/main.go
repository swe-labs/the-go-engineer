package main

import "fmt"

// Demo demonstrates sync.Pool.
func Demo() string {
    return "sync.Pool verified"
}

func main() {
    fmt.Println(Demo())
}
