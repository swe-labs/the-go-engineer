package main

import "testing"

func TestValidatePath_ValidSequential(t *testing.T) {
	modules := map[string]Module{
		"00": {ID: "00", Title: "Orientation"},
		"01": {ID: "01", Title: "Basics", Prerequisites: []string{"00"}},
		"02": {ID: "02", Title: "Advanced", Prerequisites: []string{"00", "01"}},
	}
	path := []Module{modules["00"], modules["01"], modules["02"]}
	result := validatePath(path)
	if !result.Valid {
		t.Errorf("expected valid path, got error: %s", result.Error)
	}
}

func TestValidatePath_MissingPrerequisite(t *testing.T) {
	modules := map[string]Module{
		"00": {ID: "00", Title: "Orientation"},
		"01": {ID: "01", Title: "Basics", Prerequisites: []string{"00"}},
		"02": {ID: "02", Title: "Advanced", Prerequisites: []string{"00", "01"}},
	}
	path := []Module{modules["02"]}
	result := validatePath(path)
	if result.Valid {
		t.Error("expected invalid path (missing prerequisites for module 02)")
	}
}

func TestValidatePath_SingleModuleNoPrereqs(t *testing.T) {
	m := Module{ID: "00", Title: "Orientation"}
	result := validatePath([]Module{m})
	if !result.Valid {
		t.Errorf("expected single module with no prereqs to be valid, got: %s", result.Error)
	}
}

func TestValidatePath_WrongOrder(t *testing.T) {
	modules := map[string]Module{
		"00": {ID: "00", Title: "Orientation"},
		"01": {ID: "01", Title: "Basics", Prerequisites: []string{"00"}},
	}
	path := []Module{modules["01"], modules["00"]}
	result := validatePath(path)
	if result.Valid {
		t.Error("expected invalid path (01 before 00)")
	}
}

func TestValidatePath_MultiplePrereqs(t *testing.T) {
	modules := map[string]Module{
		"00": {ID: "00", Title: "A"},
		"01": {ID: "01", Title: "B"},
		"02": {ID: "02", Title: "C", Prerequisites: []string{"00", "01"}},
	}
	path := []Module{modules["00"], modules["01"], modules["02"]}
	result := validatePath(path)
	if !result.Valid {
		t.Errorf("expected valid path, got: %s", result.Error)
	}
}

func TestValidatePath_EmptyPath(t *testing.T) {
	result := validatePath(nil)
	if !result.Valid {
		t.Error("empty path should be valid")
	}
}
