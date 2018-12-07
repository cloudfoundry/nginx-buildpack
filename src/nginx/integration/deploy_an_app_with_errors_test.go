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

	Context("an app without nginx.conf", func() {
		BeforeEach(func() {
			app = cutlass.New(filepath.Join(bpDir, "fixtures", "empty"))
			app.Buildpacks = []string{"nginx_buildpack"}
		})

		It("Logs an error", func() {
			Expect(app.Push()).ToNot(Succeed())
			Expect(app.ConfirmBuildpack(buildpackVersion)).To(Succeed())

			Eventually(app.Stdout.String).Should(ContainSubstring("Could not validate nginx.conf"))
		})
	})

	Context("an app with an invalid nginx.conf", func() {
		BeforeEach(func() {
			app = cutlass.New(filepath.Join(bpDir, "fixtures", "invalid_conf"))
			app.Buildpacks = []string{"nginx_buildpack"}
		})

		It("Logs an error", func() {
			Expect(app.Push()).ToNot(Succeed())
			Expect(app.ConfirmBuildpack(buildpackVersion)).To(Succeed())

			Eventually(app.Stdout.String).Should(ContainSubstring("nginx.conf contains syntax errors"))
		})
	})

	Context("an app with nginx.conf without {{port}}", func() {
		BeforeEach(func() {
			app = cutlass.New(filepath.Join(bpDir, "fixtures", "missing_template_port"))
			app.Buildpacks = []string{"nginx_buildpack"}
		})

		It("Logs an error", func() {
			Expect(app.Push()).ToNot(Succeed())
			Expect(app.ConfirmBuildpack(buildpackVersion)).To(Succeed())

			Eventually(app.Stdout.String).Should(ContainSubstring("nginx.conf file must be configured to respect the value of `{{port}}`"))
		})
	})
})
