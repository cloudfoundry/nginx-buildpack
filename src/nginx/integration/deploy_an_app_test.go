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

	Context("with no specified version", func() {
		BeforeEach(func() {
			app = cutlass.New(filepath.Join(bpDir, "fixtures", "unspecified_version"))
		})

		It("Uses latest mainline nginx", func() {
			PushAppAndConfirm(app)

			Eventually(app.Stdout.String).Should(ContainSubstring(`No nginx version specified - using mainline => 1.13.`))
			Eventually(app.Stdout.String).ShouldNot(ContainSubstring(`Requested nginx version:`))

			Expect(app.GetBody("/")).To(ContainSubstring("Exciting Content"))
			Eventually(app.Stdout.String).Should(ContainSubstring(`NginxLog "GET / HTTP/1.1" 200`))
		})
	})

	Context("with an nginx app specifying mainline", func() {
		BeforeEach(func() {
			app = cutlass.New(filepath.Join(bpDir, "fixtures", "mainline"))
		})

		It("Logs nginx buildpack version", func() {
			PushAppAndConfirm(app)

			Eventually(app.Stdout.String).Should(ContainSubstring(`Requested nginx version: mainline => 1.13.`))

			Expect(app.GetBody("/")).To(ContainSubstring("Exciting Content"))
			Eventually(app.Stdout.String).Should(ContainSubstring(`NginxLog "GET / HTTP/1.1" 200`))
		})
	})

	Context("with an nginx app specifying stable", func() {
		BeforeEach(func() {
			app = cutlass.New(filepath.Join(bpDir, "fixtures", "stable"))
		})

		It("Logs nginx buildpack version", func() {
			PushAppAndConfirm(app)

			Eventually(app.Stdout.String).Should(ContainSubstring(`Requested nginx version: stable => 1.12.`))
			Eventually(app.Stdout.String).Should(ContainSubstring(`Warning: usage of "stable" versions of NGINX is discouraged in most cases by the NGINX team.`))

			Expect(app.GetBody("/")).To(ContainSubstring("Exciting Content"))
			Eventually(app.Stdout.String).Should(ContainSubstring(`NginxLog "GET / HTTP/1.1" 200`))
		})
	})

	Context("with an nginx app specifying 1.12.x", func() {
		BeforeEach(func() {
			app = cutlass.New(filepath.Join(bpDir, "fixtures", "1_12_x"))
		})

		It("Logs nginx buildpack version", func() {
			PushAppAndConfirm(app)

			Eventually(app.Stdout.String).Should(ContainSubstring(`Requested nginx version: 1.12.x => 1.12.`))
			Eventually(app.Stdout.String).Should(ContainSubstring(`Warning: usage of "stable" versions of NGINX is discouraged in most cases by the NGINX team.`))

			Expect(app.GetBody("/")).To(ContainSubstring("Exciting Content"))
			Eventually(app.Stdout.String).Should(ContainSubstring(`NginxLog "GET / HTTP/1.1" 200`))
		})
	})

	Context("with an nginx app specifying an unknown version", func() {
		BeforeEach(func() {
			app = cutlass.New(filepath.Join(bpDir, "fixtures", "unavailable_version"))
		})

		It("Logs nginx buildpack version", func() {
			Expect(app.Push()).ToNot(Succeed())

			Eventually(app.Stdout.String).Should(ContainSubstring(`Available versions: mainline, stable, 1.12.x, 1.13.x`))
		})
	})
})
