package main

import (
	"fmt"
	"strings"
)

type Commit struct {
	Hash    string
	Message string
	Parent  string
}

type Branch struct {
	Name   string
	Target string
}

type SimRepo struct {
	Commits map[string]Commit
	Branches map[string]Branch
	Head     string
}

func NewSimRepo() *SimRepo {
	return &SimRepo{
		Commits:  make(map[string]Commit),
		Branches: make(map[string]Branch),
	}
}

func (r *SimRepo) Commit(branch, message string) string {
	hash := fmt.Sprintf("%08x", len(r.Commits)+1)
	parent := ""
	if b, ok := r.Branches[branch]; ok {
		parent = b.Target
	}
	r.Commits[hash] = Commit{Hash: hash, Message: message, Parent: parent}
	r.Branches[branch] = Branch{Name: branch, Target: hash}
	r.Head = hash
	return hash
}

func (r *SimRepo) Branch(name, fromBranch string) {
	target := r.Branches[fromBranch].Target
	r.Branches[name] = Branch{Name: name, Target: target}
}

func (r *SimRepo) Log(branch string) []string {
	var log []string
	hash := r.Branches[branch].Target
	for hash != "" {
		c, ok := r.Commits[hash]
		if !ok {
			break
		}
		log = append(log, fmt.Sprintf("%s %s", c.Hash[:7], c.Message))
		hash = c.Parent
	}
	return log
}

func (r *SimRepo) Merge(from, into string) (string, bool) {
	fromCommit := r.Branches[from].Target
	intoCommit := r.Branches[into].Target

	ancestor := findMergeBase(r.Commits, fromCommit, intoCommit)
	if ancestor == "" || ancestor == intoCommit {
		r.Branches[into] = Branch{Name: into, Target: fromCommit}
		return "", true
	}

	mergeHash := fmt.Sprintf("%08x", len(r.Commits)+1)
	r.Commits[mergeHash] = Commit{Hash: mergeHash, Message: fmt.Sprintf("Merge branch '%s' into %s", from, into), Parent: intoCommit + "," + fromCommit}
	r.Branches[into] = Branch{Name: into, Target: mergeHash}
	r.Head = mergeHash
	conflict := strings.Contains(r.Commits[fromCommit].Message, "conflict") && strings.Contains(r.Commits[intoCommit].Message, "conflict")
	return mergeHash, !conflict
}

func findMergeBase(commits map[string]Commit, a, b string) string {
	ancestorsA := map[string]bool{}
	for h := a; h != ""; {
		ancestorsA[h] = true
		c, ok := commits[h]
		if !ok {
			break
		}
		if c.Parent == "" {
			break
		}
		parents := strings.Split(c.Parent, ",")
		h = parents[0]
	}
	for h := b; h != ""; {
		if ancestorsA[h] {
			return h
		}
		c, ok := commits[h]
		if !ok {
			break
		}
		if c.Parent == "" {
			break
		}
		parents := strings.Split(c.Parent, ",")
		h = parents[0]
	}
	return ""
}

func main() {
	repo := NewSimRepo()

	repo.Commit("main", "Initial commit")
	repo.Commit("main", "Add README")
	repo.Branch("feature", "main")
	repo.Commit("feature", "Add feature logic")
	repo.Commit("feature", "Add feature tests")
	repo.Commit("main", "Update main config")

	fmt.Println("=== Git Branching Simulation ===")
	fmt.Println()

	fmt.Println("Main branch log:")
	for _, l := range repo.Log("main") {
		fmt.Println(" ", l)
	}
	fmt.Println()

	fmt.Println("Feature branch log:")
	for _, l := range repo.Log("feature") {
		fmt.Println(" ", l)
	}
	fmt.Println()

	hash, ok := repo.Merge("feature", "main")
	if ok {
		fmt.Printf("Merged successfully (merge commit: %s)\n", hash[:7])
	} else {
		fmt.Println("Merge conflict detected (simulated conflict)")
	}

	fmt.Println()
	fmt.Println("Main branch log after merge:")
	for _, l := range repo.Log("main") {
		fmt.Println(" ", l)
	}
}
