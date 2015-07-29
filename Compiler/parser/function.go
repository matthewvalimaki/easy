package parser

import (
	"strings"
)

type FunctionStruct struct {
	declaration, body string
	classMapping map[string]string
}

func FunctionProcessor(class *Class, functionCompilation []string, isMain bool) (FunctionStruct) {
	var functionStruct FunctionStruct
	functionStruct.classMapping = make(map[string]string)

	function := ""

	index := 0
	for _, line := range functionCompilation {
		if index == 0 {
			functionStruct.declaration = functionDeclaration(line, isMain)
		} else {
			function = function + functionInternal(class, functionStruct, line)
		}
		
		index++
	}
	
	functionStruct.body = replaceNamespaces(class, function)
	functionStruct.body = replaceClassMappings(functionStruct, function)
	
	return functionStruct
}

func replaceClassMappings(functionStruct FunctionStruct, originalFunction string) (string) {
	function := originalFunction

	for key, value := range functionStruct.classMapping {
	print(value + ":Waaa")
		function = strings.Replace(function, key, value, -1)
	}
	
	return function
}

func replaceNamespaces(class *Class, originalFunction string) (string) {
	function := originalFunction

	for _, Import := range class.imports {
		function = strings.Replace(function, Import.namespace, Import.path, -1)
	}
	
	return function
}

func functionDeclaration(line string, isMain bool) (string) {
	declaration := line

	if isMain {
		// replace the first occurrence only
		declaration = strings.Replace(declaration, "construct", "main", 1)
	}

	return "func " + strings.TrimSpace(declaration[:strings.Index(declaration, "{") - 1])
}

func functionInternal(class *Class, functionStruct FunctionStruct, line string) (string) {
	doNotAddLine := false
	internal := ""

	if strings.Contains(line, "print ") {
		internal = internal + "print(" + line[6:] + ")\n"
	} else {
		if strings.Index(line, " = new ") != -1 {
			doNotAddLine = true
			
			namespace := strings.Split(line, ".")
			functionStruct.classMapping[line[:strings.Index(line, " ")]] = namespace[len(namespace) - 1]
			
			if strings.Index(line, ".") != -1 {
				addImport(class, UseProcessor(line[strings.Index(line, "= new ") + 6:]))
			} else {
				ClassProcessor("C:\\Users\\tvalimaki\\Desktop\\epl\\Tests\\HelloWorld\\CompileCache\\src", "C:\\Users\\tvalimaki\\Desktop\\epl\\Tests\\HelloWorld\\" + strings.ToLower(line[strings.Index(line, "= new ") + 6:]) + ".easy")
			}
		} else if strings.Index(line, "fmt.") != -1 {
			addImport(class, Import{packageName: "fmt", IsGoLibrary: true})
		} else if strings.Index(line, "md5.") != -1 {		
			if !classImportsHaveNamespace(class, "EPLFramework.Kernel.Crypto.MD5") {
				doNotAddLine = true
				
				addImport(class, Import{packageName: "fmt", IsGoLibrary: true})
				addImport(class, Import{packageName: "crypto/md5", IsGoLibrary: true})
				
				internal = internal + strings.Replace(line, "md5.Sum(text)", "fmt.Sprintf(\"%x\", md5.Sum([]byte(text)))", -1)
			}
		}
	
		if !doNotAddLine {
			internal = internal + line + "\n"
		}
	}
	
	return internal
}

func addPrint(class *Class, originalCommand string) (string) {
	addImport(class, Import{packageName: "fmt", IsGoLibrary: true})

	return "fmt.Println(" + originalCommand + ")"
}