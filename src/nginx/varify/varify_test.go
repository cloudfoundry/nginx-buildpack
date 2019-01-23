package main_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("varify", func() {
	var (
		tmpDir, localModulePath, globalModulePath string
	)

	BeforeEach(func() {
		var err error
		tmpDir, err = ioutil.TempDir("", "nginx.tmpdir")
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		os.RemoveAll(tmpDir)
	})

	Describe("Run", func() {
		It("templates current port using the 'port' func", func() {
			body := runCli(tmpDir, "Hi the port is {{port}}.", []string{"PORT=8080"}, "", "", "")
			Expect(body).To(Equal("Hi the port is 8080."))
		})

		It("templates environment variables using the 'env' func", func() {
			body := runCli(tmpDir, `The env var FOO is {{env "FOO"}}`, []string{"FOO=BAR"}, "", "", "")
			Expect(body).To(Equal("The env var FOO is BAR"))
		})

		Context("templating a load_module directive using the 'module' func", func() {
			BeforeEach(func() {
				localModulePath = filepath.Join(tmpDir, "local_modules")
				globalModulePath = filepath.Join(tmpDir, "global_modules")

				Expect(os.Mkdir(localModulePath, 0744)).To(Succeed())
				Expect(os.Mkdir(globalModulePath, 0744)).To(Succeed())
				Expect(ioutil.WriteFile(filepath.Join(localModulePath, "local.so"), []byte("dummy data"), 0644)).To(Succeed())
				Expect(ioutil.WriteFile(filepath.Join(globalModulePath, "global.so"), []byte("dummy data"), 0644)).To(Succeed())
			})

			Context("when the module is in local modules directory", func() {
				It("loads the module from the local directory", func() {
					body := runCli(tmpDir, `{{module "local"}}`, nil, localModulePath, globalModulePath, "")
					Expect(body).To(Equal(fmt.Sprintf("load_module %s/local.so;", localModulePath)))
				})
			})

			Context("when the module is in global modules directory", func() {
				It("loads the module from the global directory", func() {
					body := runCli(tmpDir, `{{module "global"}}`, nil, localModulePath, globalModulePath, "")
					Expect(body).To(Equal(fmt.Sprintf("load_module %s/global.so;", globalModulePath)))
				})
			})
		})

		Context("templating a nameserver directive using the 'nameserver' func", func() {
			var defaultNameserver = "169.254.0.2"
			var nameserver = "123.245.67.89"

			It("reads nameserver from the default resolv-conf file", func() {
				var resolvConfPath = filepath.Join(tmpDir, "resolv.conf")
				Expect(ioutil.WriteFile(resolvConfPath, []byte("nameserver "+nameserver), 0644)).To(Succeed())
				body := runCli(tmpDir, "Hi the nameserver is {{nameserver}}.", nil, "", "", resolvConfPath)
				Expect(body).To(Equal("Hi the nameserver is " + nameserver + "."))
			})

			It("reads nameserver from the unusual resolv-conf file", func() {
				var resolvConfPath = filepath.Join(tmpDir, "resolv-unusual.conf")
				Expect(ioutil.WriteFile(resolvConfPath, []byte("# comment 1\n  \t  nameserver "+nameserver+"  \t  \n# comment 2"), 0644)).To(Succeed())
				body := runCli(tmpDir, "Hi the nameserver is {{nameserver}}.", nil, "", "", resolvConfPath)
				Expect(body).To(Equal("Hi the nameserver is " + nameserver + "."))
			})

			It("set the default nameserver if the resolv-conf file is empty", func() {
				var resolvConfPath = filepath.Join(tmpDir, "resolv-empty.conf")
				Expect(ioutil.WriteFile(resolvConfPath, []byte(""), 0644)).To(Succeed())
				body := runCli(tmpDir, "Hi the nameserver is {{nameserver}}.", nil, "", "", resolvConfPath)
				Expect(body).To(Equal("Hi the nameserver is " + defaultNameserver + "."))
			})

			It("set the default nameserver if the resolv-conf file does't exist", func() {
				body := runCli(tmpDir, "Hi the nameserver is {{nameserver}}.", nil, "", "", "not-existing-file.conf")
				Expect(body).To(Equal("Hi the nameserver is " + defaultNameserver + "."))
			})
		})

	})
})
