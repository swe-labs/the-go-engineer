// Package bank provides a simple bank account abstraction.
package bank

import "fmt"

type account struct {
	owner   string
	balance float64
}

func NewAccount(owner string, initialBalance float64) *account {
	return &account{owner: owner, balance: initialBalance}
}

func (a *account) Balance() float64 {
	return a.balance
}

func (a *account) Deposit(amount float64) error {
	if amount <= 0 {
		return fmt.Errorf("deposit amount must be positive")
	}
	a.balance += amount
	return nil
}

func (a *account) Withdraw(amount float64) error {
	if amount <= 0 {
		return fmt.Errorf("withdrawal amount must be positive")
	}
	if amount > a.balance {
		return fmt.Errorf("insufficient balance: have %.2f, need %.2f", a.balance, amount)
	}
	a.balance -= amount
	return nil
}
