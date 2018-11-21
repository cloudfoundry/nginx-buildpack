package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/cloudfoundry/libbuildpack"
)

func main() {
	filename := os.Args[1]
	localModulePath := os.Args[2]
	globalModulePath := os.Args[3]

	body, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("Could not read config file: %s: %s", filename, err)
	}

	fileHandle, err := os.Create(filename)
	if err != nil {
		log.Fatalf("Could not open config file for writing: %s", err)
	}
	defer fileHandle.Close()

	funcMap := template.FuncMap{
		"env": os.Getenv,
		"port": func() string {
			return os.Getenv("PORT")
		},
		"module": func(name string) string {
			pathToModules := globalModulePath
			foundLocally, err := libbuildpack.FileExists(filepath.Join(localModulePath, name+".so"))
			if err != nil {
				log.Fatalf("Error looking for module in user provided modules directory: %s", err)
			}
			if foundLocally {
				pathToModules = localModulePath
			}
			return fmt.Sprintf("load_module %s.so;", filepath.Join(pathToModules, name))
		},
	}

	t, err := template.New("conf").Option("missingkey=zero").Funcs(funcMap).Parse(string(body))
	if err != nil {
		log.Fatalf("Could not parse config file: %s", err)
	}

	if err := t.Execute(fileHandle, nil); err != nil {
		log.Fatalf("Could not write config file: %s", err)
	}
}
