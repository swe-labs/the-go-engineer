package main

import (
	"io/fs"
	"os"
	"path/filepath"
	"testing"
	"testing/fstest"
)

func TestEmbeddedFSRead(t *testing.T) {
	data, err := fs.ReadFile(embeddedFiles, "static/hello.txt")
	if err != nil {
		t.Fatal(err)
	}
	if len(data) == 0 {
		t.Error("expected non-empty embedded file")
	}
}

func TestEmbeddedFSWalk(t *testing.T) {
	count := 0
	err := fs.WalkDir(embeddedFiles, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		count++
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if count == 0 {
		t.Error("expected at least one embedded file")
	}
}

func TestDirFS(t *testing.T) {
	d := os.DirFS(".")
	data, err := fs.ReadFile(d, "main.go")
	if err != nil {
		t.Fatal(err)
	}
	if len(data) == 0 {
		t.Error("expected non-empty main.go")
	}
}

func TestMapFS(t *testing.T) {
	testFS := fstest.MapFS{
		"hello.txt":        &fstest.MapFile{Data: []byte("hello")},
		"sub/greeting.txt": &fstest.MapFile{Data: []byte("hi")},
	}

	data, err := fs.ReadFile(testFS, "hello.txt")
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "hello" {
		t.Errorf("expected 'hello', got %q", data)
	}

	entries, err := fs.ReadDir(testFS, ".")
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(entries))
	}
}

func TestSubFS(t *testing.T) {
	testFS := fstest.MapFS{
		"a/x.txt": &fstest.MapFile{Data: []byte("x")},
		"a/y.txt": &fstest.MapFile{Data: []byte("y")},
		"b/z.txt": &fstest.MapFile{Data: []byte("z")},
	}

	sub, err := fs.Sub(testFS, "a")
	if err != nil {
		t.Fatal(err)
	}

	data, err := fs.ReadFile(sub, "x.txt")
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "x" {
		t.Errorf("expected 'x', got %q", data)
	}
}

func TestValidatedFS(t *testing.T) {
	testFS := fstest.MapFS{
		"good.txt": &fstest.MapFile{Data: []byte("ok")},
	}

	err := fstest.TestFS(testFS, "good.txt")
	if err != nil {
		t.Fatal(err)
	}
}

func TestGlobFS(t *testing.T) {
	testFS := fstest.MapFS{
		"a.go":  &fstest.MapFile{Data: []byte("a")},
		"b.go":  &fstest.MapFile{Data: []byte("b")},
		"c.txt": &fstest.MapFile{Data: []byte("c")},
	}

	matches, err := fs.Glob(testFS, "*.go")
	if err != nil {
		t.Fatal(err)
	}
	if len(matches) != 2 {
		t.Errorf("expected 2 .go files, got %d", len(matches))
	}
}

func TestOsDirFSWalk(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "f1.txt"), []byte("1"), 0644)
	os.WriteFile(filepath.Join(dir, "f2.txt"), []byte("2"), 0644)
	os.Mkdir(filepath.Join(dir, "sub"), 0755)

	d := os.DirFS(dir)
	count := 0
	fs.WalkDir(d, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		count++
		return nil
	})
	if count != 4 {
		t.Errorf("expected 4 entries (., f1.txt, f2.txt, sub/), got %d", count)
	}
}
