package supply

import (
	"path/filepath"

	"github.com/cloudfoundry/libbuildpack"
)

type Manifest interface {
	InstallOnlyVersion(string, string) error
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
	if err := s.Manifest.InstallOnlyVersion("nginx", filepath.Join(s.Stager.DepDir(), "nginx")); err != nil {
		return err
	}

	return s.Stager.AddBinDependencyLink(filepath.Join(s.Stager.DepDir(), "nginx", "nginx", "sbin", "nginx"), "nginx")
}
