package main

import (
	"fmt"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

func resolvePath(p string) string {
	if filepath.IsAbs(p) {
		return p
	}
	cPath := filepath.Join(curriculumDir, p)
	if _, err := os.Stat(cPath); err == nil {
		return cPath
	}
	wPath := filepath.Join(curriculumDir, "..", p)
	if _, err := os.Stat(wPath); err == nil {
		return wPath
	}
	return cPath
}

// validatePhysicalArtifacts verifies that every item in path.core.json and
// path.electives.json has physical files on disk that satisfy minimum size contracts.
func validatePhysicalArtifacts() {
	fmt.Println("=== validate-physical-artifacts ===")
	ok := true

	var core CoreBundle
	var electives ElectiveBundle
	readJSON("path.core.json", &core)
	readJSON("path.electives.json", &electives)

	var allItems []Item
	allItems = append(allItems, core.Items...)
	allItems = append(allItems, electives.Items...)

	totalItems := len(allItems)
	readmeCount := 0

	for _, item := range allItems {
		if item.Files == nil || item.Files.ReadmePath == "" {
			errorf("%s: missing files.readme_path in metadata", item.ID)
			ok = false
			continue
		}

		resolved := resolvePath(item.Files.ReadmePath)
		info, err := os.Stat(resolved)
		if err != nil {
			errorf("%s: physical README not found at '%s'", item.ID, item.Files.ReadmePath)
			ok = false
			continue
		}

		if info.Size() < 50 {
			errorf("%s: physical README '%s' is too small (%d bytes)", item.ID, item.Files.ReadmePath, info.Size())
			ok = false
			continue
		}
		readmeCount++

		// Check main_path if specified
		if item.Files.MainPath != "" {
			mResolved := resolvePath(item.Files.MainPath)
			if _, err := os.Stat(mResolved); err != nil {
				errorf("%s: physical main_path not found at '%s'", item.ID, item.Files.MainPath)
				ok = false
			}
		}

		// Check test_path if specified
		if item.Files.TestPath != "" {
			tResolved := resolvePath(item.Files.TestPath)
			if _, err := os.Stat(tResolved); err != nil {
				errorf("%s: physical test_path not found at '%s'", item.ID, item.Files.TestPath)
				ok = false
			}
		}
	}

	fmt.Printf("  Verified %d/%d physical instructional READMEs on disk.\n", readmeCount, totalItems)
	if ok {
		fmt.Println("  All physical artifact checks passed.")
	} else {
		fmt.Println("  Physical artifact checks completed with warnings/errors.")
	}
}

// validateAssessmentArtifacts verifies that all assessments declared in
// assessments.json and Opslane milestones possess verified physical evaluation artifacts.
func validateAssessmentArtifacts() {
	fmt.Println("=== validate-assessment-artifacts ===")
	ok := true

	var rawAssessments struct {
		Assessments []struct {
			ID       string `json:"id"`
			ModuleID string `json:"module_id"`
			Slug     string `json:"slug"`
			Type     string `json:"type"`
		} `json:"assessments"`
	}
	readJSON("assessments.json", &rawAssessments)

	var core CoreBundle
	readJSON("path.core.json", &core)

	resolveSlug := func(id, slug string) string {
		if slug != "" {
			return slug
		}
		if strings.HasPrefix(id, "assessment-project-") {
			return strings.TrimPrefix(id, "assessment-project-")
		}
		if strings.HasPrefix(id, "assessment-module-") {
			return "checkpoint"
		}
		if id == "assessment-go-setup-tooling-tooling-failure-lab" {
			return "tooling-failure-lab"
		}
		return strings.TrimPrefix(id, "assessment-")
	}

	modDirMap := make(map[string]string)
	for _, m := range core.Modules {
		modDirMap[m.ID] = m.Path
	}

	verifiedCount := 0
	for _, a := range rawAssessments.Assessments {
		slug := resolveSlug(a.ID, a.Slug)

		var assessDir string
		if a.ModuleID != "" && modDirMap[a.ModuleID] != "" {
			mPath := modDirMap[a.ModuleID]
			candidates := []string{
				filepath.Join(curriculumDir, mPath, "assessments", slug),
				filepath.Join(curriculumDir, "modules", filepath.Base(mPath), "assessments", slug),
			}
			for _, c := range candidates {
				if info, err := os.Stat(c); err == nil && info.IsDir() {
					assessDir = c
					break
				}
			}
		}

		if assessDir == "" {
			// Fallback search across electives and other modules
			candidates := []string{
				filepath.Join(curriculumDir, "electives", "advanced-electives", "assessments", slug),
			}
			for _, m := range core.Modules {
				candidates = append(candidates,
					filepath.Join(curriculumDir, m.Path, "assessments", slug),
					filepath.Join(curriculumDir, "modules", filepath.Base(m.Path), "assessments", slug),
				)
			}
			for _, c := range candidates {
				if info, err := os.Stat(c); err == nil && info.IsDir() {
					assessDir = c
					break
				}
			}
		}

		if assessDir == "" {
			errorf("%s: assessment directory not found for slug '%s'", a.ID, slug)
			ok = false
			continue
		}

		// Check README exists and is non-trivial (>100 bytes)
		rPath := filepath.Join(assessDir, "README.md")
		rInfo, err := os.Stat(rPath)
		if err != nil {
			errorf("%s: missing README.md in '%s'", a.ID, assessDir)
			ok = false
			continue
		}
		if rInfo.Size() < 100 {
			errorf("%s: README.md in '%s' is too small (%d bytes)", a.ID, assessDir, rInfo.Size())
			ok = false
		}

		// For module checkpoints, labs, and capstones, require full 4-file suite
		isCheckpointOrLab := a.Type == "checkpoint" ||
			strings.Contains(a.ID, "checkpoint") ||
			strings.HasPrefix(a.ID, "assessment-module-") ||
			slug == "tooling-failure-lab" ||
			slug == "portfolio-review-pack" ||
			slug == "opslane"

		if isCheckpointOrLab {
			for _, rf := range []string{"questions.md", "rubric.md", "answer-key.md"} {
				fPath := filepath.Join(assessDir, rf)
				info, err := os.Stat(fPath)
				if err != nil {
					errorf("%s: missing required assessment file '%s'", a.ID, rf)
					ok = false
				} else if info.Size() < 50 {
					errorf("%s: assessment file '%s' is too small (%d bytes)", a.ID, rf, info.Size())
					ok = false
				}
			}
		}

		verifiedCount++
	}

	// Also verify all Opslane milestones in module-18 have full 4-file suites.
	opslaneDir := filepath.Join(curriculumDir, "modules", "18-flagship-opslane", "assessments")
	entries, err := os.ReadDir(opslaneDir)
	opslaneCount := 0
	if err != nil {
		errorf("cannot read Opslane milestone directory '%s': %v", opslaneDir, err)
		ok = false
	} else {
		for _, e := range entries {
			if !e.IsDir() || (!strings.HasPrefix(e.Name(), "0") && !strings.HasPrefix(e.Name(), "1")) {
				continue
			}
			milestoneDir := filepath.Join(opslaneDir, e.Name())
			for _, rf := range []string{"README.md", "questions.md", "rubric.md", "answer-key.md"} {
				fPath := filepath.Join(milestoneDir, rf)
				info, err := os.Stat(fPath)
				if err != nil {
					errorf("opslane milestone '%s': missing '%s'", e.Name(), rf)
					ok = false
				} else if info.Size() < 50 {
					errorf("opslane milestone '%s': '%s' too small (%d bytes)", e.Name(), rf, info.Size())
					ok = false
				}
			}
			opslaneCount++
		}
		if expected := loadExpected().Inventory.OpslaneMilestones; expected != 0 && opslaneCount != expected {
			errorf("expected %d Opslane milestones, verified %d", expected, opslaneCount)
			ok = false
		}
	}

	fmt.Printf("  Verified %d module/project assessments + %d Opslane milestones have verified physical artifacts.\n", verifiedCount, opslaneCount)
	if ok {
		fmt.Println("  All assessment artifact checks passed.")
	} else {
		fmt.Println("  Assessment artifact checks completed with warnings/errors.")
	}
}

// validateStarterParity verifies that projects which ship starter/solution
// code have valid Go implementations that parse cleanly and declare both
// production and test Go files. Projects that ship no code artifacts at all
// (environment and shell-workflow labs) are exempt from the parity contract,
// but the number of code-bearing projects is checked against expected.json so
// a silently dropped starter or solution is still caught.
func validateStarterParity() {
	fmt.Println("=== validate-starter-parity ===")
	ok := true

	var rawProjects struct {
		Projects []struct {
			ID       string `json:"id"`
			ModuleID string `json:"module_id"`
			Slug     string `json:"slug"`
		} `json:"projects"`
	}
	readJSON("projects.json", &rawProjects)

	var core CoreBundle
	readJSON("path.core.json", &core)
	var electives ElectiveBundle
	readJSON("path.electives.json", &electives)

	projectRoots := projectDirectories(core.Modules, electives.Modules)

	projectCount := 0
	fset := token.NewFileSet()

	for _, p := range rawProjects.Projects {
		projDir := locateProjectDir(projectRoots, p.Slug)
		if projDir == "" {
			errorf("%s: project directory not found for slug '%s'", p.ID, p.Slug)
			ok = false
			continue
		}

		starterDir := filepath.Join(projDir, "_starter")
		solutionDir := filepath.Join(projDir, "_solution")
		hasStarter := isDirectory(starterDir)
		hasSolution := isDirectory(solutionDir)

		if !hasStarter && !hasSolution {
			// No code artifacts shipped (e.g. environment/shell labs); these
			// projects are not part of the starter-parity inventory.
			continue
		}
		if hasStarter != hasSolution {
			errorf("%s: incomplete starter/solution pair (_starter=%t, _solution=%t)", p.ID, hasStarter, hasSolution)
			ok = false
			continue
		}

		for _, side := range []struct {
			label string
			dir   string
		}{
			{"_starter", starterDir},
			{"_solution", solutionDir},
		} {
			goFiles, testFiles, parseIssues, err := inspectGoArtifacts(side.dir, fset)
			if err != nil {
				errorf("%s: cannot read %s directory: %v", p.ID, side.label, err)
				ok = false
				continue
			}
			for _, issue := range parseIssues {
				errorf("%s: %s %s", p.ID, side.label, issue)
				ok = false
			}
			if goFiles == 0 {
				errorf("%s: %s directory has no .go files", p.ID, side.label)
				ok = false
			}
			if testFiles == 0 {
				errorf("%s: %s directory declares no Go test suite (no _test.go files)", p.ID, side.label)
				ok = false
			}
		}

		projectCount++
	}

	if expected := loadExpected().Inventory.StarterParityProjects; expected != 0 && projectCount != expected {
		errorf("expected %d starter-parity projects, verified %d", expected, projectCount)
		ok = false
	}

	fmt.Printf("  Verified Go code parity and AST validity across %d projects.\n", projectCount)
	if ok {
		fmt.Println("  All starter parity checks passed.")
	} else {
		fmt.Println("  Starter parity checks completed with warnings/errors.")
	}
}

// isDirectory reports whether path exists and is a directory.
func isDirectory(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

// inspectGoArtifacts counts Go source and _test.go files directly inside dir and
// returns a human-readable issue for each .go file that fails to parse. A
// non-nil error means the directory itself could not be read.
func inspectGoArtifacts(dir string, fset *token.FileSet) (goFiles, testFiles int, parseIssues []string, err error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return 0, 0, nil, err
	}
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".go") {
			continue
		}
		goFiles++
		if strings.HasSuffix(e.Name(), "_test.go") {
			testFiles++
		}
		fullPath := filepath.Join(dir, e.Name())
		if _, parseErr := parser.ParseFile(fset, fullPath, nil, parser.AllErrors); parseErr != nil {
			parseIssues = append(parseIssues, fmt.Sprintf("file '%s' has syntax errors: %v", e.Name(), parseErr))
		}
	}
	return goFiles, testFiles, parseIssues, nil
}

// projectDirectories returns the candidate project roots for every module,
// including elective modules, handling both absolute module paths and the
// legacy bare "<module>/projects/<slug>" layout.
func projectDirectories(coreModules, electiveModules []Module) []string {
	var roots []string
	for _, m := range append(append([]Module{}, coreModules...), electiveModules...) {
		roots = append(roots,
			filepath.Join(curriculumDir, m.Path, "projects"),
			filepath.Join(curriculumDir, "modules", filepath.Base(m.Path), "projects"),
		)
	}
	return roots
}

// locateProjectDir returns the first existing "<root>/<slug>" directory, or ""
// when no project root contains the slug.
func locateProjectDir(roots []string, slug string) string {
	for _, root := range roots {
		candidate := filepath.Join(root, slug)
		if isDirectory(candidate) {
			return candidate
		}
	}
	return ""
}
