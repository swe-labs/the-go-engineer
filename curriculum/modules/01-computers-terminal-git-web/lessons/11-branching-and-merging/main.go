package main

import (
	"fmt"
)

// Commit represents a node in the Git commit graph.
type Commit struct {
	Hash    string
	Message string
	Parent  *Commit
	Changes map[string]string
}

// Branch is a named pointer to a specific commit.
type Branch struct {
	Name string
	Tip  *Commit
}

func main() {
	fmt.Println("=== Git Branching and 3-Way Merge Simulation ===")

	// 1. Initial base commit on main
	c0 := &Commit{
		Hash:    "a1b2c3d",
		Message: "Initial commit",
		Changes: map[string]string{"README.md": "v1.0", "app.go": "package main"},
	}
	mainBranch := &Branch{Name: "main", Tip: c0}
	fmt.Printf("Initial state: branch %s at commit [%s] %q\n", mainBranch.Name, c0.Hash[:7], c0.Message)

	// 2. Branch off feature branch from c0
	featureBranch := &Branch{Name: "feature-login", Tip: mainBranch.Tip}
	fmt.Printf("Created branch %s at commit [%s]\n", featureBranch.Name, featureBranch.Tip.Hash[:7])

	// 3. Commit on main
	c1 := &Commit{
		Hash:    "f4e5d6c",
		Message: "Update app.go on main",
		Parent:  c0,
		Changes: map[string]string{"README.md": "v1.0", "app.go": "package main // v1.1"},
	}
	mainBranch.Tip = c1

	// 4. Commit on feature branch
	c2 := &Commit{
		Hash:    "9z8y7x6",
		Message: "Add login handler on feature",
		Parent:  c0,
		Changes: map[string]string{"README.md": "v1.0", "login.go": "func Login() {}"},
	}
	featureBranch.Tip = c2

	fmt.Printf("Main advanced to [%s] %q\n", mainBranch.Tip.Hash[:7], mainBranch.Tip.Message)
	fmt.Printf("Feature advanced to [%s] %q\n", featureBranch.Tip.Hash[:7], featureBranch.Tip.Message)

	// 5. 3-Way Merge Simulation
	fmt.Println("\nSimulating merge of feature-login into main...")
	base := c0
	fmt.Printf("Common Ancestor (Merge Base): [%s]\n", base.Hash[:7])

	hasConflict := false
	for file := range mainBranch.Tip.Changes {
		if _, touchedByFeature := featureBranch.Tip.Changes[file]; touchedByFeature && file != "README.md" {
			hasConflict = true
		}
	}

	if !hasConflict {
		mergeCommit := &Commit{
			Hash:    "m3r9e01",
			Message: fmt.Sprintf("Merge branch '%s' into %s", featureBranch.Name, mainBranch.Name),
			Parent:  mainBranch.Tip,
		}
		mainBranch.Tip = mergeCommit
		fmt.Printf("Merge successful! Created merge commit [%s]: %q\n", mergeCommit.Hash[:7], mergeCommit.Message)
	} else {
		fmt.Println("Conflict detected! Manual resolution required.")
	}
}
