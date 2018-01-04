package supply

import (
	"github.com/cloudfoundry/libbuildpack"
)

type Stager interface {
	DepDir() string
}

type Supplier struct {
	Stager Stager
	Log    *libbuildpack.Logger
}

func New(stager Stager, logger *libbuildpack.Logger) *Supplier {
	return &Supplier{
		Stager: stager,
		Log:    logger,
	}
}

func (s *Supplier) Run() error {
	return nil
}
