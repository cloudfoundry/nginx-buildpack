package main

import (
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"fmt"
	"path/filepath"
)

func main() {
	filename := os.Args[1]

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
			return fmt.Sprintf("load_module %s.so;", filepath.Join(os.Getenv("NGINX_MODULES"), name))
		},
	}

	t, err := template.New("conf").Funcs(funcMap).Parse(string(body))
	if err != nil {
		log.Fatalf("Could not parse config file: %s", err)
	}

	if err := t.Execute(fileHandle, nil); err != nil {
		log.Fatalf("Could not write config file: %s", err)
	}
}
