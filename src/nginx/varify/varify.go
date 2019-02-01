package main

import (
	"bufio"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/cloudfoundry/libbuildpack"
)

func main() {
	filename := os.Args[1]
	localModulePath := os.Args[2]
	globalModulePath := os.Args[3]
	resolvConfPath := "/etc/resolv.conf"
	if len(os.Args) > 4 {
		resolvConfPath = os.Args[4]
	}
	defaultNameserver := "169.254.0.2" // https://github.com/cloudfoundry/bosh-dns-release/blob/master/jobs/bosh-dns/spec#L36-L38

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
		"nameservers": func() string {
			resolvConfFile, err := os.Open(resolvConfPath)
			if err != nil {
				log.Printf("Could not open %s file for reading. "+
					"The default nameservers %s will be used. Error: %s", resolvConfPath, defaultNameserver, err)
				return defaultNameserver
			}
			defer resolvConfFile.Close()
			scanner := bufio.NewScanner(resolvConfFile)
			var nameServers strings.Builder
			for scanner.Scan() {
				var line = scanner.Text()
				// Lines that contain a semicolon (;) or hash character (#) in the first column are treated as comments.
				if len(line)>0 && line[0]!='#' && line[0]!=';' {
					line = strings.TrimSpace(line)
					matches, _ := regexp.MatchString(`nameserver \d+\.\d+\.\d+\.\d+`, line)
					if matches {
						if nameServers.Len() > 0 {
							nameServers.WriteString(" ")
						}
						nameServers.WriteString(strings.TrimSpace(line[11:]))
					} else if len(line)>0 {
						log.Printf("nameservers (%s): line '%s' has been ignored", resolvConfPath, line)
					}
				}
			}
			if nameServers.Len()>0 {
				return nameServers.String();
			} else {
				log.Printf("Could not find any nameserver in %s. "+
					"The default nameserver %s will be used.", resolvConfPath, defaultNameserver)
				return defaultNameserver
			}
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
