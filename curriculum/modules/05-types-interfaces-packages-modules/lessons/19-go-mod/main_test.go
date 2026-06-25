package main

import (
	"os"
	"os/exec"
	"testing"
)

func TestHandcraftedGoMod(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "go-mod-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	goMod := `module example.com/handcrafted

go 1.25.0

require rsc.io/quote v1.5.2
`
	if err := os.WriteFile(tmpDir+"/go.mod", []byte(goMod), 0644); err != nil {
		t.Fatal(err)
	}

	mainSrc := `package main
import "fmt"
import "rsc.io/quote"
func main() { fmt.Println(quote.Hello()) }
`
	if err := os.WriteFile(tmpDir+"/main.go", []byte(mainSrc), 0644); err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = tmpDir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("go mod tidy: %v\n%s", err, out)
	}

	cmd = exec.Command("go", "build", "-o", tmpDir+"/demo", ".")
	cmd.Dir = tmpDir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("go build: %v\n%s", err, out)
	}
}
