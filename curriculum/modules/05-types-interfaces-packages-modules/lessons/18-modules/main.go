package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	// Demonstrate module operations by creating a temporary module.
	tmpDir, err := os.MkdirTemp("", "gomod-demo")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(tmpDir)

	// Step 1: go mod init
	cmd := exec.Command("go", "mod", "init", "example.com/gomod-demo")
	cmd.Dir = tmpDir
	out, err := cmd.CombinedOutput()
	if err != nil {
		panic(fmt.Sprintf("go mod init failed: %v\n%s", err, out))
	}
	fmt.Printf("Initialized module in %s\n", tmpDir)

	// Step 2: Create main.go
	mainSrc := `package main

import "fmt"

func main() {
	fmt.Println("Module demo works!")
}
`
	if err := os.WriteFile(tmpDir+"/main.go", []byte(mainSrc), 0644); err != nil {
		panic(err)
	}

	// Step 3: go mod tidy
	cmd = exec.Command("go", "mod", "tidy")
	cmd.Dir = tmpDir
	out, err = cmd.CombinedOutput()
	if err != nil {
		panic(fmt.Sprintf("go mod tidy failed: %v\n%s", err, out))
	}
	fmt.Printf("go mod tidy succeeded\n")

	// Step 4: go build
	cmd = exec.Command("go", "build", "-o", tmpDir+"/demo", ".")
	cmd.Dir = tmpDir
	out, err = cmd.CombinedOutput()
	if err != nil {
		panic(fmt.Sprintf("go build failed: %v\n%s", err, out))
	}
	fmt.Printf("Build succeeded: %s\n", tmpDir+"/demo")

	// Step 5: Run
	cmd = exec.Command(tmpDir + "/demo")
	out, err = cmd.CombinedOutput()
	if err != nil {
		panic(fmt.Sprintf("exec failed: %v\n%s", err, out))
	}
	fmt.Printf("Output: %s", out)

	fmt.Println("\nModule creation and build demo complete.")
}
