package main

import (
	"os"
	"os/exec"
	"testing"
)

func TestModuleInitAndBuild(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "gomod-demo-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	cmd := exec.Command("go", "mod", "init", "example.com/mymod")
	cmd.Dir = tmpDir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("go mod init failed: %v\n%s", err, out)
	}

	mainSrc := `package main
import "fmt"
func main() { fmt.Println("ok") }`
	if err := os.WriteFile(tmpDir+"/main.go", []byte(mainSrc), 0644); err != nil {
		t.Fatal(err)
	}

	cmd = exec.Command("go", "mod", "tidy")
	cmd.Dir = tmpDir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("go mod tidy failed: %v\n%s", err, out)
	}

	cmd = exec.Command("go", "build", "-o", tmpDir+"/demo", ".")
	cmd.Dir = tmpDir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("go build failed: %v\n%s", err, out)
	}
}
