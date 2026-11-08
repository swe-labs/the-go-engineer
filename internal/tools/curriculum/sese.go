package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// resolveModuleDir finds the physical directory for a module path
func resolveModuleDir(modPath string) string {
	candidate1 := filepath.Join(curriculumDir, "modules", modPath)
	if info, err := os.Stat(candidate1); err == nil && info.IsDir() {
		return candidate1
	}
	candidate2 := filepath.Join(curriculumDir, modPath)
	if info, err := os.Stat(candidate2); err == nil && info.IsDir() {
		return candidate2
	}
	return candidate1
}

// deriveNextModuleLinks maps each module to the relative link of the module
// that follows it in canonical order. The terminal module has no entry.
func deriveNextModuleLinks(modules []Module) map[string]string {
	sorted := append([]Module{}, modules...)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].Order < sorted[j].Order })
	links := make(map[string]string)
	for i := 0; i+1 < len(sorted); i++ {
		links[sorted[i].ID] = "../" + filepath.Base(sorted[i+1].Path) + "/README.md"
	}
	return links
}

// validateSESENavigation verifies that module root READMEs adhere to the Single-Entry, Single-Exit (SESE) standard:
//  1. Entry Invariant: Every module root README must provide a navigation footer. Locked modules end with "Start here" or "NEXT UP";
//     unlocked modules end with "Next module" linking to the next module in the canonical order.
//  2. No Escape Hatch: Locked module root READMEs must not offer a direct "Next module" exit link.
//  3. Exit Invariant: For every locked module (see curriculum/expected.json), the assessment checkpoint must contain the exit link to the next module.
func validateSESENavigation() {
	fmt.Println("=== validate-sese-navigation ===")
	ok := true

	var core CoreBundle
	readJSON("path.core.json", &core)

	// Derive the forward-link targets from module metadata and order, so the
	// chain never has to be duplicated by hand when a module is added or renamed.
	nextModuleLinks := deriveNextModuleLinks(core.Modules)

	for _, mod := range core.Modules {
		modDir := resolveModuleDir(mod.Path)
		rootReadme := filepath.Join(modDir, "README.md")

		data, err := os.ReadFile(rootReadme)
		if err != nil {
			errorf("%s: root README not found at %s", mod.ID, rootReadme)
			ok = false
			continue
		}

		content := string(data)
		lines := strings.Split(content, "\n")

		// Check the last H2 heading in the root README
		var lastH2 string
		var lastH2LineNo int
		for i, line := range lines {
			trimmed := strings.TrimSpace(line)
			if strings.HasPrefix(trimmed, "## ") {
				lastH2 = strings.TrimPrefix(trimmed, "## ")
				lastH2 = strings.TrimSpace(lastH2)
				lastH2LineNo = i + 1
			}
		}

		// Every module must end with a navigable footer. The exact contract
		// depends on the module's position in the chain, not on whether it is
		// locked: a dead-end page is never acceptable.
		nextLink, hasNext := nextModuleLinks[mod.ID]
		if lockedModuleSet()[mod.ID] {
			if !strings.HasPrefix(lastH2, "Start here") && !strings.HasPrefix(lastH2, "NEXT UP") {
				errorf("%s root README:%d expected final heading '## Start here', got '## %s'", mod.ID, lastH2LineNo, lastH2)
				ok = false
			}
			if strings.Contains(content, "## Next module") {
				errorf("%s root README contains '## Next module' escape hatch; locked module root must only contain '## Start here'", mod.ID)
				ok = false
			}
		} else if hasNext {
			if !strings.HasPrefix(lastH2, "Next module") {
				errorf("%s root README:%d expected final heading '## Next module', got '## %s'", mod.ID, lastH2LineNo, lastH2)
				ok = false
			}
			if !strings.Contains(content, nextLink) {
				errorf("%s root README must link to the next module at '%s'", mod.ID, nextLink)
				ok = false
			}
		} else {
			if !strings.HasPrefix(lastH2, "Completion") {
				errorf("%s root README:%d is terminal and must end with '## Completion', got '## %s'", mod.ID, lastH2LineNo, lastH2)
				ok = false
			}
		}
		if !strings.Contains(content, "](") {
			errorf("%s root README has no Markdown links; the module overview is a dead end", mod.ID)
			ok = false
		}
	}

	// Verify exit gates on locked module assessments
	for modID := range lockedModuleSet() {
		nextModuleLink, hasLink := nextModuleLinks[modID]
		if !hasLink {
			errorf("%s is locked but has no next-module exit-gate link configured", modID)
			ok = false
			continue
		}
		foundExitGate := false
		var checkedFiles []string

		for _, item := range core.Items {
			if item.ModuleID == modID && item.Type == "assessment" && item.Files != nil && item.Files.ReadmePath != "" {
				fullPath := resolvePath(item.Files.ReadmePath)
				checkedFiles = append(checkedFiles, item.Files.ReadmePath)
				cData, err := os.ReadFile(fullPath)
				if err != nil {
					continue
				}
				if strings.Contains(string(cData), nextModuleLink) {
					foundExitGate = true
					break
				}
			}
		}

		if !foundExitGate {
			errorf("%s assessment exit gate missing forward link '%s' (checked %v)", modID, nextModuleLink, checkedFiles)
			ok = false
		}
	}

	if ok {
		fmt.Println("  Single-Entry, Single-Exit (SESE) navigation flow verified.")
	} else {
		fmt.Println("  SESE navigation validation completed with warnings/errors.")
	}
}

// validateProjectScripts checks cross-platform verification script parity (PowerShell .ps1 and Bash .sh)
// for all projects with verification commands.
func validateProjectScripts() {
	fmt.Println("=== validate-project-scripts ===")
	ok := true

	var rawProjects struct {
		Projects []struct {
			ID           string `json:"id"`
			ModuleID     string `json:"module_id"`
			Slug         string `json:"slug"`
			Verification any    `json:"verification"`
		} `json:"projects"`
	}
	readJSON("projects.json", &rawProjects)

	var core CoreBundle
	readJSON("path.core.json", &core)

	modPathMap := make(map[string]string)
	for _, m := range core.Modules {
		modPathMap[m.ID] = m.Path
	}

	scriptProjects := 0
	for _, p := range rawProjects.Projects {
		vMap, isMap := p.Verification.(map[string]any)
		if !isMap {
			continue
		}
		runCmd, _ := vMap["run_command"].(string)
		if !strings.Contains(runCmd, "verify") {
			continue
		}
		scriptProjects++

		mPath := modPathMap[p.ModuleID]
		projDir := filepath.Join(resolveModuleDir(mPath), "projects", p.Slug)

		// Check verify.sh
		shPath := filepath.Join(projDir, "verify.sh")
		if shInfo, err := os.Stat(shPath); err != nil {
			errorf("%s: missing verify.sh in %s", p.ID, projDir)
			ok = false
		} else if shInfo.Size() < 50 {
			errorf("%s: verify.sh is too small (%d bytes)", p.ID, shInfo.Size())
			ok = false
		}

		// Check verify.ps1
		psPath := filepath.Join(projDir, "verify.ps1")
		if psInfo, err := os.Stat(psPath); err != nil {
			errorf("%s: missing verify.ps1 in %s", p.ID, projDir)
			ok = false
		} else if psInfo.Size() < 50 {
			errorf("%s: verify.ps1 is too small (%d bytes)", p.ID, psInfo.Size())
			ok = false
		}
	}

	fmt.Printf("  Verified script parity across %d script-backed project harnesses.\n", scriptProjects)
	if ok {
		fmt.Println("  All project verification scripts verified.")
	} else {
		fmt.Println("  Project script validation completed with warnings/errors.")
	}
}
