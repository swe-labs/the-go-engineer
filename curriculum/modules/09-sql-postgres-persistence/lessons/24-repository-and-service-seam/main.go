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

type UserRepository interface {
	GetByID(ctx context.Context, id int) (*User, error)
	List(ctx context.Context) ([]User, error)
	Create(ctx context.Context, name, email string) (*User, error)
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetByID(ctx context.Context, id int) (*User, error) {
	u := &User{}
	err := r.db.QueryRowContext(ctx, "SELECT id, name, email FROM users WHERE id = ?", id).
		Scan(&u.ID, &u.Name, &u.Email)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (r *userRepository) List(ctx context.Context) ([]User, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, name, email FROM users ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, rows.Err()
}

func (r *userRepository) Create(ctx context.Context, name, email string) (*User, error) {
	result, err := r.db.ExecContext(ctx, "INSERT INTO users (name, email) VALUES (?, ?)", name, email)
	if err != nil {
		return nil, err
	}
	id, _ := result.LastInsertId()
	return r.GetByID(ctx, int(id))
}

type UserService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Register(ctx context.Context, name, email string) (*User, error) {
	if name == "" {
		return nil, fmt.Errorf("name cannot be empty")
	}
	if email == "" {
		return nil, fmt.Errorf("email cannot be empty")
	}
	return s.repo.Create(ctx, name, email)
}

func (s *UserService) GetUser(ctx context.Context, id int) (*User, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *UserService) ListUsers(ctx context.Context) ([]User, error) {
	return s.repo.List(ctx)
}

func main() {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	db.Exec(`CREATE TABLE users (id INTEGER PRIMARY KEY, name TEXT NOT NULL, email TEXT NOT NULL UNIQUE)`)

	repo := NewUserRepository(db)
	svc := NewUserService(repo)

	alice, _ := svc.Register(context.Background(), "Alice", "alice@example.com")
	bob, _ := svc.Register(context.Background(), "Bob", "bob@example.com")
	fmt.Printf("Registered: %s (id=%d), %s (id=%d)\n", alice.Name, alice.ID, bob.Name, bob.ID)

	users, _ := svc.ListUsers(context.Background())
	fmt.Printf("Total users: %d\n", len(users))
}
