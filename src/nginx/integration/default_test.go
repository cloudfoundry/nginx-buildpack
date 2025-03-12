package integration_test

import (
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/cloudfoundry/switchblade"
	"github.com/sclevine/spec"

	. "github.com/cloudfoundry/switchblade/matchers"
	. "github.com/onsi/gomega"
)

func testDefault(platform switchblade.Platform, fixtures string) func(*testing.T, spec.G, spec.S) {
	return func(t *testing.T, context spec.G, it spec.S) {
		var (
			Expect     = NewWithT(t).Expect
			Eventually = NewWithT(t).Eventually

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

		context("templated with env vars", func() {
			it("builds and runs the app", func() {
				deployment, logs, err := platform.Deploy.
					WithEnv(map[string]string{
						"OVERRIDE": `'{ "abcd": 12345 }{ \'ef\' : "ab" }'`,
					}).
					Execute(name, filepath.Join(fixtures, "default", "templated_env_vars"))
				Expect(err).NotTo(HaveOccurred())

				Expect(logs).To(ContainLines(MatchRegexp(`Installing nginx [\d\.]+`)), logs.String())

				Eventually(deployment).Should(Serve(ContainSubstring(`{ "abcd": 12345 }{ 'ef' : "ab" }`)).WithEndpoint("test"))

				cmd := exec.Command("docker", "container", "logs", deployment.Name)

				output, err := cmd.CombinedOutput()
				Expect(err).NotTo(HaveOccurred())

				Expect(string(output)).To(ContainSubstring(`NginxLog "GET /test HTTP/1.1" 200`))
			})
		})

		context("templated with env vars with include", func() {
			it("builds and runs the app", func() {
				deployment, logs, err := platform.Deploy.
					WithEnv(map[string]string{
						"OVERRIDE": `'{ "abcd": 12345 }{ \'ef\' : "ab" }'`,
					}).
					Execute(name, filepath.Join(fixtures, "default", "templated_with_include"))
				Expect(err).NotTo(HaveOccurred())

				Expect(logs).To(ContainLines(MatchRegexp(`Installing nginx [\d\.]+`)), logs.String())

				Eventually(deployment).Should(Serve(ContainSubstring(`{ "abcd": 12345 }{ 'ef' : "ab" }`)).WithEndpoint("test"))

				cmd := exec.Command("docker", "container", "logs", deployment.Name)

				output, err := cmd.CombinedOutput()
				Expect(err).NotTo(HaveOccurred())

				Expect(string(output)).To(ContainSubstring(`NginxLog "GET /test HTTP/1.1" 200`))
			})
		})

		context("with no specified pid", func() {
			it("builds and runs the app", func() {
				deployment, _, err := platform.Deploy.
					Execute(name, filepath.Join(fixtures, "default", "without_pid"))
				Expect(err).NotTo(HaveOccurred())

				Eventually(deployment).Should(Serve(ContainSubstring("Exciting Content")))

				cmd := exec.Command("docker", "container", "logs", deployment.Name)

				output, err := cmd.CombinedOutput()
				Expect(err).NotTo(HaveOccurred())

				Expect(string(output)).To(ContainSubstring(`NginxLog "GET / HTTP/1.1" 200`))
			})
		})

		context("with a specified pid", func() {
			it("builds and runs the app", func() {
				deployment, _, err := platform.Deploy.
					Execute(name, filepath.Join(fixtures, "default", "with_pid"))
				Expect(err).NotTo(HaveOccurred())

				Eventually(deployment).Should(Serve(ContainSubstring("Exciting Content")))

				cmd := exec.Command("docker", "container", "logs", deployment.Name)

				output, err := cmd.CombinedOutput()
				Expect(err).NotTo(HaveOccurred())

				Expect(string(output)).To(ContainSubstring(`NginxLog "GET / HTTP/1.1" 200`))
			})
		})

		context("with no specified version", func() {
			it("builds and runs the app and uses mainline", func() {
				deployment, logs, err := platform.Deploy.
					Execute(name, filepath.Join(fixtures, "default", "unspecified_version"))
				Expect(err).NotTo(HaveOccurred())

				Eventually(logs).Should(ContainSubstring(`No nginx version specified - using mainline => 1.27.`))
				Eventually(logs).ShouldNot(ContainSubstring(`Requested nginx version:`))

				Eventually(deployment).Should(Serve(ContainSubstring("Exciting Content")))

				cmd := exec.Command("docker", "container", "logs", deployment.Name)

				output, err := cmd.CombinedOutput()
				Expect(err).NotTo(HaveOccurred())

				Expect(string(output)).To(ContainSubstring(`NginxLog "GET / HTTP/1.1" 200`))
			})
		})

		context("with an app specifying mainline", func() {
			it("builds and runs the app", func() {
				deployment, logs, err := platform.Deploy.
					Execute(name, filepath.Join(fixtures, "default", "mainline"))
				Expect(err).NotTo(HaveOccurred())

				Eventually(logs).Should(ContainSubstring(`Requested nginx version: mainline => 1.27.`))

				Eventually(deployment).Should(Serve(ContainSubstring("Exciting Content")))

				cmd := exec.Command("docker", "container", "logs", deployment.Name)

				output, err := cmd.CombinedOutput()
				Expect(err).NotTo(HaveOccurred())

				Expect(string(output)).To(ContainSubstring(`NginxLog "GET / HTTP/1.1" 200`))
			})
		})

		context("with an app specifying stable", func() {
			it("builds and runs the app", func() {
				deployment, logs, err := platform.Deploy.
					Execute(name, filepath.Join(fixtures, "default", "stable"))
				Expect(err).NotTo(HaveOccurred())

				Eventually(logs).Should(ContainSubstring(`Requested nginx version: stable => 1.26.`))
				Eventually(logs).Should(ContainSubstring(`Warning: usage of "stable" versions of NGINX is discouraged in most cases by the NGINX team.`))

				Eventually(deployment).Should(Serve(ContainSubstring("Exciting Content")))

				cmd := exec.Command("docker", "container", "logs", deployment.Name)

				output, err := cmd.CombinedOutput()
				Expect(err).NotTo(HaveOccurred())

				Expect(string(output)).To(ContainSubstring(`NginxLog "GET / HTTP/1.1" 200`))
			})
		})

		context("with an app unavailable version", func() {
			it("the build fails and logs and error", func() {
				_, logs, err := platform.Deploy.
					Execute(name, filepath.Join(fixtures, "default", "unavailable_version"))
				Expect(err).To(HaveOccurred())

				Eventually(logs).Should(ContainSubstring(`Available versions: mainline, stable, 1.26.x, 1.27.x`))
			})
		})

		context("with using the stream module", func() {
			it("builds and runs the app", func() {
				deployment, _, err := platform.Deploy.
					Execute(name, filepath.Join(fixtures, "default", "with_stream_module"))
				Expect(err).NotTo(HaveOccurred())

				Eventually(deployment).Should(Serve(ContainSubstring("Exciting Content")))
			})
		})

		context("with an app that has no access to logging", func() {
			it("logs a warning and builds and runs the app", func() {
				deployment, logs, err := platform.Deploy.
					Execute(name, filepath.Join(fixtures, "default", "no_logging"))
				Expect(err).NotTo(HaveOccurred())

				Eventually(logs).Should(ContainSubstring(`Warning: access logging is turned off in your nginx.conf file, this may make your app difficult to debug.`))

				Eventually(deployment).Should(Serve(ContainSubstring("Exciting Content")))
			})
		})

		context("an Openresty app", func() {
			it("builds and runs the app", func() {
				deployment, _, err := platform.Deploy.
					Execute(name, filepath.Join(fixtures, "default", "openresty"))
				Expect(err).NotTo(HaveOccurred())

				Eventually(deployment).Should(Serve(ContainSubstring("<p>hello, world</p>")))

				cmd := exec.Command("docker", "container", "logs", deployment.Name)

				output, err := cmd.CombinedOutput()
				Expect(err).NotTo(HaveOccurred())

				Expect(string(output)).To(ContainSubstring(`NginxLog "GET / HTTP/1.1" 200`))
			})
		})
	}
}
