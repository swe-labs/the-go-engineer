package main

import (
	"os"
	"os/exec"
	"testing"
)

func TestGoSumIntegrity(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "gosum-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	cmd := exec.Command("go", "mod", "init", "example.com/gosum-test")
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

	cmd = exec.Command("go", "mod", "tidy")
	cmd.Dir = tmpDir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("tidy: %v\n%s", err, out)
	}

	if _, err := os.Stat(tmpDir + "/go.sum"); os.IsNotExist(err) {
		t.Fatal("go.sum was not created")
	}

	cmd = exec.Command("go", "mod", "verify")
	cmd.Dir = tmpDir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("verify failed: %v\n%s", err, out)
	}
}
