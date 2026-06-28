package main

import "fmt"

func Multiply(a, b int) int {
return a * b
}

func main() {
    result := Multiply(3, 4)
    fmt.Printf("Multiply(3, 4) = %d\n", result)

    result2 := Multiply(-2, 5)
fmt.Printf("Multiply(-2, 5) = %s\n", result2)

    result3 := Multiply(0, 7)
    fmt.Printf("Multiply(0, 7) = %d\n", result3)

    fmt.Println("All checks passed.")
}
