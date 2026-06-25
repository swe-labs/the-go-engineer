package main

import (
	"errors"
	"fmt"
)

type Account struct {
	ID      string
	Owner   string
	Balance float64
}

func NewAccount(id, owner string, initialBalance float64) (*Account, error) {
	if id == "" {
		return nil, errors.New("account ID is required")
	}
	if owner == "" {
		return nil, errors.New("owner is required")
	}
	if initialBalance < 0 {
		return nil, errors.New("initial balance cannot be negative")
	}
	return &Account{ID: id, Owner: owner, Balance: initialBalance}, nil
}

func (a *Account) Deposit(amount float64) error {
	if amount <= 0 {
		return errors.New("deposit amount must be positive")
	}
	a.Balance += amount
	return nil
}

func (a *Account) Withdraw(amount float64) error {
	if amount <= 0 {
		return errors.New("withdrawal amount must be positive")
	}
	if amount > a.Balance {
		return fmt.Errorf("insufficient balance: have %.2f, need %.2f", a.Balance, amount)
	}
	a.Balance -= amount
	return nil
}

func Transfer(from, to *Account, amount float64) error {
	if from == nil || to == nil {
		return errors.New("both accounts must be provided")
	}
	if from.ID == to.ID {
		return errors.New("cannot transfer to the same account")
	}
	if amount <= 0 {
		return errors.New("transfer amount must be positive")
	}
	if err := from.Withdraw(amount); err != nil {
		return fmt.Errorf("withdrawal failed: %w", err)
	}
	if err := to.Deposit(amount); err != nil {
		from.Balance += amount
		return fmt.Errorf("deposit failed, rollback: %w", err)
	}
	return nil
}

func main() {
	alice, _ := NewAccount("acc-1", "Alice", 500.00)
	bob, _ := NewAccount("acc-2", "Bob", 100.00)

	fmt.Printf("Before: Alice=%.2f, Bob=%.2f\n", alice.Balance, bob.Balance)
	if err := Transfer(alice, bob, 200.00); err != nil {
		fmt.Println("Transfer failed:", err)
		return
	}
	fmt.Printf("After:  Alice=%.2f, Bob=%.2f\n", alice.Balance, bob.Balance)
}
