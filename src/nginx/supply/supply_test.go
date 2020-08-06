package supply_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/cloudfoundry/libbuildpack"
	"github.com/cloudfoundry/nginx-buildpack/src/nginx/supply"
	"github.com/golang/mock/gomock"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

//go:generate mockgen -source=supply.go --destination=mocks_test.go --package=supply_test

var _ = Describe("Supply", func() {
	var (
		depDir        string
		supplier      *supply.Supplier
		logger        *libbuildpack.Logger
		mockCtrl      *gomock.Controller
		mockStager    *MockStager
		mockManifest  *MockManifest
		mockInstaller *MockInstaller
		mockCommand   *MockCommand
		buffer        *bytes.Buffer
	)

	BeforeEach(func() {
		var err error
		buffer = new(bytes.Buffer)
		logger = libbuildpack.NewLogger(buffer)

		mockCtrl = gomock.NewController(GinkgoT())
		mockStager = NewMockStager(mockCtrl)
		mockManifest = NewMockManifest(mockCtrl)
		mockInstaller = NewMockInstaller(mockCtrl)
		mockCommand = NewMockCommand(mockCtrl)
		depDir, err = ioutil.TempDir("", "nginx.depdir")
		Expect(err).ToNot(HaveOccurred())
		mockStager.EXPECT().DepDir().AnyTimes().Return(depDir)
		supplier = supply.New(mockStager, mockManifest, mockInstaller, logger, mockCommand)
	})

	AfterEach(func() {
		mockCtrl.Finish()
		os.RemoveAll(depDir)
	})

	Describe("InstallNGINX", func() {
		BeforeEach(func() {
			supplier.VersionLines = map[string]string{"": "1.13.x", "mainline": "1.13.x", "stable": "1.12.x"}
			mockManifest.EXPECT().AllDependencyVersions("nginx").Return([]string{"1.12.2", "1.12.3", "1.13.8"}).AnyTimes()
			mockStager.EXPECT().AddBinDependencyLink(gomock.Any(), gomock.Any()).AnyTimes()
		})

		Context("request unavailable version", func() {
			BeforeEach(func() {
				supplier.Config.Nginx.Version = "1.1.1"
			})

			It("Logs available versions and returns an error", func() {
				Expect(supplier.InstallNGINX()).ToNot(Succeed())
				Expect(buffer.String()).To(ContainSubstring(`Available versions: mainline, stable, 1.12.x, 1.13.x, 1.12.2, 1.12.3, 1.13.8`))
			})
		})

		Context("request mainline version", func() {
			BeforeEach(func() {
				supplier.Config.Nginx.Version = "mainline"
			})

			It("Logs the mainline version", func() {
				mockInstaller.EXPECT().InstallDependency(libbuildpack.Dependency{Name: "nginx", Version: "1.13.8"}, gomock.Any())
				Expect(supplier.InstallNGINX()).To(Succeed())
				Expect(buffer).To(ContainSubstring(`Requested nginx version: mainline => 1.13.8`))
			})
		})

		Context("request stable version", func() {
			BeforeEach(func() {
				supplier.Config.Nginx.Version = "stable"
			})

			It("Logs the stable version", func() {
				mockInstaller.EXPECT().InstallDependency(libbuildpack.Dependency{Name: "nginx", Version: "1.12.3"}, gomock.Any())
				Expect(supplier.InstallNGINX()).To(Succeed())
				Expect(buffer).To(ContainSubstring(`Requested nginx version: stable => 1.12.3`))
			})
		})

		Context("request unspecified version", func() {
			BeforeEach(func() {
				supplier.Config.Nginx.Version = ""
			})

			It("Logs the mainline version", func() {
				mockInstaller.EXPECT().InstallDependency(libbuildpack.Dependency{Name: "nginx", Version: "1.13.8"}, gomock.Any())
				Expect(supplier.InstallNGINX()).To(Succeed())
				Expect(buffer).To(ContainSubstring("No nginx version specified - using mainline => 1.13.8"))
			})
		})

		Context("request semver version", func() {
			BeforeEach(func() {
				supplier.Config.Nginx.Version = "1.12.x"
			})

			It("Logs the semver request and the matching version", func() {
				mockInstaller.EXPECT().InstallDependency(libbuildpack.Dependency{Name: "nginx", Version: "1.12.3"}, gomock.Any())
				Expect(supplier.InstallNGINX()).To(Succeed())
				Expect(buffer).To(ContainSubstring(`Requested nginx version: 1.12.x => 1.12.3`))
			})
		})

		Context("request specific version", func() {
			BeforeEach(func() {
				supplier.Config.Nginx.Version = "1.12.2"
			})

			It("Logs the specific version", func() {
				mockInstaller.EXPECT().InstallDependency(libbuildpack.Dependency{Name: "nginx", Version: "1.12.2"}, gomock.Any())
				Expect(supplier.InstallNGINX()).To(Succeed())
				Expect(buffer).To(ContainSubstring(`Requested nginx version: 1.12.2 => 1.12.2`))
			})
		})

		Describe("warns if 'stable' line is chosen", func() {
			const warning = `Warning: usage of "stable" versions of NGINX is discouraged in most cases by the NGINX team.`

			BeforeEach(func() {
				mockInstaller.EXPECT().InstallDependency(gomock.Any(), gomock.Any())
			})

			It("stable emits warning", func() {
				supplier.Config.Nginx.Version = "stable"
				Expect(supplier.InstallNGINX()).To(Succeed())
				Expect(buffer).To(ContainSubstring(warning))
			})

			It("mainline does not warn", func() {
				supplier.Config.Nginx.Version = "mainline"
				Expect(supplier.InstallNGINX()).To(Succeed())
				Expect(buffer).ToNot(ContainSubstring(warning))
			})

			It("1.13.x does not warn", func() {
				supplier.Config.Nginx.Version = "1.13.x"
				Expect(supplier.InstallNGINX()).To(Succeed())
				Expect(buffer).ToNot(ContainSubstring(warning))
			})

			It("1.12.x emits warning", func() {
				supplier.Config.Nginx.Version = "stable"
				Expect(supplier.InstallNGINX()).To(Succeed())
				Expect(buffer).To(ContainSubstring(warning))
			})

			It("1.12.2 emits warning", func() {
				supplier.Config.Nginx.Version = "stable"
				Expect(supplier.InstallNGINX()).To(Succeed())
				Expect(buffer).To(ContainSubstring(warning))
			})
		})
	})

	Describe("InstallOpenResty", func() {
		It("installs the available version of openresty", func() {
			mockManifest.EXPECT().AllDependencyVersions("openresty").Return([]string{"1.13.6.2"}).AnyTimes()
			mockStager.EXPECT().AddBinDependencyLink(gomock.Any(), gomock.Any()).AnyTimes()
			mockInstaller.EXPECT().InstallDependency(libbuildpack.Dependency{Name: "openresty", Version: "1.13.6.2"}, gomock.Any())
			Expect(supplier.InstallOpenResty()).To(Succeed())
		})
	})

	Describe("WriteProfileD", func() {
		It("writes nginx script", func() {
			mockStager.EXPECT().DepsIdx().Return("0")
			mockStager.EXPECT().WriteProfileD("nginx", "export DEP_DIR=$DEPS_DIR/0\nmkdir -p logs")

			supplier.WriteProfileD()
		})

		It("writes openresty script", func() {
			mockStager.EXPECT().DepsIdx().Return("0").Times(3)
			mockStager.EXPECT().WriteProfileD("openresty", fmt.Sprintf(
				"%s%s",
				"export LD_LIBRARY_PATH=$LD_LIBRARY_PATH:$DEPS_DIR/0/nginx/luajit/lib\n",
				"export LUA_PATH=$DEPS_DIR/0/nginx/lualib/?.lua\n",
			))
			mockStager.EXPECT().WriteProfileD("nginx", "export DEP_DIR=$DEPS_DIR/0\nmkdir -p logs")

			supplier.Config.Dist = "openresty"
			supplier.WriteProfileD()
		})
	})

	Describe("ValidateNginxConf", func() {
		var (
			buildDir string
			err      error
		)

		BeforeEach(func() {
			buildDir, err = ioutil.TempDir("", "")
			Expect(err).NotTo(HaveOccurred())

			mockStager.EXPECT().BuildDir().Return(buildDir).AnyTimes()

			mockCommand.EXPECT().Run(gomock.Any()).AnyTimes()
			mockCommand.EXPECT().Execute(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
		})

		AfterEach(func() {
			os.RemoveAll(buildDir)
		})

		It("parses the port", func() {
			ioutil.WriteFile(filepath.Join(buildDir, "nginx.conf"), []byte("{{port}}"), 0666)
			Expect(supplier.ValidateNginxConf()).To(Succeed())
		})

		It("parses the port and ignores white spaces", func() {
			ioutil.WriteFile(filepath.Join(buildDir, "nginx.conf"), []byte("{{  port  }}"), 0666)
			Expect(supplier.ValidateNginxConf()).To(Succeed())
		})

		Context("CheckAccessLogging", func() {
			It("logs a warning when access logging is not set", func() {
				ioutil.WriteFile(filepath.Join(buildDir, "nginx.conf"), []byte("some content"), 0666)
				Expect(supplier.CheckAccessLogging()).To(Succeed())
				Expect(buffer.String()).To(ContainSubstring("Warning: access logging is turned off in your nginx.conf file, this may make your app difficult to debug."))
			})

			It("logs a warning when access logging is set to off", func() {
				ioutil.WriteFile(filepath.Join(buildDir, "nginx.conf"), []byte("access_log off"), 0666)
				Expect(supplier.CheckAccessLogging()).To(Succeed())
				Expect(buffer.String()).To(ContainSubstring("Warning: access logging is turned off in your nginx.conf file, this may make your app difficult to debug."))
			})

			It("logs a warning when access logging is set to off with extra spaces", func() {
				ioutil.WriteFile(filepath.Join(buildDir, "nginx.conf"), []byte("access_log    off"), 0666)
				Expect(supplier.CheckAccessLogging()).To(Succeed())
				Expect(buffer.String()).To(ContainSubstring("Warning: access logging is turned off in your nginx.conf file, this may make your app difficult to debug."))
			})

			It("logs a warning when access logging is set to OFF", func() {
				ioutil.WriteFile(filepath.Join(buildDir, "nginx.conf"), []byte("access_log OFF"), 0666)
				Expect(supplier.CheckAccessLogging()).To(Succeed())
				Expect(buffer.String()).To(ContainSubstring("Warning: access logging is turned off in your nginx.conf file, this may make your app difficult to debug."))
			})

			It("logs a warning when access logging is set to a path", func() {
				ioutil.WriteFile(filepath.Join(buildDir, "nginx.conf"), []byte("access_log /some/path"), 0666)
				Expect(supplier.CheckAccessLogging()).To(Succeed())
				Expect(buffer.String()).ToNot(ContainSubstring("Warning: access logging is turned off in your nginx.conf file, this may make your app difficult to debug."))
			})
		})

	})
})
