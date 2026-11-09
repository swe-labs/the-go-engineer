package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestMain(m *testing.M) {
	dir := filepath.Join(".", "testdata")
	os.MkdirAll(dir, 0755)

	users := []User{
		{Name: "Alice", Age: 30},
		{Name: "Bob", Age: 25},
	}
	data, _ := json.Marshal(users)
	os.WriteFile(filepath.Join(dir, "users.json"), data, 0644)

	code := m.Run()

	os.RemoveAll(dir)
	os.Exit(code)
}

func TestLoadUsers(t *testing.T) {
	users, err := LoadUsers()
	if err != nil {
		t.Fatalf("LoadUsers failed: %v", err)
	}
	if len(users) != 2 {
		t.Fatalf("got %d users; want 2", len(users))
	}
	if users[0].Name != "Alice" || users[0].Age != 30 {
		t.Errorf("users[0] = %+v; want {Alice 30}", users[0])
	}
	if users[1].Name != "Bob" || users[1].Age != 25 {
		t.Errorf("users[1] = %+v; want {Bob 25}", users[1])
	}
}

func TestLoadUsersMissing(t *testing.T) {
	// Temporarily remove the file to test the error path
	orig := filepath.Join("testdata", "users.json")
	backup := filepath.Join("testdata", "users.json.bak")
	os.Rename(orig, backup)
	defer os.Rename(backup, orig)

	_, err := LoadUsers()
	if err == nil {
		t.Error("expected error for missing file")
	}
}
