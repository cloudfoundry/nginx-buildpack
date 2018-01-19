package main

import (
	"html/template"
	"io/ioutil"
	"log"
	"os"
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

	hash := map[string]string{"Port": os.Getenv("PORT")}

	t, err := template.New("conf").Parse(string(body))
	if err != nil {
		log.Fatalf("Could not parse config file: %s", err)
	}

	if err := t.Execute(fileHandle, hash); err != nil {
		log.Fatalf("Could not write config file: %s", err)
	}
}
