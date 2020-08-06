package supply

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"html/template"

	"github.com/cloudfoundry/libbuildpack"
)

type Command interface {
	Execute(string, io.Writer, io.Writer, string, ...string) error
	Output(string, string, ...string) (string, error)
	Run(cmd *exec.Cmd) error
}

type Manifest interface {
	DefaultVersion(depName string) (libbuildpack.Dependency, error)
	AllDependencyVersions(string) []string
	RootDir() string
}

type Installer interface {
	InstallDependency(dep libbuildpack.Dependency, outputDir string) error
	InstallOnlyVersion(string, string) error
}

type Stager interface {
	AddBinDependencyLink(string, string) error
	DepDir() string
	DepsIdx() string
	DepsDir() string
	BuildDir() string
	WriteProfileD(string, string) error
}

type Config struct {
	Nginx NginxConfig `yaml:"nginx"`
	Dist  string      `yaml:"dist"`
}

type NginxConfig struct {
	Version string `yaml:"version"`
}

type Supplier struct {
	Stager       Stager
	Manifest     Manifest
	Installer    Installer
	Log          *libbuildpack.Logger
	Config       Config
	Command      Command
	VersionLines map[string]string
}

func New(stager Stager, manifest Manifest, installer Installer, logger *libbuildpack.Logger, command Command) *Supplier {
	return &Supplier{
		Stager:    stager,
		Manifest:  manifest,
		Installer: installer,
		Log:       logger,
		Command:   command,
	}
}

func (s *Supplier) Run() error {
	s.Log.BeginStep("Supplying nginx")

	if err := s.InstallVarify(); err != nil {
		s.Log.Error("Failed to copy verify: %s", err.Error())
		return err
	}

	if err := s.Setup(); err != nil {
		s.Log.Error("Could not setup: %s", err.Error())
		return err
	}

	if s.Config.Dist == "openresty" {
		if err := s.InstallOpenResty(); err != nil {
			s.Log.Error("Could not install openresty: %s", err.Error())
			return err
		}
	} else {
		if err := s.InstallNGINX(); err != nil {
			s.Log.Error("Could not install nginx: %s", err.Error())
			return err
		}
	}

	if err := s.ValidateNginxConf(); err != nil {
		s.Log.Error("Could not validate nginx.conf: %s", err.Error())
		return err
	}

	if err := s.WriteProfileD(); err != nil {
		s.Log.Error("Could not write profile.d: %s", err.Error())
		return err
	}

	return nil
}

func (s *Supplier) WriteProfileD() error {
	if s.Config.Dist == "openresty" {
		err := s.Stager.WriteProfileD(
			"openresty",
			fmt.Sprintf(
				"export LD_LIBRARY_PATH=$LD_LIBRARY_PATH%s$DEPS_DIR/%s/nginx/luajit/lib\nexport LUA_PATH=$DEPS_DIR/%s/nginx/lualib/?.lua\n",
				string(os.PathListSeparator),
				s.Stager.DepsIdx(),
				s.Stager.DepsIdx(),
			))
		if err != nil {
			return err
		}
	}

	return s.Stager.WriteProfileD("nginx", fmt.Sprintf("export DEP_DIR=$DEPS_DIR/%s\nmkdir -p logs", s.Stager.DepsIdx()))
}

func (s *Supplier) InstallVarify() error {
	if exists, err := libbuildpack.FileExists(filepath.Join(s.Stager.DepDir(), "bin", "varify")); err != nil {
		return err
	} else if exists {
		return nil
	}

	return libbuildpack.CopyFile(filepath.Join(s.Manifest.RootDir(), "bin", "varify"), filepath.Join(s.Stager.DepDir(), "bin", "varify"))
}

func (s *Supplier) Setup() error {
	configPath := filepath.Join(s.Stager.BuildDir(), "buildpack.yml")
	if exists, err := libbuildpack.FileExists(configPath); err != nil {
		return err
	} else if exists {
		if err := libbuildpack.NewYAML().Load(configPath, &s.Config); err != nil {
			return err
		}
	}

	var m struct {
		VersionLines map[string]string `yaml:"version_lines"`
	}
	if err := libbuildpack.NewYAML().Load(filepath.Join(s.Manifest.RootDir(), "manifest.yml"), &m); err != nil {
		return err
	}
	s.VersionLines = m.VersionLines

	logsDirPath := filepath.Join(s.Stager.BuildDir(), "logs")
	if err := os.Mkdir(logsDirPath, os.ModePerm); err != nil {
		return fmt.Errorf("Could not create 'logs' directory: %v", err)
	}

	return nil
}

func (s *Supplier) ValidateNginxConf() error {
	if err := s.validateNginxConfHasPort(); err != nil {
		return err
	}

	if err := s.validateNGINXConfSyntax(); err != nil {
		return err
	}

	return s.CheckAccessLogging()
}

func (s *Supplier) CheckAccessLogging() error {
	contents, err := ioutil.ReadFile(filepath.Join(s.Stager.BuildDir(), "nginx.conf"))
	if err != nil {
		return err
	}

	isSetToOff, err := regexp.MatchString(`(?i)access_log\s+off`, string(contents))
	if err != nil {
		return err
	}

	if !strings.Contains(string(contents), "access_log") || isSetToOff {
		s.Log.Warning("Warning: access logging is turned off in your nginx.conf file, this may make your app difficult to debug.")
	}

	return nil
}

func (s *Supplier) InstallNGINX() error {
	dep, err := s.findMatchingVersion("nginx", s.Config.Nginx.Version)
	if err != nil {
		s.Log.Info(`Available versions: ` + strings.Join(s.availableVersions(), ", "))
		return fmt.Errorf("Could not determine version: %s", err)
	}
	if s.Config.Nginx.Version == "" {
		s.Log.BeginStep("No nginx version specified - using mainline => %s", dep.Version)
	} else {
		s.Log.BeginStep("Requested nginx version: %s => %s", s.Config.Nginx.Version, dep.Version)
	}

	dir := filepath.Join(s.Stager.DepDir(), "nginx")

	if s.isStableLine(dep.Version) {
		s.Log.Warning(`Warning: usage of "stable" versions of NGINX is discouraged in most cases by the NGINX team.`)
	}

	if err := s.Installer.InstallDependency(dep, dir); err != nil {
		return err
	}

	return s.Stager.AddBinDependencyLink(filepath.Join(dir, "sbin", "nginx"), "nginx")
}

func (s *Supplier) InstallOpenResty() error {
	versions := s.Manifest.AllDependencyVersions("openresty")
	if len(versions) < 1 {
		return fmt.Errorf("unable to find a version of openresty to install")
	}

	dep := libbuildpack.Dependency{Name: "openresty", Version: versions[len(versions)-1]}
	dir := filepath.Join(s.Stager.DepDir(), "nginx")
	if err := s.Installer.InstallDependency(dep, dir); err != nil {
		return err
	}

	return s.Stager.AddBinDependencyLink(filepath.Join(dir, "nginx", "sbin", "nginx"), "nginx")
}

func (s *Supplier) validateNginxConfHasPort() error {
	conf, err := ioutil.ReadFile(filepath.Join(s.Stager.BuildDir(), "nginx.conf"))
	if err != nil {
		return err
	}

	tmpDir, err := ioutil.TempDir("", "")
	if err != nil {
		s.Log.Error("Error creating temp dir: %v", err)
		return err
	}
	defer os.RemoveAll(tmpDir)

	checkConfFile := filepath.Join(tmpDir, "conf")
	fileHandle, err := os.Create(checkConfFile)
	if err != nil {
		s.Log.Error("Could not open tmp config file for writing: %s", err)
		return err
	}
	defer fileHandle.Close()

	randString := randomString(16)

	funcMap := template.FuncMap{
		"env": func(arg string) string {
			return ""
		},
		"port": func() string {
			return randString
		},
		"module": func(arg string) string {
			return ""
		},
		"nameservers": func() string {
			return ""
		},
	}

	t, err := template.New("conf").Option("missingkey=zero").Funcs(funcMap).Parse(string(conf))
	if err != nil {
		s.Log.Error("Could not parse tmp config file: %s", err)
		return err
	}

	if err := t.Execute(fileHandle, nil); err != nil {
		s.Log.Error("Could not write tmp config file: %s", err)
		return err
	}

	contents, err := ioutil.ReadFile(checkConfFile)
	if err != nil {
		s.Log.Error("Could not read temp config file: %v", err)
		return err
	}

	if !strings.Contains(string(contents), randString) {
		s.Log.Error("nginx.conf file must be configured to respect the value of `{{port}}`")
		return errors.New("no {{port}} in nginx.conf")
	}

	return nil
}

func randomString(strLength int) string {
	rand.Seed(time.Now().UnixNano())

	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	const numCharsPossible = len(letters)
	randString := ""

	for i := 0; i < strLength; i++ {
		randString += string(letters[rand.Intn(numCharsPossible)])
	}

	return randString
}

func (s *Supplier) validateNGINXConfSyntax() error {
	tmpConfDir, err := ioutil.TempDir("/tmp", "conf")
	if err != nil {
		return fmt.Errorf("Error creating temp nginx conf dir: %s", err.Error())
	}
	defer os.RemoveAll(tmpConfDir)

	if err := libbuildpack.CopyDirectory(s.Stager.BuildDir(), tmpConfDir); err != nil {
		return fmt.Errorf("Error copying nginx.conf: %s", err.Error())
	}

	nginxConfPath := filepath.Join(tmpConfDir, "nginx.conf")
	localModulePath := filepath.Join(s.Stager.BuildDir(), "modules")
	globalModulePath := filepath.Join(s.Stager.DepDir(), "nginx", "modules")
	buildpackYMLPath := filepath.Join(s.Stager.BuildDir(), "buildpack.yml")
	cmd := exec.Command(filepath.Join(s.Stager.DepDir(), "bin", "varify"), "-buildpack-yml-path", buildpackYMLPath, nginxConfPath, localModulePath, globalModulePath)
	cmd.Dir = tmpConfDir
	cmd.Stdout = ioutil.Discard
	cmd.Stderr = ioutil.Discard
	cmd.Env = append(os.Environ(), "PORT=8080")
	if err := s.Command.Run(cmd); err != nil {
		return err
	}

	nginxErr := &bytes.Buffer{}

	cmd = exec.Command(filepath.Join(s.Stager.DepDir(), "bin", "nginx"), "-t", "-c", nginxConfPath, "-p", tmpConfDir)
	cmd.Dir = tmpConfDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = nginxErr
	if s.Config.Dist == "openresty" {
		cmd.Env = append(os.Environ(), fmt.Sprintf("LD_LIBRARY_PATH=%s", filepath.Join(s.Stager.DepDir(), "nginx", "luajit", "lib")))
	}
	if err := s.Command.Run(cmd); err != nil {
		_, _ = fmt.Fprint(os.Stderr, nginxErr.String())
		return fmt.Errorf("nginx.conf contains syntax errors: %s", err.Error())
	}

	return nil
}

func (s *Supplier) availableVersions() []string {
	allVersions := s.Manifest.AllDependencyVersions("nginx")
	allNames := []string{}
	allSemver := []string{}
	for k, v := range s.VersionLines {
		if k != "" {
			allNames = append(allNames, k)
			allSemver = append(allSemver, v)
		}
	}
	sort.Strings(allNames)
	sort.Strings(allSemver)

	return append(append(allNames, allSemver...), allVersions...)
}

func (s *Supplier) findMatchingVersion(depName string, version string) (libbuildpack.Dependency, error) {
	if version == "" {
		if val, ok := s.VersionLines["mainline"]; ok {
			version = val
		} else {
			return libbuildpack.Dependency{}, fmt.Errorf("Could not find mainline version line in buildpack manifest to default to")
		}
	} else if val, ok := s.VersionLines[version]; ok {
		version = val
	}

	versions := s.Manifest.AllDependencyVersions(depName)
	if ver, err := libbuildpack.FindMatchingVersion(version, versions); err != nil {
		return libbuildpack.Dependency{}, err
	} else {
		version = ver
	}

	return libbuildpack.Dependency{Name: depName, Version: version}, nil
}

func (s *Supplier) isStableLine(version string) bool {
	stableLine := s.VersionLines["stable"]
	_, err := libbuildpack.FindMatchingVersion(stableLine, []string{version})
	return err == nil
}
