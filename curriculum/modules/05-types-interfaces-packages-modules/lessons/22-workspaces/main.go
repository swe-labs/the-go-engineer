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
	baseDir, err := os.MkdirTemp("", "workspace-demo")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(baseDir)

	appDir := baseDir + "/app"
	libDir := baseDir + "/lib"

	os.MkdirAll(appDir, 0755)
	os.MkdirAll(libDir, 0755)

	// Init lib module.
	runCmd(libDir, "go", "mod", "init", "example.com/lib")

	libSrc := `package lib

func Help() string {
	return "Help from lib!"
}
`
	os.WriteFile(libDir+"/lib.go", []byte(libSrc), 0644)

	// Init app module.
	runCmd(appDir, "go", "mod", "init", "example.com/app")

	appSrc := `package main

import (
	"fmt"
	"example.com/lib"
)

func main() {
	fmt.Println(lib.Help())
}
`
	os.WriteFile(appDir+"/main.go", []byte(appSrc), 0644)

	// Create go.work.
	runCmd(baseDir, "go", "work", "init", "./app", "./lib")

	// Build and run app with workspace.
	runCmd(baseDir, "go", "run", "./app")

	fmt.Println("Workspace demo complete.")
}
