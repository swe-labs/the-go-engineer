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
	ext = strings.ToLower(strings.TrimPrefix(ext, "."))
	switch ext {
	case "jpg", "jpeg", "png", "gif", "webp", "svg":
		return "images"
	case "pdf", "doc", "docx", "txt", "md":
		return "documents"
	case "mp3", "wav", "flac", "aac":
		return "audio"
	case "mp4", "mkv", "mov", "webm":
		return "videos"
	case "zip", "tar", "gz", "7z":
		return "archives"
	case "":
		return "unknown"
	default:
		return ext
	}
}

// Organize scans SourceDir and sorts files into subdirectories.
func (o *FileOrganizer) Organize(ctx context.Context, cfg OrganizerConfig) (OrganizeReport, error) {
	if strings.TrimSpace(cfg.SourceDir) == "" {
		return OrganizeReport{}, ErrEmptySource
	}

	srcInfo, err := os.Stat(cfg.SourceDir)
	if err != nil || !srcInfo.IsDir() {
		return OrganizeReport{}, ErrInvalidSource
	}

	if cfg.TargetDir == "" {
		cfg.TargetDir = cfg.SourceDir
	}
	if cfg.Mode == "" {
		cfg.Mode = ModeExtension
	}
	if cfg.Policy == "" {
		cfg.Policy = PolicyRename
	}

	var report OrganizeReport

	err = filepath.WalkDir(cfg.SourceDir, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			report.Errors = append(report.Errors, fmt.Sprintf("access error on %s: %v", path, walkErr))
			return nil
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		// Skip directories and the root itself
		if d.IsDir() {
			if path != cfg.SourceDir && filepath.Dir(path) != cfg.SourceDir {
				// Don't descend deeply into already organized subdirectories
				return filepath.SkipDir
			}
			return nil
		}

		// Skip hidden files
		if strings.HasPrefix(d.Name(), ".") {
			return nil
		}

		report.FilesScanned++

		fileInfo, err := d.Info()
		if err != nil {
			report.Errors = append(report.Errors, fmt.Sprintf("info error on %s: %v", path, err))
			return nil
		}

		// Determine destination subfolder
		var subfolder string
		if cfg.Mode == ModeDate {
			subfolder = fileInfo.ModTime().UTC().Format("2006-01")
		} else {
			subfolder = CategoryFromExtension(filepath.Ext(d.Name()))
		}

		destDir := filepath.Join(cfg.TargetDir, subfolder)
		destPath := filepath.Join(destDir, d.Name())

		// If file is already in its target folder, skip
		if filepath.Dir(path) == destDir {
			report.FilesSkipped++
			return nil
		}

		// Handle destination collision
		if _, err := os.Stat(destPath); err == nil {
			switch cfg.Policy {
			case PolicySkip:
				report.FilesSkipped++
				return nil
			case PolicyRename:
				destPath = resolveCollisionPath(destDir, d.Name())
			case PolicyOverwrite:
				// Proceed with overwrite
			}
		}

		if cfg.DryRun {
			report.FilesMoved++
			return nil
		}

		if err := os.MkdirAll(destDir, 0755); err != nil {
			report.Errors = append(report.Errors, fmt.Sprintf("mkdir error on %s: %v", destDir, err))
			return nil
		}

		// Move file safely (fallback to copy+delete for cross-device moves)
		if err := moveFile(path, destPath); err != nil {
			report.Errors = append(report.Errors, fmt.Sprintf("move error from %s to %s: %v", path, destPath, err))
			return nil
		}

		report.FilesMoved++
		return nil
	})

	return report, err
}

func resolveCollisionPath(dir, filename string) string {
	ext := filepath.Ext(filename)
	base := strings.TrimSuffix(filename, ext)

	counter := 1
	for {
		candidate := filepath.Join(dir, fmt.Sprintf("%s_%d%s", base, counter, ext))
		if _, err := os.Stat(candidate); errors.Is(err, os.ErrNotExist) {
			return candidate
		}
		counter++
	}
}

func moveFile(src, dst string) error {
	err := os.Rename(src, dst)
	if err == nil {
		return nil
	}

	// Fallback for cross-device rename
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		return err
	}

	_ = in.Close()
	return os.Remove(src)
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
	if len(report.Errors) > 0 {
		fmt.Printf("Errors encountered (%d):\n", len(report.Errors))
		for _, e := range report.Errors {
			fmt.Printf(" - %s\n", e)
		}
	}
}
