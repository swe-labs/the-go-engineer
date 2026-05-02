package curriculumvalidator

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"unicode"
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
var curriculumPrefixPattern = regexp.MustCompile(`^[A-Z]{2,6}$`)
var sectionIDPattern = regexp.MustCompile(`^s\d{2}$`)
var curriculumReferencePattern = regexp.MustCompile(`\b[A-Z]{2,6}\.\d+\b`)
var markdownAlertPattern = regexp.MustCompile(`^>\s*\[!([A-Z]+)\]`)
var markdownReadmeLinkPattern = regexp.MustCompile(`\[[^\]]+\]\([^)]*README\.md(?:#[^)]*)?\)`)
var legacyReferenceHeadingPattern = regexp.MustCompile(`^#{1,6}\s+.*\b(?:Forward|Backward) Reference\b`)
var staleCoverageCommandPattern = regexp.MustCompile(`go test\s+-coverprofile\s+coverage\.out\s+\./\.\.\.`)
var kebabSlugPattern = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)
var windowsAbsolutePathPattern = regexp.MustCompile(`^[A-Za-z]:[\\/]`)

const expectedSchemaVersion = 1

var (
	curriculumTopLevelFields = []string{"schema_version", "sections", "items"}
	curriculumSectionFields  = []string{"id", "number", "slug", "title", "path_prefix", "phase", "summary", "status", "entry_points", "outputs", "prerequisites"}
	curriculumItemFields     = []string{"id", "section_id", "slug", "title", "type", "subtype", "level", "status", "verification_mode", "path", "prerequisites", "run_command", "test_command", "starter_path", "next_items"}
	curriculumSectionArrays  = []string{"entry_points", "outputs", "prerequisites"}
	curriculumItemArrays     = []string{"prerequisites", "next_items"}
)

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
	standardsErrors := validateRepositoryStandards(root, report)

	return Result{
		FilesScanned:     filesScanned,
		V2SectionCount:   v2SectionCount,
		V2ItemCount:      v2ItemCount,
		PlaceholderCount: v2PlaceholderCount,
		HasV2:            hasV2,
		ErrorCount:       runErrors + v2Errors + markdownErrors + standardsErrors,
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
	if strings.Contains(path, "\\") || strings.HasPrefix(path, "//") || windowsAbsolutePathPattern.MatchString(path) {
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
		targets = append(targets, filepath.ToSlash(filepath.Clean(strings.TrimPrefix(target, "./"))))
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

		cleanPath := filepath.ToSlash(filepath.Clean(relPath))
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
				target := filepath.ToSlash(filepath.Clean(strings.TrimPrefix(trimmedMatch, "./")))
				alternateTarget := filepath.ToSlash(filepath.Clean(filepath.Join(filepath.Dir(cleanPath), strings.TrimPrefix(trimmedMatch, "./"))))

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

	errorsFound := validateV2CurriculumJSONContract(data, cur, report)
	placeholderCount := 0
	sectionIDs := make(map[string]V2Section, len(cur.Sections))
	itemIDs := make(map[string]V2Item, len(cur.Items))
	prefixOwners := make(map[string]string)

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
		if !sectionIDPattern.MatchString(s.ID) {
			report(fmt.Sprintf("Invalid v2 section id format: %s", s.ID))
			errorsFound++
		}
		if s.Number != "" && sectionNumberFromID(s.ID) != "" && s.Number != sectionNumberFromID(s.ID) {
			report(fmt.Sprintf("Invalid v2 section number alignment: %s -> %s", s.ID, s.Number))
			errorsFound++
		}
		if s.Slug != "" && !kebabSlugPattern.MatchString(s.Slug) {
			report(fmt.Sprintf("Invalid v2 section slug format: %s -> %s", s.ID, s.Slug))
			errorsFound++
		}
		if s.PathPrefix != "" && s.Number != "" && s.Slug != "" && s.PathPrefix != s.Number+"-"+s.Slug {
			report(fmt.Sprintf("Invalid v2 section path_prefix alignment: %s -> %s (expected %s-%s)", s.ID, s.PathPrefix, s.Number, s.Slug))
			errorsFound++
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

		errorsFound += validateNoDuplicateStrings(fmt.Sprintf("v2 section %s entry_points", s.ID), s.EntryPoints, report)
		errorsFound += validateNoDuplicateStrings(fmt.Sprintf("v2 section %s outputs", s.ID), s.Outputs, report)
		errorsFound += validateNoDuplicateStrings(fmt.Sprintf("v2 section %s prerequisites", s.ID), s.Prerequisites, report)

		sectionIDs[s.ID] = s
	}

	for _, item := range cur.Items {
		if item.ID == "" {
			report("Invalid v2 item: missing id")
			errorsFound++
			continue
		}
		if !isValidCurriculumItemID(item.ID) {
			report(fmt.Sprintf("Invalid v2 item id format: %s", item.ID))
			errorsFound++
		}
		prefix, _, ok := splitCurriculumID(item.ID)
		if ok {
			if owner, exists := prefixOwners[prefix]; exists && owner != item.SectionID {
				report(fmt.Sprintf("Invalid v2 item prefix ownership: %s is used in %s and %s", prefix, owner, item.SectionID))
				errorsFound++
			} else if !exists {
				prefixOwners[prefix] = item.SectionID
			}
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
		if item.Slug != "" && !kebabSlugPattern.MatchString(item.Slug) {
			report(fmt.Sprintf("Invalid v2 item slug format: %s -> %s", item.ID, item.Slug))
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

		errorsFound += validateNoDuplicateStrings(fmt.Sprintf("v2 item %s prerequisites", item.ID), item.Prerequisites, report)
		errorsFound += validateNoDuplicateStrings(fmt.Sprintf("v2 item %s next_items", item.ID), item.NextItems, report)

		itemIDs[item.ID] = item
	}

	errorsFound += validateV2ArchitectureContract(cur.Sections, report)
	errorsFound += validateV2ItemOrdering(cur.Items, sectionIDs, report)

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
			if prereqID == item.ID {
				report(fmt.Sprintf("Invalid v2 prerequisite: %s cannot reference itself", item.ID))
				errorsFound++
				continue
			}
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
			if nextID == item.ID {
				report(fmt.Sprintf("Invalid v2 next item: %s cannot reference itself", item.ID))
				errorsFound++
				continue
			}
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
	errorsFound += validateV2LessonSourceStandards(root, cur.Items, report)
	errorsFound += validateV2ReadmeNavigation(root, cur.Items, report)
	errorsFound += validateFlagshipProjects(root, sectionIDs, cur.Items, report)
	errorsFound += validateV2SectionLabels(root, sectionIDs, cur.Items, report)
	errorsFound += validateSectionReadmeTrackLabels(root, sectionIDs, report)
	errorsFound += validateV2TextEncoding(root, sectionIDs, cur.Items, report)
	errorsFound += validateFoundationsReadmeContracts(root, cur.Items, report)
	errorsFound += validateEngineeringReadmeContracts(root, cur.Items, report)

	return len(cur.Sections), len(cur.Items), placeholderCount, errorsFound, true, nil
}

func validateV2CurriculumJSONContract(data []byte, cur V2Curriculum, report func(string)) int {
	errorsFound := 0
	errorsFound += validateDuplicateJSONKeys(data, report)
	errorsFound += validateV2CurriculumRawShape(data, report)

	canonical, err := canonicalV2CurriculumJSON(cur)
	if err != nil {
		report(fmt.Sprintf("Invalid v2 curriculum JSON contract: failed to render canonical curriculum.v2.json -> %v", err))
		return errorsFound + 1
	}
	if !bytes.Equal(data, canonical) {
		report("Invalid v2 curriculum JSON formatting: curriculum.v2.json must use canonical field order, two-space indentation, [] arrays, and one trailing newline")
		errorsFound++
	}

	return errorsFound
}

func validateDuplicateJSONKeys(data []byte, report func(string)) int {
	decoder := json.NewDecoder(bytes.NewReader(data))
	errorsFound, err := scanJSONValueForDuplicateKeys(decoder, "$", report)
	if err != nil {
		report(fmt.Sprintf("Invalid v2 curriculum JSON structure: %v", err))
		return errorsFound + 1
	}

	if decoder.More() {
		report("Invalid v2 curriculum JSON structure: trailing JSON data")
		return errorsFound + 1
	}

	return errorsFound
}

func scanJSONValueForDuplicateKeys(decoder *json.Decoder, path string, report func(string)) (int, error) {
	token, err := decoder.Token()
	if err != nil {
		return 0, err
	}

	delim, ok := token.(json.Delim)
	if !ok {
		return 0, nil
	}

	errorsFound := 0
	switch delim {
	case '{':
		seen := map[string]bool{}
		for decoder.More() {
			keyToken, err := decoder.Token()
			if err != nil {
				return errorsFound, err
			}
			key, ok := keyToken.(string)
			if !ok {
				return errorsFound, fmt.Errorf("%s object key is not a string", path)
			}

			childPath := path + "." + key
			if seen[key] {
				report(fmt.Sprintf("Duplicate JSON key in curriculum.v2.json: %s", childPath))
				errorsFound++
			}
			seen[key] = true

			childErrors, err := scanJSONValueForDuplicateKeys(decoder, childPath, report)
			errorsFound += childErrors
			if err != nil {
				return errorsFound, err
			}
		}

		if _, err := decoder.Token(); err != nil {
			return errorsFound, err
		}
	case '[':
		idx := 0
		for decoder.More() {
			childErrors, err := scanJSONValueForDuplicateKeys(decoder, fmt.Sprintf("%s[%d]", path, idx), report)
			errorsFound += childErrors
			if err != nil {
				return errorsFound, err
			}
			idx++
		}

		if _, err := decoder.Token(); err != nil {
			return errorsFound, err
		}
	default:
		return errorsFound, fmt.Errorf("%s has unexpected delimiter %q", path, delim)
	}

	return errorsFound, nil
}

func validateV2CurriculumRawShape(data []byte, report func(string)) int {
	var top map[string]json.RawMessage
	if err := json.Unmarshal(data, &top); err != nil {
		report(fmt.Sprintf("Invalid v2 curriculum JSON structure: %v", err))
		return 1
	}

	errorsFound := validateJSONObjectFields("$", top, curriculumTopLevelFields, report)
	errorsFound += validateJSONArrayRawField("$", top, "sections", report)
	errorsFound += validateJSONArrayRawField("$", top, "items", report)

	var sections []map[string]json.RawMessage
	if raw, ok := top["sections"]; ok {
		if err := json.Unmarshal(raw, &sections); err != nil {
			report(fmt.Sprintf("Invalid v2 curriculum JSON structure: $.sections must be an array of objects -> %v", err))
			errorsFound++
		}
	}
	for idx, section := range sections {
		path := fmt.Sprintf("$.sections[%d]", idx)
		errorsFound += validateJSONObjectFields(path, section, curriculumSectionFields, report)
		for _, field := range curriculumSectionArrays {
			errorsFound += validateJSONArrayRawField(path, section, field, report)
		}
	}

	var items []map[string]json.RawMessage
	if raw, ok := top["items"]; ok {
		if err := json.Unmarshal(raw, &items); err != nil {
			report(fmt.Sprintf("Invalid v2 curriculum JSON structure: $.items must be an array of objects -> %v", err))
			errorsFound++
		}
	}
	for idx, item := range items {
		path := fmt.Sprintf("$.items[%d]", idx)
		errorsFound += validateJSONObjectFields(path, item, curriculumItemFields, report)
		for _, field := range curriculumItemArrays {
			errorsFound += validateJSONArrayRawField(path, item, field, report)
		}
	}

	return errorsFound
}

func validateJSONObjectFields(path string, object map[string]json.RawMessage, expected []string, report func(string)) int {
	errorsFound := 0
	expectedSet := make(map[string]bool, len(expected))
	for _, field := range expected {
		expectedSet[field] = true
		if _, exists := object[field]; !exists {
			report(fmt.Sprintf("Invalid v2 curriculum JSON object: %s missing field %s", path, field))
			errorsFound++
		}
	}

	for field := range object {
		if expectedSet[field] {
			continue
		}
		report(fmt.Sprintf("Invalid v2 curriculum JSON object: %s has unknown field %s", path, field))
		errorsFound++
	}

	return errorsFound
}

func validateJSONArrayRawField(path string, object map[string]json.RawMessage, field string, report func(string)) int {
	raw, exists := object[field]
	if !exists {
		return 0
	}

	trimmed := bytes.TrimSpace(raw)
	if bytes.Equal(trimmed, []byte("null")) {
		report(fmt.Sprintf("Invalid v2 curriculum JSON array: %s.%s must be [] instead of null", path, field))
		return 1
	}
	if len(trimmed) == 0 || trimmed[0] != '[' {
		report(fmt.Sprintf("Invalid v2 curriculum JSON array: %s.%s must be an array", path, field))
		return 1
	}

	return 0
}

func canonicalV2CurriculumJSON(cur V2Curriculum) ([]byte, error) {
	normalized := normalizeV2Curriculum(cur)
	var output bytes.Buffer
	encoder := json.NewEncoder(&output)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(normalized); err != nil {
		return nil, err
	}
	return output.Bytes(), nil
}

func normalizeV2Curriculum(cur V2Curriculum) V2Curriculum {
	normalized := cur
	normalized.Sections = make([]V2Section, len(cur.Sections))
	for idx, section := range cur.Sections {
		normalized.Sections[idx] = section
		normalized.Sections[idx].EntryPoints = normalizeStringSlice(section.EntryPoints)
		normalized.Sections[idx].Outputs = normalizeStringSlice(section.Outputs)
		normalized.Sections[idx].Prerequisites = normalizeStringSlice(section.Prerequisites)
	}

	normalized.Items = make([]V2Item, len(cur.Items))
	for idx, item := range cur.Items {
		normalized.Items[idx] = item
		normalized.Items[idx].Prerequisites = normalizeStringSlice(item.Prerequisites)
		normalized.Items[idx].NextItems = normalizeStringSlice(item.NextItems)
	}

	return normalized
}

func normalizeStringSlice(values []string) []string {
	if values == nil {
		return []string{}
	}
	return values
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

func isValidCurriculumItemID(id string) bool {
	prefix, _, ok := splitCurriculumID(id)
	return ok && curriculumPrefixPattern.MatchString(prefix)
}

func sectionNumberFromID(id string) string {
	if !sectionIDPattern.MatchString(id) {
		return ""
	}
	return strings.TrimPrefix(id, "s")
}

func validateNoDuplicateStrings(label string, values []string, report func(string)) int {
	errorsFound := 0
	seen := make(map[string]bool, len(values))
	for _, value := range values {
		if value == "" {
			report(fmt.Sprintf("Invalid %s: empty value", label))
			errorsFound++
			continue
		}
		if seen[value] {
			report(fmt.Sprintf("Invalid %s: duplicate value %s", label, value))
			errorsFound++
			continue
		}
		seen[value] = true
	}
	return errorsFound
}

func validateV2ItemOrdering(items []V2Item, sections map[string]V2Section, report func(string)) int {
	errorsFound := 0
	sectionOrder := canonicalSectionOrder()
	lastSectionIndex := -1
	var lastItem V2Item
	hasLastItem := false

	for _, item := range items {
		sectionIndex, ok := sectionOrder[item.SectionID]
		if !ok {
			continue
		}

		if sectionIndex < lastSectionIndex {
			report(fmt.Sprintf("Invalid v2 item order: %s appears after a later section item", item.ID))
			errorsFound++
		}

		if hasLastItem && sectionIndex == lastSectionIndex && compareCurriculumItemOrder(lastItem, item) > 0 {
			report(fmt.Sprintf("Invalid v2 item order: %s should appear before %s by path and ID order", item.ID, lastItem.ID))
			errorsFound++
		}

		if section, exists := sections[item.SectionID]; exists {
			cleanItemPath := filepath.ToSlash(filepath.Clean(item.Path))
			cleanSectionPath := filepath.ToSlash(filepath.Clean(section.PathPrefix))
			if cleanItemPath != cleanSectionPath && !strings.HasPrefix(cleanItemPath, cleanSectionPath+"/") {
				report(fmt.Sprintf("Invalid v2 item path prefix: %s -> %s (expected under %s)", item.ID, item.Path, section.PathPrefix))
				errorsFound++
			}
		}

		lastSectionIndex = sectionIndex
		lastItem = item
		hasLastItem = true
	}

	return errorsFound
}

func canonicalSectionOrder() map[string]int {
	order := make(map[string]int, len(canonicalV2Sections))
	for idx, section := range canonicalV2Sections {
		order[section.ID] = idx
	}
	return order
}

func compareCurriculumItemOrder(a, b V2Item) int {
	if pathCompare := compareNaturalPath(a.Path, b.Path); pathCompare != 0 {
		return pathCompare
	}
	return compareCurriculumID(a.ID, b.ID)
}

func compareCurriculumID(a, b string) int {
	aPrefix, aNumber, aOK := splitCurriculumID(a)
	bPrefix, bNumber, bOK := splitCurriculumID(b)
	if aOK && bOK {
		if aPrefix != bPrefix {
			return strings.Compare(aPrefix, bPrefix)
		}
		switch {
		case aNumber < bNumber:
			return -1
		case aNumber > bNumber:
			return 1
		default:
			return 0
		}
	}
	return strings.Compare(a, b)
}

func compareNaturalPath(a, b string) int {
	aParts := strings.Split(filepath.ToSlash(filepath.Clean(a)), "/")
	bParts := strings.Split(filepath.ToSlash(filepath.Clean(b)), "/")
	for idx := 0; idx < len(aParts) && idx < len(bParts); idx++ {
		if partCompare := compareNaturalPart(aParts[idx], bParts[idx]); partCompare != 0 {
			return partCompare
		}
	}

	switch {
	case len(aParts) < len(bParts):
		return -1
	case len(aParts) > len(bParts):
		return 1
	default:
		return 0
	}
}

func compareNaturalPart(a, b string) int {
	aNumber, aRest, aHasNumber := leadingNumber(a)
	bNumber, bRest, bHasNumber := leadingNumber(b)
	if aHasNumber && bHasNumber && aNumber != bNumber {
		if aNumber < bNumber {
			return -1
		}
		return 1
	}
	if aHasNumber && bHasNumber {
		return strings.Compare(aRest, bRest)
	}
	return strings.Compare(a, b)
}

func leadingNumber(value string) (int, string, bool) {
	idx := 0
	for idx < len(value) {
		r := rune(value[idx])
		if !unicode.IsDigit(r) {
			break
		}
		idx++
	}
	if idx == 0 {
		return 0, value, false
	}
	number, err := strconv.Atoi(value[:idx])
	if err != nil {
		return 0, value, false
	}
	return number, value[idx:], true
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

func validateV2LessonSourceStandards(root string, items []V2Item, report func(string)) int {
	errorsFound := 0
	seen := map[string]bool{}

	for _, item := range items {
		if isPlaceholderItem(item) {
			continue
		}

		paths := curriculumGoSurfacePaths(root, item)
		for _, relPath := range paths {
			if seen[relPath] {
				continue
			}
			seen[relPath] = true

			data, cleanPath, err := readRepoFile(root, relPath)
			if err != nil {
				report(fmt.Sprintf("Failed to read v2 source standards surface: %s -> %v", filepath.ToSlash(relPath), err))
				errorsFound++
				continue
			}

			mainPath := filepath.ToSlash(filepath.Join(item.Path, "main.go"))
			if cleanPath == filepath.ToSlash(filepath.Clean(mainPath)) {
				errorsFound += validateCompletedLessonMainHeader(data, cleanPath, item, report)
			}
			errorsFound += validateMachineRoleComments(data, cleanPath, report)
		}
	}

	return errorsFound
}

func curriculumGoSurfacePaths(root string, item V2Item) []string {
	roots := []string{item.Path}
	if strings.TrimSpace(item.StarterPath) != "" {
		roots = append(roots, item.StarterPath)
	}

	seen := map[string]bool{}
	var paths []string
	for _, relRoot := range roots {
		absRoot, cleanRoot, ok := repoPath(root, relRoot)
		if !ok {
			continue
		}
		info, err := os.Stat(absRoot)
		if err != nil {
			continue
		}
		if !info.IsDir() {
			if filepath.Ext(cleanRoot) == ".go" {
				paths = append(paths, cleanRoot)
			}
			continue
		}

		_ = filepath.WalkDir(absRoot, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return nil
			}
			if d.IsDir() {
				if shouldSkipStandardsDir(d.Name()) {
					return filepath.SkipDir
				}
				return nil
			}
			if filepath.Ext(path) != ".go" {
				return nil
			}

			relPath, err := filepath.Rel(root, path)
			if err != nil {
				return nil
			}
			cleanPath := filepath.ToSlash(filepath.Clean(relPath))
			if seen[cleanPath] {
				return nil
			}
			seen[cleanPath] = true
			paths = append(paths, cleanPath)
			return nil
		})
	}

	sort.Strings(paths)
	return paths
}

func validateCompletedLessonMainHeader(data []byte, cleanPath string, item V2Item, report func(string)) int {
	text := strings.ReplaceAll(string(data), "\r\n", "\n")
	errorsFound := 0

	if !strings.HasPrefix(text, "// Copyright (c) 2026 Rasel Hossen\n// Licensed under The Go Engineer License v1.0\n") {
		report(fmt.Sprintf("Invalid v2 lesson source header: %s -> %s must start with copyright and license lines", item.ID, cleanPath))
		errorsFound++
	}

	if sectionNumber := canonicalSectionNumber(item.SectionID); sectionNumber != "" {
		sectionLabel := "Section " + sectionNumber + ":"
		if !strings.Contains(text, sectionLabel) {
			report(fmt.Sprintf("Invalid v2 lesson source header: %s -> %s missing %s", item.ID, cleanPath, sectionLabel))
			errorsFound++
		}
	}
	if !strings.Contains(text, "WHAT YOU'LL LEARN:") {
		report(fmt.Sprintf("Invalid v2 lesson source header: %s -> %s missing WHAT YOU'LL LEARN", item.ID, cleanPath))
		errorsFound++
	}
	if !strings.Contains(text, "WHY THIS MATTERS:") {
		report(fmt.Sprintf("Invalid v2 lesson source header: %s -> %s missing WHY THIS MATTERS", item.ID, cleanPath))
		errorsFound++
	}
	if strings.TrimSpace(item.RunCommand) != "" {
		runCount := countLessonRunHeaders(data)
		if runCount != 1 {
			report(fmt.Sprintf("Invalid v2 lesson run header count: %s -> %s has %d RUN headers (expected 1)", item.ID, cleanPath, runCount))
			errorsFound++
		}
	}
	if !strings.Contains(text, "KEY TAKEAWAY:") {
		report(fmt.Sprintf("Invalid v2 lesson source footer: %s -> %s missing KEY TAKEAWAY", item.ID, cleanPath))
		errorsFound++
	}
	if len(item.NextItems) > 0 && !strings.HasPrefix(item.NextItems[0], "s") && !strings.Contains(text, "NEXT UP:") {
		report(fmt.Sprintf("Invalid v2 lesson source footer: %s -> %s missing NEXT UP", item.ID, cleanPath))
		errorsFound++
	}

	return errorsFound
}

func canonicalSectionNumber(sectionID string) string {
	for _, section := range canonicalV2Sections {
		if section.ID == sectionID {
			return section.Number
		}
	}
	return ""
}

func countLessonRunHeaders(data []byte) int {
	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	count := 0
	for scanner.Scan() {
		if strings.HasPrefix(strings.TrimSpace(scanner.Text()), "// RUN:") {
			count++
		}
	}
	return count
}

func validateMachineRoleComments(data []byte, cleanPath string, report func(string)) int {
	fileSet := token.NewFileSet()
	parsed, err := parser.ParseFile(fileSet, cleanPath, data, parser.ParseComments)
	if err != nil {
		report(fmt.Sprintf("Invalid Go source syntax in v2 source standards surface: %s -> %v", cleanPath, err))
		return 1
	}

	errorsFound := 0
	for _, decl := range parsed.Decls {
		switch typed := decl.(type) {
		case *ast.FuncDecl:
			if shouldSkipMachineRoleFunction(typed.Name.Name) {
				continue
			}
			symbols := []string{typed.Name.Name}
			if receiver := receiverTypeName(typed.Recv); receiver != "" {
				symbols = append([]string{receiver + "." + typed.Name.Name}, symbols...)
			}
			if hasMachineRoleComment(typed.Doc, symbols) {
				continue
			}
			line := fileSet.Position(typed.Pos()).Line
			report(fmt.Sprintf("Missing Machine Role comment: %s:%d -> %s", cleanPath, line, symbols[0]))
			errorsFound++
		case *ast.GenDecl:
			if typed.Tok == token.IMPORT {
				continue
			}
			for _, spec := range typed.Specs {
				switch spec := spec.(type) {
				case *ast.TypeSpec:
					if hasMachineRoleComment(preferredDoc(spec.Doc, typed.Doc), []string{spec.Name.Name}) {
						continue
					}
					line := fileSet.Position(spec.Pos()).Line
					report(fmt.Sprintf("Missing Machine Role comment: %s:%d -> %s", cleanPath, line, spec.Name.Name))
					errorsFound++
				case *ast.ValueSpec:
					if !shouldRequireMachineRoleValueSpec(typed.Tok, spec) {
						continue
					}
					for _, name := range spec.Names {
						if hasMachineRoleComment(preferredDoc(spec.Doc, typed.Doc), []string{name.Name}) {
							continue
						}
						line := fileSet.Position(name.Pos()).Line
						report(fmt.Sprintf("Missing Machine Role comment: %s:%d -> %s", cleanPath, line, name.Name))
						errorsFound++
					}
				}
			}
		}
	}

	return errorsFound
}

func shouldRequireMachineRoleValueSpec(tokenType token.Token, spec *ast.ValueSpec) bool {
	if tokenType == token.CONST {
		return false
	}

	for _, name := range spec.Names {
		if ast.IsExported(name.Name) {
			return true
		}
	}
	if isMachineRoleValueType(spec.Type) {
		return true
	}
	for _, value := range spec.Values {
		if isMachineRoleValueType(value) {
			return true
		}
	}
	return false
}

func isMachineRoleValueType(expr ast.Expr) bool {
	switch typed := expr.(type) {
	case *ast.ArrayType, *ast.MapType, *ast.ChanType, *ast.FuncType:
		return true
	case *ast.UnaryExpr:
		return isMachineRoleValueType(typed.X)
	case *ast.CompositeLit:
		return isMachineRoleValueType(typed.Type)
	case *ast.CallExpr:
		return isMachineRoleValueType(typed.Fun)
	case *ast.SelectorExpr:
		return typed.Sel.Name == "Mutex" ||
			typed.Sel.Name == "RWMutex" ||
			typed.Sel.Name == "WaitGroup" ||
			typed.Sel.Name == "Pool" ||
			typed.Sel.Name == "Map" ||
			typed.Sel.Name == "Context"
	case *ast.Ident:
		return typed.Name == "Mutex" ||
			typed.Name == "RWMutex" ||
			typed.Name == "WaitGroup" ||
			typed.Name == "Pool" ||
			typed.Name == "Map" ||
			typed.Name == "Context"
	default:
		return false
	}
}

func shouldSkipMachineRoleFunction(name string) bool {
	if name == "main" || name == "init" {
		return true
	}
	for _, prefix := range []string{"Test", "Benchmark", "Example", "Fuzz"} {
		if strings.HasPrefix(name, prefix) {
			return true
		}
	}
	return false
}

func receiverTypeName(receiver *ast.FieldList) string {
	if receiver == nil || len(receiver.List) == 0 {
		return ""
	}
	return exprTypeName(receiver.List[0].Type)
}

func exprTypeName(expr ast.Expr) string {
	switch typed := expr.(type) {
	case *ast.Ident:
		return typed.Name
	case *ast.StarExpr:
		return exprTypeName(typed.X)
	case *ast.SelectorExpr:
		return typed.Sel.Name
	case *ast.IndexExpr:
		return exprTypeName(typed.X)
	case *ast.IndexListExpr:
		return exprTypeName(typed.X)
	default:
		return ""
	}
}

func preferredDoc(primary, fallback *ast.CommentGroup) *ast.CommentGroup {
	if primary != nil {
		return primary
	}
	return fallback
}

func hasMachineRoleComment(group *ast.CommentGroup, symbols []string) bool {
	if group == nil {
		return false
	}

	for _, comment := range group.List {
		line := strings.TrimSpace(comment.Text)
		line = strings.TrimPrefix(line, "//")
		line = strings.TrimPrefix(line, "/*")
		line = strings.TrimSuffix(line, "*/")
		line = strings.TrimSpace(line)
		for _, symbol := range symbols {
			if !strings.HasPrefix(line, symbol+" (") {
				continue
			}
			_, role, ok := strings.Cut(line, "):")
			if ok && strings.TrimSpace(role) != "" {
				return true
			}
		}
	}

	return false
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
		resolved := filepath.ToSlash(filepath.Clean(filepath.Join(filepath.Dir(cleanPath), filepath.FromSlash(target))))
		if !pathExists(root, resolved) {
			report(fmt.Sprintf("Broken local doc link: %s -> %s", cleanPath, target))
			errorsFound++
		}
	}

	return errorsFound
}

func validateRepositoryStandards(root string, report func(string)) int {
	errorsFound := 0
	errorsFound += validateCodeStandardsContract(root, report)
	errorsFound += validateMarkdownStandards(root, report)
	errorsFound += validateGoSourceCommentStandards(root, report)
	return errorsFound
}

func validateCodeStandardsContract(root string, report func(string)) int {
	data, cleanPath, err := readRepoFile(root, "CODE-STANDARDS.md")
	if err != nil {
		report(fmt.Sprintf("Invalid code standards contract: CODE-STANDARDS.md is required -> %v", err))
		return 1
	}

	text := string(data)
	errorsFound := 0
	requiredSnippets := []struct {
		name    string
		snippet string
	}{
		{"standard layers", "## Standard Layers"},
		{"production level taxonomy", "Level: Foundation | Core | Production | Stretch"},
		{"canonical coverage command", "go test -coverprofile=coverage.out ./..."},
		{"machine role doc-comment compatibility", "Machine Role comments can satisfy this requirement"},
		{"cross-reference alert rules", "do not use legacy `Forward Reference` or `Backward Reference` labels"},
		{"curriculum registry standard", "## Curriculum Registry Standard"},
		{"lesson proof surface", "## Lesson Proof Surface"},
		{"production-shaped code", "## Production-Shaped Code"},
	}

	for _, required := range requiredSnippets {
		if strings.Contains(text, required.snippet) {
			continue
		}
		report(fmt.Sprintf("Invalid code standards contract: %s missing %s", cleanPath, required.name))
		errorsFound++
	}

	if strings.Contains(text, "AGENTS.md") {
		report(fmt.Sprintf("Invalid code standards contract: %s references maintainer-only AGENTS.md", cleanPath))
		errorsFound++
	}

	return errorsFound
}

func validateMarkdownStandards(root string, report func(string)) int {
	errorsFound := 0

	walkErr := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			if shouldSkipStandardsDir(d.Name()) {
				return filepath.SkipDir
			}
			return nil
		}

		if filepath.Ext(path) != ".md" {
			return nil
		}

		relPath, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}

		cleanPath := filepath.ToSlash(filepath.Clean(relPath))
		errorsFound += validateMarkdownStandardsFile(root, cleanPath, report)
		return nil
	})
	if walkErr != nil {
		report(fmt.Sprintf("Failed to scan repository standards: %v", walkErr))
		return errorsFound + 1
	}

	return errorsFound
}

func validateMarkdownStandardsFile(root, relPath string, report func(string)) int {
	data, cleanPath, err := readRepoFile(root, relPath)
	if err != nil {
		report(fmt.Sprintf("Failed to read markdown standards surface: %s -> %v", filepath.ToSlash(relPath), err))
		return 1
	}

	lines := strings.Split(strings.ReplaceAll(string(data), "\r\n", "\n"), "\n")
	errorsFound := 0
	for idx := 0; idx < len(lines); idx++ {
		lineNo := idx + 1
		line := strings.TrimSpace(lines[idx])

		if staleCoverageCommandPattern.MatchString(line) {
			report(fmt.Sprintf("Invalid verification command: %s:%d uses coverprofile without '='; use go test -coverprofile=coverage.out ./...", cleanPath, lineNo))
			errorsFound++
		}

		if strings.Contains(line, "Foundation | Core | Stretch") && !strings.Contains(line, "Production") {
			report(fmt.Sprintf("Invalid level taxonomy: %s:%d omits Production", cleanPath, lineNo))
			errorsFound++
		}

		if legacyReferenceHeadingPattern.MatchString(line) {
			report(fmt.Sprintf("Invalid markdown cross-reference heading: %s:%d uses Forward/Backward Reference heading; use [!TIP] or [!NOTE]", cleanPath, lineNo))
			errorsFound++
		}

		matches := markdownAlertPattern.FindStringSubmatch(line)
		if len(matches) != 2 {
			continue
		}

		blockText := line + "\n"
		for next := idx + 1; next < len(lines); next++ {
			nextLine := strings.TrimSpace(lines[next])
			if !strings.HasPrefix(nextLine, ">") {
				break
			}
			blockText += nextLine + "\n"
		}

		errorsFound += validateMarkdownAlertBlock(cleanPath, lineNo, matches[1], blockText, report)
	}

	return errorsFound
}

func validateMarkdownAlertBlock(cleanPath string, lineNo int, alertType, blockText string, report func(string)) int {
	if !curriculumReferencePattern.MatchString(blockText) {
		return 0
	}

	errorsFound := 0
	if alertType != "NOTE" && alertType != "TIP" {
		report(fmt.Sprintf("Invalid markdown cross-reference alert: %s:%d uses [!%s] for curriculum reference; use [!NOTE] or [!TIP]", cleanPath, lineNo, alertType))
		errorsFound++
	}

	if !markdownReadmeLinkPattern.MatchString(blockText) {
		report(fmt.Sprintf("Invalid markdown cross-reference link: %s:%d references a curriculum ID without a clickable README.md link", cleanPath, lineNo))
		errorsFound++
	}

	return errorsFound
}

func validateGoSourceCommentStandards(root string, report func(string)) int {
	errorsFound := 0

	walkErr := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			if shouldSkipStandardsDir(d.Name()) {
				return filepath.SkipDir
			}
			return nil
		}

		if filepath.Ext(path) != ".go" {
			return nil
		}

		relPath, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}

		cleanPath := filepath.ToSlash(filepath.Clean(relPath))
		data, _, err := readRepoFile(root, cleanPath)
		if err != nil {
			report(fmt.Sprintf("Failed to read Go source standards surface: %s -> %v", cleanPath, err))
			return nil
		}

		scanner := bufio.NewScanner(strings.NewReader(string(data)))
		lineNo := 0
		for scanner.Scan() {
			lineNo++
			line := strings.TrimSpace(scanner.Text())
			if !strings.HasPrefix(line, "//") {
				continue
			}
			if strings.Contains(line, "[!NOTE]") || strings.Contains(line, "[!TIP]") {
				report(fmt.Sprintf("Invalid Go source cross-reference comment: %s:%d uses markdown alert syntax", cleanPath, lineNo))
				errorsFound++
			}
		}
		if err := scanner.Err(); err != nil {
			report(fmt.Sprintf("Failed to scan Go source standards surface: %s -> %v", cleanPath, err))
			errorsFound++
		}

		return nil
	})
	if walkErr != nil {
		report(fmt.Sprintf("Failed to scan Go source standards: %v", walkErr))
		return errorsFound + 1
	}

	return errorsFound
}

func shouldSkipStandardsDir(name string) bool {
	switch name {
	case ".git", "vendor", ".agents", ".cache", ".opencode", "temp":
		return true
	default:
		return false
	}
}
