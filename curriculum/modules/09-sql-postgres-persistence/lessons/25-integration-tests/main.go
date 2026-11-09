package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	_ "modernc.org/sqlite"
)

type User struct {
	ID    int
	Name  string
	Email string
}

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetByID(ctx context.Context, id int) (*User, error) {
	u := &User{}
	err := r.db.QueryRowContext(ctx, "SELECT id, name, email FROM users WHERE id = ?", id).
		Scan(&u.ID, &u.Name, &u.Email)
	return u, err
}

func (r *UserRepository) Create(ctx context.Context, name, email string) (*User, error) {
	result, err := r.db.ExecContext(ctx, "INSERT INTO users (name, email) VALUES (?, ?)", name, email)
	if err != nil {
		return nil, err
	}
	id, _ := result.LastInsertId()
	return r.GetByID(ctx, int(id))
}

func main() {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	db.Exec(`CREATE TABLE users (id INTEGER PRIMARY KEY, name TEXT NOT NULL, email TEXT NOT NULL UNIQUE)`)

	repo := NewUserRepository(db)
	u, _ := repo.Create(context.Background(), "Alice", "alice@example.com")
	fmt.Printf("Created user %d: %s\n", u.ID, u.Name)
}
