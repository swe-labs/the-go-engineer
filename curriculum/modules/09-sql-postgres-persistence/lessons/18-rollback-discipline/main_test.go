package main

import (
	"database/sql"
	"errors"
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

	_, err = db.Exec(`CREATE TABLE items (
		id INTEGER PRIMARY KEY,
		name TEXT NOT NULL,
		quantity INTEGER NOT NULL DEFAULT 0
	)`)
	if err != nil {
		t.Fatal(err)
	}
	return db
}

func TestWithTxSuccess(t *testing.T) {
	db := setupTestDB(t)

	err := WithTx(db, func(tx *sql.Tx) error {
		_, err := tx.Exec("INSERT INTO items (name, quantity) VALUES (?, ?)", "Widget", 10)
		return err
	})
	if err != nil {
		t.Fatal(err)
	}

	var count int
	db.QueryRow("SELECT COUNT(*) FROM items").Scan(&count)
	if count != 1 {
		t.Errorf("expected 1 item, got %d", count)
	}
}

func TestWithTxRollbackOnError(t *testing.T) {
	db := setupTestDB(t)

	expectedErr := errors.New("something went wrong")
	err := WithTx(db, func(tx *sql.Tx) error {
		tx.Exec("INSERT INTO items (name, quantity) VALUES (?, ?)", "Widget", 10)
		return expectedErr
	})
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected %v, got %v", expectedErr, err)
	}

	var count int
	db.QueryRow("SELECT COUNT(*) FROM items").Scan(&count)
	if count != 0 {
		t.Errorf("expected 0 items after rollback, got %d", count)
	}
}

func TestWithTxRollbackOnPrepareError(t *testing.T) {
	db := setupTestDB(t)

	err := WithTx(db, func(tx *sql.Tx) error {
		tx.Exec("INSERT INTO items (name, quantity) VALUES (?, ?)", "Widget", 10)
		_, err := tx.Exec("INVALID SQL")
		return err
	})
	if err == nil {
		t.Fatal("expected error from invalid SQL, got nil")
	}

	var count int
	db.QueryRow("SELECT COUNT(*) FROM items").Scan(&count)
	if count != 0 {
		t.Errorf("expected 0 items after rollback, got %d", count)
	}
}

func TestWithTxMultipleOperations(t *testing.T) {
	db := setupTestDB(t)

	err := WithTx(db, func(tx *sql.Tx) error {
		if _, err := tx.Exec("INSERT INTO items (name, quantity) VALUES (?, ?)", "A", 1); err != nil {
			return err
		}
		if _, err := tx.Exec("INSERT INTO items (name, quantity) VALUES (?, ?)", "B", 2); err != nil {
			return err
		}
		if _, err := tx.Exec("INSERT INTO items (name, quantity) VALUES (?, ?)", "C", 3); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}

	var count int
	db.QueryRow("SELECT COUNT(*) FROM items").Scan(&count)
	if count != 3 {
		t.Errorf("expected 3 items, got %d", count)
	}
}

func TestWithTxBeginError(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	db.Close() // force Begin to fail

	err = WithTx(db, func(tx *sql.Tx) error {
		return nil
	})
	if err == nil {
		t.Fatal("expected error from closed db, got nil")
	}
}
