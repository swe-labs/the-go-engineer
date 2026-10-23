package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

var (
	ErrNegativeAmount     = errors.New("amount must be strictly positive")
	ErrInsufficientFunds  = errors.New("insufficient funds for withdrawal")
	ErrAccountFrozen      = errors.New("account is frozen")
	ErrBufferEmpty        = errors.New("buffer is empty")
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
	b.mu.Lock()
	defer b.mu.Unlock()

	if len(b.items) >= b.capacity {
		b.items = b.items[1:]
	}
	b.items = append(b.items, item)
}

// All returns a point-in-time snapshot of elements in the buffer.
func (b *TransactionBuffer[T]) All() []T {
	b.mu.RLock()
	defer b.mu.RUnlock()

	res := make([]T, len(b.items))
	copy(res, b.items)
	return res
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
	if initialDeposit < 0 {
		return nil, ErrNegativeAmount
	}
	if overdraftLimit < 0 {
		overdraftLimit = 0
	}

	acc := &CheckingAccount{
		id:             id,
		owner:          owner,
		balanceCents:   initialDeposit,
		overdraftLimit: overdraftLimit,
		history:        NewTransactionBuffer[Transaction](50),
	}

	if initialDeposit > 0 {
		acc.history.Push(Transaction{
			ID:        fmt.Sprintf("tx-%d", 1),
			Type:      TxDeposit,
			Amount:    initialDeposit,
			Balance:   initialDeposit,
			Timestamp: time.Now().UTC(),
		})
	}

	return acc, nil
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
	if amountCents <= 0 {
		return ErrNegativeAmount
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if c.isFrozen {
		return ErrAccountFrozen
	}

	c.balanceCents += amountCents
	c.history.Push(Transaction{
		ID:        fmt.Sprintf("tx-%d", c.history.Len()+1),
		Type:      TxDeposit,
		Amount:    amountCents,
		Balance:   c.balanceCents,
		Timestamp: time.Now().UTC(),
	})

	return nil
}

func (c *CheckingAccount) Withdraw(amountCents int) error {
	if amountCents <= 0 {
		return ErrNegativeAmount
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if c.isFrozen {
		return ErrAccountFrozen
	}

	maxAllowed := c.balanceCents + c.overdraftLimit
	if amountCents > maxAllowed {
		return ErrInsufficientFunds
	}

	c.balanceCents -= amountCents
	c.history.Push(Transaction{
		ID:        fmt.Sprintf("tx-%d", c.history.Len()+1),
		Type:      TxWithdraw,
		Amount:    amountCents,
		Balance:   c.balanceCents,
		Timestamp: time.Now().UTC(),
	})

	return nil
}

func (c *CheckingAccount) Transactions() []Transaction {
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
	if initialDeposit < minBalance {
		return nil, ErrBelowMinimumBalance
	}

	acc := &SavingsAccount{
		id:             id,
		owner:          owner,
		balanceCents:   initialDeposit,
		minimumBalance: minBalance,
		history:        NewTransactionBuffer[Transaction](50),
	}

	if initialDeposit > 0 {
		acc.history.Push(Transaction{
			ID:        fmt.Sprintf("tx-%d", 1),
			Type:      TxDeposit,
			Amount:    initialDeposit,
			Balance:   initialDeposit,
			Timestamp: time.Now().UTC(),
		})
	}

	return acc, nil
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
	if amountCents <= 0 {
		return ErrNegativeAmount
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.balanceCents += amountCents
	s.history.Push(Transaction{
		ID:        fmt.Sprintf("tx-%d", s.history.Len()+1),
		Type:      TxDeposit,
		Amount:    amountCents,
		Balance:   s.balanceCents,
		Timestamp: time.Now().UTC(),
	})

	return nil
}

func (s *SavingsAccount) Withdraw(amountCents int) error {
	if amountCents <= 0 {
		return ErrNegativeAmount
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.balanceCents-amountCents < s.minimumBalance {
		return ErrBelowMinimumBalance
	}

	s.balanceCents -= amountCents
	s.history.Push(Transaction{
		ID:        fmt.Sprintf("tx-%d", s.history.Len()+1),
		Type:      TxWithdraw,
		Amount:    amountCents,
		Balance:   s.balanceCents,
		Timestamp: time.Now().UTC(),
	})

	return nil
}

// AccrueInterest applies an interest rate percentage to the balance.
func (s *SavingsAccount) AccrueInterest(ratePercent float64) error {
	if ratePercent <= 0 {
		return errors.New("interest rate must be positive")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	interest := int(float64(s.balanceCents) * (ratePercent / 100.0))
	if interest > 0 {
		s.balanceCents += interest
		s.history.Push(Transaction{
			ID:        fmt.Sprintf("tx-%d", s.history.Len()+1),
			Type:      TxInterest,
			Amount:    interest,
			Balance:   s.balanceCents,
			Timestamp: time.Now().UTC(),
		})
	}

	return nil
}

func (s *SavingsAccount) Transactions() []Transaction {
	return s.history.All()
}

// Transfer moves funds atomically from sender to receiver.
func Transfer(from Account, to Account, amountCents int) error {
	if from == nil || to == nil {
		return errors.New("from and to accounts must be non-nil")
	}
	if from.ID() == to.ID() {
		return errors.New("cannot transfer to the same account")
	}

	if err := from.Withdraw(amountCents); err != nil {
		return fmt.Errorf("transfer withdrawal failed: %w", err)
	}

	if err := to.Deposit(amountCents); err != nil {
		// Roll back withdrawal
		_ = from.Deposit(amountCents)
		return fmt.Errorf("transfer deposit failed: %w", err)
	}

	return nil
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
