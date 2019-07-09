package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"github.com/miekg/dns"

	"github.com/cloudfoundry/libbuildpack"
)

func main() {
	filename := os.Args[1]
	localModulePath := os.Args[2]
	globalModulePath := os.Args[3]
	resolvConfPath := "/etc/resolv.conf"
	if len(os.Args) > 4 && len(os.Args[4]) > 0 {
		resolvConfPath = os.Args[4]
	}
	// https://github.com/cloudfoundry/bosh-dns-release/blob/master/jobs/bosh-dns/spec#L36-L38
	defaultNameServer := "169.254.0.2"
	if len(os.Args) > 5 && len(os.Args[5]) > 0 {
		defaultNameServer = os.Args[5]
	}
	nameServers, err := readNameServers(resolvConfPath, defaultNameServer)
	if err != nil {
		nameServers = []string{defaultNameServer}
		log.Printf("Could not open %s file for reading. "+
			"The default nameservers %s will be used. Error: %s", resolvConfPath, defaultNameServer, err)
	}

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
			return strings.Join(nameServers, " ")
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

func readNameServers(resolvConfPath string, defaultNameServer string) ([]string, error) {
	result := []string{defaultNameServer}

	config, err := dns.ClientConfigFromFile(resolvConfPath)
	if err != nil {
		return []string{}, err
	}

	if len(config.Servers) > 0 {
		return config.Servers, nil
	}

	return result, nil
}
