package main

import (
	"database/sql"
	"testing"
	"time"

	_ "modernc.org/sqlite"
)

func setupTestDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })

	db.Exec("PRAGMA foreign_keys = ON")

	_, err = db.Exec(`
		CREATE TABLE users (
			id INTEGER PRIMARY KEY,
			name TEXT NOT NULL,
			deleted_at TIMESTAMP
		);
		CREATE TABLE orders (
			id INTEGER PRIMARY KEY,
			user_id INTEGER NOT NULL,
			total REAL NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id)
		);
		INSERT INTO users VALUES (1, 'Alice', NULL);
		INSERT INTO users VALUES (2, 'Bob', NULL);
		INSERT INTO users VALUES (3, 'Carol', NULL);
		INSERT INTO orders VALUES (101, 1, 50.00, '2024-01-15');
		INSERT INTO orders VALUES (102, 1, 30.00, '2024-06-20');
		INSERT INTO orders VALUES (103, 3, 20.00, '2024-03-10');
	`)
	if err != nil {
		t.Fatal(err)
	}
	return db
}

func TestHardDelete(t *testing.T) {
	db := setupTestDB(t)

	result, err := db.Exec("DELETE FROM orders WHERE id = ?", 101)
	if err != nil {
		t.Fatal(err)
	}
	rows, _ := result.RowsAffected()
	if rows != 1 {
		t.Errorf("expected 1 row deleted, got %d", rows)
	}

	var count int
	db.QueryRow("SELECT COUNT(*) FROM orders").Scan(&count)
	if count != 2 {
		t.Errorf("expected 2 orders remaining, got %d", count)
	}
}

func TestHardDeleteNoMatch(t *testing.T) {
	db := setupTestDB(t)

	result, err := db.Exec("DELETE FROM orders WHERE id = ?", 999)
	if err != nil {
		t.Fatal(err)
	}
	rows, _ := result.RowsAffected()
	if rows != 0 {
		t.Errorf("expected 0 rows, got %d", rows)
	}
}

func TestSoftDelete(t *testing.T) {
	db := setupTestDB(t)

	result, err := db.Exec("UPDATE users SET deleted_at = ? WHERE id = ?", time.Now().UTC(), 1)
	if err != nil {
		t.Fatal(err)
	}
	rows, _ := result.RowsAffected()
	if rows != 1 {
		t.Errorf("expected 1 user soft deleted, got %d", rows)
	}

	var active, deleted int
	db.QueryRow("SELECT COUNT(*) FROM users WHERE deleted_at IS NULL").Scan(&active)
	db.QueryRow("SELECT COUNT(*) FROM users WHERE deleted_at IS NOT NULL").Scan(&deleted)
	if active != 2 {
		t.Errorf("expected 2 active users, got %d", active)
	}
	if deleted != 1 {
		t.Errorf("expected 1 deleted user, got %d", deleted)
	}
}

func TestSoftDeleteRestore(t *testing.T) {
	db := setupTestDB(t)

	db.Exec("UPDATE users SET deleted_at = ? WHERE id = ?", time.Now().UTC(), 2)
	result, err := db.Exec("UPDATE users SET deleted_at = NULL WHERE id = ?", 2)
	if err != nil {
		t.Fatal(err)
	}
	rows, _ := result.RowsAffected()
	if rows != 1 {
		t.Errorf("expected 1 user restored, got %d", rows)
	}

	var active int
	db.QueryRow("SELECT COUNT(*) FROM users WHERE deleted_at IS NULL").Scan(&active)
	if active != 3 {
		t.Errorf("expected 3 active users, got %d", active)
	}
}

func TestBatchDelete(t *testing.T) {
	db := setupTestDB(t)

	result, err := db.Exec("DELETE FROM orders WHERE id IN (?, ?)", 101, 103)
	if err != nil {
		t.Fatal(err)
	}
	rows, _ := result.RowsAffected()
	if rows != 2 {
		t.Errorf("expected 2 rows deleted, got %d", rows)
	}

	var count int
	db.QueryRow("SELECT COUNT(*) FROM orders").Scan(&count)
	if count != 1 {
		t.Errorf("expected 1 order remaining, got %d", count)
	}
}

func TestForeignKeyPreventsDelete(t *testing.T) {
	db := setupTestDB(t)

	_, err := db.Exec("DELETE FROM users WHERE id = ?", 1)
	if err == nil {
		t.Error("expected foreign key error, got nil")
	}
}

func TestDeleteAllOrdersByUser(t *testing.T) {
	db := setupTestDB(t)

	result, err := db.Exec("DELETE FROM orders WHERE user_id = ?", 1)
	if err != nil {
		t.Fatal(err)
	}
	rows, _ := result.RowsAffected()
	if rows != 2 {
		t.Errorf("expected 2 orders deleted, got %d", rows)
	}

	var count int
	db.QueryRow("SELECT COUNT(*) FROM orders").Scan(&count)
	if count != 1 {
		t.Errorf("expected 1 order remaining, got %d", count)
	}
}

func TestDeleteWithTransaction(t *testing.T) {
	db := setupTestDB(t)

	tx, err := db.Begin()
	if err != nil {
		t.Fatal(err)
	}

	tx.Exec("DELETE FROM orders WHERE user_id = ?", 1)
	tx.Exec("DELETE FROM orders WHERE user_id = ?", 3)
	err = tx.Commit()
	if err != nil {
		t.Fatal(err)
	}

	var count int
	db.QueryRow("SELECT COUNT(*) FROM orders").Scan(&count)
	if count != 0 {
		t.Errorf("expected 0 orders after transaction, got %d", count)
	}
}
