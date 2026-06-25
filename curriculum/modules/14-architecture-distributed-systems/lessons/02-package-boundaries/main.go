package main

import (
	"fmt"
)

type ExportedService struct {
	repo Repository
}

func NewExportedService(repo Repository) *ExportedService {
	return &ExportedService{repo: repo}
}

type Repository interface {
	Save(item string) error
	Find(id string) (string, error)
}

type inMemoryStore struct {
	data map[string]string
}

func newInMemoryStore() *inMemoryStore {
	return &inMemoryStore{data: make(map[string]string)}
}

func (s *inMemoryStore) Save(item string) error {
	s.data[item] = item
	return nil
}

func (s *inMemoryStore) Find(id string) (string, error) {
	val, ok := s.data[id]
	if !ok {
		return "", fmt.Errorf("not found: %s", id)
	}
	return val, nil
}

func main() {
	repo := newInMemoryStore()
	svc := NewExportedService(repo)

	_ = svc.repo.Save("user-1")
	result, _ := svc.repo.Find("user-1")
	fmt.Println("Found:", result)

	_, err := svc.repo.Find("user-2")
	fmt.Println("Error:", err)
}
