package supply_test

import (
	"io/ioutil"
	"nginx/supply"
	"os"

	"bytes"

	"github.com/cloudfoundry/libbuildpack"
	"github.com/golang/mock/gomock"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

//go:generate mockgen -source=supply.go --destination=mocks_test.go --package=supply_test

var _ = Describe("Supply", func() {
	var (
		depDir       string
		supplier     *supply.Supplier
		logger       *libbuildpack.Logger
		mockCtrl     *gomock.Controller
		mockStager   *MockStager
		mockManifest *MockManifest
		buffer       *bytes.Buffer
	)

	BeforeEach(func() {
		var err error
		buffer = new(bytes.Buffer)
		logger = libbuildpack.NewLogger(buffer)

		mockCtrl = gomock.NewController(GinkgoT())
		mockStager = NewMockStager(mockCtrl)
		mockManifest = NewMockManifest(mockCtrl)
		depDir, err = ioutil.TempDir("", "nginx.depdir")
		Expect(err).ToNot(HaveOccurred())
		mockStager.EXPECT().DepDir().AnyTimes().Return(depDir)
		supplier = supply.New(mockStager, mockManifest, logger)
	})

	AfterEach(func() {
		mockCtrl.Finish()
		os.RemoveAll(depDir)
	})

	Describe("InstallNginx", func() {
		BeforeEach(func() {
			supplier.VersionLines = map[string]string{"": "1.13.x", "mainline": "1.13.x", "stable": "1.12.x"}
			mockManifest.EXPECT().AllDependencyVersions("nginx").Return([]string{"1.12.2", "1.12.3", "1.13.8"}).AnyTimes()
			mockStager.EXPECT().AddBinDependencyLink(gomock.Any(), gomock.Any()).AnyTimes()
		})
		Context("request unavailable version", func() {
			BeforeEach(func() {
				supplier.Config.Version = "1.1.1"
			})
			It("Logs available versions and returns an error", func() {
				Expect(supplier.InstallNginx()).ToNot(Succeed())
				Expect(buffer.String()).To(ContainSubstring(`Available versions: mainline, stable, 1.12.x, 1.13.x, 1.12.2, 1.12.3, 1.13.8`))
			})
		})
		Context("request mainline version", func() {
			BeforeEach(func() {
				supplier.Config.Version = "mainline"
			})
			It("Logs the mainline version", func() {
				mockManifest.EXPECT().InstallDependency(libbuildpack.Dependency{Name: "nginx", Version: "1.13.8"}, gomock.Any())
				Expect(supplier.InstallNginx()).To(Succeed())
				Expect(buffer).To(ContainSubstring(`Requested nginx version: mainline => 1.13.8`))
			})
		})
		Context("request stable version", func() {
			BeforeEach(func() {
				supplier.Config.Version = "stable"
			})
			It("Logs the stable version", func() {
				mockManifest.EXPECT().InstallDependency(libbuildpack.Dependency{Name: "nginx", Version: "1.12.3"}, gomock.Any())
				Expect(supplier.InstallNginx()).To(Succeed())
				Expect(buffer).To(ContainSubstring(`Requested nginx version: stable => 1.12.3`))
			})
		})
		Context("request unspecified version", func() {
			BeforeEach(func() {
				supplier.Config.Version = ""
			})
			It("Logs the mainline version", func() {
				mockManifest.EXPECT().InstallDependency(libbuildpack.Dependency{Name: "nginx", Version: "1.13.8"}, gomock.Any())
				Expect(supplier.InstallNginx()).To(Succeed())
				Expect(buffer).To(ContainSubstring("No nginx version specified - using mainline => 1.13.8"))
			})
		})
		Context("request semver version", func() {
			BeforeEach(func() {
				supplier.Config.Version = "1.12.x"
			})
			It("Logs the semver request and the matching version", func() {
				mockManifest.EXPECT().InstallDependency(libbuildpack.Dependency{Name: "nginx", Version: "1.12.3"}, gomock.Any())
				Expect(supplier.InstallNginx()).To(Succeed())
				Expect(buffer).To(ContainSubstring(`Requested nginx version: 1.12.x => 1.12.3`))
			})
		})
		Context("request specific version", func() {
			BeforeEach(func() {
				supplier.Config.Version = "1.12.2"
			})
			It("Logs the specific version", func() {
				mockManifest.EXPECT().InstallDependency(libbuildpack.Dependency{Name: "nginx", Version: "1.12.2"}, gomock.Any())
				Expect(supplier.InstallNginx()).To(Succeed())
				Expect(buffer).To(ContainSubstring(`Requested nginx version: 1.12.2 => 1.12.2`))
			})
		})

		Describe("warns if 'stable' line is chosen", func() {
			const warning = `Warning: usage of "stable" versions of NGINX is discouraged in most cases by the NGINX team.`

			BeforeEach(func() {
				mockManifest.EXPECT().InstallDependency(gomock.Any(), gomock.Any())
			})

			It("stable emits warning", func() {
				supplier.Config.Version = "stable"
				Expect(supplier.InstallNginx()).To(Succeed())
				Expect(buffer).To(ContainSubstring(warning))
			})

			It("mainline does not warn", func() {
				supplier.Config.Version = "mainline"
				Expect(supplier.InstallNginx()).To(Succeed())
				Expect(buffer).ToNot(ContainSubstring(warning))
			})

			It("1.13.x does not warn", func() {
				supplier.Config.Version = "1.13.x"
				Expect(supplier.InstallNginx()).To(Succeed())
				Expect(buffer).ToNot(ContainSubstring(warning))
			})

			It("1.12.x emits warning", func() {
				supplier.Config.Version = "stable"
				Expect(supplier.InstallNginx()).To(Succeed())
				Expect(buffer).To(ContainSubstring(warning))
			})

			It("1.12.2 emits warning", func() {
				supplier.Config.Version = "stable"
				Expect(supplier.InstallNginx()).To(Succeed())
				Expect(buffer).To(ContainSubstring(warning))
			})
		})
	})
})
