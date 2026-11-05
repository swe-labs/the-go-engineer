package main

import (
	"fmt"
	"os"
	"path/filepath"
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

// validateSESENavigation verifies that module root READMEs adhere to the Single-Entry, Single-Exit (SESE) standard:
// 1. Entry Invariant: Module root READMEs must end with "Start here" or "NEXT UP" pointing to Lesson 01.
// 2. No Escape Hatch: Module root READMEs must not offer direct exit links jumping to downstream modules.
// 3. Exit Invariant: For completed modules (00, 01), the assessment checkpoint must contain the exit link to the next module.
func validateSESENavigation() {
	fmt.Println("=== validate-sese-navigation ===")
	ok := true

	var core CoreBundle
	readJSON("path.core.json", &core)

	// Modules audited and locked into SESE
	lockedModules := map[string]string{
		"module-00": "../01-computers-terminal-git-web/README.md",
		"module-01": "../02-go-setup-tooling/README.md",
		"module-02": "../03-programming-fundamentals/README.md",
	}

	for _, mod := range core.Modules {
		modDir := resolveModuleDir(mod.Path)
		rootReadme := filepath.Join(modDir, "README.md")

		data, err := os.ReadFile(rootReadme)
		if err != nil {
			warnf("%s: root README not found at %s", mod.ID, rootReadme)
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

		// For locked modules, assert strictly
		if _, isLocked := lockedModules[mod.ID]; isLocked {
			if !strings.HasPrefix(lastH2, "Start here") && !strings.HasPrefix(lastH2, "NEXT UP") {
				warnf("%s root README:%d expected final heading '## Start here', got '## %s'", mod.ID, lastH2LineNo, lastH2)
				ok = false
			}

			// Must not contain direct skip link to next module at root
			if strings.Contains(content, "## Next module") {
				warnf("%s root README contains '## Next module' escape hatch; module root must only contain '## Start here'", mod.ID)
				ok = false
			}
		}
	}

	// Verify exit gates on locked module assessments
	for modID, nextModuleLink := range lockedModules {
		var assessDir string
		for _, mod := range core.Modules {
			if mod.ID == modID {
				assessDir = filepath.Join(resolveModuleDir(mod.Path), "assessments")
				break
			}
		}
		if assessDir == "" {
			continue
		}

		// Find the checkpoint assessment README
		checkpointReadme := filepath.Join(assessDir, "checkpoint", "README.md")
		if modID == "module-01" {
			checkpointReadme = filepath.Join(assessDir, "shell-git-lab", "README.md")
		}
		if modID == "module-02" {
			checkpointReadme = filepath.Join(assessDir, "tooling-failure-lab", "README.md")
		}

		cData, err := os.ReadFile(checkpointReadme)
		if err != nil {
			warnf("%s missing exit assessment README at %s", modID, checkpointReadme)
			ok = false
			continue
		}

		cContent := string(cData)
		if !strings.Contains(cContent, nextModuleLink) {
			warnf("%s assessment exit gate at %s missing forward link '%s'", modID, checkpointReadme, nextModuleLink)
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
			warnf("%s: missing verify.sh in %s", p.ID, projDir)
			ok = false
		} else if shInfo.Size() < 50 {
			warnf("%s: verify.sh is too small (%d bytes)", p.ID, shInfo.Size())
			ok = false
		}

		// Check verify.ps1
		psPath := filepath.Join(projDir, "verify.ps1")
		if psInfo, err := os.Stat(psPath); err != nil {
			warnf("%s: missing verify.ps1 in %s", p.ID, projDir)
			ok = false
		} else if psInfo.Size() < 50 {
			warnf("%s: verify.ps1 is too small (%d bytes)", p.ID, psInfo.Size())
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
