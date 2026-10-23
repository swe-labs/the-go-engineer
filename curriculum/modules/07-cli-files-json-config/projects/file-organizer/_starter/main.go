package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	ErrInvalidSource = errors.New("source directory does not exist or is not a directory")
	ErrEmptySource   = errors.New("source directory path cannot be empty")
)

type OrganizeMode string

const (
	ModeExtension OrganizeMode = "extension"
	ModeDate      OrganizeMode = "date"
)

type CollisionPolicy string

const (
	PolicyRename    CollisionPolicy = "rename"
	PolicySkip      CollisionPolicy = "skip"
	PolicyOverwrite CollisionPolicy = "overwrite"
)

type OrganizerConfig struct {
	SourceDir string
	TargetDir string
	Mode      OrganizeMode
	Policy    CollisionPolicy
	DryRun    bool
}

type OrganizeReport struct {
	FilesScanned int
	FilesMoved   int
	FilesSkipped int
	Errors       []string
}

type FileOrganizer struct{}

func NewFileOrganizer() *FileOrganizer {
	return &FileOrganizer{}
}

// CategoryFromExtension maps file extensions to semantic directory categories.
func CategoryFromExtension(ext string) string {
	// TODO: Map image, document, audio, video, archive extensions to semantic folder names
	return "other"
}

// Organize scans SourceDir and sorts files into subdirectories.
func (o *FileOrganizer) Organize(ctx context.Context, cfg OrganizerConfig) (OrganizeReport, error) {
	// TODO: Validate SourceDir exists and is a directory
	// TODO: Walk SourceDir using filepath.WalkDir
	// TODO: For each file, determine target subfolder based on cfg.Mode (extension vs date)
	// TODO: Handle destination collisions based on cfg.Policy (rename, skip, overwrite)
	// TODO: If not DryRun, move file using moveFile
	return OrganizeReport{}, errors.New("not implemented")
}

func resolveCollisionPath(dir, filename string) string {
	// TODO: Incrementally append _1, _2 before file extension until unused path found
	return filepath.Join(dir, filename)
}

func moveFile(src, dst string) error {
	// TODO: Attempt os.Rename, fallback to io.Copy + os.Remove for cross-device support
	return errors.New("not implemented")
}

func main() {
	srcDir := flag.String("src", ".", "Source directory to scan")
	targetDir := flag.String("target", "", "Target directory for organized files")
	mode := flag.String("mode", "extension", "Organization mode: extension or date")
	policy := flag.String("policy", "rename", "Collision policy: rename, skip, or overwrite")
	dryRun := flag.Bool("dry-run", false, "Simulate file moves without modifying filesystem")
	flag.Parse()

	organizer := NewFileOrganizer()
	cfg := OrganizerConfig{
		SourceDir: *srcDir,
		TargetDir: *targetDir,
		Mode:      OrganizeMode(*mode),
		Policy:    CollisionPolicy(*policy),
		DryRun:    *dryRun,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	slog.Info("organizing files", "src", cfg.SourceDir, "mode", cfg.Mode, "dry_run", cfg.DryRun)
	report, err := organizer.Organize(ctx, cfg)
	if err != nil {
		slog.Error("organization failed", "error", err)
		os.Exit(1)
	}

	fmt.Printf("Files Scanned: %d\n", report.FilesScanned)
	fmt.Printf("Files Moved:   %d\n", report.FilesMoved)
	fmt.Printf("Files Skipped: %d\n", report.FilesSkipped)
}
