//go:build integration

package main

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	_ "modernc.org/sqlite"
)

func setupDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })

	_, err = db.Exec(`CREATE TABLE users (
		id INTEGER PRIMARY KEY,
		name TEXT NOT NULL,
		email TEXT NOT NULL UNIQUE
	)`)
	if err != nil {
		t.Fatal(err)
	}
	return db
}

func TestIntegrationCreateUser(t *testing.T) {
	db := setupDB(t)
	repo := NewUserRepository(db)

	u, err := repo.Create(context.Background(), "Alice", "alice@example.com")
	if err != nil {
		t.Fatal(err)
	}
	if u.Name != "Alice" {
		t.Errorf("name = %q, want %q", u.Name, "Alice")
	}
	if u.ID == 0 {
		t.Error("expected non-zero ID")
	}
}

func TestIntegrationGetByID(t *testing.T) {
	db := setupDB(t)
	repo := NewUserRepository(db)

	created, _ := repo.Create(context.Background(), "Bob", "bob@example.com")
	got, err := repo.GetByID(context.Background(), created.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.Email != "bob@example.com" {
		t.Errorf("email = %q, want %q", got.Email, "bob@example.com")
	}
}

func TestIntegrationGetByIDNotFound(t *testing.T) {
	db := setupDB(t)
	repo := NewUserRepository(db)

	_, err := repo.GetByID(context.Background(), 999)
	if !errors.Is(err, sql.ErrNoRows) {
		t.Errorf("expected ErrNoRows, got %v", err)
	}
}

func TestIntegrationUniqueConstraint(t *testing.T) {
	db := setupDB(t)
	repo := NewUserRepository(db)

	repo.Create(context.Background(), "Alice", "alice@example.com")
	_, err := repo.Create(context.Background(), "Alice2", "alice@example.com")
	if err == nil {
		t.Fatal("expected unique constraint error")
	}
}

func TestIntegrationMultipleUsers(t *testing.T) {
	db := setupDB(t)
	repo := NewUserRepository(db)

	users := []struct{ name, email string }{
		{"Alice", "alice@example.com"},
		{"Bob", "bob@example.com"},
		{"Carol", "carol@example.com"},
	}

	for _, u := range users {
		created, err := repo.Create(context.Background(), u.name, u.email)
		if err != nil {
			t.Fatalf("failed to create %s: %v", u.name, err)
		}
		if created.Name != u.name {
			t.Errorf("name = %q, want %q", created.Name, u.name)
		}
	}
}
