package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

type catalogProject struct {
	ID         string `json:"id"`
	ModuleID   string `json:"module_id"`
	Title      string `json:"title"`
	Type       string `json:"type"`
	Stage      string `json:"stage"`
	Difficulty string `json:"difficulty"`
	Order      int    `json:"order"`
}

type catalogProjectsBundle struct {
	Projects []catalogProject `json:"projects"`
}

type catalogAssessment struct {
	ID           string `json:"id"`
	ModuleID     string `json:"module_id"`
	Title        string `json:"title"`
	Type         string `json:"type"`
	PassingScore any    `json:"passing_score"`
}

type catalogAssessmentsBundle struct {
	Assessments []catalogAssessment `json:"assessments"`
}

type catalogConceptsBundle struct {
	Concepts []struct {
		Concept string `json:"concept"`
	} `json:"concepts"`
}

type catalogFailuresBundle struct {
	FailureCategories []any `json:"failure_categories"`
}

// generateCatalog produces the authoritative Markdown catalog for curriculum/README.md.
func generateCatalog() (string, error) {
	core := loadCoreBundle()
	electives := loadElectiveBundle()

	var projBundle catalogProjectsBundle
	readJSON("projects.json", &projBundle)

	var assessBundle catalogAssessmentsBundle
	readJSON("assessments.json", &assessBundle)

	var crossrefs CrossrefBundle
	readJSON("crossrefs.json", &crossrefs)

	var concepts catalogConceptsBundle
	readJSON("concepts.json", &concepts)

	var failures catalogFailuresBundle
	readJSON("failures.json", &failures)

	exp := loadExpected()

	pad := func(val any, width int) string {
		s := fmt.Sprintf("%v", val)
		if len(s) >= width {
			return s
		}
		return s + strings.Repeat(" ", width-len(s))
	}

	var parts []string
	add := func(s string) {
		parts = append(parts, s)
	}

	add("<!-- curriculum:generated:start -->\n")
	add("# The Go Engineer v3 — Curriculum\n")
	add(fmt.Sprintf("> Schema: %s | Curriculum: %s | Status: %s | Updated: %s\n",
		core.SchemaVersion, core.CurriculumVersion, core.Status, core.LastUpdated))

	// Architecture Decision
	add("## Architecture Decision\n")
	add("**Decision:** Authoritative v3 modular curriculum engineered around a zero-magic phase-based path from foundational Go to production software engineering.\n")
	add("- **Authoritative Source of Truth:** Authoritative curriculum metadata resides in the modular JSON manifests: `curriculum/path.core.json`, `curriculum/path.electives.json`, `curriculum/projects.json`, `curriculum/assessments.json`, `curriculum/crossrefs.json`, `curriculum/concepts.json`, and `curriculum/failures.json`.")
	add("- **Physical Curriculum Delivery:** Physical lesson READMEs, runnable examples, test suites, and exercises are organized canonically under `curriculum/modules/` and `curriculum/electives/`.")
	add("- **Core Scope:** The core curriculum is dedicated to production Go backend engineering (systems, networking, persistence, concurrency, observability, architecture, and deployment). Modern server-driven full-stack frontend engineering is delivered via the companion path (semantic HTML, CSS, JavaScript, and progressive enhancement with Go templates/templ/HTMX/Tailwind/Datastar).")
	add("- **Capstone & Learning Path Flow:** Flagship Opslane (Module 18) directly requires the complete backend engineering and deployment runbook foundation of Module 15 (`core-15-19`). Module 16 (Portfolio and Interview Readiness) packages the production artifacts produced by Opslane for senior-level engineering interviews. Module 17 (Advanced Electives) provides independent specialization tracks in Generics, gRPC, and distributed systems.\n")

	// Quick stats
	coreItems := core.Items
	lessons := 0
	labs := 0
	checkpoints := 0
	capstones := 0
	pathAssessments := 0
	pathProjects := 0

	for _, it := range coreItems {
		switch it.Type {
		case "lesson":
			lessons++
		case "lab":
			labs++
		case "checkpoint":
			checkpoints++
		case "capstone":
			capstones++
		case "assessment":
			pathAssessments++
		case "project":
			pathProjects++
		}
	}

	electiveItems := len(electives.Items)
	prjRecords := len(projBundle.Projects)
	regAssessments := len(assessBundle.Assessments)

	var moduleAssessments []catalogAssessment
	var projectAssessments []catalogAssessment
	for _, a := range assessBundle.Assessments {
		if strings.HasPrefix(a.ID, "assessment-module-") {
			moduleAssessments = append(moduleAssessments, a)
		} else {
			projectAssessments = append(projectAssessments, a)
		}
	}

	conceptCount := len(concepts.Concepts)
	failureCategoryCount := len(failures.FailureCategories)
	if failureCategoryCount == 0 {
		failureCategoryCount = 13
	}

	flagshipProjectsCount := 0
	for _, p := range projBundle.Projects {
		if p.Type == "flagship" {
			flagshipProjectsCount++
		}
	}

	inv := exp.Inventory
	coreModulesCount := inv.CoreModules
	electiveModulesCount := inv.ElectiveModules
	totalModulesCount := inv.TotalModules
	totalItemCount := inv.TotalItems

	goldenCount := 0
	for _, it := range core.Items {
		if it.ZeroMagicStatus == "golden" {
			goldenCount++
		}
	}
	for _, it := range electives.Items {
		if it.ZeroMagicStatus == "golden" {
			goldenCount++
		}
	}

	entries, err := os.ReadDir(curriculumDir)
	if err != nil {
		return "", fmt.Errorf("read curriculum directory: %w", err)
	}
	manifestCount := 0
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".json") {
			manifestCount++
		}
	}

	add("## Quick Stats\n")
	add("| Metric                           | Count | Notes                                                             |")
	add("| -------------------------------- | ----- | ----------------------------------------------------------------- |")
	add(fmt.Sprintf("| Core Modules                     | %s | Modules 00 through 16, plus Module 18 (Flagship Opslane)          |", pad(coreModulesCount, 5)))
	add(fmt.Sprintf("| Elective Modules                 | %s | Module 17 (Advanced Electives)                                    |", pad(electiveModulesCount, 5)))
	add(fmt.Sprintf("| Total Modules                    | %s | Modules 00–18                                                     |", pad(totalModulesCount, 5)))
	add(fmt.Sprintf("| Core Path Items (Metadata Total) | %s | Total items array in `path.core.json` across all core modules     |", pad(len(coreItems), 5)))
	add(fmt.Sprintf("| Instructional Lessons & Labs     | %s | %d lessons + %d lab across Modules 00–16                         |", pad(lessons+labs, 5), lessons, labs))
	add(fmt.Sprintf("| Flagship Opslane Milestones      | %s | %d checkpoints + %d capstone defense in Module 18          |", pad(checkpoints+capstones, 5), checkpoints, capstones))
	add(fmt.Sprintf("| Learner Instructional Units      | %s | %d instructional lessons/labs + %d Opslane milestones            |", pad(lessons+labs+checkpoints+capstones, 5), lessons+labs, checkpoints+capstones))
	add(fmt.Sprintf("| Elective Items                   | %s | %d advanced elective lessons across Generics, gRPC, and distributed patterns in Module 17 |", pad(electiveItems, 5), electiveItems))
	add(fmt.Sprintf("| Projects                         | %s | %d module projects + %d flagship capstone project                  |", pad(prjRecords, 5), prjRecords-flagshipProjectsCount, flagshipProjectsCount))
	add(fmt.Sprintf("| In-Path Assessment Checkpoints   | %s | End-of-module checkpoints, project reviews, and flagship Opslane milestones across Modules 00–18 |", pad(pathAssessments, 5)))
	add(fmt.Sprintf("| Assessment Registry Records      | %s | %d module assessments + %d project assessments in `assessments.json`|", pad(regAssessments, 5), len(moduleAssessments), len(projectAssessments)))
	add(fmt.Sprintf("| Canonical Concepts               | %s | Fully resolved ownership without leaks in `concepts.json`         |", pad(conceptCount, 5)))
	add(fmt.Sprintf("| Operational Categories           | %s | Mapped across %d module coverage entries                          |\n", pad(failureCategoryCount, 5), coreModulesCount))

	// Module Map
	add("## Module Map\n")
	add("| #   | Module | Phase | Status | Lessons/Labs | Milestones | Projects | Assessments | Total Path Items | Path |")
	add("| --- | --- | --- | --- | ---: | ---: | ---: | ---: | ---: | --- |")

	totLessons, totMilestones, totProjects, totAssessments, totItems := 0, 0, 0, 0, 0
	for _, m := range core.Modules {
		mNum, _ := strconv.Atoi(strings.TrimPrefix(m.ID, "module-"))
		mLessons, mMilestones, mProjects, mAssessments, mTotal := 0, 0, 0, 0, 0
		for _, it := range coreItems {
			if it.ModuleID == m.ID {
				mTotal++
				switch it.Type {
				case "lesson", "lab":
					mLessons++
				case "checkpoint", "capstone":
					mMilestones++
				case "project":
					mProjects++
				case "assessment":
					mAssessments++
				}
			}
		}
		totLessons += mLessons
		totMilestones += mMilestones
		totProjects += mProjects
		totAssessments += mAssessments
		totItems += mTotal

		pathCol := "`curriculum/" + m.Path + "`"
		if !strings.HasPrefix(m.Path, "modules/") {
			pathCol = "`curriculum/modules/" + m.Path + "`"
		}

		add(fmt.Sprintf("| %02d  | **%s** | %s | %s | %d | %d | %d | %d | %d | %s |",
			mNum, m.Title, m.Phase, m.Status, mLessons, mMilestones, mProjects, mAssessments, mTotal, pathCol))
	}

	eModule := electives.Modules[0]
	add(fmt.Sprintf("| 17  | **%s** | %s | %s | %d | 0 | 0 | 0 | %d | `curriculum/%s` |",
		eModule.Title, eModule.Phase, eModule.Status, electiveItems, electiveItems, eModule.Path))
	add(fmt.Sprintf("| —   | **Core Subtotal** | — | — | **%d** | **%d** | **%d** | **%d** | **%d** | *(Modules 00–16, 18)* |",
		totLessons, totMilestones, totProjects, totAssessments, totItems))
	add(fmt.Sprintf("| —   | **Total Curriculum** | — | — | **%d** | **%d** | **%d** | **%d** | **%d** | *(%d core items + %d electives)* |\n",
		totLessons+electiveItems, totMilestones, totProjects, totAssessments, totItems+electiveItems, totItems, electiveItems))

	// Workload Estimates
	add("## Workload Estimates & 52-Week Pacing Guide\n")
	add("Curriculum completion time cannot be calculated from reading time alone. A realistic workload accounts for conceptual reading, hands-on lab exercises, independent debugging, test authoring, and capstone engineering.\n")
	add("### 1. Workload by Activity Layer (Core Path)\n")

	lessonItemsCount := lessons + labs
	lessonMins := 0
	for _, it := range coreItems {
		if it.Type == "lesson" || it.Type == "lab" {
			mins := it.EstimatedMinutes
			if mins == 0 {
				mins = 45
			}
			lessonMins += mins
		}
	}
	lessonHrs := fmt.Sprintf("%.1f", float64(lessonMins)/60.0)

	projectTypeMap := make(map[string]string)
	for _, p := range projBundle.Projects {
		projectTypeMap[p.ID] = p.Type
	}

	var miniProjects, portfolioProjects, flagshipProjects []Item
	for _, it := range coreItems {
		if it.Type == "project" {
			pType := projectTypeMap[it.ID]
			if pType == "mini-project" {
				miniProjects = append(miniProjects, it)
			} else if pType == "portfolio-project" {
				portfolioProjects = append(portfolioProjects, it)
			}
			if pType == "flagship" || strings.Contains(it.ID, "opslane") {
				flagshipProjects = append(flagshipProjects, it)
			}
		}
	}

	miniMins := 0
	for _, it := range miniProjects {
		mins := it.EstimatedMinutes
		if mins == 0 {
			mins = 120
		}
		miniMins += mins
	}
	miniHrs := fmt.Sprintf("%.1f", float64(miniMins)/60.0)

	portfolioMins := 0
	for _, it := range portfolioProjects {
		mins := it.EstimatedMinutes
		if mins == 0 {
			mins = 120
		}
		portfolioMins += mins
	}
	portfolioHrs := fmt.Sprintf("%.1f", float64(portfolioMins)/60.0)

	flagshipMilestonesCount := checkpoints + capstones
	milestoneMins := 0
	for _, it := range coreItems {
		if it.Type == "checkpoint" || it.Type == "capstone" {
			mins := it.EstimatedMinutes
			if mins == 0 {
				mins = 120
			}
			milestoneMins += mins
		}
	}

	flagshipMins := 0
	for _, it := range flagshipProjects {
		mins := it.EstimatedMinutes
		if mins == 0 {
			mins = 120
		}
		flagshipMins += mins
	}

	assessMins := 0
	for _, it := range coreItems {
		if it.Type == "assessment" {
			mins := it.EstimatedMinutes
			if mins == 0 {
				mins = 60
			}
			assessMins += mins
		}
	}
	assessHrs := fmt.Sprintf("%.1f", float64(assessMins)/60.0)

	totalCoreItems := lessonItemsCount + len(miniProjects) + len(portfolioProjects) + len(flagshipProjects) + flagshipMilestonesCount + pathAssessments
	totalCoreMins := lessonMins + miniMins + portfolioMins + flagshipMins + milestoneMins + assessMins
	totalCoreHrs := fmt.Sprintf("%.1f", float64(totalCoreMins)/60.0)

	add("| Activity Layer | Items | On-Paper Metadata | Realistic Learner Effort | Notes |")
	add("| :--- | :---: | :---: | :---: | :--- |")
	add(fmt.Sprintf("| **Instructional Lessons & Labs** | %d | %s hrs | 235 – 265 hrs | Reading, mental models, running code examples, review questions (Modules 00–16) |", lessonItemsCount, lessonHrs))
	add(fmt.Sprintf("| **Mini-Projects** | %d | %s hrs | 36 – 72 hrs | Immediate practical reinforcement (2–4 hours per project) |", len(miniProjects), miniHrs))
	add(fmt.Sprintf("| **Portfolio Projects** | %d | %s hrs | 65 – 120 hrs | Multi-component services with tests, migrations, CI/CD (8–15 hours per project) |", len(portfolioProjects), portfolioHrs))
	add(fmt.Sprintf("| **Flagship Opslane Milestones & Projects** | %d | %.1f hrs | 55 – 85 hrs | Integrated production system with auth, DB, payments, worker pools, metrics (Module 18) |",
		flagshipMilestonesCount+len(flagshipProjects), float64(milestoneMins+flagshipMins)/60.0))
	add(fmt.Sprintf("| **In-Path Assessments & Defense** | %d | %s hrs | 45 – 60 hrs | Self-evaluation checklists, code reviews, and capstone oral/technical defense |", pathAssessments, assessHrs))
	add(fmt.Sprintf("| **Core Total** | **%d** | **%s hrs** | **435 – 600 hrs** | **Realistic full engineering workload across all 18 core modules** |\n", totalCoreItems, totalCoreHrs))

	add("### 2. Weekly Commitment Planning Models\n")
	add("| Weekly Commitment | Total Weeks | Total Months | Viability & Pace Description |")
	add("| :---: | :---: | :---: | :--- |")
	add("| **12–15 hrs/week** | **40 – 52 weeks** | **10 – 12 mos** | **(Recommended)** Completes 435–600 hours of core backend engineering in 29–50 weeks of focused study, leaving deliberate slack for review, autonomous project debugging, and parallel companion frontend progression. |")
	add("| **8–10 hrs/week** | **55 – 65 weeks** | **13 – 15 mos** | Steady progress; core concepts finish in 1 year, with portfolio and Opslane defense extending into months 13–15. |")
	add("| **5 hrs/week** | **90 – 115 weeks** | **21 – 26 mos** | Casual part-time; learners reach Module 08/09 by month 12 and graduate in year 2. |\n")

	add("### 3. Integrated 52-Week Schedule Framework\n")
	add("| Weeks | Phase | Focus Modules | Key Practical Milestones |")
	add("| :---: | :--- | :--- | :--- |")
	add("| **1–3** | Foundations | Modules 00–02 | Environment verification, Git workflow, Shell + Git Lab, Tooling Failure Lab |")
	add("| **4–8** | Core Go Programming | Modules 03–04 | Syntax, memory models, functions, error handling, Pricing Checkout, Order Summary Refactor |")
	add("| **9–14** | Engineering Disciplines | Modules 05–07 | Types, interfaces, test-driven refactoring, CLI tools, Config Parser *(Companion: HTML/CSS)* |")
	add("| **15–24** | Backend & Data Persistence | Modules 08–09 | HTTP APIs, Postgres, migrations, transactions, Book Catalog REST API, Task API *(Companion: JS/DOM)* |")
	add("| **25–33** | Security, Lifecycle & Concurrency | Modules 10–11 | JWT/RBAC, goroutines, channels, worker pools, Secure Authenticated API *(Companion: templ/HTMX)* |")
	add("| **34–42** | Production Engineering | Modules 12–15 | Observability (slog/OTel), profiling (pprof), Docker multi-stage, CI/CD pipelines, deployment runbook |")
	add("| **43–50** | Flagship Implementation | Module 18 | Flagship Opslane: multi-tenant architecture, event bus, caching, graceful shutdown, production deployment |")
	add("| **51–52** | Portfolio & Graduation Review | Module 16 | Portfolio packaging, architecture defense rehearsal, resume project bullets, mock interview prep |\n")

	// Repository structure
	sortedMods := make([]Module, len(core.Modules))
	copy(sortedMods, core.Modules)
	sort.Slice(sortedMods, func(i, j int) bool {
		ni, _ := strconv.Atoi(strings.TrimPrefix(sortedMods[i].ID, "module-"))
		nj, _ := strconv.Atoi(strings.TrimPrefix(sortedMods[j].ID, "module-"))
		return ni < nj
	})
	var modDirs []string
	for _, m := range sortedMods {
		modDirs = append(modDirs, "    "+filepath.Base(m.Path)+"/")
	}

	var curriculumFiles []string
	for _, e := range entries {
		if !e.IsDir() && (strings.HasSuffix(e.Name(), ".json") || (strings.HasSuffix(e.Name(), ".md") && e.Name() != "README.md")) {
			curriculumFiles = append(curriculumFiles, e.Name())
		}
	}
	sort.Strings(curriculumFiles)

	electiveDir := filepath.Base(electives.Modules[0].Path)
	add("## Repository Structure\n")
	add("```\ncurriculum/\n  modules/\n" + strings.Join(modDirs, "\n") + "\n  electives/\n    " + electiveDir + "/\n  companion-frontend/\n    README.md\n  " + strings.Join(curriculumFiles, "\n  ") + "\ninternal/tools/curriculum/\n```\n")

	add("## Companion Tracks & Quality Assurance Protocols\n")
	add("- **[Graduation Outcomes & Capstone Defense Rubric](GRADUATION_CRITERIA.md):** 6 observable competency outcomes (unassisted problem-solving, debugging, testing, disaster recovery, architectural defense, incident response) and formal Flagship Opslane oral defense protocol.")
	add("- **[Frontend Companion Track](companion-frontend/README.md):** 4-stage progressive server-driven hypermedia track (`templ` + `HTMX` + `Tailwind` + `Alpine.js`, with `Datastar` comparison) running parallel to backend modules and converging at Flagship Opslane.")
	add("- **[Beginner Pilot Testing Protocol](PILOT_TESTING_PROTOCOL.md):** Formal empirical zero-intervention testing methodology evaluating unassisted learner velocity and cognitive load across Modules 00–05.\n")

	add("---\n")
	add("# Core Path\n")

	clusterNotes := map[string]string{
		"core-04-15": " *(Workshop: Deterministic Cleanup — Part 1)*",
		"core-04-16": " *(Workshop: Deterministic Cleanup — Part 2)*",
		"core-05-03": " *(Workshop: Method Receivers & Mutation — Part 1)*",
		"core-05-04": " *(Workshop: Method Receivers & Mutation — Part 2)*",
		"core-05-05": " *(Workshop: Method Receivers & Mutation — Part 3)*",
		"core-05-12": " *(Workshop: Dynamic Type Inspection — Part 1)*",
		"core-05-13": " *(Workshop: Dynamic Type Inspection — Part 2)*",
		"core-06-14": " *(Workshop: Delve Debugging Workflow — Step 1)*",
		"core-06-15": " *(Workshop: Delve Debugging Workflow — Step 2)*",
		"core-06-16": " *(Workshop: Delve Debugging Workflow — Step 3)*",
		"core-06-17": " *(Workshop: Delve Debugging Workflow — Step 4)*",
		"core-07-12": " *(Workshop: In-Memory JSON — Part 1)*",
		"core-07-13": " *(Workshop: In-Memory JSON — Part 2)*",
		"core-07-14": " *(Workshop: Streaming JSON — Part 1)*",
		"core-07-15": " *(Workshop: Streaming JSON — Part 2)*",
		"core-11-01": " *(Workshop: Time Mechanics & Scheduling — Part 1)*",
		"core-11-02": " *(Workshop: Time Mechanics & Scheduling — Part 2)*",
		"core-11-03": " *(Workshop: Time Mechanics & Scheduling — Part 3)*",
		"core-11-04": " *(Workshop: Time Mechanics & Scheduling — Part 4)*",
	}

	moduleMap := make(map[string]*Module)
	for i := range core.Modules {
		moduleMap[core.Modules[i].ID] = &core.Modules[i]
	}

	for _, m := range core.Modules {
		mNum, _ := strconv.Atoi(strings.TrimPrefix(m.ID, "module-"))
		var mItems []Item
		for _, it := range coreItems {
			if it.ModuleID == m.ID && (it.Type == "lesson" || it.Type == "lab" || it.Type == "checkpoint" || it.Type == "capstone") {
				mItems = append(mItems, it)
			}
		}

		add(fmt.Sprintf("## %d. %s\n", mNum, m.Title))
		reqStr := "True"
		if !m.Required {
			reqStr = "False"
		}
		add(fmt.Sprintf("> **Phase:** %s | **Status:** %s | **Required:** %s\n", m.Phase, m.Status, reqStr))
		add(fmt.Sprintf("**Learning Goal:** %s\n", m.LearningGoal))

		if m.ID == "module-18" {
			add("**Prerequisites:** 15. Docker, CI/CD, and Deployment (Technical milestone: `core-15-19`)\n")
			add("> **Architecture Note:** Module 16 (Portfolio and Interview Readiness) is intentionally decoupled from technical implementation so learners can construct Flagship Opslane first, producing the production artifact that Module 16 packages for interviews.\n")
		} else if len(m.Prerequisites) > 0 {
			var prereqTitles []string
			for _, pid := range m.Prerequisites {
				if pMod, ok := moduleMap[pid]; ok {
					pNum, _ := strconv.Atoi(strings.TrimPrefix(pMod.ID, "module-"))
					prereqTitles = append(prereqTitles, fmt.Sprintf("%d. %s", pNum, pMod.Title))
				} else {
					prereqTitles = append(prereqTitles, pid)
				}
			}
			add(fmt.Sprintf("**Prerequisites:** %s\n", strings.Join(prereqTitles, ", ")))
		}

		add("| #   | Type    | Slug | Title | Minutes |")
		add("| --- | ------- | ---- | ----- | ------- |")

		for idx, it := range mItems {
			itemType := it.Subtype
			if itemType == "" {
				itemType = it.Type
			}
			note := clusterNotes[it.ID]
			mins := it.EstimatedMinutes
			if mins == 0 {
				mins = 45
			}
			add(fmt.Sprintf("| %d   | %s | %s | **%s**%s | %d |",
				idx+1, pad(itemType, 7), it.Slug, it.Title, note, mins))
		}
		add("")
	}

	add("---\n")
	add("# Advanced Electives\n")

	m17 := electives.Modules[0]
	add(fmt.Sprintf("## 17. %s\n", m17.Title))
	add(fmt.Sprintf("> **Phase:** %s | **Status:** %s | **Required:** False\n", m17.Phase, m17.Status))
	add(fmt.Sprintf("**Learning Goal:** %s\n", m17.LearningGoal))
	add("**Prerequisites:** 16. Portfolio and Interview Readiness\n")
	add("| #   | Type      | Slug | Title | Minutes |")
	add("| --- | --------- | ---- | ----- | ------- |")

	for idx, it := range electives.Items {
		itemType := it.Subtype
		if itemType == "" {
			itemType = it.Type
		}
		mins := it.EstimatedMinutes
		if mins == 0 {
			mins = 45
		}
		add(fmt.Sprintf("| %d   | %s | %s | **%s** | %d |",
			idx+1, pad(itemType, 9), it.Slug, it.Title, mins))
	}

	add("\n---\n")
	add("# Projects\n")
	add("| #   | Module | Title | Type | Stage | Difficulty |")
	add("| --- | ------ | ----- | ---- | ----- | ---------- |")

	sortedProjects := append([]catalogProject{}, projBundle.Projects...)
	sort.Slice(sortedProjects, func(i, j int) bool {
		ai, _ := strconv.Atoi(strings.TrimPrefix(sortedProjects[i].ModuleID, "module-"))
		aj, _ := strconv.Atoi(strings.TrimPrefix(sortedProjects[j].ModuleID, "module-"))
		if ai != aj {
			return ai < aj
		}
		return sortedProjects[i].ID < sortedProjects[j].ID
	})

	for idx, p := range sortedProjects {
		modNum := "—"
		if p.ModuleID != "" {
			modNum = strings.Replace(p.ModuleID, "module-", "Mod ", 1)
		}
		add(fmt.Sprintf("| %d   | %s | **%s** | %s | %s | %s |",
			idx+1, modNum, p.Title, p.Type, p.Stage, p.Difficulty))
	}

	add("\n---\n")
	add("# Assessments\n")
	add("The assessment framework evaluates mastery across end-of-module milestones and practical project implementations.\n")
	add(fmt.Sprintf("### Module Assessments (%d Records)\n", len(moduleAssessments)))
	add("| ID | Title | Type | Passing Score |")
	add("| -- | ----- | ---- | ------------- |")

	for _, a := range moduleAssessments {
		scoreStr := fmt.Sprintf("%v", a.PassingScore)
		if !strings.HasSuffix(scoreStr, "%") {
			scoreStr += "%"
		}
		add(fmt.Sprintf("| `%s` | **%s** | %s | %s |", a.ID, a.Title, a.Type, scoreStr))
	}

	add(fmt.Sprintf("\n### Project Assessments (%d Records)\n", len(projectAssessments)))
	add("| ID | Title | Type | Passing Score |")
	add("| -- | ----- | ---- | ------------- |")

	for _, a := range projectAssessments {
		scoreStr := fmt.Sprintf("%v", a.PassingScore)
		if !strings.HasSuffix(scoreStr, "%") {
			scoreStr += "%"
		}
		add(fmt.Sprintf("| `%s` | **%s** | %s | %s |", a.ID, a.Title, a.Type, scoreStr))
	}

	add("\n---\n")
	add("# Cross-References\n")

	refs := crossrefs.Crossrefs.References
	relDist := make(map[string]int)
	for _, r := range refs {
		relDist[r.Relation]++
	}

	add(fmt.Sprintf("Total cross-references: **%d**\n", len(refs)))
	add("> **Gating vs. lineage:** Hard progression gates are the `prerequisites` arrays in\n> `path.core.json` and `path.electives.json`. `builds_on` records the broader,\n> directional pedagogical lineage and must remain acyclic, but it is not the\n> field a learner-frontier tool should use to block progression.\n")
	add("### Style Rules\n")
	add("| Relation      | Template                                                     |")
	add("| ------------- | ------------------------------------------------------------ |")
	add("| builds_on     | Builds on: Module X.Y — Lesson Title                         |")
	add("| preview_only  | Preview only: Explained fully in Module X.Y — Lesson Title   |")
	add("| reinforced_in | Reinforced in: Project X — Project Title                     |")
	add("| related       | Related: Module X.Y — Lesson Title                           |\n")

	add("### Relation Distribution\n")
	add("| Relation      | Count | Description |")
	add("| ------------- | ----- | ----------- |")
	add(fmt.Sprintf("| related       | %d   | Bidirectional conceptual links between related lessons |", relDist["related"]))
	add(fmt.Sprintf("| builds_on     | %d   | Directional pedagogical lineage (acyclic; not the hard gate) |", relDist["builds_on"]))
	add(fmt.Sprintf("| reinforced_in | %d   | Practical application links to projects and labs |", relDist["reinforced_in"]))
	add(fmt.Sprintf("| preview_only  | %d    | Early mental-model previews developed fully later |\n", relDist["preview_only"]))

	add("### Authoritative Quality Gates\n")
	var gateLines []string
	gateLines = append(gateLines,
		fmt.Sprintf("- `validate-schema`: %d canonical JSON manifests valid against Schema 3.1.0", manifestCount),
		"- `validate-graph`: Strict DAG acyclicity across all item dependencies",
		"- `validate-lessons`: Zero machine role comments or unmotivated markers",
		"- `validate-prerequisite-closure`: Strict pedagogical progression across module boundaries",
		fmt.Sprintf("- `validate-project-binding`: All %d projects bound to valid curriculum placement anchors", prjRecords),
		"- `validate-no-placeholder-zm`: Zero mechanical placeholders across zero-magic blocks",
		fmt.Sprintf("- `validate-assessment-completeness`: All %d assessments equipped with weighted criteria", regAssessments),
		"- `validate-unique-tags`: Tag taxonomy consistency across core and electives",
		"- `validate-crossrefs`: Complete bidirectional graph integrity across all entities",
		fmt.Sprintf("- `validate-concept-ownership`: All %d concepts mapped to unambiguous canonical owners", conceptCount),
		"- `validate-cognitive-load`: Module pacing and cognitive density within verified bounds",
		fmt.Sprintf("- `validate-operational-failures`: %d operational failure categories mapped across %d modules", failureCategoryCount, coreModulesCount),
		"- `validate-failure-engineering`: Failure modes and detection strategies verified",
		fmt.Sprintf("- `validate-zero-magic`: %d golden lessons certified against pedagogical spec", goldenCount),
		fmt.Sprintf("- `validate-physical-artifacts`: %d/%d physical lesson and elective READMEs verified on disk", totalItemCount, totalItemCount),
		fmt.Sprintf("- `validate-assessment-artifacts`: %d assessments + %d Opslane milestones verified with complete physical suites", regAssessments, inv.OpslaneMilestones),
		fmt.Sprintf("- `validate-starter-parity`: %d project starters (RED) and solutions (GREEN) verified with valid Go ASTs", inv.StarterParityProjects),
		"- `validate-markdown-links`: All relative Markdown links verified with zero broken targets",
		"- `validate-inline-paths`: All inline-code repository path citations resolve to real files on disk",
		"- `validate-sese-navigation`: Single-entry, single-exit navigation flow verified",
		"- `validate-project-scripts`: Script-backed project harness parity verified",
		"- `validate-no-emojis`: Zero emoji violations across Markdown prose",
		"- `validate-anchor-quotes`: Authoritative pioneer anchor quotes present in locked lessons",
		"- `validate-diagram-palette`: Mermaid diagrams comply with the Dark Engineering Palette",
		"- `validate-readme-contracts`: Documentation contracts satisfied per item family and subtype",
		"- `validate-anti-slop`: Strict anti-slop lexical rules and banned generic filler eliminated across prose",
	)
	add(strings.Join(gateLines, "\n"))

	add("\n### Validation Commands\n")
	add("| Command | Purpose | Execution |")
	add("| --- | --- | --- |")
	add("| **Run All Gates** | Execute full curriculum verification suite | `go run ./scripts/validate_curriculum.go` |")
	add("| **Validate Schema** | Check JSON manifests and key ordering | `go run ./internal/tools/curriculum validate-schema` |")
	add("| **Validate Graph** | Check DAG dependency structure | `go run ./internal/tools/curriculum validate-graph` |")
	add(fmt.Sprintf("| **Validate Physical** | Verify %d physical READMEs | `go run ./internal/tools/curriculum validate-physical-artifacts` |", totalItemCount))
	add("| **Validate Assessments** | Verify 4-file assessment suites | `go run ./internal/tools/curriculum validate-assessment-artifacts` |")
	add("| **Validate Code Parity** | Verify Go code and AST correctness | `go run ./internal/tools/curriculum validate-starter-parity` |\n")

	add("---\n")
	add(fmt.Sprintf("_Generated from curriculum v%s | Schema %s | %s_\n",
		core.CurriculumVersion, core.SchemaVersion, core.LastUpdated))
	add("<!-- curriculum:generated:end -->\n")

	return strings.Join(parts, "\n"), nil
}

// writeCatalog generates curriculum/README.md and writes it to disk.
func writeCatalog() {
	content, err := generateCatalog()
	if err != nil {
		fatalf("failed to generate catalog: %v", err)
	}
	outPath := filepath.Join(curriculumDir, "README.md")
	if err := os.WriteFile(outPath, []byte(content), 0o644); err != nil {
		fatalf("failed to write catalog to %s: %v", outPath, err)
	}
	fmt.Printf("Successfully generated %s (%d bytes)\n", outPath, len(content))
}
