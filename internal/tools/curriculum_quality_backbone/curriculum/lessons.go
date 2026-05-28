package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

func ValidateRepositoryContent(cur *Curriculum, report *Report, cfg Config) {
	ValidateDeclaredFilesExist(cur, report, cfg)
	ValidateReadmeQuality(cur, report, cfg)
	ValidateCodeQuality(cur, report, cfg)
	ValidateAssetQuality(cur, report, cfg)
}

func ValidateDeclaredFilesExist(cur *Curriculum, report *Report, cfg Config) {
	check := "repository-files"
	for _, item := range cur.Items {
		id := str(item, "id")
		files := obj(item, "files")
		if files == nil {
			report.Error(check, id, "missing files block")
			continue
		}
		mustExist := []string{"readme_path"}
		cc := obj(item, "content_contract")
		if boolVal(cc, "runnable_required") {
			mustExist = append(mustExist, "main_path")
		}
		if boolVal(cc, "tests_required") {
			mustExist = append(mustExist, "test_path")
		}
		if boolVal(cc, "visual_model_required") {
			mustExist = append(mustExist, "assets_dir")
		}
		// Strict repository mode requires every declared non-empty file path to exist.
		if cfg.StrictRepository {
			for _, key := range []string{"main_path", "test_path", "starter_path", "solution_path", "assets_dir"} {
				if str(files, key) != "" {
					mustExist = append(mustExist, key)
				}
			}
		}
		seen := map[string]bool{}
		for _, key := range mustExist {
			if seen[key] {
				continue
			}
			seen[key] = true
			rel := str(files, key)
			if rel == "" {
				report.Error(check, id, "%s is required but empty", key)
				continue
			}
			abs := filepath.Join(cfg.Root, filepath.FromSlash(rel))
			if strings.HasSuffix(key, "dir") {
				if !dirExists(abs) {
					report.Error(check, id, "%s directory not found: %s", key, rel)
				}
			} else if !fileExists(abs) {
				report.Error(check, id, "%s file not found: %s", key, rel)
			}
		}
		for _, rv := range list(files, "diagram_paths") {
			rel, _ := rv.(string)
			if rel == "" {
				continue
			}
			if !fileExists(filepath.Join(cfg.Root, filepath.FromSlash(rel))) {
				report.Error(check, id, "diagram file not found: %s", rel)
			}
		}
	}
}

func ValidateReadmeQuality(cur *Curriculum, report *Report, cfg Config) {
	check := "readme-quality"
	requiredHeadings := []string{
		"# ",
		"## Problem",
		"## Mental Model",
		"## Why It Matters",
		"## How It Works",
		"## Common Mistakes",
		"## Practice",
		"## Review Questions",
		"NEXT UP:",
	}
	for _, item := range cur.Items {
		id := str(item, "id")
		files := obj(item, "files")
		rel := str(files, "readme_path")
		if rel == "" {
			continue
		}
		abs := filepath.Join(cfg.Root, filepath.FromSlash(rel))
		data, err := os.ReadFile(abs)
		if err != nil {
			if cfg.StrictRepository {
				report.Error(check, id, "cannot read README: %s", rel)
			}
			continue
		}
		content := string(data)
		if len(strings.Fields(content)) < 400 {
			report.Error(check, id, "README is too short for world-class lesson documentation: %d words", len(strings.Fields(content)))
		}
		for _, h := range requiredHeadings {
			if !strings.Contains(content, h) {
				report.Error(check, id, "README missing required section or footer %q", h)
			}
		}
		if !containsAny(content, []string{"analogy", "like a", "think of", "mental model"}) {
			report.Error(check, id, "README needs an analogy or learner-friendly mental model")
		}
		if !containsAny(content, []string{"production", "real world", "at work", "in a service"}) {
			report.Error(check, id, "README needs real-world/production context")
		}
		if !containsAny(content, []string{"```go", "go run", "go test"}) {
			report.Error(check, id, "README should include Go code or runnable commands")
		}
		scanStrictAuthoredContent(report, check, id, content)
		// Validate NEXT UP points to one of the declared next item IDs.
		if nextIDs := stringsList(item, "next_item_ids"); len(nextIDs) > 0 {
			ok := false
			for _, nid := range nextIDs {
				next := cur.ItemByID[nid]
				if strings.Contains(content, nid) || strings.Contains(content, str(next, "title")) {
					ok = true
					break
				}
			}
			if !ok {
				report.Error(check, id, "NEXT UP footer must reference one of next_item_ids %v", nextIDs)
			}
		}
	}
}

func ValidateCodeQuality(cur *Curriculum, report *Report, cfg Config) {
	check := "code-quality"
	goFiles := map[string]string{}
	for _, item := range cur.Items {
		id := str(item, "id")
		files := obj(item, "files")
		for _, key := range []string{"main_path", "test_path"} {
			rel := str(files, key)
			if rel == "" || !strings.HasSuffix(rel, ".go") {
				continue
			}
			abs := filepath.Join(cfg.Root, filepath.FromSlash(rel))
			if !fileExists(abs) {
				continue
			}
			goFiles[abs] = id
			data, err := os.ReadFile(abs)
			if err != nil {
				report.Error(check, id, "cannot read %s", rel)
				continue
			}
			content := string(data)
			if strings.Contains(content, "panic(\"TODO") || strings.Contains(content, "TODO:") {
				report.Error(check, id, "Go file contains TODO placeholder: %s", rel)
			}
			if strings.Contains(content, "fmt.Println(\"TODO") {
				report.Error(check, id, "Go file contains placeholder print: %s", rel)
			}
			if strings.HasSuffix(rel, "_test.go") || key == "test_path" {
				if !regexp.MustCompile(`func\s+Test[A-Za-z0-9_]+\s*\(`).MatchString(content) {
					report.Error(check, id, "test file has no Test* function: %s", rel)
				}
			}
		}
		verification := obj(item, "verification")
		if boolVal(obj(item, "content_contract"), "runnable_required") && str(verification, "run_command") == "" {
			report.Error(check, id, "runnable_required but verification.run_command is empty")
		}
		if boolVal(obj(item, "content_contract"), "tests_required") && str(verification, "test_command") == "" {
			report.Error(check, id, "tests_required but verification.test_command is empty")
		}
	}
	if len(goFiles) == 0 {
		return
	}
	// gofmt check on all actual Go files.
	paths := make([]string, 0, len(goFiles))
	for p := range goFiles {
		paths = append(paths, p)
	}
	cmd := exec.Command("gofmt", append([]string{"-l"}, paths...)...)
	var out bytes.Buffer
	cmd.Stdout = &out
	_ = cmd.Run()
	for _, line := range strings.Split(strings.TrimSpace(out.String()), "\n") {
		if strings.TrimSpace(line) == "" {
			continue
		}
		report.Error(check, goFiles[line], "gofmt required for %s", line)
	}
}

func ValidateAssetQuality(cur *Curriculum, report *Report, cfg Config) {
	check := "asset-quality"
	for _, item := range cur.Items {
		id := str(item, "id")
		cc := obj(item, "content_contract")
		files := obj(item, "files")
		if !boolVal(cc, "visual_model_required") && !cfg.StrictRepository {
			continue
		}
		assetRel := str(files, "assets_dir")
		if assetRel == "" {
			report.Error(check, id, "visual model required but assets_dir is empty")
			continue
		}
		assetAbs := filepath.Join(cfg.Root, filepath.FromSlash(assetRel))
		if !dirExists(assetAbs) {
			report.Error(check, id, "assets_dir not found: %s", assetRel)
			continue
		}
		entries, _ := os.ReadDir(assetAbs)
		foundVisual := false
		for _, e := range entries {
			name := strings.ToLower(e.Name())
			if strings.HasSuffix(name, ".png") || strings.HasSuffix(name, ".jpg") || strings.HasSuffix(name, ".jpeg") || strings.HasSuffix(name, ".svg") || strings.HasSuffix(name, ".mermaid") || strings.HasSuffix(name, ".mmd") {
				foundVisual = true
			}
		}
		if !foundVisual {
			report.Error(check, id, "assets_dir has no visual asset (.png/.jpg/.svg/.mermaid/.mmd)")
		}
	}
}

func containsAny(s string, needles []string) bool {
	lower := strings.ToLower(s)
	for _, n := range needles {
		if strings.Contains(lower, strings.ToLower(n)) {
			return true
		}
	}
	return false
}

func runCommandForCI(command string, root string) error {
	parts := strings.Fields(command)
	if len(parts) == 0 {
		return nil
	}
	cmd := exec.Command(parts[0], parts[1:]...)
	cmd.Dir = root
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s failed: %s", command, string(out))
	}
	return nil
}
