package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func ValidateRepository(cfg Config) ValidationResult {
	var result ValidationResult
	metadata, loadResult := loadMetadata(cfg)
	result.Merge(loadResult)
	if !loadResult.OK() {
		return result
	}
	result.Merge(validatePhysicalModules(cfg, metadata))
	result.Merge(validatePhysicalItems(cfg, metadata))
	result.Merge(validatePhysicalProjects(cfg, metadata))
	result.Merge(validatePhysicalAssessments(cfg, metadata))
	result.Merge(validateForbiddenNames(cfg))
	return result
}

func validatePhysicalModules(cfg Config, m Metadata) ValidationResult {
	var r ValidationResult
	for _, module := range allModules(m) {
		if module.Path == "" {
			r.Errorf("%s missing path", module.ID)
			continue
		}
		abs := filepath.Join(cfg.Root, filepath.FromSlash(module.Path))
		if cfg.Strict && !dirExists(abs) {
			r.Errorf("%s module directory missing: %s", module.ID, module.Path)
		}
		readme := filepath.Join(abs, "README.md")
		if cfg.Strict && !fileExists(readme) {
			r.Errorf("%s module README missing: %s/README.md", module.ID, module.Path)
		}
	}
	return r
}

func validatePhysicalItems(cfg Config, m Metadata) ValidationResult {
	var r ValidationResult
	for _, item := range allItems(m) {
		if item.Files.ReadmePath == "" {
			r.Errorf("%s missing files.readme_path", item.ID)
			continue
		}
		readme := filepath.Join(cfg.Root, filepath.FromSlash(item.Files.ReadmePath))
		if !canonicalContentPath(item.Files.ReadmePath) {
			r.Errorf("%s readme_path is not canonical: %s", item.ID, item.Files.ReadmePath)
		}
		if cfg.Strict {
			if !fileExists(readme) {
				r.Errorf("%s README missing: %s", item.ID, item.Files.ReadmePath)
			} else {
				validateReadmeContent(item.ID, readme, &r)
			}
			validateCodeFiles(cfg, item, &r)
			validateAssetFiles(cfg, item, &r)
		}
	}
	return r
}

func validateCodeFiles(cfg Config, item Item, r *ValidationResult) {
	if item.Files.MainPath != "" {
		mainPath := filepath.Join(cfg.Root, filepath.FromSlash(item.Files.MainPath))
		if !fileExists(mainPath) {
			r.Errorf("%s main.go missing: %s", item.ID, item.Files.MainPath)
		} else {
			validateGoFile(item.ID, mainPath, false, r)
		}
	}
	if item.Files.TestPath != "" {
		testPath := filepath.Join(cfg.Root, filepath.FromSlash(item.Files.TestPath))
		if !fileExists(testPath) {
			r.Errorf("%s main_test.go missing: %s", item.ID, item.Files.TestPath)
		} else {
			validateGoFile(item.ID, testPath, true, r)
		}
	}
	if item.Files.StarterPath != "" {
		starter := filepath.Join(cfg.Root, filepath.FromSlash(item.Files.StarterPath))
		if !dirExists(starter) {
			r.Errorf("%s starter directory missing: %s", item.ID, item.Files.StarterPath)
		}
	}
	if item.Files.SolutionPath != "" {
		solution := filepath.Join(cfg.Root, filepath.FromSlash(item.Files.SolutionPath))
		if !dirExists(solution) {
			r.Errorf("%s solution directory missing: %s", item.ID, item.Files.SolutionPath)
		}
	}
}

func validateAssetFiles(cfg Config, item Item, r *ValidationResult) {
	if item.Files.AssetsDir != "" {
		assets := filepath.Join(cfg.Root, filepath.FromSlash(item.Files.AssetsDir))
		if !dirExists(assets) {
			r.Errorf("%s assets directory missing: %s", item.ID, item.Files.AssetsDir)
		}
	}
	for _, rel := range item.Files.DiagramPaths {
		if rel == "" {
			continue
		}
		if !fileExists(filepath.Join(cfg.Root, filepath.FromSlash(rel))) {
			r.Errorf("%s diagram missing: %s", item.ID, rel)
		}
	}
}

func validateReadmeContent(owner string, path string, r *ValidationResult) {
	data, err := os.ReadFile(path)
	if err != nil {
		r.Errorf("%s cannot read README: %v", owner, err)
		return
	}
	text := string(data)
	required := []string{"## Mission", "## Prerequisites", "## Mental Model", "## Visual Model", "## Machine View", "## Run Instructions", "## Try It", "## In Production", "## Thinking Questions", "## Next Step"}
	last := -1
	for _, heading := range required {
		idx := strings.Index(text, heading)
		if idx < 0 {
			r.Errorf("%s README missing heading %s", owner, heading)
			continue
		}
		if idx < last {
			r.Errorf("%s README heading out of order: %s", owner, heading)
		}
		last = idx
	}
	lower := strings.ToLower(text)
	forbidden := []string{"todo", "tbd", "placeholder", "lorem ipsum", "coming soon"}
	for _, token := range forbidden {
		if strings.Contains(lower, token) {
			r.Errorf("%s README contains forbidden placeholder token %q", owner, token)
		}
	}
	if !strings.Contains(text, "```") {
		r.Errorf("%s README should include at least one fenced command or code example", owner)
	}
}

func validateGoFile(owner string, path string, isTest bool, r *ValidationResult) {
	data, err := os.ReadFile(path)
	if err != nil {
		r.Errorf("%s cannot read Go file: %v", owner, err)
		return
	}
	text := string(data)
	if !strings.Contains(text, "package ") {
		r.Errorf("%s Go file missing package declaration: %s", owner, path)
	}
	if isTest {
		matched, _ := regexp.MatchString(`func\s+Test[A-Za-z0-9_]+\s*\(`, text)
		if !matched {
			r.Errorf("%s test file must contain at least one Test* function", owner)
		}
	}
	if strings.Contains(strings.ToLower(text), "todo") {
		r.Errorf("%s Go file contains TODO: %s", owner, path)
	}
}

func validatePhysicalProjects(cfg Config, m Metadata) ValidationResult {
	var r ValidationResult
	for _, project := range m.Projects.Projects {
		if project.Files.ReadmePath == "" {
			r.Errorf("%s project missing files.readme_path", project.ID)
			continue
		}
		if !canonicalContentPath(project.Files.ReadmePath) {
			r.Errorf("%s project readme path not canonical: %s", project.ID, project.Files.ReadmePath)
		}
		if cfg.Strict {
			path := filepath.Join(cfg.Root, filepath.FromSlash(project.Files.ReadmePath))
			if !fileExists(path) {
				r.Errorf("%s project README missing: %s", project.ID, project.Files.ReadmePath)
			}
		}
	}
	return r
}

func validatePhysicalAssessments(cfg Config, m Metadata) ValidationResult {
	var r ValidationResult
	for _, assessment := range m.Assessments.Assessments {
		readme := assessment.Files["readme_path"]
		if readme == "" {
			r.Errorf("%s assessment missing files.readme_path", assessment.ID)
			continue
		}
		if !canonicalContentPath(readme) {
			r.Errorf("%s assessment readme path not canonical: %s", assessment.ID, readme)
		}
		if cfg.Strict {
			if !fileExists(filepath.Join(cfg.Root, filepath.FromSlash(readme))) {
				r.Errorf("%s assessment README missing: %s", assessment.ID, readme)
			}
			for _, key := range []string{"questions_path", "answer_key_path", "rubric_path"} {
				if rel := assessment.Files[key]; rel != "" && !fileExists(filepath.Join(cfg.Root, filepath.FromSlash(rel))) {
					r.Errorf("%s assessment %s missing: %s", assessment.ID, key, rel)
				}
			}
		}
	}
	return r
}

func validateForbiddenNames(cfg Config) ValidationResult {
	var r ValidationResult
	forbidden := map[string]bool{"codex": true, "ai": true, "agent": true, "bot": true, "llm": true}
	_ = filepath.WalkDir(cfg.Root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if !d.IsDir() {
			return nil
		}
		base := strings.ToLower(d.Name())
		if forbidden[base] {
			r.Errorf("forbidden folder name %s at %s", d.Name(), path)
		}
		if base == ".git" {
			return filepath.SkipDir
		}
		return nil
	})
	if fileExists(filepath.Join(cfg.Root, "AGENTS.md")) {
		r.Errorf("AGENTS.md must not be kept at root; use tools/authoring or the packaged Skill")
	}
	if fileExists(filepath.Join(cfg.Root, ".env")) {
		r.Errorf(".env must not be committed")
	}
	fmt.Print("")
	return r
}
