package main

import (
	"errors"
	"fmt"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type userDB struct {
	hashes map[string][]byte
}

func newUserDB() *userDB {
	return &userDB{hashes: make(map[string][]byte)}
}

func (db *userDB) CreateUser(username, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hash failed: %w", err)
	}
	db.hashes[username] = hash
	return nil
}

func (db *userDB) Authenticate(username, password string) error {
	hash, ok := db.hashes[username]
	if !ok {
		return errors.New("user not found")
	}
	if err := bcrypt.CompareHashAndPassword(hash, []byte(password)); err != nil {
		return errors.New("invalid password")
	}
	return nil
}

func CostBenchmark(password string, costs []int) {
	for _, cost := range costs {
		start := time.Now()
		hash, err := bcrypt.GenerateFromPassword([]byte(password), cost)
		duration := time.Since(start)
		if err != nil {
			log.Printf("Cost %d: error - %v", cost, err)
			continue
		}
		fmt.Printf("Cost %2d: %v (hash length: %d bytes)\n", cost, duration.Round(time.Millisecond), len(hash))
	}
}

func HashPassword(password string, cost int) ([]byte, error) {
	if len(password) < 8 {
		return nil, errors.New("password too short: minimum 8 characters")
	}
	if cost < 4 || cost > 31 {
		return nil, errors.New("cost must be between 4 and 31")
	}
	return bcrypt.GenerateFromPassword([]byte(password), cost)
}

func VerifyPassword(hash []byte, password string) bool {
	return bcrypt.CompareHashAndPassword(hash, []byte(password)) == nil
}

func main() {
	db := newUserDB()

	db.CreateUser("alice", "correct-horse-battery-staple")
	db.CreateUser("bob", "p@ssw0rd-rocks!")

	if err := db.Authenticate("alice", "correct-horse-battery-staple"); err != nil {
		log.Fatal("Unexpected auth failure:", err)
	}
	fmt.Println("alice: authenticated successfully")

	if err := db.Authenticate("alice", "wrong-password"); err != nil {
		fmt.Println("alice with wrong password: rejected -", err)
	}

	if err := db.Authenticate("eve", "password"); err != nil {
		fmt.Println("eve: rejected -", err)
	}

	fmt.Println("\n=== Cost benchmark ===")
	CostBenchmark("test-password-123", []int{4, 6, 8})
}
