package main

import (
	"fmt"
	"strings"
)

var commitSeq int

func initRepo() {
	commitSeq = 0
}

func nextHash() string {
	commitSeq++
	return fmt.Sprintf("%08x", commitSeq)
}

func commit(branch, message string) (string, string) {
	hash := nextHash()
	entry := hash + ":" + message
	if branch == "" {
		return entry, hash
	}
	return branch + "," + entry, hash
}

func branchFrom(source string) string {
	return source
}

func logBranch(branch string) []string {
	if branch == "" {
		return []string{}
	}
	parts := strings.Split(branch, ",")
	var result [30]string
	n := 0
	for i := len(parts) - 1; i >= 0; i-- {
		p := parts[i]
		colonAt := 0
		for j := 0; j < len(p); j++ {
			if p[j] == ':' {
				colonAt = j
				break
			}
		}
		if colonAt == 0 {
			continue
		}
		hash := p[:colonAt]
		msg := p[colonAt+1:]
		result[n] = fmt.Sprintf("%s %s", hash[:7], msg)
		n++
	}
	return result[:n]
}

func merge(source, target string) (string, bool) {
	if source == "" {
		return target, false
	}
	if target == "" {
		return source, true
	}
	sourceParts := strings.Split(source, ",")
	targetParts := strings.Split(target, ",")
	targetTipHash := targetParts[len(targetParts)-1][:8]
	ff := false
	for _, p := range sourceParts {
		if len(p) >= 8 && p[:8] == targetTipHash {
			ff = true
			break
		}
	}
	if ff {
		return source, true
	}
	merged := target
	for _, sp := range sourceParts {
		dup := false
		for _, tp := range targetParts {
			if sp == tp {
				dup = true
				break
			}
		}
		if !dup {
			merged = merged + "," + sp
		}
	}
	hash := nextHash()
	merged = merged + "," + hash + ":Merge branch"
	return merged, true
}

func main() {
	initRepo()

	mainBranch, _ := commit("", "Initial commit")
	mainBranch, _ = commit(mainBranch, "Add README")
	featureBranch := branchFrom(mainBranch)
	featureBranch, _ = commit(featureBranch, "Add feature logic")
	featureBranch, _ = commit(featureBranch, "Add feature tests")
	mainBranch, _ = commit(mainBranch, "Update main config")

	fmt.Println("=== Git Branching Simulation ===")
	fmt.Println()

	fmt.Println("Main branch log:")
	for _, l := range logBranch(mainBranch) {
		fmt.Println(" ", l)
	}
	fmt.Println()

	fmt.Println("Feature branch log:")
	for _, l := range logBranch(featureBranch) {
		fmt.Println(" ", l)
	}
	fmt.Println()

	merged, ok := merge(featureBranch, mainBranch)
	if ok {
		fmt.Println("Merged successfully")
	} else {
		fmt.Println("Merge conflict detected (simulated conflict)")
	}
	fmt.Println()

	fmt.Println("Main branch log after merge:")
	for _, l := range logBranch(merged) {
		fmt.Println(" ", l)
	}
}
