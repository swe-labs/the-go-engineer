package main

import "fmt"

// Demo demonstrates gRPC interceptors.
func Demo() string {
    return "gRPC interceptors verified"
}

func main() {
    fmt.Println(Demo())
}
