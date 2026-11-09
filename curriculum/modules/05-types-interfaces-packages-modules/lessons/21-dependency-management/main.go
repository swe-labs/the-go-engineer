package main

import (
	"fmt"
	"os"
	"os/exec"
)

func runCmd(dir, name string, args ...string) {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	fmt.Printf("> %s %v\n%s", name, args, out)
	if err != nil {
		panic(fmt.Sprintf("%v: %v\n%s", name, err, out))
	}
}

func main() {
	tmpDir, err := os.MkdirTemp("", "dep-demo")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(tmpDir)

	runCmd(tmpDir, "go", "mod", "init", "example.com/dep-demo")

	mainSrc := `package main
import "fmt"
import "rsc.io/quote"
func main() { fmt.Println(quote.Hello()) }
`
	if err := os.WriteFile(tmpDir+"/main.go", []byte(mainSrc), 0644); err != nil {
		panic(err)
	}

	runCmd(tmpDir, "go", "get", "rsc.io/quote@v1.5.2")

	data, _ := os.ReadFile(tmpDir + "/go.mod")
	fmt.Printf("go.mod:\n%s\n", data)

	runCmd(tmpDir, "go", "list", "-m", "all")

	runCmd(tmpDir, "go", "mod", "tidy")

	fmt.Println("Dependency management demo complete.")
}
