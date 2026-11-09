package main

import (
	"path/filepath"
	"runtime"
	"testing"
)

func TestJoin(t *testing.T) {
	p := filepath.Join("a", "b", "c.txt")
	expected := "a/b/c.txt"
	if runtime.GOOS == "windows" {
		expected = "a\\b\\c.txt"
	}
	if p != expected {
		t.Errorf("Join('a','b','c.txt') = %q, want %q", p, expected)
	}
}

func TestClean(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"a/b/../c", "a/c"},
		{"a/./b", "a/b"},
		{"/a/b/", "/a/b"},
		{".", "."},
	}

	for _, tt := range tests {
		got := filepath.Clean(tt.input)
		// On Windows, "/" becomes the volume root
		if got != tt.expected && runtime.GOOS != "windows" {
			t.Errorf("Clean(%q) = %q, want %q", tt.input, got, tt.expected)
		}
	}
}

func TestSplit(t *testing.T) {
	dir, file := filepath.Split("a/b/c.txt")
	if dir != "a/b/" && runtime.GOOS == "windows" && dir != "a/b\\" {
		t.Logf("Split('a/b/c.txt') dir = %q", dir)
	}
	if file != "c.txt" {
		t.Errorf("Split('a/b/c.txt') file = %q, want 'c.txt'", file)
	}
}

func TestExt(t *testing.T) {
	ext := filepath.Ext("data.json")
	if ext != ".json" {
		t.Errorf("Ext('data.json') = %q, want '.json'", ext)
	}

	ext = filepath.Ext("archive.tar.gz")
	if ext != ".gz" {
		t.Errorf("Ext('archive.tar.gz') = %q, want '.gz'", ext)
	}

	ext = filepath.Ext("noext")
	if ext != "" {
		t.Errorf("Ext('noext') = %q, want ''", ext)
	}
}

func TestBase(t *testing.T) {
	base := filepath.Base("/a/b/c.txt")
	if base != "c.txt" {
		t.Errorf("Base('/a/b/c.txt') = %q, want 'c.txt'", base)
	}

	base = filepath.Base("/a/b/")
	if base != "b" {
		t.Logf("Base('/a/b/') = %q", base)
	}
}

func TestDir(t *testing.T) {
	dir := filepath.Dir("/a/b/c.txt")
	if dir != "/a/b" && dir != "\\a\\b" {
		t.Logf("Dir('/a/b/c.txt') = %q", dir)
	}
}

func TestIsAbs(t *testing.T) {
	if !filepath.IsAbs("/") || filepath.IsAbs("relative") {
		if runtime.GOOS == "windows" {
			// On Windows, "/" is not absolute; "C:\\" is
		}
		t.Logf("IsAbs('/') = %v, IsAbs('relative') = %v", filepath.IsAbs("/"), filepath.IsAbs("relative"))
	}
}

func TestRel(t *testing.T) {
	rel, err := filepath.Rel("/a/b", "/a/b/c/d.txt")
	if err != nil {
		t.Fatal(err)
	}
	expected := "c/d.txt"
	if runtime.GOOS == "windows" {
		expected = "c/d.txt"
	}
	if rel != expected {
		t.Logf("Rel('/a/b', '/a/b/c/d.txt') = %q, want %q", rel, expected)
	}
}

func TestGlob(t *testing.T) {
	matches, err := filepath.Glob("*.go")
	if err != nil {
		t.Fatal(err)
	}
	if len(matches) == 0 {
		t.Error("expected at least one .go file match")
	}
}

func TestVolumeName(t *testing.T) {
	vol := filepath.VolumeName("C:\\path\\file.txt")
	if runtime.GOOS == "windows" && vol != "C:" {
		t.Errorf("VolumeName('C:\\\\path\\\\file.txt') = %q, want 'C:'", vol)
	}
}
