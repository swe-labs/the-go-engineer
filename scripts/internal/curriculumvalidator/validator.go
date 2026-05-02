package curriculumvalidator

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type V2Section struct {
	ID            string   `json:"id"`
	Number        string   `json:"number"`
	Slug          string   `json:"slug"`
	Title         string   `json:"title"`
	PathPrefix    string   `json:"path_prefix"`
	Phase         string   `json:"phase"`
	Summary       string   `json:"summary"`
	Status        string   `json:"status"`
	EntryPoints   []string `json:"entry_points"`
	Outputs       []string `json:"outputs"`
	Prerequisites []string `json:"prerequisites"`
}

type V2Item struct {
	ID               string   `json:"id"`
	SectionID        string   `json:"section_id"`
	Slug             string   `json:"slug"`
	Title            string   `json:"title"`
	Type             string   `json:"type"`
	Subtype          string   `json:"subtype"`
	Level            string   `json:"level"`
	Status           string   `json:"status"`
	VerificationMode string   `json:"verification_mode"`
	Path             string   `json:"path"`
	Prerequisites    []string `json:"prerequisites"`
	RunCommand       string   `json:"run_command"`
	TestCommand      string   `json:"test_command"`
	StarterPath      string   `json:"starter_path"`
	NextItems        []string `json:"next_items"`
}

type V2Curriculum struct {
	SchemaVersion int         `json:"schema_version"`
	Sections      []V2Section `json:"sections"`
	Items         []V2Item    `json:"items"`
}

type Result struct {
	FilesScanned     int
	V2SectionCount   int
	V2ItemCount      int
	PlaceholderCount int
	HasV2            bool
	ErrorCount       int
}

var runPathPattern = regexp.MustCompile(`\./[A-Za-z0-9._/\-]+(?:/\.\.\.)?`)
var nextUpFooterPattern = regexp.MustCompile(`NEXT UP:\s*([A-Z]{2,6}\.\d+)\s*->\s*([A-Za-z0-9._/\-]+)`)
var markdownLinkPattern = regexp.MustCompile(`\[[^\]]+\]\(([^)]+)\)`)
var flagshipPrefixPattern = regexp.MustCompile(`^[A-Z]{3,6}$`)

const expectedSchemaVersion = 1

type canonicalSectionContract struct {
	ID          string
	Number      string
	Slug        string
	Title       string
	PathPrefix  string
	Phase       string
	Status      string
	EntryPoints []string
	Outputs     []string
}

var (
	allowedItemTypes = map[string]bool{
		"lesson":       true,
		"drill":        true,
		"exercise":     true,
		"checkpoint":   true,
		"mini_project": true,
		"capstone":     true,
		"reference":    true,
	}
	allowedLessonSubtypes = map[string]bool{
		"concept":     true,
		"pattern":     true,
		"integration": true,
	}
	allowedLevels = map[string]bool{
		"foundation": true,
		"core":       true,
		"stretch":    true,
		"production": true,
	}
	levelDisplayLabels = map[string]string{
		"foundation": "Foundation",
		"core":       "Core",
		"production": "Production",
		"stretch":    "Stretch",
	}
	allowedVerificationModes = map[string]bool{
		"run":    true,
		"test":   true,
		"rubric": true,
		"mixed":  true,
	}
	allowedSectionStatuses = map[string]bool{
		"stable": true,
	}
	allowedSectionPhases = map[string]bool{
		"foundations":      true,
		"engineering-core": true,
		"systems":          true,
	}
	canonicalV2Sections = []canonicalSectionContract{
		canonicalSection("s00", "00", "how-computers-work", "How Computers Work", "00-how-computers-work", "foundations", []string{"HC.1"}, []string{"HC.5"}),
		canonicalSection("s01", "01", "getting-started", "Getting Started", "01-getting-started", "foundations", []string{"GT.1"}, []string{"GT.6"}),
		canonicalSection("s02", "02", "language-basics", "Language Basics", "02-language-basics", "foundations", []string{"LB.1"}, []string{"LB.4", "CF.7", "DS.6"}),
		canonicalSection("s03", "03", "functions-errors", "Functions and Errors", "03-functions-errors", "foundations", []string{"FE.1"}, []string{"FE.10"}),
		canonicalSection("s04", "04", "types-design", "Types and Design", "04-types-design", "foundations", []string{"TI.1", "CO.1", "ST.1"}, []string{"TI.15", "CO.3", "ST.6"}),
		canonicalSection("s05", "05", "packages-io", "Packages and I/O", "05-packages-io", "engineering-core", []string{"MP.1", "CL.1", "EN.1", "FS.1"}, []string{"MP.4", "CL.4", "EN.6", "FS.8"}),
		canonicalSection("s06", "06", "backend-db", "Backend, APIs & Databases", "06-backend-db", "engineering-core", []string{"HS.1", "API.1", "DB.1"}, []string{"HS.10", "API.9", "DB.8"}),
		canonicalSection("s07", "07", "concurrency", "Concurrency", "07-concurrency", "engineering-core", []string{"GC.0", "SY.1", "CT.1", "TM.1", "CP.1"}, []string{"GC.7", "SY.6", "CT.5", "TM.7", "CP.5"}),
		canonicalSection("s08", "08", "quality-test", "Quality & Testing", "08-quality-test", "engineering-core", []string{"TE.1", "PR.1"}, []string{"TE.10", "PR.6"}),
		canonicalSection("s09", "09", "architecture", "Architecture & Security", "09-architecture", "systems", []string{"PD.1", "ARCH.1", "SEC.1"}, []string{"PD.3", "ARCH.9", "SEC.11"}),
		canonicalSection("s10", "10", "production", "Production Operations", "10-production", "systems", []string{"SL.1", "GS.1", "CFG.1", "OPS.1", "DOCKER.1"}, []string{"SL.5", "GS.3", "CFG.5", "OPS.5", "DOCKER.3", "DEPLOY.3", "CG.3"}),
		canonicalSection("s11", "11", "flagship", "Flagship", "11-flagship", "systems", []string{"OPSL.1"}, []string{"OPSL.10"}),
	}
	canonicalSectionReadmeTracks = map[string][]string{
		"s05": {"MP.1-MP.4", "CL.1-CL.4", "EN.1-EN.6", "FS.1-FS.8"},
		"s08": {"TE.1-TE.10", "PR.1-PR.6"},
		"s09": {"PD.1-PD.3", "ARCH.1-ARCH.9", "SEC.1-SEC.11"},
		"s10": {"SL.1-SL.5", "GS.1-GS.3", "CFG.1-CFG.5", "OPS.1-OPS.5", "DOCKER.1-DOCKER.3", "DEPLOY.1-DEPLOY.3", "CG.1-CG.3"},
	}
	forbiddenSectionReadmeLabels = map[string][]string{
		"s05": {"PKG.1", "IO.1"},
		"s08": {"Track A", "Track B"},
		"s09": {"Track GR"},
	}
)

func canonicalSection(id, number, slug, title, pathPrefix, phase string, entryPoints, outputs []string) canonicalSectionContract {
	return canonicalSectionContract{
		ID:          id,
		Number:      number,
		Slug:        slug,
		Title:       title,
		PathPrefix:  pathPrefix,
		Phase:       phase,
		Status:      "stable",
		EntryPoints: entryPoints,
		Outputs:     outputs,
	}
}

func Validate(root string, report func(string)) (Result, error) {
	if report == nil {
		report = func(string) {}
	}

	filesScanned, runErrors, err := validateRunPaths(root, report)
	if err != nil {
		return Result{}, err
	}

	v2SectionCount, v2ItemCount, v2PlaceholderCount, v2Errors, hasV2, err := validateV2Curriculum(root, report)
	if err != nil {
		return Result{}, err
	}

	markdownErrors := validateMarkdownSurfaces(root, report)

	return Result{
		FilesScanned:     filesScanned,
		V2SectionCount:   v2SectionCount,
		V2ItemCount:      v2ItemCount,
		PlaceholderCount: v2PlaceholderCount,
		HasV2:            hasV2,
		ErrorCount:       runErrors + v2Errors + markdownErrors,
	}, nil
}

func isPlaceholderPath(path string) bool {
	placeholders := []string{
		"./SECTION/LESSON",
		"./path/to/",
		"./my-folder",
		"./my-messy-folder",
		"./00-how-computers-work/...",
	}

	for _, placeholder := range placeholders {
		if strings.Contains(path, placeholder) {
			return true
		}
	}

	return false
}

func isPlaceholderItem(item V2Item) bool {
	return strings.EqualFold(strings.TrimSpace(item.Status), "placeholder")
}

func isImplementedItem(item V2Item) bool {
	status := strings.TrimSpace(item.Status)
	return status == "" || strings.EqualFold(status, "implemented")
}

func cleanRepoPath(path string) (string, bool) {
	path = strings.TrimSpace(path)
	if path == "" {
		return "", false
	}

	nativePath := filepath.FromSlash(path)
	if filepath.IsAbs(nativePath) || filepath.VolumeName(nativePath) != "" {
		return "", false
	}

	clean := filepath.ToSlash(filepath.Clean(nativePath))
	if clean == "." || clean == ".." || strings.HasPrefix(clean, "../") {
		return "", false
	}

	return clean, true
}

func repoPath(root, path string) (string, string, bool) {
	clean, ok := cleanRepoPath(path)
	if !ok {
		return "", "", false
	}

	return filepath.Join(root, filepath.FromSlash(clean)), clean, true
}

func readRepoFile(root, path string) ([]byte, string, error) {
	abs, clean, ok := repoPath(root, path)
	if !ok {
		return nil, "", fmt.Errorf("path escapes repository root: %s", path)
	}

	data, err := os.ReadFile(abs)
	return data, clean, err
}

func pathExists(root, path string) bool {
	abs, _, ok := repoPath(root, path)
	if !ok {
		return false
	}

	_, err := os.Stat(abs)
	return err == nil
}

func extractCommandTargets(command string) ([]string, error) {
	command = strings.TrimSpace(command)
	if command == "" {
		return nil, errors.New("command is empty")
	}

	matches := runPathPattern.FindAllString(command, -1)
	if len(matches) == 0 {
		return nil, fmt.Errorf("command does not contain a concrete ./path target: %q", command)
	}

	targets := make([]string, 0, len(matches))
	for _, match := range matches {
		if match == "./..." || isPlaceholderPath(match) {
			continue
		}
		target := strings.TrimSuffix(match, "/...")
		targets = append(targets, filepath.Clean(strings.TrimPrefix(target, "./")))
	}

	if len(targets) == 0 {
		return nil, fmt.Errorf("command does not contain a concrete ./path target: %q", command)
	}

	return targets, nil
}

func shouldScanRunPaths(path string) bool {
	if filepath.Ext(path) == ".go" {
		if strings.HasSuffix(filepath.Base(path), "_test.go") {
			return false
		}
		return true
	}

	if filepath.Base(path) != "README.md" {
		return false
	}

	first := strings.Split(filepath.ToSlash(path), "/")[0]
	if len(first) < 2 {
		return false
	}

	return first[0] >= '0' && first[0] <= '9' && first[1] >= '0' && first[1] <= '9'
}

func validateRunPaths(root string, report func(string)) (int, int, error) {
	filesScanned := 0
	errorsFound := 0

	walkErr := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			switch d.Name() {
			case ".git", "vendor":
				return filepath.SkipDir
			default:
				return nil
			}
		}

		relPath, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}

		cleanPath := filepath.Clean(relPath)
		if !shouldScanRunPaths(cleanPath) {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}

		filesScanned++
		scanner := bufio.NewScanner(file)
		lineNo := 0
		for scanner.Scan() {
			lineNo++
			line := scanner.Text()
			if !strings.Contains(line, "go run ") && !strings.Contains(line, "go test ") {
				continue
			}

			for _, match := range runPathPattern.FindAllString(line, -1) {
				if match == "./..." || isPlaceholderPath(match) {
					continue
				}

				trimmedMatch := strings.TrimSuffix(match, "/...")
				target := filepath.Clean(strings.TrimPrefix(trimmedMatch, "./"))
				alternateTarget := filepath.Clean(filepath.Join(filepath.Dir(cleanPath), strings.TrimPrefix(trimmedMatch, "./")))

				if pathExists(root, target) || pathExists(root, alternateTarget) {
					continue
				}

				report(fmt.Sprintf("Invalid run path: %s:%d -> %s", cleanPath, lineNo, match))
				errorsFound++
			}
		}

		scanErr := scanner.Err()
		closeErr := file.Close()
		if scanErr != nil {
			return scanErr
		}
		return closeErr
	})
	if walkErr != nil {
		return 0, 0, fmt.Errorf("Failed to scan run paths: %v", walkErr)
	}

	return filesScanned, errorsFound, nil
}

func validateV2Curriculum(root string, report func(string)) (int, int, int, int, bool, error) {
	if !pathExists(root, "curriculum.v2.json") {
		return 0, 0, 0, 0, false, nil
	}

	data, _, err := readRepoFile(root, "curriculum.v2.json")
	if err != nil {
		return 0, 0, 0, 0, false, fmt.Errorf("Failed to read curriculum.v2.json: %v", err)
	}

	var cur V2Curriculum
	if err := json.Unmarshal(data, &cur); err != nil {
		return 0, 0, 0, 0, false, fmt.Errorf("Failed to parse curriculum.v2.json: %v", err)
	}

	errorsFound := 0
	placeholderCount := 0
	sectionIDs := make(map[string]V2Section, len(cur.Sections))
	itemIDs := make(map[string]V2Item, len(cur.Items))

	if cur.SchemaVersion != expectedSchemaVersion {
		report(fmt.Sprintf("Invalid v2 schema_version: %d (expected %d)", cur.SchemaVersion, expectedSchemaVersion))
		errorsFound++
	}

	for _, s := range cur.Sections {
		if s.ID == "" {
			report("Invalid v2 section: missing id")
			errorsFound++
			continue
		}

		if _, exists := sectionIDs[s.ID]; exists {
			report(fmt.Sprintf("Duplicate v2 section id: %s", s.ID))
			errorsFound++
			continue
		}

		if s.Number == "" || s.Slug == "" || s.Title == "" || s.PathPrefix == "" || s.Phase == "" || s.Summary == "" {
			report(fmt.Sprintf("Invalid v2 section metadata: %s requires number, slug, title, path_prefix, phase, and summary", s.ID))
			errorsFound++
		}

		if s.PathPrefix != "" && !pathExists(root, s.PathPrefix) {
			report(fmt.Sprintf("Invalid v2 section path_prefix: %s -> %s", s.ID, s.PathPrefix))
			errorsFound++
		}

		if strings.TrimSpace(s.Phase) != "" && !allowedSectionPhases[s.Phase] {
			report(fmt.Sprintf("Invalid v2 section phase: %s -> %s", s.ID, s.Phase))
			errorsFound++
		}

		if strings.TrimSpace(s.Status) != "" && !allowedSectionStatuses[s.Status] {
			report(fmt.Sprintf("Invalid v2 section status: %s -> %s", s.ID, s.Status))
			errorsFound++
		}

		sectionIDs[s.ID] = s
	}

	for _, item := range cur.Items {
		if item.ID == "" {
			report("Invalid v2 item: missing id")
			errorsFound++
			continue
		}

		if _, exists := itemIDs[item.ID]; exists {
			report(fmt.Sprintf("Duplicate v2 item id: %s", item.ID))
			errorsFound++
			continue
		}

		if item.SectionID == "" || item.Slug == "" || item.Title == "" || item.Type == "" || item.Level == "" || item.VerificationMode == "" || item.Path == "" {
			report(fmt.Sprintf("Invalid v2 item metadata: %s requires section_id, slug, title, type, level, verification_mode, and path", item.ID))
			errorsFound++
		}

		if _, exists := sectionIDs[item.SectionID]; !exists {
			report(fmt.Sprintf("Invalid v2 section linkage: %s -> %s", item.ID, item.SectionID))
			errorsFound++
		}

		if !allowedItemTypes[item.Type] {
			report(fmt.Sprintf("Invalid v2 item type: %s -> %s", item.ID, item.Type))
			errorsFound++
		}

		if item.Type == "lesson" {
			if !allowedLessonSubtypes[item.Subtype] {
				report(fmt.Sprintf("Invalid v2 lesson subtype: %s -> %s", item.ID, item.Subtype))
				errorsFound++
			}
		} else if item.Subtype != "" {
			report(fmt.Sprintf("Unexpected v2 subtype for non-lesson item: %s -> %s", item.ID, item.Subtype))
			errorsFound++
		}

		if !allowedLevels[item.Level] {
			report(fmt.Sprintf("Invalid v2 item level: %s -> %s", item.ID, item.Level))
			errorsFound++
		}

		if !allowedVerificationModes[item.VerificationMode] {
			report(fmt.Sprintf("Invalid v2 verification mode: %s -> %s", item.ID, item.VerificationMode))
			errorsFound++
		}

		if !isImplementedItem(item) && !isPlaceholderItem(item) {
			report(fmt.Sprintf("Invalid v2 item status: %s -> %s", item.ID, item.Status))
			errorsFound++
		}

		// Placeholder items: skip deep validation (path, run commands, etc.)
		if isPlaceholderItem(item) {
			if !pathExists(root, item.Path) {
				report(fmt.Sprintf("Invalid v2 placeholder path: %s -> %s", item.ID, item.Path))
				errorsFound++
				itemIDs[item.ID] = item
				continue
			}

			placeholderCount++
			report(fmt.Sprintf("Warning: placeholder item: %s (%s) - not yet implemented", item.ID, item.Title))
			itemIDs[item.ID] = item
			continue
		}

		if !pathExists(root, item.Path) {
			report(fmt.Sprintf("Invalid v2 item path: %s -> %s", item.ID, item.Path))
			errorsFound++
		}

		if section, exists := sectionIDs[item.SectionID]; exists && section.PathPrefix != "" {
			itemPath := filepath.ToSlash(filepath.Clean(item.Path))
			allowedPrefixes := allowedPathPrefixesForSection(section)
			if !matchesAnyPrefix(itemPath, allowedPrefixes) {
				report(fmt.Sprintf("Invalid v2 section path alignment: %s -> %s (expected prefix %s)", item.ID, item.Path, strings.Join(allowedPrefixes, " or ")))
				errorsFound++
			}
		}

		if item.RunCommand != "" {
			targets, err := extractCommandTargets(item.RunCommand)
			if err != nil {
				report(fmt.Sprintf("Invalid v2 run command: %s -> %v", item.ID, err))
				errorsFound++
			} else {
				for _, target := range targets {
					if !pathExists(root, target) {
						report(fmt.Sprintf("Invalid v2 run command target: %s -> %s", item.ID, item.RunCommand))
						errorsFound++
						break
					}
				}
			}
		}

		if item.TestCommand != "" {
			targets, err := extractCommandTargets(item.TestCommand)
			if err != nil {
				report(fmt.Sprintf("Invalid v2 test command: %s -> %v", item.ID, err))
				errorsFound++
			} else {
				for _, target := range targets {
					if !pathExists(root, target) {
						report(fmt.Sprintf("Invalid v2 test command target: %s -> %s", item.ID, item.TestCommand))
						errorsFound++
						break
					}
				}
			}
		}

		switch item.VerificationMode {
		case "run":
			if strings.TrimSpace(item.RunCommand) == "" {
				report(fmt.Sprintf("Invalid v2 run contract: %s requires run_command", item.ID))
				errorsFound++
			}
		case "test":
			if strings.TrimSpace(item.TestCommand) == "" {
				report(fmt.Sprintf("Invalid v2 test contract: %s requires test_command", item.ID))
				errorsFound++
			}
		case "mixed":
			if strings.TrimSpace(item.RunCommand) == "" && strings.TrimSpace(item.TestCommand) == "" {
				report(fmt.Sprintf("Invalid v2 mixed contract: %s requires run_command or test_command", item.ID))
				errorsFound++
			}
		case "rubric":
			readmePath := filepath.ToSlash(filepath.Join(item.Path, "README.md"))
			if !pathExists(root, readmePath) {
				report(fmt.Sprintf("Invalid v2 rubric contract: %s requires README proof surface", item.ID))
				errorsFound++
			}
		}

		if item.StarterPath != "" && !pathExists(root, item.StarterPath) {
			report(fmt.Sprintf("Invalid v2 starter path: %s -> %s", item.ID, item.StarterPath))
			errorsFound++
		}

		itemIDs[item.ID] = item
	}

	errorsFound += validateV2ArchitectureContract(cur.Sections, report)

	for _, s := range cur.Sections {
		for _, prereqID := range s.Prerequisites {
			if _, exists := sectionIDs[prereqID]; !exists {
				report(fmt.Sprintf("Invalid v2 section prerequisite: %s -> %s", s.ID, prereqID))
				errorsFound++
			}
		}

		for _, entryID := range s.EntryPoints {
			if _, exists := itemIDs[entryID]; !exists {
				report(fmt.Sprintf("Invalid v2 section entry point: %s -> %s", s.ID, entryID))
				errorsFound++
			}
		}

		for _, outputID := range s.Outputs {
			if _, exists := itemIDs[outputID]; !exists {
				report(fmt.Sprintf("Invalid v2 section output: %s -> %s", s.ID, outputID))
				errorsFound++
			}
		}
	}

	for _, item := range cur.Items {
		for _, prereqID := range item.Prerequisites {
			if _, itemExists := itemIDs[prereqID]; itemExists {
				continue
			}
			if _, sectionExists := sectionIDs[prereqID]; sectionExists {
				continue
			}
			report(fmt.Sprintf("Invalid v2 prerequisite: %s -> %s", item.ID, prereqID))
			errorsFound++
		}

		for _, nextID := range item.NextItems {
			if _, itemExists := itemIDs[nextID]; itemExists {
				continue
			}
			if _, sectionExists := sectionIDs[nextID]; sectionExists {
				continue
			}
			report(fmt.Sprintf("Invalid v2 next item: %s -> %s", item.ID, nextID))
			errorsFound++
		}
	}

	errorsFound += validateV2LessonNavigation(root, cur.Items, report)
	errorsFound += validateV2LessonSourceHeaders(root, cur.Items, report)
	errorsFound += validateV2ReadmeNavigation(root, cur.Items, report)
	errorsFound += validateFlagshipProjects(root, sectionIDs, cur.Items, report)
	errorsFound += validateV2SectionLabels(root, sectionIDs, cur.Items, report)
	errorsFound += validateSectionReadmeTrackLabels(root, sectionIDs, report)
	errorsFound += validateV2TextEncoding(root, sectionIDs, cur.Items, report)
	errorsFound += validateFoundationsReadmeContracts(root, cur.Items, report)
	errorsFound += validateEngineeringReadmeContracts(root, cur.Items, report)

	return len(cur.Sections), len(cur.Items), placeholderCount, errorsFound, true, nil
}

func sameStringSlice(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

func validateV2ArchitectureContract(sections []V2Section, report func(string)) int {
	errorsFound := 0

	if len(sections) != len(canonicalV2Sections) {
		report(fmt.Sprintf("Invalid v2 architecture contract: expected %d sections, found %d", len(canonicalV2Sections), len(sections)))
		errorsFound++
	}

	seen := make(map[string]bool, len(sections))
	for idx, section := range sections {
		if idx >= len(canonicalV2Sections) {
			report(fmt.Sprintf("Invalid v2 architecture contract: unexpected section %s at position %d", section.ID, idx))
			errorsFound++
			continue
		}

		expected := canonicalV2Sections[idx]
		seen[section.ID] = true

		if section.ID != expected.ID {
			report(fmt.Sprintf("Invalid v2 architecture contract: section position %d -> %s (expected %s)", idx, section.ID, expected.ID))
			errorsFound++
			continue
		}

		if section.Number != expected.Number {
			report(fmt.Sprintf("Invalid v2 section number: %s -> %s (expected %s)", section.ID, section.Number, expected.Number))
			errorsFound++
		}
		if section.Slug != expected.Slug {
			report(fmt.Sprintf("Invalid v2 section slug: %s -> %s (expected %s)", section.ID, section.Slug, expected.Slug))
			errorsFound++
		}
		if section.Title != expected.Title {
			report(fmt.Sprintf("Invalid v2 section title: %s -> %s (expected %s)", section.ID, section.Title, expected.Title))
			errorsFound++
		}
		if section.PathPrefix != expected.PathPrefix {
			report(fmt.Sprintf("Invalid v2 section path_prefix: %s -> %s (expected %s)", section.ID, section.PathPrefix, expected.PathPrefix))
			errorsFound++
		}
		if section.Phase != expected.Phase {
			report(fmt.Sprintf("Invalid v2 section phase: %s -> %s (expected %s)", section.ID, section.Phase, expected.Phase))
			errorsFound++
		}
		if section.Status != expected.Status {
			report(fmt.Sprintf("Invalid v2 section status: %s -> %s (expected %s)", section.ID, section.Status, expected.Status))
			errorsFound++
		}
		if !sameStringSlice(section.EntryPoints, expected.EntryPoints) {
			report(fmt.Sprintf("Invalid v2 section entry points: %s -> %s (expected %s)", section.ID, strings.Join(section.EntryPoints, ", "), strings.Join(expected.EntryPoints, ", ")))
			errorsFound++
		}
		if !sameStringSlice(section.Outputs, expected.Outputs) {
			report(fmt.Sprintf("Invalid v2 section outputs: %s -> %s (expected %s)", section.ID, strings.Join(section.Outputs, ", "), strings.Join(expected.Outputs, ", ")))
			errorsFound++
		}
	}

	for _, expected := range canonicalV2Sections {
		if !seen[expected.ID] {
			report(fmt.Sprintf("Invalid v2 architecture contract: missing section %s", expected.ID))
			errorsFound++
		}
	}

	return errorsFound
}

func validateFlagshipProjects(root string, sections map[string]V2Section, items []V2Item, report func(string)) int {
	stage, exists := sections["s11"]
	if !exists {
		return 0
	}

	errorsFound := 0
	reservedPrefixes := make(map[string]bool)
	type groupedItem struct {
		item   V2Item
		number int
	}

	projectItems := make(map[string][]groupedItem)
	projectRoots := make(map[string]string)
	rootOwners := make(map[string]string)

	for _, item := range items {
		prefix, _, ok := splitCurriculumID(item.ID)
		if !ok {
			continue
		}

		if item.SectionID != "s11" {
			reservedPrefixes[prefix] = true
		}
	}

	for _, item := range items {
		prefix, number, ok := splitCurriculumID(item.ID)
		if !ok || item.SectionID != "s11" {
			continue
		}

		if reservedPrefixes[prefix] {
			report(fmt.Sprintf("Invalid flagship project prefix: %s -> %s is already used outside s11", item.ID, prefix))
			errorsFound++
		}

		projectItems[prefix] = append(projectItems[prefix], groupedItem{item: item, number: number})

		projectRoot, ok := flagshipProjectRoot(item.Path)
		if !ok {
			report(fmt.Sprintf("Invalid flagship project path: %s -> %s", item.ID, item.Path))
			errorsFound++
			continue
		}

		if existingRoot, exists := projectRoots[prefix]; exists && existingRoot != projectRoot {
			report(fmt.Sprintf("Invalid flagship project root alignment: %s -> %s and %s", prefix, existingRoot, projectRoot))
			errorsFound++
		} else {
			projectRoots[prefix] = projectRoot
		}

		if existingPrefix, exists := rootOwners[projectRoot]; exists && existingPrefix != prefix {
			report(fmt.Sprintf("Invalid flagship project root reuse: %s and %s both map to %s", existingPrefix, prefix, projectRoot))
			errorsFound++
		} else {
			rootOwners[projectRoot] = prefix
		}
	}

	if len(stage.EntryPoints) != 1 {
		report(fmt.Sprintf("Invalid flagship stage contract: s11 requires exactly 1 entry point, found %d", len(stage.EntryPoints)))
		errorsFound++
	}
	if len(stage.Outputs) != 1 {
		report(fmt.Sprintf("Invalid flagship stage contract: s11 requires exactly 1 output, found %d", len(stage.Outputs)))
		errorsFound++
	}

	canonicalPrefix := ""
	if len(stage.EntryPoints) == 1 {
		prefix, number, ok := splitCurriculumID(stage.EntryPoints[0])
		if !ok || number != 1 {
			report(fmt.Sprintf("Invalid flagship stage entry point: s11 -> %s", stage.EntryPoints[0]))
			errorsFound++
		} else {
			canonicalPrefix = prefix
		}
	}
	if len(stage.Outputs) == 1 {
		outputPrefix, _, ok := splitCurriculumID(stage.Outputs[0])
		if !ok {
			report(fmt.Sprintf("Invalid flagship stage output: s11 -> %s", stage.Outputs[0]))
			errorsFound++
		} else if canonicalPrefix != "" && outputPrefix != canonicalPrefix {
			report(fmt.Sprintf("Invalid flagship stage contract: s11 entry prefix %s does not match output prefix %s", canonicalPrefix, outputPrefix))
			errorsFound++
		}
	}

	for prefix, grouped := range projectItems {
		if !flagshipPrefixPattern.MatchString(prefix) {
			report(fmt.Sprintf("Invalid flagship project prefix: %s must be 3-6 uppercase letters", prefix))
			errorsFound++
		}

		sort.Slice(grouped, func(i, j int) bool {
			return grouped[i].number < grouped[j].number
		})

		for idx, entry := range grouped {
			expectedNumber := idx + 1
			if entry.number != expectedNumber {
				report(fmt.Sprintf("Invalid flagship module numbering: %s expected %s.%d", entry.item.ID, prefix, expectedNumber))
				errorsFound++
			}
		}

		for idx, entry := range grouped {
			if idx == len(grouped)-1 {
				if len(entry.item.NextItems) != 0 {
					report(fmt.Sprintf("Invalid flagship module chain: %s must terminate the project chain", entry.item.ID))
					errorsFound++
				}
				continue
			}

			expectedNext := grouped[idx+1].item.ID
			if len(entry.item.NextItems) != 1 || entry.item.NextItems[0] != expectedNext {
				report(fmt.Sprintf("Invalid flagship module chain: %s must point to %s", entry.item.ID, expectedNext))
				errorsFound++
			}
		}

		projectRoot := projectRoots[prefix]
		if projectRoot == "" {
			continue
		}

		if !pathExists(root, filepath.ToSlash(filepath.Join(projectRoot, "README.md"))) {
			report(fmt.Sprintf("Missing flagship project README: %s -> %s/README.md", prefix, filepath.ToSlash(projectRoot)))
			errorsFound++
		}
		if !pathExists(root, filepath.ToSlash(filepath.Join(projectRoot, "MODULES.md"))) {
			report(fmt.Sprintf("Missing flagship project module map: %s -> %s/MODULES.md", prefix, filepath.ToSlash(projectRoot)))
			errorsFound++
		}

		implementedProject := false
		for _, entry := range grouped {
			if isImplementedItem(entry.item) {
				implementedProject = true
				break
			}
		}
		if implementedProject && !pathExists(root, filepath.ToSlash(filepath.Join(projectRoot, "scripts", "progress.go"))) {
			report(fmt.Sprintf("Missing flagship progress checker: %s -> %s/scripts/progress.go", prefix, filepath.ToSlash(projectRoot)))
			errorsFound++
		}

		if canonicalPrefix == prefix && len(stage.Outputs) == 1 {
			expectedFinal := grouped[len(grouped)-1].item.ID
			if stage.Outputs[0] != expectedFinal {
				report(fmt.Sprintf("Invalid flagship stage output: s11 -> %s (expected %s)", stage.Outputs[0], expectedFinal))
				errorsFound++
			}
		}
	}

	return errorsFound
}

func splitCurriculumID(id string) (string, int, bool) {
	prefix, suffix, found := strings.Cut(id, ".")
	if !found || prefix == "" || suffix == "" {
		return "", 0, false
	}

	number, err := strconv.Atoi(suffix)
	if err != nil {
		return "", 0, false
	}

	return prefix, number, true
}

func flagshipProjectRoot(itemPath string) (string, bool) {
	parts := strings.Split(filepath.ToSlash(filepath.Clean(itemPath)), "/")
	if len(parts) < 2 || parts[0] != "11-flagship" {
		return "", false
	}

	return filepath.ToSlash(filepath.Join(parts[0], parts[1])), true
}

func allowedPathPrefixesForSection(section V2Section) []string {
	return []string{filepath.ToSlash(filepath.Clean(section.PathPrefix))}
}

func matchesAnyPrefix(itemPath string, prefixes []string) bool {
	for _, prefix := range prefixes {
		if itemPath == prefix || strings.HasPrefix(itemPath, prefix+"/") {
			return true
		}
	}

	return false
}

func isFoundationsSection(sectionID string) bool {
	switch sectionID {
	case "s00", "s01", "s02", "s03", "s04":
		return true
	default:
		return false
	}
}

func isEngineeringSection(sectionID string) bool {
	switch sectionID {
	case "s05", "s06", "s07", "s08", "s09", "s10":
		return true
	default:
		return false
	}
}

func validateFoundationsReadmeContracts(root string, items []V2Item, report func(string)) int {
	errorsFound := 0

	for _, item := range items {
		if !isFoundationsSection(item.SectionID) {
			continue
		}
		if isPlaceholderItem(item) {
			continue
		}

		itemPath := filepath.ToSlash(filepath.Clean(item.Path))
		readmePath := filepath.ToSlash(filepath.Join(itemPath, "README.md"))
		if !pathExists(root, readmePath) {
			report(fmt.Sprintf("Missing foundations README: %s -> %s", item.ID, readmePath))
			errorsFound++
			continue
		}

		errorsFound += validateRequiredHeadingsForItem(root, readmePath, item, report)
		errorsFound += validateFoundationsVisualModelUsesMermaid(root, readmePath, item.ID, report)

		if item.Type == "lesson" && item.VerificationMode == "run" {
			mainPath := filepath.ToSlash(filepath.Join(itemPath, "main.go"))
			if !pathExists(root, mainPath) {
				report(fmt.Sprintf("Missing foundations lesson main.go: %s -> %s", item.ID, mainPath))
				errorsFound++
			}
		}
	}

	return errorsFound
}

func validateEngineeringReadmeContracts(root string, items []V2Item, report func(string)) int {
	errorsFound := 0

	for _, item := range items {
		if !isEngineeringSection(item.SectionID) {
			continue
		}
		if isPlaceholderItem(item) {
			continue
		}

		itemPath := filepath.ToSlash(filepath.Clean(item.Path))
		readmePath := filepath.ToSlash(filepath.Join(itemPath, "README.md"))
		if !pathExists(root, readmePath) {
			report(fmt.Sprintf("Invalid engineering README contract: %s -> %s missing entirely", item.ID, filepath.ToSlash(readmePath)))
			errorsFound++
			continue
		}

		errorsFound += validateEngineeringHeadings(root, readmePath, item, report)
	}

	return errorsFound
}

func validateEngineeringHeadings(root, readmePath string, item V2Item, report func(string)) int {
	return validateOrderedMarkdownHeadings(root, readmePath, item.ID, item, "engineering", report)
}

func validateRequiredHeadingsForItem(root, readmePath string, item V2Item, report func(string)) int {
	return validateOrderedMarkdownHeadings(root, readmePath, item.ID, item, "foundations", report)
}

func requiredReadmeHeadingsForItem(item V2Item) []string {
	requiredHeadings := []string{
		"## Mission",
		"## Prerequisites",
		"## Mental Model",
		"## Visual Model",
		"## Machine View",
		"## Run Instructions",
	}

	if item.Type == "lesson" {
		requiredHeadings = append(requiredHeadings, "## Code Walkthrough")
	} else {
		requiredHeadings = append(requiredHeadings, "## Solution Walkthrough")
	}

	requiredHeadings = append(requiredHeadings, "## Try It")

	if item.Type != "lesson" {
		requiredHeadings = append(requiredHeadings, "## Verification Surface")
	}

	requiredHeadings = append(requiredHeadings,
		"## In Production",
		"## Thinking Questions",
		"## Next Step",
	)

	return requiredHeadings
}

func validateOrderedMarkdownHeadings(root, relPath, itemID string, item V2Item, contractName string, report func(string)) int {
	data, cleanPath, err := readRepoFile(root, relPath)
	if err != nil {
		report(fmt.Sprintf("Failed to read %s README: %s -> %v", contractName, filepath.ToSlash(relPath), err))
		return 1
	}

	headings := markdownHeadings(string(data))
	requiredHeadings := requiredReadmeHeadingsForItem(item)
	errorsFound := 0

	lastIndex := -1
	for _, required := range requiredHeadings {
		foundIndex := -1
		for idx := lastIndex + 1; idx < len(headings); idx++ {
			if headings[idx] == required {
				foundIndex = idx
				break
			}
		}
		if foundIndex >= 0 {
			lastIndex = foundIndex
			continue
		}

		if containsString(headings, required) {
			report(fmt.Sprintf("Invalid %s README contract: %s -> %s has %s out of order", contractName, itemID, cleanPath, required))
			errorsFound++
			continue
		}

		report(fmt.Sprintf("Invalid %s README contract: %s -> %s missing %s", contractName, itemID, cleanPath, required))
		errorsFound++
	}

	return errorsFound
}

func markdownHeadings(text string) []string {
	scanner := bufio.NewScanner(strings.NewReader(strings.ReplaceAll(text, "\r\n", "\n")))
	inFence := false
	headings := []string{}

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "```") {
			inFence = !inFence
			continue
		}
		if inFence || !strings.HasPrefix(line, "## ") {
			continue
		}

		headings = append(headings, line)
	}

	return headings
}

func markdownSectionText(text, heading string) (string, bool) {
	scanner := bufio.NewScanner(strings.NewReader(strings.ReplaceAll(text, "\r\n", "\n")))
	inFence := false
	inSection := false
	var section strings.Builder

	for scanner.Scan() {
		rawLine := scanner.Text()
		line := strings.TrimSpace(rawLine)
		if strings.HasPrefix(line, "```") {
			if inSection {
				section.WriteString(rawLine)
				section.WriteByte('\n')
			}
			inFence = !inFence
			continue
		}

		if !inFence && strings.HasPrefix(line, "## ") {
			if line == heading {
				inSection = true
				section.Reset()
				continue
			}
			if inSection {
				return section.String(), true
			}
		}

		if inSection {
			section.WriteString(rawLine)
			section.WriteByte('\n')
		}
	}

	return section.String(), inSection
}

func containsString(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}

	return false
}
func validateFoundationsVisualModelUsesMermaid(root, relPath, itemID string, report func(string)) int {
	data, cleanPath, err := readRepoFile(root, relPath)
	if err != nil {
		report(fmt.Sprintf("Failed to read foundations README: %s -> %v", filepath.ToSlash(relPath), err))
		return 1
	}

	text := strings.ReplaceAll(string(data), "\r\n", "\n")
	sectionText, ok := markdownSectionText(text, "## Visual Model")
	if !ok {
		return 0
	}

	if !strings.Contains(sectionText, "```mermaid") {
		report(fmt.Sprintf("Invalid foundations README contract: %s -> %s Visual Model must include a Mermaid diagram", itemID, cleanPath))
		return 1
	}

	return 0
}

var mojibakeMarkers = []string{
	"\uFFFD",                                     // Unicode replacement rune from undecodable bytes.
	"\u00c3\u0192\u00c2\u00a2",                   // Double-encoded UTF-8 often seen from smart punctuation.
	"\u00c3\u0192\u00c2\u00b0",                   // Double-encoded UTF-8 from non-ASCII text.
	"\u00c3\u0192\u00c6\u2019",                   // Double-encoded Latin-1/Windows-1252 marker.
	"\u00c3\u0192\u00e2\u20ac\u0161",             // Double-encoded Windows-1252 punctuation marker.
	"\u00c3\u00b0\u00c5\u00b8",                   // Broken emoji or symbol prefix from UTF-8/Latin-1 confusion.
	"\u00e2\u0153",                               // Broken checkmark/cross sequence.
	"\u00e2\u0161",                               // Broken warning-symbol sequence.
	"\u00ef\u00b8",                               // Broken emoji variation selector.
	"\u00c3\u00a2\u201a\u00ac\u00e2\u20ac\u009d", // Double-encoded dash/quote sequence.
	"\u00c3\u00a2\u20ac\u00a0",                   // Broken dagger/bullet punctuation.
	"\u00c3\u00a2\u00c5\u201c",                   // Broken opening-quote sequence.
	"\u00c3\u00a2\u00c2\u009d",                   // Broken cross/checkmark sequence.
}

func validateV2TextEncoding(root string, sections map[string]V2Section, items []V2Item, report func(string)) int {
	errorsFound := 0

	for _, item := range items {
		section, exists := sections[item.SectionID]
		if !exists {
			continue
		}
		if section.Number < "09" {
			continue
		}

		candidateFiles := []string{
			filepath.Join(item.Path, "main.go"),
			filepath.Join(item.Path, "README.md"),
		}
		if item.StarterPath != "" {
			candidateFiles = append(candidateFiles, filepath.Join(item.StarterPath, "main.go"))
		}

		for _, rel := range candidateFiles {
			if !pathExists(root, rel) {
				continue
			}

			data, cleanPath, err := readRepoFile(root, rel)
			if err != nil {
				report(fmt.Sprintf("Failed to read v2 text surface: %s -> %v", filepath.ToSlash(rel), err))
				errorsFound++
				continue
			}

			text := string(data)
			for _, marker := range mojibakeMarkers {
				if strings.Contains(text, marker) {
					report(fmt.Sprintf("Possible mojibake in v2 text surface: %s -> %s", item.ID, cleanPath))
					errorsFound++
					break
				}
			}
		}
	}

	return errorsFound
}

func validateV2SectionLabels(root string, sections map[string]V2Section, items []V2Item, report func(string)) int {
	errorsFound := 0

	for _, item := range items {
		section, exists := sections[item.SectionID]
		if !exists {
			continue
		}
		if section.Number < "09" {
			continue
		}

		expectedSectionLabel := fmt.Sprintf("Section %s", section.Number)
		expectedStageLabel := fmt.Sprintf("Stage %s", section.Number)
		candidateFiles := []string{
			filepath.Join(item.Path, "main.go"),
		}
		if item.StarterPath != "" {
			candidateFiles = append(candidateFiles, filepath.Join(item.StarterPath, "main.go"))
		}

		for _, rel := range candidateFiles {
			if !pathExists(root, rel) {
				continue
			}

			data, cleanPath, err := readRepoFile(root, rel)
			if err != nil {
				report(fmt.Sprintf("Failed to read v2 section label surface: %s -> %v", filepath.ToSlash(rel), err))
				errorsFound++
				continue
			}

			text := string(data)
			if !strings.Contains(text, expectedSectionLabel) && !strings.Contains(text, expectedStageLabel) {
				report(fmt.Sprintf("Invalid v2 section label: %s -> %s (expected %s or %s)", item.ID, cleanPath, expectedSectionLabel, expectedStageLabel))
				errorsFound++
			}
		}
	}

	return errorsFound
}

func validateSectionReadmeTrackLabels(root string, sections map[string]V2Section, report func(string)) int {
	errorsFound := 0

	for sectionID, expectedLabels := range canonicalSectionReadmeTracks {
		section, exists := sections[sectionID]
		if !exists || section.PathPrefix == "" {
			continue
		}

		readmePath := filepath.ToSlash(filepath.Join(section.PathPrefix, "README.md"))
		data, cleanPath, err := readRepoFile(root, readmePath)
		if err != nil {
			report(fmt.Sprintf("Invalid section README contract: %s -> %s read failure: %v", sectionID, readmePath, err))
			errorsFound++
			continue
		}

		text := string(data)
		for _, label := range expectedLabels {
			if !strings.Contains(text, label) {
				report(fmt.Sprintf("Invalid section README contract: %s -> %s missing canonical track label %s", sectionID, cleanPath, label))
				errorsFound++
			}
		}

		for _, label := range forbiddenSectionReadmeLabels[sectionID] {
			if strings.Contains(text, label) {
				report(fmt.Sprintf("Invalid section README contract: %s -> %s contains non-canonical track label %s", sectionID, cleanPath, label))
				errorsFound++
			}
		}
	}

	return errorsFound
}

func validateV2LessonNavigation(root string, items []V2Item, report func(string)) int {
	errorsFound := 0
	itemIDs := make(map[string]V2Item, len(items))
	for _, item := range items {
		itemIDs[item.ID] = item
	}

	for _, item := range items {
		if len(item.NextItems) == 0 {
			continue
		}
		if isPlaceholderItem(item) {
			continue
		}

		expectedNextID := item.NextItems[0]
		if strings.HasPrefix(expectedNextID, "s") {
			continue
		}
		expectedNextItem, ok := itemIDs[expectedNextID]
		if !ok {
			continue
		}

		mainPath := filepath.ToSlash(filepath.Join(item.Path, "main.go"))
		if !pathExists(root, mainPath) {
			continue
		}

		data, cleanPath, err := readRepoFile(root, mainPath)
		if err != nil {
			report(fmt.Sprintf("Failed to read v2 lesson source: %s -> %v", item.ID, err))
			errorsFound++
			continue
		}

		expectedNextPath := filepath.ToSlash(expectedNextItem.Path)
		expectedFooter := fmt.Sprintf("NEXT UP: %s -> %s", expectedNextID, expectedNextPath)
		match := nextUpFooterPattern.FindSubmatch(data)
		if len(match) < 3 {
			report(fmt.Sprintf("Missing v2 lesson navigation footer: %s -> %s (expected %q)", item.ID, cleanPath, expectedFooter))
			errorsFound++
			continue
		}

		actualNextID := string(match[1])
		if actualNextID != expectedNextID {
			report(fmt.Sprintf("Invalid v2 lesson navigation footer: %s -> %s (expected %s)", item.ID, actualNextID, expectedNextID))
			errorsFound++
			continue
		}
		actualNextPath := string(match[2])
		if actualNextPath != expectedNextPath {
			report(fmt.Sprintf("Invalid v2 lesson navigation footer path: %s -> %s (expected %s)", item.ID, actualNextPath, expectedNextPath))
			errorsFound++
		}
	}

	return errorsFound
}

func validateV2LessonSourceHeaders(root string, items []V2Item, report func(string)) int {
	errorsFound := 0

	for _, item := range items {
		if isPlaceholderItem(item) {
			continue
		}

		expectedLevel, ok := levelDisplayLabels[item.Level]
		if !ok {
			continue
		}

		mainPath := filepath.ToSlash(filepath.Join(item.Path, "main.go"))
		if !pathExists(root, mainPath) {
			continue
		}

		data, cleanPath, err := readRepoFile(root, mainPath)
		if err != nil {
			report(fmt.Sprintf("Failed to read v2 lesson source header: %s -> %v", item.ID, err))
			errorsFound++
			continue
		}

		actualLevel, found := lessonCommentHeaderValue(data, "Level:")
		if !found {
			report(fmt.Sprintf("Missing v2 lesson level header: %s -> %s (expected Level: %s)", item.ID, cleanPath, expectedLevel))
			errorsFound++
			continue
		}

		if actualLevel != expectedLevel {
			report(fmt.Sprintf("Invalid v2 lesson level header: %s -> %s has Level: %s (expected Level: %s)", item.ID, cleanPath, actualLevel, expectedLevel))
			errorsFound++
		}

		expectedRun := strings.TrimSpace(item.RunCommand)
		if expectedRun == "" {
			continue
		}

		actualRun, found := lessonRunHeader(data)
		if !found {
			report(fmt.Sprintf("Missing v2 lesson run header: %s -> %s (expected RUN: %s)", item.ID, cleanPath, expectedRun))
			errorsFound++
			continue
		}

		if actualRun != expectedRun {
			report(fmt.Sprintf("Invalid v2 lesson run header: %s -> %s has RUN: %s (expected RUN: %s)", item.ID, cleanPath, actualRun, expectedRun))
			errorsFound++
		}
	}

	return errorsFound
}

func lessonCommentHeaderValue(data []byte, prefix string) (string, bool) {
	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		value, ok := strings.CutPrefix(line, "// "+prefix)
		if !ok {
			continue
		}
		return strings.TrimSpace(value), true
	}

	return "", false
}

func lessonRunHeader(data []byte) (string, bool) {
	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		run, ok := strings.CutPrefix(line, "// RUN:")
		if !ok {
			continue
		}

		run = strings.TrimSpace(run)
		if run != "" {
			return run, true
		}

		for scanner.Scan() {
			next := strings.TrimSpace(scanner.Text())
			if next == "//" || next == "" {
				continue
			}
			value, ok := strings.CutPrefix(next, "//")
			if !ok {
				return "", true
			}
			return strings.TrimSpace(value), true
		}

		return "", true
	}

	return "", false
}

func validateV2ReadmeNavigation(root string, items []V2Item, report func(string)) int {
	errorsFound := 0
	itemIDs := make(map[string]V2Item, len(items))
	pathCounts := make(map[string]int, len(items))
	for _, item := range items {
		itemIDs[item.ID] = item
		if !isPlaceholderItem(item) {
			if cleanPath, ok := cleanRepoPath(item.Path); ok {
				pathCounts[cleanPath]++
			}
		}
	}

	for _, item := range items {
		if len(item.NextItems) == 0 {
			continue
		}
		if isPlaceholderItem(item) {
			continue
		}
		cleanItemPath, ok := cleanRepoPath(item.Path)
		if !ok {
			continue
		}
		if pathCounts[cleanItemPath] > 1 {
			continue
		}

		expectedNextID := item.NextItems[0]
		if strings.HasPrefix(expectedNextID, "s") {
			continue
		}
		expectedNextItem, ok := itemIDs[expectedNextID]
		if !ok {
			continue
		}

		readmePath := filepath.ToSlash(filepath.Join(item.Path, "README.md"))
		if !pathExists(root, readmePath) {
			continue
		}

		data, cleanPath, err := readRepoFile(root, readmePath)
		if err != nil {
			report(fmt.Sprintf("Failed to read v2 README navigation: %s -> %v", item.ID, err))
			errorsFound++
			continue
		}

		expectedLine, err := expectedReadmeNavigationLine(item, expectedNextID, expectedNextItem)
		if err != nil {
			report(fmt.Sprintf("Failed to resolve v2 README navigation target: %s -> %v", item.ID, err))
			errorsFound++
			continue
		}

		text := string(data)
		if !strings.Contains(text, expectedLine) {
			report(fmt.Sprintf("Invalid v2 README navigation footer: %s -> %s (expected %q)", item.ID, cleanPath, expectedLine))
			errorsFound++
		}
	}

	return errorsFound
}

func expectedReadmeNavigationLine(item V2Item, expectedNextID string, expectedNextItem V2Item) (string, error) {
	itemPath, ok := cleanRepoPath(item.Path)
	if !ok {
		return "", fmt.Errorf("invalid item path %q", item.Path)
	}
	nextPath, ok := cleanRepoPath(expectedNextItem.Path)
	if !ok {
		return "", fmt.Errorf("invalid next item path %q", expectedNextItem.Path)
	}

	linkTarget, err := filepath.Rel(filepath.FromSlash(itemPath), filepath.Join(filepath.FromSlash(nextPath), "README.md"))
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Next: `%s` -> [`%s`](%s)", expectedNextID, nextPath, filepath.ToSlash(linkTarget)), nil
}

func validateMarkdownSurfaces(root string, report func(string)) int {
	errorsFound := 0

	walkErr := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			switch d.Name() {
			case ".git", "vendor", ".agents", ".cache", ".github", ".opencode", "temp":
				return filepath.SkipDir
			default:
				return nil
			}
		}

		if filepath.Ext(path) != ".md" {
			return nil
		}

		relPath, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}

		cleanPath := filepath.ToSlash(filepath.Clean(relPath))
		errorsFound += validateMarkdownLocalLinks(root, cleanPath, report)
		errorsFound += validateMarkdownTextHealth(root, cleanPath, report)
		errorsFound += validateMarkdownReferenceAlerts(root, cleanPath, report)

		return nil
	})
	if walkErr != nil {
		report(fmt.Sprintf("Failed to scan markdown surfaces: %v", walkErr))
		return errorsFound + 1
	}

	return errorsFound
}

func validateMarkdownTextHealth(root, relPath string, report func(string)) int {
	data, cleanPath, err := readRepoFile(root, relPath)
	if err != nil {
		report(fmt.Sprintf("Failed to read markdown surface: %s -> %v", filepath.ToSlash(relPath), err))
		return 1
	}

	text := string(data)
	errorsFound := 0

	for _, marker := range mojibakeMarkers {
		if strings.Contains(text, marker) {
			report(fmt.Sprintf("Possible mojibake in markdown surface: %s", cleanPath))
			errorsFound++
			break
		}
	}

	return errorsFound
}

func validateMarkdownReferenceAlerts(root, relPath string, report func(string)) int {
	data, cleanPath, err := readRepoFile(root, relPath)
	if err != nil {
		report(fmt.Sprintf("Failed to read markdown surface: %s -> %v", filepath.ToSlash(relPath), err))
		return 1
	}

	errorsFound := 0
	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	lineNo := 0
	for scanner.Scan() {
		lineNo++
		line := strings.TrimSpace(scanner.Text())
		switch {
		case strings.HasPrefix(line, "> **Forward Reference:**"):
			report(fmt.Sprintf("Invalid markdown cross-reference alert: %s:%d uses Forward Reference label; use [!TIP] or [!NOTE]", cleanPath, lineNo))
			errorsFound++
		case strings.HasPrefix(line, "> **Backward Reference:**"):
			report(fmt.Sprintf("Invalid markdown cross-reference alert: %s:%d uses Backward Reference label; use [!NOTE]", cleanPath, lineNo))
			errorsFound++
		}
	}
	if err := scanner.Err(); err != nil {
		report(fmt.Sprintf("Failed to scan markdown surface: %s -> %v", cleanPath, err))
		return errorsFound + 1
	}

	return errorsFound
}

func validateMarkdownLocalLinks(root, relPath string, report func(string)) int {
	data, cleanPath, err := readRepoFile(root, relPath)
	if err != nil {
		report(fmt.Sprintf("Failed to read markdown surface: %s -> %v", filepath.ToSlash(relPath), err))
		return 1
	}

	errorsFound := 0
	for _, match := range markdownLinkPattern.FindAllStringSubmatch(string(data), -1) {
		if len(match) < 2 {
			continue
		}

		target := strings.TrimSpace(match[1])
		if target == "" ||
			strings.HasPrefix(target, "http://") ||
			strings.HasPrefix(target, "https://") ||
			strings.HasPrefix(target, "#") ||
			strings.HasPrefix(target, "mailto:") {
			continue
		}

		target = strings.SplitN(target, "#", 2)[0]
		resolved := filepath.Clean(filepath.Join(filepath.Dir(cleanPath), filepath.FromSlash(target)))
		if !pathExists(root, resolved) {
			report(fmt.Sprintf("Broken local doc link: %s -> %s", cleanPath, target))
			errorsFound++
		}
	}

	return errorsFound
}
