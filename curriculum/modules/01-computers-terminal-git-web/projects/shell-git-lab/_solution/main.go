package main

import (
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"
)

type ShellCommand struct {
	Name string
	Args []string
}

type GitCommand struct {
	Name string
	Args []string
}

type File struct {
	Name     string
	IsDir    bool
	Children []*File
	Content  string
}

type Repo struct {
	Files    map[string]string
	Staged   map[string]string
	Commits  []string
	Branches map[string]string
	Head     string
	Current  string
}

func NewRepo() *Repo {
	return &Repo{
		Files:    make(map[string]string),
		Staged:   make(map[string]string),
		Commits:  make([]string, 0),
		Branches: map[string]string{"main": ""},
		Head:     "main",
		Current:  "main",
	}
}

func (r *Repo) GitInit() {
	r.Branches["main"] = ""
	r.Current = "main"
	r.Head = "main"
}

func (r *Repo) GitAdd(filename string) error {
	content, ok := r.Files[filename]
	if !ok {
		return fmt.Errorf("file not found: %s", filename)
	}
	r.Staged[filename] = content
	return nil
}

func (r *Repo) GitCommit(message string) error {
	if len(r.Staged) == 0 {
		return errors.New("nothing staged for commit")
	}
	files := make([]string, 0, len(r.Staged))
	for name := range r.Staged {
		files = append(files, name)
	}
	sort.Strings(files)

	msg := message
	if msg == "" {
		msg = strings.Join(files, ", ")
	}
	r.Commits = append(r.Commits, msg)
	r.Branches[r.Current] = msg
	r.Staged = make(map[string]string)
	return nil
}

func (r *Repo) GitBranch(name string) error {
	if _, ok := r.Branches[name]; ok {
		return fmt.Errorf("branch already exists: %s", name)
	}
	r.Branches[name] = r.Branches[r.Current]
	return nil
}

func (r *Repo) GitCheckout(branch string) error {
	if _, ok := r.Branches[branch]; !ok {
		return fmt.Errorf("branch not found: %s", branch)
	}
	r.Current = branch
	return nil
}

func (r *Repo) GitMerge(branch string) error {
	if _, ok := r.Branches[branch]; !ok {
		return fmt.Errorf("branch not found: %s", branch)
	}
	msg := fmt.Sprintf("Merge branch '%s' into %s", branch, r.Current)
	r.Commits = append(r.Commits, msg)
	r.Branches[r.Current] = msg
	return nil
}

func (r *Repo) ListFiles() []string {
	names := make([]string, 0, len(r.Files))
	for name := range r.Files {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func (r *Repo) FileExists(filename string) bool {
	_, ok := r.Files[filename]
	return ok
}

func (r *Repo) ReadFile(filename string) (string, error) {
	content, ok := r.Files[filename]
	if !ok {
		return "", fmt.Errorf("file not found: %s", filename)
	}
	return content, nil
}

func (r *Repo) CreateFile(name, content string) {
	r.Files[name] = content
}

func main() {
	repo := NewRepo()
	repo.GitInit()

	repo.CreateFile("README.md", "# My Project")
	repo.CreateFile("main.go", "package main\n\nfunc main() {}")

	err := repo.GitAdd("README.md")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error staging file: %v\n", err)
		os.Exit(1)
	}

	err = repo.GitCommit("Initial commit")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error committing: %v\n", err)
		os.Exit(1)
	}

	err = repo.GitBranch("feature")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating branch: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Repository state:")
	fmt.Printf("  Current branch: %s\n", repo.Current)
	fmt.Printf("  Files: %s\n", strings.Join(repo.ListFiles(), ", "))
	fmt.Printf("  Commits: %d\n", len(repo.Commits))
	fmt.Printf("  Branches: ")
	names := make([]string, 0, len(repo.Branches))
	for name := range repo.Branches {
		names = append(names, name)
	}
	fmt.Printf("%s\n", strings.Join(names, ", "))
}
