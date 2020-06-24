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

	Context("with templated json return value", func() {
		BeforeEach(func() {
			app = cutlass.New(filepath.Join(bpDir, "fixtures", "templated_env_vars"))
		})

		It("Deploys successfully", func() {
			env := `'{ "abcd": 12345 }{ \'ef\' : "ab" }'`
			app.SetEnv("OVERRIDE", env)
			PushAppAndConfirm(app)

			Expect(app.GetBody("/test")).To(ContainSubstring(`{ "abcd": 12345 }{ 'ef' : "ab" }`))
			Eventually(app.Stdout.String).Should(ContainSubstring(`NginxLog "GET /test HTTP/1.1" 200`))
		})
	})

	Context("with no specified pid", func() {
		BeforeEach(func() {
			app = cutlass.New(filepath.Join(bpDir, "fixtures", "without_pid"))
		})

		It("Deploys successfully", func() {
			PushAppAndConfirm(app)

			Expect(app.GetBody("/")).To(ContainSubstring("Exciting Content"))
			Eventually(app.Stdout.String).Should(ContainSubstring(`NginxLog "GET / HTTP/1.1" 200`))
		})
	})

	Context("with a specified pid", func() {
		BeforeEach(func() {
			app = cutlass.New(filepath.Join(bpDir, "fixtures", "with_pid"))
		})

		It("Deploys successfully", func() {
			PushAppAndConfirm(app)

			Expect(app.GetBody("/")).To(ContainSubstring("Exciting Content"))
			Eventually(app.Stdout.String).Should(ContainSubstring(`NginxLog "GET / HTTP/1.1" 200`))
		})
	})

	Context("with no specified version", func() {
		BeforeEach(func() {
			app = cutlass.New(filepath.Join(bpDir, "fixtures", "unspecified_version"))
		})

		It("Uses latest mainline nginx", func() {
			PushAppAndConfirm(app)

			Eventually(app.Stdout.String).Should(ContainSubstring(`No nginx version specified - using mainline => 1.19.`))
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

			Eventually(app.Stdout.String).Should(ContainSubstring(`Requested nginx version: mainline => 1.19.`))

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

			Eventually(app.Stdout.String).Should(ContainSubstring(`Requested nginx version: stable => 1.18.`))
			Eventually(app.Stdout.String).Should(ContainSubstring(`Warning: usage of "stable" versions of NGINX is discouraged in most cases by the NGINX team.`))

			Expect(app.GetBody("/")).To(ContainSubstring("Exciting Content"))
			Eventually(app.Stdout.String).Should(ContainSubstring(`NginxLog "GET / HTTP/1.1" 200`))
		})
	})

	Context("with an nginx app specifying an unknown version", func() {
		BeforeEach(func() {
			app = cutlass.New(filepath.Join(bpDir, "fixtures", "unavailable_version"))
		})

		It("Logs nginx buildpack versions", func() {
			Expect(app.Push()).ToNot(Succeed())

			Eventually(app.Stdout.String).Should(ContainSubstring(`Available versions: mainline, stable, 1.18.x, 1.19.x`))
		})
	})

	Context("with an nginx app that uses the stream module", func() {
		BeforeEach(func() {
			app = cutlass.New(filepath.Join(bpDir, "fixtures", "with_stream_module"))
		})

		It("Pushes the app successfully", func() {
			PushAppAndConfirm(app)
		})
	})

	Context("an app without access logging", func() {
		const warning = `Warning: access logging is turned off in your nginx.conf file, this may make your app difficult to debug.`
		BeforeEach(func() {
			app = cutlass.New(filepath.Join(bpDir, "fixtures", "no_logging"))
			app.Buildpacks = []string{"nginx_buildpack"}
		})
		AfterEach(func() {
			app.Destroy()
		})

		It("Logs a warning", func() {
			Expect(app.Push()).To(Succeed())

			Eventually(app.Stdout.String).Should(ContainSubstring(warning))
		})
	})

	Context("an OpenResty app", func() {
		BeforeEach(func() {
			app = cutlass.New(filepath.Join(bpDir, "fixtures", "openresty"))
		})

		It("Deploys successfully", func() {
			PushAppAndConfirm(app)

			Expect(app.GetBody("/")).To(ContainSubstring("<p>hello, world</p>"))
			Eventually(app.Stdout.String).Should(ContainSubstring(`NginxLog "GET / HTTP/1.1" 200`))
		})
	})
})
