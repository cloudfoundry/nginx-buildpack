package integration_test

import (
	"path/filepath"

	"github.com/cloudfoundry/libbuildpack/cutlass"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("pushing an app a second time", func() {
	var app *cutlass.App
	AfterEach(func() {
		if app != nil {
			app.Destroy()
		}
		app = nil
	})

	BeforeEach(func() {
		if cutlass.Cached {
			Skip("running uncached tests")
		}

		app = cutlass.New(filepath.Join(bpDir, "fixtures", "mainline"))
		app.Buildpacks = []string{"nginx_buildpack"}
	})

	Regexp := `\[.*/nginx\_[\d\.]+\_linux\_x64\_(cflinuxfs.*_)?[\da-f]+\.tgz\]`
	DownloadRegexp := "Download " + Regexp
	CopyRegexp := "Copy " + Regexp

	It("uses the cache for manifest dependencies", func() {
		PushAppAndConfirm(app)
		Eventually(app.Stdout.String).Should(MatchRegexp(DownloadRegexp))
		Expect(app.Stdout.String()).ToNot(MatchRegexp(CopyRegexp))

		app.Stdout.Reset()
		PushAppAndConfirm(app)
		Eventually(app.Stdout.String).Should(MatchRegexp(CopyRegexp))
		Expect(app.Stdout.String()).ToNot(MatchRegexp(DownloadRegexp))
	})
})
