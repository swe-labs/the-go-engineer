package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
)

var (
	Version   = "dev"
	Commit    = "none"
	BuildTime = "unknown"
)

type Artifact struct {
	Name     string
	Version  string
	GOOS     string
	GOARCH   string
	Checksum string
	Data     []byte
}

func NewArtifact(name, version, goos, goarch string, data []byte) *Artifact {
	return &Artifact{
		Name:    name,
		Version: version,
		GOOS:    goos,
		GOARCH:  goarch,
		Data:    data,
	}
}

func (a *Artifact) BinaryName() string {
	ext := ""
	if a.GOOS == "windows" {
		ext = ".exe"
	}
	if a.Name == "" {
		a.Name = "app"
	}
	return fmt.Sprintf("%s-%s-%s-%s%s", a.Name, a.Version, a.GOOS, a.GOARCH, ext)
}

func (a *Artifact) ComputeChecksum() string {
	h := sha256.Sum256(a.Data)
	return hex.EncodeToString(h[:])
}

func (a *Artifact) ChecksumLine() string {
	if a.Checksum == "" {
		a.Checksum = a.ComputeChecksum()
	}
	return fmt.Sprintf("%s  %s", a.Checksum, a.BinaryName())
}

type ChecksumEntry struct {
	Checksum   string
	BinaryName string
}

type ChecksumFile struct {
	Entries []ChecksumEntry
}

func ParseChecksums(content string) (*ChecksumFile, error) {
	cf := &ChecksumFile{}
	for _, line := range strings.Split(strings.TrimSpace(content), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) < 2 {
			return nil, fmt.Errorf("invalid checksum line: %s", line)
		}
		cf.Entries = append(cf.Entries, ChecksumEntry{
			Checksum:   parts[0],
			BinaryName: strings.Join(parts[1:], " "),
		})
	}
	return cf, nil
}

func (cf *ChecksumFile) Verify(artifact *Artifact) bool {
	artifact.Checksum = artifact.ComputeChecksum()
	wantName := artifact.BinaryName()
	for _, entry := range cf.Entries {
		if entry.BinaryName == wantName && entry.Checksum == artifact.Checksum {
			return true
		}
	}
	return false
}

func main() {
	fmt.Println("Release artifacts")
	fmt.Printf("Version: %s, Commit: %s, Built: %s\n", Version, Commit, BuildTime)

	app := NewArtifact("myapp", "v1.2.3", "linux", "amd64", []byte("binary-content"))
	app.Checksum = app.ComputeChecksum()
	fmt.Printf("Binary: %s\n", app.BinaryName())
	fmt.Printf("SHA256: %s\n", app.Checksum)

	cf := &ChecksumFile{
		Entries: []ChecksumEntry{
			{Checksum: app.Checksum, BinaryName: app.BinaryName()},
		},
	}
	if cf.Verify(app) {
		fmt.Println("Integrity: OK")
	}
}
