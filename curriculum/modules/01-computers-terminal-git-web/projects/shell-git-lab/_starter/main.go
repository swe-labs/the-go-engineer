package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

// ShellCommand represents a shell command to execute.
type ShellCommand struct {
	Name string
	Args []string
}

// GitCommand represents a Git command to execute.
type GitCommand struct {
	Name string
	Args []string
}

// File represents a simulated file or directory.
type File struct {
	Name     string
	IsDir    bool
	Children []*File
	Content  string
}

// Repo represents a simulated Git repository.
type Repo struct {
	Files    map[string]string
	Staged   map[string]string
	Commits  []string
	Branches map[string]string
	Head     string
	Current  string
}

// NewRepo creates a new simulated Git repository.
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

// GitInit initializes the repository.
func (r *Repo) GitInit() {
	r.Branches["main"] = ""
	r.Current = "main"
	r.Head = "main"
}

// TODO: Implement GitAdd to stage a file by name.
// It should store the file content in r.Staged if the file exists in r.Files.
// Return an error if the file does not exist.
func (r *Repo) GitAdd(filename string) error {
	// TODO: Implement this function
	return nil
}

// TODO: Implement GitCommit to create a new commit from staged files.
// It should:
//   - Return an error if nothing is staged
//   - Create a commit message combining staged file names
//   - Append the commit message to r.Commits
//   - Update r.Branches[r.Current] to point to the new commit
//   - Clear r.Staged
func (r *Repo) GitCommit(message string) error {
	// TODO: Implement this function
	return nil
}

// TODO: Implement GitBranch to create a new branch at the current commit.
// It should create a new entry in r.Branches pointing to the same commit as r.Current.
// Return an error if the branch already exists.
func (r *Repo) GitBranch(name string) error {
	// TODO: Implement this function
	return nil
}

// TODO: Implement GitCheckout to switch to an existing branch.
// It should update r.Current to the given branch name.
// Return an error if the branch does not exist.
func (r *Repo) GitCheckout(branch string) error {
	// TODO: Implement this function
	return nil
}

// TODO: Implement GitMerge to merge a branch into the current branch.
// It should:
//   - Return an error if the branch does not exist
//   - Create a merge commit message
//   - Append the merge commit to r.Commits
//   - Update the current branch's commit pointer
func (r *Repo) GitMerge(branch string) error {
	// TODO: Implement this function
	return nil
}

// TODO: Implement ListFiles to return a sorted list of file names in r.Files.
func (r *Repo) ListFiles() []string {
	// TODO: Implement this function
	return nil
}

// TODO: Implement FileExists to check if a file exists in r.Files.
func (r *Repo) FileExists(filename string) bool {
	// TODO: Implement this function
	return false
}

// TODO: Implement ReadFile to return the content of a file from r.Files.
// Return an error if the file does not exist.
func (r *Repo) ReadFile(filename string) (string, error) {
	// TODO: Implement this function
	return "", errors.New("not implemented")
}

// CreateFile adds a file to the repository's working directory.
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
