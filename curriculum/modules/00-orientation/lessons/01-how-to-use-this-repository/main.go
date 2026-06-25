package main

import (
	"fmt"
)

type DirEntry struct {
	Name     string
	Purpose  string
	Contents []string
}

var repoDirs = []DirEntry{
	{
		Name:    "metadata/",
		Purpose: "defines the curriculum map: what lessons exist, their order, and how they connect",
		Contents: []string{
			"curriculum.json         -- full lesson graph",
			"modules/                -- module definitions",
		},
	},
	{
		Name:    "curriculum/",
		Purpose: "holds learner-facing content: READMEs, code, starter and solution files",
		Contents: []string{
			"modules/00-orientation/ -- Module 0: getting started",
			"modules/01-go-basics/   -- Module 1: Go fundamentals",
		},
	},
	{
		Name:    "tools/",
		Purpose: "validators, auditors, and generators that check curriculum correctness",
		Contents: []string{
			"validate/   -- checks metadata against filesystem",
			"audit/      -- reviews consistency and completeness",
		},
	},
	{
		Name:    "docs/",
		Purpose: "maintainer documentation: style guide, review process, conventions",
		Contents: []string{
			"CONTRIBUTING.md",
			"style-guide.md",
		},
	},
	{
		Name:    "dist/",
		Purpose: "generated output: exports, bundles, release artifacts (git-ignored)",
		Contents: []string{
			"(generated at build time)",
		},
	},
}

func Lookup(name string) *DirEntry {
	for i := range repoDirs {
		if repoDirs[i].Name == name || repoDirs[i].Name == name+"/" {
			return &repoDirs[i]
		}
	}
	return nil
}

func main() {
	fmt.Println("The Go Engineer — Repository Map")
	fmt.Println()
	for _, d := range repoDirs {
		fmt.Printf("  %-14s %s\n", d.Name, d.Purpose)
	}
	fmt.Println()
	fmt.Println("Use Lookup(name) to get details about any top-level directory.")
}
