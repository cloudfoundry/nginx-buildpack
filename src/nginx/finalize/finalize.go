package finalize

import (
	"github.com/cloudfoundry/libbuildpack"
)

type Finalizer struct {
	BuildDir string
	DepDir   string
	Log      *libbuildpack.Logger
}

func Run(sf *Finalizer) error {
	return nil
}
