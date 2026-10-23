package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

var (
	ErrNegativeAmount      = errors.New("amount must be strictly positive")
	ErrInsufficientFunds   = errors.New("insufficient funds for withdrawal")
	ErrAccountFrozen       = errors.New("account is frozen")
	ErrBufferEmpty         = errors.New("buffer is empty")
	ErrBelowMinimumBalance = errors.New("withdrawal breaches minimum balance")
)

// TransactionType denotes the nature of an account transaction.
type TransactionType string

const (
	TxDeposit  TransactionType = "DEPOSIT"
	TxWithdraw TransactionType = "WITHDRAW"
	TxInterest TransactionType = "INTEREST"
)

// Transaction records a single ledger movement.
type Transaction struct {
	ID        string          `json:"id"`
	Type      TransactionType `json:"type"`
	Amount    int             `json:"amount_cents"`
	Balance   int             `json:"balance_after"`
	Timestamp time.Time       `json:"timestamp"`
}

// Account defines the common interface for bank accounts.
type Account interface {
	ID() string
	Owner() string
	Balance() int
	Deposit(amountCents int) error
	Withdraw(amountCents int) error
	Transactions() []Transaction
}

// TransactionBuffer is a generic ring buffer for storing transaction histories.
type TransactionBuffer[T any] struct {
	mu       sync.RWMutex
	items    []T
	capacity int
}

// NewTransactionBuffer initializes a type-safe generic buffer with maximum capacity.
func NewTransactionBuffer[T any](capacity int) *TransactionBuffer[T] {
	if capacity <= 0 {
		capacity = 100
	}
	return &TransactionBuffer[T]{
		items:    make([]T, 0, capacity),
		capacity: capacity,
	}
}

// Push adds an item to the buffer, evicting the oldest element if full.
func (b *TransactionBuffer[T]) Push(item T) {
	// TODO: Under b.mu.Lock(), append item, evicting oldest if len >= capacity
}

// All returns a point-in-time snapshot of elements in the buffer.
func (b *TransactionBuffer[T]) All() []T {
	// TODO: Under b.mu.RLock(), return a copied slice of items
	return nil
}

// Len returns the current count of elements in the buffer.
func (b *TransactionBuffer[T]) Len() int {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return len(b.items)
}

// CheckingAccount provides everyday transactional banking with overdraft allowance.
type CheckingAccount struct {
	mu             sync.RWMutex
	id             string
	owner          string
	balanceCents   int
	overdraftLimit int
	isFrozen       bool
	history        *TransactionBuffer[Transaction]
}

// NewCheckingAccount initializes a CheckingAccount.
func NewCheckingAccount(id, owner string, initialDeposit, overdraftLimit int) (*CheckingAccount, error) {
	// TODO: Validate initialDeposit >= 0, initialize struct with history buffer
	return nil, errors.New("not implemented")
}

func (c *CheckingAccount) ID() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.id
}

func (c *CheckingAccount) Owner() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.owner
}

func (c *CheckingAccount) Balance() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.balanceCents
}

func (c *CheckingAccount) Deposit(amountCents int) error {
	// TODO: Validate positive amount, update balanceCents and push Transaction under Lock
	return errors.New("not implemented")
}

func (c *CheckingAccount) Withdraw(amountCents int) error {
	// TODO: Validate positive amount, check overdraft limit (balance + overdraftLimit >= amount)
	// TODO: Deduct balance and push Transaction under Lock
	return errors.New("not implemented")
}

func (c *CheckingAccount) Transactions() []Transaction {
	if c.history == nil {
		return nil
	}
	return c.history.All()
}

// SavingsAccount holds funds with interest accrual and minimum balance restrictions.
type SavingsAccount struct {
	mu             sync.RWMutex
	id             string
	owner          string
	balanceCents   int
	minimumBalance int
	history        *TransactionBuffer[Transaction]
}

// NewSavingsAccount initializes a SavingsAccount.
func NewSavingsAccount(id, owner string, initialDeposit, minBalance int) (*SavingsAccount, error) {
	// TODO: Validate initialDeposit >= minBalance
	return nil, errors.New("not implemented")
}

func (s *SavingsAccount) ID() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.id
}

func (s *SavingsAccount) Owner() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.owner
}

func (s *SavingsAccount) Balance() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.balanceCents
}

func (s *SavingsAccount) Deposit(amountCents int) error {
	// TODO: Validate amount > 0, update balance and record transaction
	return errors.New("not implemented")
}

func (s *SavingsAccount) Withdraw(amountCents int) error {
	// TODO: Validate amount > 0 and balance - amount >= minimumBalance
	return errors.New("not implemented")
}

// AccrueInterest applies an interest rate percentage to the balance.
func (s *SavingsAccount) AccrueInterest(ratePercent float64) error {
	// TODO: Validate ratePercent > 0, compute interest, update balance and push Transaction
	return errors.New("not implemented")
}

func (s *SavingsAccount) Transactions() []Transaction {
	if s.history == nil {
		return nil
	}
	return s.history.All()
}

// Transfer moves funds atomically from sender to receiver.
func Transfer(from Account, to Account, amountCents int) error {
	// TODO: Check non-nil accounts and different account IDs
	// TODO: Withdraw from source; if deposit into destination fails, rollback withdrawal
	return errors.New("not implemented")
}

func main() {
	checking, _ := NewCheckingAccount("chk-101", "Alice Doe", 50000, 10000)
	savings, _ := NewSavingsAccount("sav-202", "Alice Doe", 100000, 20000)

	fmt.Printf("Checking Initial Balance: $%0.2f\n", float64(checking.Balance())/100)
	fmt.Printf("Savings Initial Balance: $%0.2f\n", float64(savings.Balance())/100)

	_ = checking.Withdraw(30000)
	_ = savings.AccrueInterest(2.5)

	_ = Transfer(savings, checking, 15000)

	fmt.Printf("Checking Final Balance: $%0.2f\n", float64(checking.Balance())/100)
	fmt.Printf("Savings Final Balance: $%0.2f\n", float64(savings.Balance())/100)
}
