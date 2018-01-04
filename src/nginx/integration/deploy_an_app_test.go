package integration_test

import (
	"path/filepath"

	"github.com/cloudfoundry/libbuildpack/cutlass"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("CF Nginx Buildpack", func() {
	var app *cutlass.App

	AfterEach(func() {
		if app != nil {
			app.Destroy()
		}
		app = nil
	})

	Context("with a simple nginx app", func() {

		BeforeEach(func() {
			app = cutlass.New(filepath.Join(bpDir, "fixtures", "simple"))
		})

		It("Logs nginx buildpack version", func() {
			Expect(app.Push()).ToNot(Succeed())
			Expect(app.ConfirmBuildpack(buildpackVersion)).To(Succeed())
		})
	})
})
