package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(2)
	}

	command := os.Args[1]
	fs := flag.NewFlagSet(command, flag.ExitOnError)
	rootFlag := fs.String("root", "", "repository root; defaults to nearest ancestor containing metadata/path.core.json")
	strictFlag := fs.Bool("strict-repository", false, "require all learner-facing files and strict content checks")
	_ = fs.Parse(os.Args[2:])

	root, err := discoverRoot(*rootFlag)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(2)
	}
	cfg := Config{
		Root:          root,
		MetadataDir:   filepath.Join(root, "metadata"),
		CurriculumDir: filepath.Join(root, "curriculum"),
		Strict:        *strictFlag,
	}

	var result ValidationResult
	switch command {
	case "validate-metadata":
		result = ValidateMetadata(cfg)
	case "validate-repository":
		result = ValidateRepository(cfg)
	case "validate-all":
		result = ValidateMetadata(cfg)
		if result.OK() {
			result.Merge(ValidateRepository(cfg))
		}
	case "help", "-h", "--help":
		printUsage()
		return
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", command)
		printUsage()
		os.Exit(2)
	}

	result.Print()
	if !result.OK() {
		os.Exit(1)
	}
	fmt.Println("VALIDATION PASSED")
}

func printUsage() {
	fmt.Println(`usage:
  go run ./tools/validate/curriculum validate-metadata [--root .]
  go run ./tools/validate/curriculum validate-repository [--root .] [--strict-repository]
  go run ./tools/validate/curriculum validate-all [--root .] [--strict-repository]

commands:
  validate-metadata     validate metadata graph, contracts, concepts, projects, assessments, and migration
  validate-repository   validate learner-facing files under curriculum/
  validate-all          run metadata first, then repository validation`)
}
