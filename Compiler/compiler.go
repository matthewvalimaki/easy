package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"./parser"
)

func main() {
	cacheDir := executionPath() + string(filepath.Separator) + "CompileCache" + string(filepath.Separator) + "src"
	fileToProcess := executionPath() + string(filepath.Separator) + "main.easy"
	
	// process classes
	mainFile := parser.ClassProcessor(cacheDir, fileToProcess)
	
	// build
	build(mainFile)
}

func executionPath() (string) {
    pwd, err := os.Getwd()
    if err != nil {
        os.Exit(1)
    }
	
	return pwd
}

func build(mainFile string) {
	originalGoPath := os.Getenv("GOPATH")
	os.Setenv("GOPATH", originalGoPath + ";" + executionPath() + string(filepath.Separator) + "CompileCache")

	cmd := exec.Command("go", "build", mainFile)
	stdout, err := cmd.Output()
	
	if err != nil {
	print("error")
		println(err.Error())
        return
    }
	
	print(string(stdout))
	
	os.Setenv("GOPATH", originalGoPath)
}