package brats_test

import (
	"github.com/cloudfoundry/libbuildpack/bratshelper"
	"github.com/cloudfoundry/libbuildpack/cutlass"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Nginx buildpack", func() {
	bratshelper.UnbuiltBuildpack("nginx", CopyBrats)
	bratshelper.DeployingAnAppWithAnUpdatedVersionOfTheSameBuildpack(CopyBrats)
	bratshelper.DeployAppWithExecutableProfileScript("nginx", CopyBrats)
	bratshelper.DeployAnAppWithSensitiveEnvironmentVariables(CopyBrats)
	bratshelper.ForAllSupportedVersions("nginx", CopyBrats, func(nginxVersion string, app *cutlass.App) {
		bratshelper.PushApp(app)

		By("installs the correct version of Nginx", func() {
			Expect(app.Stdout.String()).To(ContainSubstring("Installing nginx " + nginxVersion))
		})
		By("runs a simple webserver", func() {
			Expect(app.GetBody("/")).To(ContainSubstring("Hello World!"))
		})
	})
})
