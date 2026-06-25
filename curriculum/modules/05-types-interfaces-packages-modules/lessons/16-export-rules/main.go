package main

import (
	"fmt"
	"log"

	"github.com/rasel9t6/the-go-engineer/curriculum/modules/05-types-interfaces-packages-modules/lessons/16-export-rules/bank"
)

func main() {
	acc := bank.NewAccount("Alice", 100.0)
	fmt.Printf("Balance: %.2f\n", acc.Balance())

	if err := acc.Deposit(50.0); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("After deposit: %.2f\n", acc.Balance())

	if err := acc.Withdraw(30.0); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("After withdrawal: %.2f\n", acc.Balance())

	if err := acc.Withdraw(200.0); err != nil {
		fmt.Println("Expected error:", err)
	}
}
