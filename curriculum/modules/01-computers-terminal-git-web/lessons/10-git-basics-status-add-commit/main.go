package main

import (
	"fmt"
	"sort"
)

// Repo simulates a basic Git repository with working directory, staging, and history.
type Repo struct {
	Working map[string]string
	staging map[string]string
	History []Commit
}

// Commit represents a saved snapshot.
type Commit struct {
	Message string
	Files   map[string]string
}

// NewRepo creates an empty repository.
func NewRepo() *Repo {
	return &Repo{
		Working: make(map[string]string),
		staging: make(map[string]string),
		History: []Commit{},
	}
}

// StageNames returns sorted names of staged files.
func (r *Repo) StageNames() []string {
	names := make([]string, 0, len(r.staging))
	for name := range r.staging {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// Stage moves a file from working directory to staging area.
func (r *Repo) Stage(name string) error {
	content, ok := r.Working[name]
	if !ok {
		return fmt.Errorf("file not found: %s", name)
	}
	r.staging[name] = content
	return nil
}

// Unstage removes a file from staging area.
func (r *Repo) Unstage(name string) {
	delete(r.staging, name)
}

// Commit saves all staged files to history.
func (r *Repo) Commit(msg string) error {
	if len(r.staging) == 0 {
		return fmt.Errorf("nothing staged")
	}
	files := make(map[string]string)
	for k, v := range r.staging {
		files[k] = v
	}
	r.History = append(r.History, Commit{Message: msg, Files: files})
	r.staging = make(map[string]string)
	return nil
}

func main() {
	repo := NewRepo()
	repo.Working["hello.txt"] = "Hello, Git!"

	fmt.Println("Staged:", repo.StageNames())

	repo.Stage("hello.txt")
	fmt.Println("After stage:", repo.StageNames())

	repo.Commit("feat: add hello.txt")
	fmt.Printf("History: %d commits\n", len(repo.History))

	repo.Working["main.go"] = "package main"
	repo.Stage("main.go")
	repo.Commit("feat: add main.go")
	fmt.Printf("History: %d commits\n", len(repo.History))
}
