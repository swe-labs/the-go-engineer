package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	tmpDir, err := os.MkdirTemp("", "go-mod-demo")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(tmpDir)

	// Write go.mod by hand.
	goMod := `module example.com/handcrafted

go 1.25.0

require rsc.io/quote v1.5.2
`
	if err := os.WriteFile(tmpDir+"/go.mod", []byte(goMod), 0644); err != nil {
		panic(err)
	}

	// Write main.go.
	mainSrc := `package main

import (
	"fmt"
	"rsc.io/quote"
)

func main() {
	fmt.Println(quote.Hello())
}
`
	if err := os.WriteFile(tmpDir+"/main.go", []byte(mainSrc), 0644); err != nil {
		panic(err)
	}

	fmt.Println("Module created. Running go mod tidy...")
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = tmpDir
	if out, err := cmd.CombinedOutput(); err != nil {
		panic(fmt.Sprintf("go mod tidy: %v\n%s", err, out))
	}

	fmt.Println("Building...")
	cmd = exec.Command("go", "build", "-o", tmpDir+"/demo", ".")
	cmd.Dir = tmpDir
	if out, err := cmd.CombinedOutput(); err != nil {
		panic(fmt.Sprintf("go build: %v\n%s", err, out))
	}

	fmt.Println("Running...")
	cmd = exec.Command(tmpDir + "/demo")
	out, err := cmd.CombinedOutput()
	if err != nil {
		panic(fmt.Sprintf("run: %v\n%s", err, out))
	}
	fmt.Printf("Output: %s", out)

	fmt.Println("\ngo.mod demo complete.")
}
