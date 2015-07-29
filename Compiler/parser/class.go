package parser

import (
	"fmt"
	"os"
	"bufio"
	"regexp"
	"strings"
	"path/filepath"
)

type Class struct {
	packageName, path, fileName, pathWithFile, className string
	imports []Import
	constants []string
	functions []FunctionStruct
}

var isMainPackage bool = true
var isMainFunction bool = true
var pathWithFile string
var processingFunction = false
var functionBracetCount = 0
var functionCompiler []string

func ClassProcessor(cacheDir string, path string) (string) {
	class := new(Class)

	inFile, _ := os.Open(path)
	defer inFile.Close()
	scanner := bufio.NewScanner(inFile)
	scanner.Split(bufio.ScanLines) 

	var text string
	lineNumber := 0
	
	for scanner.Scan() {
		text = strings.TrimSpace(scanner.Text())
		
		if text != "" {
			fmt.Println(text)
			
			processLine(class, lineNumber, strings.TrimSpace(text), cacheDir)
			
			lineNumber++
		}
	}
	
	// create files
	export(class, cacheDir)
	
	return mainFile(class)
}

func mainFile(class *Class) (string) {
	return class.path + string(filepath.Separator) + class.fileName + ".go"
}

func export(class *Class, cacheDir string) {
	pathWithFile := createPackage(class.path, class.fileName, class.packageName, class.className)
	class.pathWithFile = pathWithFile
	
	if len(class.imports) > 0 {
		processImports(cacheDir, class, class.imports)
		exportImports(class)
	}
	
	if len(class.constants) > 0 {
		exportConstants(class)
	}
	
	exportFunctions(class)
}

func addImport(class *Class, importStruct Import) {
	class.imports = append(class.imports, importStruct)
}

func addConstant(class *Class, constant string) {
	class.constants = append(class.constants, constant)
}

func addFunc(class *Class, funcStruct FunctionStruct) {
	class.functions = append(class.functions, funcStruct)
}

func processImports(cacheDir string, class *Class, imports []Import) {
	for _, Import := range imports {
		if Import.IsGoLibrary == false {
			ClassProcessor(cacheDir, "C:\\Users\\tvalimaki\\Desktop\\epl\\" + string(filepath.Separator) + Import.path + string(filepath.Separator) + Import.packageName + ".easy")
		}
	}
}

func processLine(class *Class, lineNumber int, line string, cacheDir string) {
	if lineNumber == 0 {
		path, fileName, packageName, className := createPackagePath(line, cacheDir)
		class.packageName = packageName
		class.path = path
		class.fileName = fileName
		class.className = className
	} else {
		isFuncDeclaration, _ := regexp.MatchString("^[[:alnum:]]*\\(.*\\)", line)
	
		if strings.Contains(line, "use ") {
			addImport(class, UseProcessor(line))
		} else if strings.Contains(line, "const ") {
			addConstant(class, ConstProcessor(line))
		} else if isFuncDeclaration && processingFunction == false {
			processingFunction = true
			
			functionCompiler = append(functionCompiler, line)
			functionBracetCount = strings.Count(line, "{") - strings.Count(line, "}")
		} else {
			if processingFunction {
				if functionBracetCount > 0 {
					newBracetCount := strings.Count(line, "{") - strings.Count(line, "}")
					
					if newBracetCount == 0 {
						functionCompiler = append(functionCompiler, line)
					} else if strings.Count(line, "{") > 0 {
						functionBracetCount = functionBracetCount + newBracetCount
					} else {
						functionBracetCount = functionBracetCount - strings.Count(line, "}")		
						
						if functionBracetCount == 0 {
							addFunc(class, FunctionProcessor(class, functionCompiler, isMainFunction))
							
							isMainFunction = false
							functionCompiler = nil
							processingFunction = false
						}
					}
				}
			}
		}
	}
}

func createPackagePath(namespace string, cacheDir string) (string, string, string, string) {
	directories := strings.Split(namespace, ".")
	var packageName string
	var fileName string
	var className string
	
	if len(directories) == 1 {
		className = namespace
		packageName = strings.ToLower(namespace)
		fileName = strings.ToLower(namespace)
	} else {
		className = directories[len(directories) - 1]
		packageName = strings.ToLower(directories[len(directories) - 2])
		fileName = strings.ToLower(directories[len(directories) - 1])
		// remove last as that is the file name
		directories = directories[:len(directories) - 1]
	}
		
	if isMainPackage {
		packageName = "main"
		isMainPackage = false
	}
	
	currentDirectory := cacheDir	
	
	for _, value := range directories {
		currentDirectory = currentDirectory + string(filepath.Separator) + strings.ToLower(value)
		
		os.MkdirAll(currentDirectory, 0777)
	}
	
	return currentDirectory, fileName, packageName, className
}

func createPackage(path string, fileName string, packageName string, className string) (string) {
	pathWithFile := path + string(filepath.Separator) + fileName + ".go"
	
	contents := "package " + className + "\n"
	
	// write contents
	addLine(pathWithFile, contents)
	
	return pathWithFile
}

func exportImports(class *Class) {
	addLine(class.pathWithFile, "import (\n")

	for _, Import := range class.imports {
		if Import.IsGoLibrary {
			addLine(class.pathWithFile, "\"" + Import.packageName + "\"\n")
		} else {
			addLine(class.pathWithFile, "\"" + Import.path + "\"\n")
		}
	}
	
	addLine(class.pathWithFile, ")\n")
}

func exportConstants(class *Class) {
	for _, constant := range class.constants {
		addLine(class.pathWithFile, constant)
	}
}

func exportFunctions(class *Class) {
	for _, FunctionStruct := range class.functions {
		addLine(class.pathWithFile, FunctionStruct.declaration)
		addLine(class.pathWithFile, " {\n")
		addLine(class.pathWithFile, FunctionStruct.body + "\n")
		addLine(class.pathWithFile, "}\n")
	}
}

func addLine(pathWithFile string, line string) {
	f, _ := os.OpenFile(pathWithFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	_, err := f.WriteString(line)
	if err != nil {
		fmt.Println("Cannot append to file: " + pathWithFile)
	}
	f.Close()
}

func classImportsHaveNamespace(class *Class, namespace string) (bool) {
	for _, Import := range class.imports {
		if Import.namespace == namespace {
			return true
		
			break
		}
	}
	
	return false
}

