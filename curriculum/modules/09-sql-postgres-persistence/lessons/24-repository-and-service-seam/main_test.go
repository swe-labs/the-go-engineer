package main

import (
	"context"
	"database/sql"
	"errors"
	"sync"
	"testing"
)

type fakeUserRepo struct {
	mu     sync.Mutex
	users  map[int]User
	nextID int
}

func newFakeUserRepo() *fakeUserRepo {
	return &fakeUserRepo{users: make(map[int]User), nextID: 1}
}

func (f *fakeUserRepo) GetByID(_ context.Context, id int) (*User, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	u, ok := f.users[id]
	if !ok {
		return nil, sql.ErrNoRows
	}
	return &u, nil
}

func (f *fakeUserRepo) List(_ context.Context) ([]User, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	var users []User
	for _, u := range f.users {
		users = append(users, u)
	}
	return users, nil
}

func (f *fakeUserRepo) Create(_ context.Context, name, email string) (*User, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	// Simulate UNIQUE constraint
	for _, u := range f.users {
		if u.Email == email {
			return nil, errors.New("UNIQUE constraint failed: email already exists")
		}
	}
	id := f.nextID
	f.nextID++
	u := User{ID: id, Name: name, Email: email}
	f.users[id] = u
	return &u, nil
}

func TestUserServiceRegister(t *testing.T) {
	repo := newFakeUserRepo()
	svc := NewUserService(repo)

	u, err := svc.Register(context.Background(), "Alice", "alice@example.com")
	if err != nil {
		t.Fatal(err)
	}
	if u.Name != "Alice" {
		t.Errorf("expected Alice, got %s", u.Name)
	}
	if u.ID != 1 {
		t.Errorf("expected id 1, got %d", u.ID)
	}
}

func TestUserServiceRegisterEmptyName(t *testing.T) {
	svc := NewUserService(newFakeUserRepo())

	_, err := svc.Register(context.Background(), "", "alice@example.com")
	if err == nil {
		t.Fatal("expected error for empty name")
	}
}

func TestUserServiceRegisterEmptyEmail(t *testing.T) {
	svc := NewUserService(newFakeUserRepo())

	_, err := svc.Register(context.Background(), "Alice", "")
	if err == nil {
		t.Fatal("expected error for empty email")
	}
}

func TestUserServiceGetUser(t *testing.T) {
	repo := newFakeUserRepo()
	svc := NewUserService(repo)

	created, _ := svc.Register(context.Background(), "Alice", "alice@example.com")
	got, err := svc.GetUser(context.Background(), created.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.Name != "Alice" {
		t.Errorf("expected Alice, got %s", got.Name)
	}
}

func TestUserServiceGetUserNotFound(t *testing.T) {
	svc := NewUserService(newFakeUserRepo())

	_, err := svc.GetUser(context.Background(), 999)
	if !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("expected ErrNoRows, got %v", err)
	}
}

func TestUserServiceListUsers(t *testing.T) {
	repo := newFakeUserRepo()
	svc := NewUserService(repo)

	svc.Register(context.Background(), "Alice", "alice@example.com")
	svc.Register(context.Background(), "Bob", "bob@example.com")

	users, err := svc.ListUsers(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(users) != 2 {
		t.Errorf("expected 2 users, got %d", len(users))
	}
}

func TestUserServiceDuplicateEmail(t *testing.T) {
	repo := newFakeUserRepo()
	svc := NewUserService(repo)

	svc.Register(context.Background(), "Alice", "alice@example.com")
	_, err := svc.Register(context.Background(), "Alice2", "alice@example.com")
	if err == nil {
		t.Fatal("expected error for duplicate email")
	}
}

func TestUserServiceEmptyList(t *testing.T) {
	svc := NewUserService(newFakeUserRepo())

	users, err := svc.ListUsers(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(users) != 0 {
		t.Errorf("expected empty list, got %d users", len(users))
	}
}

func TestSQLiteRepository(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })

	db.Exec(`CREATE TABLE users (id INTEGER PRIMARY KEY, name TEXT NOT NULL, email TEXT NOT NULL UNIQUE)`)

	repo := NewUserRepository(db)
	svc := NewUserService(repo)

	u, err := svc.Register(context.Background(), "Alice", "alice@example.com")
	if err != nil {
		t.Fatal(err)
	}
	if u.Name != "Alice" {
		t.Errorf("expected Alice, got %s", u.Name)
	}

	got, err := svc.GetUser(context.Background(), u.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.Email != "alice@example.com" {
		t.Errorf("expected alice@example.com, got %s", got.Email)
	}
}
