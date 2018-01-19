package supply

import (
	"path/filepath"

	"github.com/cloudfoundry/libbuildpack"
)

type Manifest interface {
	InstallOnlyVersion(string, string) error
	RootDir() string
}
type Stager interface {
	AddBinDependencyLink(string, string) error
	DepDir() string
}

type Supplier struct {
	Stager   Stager
	Manifest Manifest
	Log      *libbuildpack.Logger
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

func (s *Supplier) InstallNginx() error {
	if err := s.Manifest.InstallOnlyVersion("nginx", filepath.Join(s.Stager.DepDir(), "nginx")); err != nil {
		return err
	}

	return s.Stager.AddBinDependencyLink(filepath.Join(s.Stager.DepDir(), "nginx", "nginx", "sbin", "nginx"), "nginx")
}
