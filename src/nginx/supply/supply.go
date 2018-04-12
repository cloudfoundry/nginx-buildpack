package supply

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/cloudfoundry/libbuildpack"
)

type Command interface {
	Execute(string, io.Writer, io.Writer, string, ...string) error
	Output(string, string, ...string) (string, error)
	Run(cmd *exec.Cmd) error
}

type Manifest interface {
	InstallOnlyVersion(string, string) error
	DefaultVersion(depName string) (libbuildpack.Dependency, error)
	AllDependencyVersions(string) []string
	InstallDependency(dep libbuildpack.Dependency, outputDir string) error
	RootDir() string
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
	Version string `yaml:"version"`
}

type Supplier struct {
	Stager       Stager
	Manifest     Manifest
	Log          *libbuildpack.Logger
	Config       Config
	Command      Command
	VersionLines map[string]string
}

func New(stager Stager, manifest Manifest, logger *libbuildpack.Logger, command Command) *Supplier {
	return &Supplier{
		Stager:   stager,
		Manifest: manifest,
		Log:      logger,
		Command:  command,
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

	if err := s.InstallNginx(); err != nil {
		s.Log.Error("Could not install nginx: %s", err.Error())
		return err
	}

	if err := s.validateNginxConf(); err != nil {
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
	return s.Stager.WriteProfileD("nginx", fmt.Sprintf("export NGINX_MODULES=%s", filepath.Join("$DEPS_DIR", s.Stager.DepsIdx(), "nginx", "nginx", "modules")))
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
	configPath := filepath.Join(s.Stager.BuildDir(), "nginx.yml")
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

	return nil
}

func (s *Supplier) validateNginxConf() error {
	if err := s.validateNginxConfExists(); err != nil {
		return err
	}
	if err := s.validateNginxConfHasPort(); err != nil {
		return err
	}
	return s.validateNginxConfSyntax()
}

func (s *Supplier) validateNginxConfExists() error {
	if exists, err := libbuildpack.FileExists(filepath.Join(s.Stager.BuildDir(), "nginx.conf")); err != nil {
		return err
	} else if !exists {
		s.Log.Error("nginx.conf file must be present at the app root")
		return errors.New("no nginx")
	}
	return nil
}

func (s *Supplier) validateNginxConfHasPort() error {
	conf, err := ioutil.ReadFile(filepath.Join(s.Stager.BuildDir(), "nginx.conf"))
	if err != nil {
		return err
	}
	if portFound, err := regexp.Match("{{.Port}}", conf); err != nil {
		return err
	} else if !portFound {
		s.Log.Error("nginx.conf file must be configured to respect the value of `{{.Port}}`")
		return errors.New("no .Port")
	}
	return nil
}

func (s *Supplier) validateNginxConfSyntax() error {

	tmpConfDir, err := ioutil.TempDir("/tmp", "conf")
	if err != nil {
		return fmt.Errorf("Error creating temp nginx conf dir: %s", err.Error())
	}
	defer os.RemoveAll(tmpConfDir)

	if err := libbuildpack.CopyDirectory(s.Stager.BuildDir(), tmpConfDir); err != nil {
		return fmt.Errorf("Error copying nginx.conf: %s", err.Error())
	}

	nginxConfPath := filepath.Join(tmpConfDir, "nginx.conf")
	cmd := exec.Command(filepath.Join(s.Stager.DepDir(), "bin", "varify"), nginxConfPath)
	cmd.Dir = tmpConfDir
	cmd.Stdout = ioutil.Discard
	cmd.Stderr = ioutil.Discard
	cmd.Env = append(os.Environ(), "PORT=8080", fmt.Sprintf("NGINX_MODULES=%s", filepath.Join(s.Stager.DepDir(), "nginx", "nginx", "modules")))
	if err := s.Command.Run(cmd); err != nil {
		return err
	}

	nginxExecDir := filepath.Join(s.Stager.DepDir(), "nginx", "nginx", "sbin")
	nginxErr := bytes.Buffer{}
	if err := s.Command.Execute(tmpConfDir, os.Stdout, &nginxErr, filepath.Join(nginxExecDir, "nginx"), "-t", "-c", nginxConfPath, "-p", tmpConfDir); err != nil {
		fmt.Fprint(os.Stderr, nginxErr.String())
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
	if val, ok := s.VersionLines[version]; ok {
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

func (s *Supplier) InstallNginx() error {
	dep, err := s.findMatchingVersion("nginx", s.Config.Version)
	if err != nil {
		s.Log.Info(`Available versions: ` + strings.Join(s.availableVersions(), ", "))
		return fmt.Errorf("Could not determine version: %s", err)
	}
	if s.Config.Version == "" {
		s.Log.BeginStep("No nginx version specified - using mainline => %s", dep.Version)
	} else {
		s.Log.BeginStep("Requested nginx version: %s => %s", s.Config.Version, dep.Version)
	}

	dir := filepath.Join(s.Stager.DepDir(), "nginx")

	if s.isStableLine(dep.Version) {
		s.Log.Warning(`Warning: usage of "stable" versions of NGINX is discouraged in most cases by the NGINX team.`)
	}

	if err := s.Manifest.InstallDependency(dep, dir); err != nil {
		return err
	}

	return s.Stager.AddBinDependencyLink(filepath.Join(dir, "nginx", "sbin", "nginx"), "nginx")
}
