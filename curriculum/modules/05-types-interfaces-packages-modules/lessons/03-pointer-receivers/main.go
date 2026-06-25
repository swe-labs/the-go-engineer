package main

import (
	"errors"
	"fmt"
)

type Counter struct {
	value int
}

func (c *Counter) Increment() {
	c.value++
}

func (c *Counter) Add(n int) {
	c.value += n
}

func (c *Counter) Value() int {
	if c == nil {
		return 0
	}
	return c.value
}

func (c *Counter) Reset() {
	c.value = 0
}

type BankAccount struct {
	balance float64
}

func (a *BankAccount) Deposit(amount float64) {
	a.balance += amount
}

func (a *BankAccount) Withdraw(amount float64) error {
	if amount > a.balance {
		return errors.New("insufficient funds")
	}
	a.balance -= amount
	return nil
}

func (a *BankAccount) Balance() float64 {
	return a.balance
}

func main() {
	c := Counter{}
	c.Increment()
	c.Increment()
	c.Add(3)
	fmt.Println("value:", c.Value())

	c.Reset()
	fmt.Println("after reset:", c.Value())

	var nilCounter *Counter
	fmt.Println("nil receiver:", nilCounter.Value())

	account := BankAccount{}
	account.Deposit(100)
	err := account.Withdraw(30)
	if err != nil {
		fmt.Println("withdraw error:", err)
	}
	fmt.Println("balance:", account.Balance())

	err = account.Withdraw(100)
	if err != nil {
		fmt.Println("withdraw error:", err)
	}
	fmt.Println("balance after failed withdraw:", account.Balance())
}
