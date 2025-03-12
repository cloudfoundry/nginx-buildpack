package integration_test

import (
	"path/filepath"
	"testing"

	"github.com/cloudfoundry/switchblade"
	"github.com/sclevine/spec"

	. "github.com/onsi/gomega"
)

func testErrors(platform switchblade.Platform, fixtures string) func(*testing.T, spec.G, spec.S) {
	return func(t *testing.T, context spec.G, it spec.S) {
		var (
			Expect = NewWithT(t).Expect

			name string
		)

		it.Before(func() {
			var err error
			name, err = switchblade.RandomName()
			Expect(err).NotTo(HaveOccurred())
		})

		it.After(func() {
			Expect(platform.Delete.Execute(name)).To(Succeed())
		})

		context("an app without an nginx.conf", func() {
			it("logs an error", func() {
				_, logs, err := platform.Deploy.
					WithBuildpacks("nginx_buildpack").
					Execute(name, filepath.Join(fixtures, "errors", "empty"))
				Expect(err).To(HaveOccurred())

				Expect(logs).To(ContainSubstring("Could not validate nginx.conf"), logs.String())
			})
		})

		context("an app with an invalid nginx.conf", func() {
			it("logs an error", func() {
				_, logs, err := platform.Deploy.
					WithBuildpacks("nginx_buildpack").
					Execute(name, filepath.Join(fixtures, "errors", "invalid_conf"))
				Expect(err).To(HaveOccurred())

				Expect(logs).To(ContainSubstring("nginx.conf contains syntax errors"), logs.String())
			})
		})

		context("an app with nginx.conf without {{port}}", func() {
			it("logs an error", func() {
				_, logs, err := platform.Deploy.
					WithBuildpacks("nginx_buildpack").
					Execute(name, filepath.Join(fixtures, "errors", "missing_template_port"))
				Expect(err).To(HaveOccurred())

				Expect(logs).To(ContainSubstring("The listen port value in nginx.conf must be configured to the template `{{port}}`"), logs.String())
			})
		})
	}
}
