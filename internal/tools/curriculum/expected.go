package main

// ExpectedBundle mirrors curriculum/expected.json — the single source of truth
// for curriculum inventory counts, the canonical last_updated date, and the
// locked-module set. Validators read these values instead of hardcoding magic
// numbers in Go, so a count change (e.g. 20 -> 26 electives) is edited in one
// obvious place.
type ExpectedInventory struct {
	CoreItems             int `json:"core_items"`
	CoreModules           int `json:"core_modules"`
	ElectiveItems         int `json:"elective_items"`
	ElectiveModules       int `json:"elective_modules"`
	TotalItems            int `json:"total_items"`
	TotalModules          int `json:"total_modules"`
	Assessments           int `json:"assessments"`
	Projects              int `json:"projects"`
	Concepts              int `json:"concepts"`
	StarterParityProjects int `json:"starter_parity_projects"`
	OpslaneMilestones     int `json:"opslane_milestones"`
}

type ExpectedBundle struct {
	SchemaVersion     string            `json:"schema_version"`
	CurriculumVersion string            `json:"curriculum_version"`
	LastUpdated       string            `json:"last_updated"`
	Status            string            `json:"status"`
	Inventory         ExpectedInventory `json:"inventory"`
	LockedModules     []string          `json:"locked_modules"`
}

var expectedCache *ExpectedBundle

// loadExpected returns the parsed expected.json manifest, cached after first load.
func loadExpected() *ExpectedBundle {
	if expectedCache != nil {
		return expectedCache
	}
	var e ExpectedBundle
	readJSON("expected.json", &e)
	expectedCache = &e
	return expectedCache
}

// lockedModuleSet returns the authoritative set of locked module IDs derived
// from expected.json. This is the single source of truth that replaces the
// hardcoded maps previously duplicated across contracts.go, sese.go, quotes.go,
// and palette.go.
func lockedModuleSet() map[string]bool {
	set := make(map[string]bool)
	for _, id := range loadExpected().LockedModules {
		set[id] = true
	}
	return set
}
