package finalize

import (
	"errors"
	"io/ioutil"
	"path/filepath"
	"regexp"

	"github.com/cloudfoundry/libbuildpack"
)

type Finalizer struct {
	BuildDir string
	DepDir   string
	Log      *libbuildpack.Logger
}

func Run(sf *Finalizer) error {

	if exists, err := libbuildpack.FileExists(filepath.Join(sf.BuildDir, "nginx.conf")); err != nil {
		return err
	} else if !exists {
		sf.Log.Error("nginx.conf file must be present at the app root")
		return errors.New("no nginx")
	}

	conf, err := ioutil.ReadFile(filepath.Join(sf.BuildDir, "nginx.conf"))
	if err != nil {
		return err
	}
	if portFound, err := regexp.Match("{{.Port}}", conf); err != nil {
		return err
	} else if !portFound {
		sf.Log.Error("nginx.conf file must be configured to respect the value of `{{.Port}}`")
		return errors.New("no .Port")
	}
	return nil
}
