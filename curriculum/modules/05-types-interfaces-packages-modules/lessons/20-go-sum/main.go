package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	tmpDir, err := os.MkdirTemp("", "gosum-demo")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(tmpDir)

	// Init module.
	cmd := exec.Command("go", "mod", "init", "example.com/gosum-demo")
	cmd.Dir = tmpDir
	if out, err := cmd.CombinedOutput(); err != nil {
		panic(fmt.Sprintf("init: %v\n%s", err, out))
	}

	// Create main.go with an import.
	mainSrc := `package main
import "fmt"
import "rsc.io/quote"
func main() { fmt.Println(quote.Hello()) }
`
	if err := os.WriteFile(tmpDir+"/main.go", []byte(mainSrc), 0644); err != nil {
		panic(err)
	}

	// Run go mod tidy to populate go.sum.
	cmd = exec.Command("go", "mod", "tidy")
	cmd.Dir = tmpDir
	if out, err := cmd.CombinedOutput(); err != nil {
		panic(fmt.Sprintf("tidy: %v\n%s", err, out))
	}

	// Read and display go.sum.
	data, err := os.ReadFile(tmpDir + "/go.sum")
	if err != nil {
		panic(err)
	}
	fmt.Println("=== go.sum contents ===")
	fmt.Println(string(data))

	// Run go mod verify.
	cmd = exec.Command("go", "mod", "verify")
	cmd.Dir = tmpDir
	out, err := cmd.CombinedOutput()
	if err != nil {
		panic(fmt.Sprintf("verify: %v\n%s", err, out))
	}
	fmt.Printf("go mod verify: %s", out)

	fmt.Println("go.sum demo complete.")
}
