package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "modernc.org/sqlite"
)

func main() {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(`CREATE TABLE accounts (
		id INTEGER PRIMARY KEY,
		name TEXT NOT NULL,
		balance REAL NOT NULL
	)`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`INSERT INTO accounts VALUES
		(1, 'Alice', 1000.00),
		(2, 'Bob', 500.00)`)
	if err != nil {
		log.Fatal(err)
	}

	err = Transfer(db, 1, 2, 200.00)
	if err != nil {
		log.Fatal(err)
	}

	var aliceBalance, bobBalance float64
	db.QueryRow("SELECT balance FROM accounts WHERE id = 1").Scan(&aliceBalance)
	db.QueryRow("SELECT balance FROM accounts WHERE id = 2").Scan(&bobBalance)
	fmt.Printf("Alice: $%.2f, Bob: $%.2f\n", aliceBalance, bobBalance)
}

func Transfer(db *sql.DB, fromID, toID int, amount float64) (err error) {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	var fromBalance float64
	err = tx.QueryRow("SELECT balance FROM accounts WHERE id = ?", fromID).Scan(&fromBalance)
	if err != nil {
		return fmt.Errorf("read from account: %w", err)
	}
	if amount <= 0 {
		err = fmt.Errorf("transfer amount must be positive: %.2f", amount)
		return err
	}
	if fromBalance < amount {
		err = fmt.Errorf("insufficient funds: %.2f < %.2f", fromBalance, amount)
		return err
	}

	_, err = tx.Exec("UPDATE accounts SET balance = balance - ? WHERE id = ?", amount, fromID)
	if err != nil {
		return fmt.Errorf("debit: %w", err)
	}

	var toExists int
	err = tx.QueryRow("SELECT COUNT(*) FROM accounts WHERE id = ?", toID).Scan(&toExists)
	if err != nil {
		return fmt.Errorf("check target account: %w", err)
	}
	if toExists == 0 {
		err = fmt.Errorf("target account %d does not exist", toID)
		return err
	}

	_, err = tx.Exec("UPDATE accounts SET balance = balance + ? WHERE id = ?", amount, toID)
	if err != nil {
		return fmt.Errorf("credit: %w", err)
	}

	return tx.Commit()
}
