package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// contractViolation is one issue found while validating an item README against
// its resolved contract family.
type contractViolation struct {
	ItemID string
	Family string
	Kind   string // contract-id | heading | forbidden-heading | code | file
	Detail string
	Sev    severity
}

// validateReadmeContracts resolves every curriculum item to a README contract
// family from readme.contracts.json and verifies:
//   - the declared contract_id exists and matches the family resolved from the
//     item's type/subtype,
//   - the README satisfies the family's required heading core (after alias
//     normalization) and omits any subtype-forbidden headings,
//   - code_policy (required/forbidden/optional) and content_contract runnable
//     expectations,
//   - file_policy (readme/main/test/starter/solution) derived from
//     content_contract flags.
//
// Structural failures always block. Content-shape failures are tiered by
// severity policy: locked modules always block, and other items follow the
// registry's by_status map (draft/scaffolded skip, review/reviewed warn,
// ready-and-later block). --strict promotes all content-shape warnings to
// errors; --advisory skips content-shape entirely.
func validateReadmeContracts() {
	fmt.Println("=== validate-readme-contracts ===")
	reg := loadRegistry()

	var core CoreBundle
	var electives ElectiveBundle
	readJSON("path.core.json", &core)
	readJSON("path.electives.json", &electives)

	items := make([]Item, 0, len(core.Items)+len(electives.Items))
	items = append(items, core.Items...)
	items = append(items, electives.Items...)

	aliases := reg.aliasesOf() // aliasID -> canonical family id

	var violations []contractViolation
	add := func(v contractViolation) {
		if v.Sev == sevSkip {
			return
		}
		violations = append(violations, v)
	}

	checked := 0
	for _, item := range items {
		c, matched := reg.resolveItemContract(item)
		if !matched {
			add(contractViolation{item.ID, "", "contract-id",
				"no single contract family matches the item's type/subtype", sevError})
			continue
		}

		// 1. contract_id integrity (always structural).
		if item.ReadmeContract != nil {
			declared := item.ReadmeContract.ContractID
			if declared != "" {
				canonical := declared
				if aliasTarget, isAlias := aliases[declared]; isAlias {
					canonical = aliasTarget
				}
				if reg.contractByID(canonical) == nil {
					add(contractViolation{item.ID, c.ID, "contract-id",
						fmt.Sprintf("declared contract_id '%s' is not defined in readme.contracts.json", declared), sevError})
				} else if canonical != c.ID {
					add(contractViolation{item.ID, c.ID, "contract-id",
						fmt.Sprintf("declared contract_id '%s' does not match resolved family '%s'", declared, c.ID), sevError})
				}
			}
		}

		// 2. README must exist (structural).
		if item.Files == nil || item.Files.ReadmePath == "" {
			add(contractViolation{item.ID, c.ID, "file", "missing files.readme_path in metadata", sevError})
			continue
		}
		readme, err := os.ReadFile(resolvePath(item.Files.ReadmePath))
		if err != nil {
			add(contractViolation{item.ID, c.ID, "file", fmt.Sprintf("cannot read README: %v", err), sevError})
			continue
		}
		checked++

		sevs := contentShapeSeverity(item)

		// Layered subtype modifier.
		codePolicy := c.CodePolicy
		required := append([]string{}, c.RequiredHeadings...)
		var forbidden []string
		if mod, hasMod := c.SubtypeModifiers[item.Subtype]; hasMod {
			if mod.CodePolicy != "" {
				codePolicy = mod.CodePolicy
			}
			required = append(required, mod.RequiredHeadings...)
			forbidden = mod.ForbiddenHeadings
		}
		if c.ID == "lesson.v3" {
			required = sortLessonRequired(required)
		}

		headings, hasGo := scanReadmeHeadings(readme)

		// 3. Required headings (content-shape).
		found := normalizeHeadings(headings, c.HeadingAliases)
		lastReqIdx := -1
		var lastReqHeading string
		for _, req := range required {
			if !headingSatisfied(req, found) {
				add(contractViolation{item.ID, c.ID, "heading",
					fmt.Sprintf("missing required heading '## %s'", req), sevs})
			} else {
				reqIdx := findHeadingIndex(req, headings, c.HeadingAliases)
				if reqIdx != -1 {
					if reqIdx < lastReqIdx {
						add(contractViolation{item.ID, c.ID, "heading-order",
							fmt.Sprintf("required heading '## %s' appears before preceding required heading '## %s'", req, lastReqHeading), sevs})
					} else {
						lastReqIdx = reqIdx
						lastReqHeading = req
					}
				}
			}
		}

		// 4. Forbidden headings (content-shape).
		for _, fb := range forbidden {
			for _, h := range headings {
				if headingMatches(h, fb) {
					add(contractViolation{item.ID, c.ID, "forbidden-heading",
						fmt.Sprintf("contains forbidden heading '## %s'", h), sevs})
				}
			}
		}

		// 5. Code policy (content-shape).
		switch codePolicy {
		case "required":
			if !hasGo {
				add(contractViolation{item.ID, c.ID, "code",
					"code_policy=required but the README has no Go code block", sevs})
			}
		case "forbidden":
			if hasGo {
				add(contractViolation{item.ID, c.ID, "code",
					"code_policy=forbidden but the README contains a Go code block", sevs})
			}
		}
		if item.ContentContract != nil && item.ContentContract.RunnableRequired && !hasGo {
			add(contractViolation{item.ID, c.ID, "code",
				"content_contract.runnable_required but the README has no Go code block", sevs})
		}

		// 6. File policy (content-shape), derived from content_contract flags.
		for _, key := range sortedKeys(c.FilePolicy) {
			token := c.FilePolicy[key]
			need := token == "required"
			if strings.HasPrefix(token, "when:") {
				need = contentFlag(item, strings.TrimPrefix(token, "when:"))
			}
			if !need {
				continue
			}
			field, path := fileField(item, key)
			if path == "" {
				add(contractViolation{item.ID, c.ID, "file",
					fmt.Sprintf("file_policy requires '%s' but files.%s is empty", key, field), sevs})
			} else if !pathExists(resolvePath(path)) {
				add(contractViolation{item.ID, c.ID, "file",
					fmt.Sprintf("file_policy requires '%s' but '%s' is missing on disk", key, path), sevs})
			}
		}
	}

	// Document-level contracts (root README, module overviews, flagship, and
	// companion stages) govern whole documents rather than registry items, so
	// they are checked separately through their selector "doc" role token.
	docChecked := validateDocumentContracts(reg, core, electives, add)

	// The curriculum catalog is machine-generated. Enforce its generated
	// boundary markers so a generator run can never silently overwrite
	// human-authored prose without a visible, checkable signal.
	catalogPath := filepath.Join(curriculumDir, "README.md")
	if catalogData, err := os.ReadFile(catalogPath); err != nil {
		add(contractViolation{"curriculum/README.md", "root.v3", "generated-markers",
			fmt.Sprintf("cannot read generated catalog: %v", err), sevError})
	} else if !hasGeneratedMarkers(catalogData) {
		add(contractViolation{"curriculum/README.md", "root.v3", "generated-markers",
			"generated catalog must be wrapped in curriculum:generated start/end markers", sevError})
	}

	reportContractViolations(violations)

	fmt.Printf("  Verified README contracts across %d items and %d documents (%d violations).\n", checked, docChecked, len(violations))
	if len(violations) == 0 {
		fmt.Println("  All items comply with README contract specifications.")
	} else {
		fmt.Println("  README contract issues reported above.")
	}
}

// hasGeneratedMarkers reports whether a file carries both curriculum generated
// boundary markers, the signal that the file is machine-generated and safe to
// regenerate.
func hasGeneratedMarkers(data []byte) bool {
	s := string(data)
	return strings.Contains(s, "<!-- curriculum:generated:start -->") &&
		strings.Contains(s, "<!-- curriculum:generated:end -->")
}

// validateDocumentContracts checks document-level contract families (those with
// a selector "doc" role token) against their target READMEs: the repository
// root, each module overview, the flagship overview, and the companion-stage
// guides. It returns the number of documents checked and reports any heading
// violations through add.
func validateDocumentContracts(reg *Registry, core CoreBundle, electives ElectiveBundle, add func(contractViolation)) int {
	companionRoot := filepath.Join(curriculumDir, "companion-frontend")

	type target struct {
		id       string // report identifier (module id or path label)
		path     string
		moduleID string
		family   *Contract
	}

	var targets []target
	modules := append(append([]Module{}, core.Modules...), electives.Modules...)

	for i := range reg.Contracts {
		c := &reg.Contracts[i]
		switch c.Selector.Doc {
		case "":
			continue
		case "repository-root":
			targets = append(targets, target{
				"repository-root", filepath.Join(filepath.Dir(curriculumDir), "README.md"), "", c,
			})
		case "module", "flagship-module":
			for _, m := range modules {
				if m.ReadmeContract == nil || m.ReadmeContract.ContractID == "" {
					continue
				}
				fam := reg.contractByID(canonicalContractID(reg, m.ReadmeContract.ContractID))
				if fam == nil || fam.ID != c.ID {
					continue
				}
				targets = append(targets, target{
					m.ID, filepath.Join(resolveModuleDir(m.Path), "README.md"), m.ID, c,
				})
			}
		case "companion-stage":
			for _, p := range companionReadmes(companionRoot) {
				targets = append(targets, target{
					"companion:" + filepath.Base(filepath.Dir(p)), p, "", c,
				})
			}
		}
	}

	checked := 0
	for _, t := range targets {
		data, err := os.ReadFile(t.path)
		if err != nil {
			add(contractViolation{t.id, t.family.ID, "file",
				fmt.Sprintf("cannot read document README: %v", err), sevError})
			continue
		}
		checked++
		sev := documentShapeSeverity(t.moduleID)
		headings, _ := scanReadmeHeadings(data)
		found := normalizeHeadings(headings, t.family.HeadingAliases)
		lastReqIdx := -1
		var lastReqHeading string
		for _, req := range t.family.RequiredHeadings {
			if !headingSatisfied(req, found) {
				add(contractViolation{t.id, t.family.ID, "heading",
					fmt.Sprintf("missing required heading '## %s'", req), sev})
			} else {
				reqIdx := findHeadingIndex(req, headings, t.family.HeadingAliases)
				if reqIdx != -1 {
					if reqIdx < lastReqIdx {
						add(contractViolation{t.id, t.family.ID, "heading-order",
							fmt.Sprintf("required heading '## %s' appears before preceding required heading '## %s'", req, lastReqHeading), sev})
					} else {
						lastReqIdx = reqIdx
						lastReqHeading = req
					}
				}
			}
		}
	}
	return checked
}

// canonicalContractID resolves a declared contract_id (possibly a selector
// alias such as elective-module.v3) to its canonical family id.
func canonicalContractID(reg *Registry, id string) string {
	if alias, ok := reg.aliasesOf()[id]; ok {
		return alias
	}
	return id
}

// documentShapeSeverity mirrors contentShapeSeverity for document-level
// contracts, which carry no item status: locked modules still block, and every
// other document follows the same advisory/strict rollout switch.
func documentShapeSeverity(moduleID string) severity {
	if advisoryMode {
		return sevSkip
	}
	if strictMode {
		return sevError
	}
	sp := loadRegistry().SeverityPolicy
	if sp.Enforcement != "enforce" {
		return sevWarning
	}
	if moduleID != "" && lockedModuleSet()[moduleID] {
		return sevError
	}
	return sevWarning
}

// companionReadmes returns every README.md under the companion track directory,
// sorted for deterministic reporting.
func companionReadmes(root string) []string {
	var out []string
	_ = filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if !d.IsDir() && d.Name() == "README.md" {
			out = append(out, path)
		}
		return nil
	})
	sort.Strings(out)
	return out
}

// scanReadmeHeadings returns all H2 headings and whether a Go code fence exists.
func scanReadmeHeadings(data []byte) (headings []string, hasGo bool) {
	scanner := bufio.NewScanner(bytes.NewReader(data))
	scanner.Buffer(make([]byte, 1024*1024), 1024*1024)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "```go") || strings.HasPrefix(line, "```golang") {
			hasGo = true
		}
		if strings.HasPrefix(line, "## ") {
			headings = append(headings, strings.TrimSpace(strings.TrimPrefix(line, "## ")))
		}
	}
	return headings, hasGo
}

func normalizeHeading(s string) string {
	s = strings.TrimSpace(s)
	s = strings.TrimSpace(strings.TrimLeft(s, "#"))
	return strings.ToLower(s)
}

// normalizeHeadings converts raw H2 headings into a normalized set, applying the
// contract's alias map so variants count toward their canonical heading.
func normalizeHeadings(headings []string, aliases map[string]string) map[string]bool {
	lookup := make(map[string]string, len(aliases))
	for variant, canonical := range aliases {
		lookup[normalizeHeading(variant)] = canonical
	}
	found := make(map[string]bool, len(headings)*2)
	for _, h := range headings {
		n := normalizeHeading(h)
		found[n] = true
		if canon, ok := lookup[n]; ok {
			found[normalizeHeading(canon)] = true
		}
	}
	return found
}

// headingEnumerator strips a leading ordered/unordered list marker that authors
// sometimes bake into a heading, e.g. "7. 4-Tier Assessment Rubric" -> "4-Tier...".
var headingEnumerator = regexp.MustCompile(`^\s*(?:\d+[.)]?|[-*+])\s+`)

// headingTokens splits a normalized heading into lowercased alphanumeric words,
// discarding punctuation so "4-tier" -> ["4","tier"] and "rubric" -> ["rubric"].
func headingTokens(s string) []string {
	return strings.FieldsFunc(s, func(r rune) bool {
		return !(r >= 'a' && r <= 'z') && !(r >= '0' && r <= '9')
	})
}

// headingMatches reports whether a present heading satisfies a required one.
// Matching is whole-word and order-preserving: the required heading must be a
// contiguous token subsequence of the present heading (case-insensitive, after
// stripping any list marker). This prevents substring false positives such as
// "Submission checklist" satisfying "Mission", while still allowing a prefixed
// heading such as "Grading rubric" to satisfy "Rubric".
func headingMatches(present, required string) bool {
	p := headingEnumerator.ReplaceAllString(normalizeHeading(present), "")
	r := headingEnumerator.ReplaceAllString(normalizeHeading(required), "")
	if p == r {
		return true
	}
	pt := headingTokens(p)
	rt := headingTokens(r)
	if len(rt) == 0 || len(rt) > len(pt) {
		return false
	}
	for i := 0; i+len(rt) <= len(pt); i++ {
		match := true
		for j := range rt {
			if pt[i+j] != rt[j] {
				match = false
				break
			}
		}
		if match {
			return true
		}
	}
	return false
}

func headingSatisfied(required string, found map[string]bool) bool {
	for f := range found {
		if headingMatches(f, required) {
			return true
		}
	}
	return false
}

// findHeadingIndex returns the index of the first heading in headings that satisfies required.
func findHeadingIndex(required string, headings []string, aliases map[string]string) int {
	for idx, h := range headings {
		if headingMatches(h, required) {
			return idx
		}
		for variant, canon := range aliases {
			if headingMatches(required, canon) && headingMatches(h, variant) {
				return idx
			}
		}
	}
	return -1
}

var canonicalLessonSequence = []string{
	"Learning objective",
	"Why this matters",
	"Mental model",
	"Visual model",
	"Core idea",
	"Under the hood",
	"How Go uses it",
	"Go example",
	"Step-by-step execution",
	"Runtime evidence",
	"Common mistakes",
	"Debugging walkthrough",
	"Production notes",
	"Performance implications",
	"Practice task",
	"Tests / verification",
	"Review questions",
	"NEXT UP",
}

func lessonHeadingRank(h string) int {
	for i, c := range canonicalLessonSequence {
		if headingMatches(h, c) {
			return i
		}
	}
	return 999
}

func sortLessonRequired(reqs []string) []string {
	res := append([]string{}, reqs...)
	sort.SliceStable(res, func(i, j int) bool {
		return lessonHeadingRank(res[i]) < lessonHeadingRank(res[j])
	})
	return res
}

func contentFlag(item Item, flag string) bool {
	cc := item.ContentContract
	if cc == nil {
		return false
	}
	switch flag {
	case "readme_required":
		return cc.ReadmeRequired
	case "runnable_required":
		return cc.RunnableRequired
	case "tests_required":
		return cc.TestsRequired
	case "visual_model_required":
		return cc.VisualModelRequired
	case "machine_view_required":
		return cc.MachineViewRequired
	case "common_mistakes_required":
		return cc.CommonMistakesRequired
	case "production_notes_required":
		return cc.ProductionNotesRequired
	case "review_questions_required":
		return cc.ReviewQuestionsRequired
	case "portfolio_artifact_required":
		return cc.PortfolioArtifactRequired
	}
	return false
}

func fileField(item Item, key string) (field, path string) {
	if item.Files == nil {
		return key + "_path", ""
	}
	switch key {
	case "readme":
		return "readme_path", item.Files.ReadmePath
	case "main":
		return "main_path", item.Files.MainPath
	case "test":
		return "test_path", item.Files.TestPath
	case "starter":
		return "starter_path", item.Files.StarterPath
	case "solution":
		return "solution_path", item.Files.SolutionPath
	}
	return key + "_path", ""
}

const maxContractPrint = 40

// reportContractViolations prints violations (errors first), counts them, prints
// a by-kind summary, and writes a JSON report when --report is set.
func reportContractViolations(violations []contractViolation) {
	sort.SliceStable(violations, func(i, j int) bool {
		if violations[i].Sev != violations[j].Sev {
			return violations[i].Sev > violations[j].Sev
		}
		return violations[i].ItemID < violations[j].ItemID
	})

	byKind := make(map[string]int)
	errTotal, warnTotal := 0, 0
	printedErr, printedWarn := 0, 0
	for _, v := range violations {
		byKind[v.Kind]++
		switch v.Sev {
		case sevError:
			errTotal++
			errorCount++
			recordValidation("error", v.ItemID, v.Family, v.Kind, v.Detail)
			if printedErr < maxContractPrint {
				fmt.Fprintf(os.Stderr, "ERROR: [%s/%s] %s: %s\n", v.Family, v.Kind, v.ItemID, v.Detail)
				printedErr++
			}
		case sevWarning:
			warnTotal++
			warningCount++
			recordValidation("warning", v.ItemID, v.Family, v.Kind, v.Detail)
			if printedWarn < maxContractPrint {
				fmt.Fprintf(os.Stderr, "WARN: [%s/%s] %s: %s\n", v.Family, v.Kind, v.ItemID, v.Detail)
				printedWarn++
			}
		}
	}

	if len(violations) > 0 {
		fmt.Printf("  contract violations by kind:")
		for _, k := range sortedKeys(byKind) {
			fmt.Printf(" %s=%d", k, byKind[k])
		}
		fmt.Printf(" (errors=%d warnings=%d)\n", errTotal, warnTotal)
		if errTotal > maxContractPrint || warnTotal > maxContractPrint {
			fmt.Printf("  (output truncated at %d per severity; use --report for the full list)\n", maxContractPrint)
		}
	}
}

// sortedKeys returns the sorted keys of any string-keyed map.
func sortedKeys[V any](m map[string]V) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
