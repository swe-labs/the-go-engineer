package main

import "testing"

func TestDemoScope(t *testing.T) {
	result := demoScope()
	expected := "block local global"
	if result != expected {
		t.Errorf("demoScope() = %q; want %q", result, expected)
	}
}

func TestPackageVarUnchanged(t *testing.T) {
	if name != "global" {
		t.Errorf("package-level name = %q; want %q", name, "global")
	}
}
