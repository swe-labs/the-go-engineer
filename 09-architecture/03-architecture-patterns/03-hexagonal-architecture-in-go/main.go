// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 09: Architecture & Security
// Title: Hexagonal architecture in Go
// Level: Production
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Learn how ports and adapters keep domain rules independent of transport and storage details.
//
// WHY THIS MATTERS:
//   - Hexagonal architecture isolates the core from delivery and persistence concerns through explicit ports.
//
// RUN:
//   go run ./09-architecture/03-architecture-patterns/03-hexagonal-architecture-in-go
//
// KEY TAKEAWAY:
//   - The domain core defines port interfaces. Adapters implement them. The core never imports an adapter.
// ============================================================================

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

// User (Struct): core domain entity with zero external dependencies.
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// UserRepository (Interface): outbound port that adapters must satisfy.
type UserRepository interface {
	Save(user User) error
	FindByID(id int) (User, bool)
	FindAll() ([]User, error)
}

// UserService (Struct): application service that orchestrates domain logic via port interfaces.
type UserService struct {
	repo UserRepository
}

// NewUserService (Constructor): injects a UserRepository into a new UserService.
func NewUserService(repo UserRepository) *UserService {
	return &UserService{repo: repo}
}

// UserService.Register (Method): persists a user through the repository port.
func (s *UserService) Register(id int, name, email string) error {
	return s.repo.Save(User{ID: id, Name: name, Email: email})
}

// UserService.Lookup (Method): retrieves a user by ID through the repository port.
func (s *UserService) Lookup(id int) (User, bool) {
	return s.repo.FindByID(id)
}

// UserService.List (Method): returns all users through the repository port.
func (s *UserService) List() ([]User, error) {
	return s.repo.FindAll()
}

// InMemoryUserRepository (Struct): in-memory adapter backed by a mutex-protected map.
type InMemoryUserRepository struct {
	mu   sync.RWMutex
	data map[int]User
}

// NewInMemoryUserRepository (Constructor): returns an initialized in-memory repository.
func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{data: make(map[int]User)}
}

// InMemoryUserRepository.Save (Method): upserts a user under a write lock.
func (r *InMemoryUserRepository) Save(user User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.data[user.ID] = user
	return nil
}

// InMemoryUserRepository.FindByID (Method): reads a user by ID under a read lock.
func (r *InMemoryUserRepository) FindByID(id int) (User, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	u, ok := r.data[id]
	return u, ok
}

// InMemoryUserRepository.FindAll (Method): returns a copy of all users under a read lock.
func (r *InMemoryUserRepository) FindAll() ([]User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]User, 0, len(r.data))
	for _, u := range r.data {
		result = append(result, u)
	}
	return result, nil
}

// FileUserRepository (Struct): file-based adapter backed by a JSON-lines file.
type FileUserRepository struct {
	mu   sync.RWMutex
	path string
}

// NewFileUserRepository (Constructor): creates a repository backed by a JSON-lines file at the given path.
func NewFileUserRepository(path string) *FileUserRepository {
	return &FileUserRepository{path: path}
}

// FileUserRepository.Save (Method): appends a user as a JSON line under a write lock.
func (r *FileUserRepository) Save(user User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	f, err := os.OpenFile(r.path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(user)
}

// FileUserRepository.FindByID (Method): scans all stored lines for a matching ID under a read lock.
func (r *FileUserRepository) FindByID(id int) (User, bool) {
	users, err := r.loadAll()
	if err != nil {
		return User{}, false
	}
	for _, u := range users {
		if u.ID == id {
			return u, true
		}
	}
	return User{}, false
}

// FileUserRepository.FindAll (Method): returns all users by delegating to loadAll.
func (r *FileUserRepository) FindAll() ([]User, error) {
	return r.loadAll()
}

// FileUserRepository.loadAll (Method): reads and decodes every JSON line from the file under a read lock.
func (r *FileUserRepository) loadAll() ([]User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	f, err := os.Open(r.path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	defer f.Close()
	var users []User
	dec := json.NewDecoder(f)
	for dec.More() {
		var u User
		if err := dec.Decode(&u); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func main() {
	fmt.Println("=== ARCH.3 Hexagonal architecture in Go ===")
	fmt.Println()

	// ---- Adapter 1: InMemoryUserRepository ----
	fmt.Println("--- InMemory adapter ---")
	memRepo := NewInMemoryUserRepository()
	memSvc := NewUserService(memRepo)
	_ = memSvc.Register(1, "Alice", "alice@example.com")
	_ = memSvc.Register(2, "Bob", "bob@example.com")
	users, _ := memSvc.List()
	for _, u := range users {
		fmt.Printf("  %+v\n", u)
	}

	// ---- Adapter 2: FileUserRepository ----
	fmt.Println()
	fmt.Println("--- File adapter ---")
	tmpDir := os.TempDir()
	filePath := filepath.Join(tmpDir, "arch3_users.jsonl")
	_ = os.Remove(filePath) // clean up any previous run
	fileRepo := NewFileUserRepository(filePath)
	fileSvc := NewUserService(fileRepo)
	_ = fileSvc.Register(3, "Carol", "carol@example.com")
	_ = fileSvc.Register(4, "Dave", "dave@example.com")
	users2, _ := fileSvc.List()
	for _, u := range users2 {
		fmt.Printf("  %+v\n", u)
	}

	// ---- Cross-adapter proof ----
	fmt.Println()
	fmt.Println("--- Proof: same UserService works with both adapters ---")
	svcA := NewUserService(memRepo)
	alice, ok := svcA.Lookup(1)
	fmt.Printf("  InMemory  -> id=1: %+v (found=%v)\n", alice, ok)

	svcB := NewUserService(fileRepo)
	dave, ok := svcB.Lookup(4)
	fmt.Printf("  File      -> id=4: %+v (found=%v)\n", dave, ok)

	_ = os.Remove(filePath) // clean up temp file

	fmt.Println()
	fmt.Println("The core (UserService) never imports an adapter. Both storage strategies satisfy the same port contract.")
	fmt.Println()
	fmt.Println("---------------------------------------------------")
	fmt.Println("NEXT UP: ARCH.4 -> 09-architecture/03-architecture-patterns/04-repository-pattern-deep-dive")
	fmt.Println("Current: ARCH.3 (hexagonal architecture in go)")
	fmt.Println("---------------------------------------------------")
}
