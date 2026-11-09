package main

import (
	"crypto/sha256"
	"encoding/hex"
	"runtime"
	"testing"
)

func TestNewArtifact(t *testing.T) {
	a := NewArtifact("testapp", "v1.0", "linux", "amd64", []byte("data"))
	if a.Name != "testapp" {
		t.Errorf("unexpected name: %s", a.Name)
	}
	if a.Version != "v1.0" {
		t.Errorf("unexpected version: %s", a.Version)
	}
}

func TestBinaryName(t *testing.T) {
	tests := []struct {
		artifact Artifact
		want     string
	}{
		{Artifact{Name: "app", Version: "v1", GOOS: "linux", GOARCH: "amd64"}, "app-v1-linux-amd64"},
		{Artifact{Name: "app", Version: "v1", GOOS: "windows", GOARCH: "amd64"}, "app-v1-windows-amd64.exe"},
		{Artifact{Name: "", Version: "v1", GOOS: "linux", GOARCH: "arm64"}, "app-v1-linux-arm64"},
	}
	for _, tc := range tests {
		got := tc.artifact.BinaryName()
		if got != tc.want {
			t.Errorf("BinaryName(%+v) = %q, want %q", tc.artifact, got, tc.want)
		}
	}
}

func TestComputeChecksum(t *testing.T) {
	a := NewArtifact("test", "v1", "linux", "amd64", []byte("hello"))
	h := sha256.Sum256([]byte("hello"))
	want := hex.EncodeToString(h[:])
	got := a.ComputeChecksum()
	if got != want {
		t.Errorf("checksum mismatch: got %s, want %s", got, want)
	}
}

func TestChecksumLine(t *testing.T) {
	a := NewArtifact("app", "v1", "linux", "amd64", []byte("test data"))
	a.Checksum = a.ComputeChecksum()
	line := a.ChecksumLine()
	expectedLine := a.Checksum + "  " + a.BinaryName()
	if line != expectedLine {
		t.Errorf("unexpected checksum line: %s", line)
	}
}

func TestParseChecksums(t *testing.T) {
	content := "abc123  myapp-v1-linux-amd64\ndef456  myapp-v1-windows-amd64.exe\n"
	cf, err := ParseChecksums(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cf.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(cf.Entries))
	}
	if cf.Entries[0].Checksum != "abc123" {
		t.Errorf("unexpected first checksum: %s", cf.Entries[0].Checksum)
	}
}

func TestParseChecksumsInvalidLine(t *testing.T) {
	_, err := ParseChecksums("onlyonefield")
	if err == nil {
		t.Error("expected error for invalid line")
	}
}

func TestParseChecksumsEmptyContent(t *testing.T) {
	cf, err := ParseChecksums("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cf.Entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(cf.Entries))
	}
}

func TestVerifyChecksum(t *testing.T) {
	data := []byte("original content")
	a := NewArtifact("app", "v1", "linux", "amd64", data)
	a.Checksum = a.ComputeChecksum()

	cf := &ChecksumFile{
		Entries: []ChecksumEntry{
			{Checksum: a.Checksum, BinaryName: a.BinaryName()},
		},
	}
	if !cf.Verify(a) {
		t.Error("expected verification to pass")
	}
}

func TestVerifyChecksumTampered(t *testing.T) {
	original := NewArtifact("app", "v1", "linux", "amd64", []byte("original"))
	original.Checksum = original.ComputeChecksum()

	cf := &ChecksumFile{
		Entries: []ChecksumEntry{
			{Checksum: original.Checksum, BinaryName: original.BinaryName()},
		},
	}

	tampered := NewArtifact("app", "v1", "linux", "amd64", []byte("tampered"))
	if cf.Verify(tampered) {
		t.Error("expected verification to fail for tampered artifact")
	}
}

func TestVerifyChecksumWrongName(t *testing.T) {
	a := NewArtifact("app", "v1", "linux", "amd64", []byte("data"))
	a.Checksum = a.ComputeChecksum()

	cf := &ChecksumFile{
		Entries: []ChecksumEntry{
			{Checksum: a.Checksum, BinaryName: "different-name"},
		},
	}
	if cf.Verify(a) {
		t.Error("expected verification to fail when binary name differs")
	}
}

func TestVersionLdflags(t *testing.T) {
	if Version != "dev" {
		t.Errorf("expected default version dev, got %s", Version)
	}
	if Commit != "none" {
		t.Errorf("expected default commit none, got %s", Commit)
	}
}

func TestArtifactEmptyNameDefault(t *testing.T) {
	a := NewArtifact("", "v1", runtime.GOOS, runtime.GOARCH, []byte("x"))
	name := a.BinaryName()
	if name == "" {
		t.Error("binary name should not be empty")
	}
}
