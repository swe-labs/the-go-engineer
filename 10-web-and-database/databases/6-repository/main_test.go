// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

package main

import (
	"database/sql"
	"path/filepath"
	"testing"

	_ "github.com/mattn/go-sqlite3"

	"github.com/rasel9t6/the-go-engineer/10-web-and-database/databases/6-repository/repository"
)

func openRepositoryTestDB(t *testing.T) *sql.DB {
	t.Helper()

	dbPath := filepath.Join(t.TempDir(), "users.db")
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		t.Fatalf("open test database: %v", err)
	}

	t.Cleanup(func() {
		_ = db.Close()
	})

	if err := db.Ping(); err != nil {
		t.Fatalf("ping test database: %v", err)
	}

	if _, err := db.Exec(repositorySchema); err != nil {
		t.Fatalf("create repository schema: %v", err)
	}

	return db
}

func TestRepositoryCreateUserAndGetByEmail(t *testing.T) {
	db := openRepositoryTestDB(t)
	repo := repository.NewSQLUserRepository(db)

	id, err := repo.CreateUser("Alice Repo", "alice.repo@example.com", "supersecret", "alice.png")
	if err != nil {
		t.Fatalf("create user: %v", err)
	}
	if id <= 0 {
		t.Fatalf("expected positive user id, got %d", id)
	}

	user, err := repo.GetUserByEmail("alice.repo@example.com")
	if err != nil {
		t.Fatalf("get user by email: %v", err)
	}
	if user.Name != "Alice Repo" {
		t.Fatalf("unexpected user name: %q", user.Name)
	}
	if user.Profile.Avatar != "alice.png" {
		t.Fatalf("unexpected avatar: %q", user.Profile.Avatar)
	}
	if user.HashedPassword == "supersecret" || user.HashedPassword == "" {
		t.Fatalf("expected stored password to be hashed")
	}
}

func TestRepositoryGetUsersReturnsCreatedRows(t *testing.T) {
	db := openRepositoryTestDB(t)
	repo := repository.NewSQLUserRepository(db)

	if _, err := repo.CreateUser("Alice Repo", "alice@example.com", "secret123", "alice.png"); err != nil {
		t.Fatalf("create first user: %v", err)
	}
	if _, err := repo.CreateUser("Bob Repo", "bob@example.com", "secret456", "bob.png"); err != nil {
		t.Fatalf("create second user: %v", err)
	}

	users, err := repo.GetUsers()
	if err != nil {
		t.Fatalf("get users: %v", err)
	}
	if len(users) != 2 {
		t.Fatalf("expected 2 users, got %d", len(users))
	}
}
