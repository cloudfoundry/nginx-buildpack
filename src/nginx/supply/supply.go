package supply

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/cloudfoundry/libbuildpack"
)

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
	BuildDir() string
}

type Config struct {
	Version string `yaml:"version"`
}

type Supplier struct {
	Stager       Stager
	Manifest     Manifest
	Log          *libbuildpack.Logger
	Config       Config
	VersionLines map[string]string
}

func New(stager Stager, manifest Manifest, logger *libbuildpack.Logger) *Supplier {
	return &Supplier{
		Stager:   stager,
		Manifest: manifest,
		Log:      logger,
	}
}

func (s *Supplier) Run() error {
	s.Log.BeginStep("Supplying nginx")

	if err := s.InstallVarify(); err != nil {
		s.Log.Error("Failed to copy verify: %s", err)
		return err
	}
	if err := s.Setup(); err != nil {
		s.Log.Error("Could not setup: %s", err)
		return err
	}

	if err := s.InstallNginx(); err != nil {
		s.Log.Error("Could not install nginx: %s", err)
		return err
	}

	return nil
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
