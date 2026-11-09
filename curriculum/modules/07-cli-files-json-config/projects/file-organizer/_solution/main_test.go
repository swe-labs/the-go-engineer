package main

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func createTestFile(t *testing.T, path, content string) {
	t.Helper()
	err := os.WriteFile(path, []byte(content), 0644)
	if err != nil {
		t.Fatalf("failed to create test file %s: %v", path, err)
	}
}

func TestOrganize_ExtensionMode(t *testing.T) {
	tempDir := t.TempDir()
	srcDir := filepath.Join(tempDir, "source")
	targetDir := filepath.Join(tempDir, "target")
	_ = os.MkdirAll(srcDir, 0755)

	createTestFile(t, filepath.Join(srcDir, "doc.pdf"), "pdf content")
	createTestFile(t, filepath.Join(srcDir, "avatar.png"), "image content")
	createTestFile(t, filepath.Join(srcDir, "song.mp3"), "audio content")

	organizer := NewFileOrganizer()
	report, err := organizer.Organize(t.Context(), OrganizerConfig{
		SourceDir: srcDir,
		TargetDir: targetDir,
		Mode:      ModeExtension,
		Policy:    PolicyRename,
	})

	if err != nil {
		t.Fatalf("organize failed: %v", err)
	}
	if report.FilesMoved != 3 {
		t.Fatalf("expected 3 files moved, got %d", report.FilesMoved)
	}

	// Verify locations
	if _, err := os.Stat(filepath.Join(targetDir, "documents", "doc.pdf")); err != nil {
		t.Fatal("expected doc.pdf in documents/")
	}
	if _, err := os.Stat(filepath.Join(targetDir, "images", "avatar.png")); err != nil {
		t.Fatal("expected avatar.png in images/")
	}
	if _, err := os.Stat(filepath.Join(targetDir, "audio", "song.mp3")); err != nil {
		t.Fatal("expected song.mp3 in audio/")
	}
}

func TestOrganize_CollisionRename(t *testing.T) {
	tempDir := t.TempDir()
	srcDir := filepath.Join(tempDir, "source")
	targetDir := filepath.Join(tempDir, "target")
	_ = os.MkdirAll(srcDir, 0755)
	_ = os.MkdirAll(filepath.Join(targetDir, "documents"), 0755)

	// Existing file in target
	createTestFile(t, filepath.Join(targetDir, "documents", "doc.txt"), "original")
	// File in source with same name
	createTestFile(t, filepath.Join(srcDir, "doc.txt"), "new file")

	organizer := NewFileOrganizer()
	report, err := organizer.Organize(t.Context(), OrganizerConfig{
		SourceDir: srcDir,
		TargetDir: targetDir,
		Mode:      ModeExtension,
		Policy:    PolicyRename,
	})

	if err != nil {
		t.Fatalf("organize failed: %v", err)
	}
	if report.FilesMoved != 1 {
		t.Fatalf("expected 1 file moved, got %d", report.FilesMoved)
	}

	// Original still exists
	origContent, _ := os.ReadFile(filepath.Join(targetDir, "documents", "doc.txt"))
	if string(origContent) != "original" {
		t.Fatalf("original file was corrupted: %s", string(origContent))
	}

	// Renamed copy exists
	renamedContent, err := os.ReadFile(filepath.Join(targetDir, "documents", "doc_1.txt"))
	if err != nil || string(renamedContent) != "new file" {
		t.Fatalf("expected doc_1.txt with 'new file', got %s (err: %v)", string(renamedContent), err)
	}
}

func TestOrganize_CollisionSkip(t *testing.T) {
	tempDir := t.TempDir()
	srcDir := filepath.Join(tempDir, "source")
	targetDir := filepath.Join(tempDir, "target")
	_ = os.MkdirAll(srcDir, 0755)
	_ = os.MkdirAll(filepath.Join(targetDir, "documents"), 0755)

	createTestFile(t, filepath.Join(targetDir, "documents", "data.txt"), "stay")
	createTestFile(t, filepath.Join(srcDir, "data.txt"), "skip me")

	organizer := NewFileOrganizer()
	report, err := organizer.Organize(t.Context(), OrganizerConfig{
		SourceDir: srcDir,
		TargetDir: targetDir,
		Mode:      ModeExtension,
		Policy:    PolicySkip,
	})

	if err != nil {
		t.Fatalf("organize failed: %v", err)
	}
	if report.FilesSkipped != 1 {
		t.Fatalf("expected 1 file skipped, got %d", report.FilesSkipped)
	}
	if report.FilesMoved != 0 {
		t.Fatalf("expected 0 files moved, got %d", report.FilesMoved)
	}
}

func TestOrganize_DryRun(t *testing.T) {
	tempDir := t.TempDir()
	srcDir := filepath.Join(tempDir, "source")
	targetDir := filepath.Join(tempDir, "target")
	_ = os.MkdirAll(srcDir, 0755)

	createTestFile(t, filepath.Join(srcDir, "log.txt"), "log data")

	organizer := NewFileOrganizer()
	report, err := organizer.Organize(t.Context(), OrganizerConfig{
		SourceDir: srcDir,
		TargetDir: targetDir,
		Mode:      ModeExtension,
		DryRun:    true,
	})

	if err != nil {
		t.Fatalf("dry-run failed: %v", err)
	}
	if report.FilesMoved != 1 {
		t.Fatalf("expected 1 simulated move, got %d", report.FilesMoved)
	}

	// File must still be in sourceDir!
	if _, err := os.Stat(filepath.Join(srcDir, "log.txt")); err != nil {
		t.Fatal("file was moved during dry run")
	}
}

func TestOrganize_InvalidDirectory(t *testing.T) {
	organizer := NewFileOrganizer()
	_, err := organizer.Organize(context.Background(), OrganizerConfig{
		SourceDir: "/non/existent/path/987123",
	})

	if !errors.Is(err, ErrInvalidSource) {
		t.Fatalf("expected ErrInvalidSource, got %v", err)
	}
}
