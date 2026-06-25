package main

import (
	"os"
	"os/exec"
	"testing"
)

func TestWorkspaceDemo(t *testing.T) {
	baseDir, err := os.MkdirTemp("", "workspace-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(baseDir)

	appDir := baseDir + "/app"
	libDir := baseDir + "/lib"

	os.MkdirAll(appDir, 0755)
	os.MkdirAll(libDir, 0755)

	cmd := exec.Command("go", "mod", "init", "example.com/lib")
	cmd.Dir = libDir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("lib init: %v\n%s", err, out)
	}

	libSrc := `package lib
func Help() string { return "lib help" }
`
	os.WriteFile(libDir+"/lib.go", []byte(libSrc), 0644)

	cmd = exec.Command("go", "mod", "init", "example.com/app")
	cmd.Dir = appDir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("app init: %v\n%s", err, out)
	}

	appSrc := `package main
import "fmt"
import "example.com/lib"
func main() { fmt.Println(lib.Help()) }
`
	os.WriteFile(appDir+"/main.go", []byte(appSrc), 0644)

	cmd = exec.Command("go", "work", "init", "./app", "./lib")
	cmd.Dir = baseDir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("work init: %v\n%s", err, out)
	}

	if _, err := os.Stat(baseDir + "/go.work"); os.IsNotExist(err) {
		t.Fatal("go.work was not created")
	}

	cmd = exec.Command("go", "run", "./app")
	cmd.Dir = baseDir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("run with workspace: %v\n%s", err, out)
	}
}
