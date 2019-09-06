package main

import (
	"bytes"
	"flag"
	"fmt"
	htmlTemplate "html/template"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	textTemplate "text/template"

	"gopkg.in/yaml.v2"

	"github.com/miekg/dns"

	"github.com/cloudfoundry/libbuildpack"
)

func main() {

	buildpackYMLPath := flag.String("buildpack-yml-path", "", "path to buildpack.yml file")

	flag.Parse()

	flag.Args()

	filename := flag.Args()[0]
	localModulePath := flag.Args()[1]
	globalModulePath := flag.Args()[2]
	resolvConfPath := "/etc/resolv.conf"
	if len(flag.Args()) > 3 && len(flag.Args()[3]) > 0 {
		resolvConfPath = flag.Args()[3]
	}
	// https://github.com/cloudfoundry/bosh-dns-release/blob/master/jobs/bosh-dns/spec#L36-L38
	defaultNameServer := "169.254.0.2"
	if len(flag.Args()) > 4 && len(flag.Args()[4]) > 0 {
		defaultNameServer = flag.Args()[4]
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

	confBuf := bytes.Buffer{}
	tempConfWriter := io.Writer(&confBuf)

	fileHandle, err := os.Create(filename)
	if err != nil {
		log.Fatalf("Could not open config file for writing: %s", err)
	}
	defer fileHandle.Close()

	plainTextEnvVars, err := getPlaintextEnvVars(*buildpackYMLPath)
	if err != nil {
		log.Fatalf("Unable to read buildpath.yml path '%s'", *buildpackYMLPath)
	}

	plainTextFuncMap := textTemplate.FuncMap{
		"env":         safeEnv(plainTextEnvVars),
		"port":        noArgIdentity("port"),
		"module":      singleArgIdentity("module"),
		"nameservers": noArgIdentity("nameservers"),
	}

	htmlFuncMap := htmlTemplate.FuncMap{
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

	textT, err := textTemplate.New("tempconf").Option("missingkey=zero").Funcs(plainTextFuncMap).Parse(string(body))
	if err != nil {
		log.Fatalf("Could not parse config file: %s", err)
	}
	if err := textT.Execute(tempConfWriter, nil); err != nil {
		log.Fatalf("Could not write temp config to buffer: %s", err)
	}

	fmt.Printf("plain text tempate: %s", confBuf.String())

	htmlT, err := htmlTemplate.New("tempconf").Option("missingkey=zero").Funcs(htmlFuncMap).Parse(confBuf.String())
	if err != nil {
		log.Fatalf("Could not parse config file: %s", err)
	}

	if err := htmlT.Execute(fileHandle, nil); err != nil {
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

type BuildpackYML struct {
	Nginx struct {
		PlaintextEnvVars []string `yaml:"plaintext_env_vars"`
	} `yaml:"nginx"`
}

func getPlaintextEnvVars(bpYMLPath string) ([]string, error) {
	var bpYML BuildpackYML
	exists, err := libbuildpack.FileExists(bpYMLPath)
	if err != nil {
		return []string{}, err
	} else if bpYMLPath == "" || !exists {
		return []string{}, nil
	}

	bpYMLContents, err := ioutil.ReadFile(bpYMLPath)
	if err != nil {
		return []string{}, err
	}

	if err = yaml.Unmarshal(bpYMLContents, &bpYML); err != nil {
		return []string{}, err
	}

	return bpYML.Nginx.PlaintextEnvVars, nil
}

func safeEnv(keys []string) func(string) string {
	return func(key string) string {
		for _, safeKey := range keys {
			if key == safeKey {
				return os.Getenv(key)
			}
		}
		return fmt.Sprintf(`{{env "%s"}}`, key)
	}
}

func singleArgIdentity(key string) func(string) string {
	return func(val string) string {
		return fmt.Sprintf(`{{%s "%s"}}`, key, val)
	}
}

func noArgIdentity(key string) func() string {
	return func() string {
		return fmt.Sprintf(`{{%s}}`, key)
	}
}
