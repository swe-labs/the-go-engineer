package main

import (
	"errors"
	"fmt"
	"sync"
)

type User struct {
	ID    string
	Email string
}

type UserRepository interface {
	Save(user User) error
	FindByID(id string) (User, error)
}

type EmailService interface {
	SendWelcome(email string) error
}

type UserService struct {
	repo  UserRepository
	email EmailService
}

func NewUserService(repo UserRepository, email EmailService) *UserService {
	return &UserService{repo: repo, email: email}
}

func (s *UserService) RegisterUser(id, email string) error {
	if email == "" {
		return errors.New("email cannot be empty")
	}
	user := User{ID: id, Email: email}
	if err := s.email.SendWelcome(email); err != nil {
		return fmt.Errorf("welcome email failed: %w", err)
	}
	if err := s.repo.Save(user); err != nil {
		return fmt.Errorf("save failed: %w", err)
	}
	return nil
}

type InMemoryUserRepoAdapter struct {
	mu    sync.RWMutex
	users map[string]User
}

func NewInMemoryUserRepoAdapter() *InMemoryUserRepoAdapter {
	return &InMemoryUserRepoAdapter{users: make(map[string]User)}
}

func (a *InMemoryUserRepoAdapter) Save(user User) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.users[user.ID] = user
	return nil
}

func (a *InMemoryUserRepoAdapter) FindByID(id string) (User, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	u, ok := a.users[id]
	if !ok {
		return User{}, errors.New("user not found")
	}
	return u, nil
}

type ConsoleEmailAdapter struct{}

func (ConsoleEmailAdapter) SendWelcome(email string) error {
	fmt.Printf("Welcome email sent to %s\n", email)
	return nil
}

type LoggerMiddleware struct {
	next UserRepository
	logs []string
}

func (m *LoggerMiddleware) Save(user User) error {
	m.logs = append(m.logs, fmt.Sprintf("saving user %s", user.ID))
	return m.next.Save(user)
}

func (m *LoggerMiddleware) FindByID(id string) (User, error) {
	m.logs = append(m.logs, fmt.Sprintf("finding user %s", id))
	return m.next.FindByID(id)
}

func main() {
	repo := NewInMemoryUserRepoAdapter()
	email := ConsoleEmailAdapter{}
	svc := NewUserService(repo, email)

	if err := svc.RegisterUser("u1", "alice@example.com"); err != nil {
		fmt.Println("Error:", err)
		return
	}
	user, _ := repo.FindByID("u1")
	fmt.Printf("Registered user: %s (%s)\n", user.ID, user.Email)
}
