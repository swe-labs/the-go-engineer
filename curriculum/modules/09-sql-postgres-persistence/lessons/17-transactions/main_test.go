package main

import (
	"database/sql"
	"testing"

	_ "modernc.org/sqlite"
)

func setupTestDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })

	_, err = db.Exec(`CREATE TABLE accounts (
		id INTEGER PRIMARY KEY,
		name TEXT NOT NULL,
		balance REAL NOT NULL
	)`)
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec(`INSERT INTO accounts VALUES
		(1, 'Alice', 1000.00),
		(2, 'Bob', 500.00)`)
	if err != nil {
		t.Fatal(err)
	}
	return db
}

func TestTransferSuccess(t *testing.T) {
	db := setupTestDB(t)

	err := Transfer(db, 1, 2, 200.00)
	if err != nil {
		t.Fatal(err)
	}

	var aliceBal, bobBal float64
	db.QueryRow("SELECT balance FROM accounts WHERE id = 1").Scan(&aliceBal)
	db.QueryRow("SELECT balance FROM accounts WHERE id = 2").Scan(&bobBal)

	if aliceBal != 800.00 {
		t.Errorf("Alice balance = %.2f, want 800.00", aliceBal)
	}
	if bobBal != 700.00 {
		t.Errorf("Bob balance = %.2f, want 700.00", bobBal)
	}
}

func TestTransferInsufficientFunds(t *testing.T) {
	db := setupTestDB(t)

	err := Transfer(db, 2, 1, 1000.00)
	if err == nil {
		t.Fatal("expected insufficient funds error, got nil")
	}

	var bobBal float64
	db.QueryRow("SELECT balance FROM accounts WHERE id = 2").Scan(&bobBal)
	if bobBal != 500.00 {
		t.Errorf("Bob balance should remain 500.00, got %.2f", bobBal)
	}
}

func TestTransferInvalidAccount(t *testing.T) {
	db := setupTestDB(t)

	err := Transfer(db, 1, 999, 50.00)
	if err == nil {
		t.Fatal("expected error for invalid account, got nil")
	}

	var aliceBal float64
	db.QueryRow("SELECT balance FROM accounts WHERE id = 1").Scan(&aliceBal)
	if aliceBal != 1000.00 {
		t.Errorf("Alice balance should remain 1000.00, got %.2f", aliceBal)
	}
}

func TestTransferZeroAmount(t *testing.T) {
	db := setupTestDB(t)

	err := Transfer(db, 1, 2, 0)
	if err == nil {
		t.Fatal("expected error for zero transfer")
	}
}

func TestTransferAtomicity(t *testing.T) {
	db := setupTestDB(t)

	err := Transfer(db, 1, 2, -50.00)
	if err == nil {
		t.Fatal("expected error for negative transfer")
	}

	var aliceBal float64
	db.QueryRow("SELECT balance FROM accounts WHERE id = 1").Scan(&aliceBal)
	if aliceBal != 1000.00 {
		t.Errorf("Alice balance should remain 1000.00 after failed transfer, got %.2f", aliceBal)
	}
}
