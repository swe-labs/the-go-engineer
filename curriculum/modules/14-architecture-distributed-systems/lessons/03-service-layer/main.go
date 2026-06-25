package main

import (
	"errors"
	"fmt"
)

var errUserNotFound = errors.New("user not found")

type User struct {
	ID    string
	Name  string
	Email string
}

type UserRepository interface {
	Save(user User) error
	FindByID(id string) (User, error)
}

type UserService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Register(name, email string) (User, error) {
	if name == "" {
		return User{}, errors.New("name cannot be empty")
	}
	if email == "" {
		return User{}, errors.New("email cannot be empty")
	}
	user := User{ID: fmt.Sprintf("usr_%s", email), Name: name, Email: email}
	if err := s.repo.Save(user); err != nil {
		return User{}, fmt.Errorf("save failed: %w", err)
	}
	return user, nil
}

func (s *UserService) GetUser(id string) (User, error) {
	return s.repo.FindByID(id)
}

type inMemoryUserRepo struct {
	users map[string]User
}

func newInMemoryUserRepo() *inMemoryUserRepo {
	return &inMemoryUserRepo{users: make(map[string]User)}
}

func (r *inMemoryUserRepo) Save(user User) error {
	r.users[user.ID] = user
	return nil
}

func (r *inMemoryUserRepo) FindByID(id string) (User, error) {
	user, ok := r.users[id]
	if !ok {
		return User{}, errUserNotFound
	}
	return user, nil
}

func main() {
	repo := newInMemoryUserRepo()
	svc := NewUserService(repo)

	user, err := svc.Register("Alice", "alice@example.com")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Printf("Registered: %+v\n", user)

	fetched, err := svc.GetUser(user.ID)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Printf("Fetched: %+v\n", fetched)
}
