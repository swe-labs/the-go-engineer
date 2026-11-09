package main

import (
	"os"
	"os/exec"
	"testing"
)

func TestDependencyManagement(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "dep-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	cmd := exec.Command("go", "mod", "init", "example.com/dep-test")
	cmd.Dir = tmpDir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("init: %v\n%s", err, out)
	}

	mainSrc := `package main
import "fmt"
import "rsc.io/quote"
func main() { fmt.Println(quote.Hello()) }
`
	if err := os.WriteFile(tmpDir+"/main.go", []byte(mainSrc), 0644); err != nil {
		t.Fatal(err)
	}

	cmd = exec.Command("go", "get", "rsc.io/quote@v1.5.2")
	cmd.Dir = tmpDir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("go get: %v\n%s", err, out)
	}

	cmd = exec.Command("go", "build", "-o", tmpDir+"/demo", ".")
	cmd.Dir = tmpDir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("build: %v\n%s", err, out)
	}
}
