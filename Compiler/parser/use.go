package parser

import (
	"strings"
)

type Import struct {
	packageName, parentPackageName, path, namespace string
	IsGoLibrary bool
}

func UseProcessor(line string) (Import) {
	namespace := strings.TrimSpace(strings.Replace(line, "use", "", -1))
	path := strings.Replace(namespace, ".", "/", -1)
	path = strings.ToLower(path)
	
	packageName := path[strings.LastIndex(path, "/") + 1:]
	path = strings.Replace(path, "/" + packageName, "", 1)
	parentPackageName := path[strings.LastIndex(path, "/") + 1:]
	
	return Import{packageName: packageName, parentPackageName: parentPackageName, path: path, namespace: namespace, IsGoLibrary: false};
}