package main

import "fmt"

type Bank struct {
	accounts map[string]int
}

func NewBank() *Bank {
	return &Bank{accounts: make(map[string]int)}
}

func (b *Bank) Deposit(name string, amount int) {
	if b == nil || b.accounts == nil {
		return
	}
	b.accounts[name] += amount
}

func (b *Bank) Withdraw(name string, amount int) bool {
	if b == nil || b.accounts == nil {
		return false
	}
	if b.accounts[name] < amount {
		return false
	}
	b.accounts[name] -= amount
	return true
}

func (b *Bank) Balance(name string) (int, bool) {
	if b == nil || b.accounts == nil {
		return 0, false
	}
	bal, ok := b.accounts[name]
	return bal, ok
}

func main() {
	bank := NewBank()

	bank.Deposit("Alice", 200)
	bank.Deposit("Bob", 100)

	fmt.Println("Alice balance:", balStr(bank.Balance("Alice")))
	fmt.Println("Bob balance:", balStr(bank.Balance("Bob")))

	bank.Withdraw("Alice", 50)
	fmt.Println("After withdrawal, Alice:", balStr(bank.Balance("Alice")))

	bank.Withdraw("Bob", 999)
	fmt.Println("Bob attempted large withdraw:", balStr(bank.Balance("Bob")))
}

func balStr(bal int, ok bool) string {
	if !ok {
		return "N/A"
	}
	return fmt.Sprintf("$%d", bal)
}
