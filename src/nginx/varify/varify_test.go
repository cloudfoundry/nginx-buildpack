package main_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
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
			body, _ := runCli(tmpDir, "Hi the port is {{port}}.", []string{"PORT=8080"}, "", "", "", "", "", 0)
			Expect(body).To(Equal("Hi the port is 8080."))
		})

		It("templates environment variables using the 'env' func", func() {
			body, _ := runCli(tmpDir, `The env var FOO is {{env "FOO"}}`, []string{"FOO=BAR"}, "", "", "", "", "", 0)
			Expect(body).To(Equal("The env var FOO is BAR"))
		})

		It("plaintext templates environment variables listed in buildpack.yml", func() {
			textBody := `The env var FOO is {{env "FOO"}}`
			bpYMLPath := filepath.Join(tmpDir, "buildpack.yml")
			contents := `---
nginx:
  version: stable
  plaintext_env_vars:
    - "FOO"
`
			defer os.RemoveAll(bpYMLPath)
			Expect(ioutil.WriteFile(bpYMLPath, []byte(contents), os.ModePerm)).To(Succeed())
			body, _ := runCli(tmpDir, textBody, []string{`FOO={"abcd":1234}`}, "", "", "", "", bpYMLPath, 0)
			Expect(body).To(Equal(`The env var FOO is {"abcd":1234}`))
		})

		Describe("templating conf with include files", func() {
			It("parses include file", func() {
				const nginxConfStr = `
	include    custom.conf;
`
				const customConfStr = `
  server {
		listen       {{port}};
	}
`
				Expect(ioutil.WriteFile(filepath.Join(tmpDir, "custom.conf"), []byte(customConfStr), os.ModePerm)).To(Succeed())
				body, _ := runCli(tmpDir, nginxConfStr, []string{"PORT=8080"}, "", "", "", "", "", 0)
				Expect(body).To(Equal(nginxConfStr))

				contents, err := ioutil.ReadFile(filepath.Join(tmpDir, "custom.conf"))
				Expect(err).NotTo(HaveOccurred())
				Expect(string(contents)).To(Equal(`
  server {
		listen       8080;
	}
`))
			})
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
					body, _ := runCli(tmpDir, `{{module "local"}}`, nil, localModulePath, globalModulePath, "", "", "", 0)
					Expect(body).To(Equal(fmt.Sprintf("load_module %s/local.so;", localModulePath)))
				})
			})

			Context("when the module is in global modules directory", func() {
				It("loads the module from the global directory", func() {
					body, _ := runCli(tmpDir, `{{module "global"}}`, nil, localModulePath, globalModulePath, "", "", "", 0)
					Expect(body).To(Equal(fmt.Sprintf("load_module %s/global.so;", globalModulePath)))
				})
			})
		})

		Context("templating a nameservers directive using the 'nameservers' func", func() {
			var defaultNameServer = "169.254.0.123"
			var nameserver1 = "123.245.67.89"
			var nameserver2 = "89.67.245.123"

			It("reads nameservers from the simple resolv-conf file", func() {
				var resolvConfPath = filepath.Join(tmpDir, "resolv-simple.conf")
				Expect(ioutil.WriteFile(resolvConfPath, []byte("nameserver "+nameserver1), 0644)).To(Succeed())
				body, _ := runCli(tmpDir, "Hi the nameservers are {{nameservers}}.", nil, "", "", resolvConfPath, defaultNameServer, "", 0)
				Expect(body).To(Equal("Hi the nameservers are " + nameserver1 + "."))
			})

			It("reads nameservers from the unusual resolv-conf file", func() {
				var resolvConfPath = filepath.Join(tmpDir, "resolv-unusual.conf")
				Expect(ioutil.WriteFile(resolvConfPath, []byte("# comment 1\n  \t  nameserver "+nameserver1+"  \t  \n# comment 2"), 0644)).To(Succeed())
				body, _ := runCli(tmpDir, "Hi the nameservers are {{nameservers}}.", nil, "", "", resolvConfPath, defaultNameServer, "", 0)
				Expect(body).To(Equal("Hi the nameservers are " + nameserver1 + "."))
			})

			It("reads nameservers from the resolv-conf file with multiple entries", func() {
				var resolvConfPath = filepath.Join(tmpDir, "resolv-multiple.conf")
				Expect(ioutil.WriteFile(resolvConfPath, []byte("nameserver "+nameserver1+"\nnameserver "+nameserver2), 0644)).To(Succeed())
				body, _ := runCli(tmpDir, "Hi the nameservers are {{nameservers}}.", nil, "", "", resolvConfPath, defaultNameServer, "", 0)
				Expect(body).To(Equal("Hi the nameservers are " + nameserver1 + " " + nameserver2 + "."))
			})

			It("set the default nameservers if the resolv-conf file is empty", func() {
				var resolvConfPath = filepath.Join(tmpDir, "resolv-empty.conf")
				Expect(ioutil.WriteFile(resolvConfPath, []byte(""), 666)).To(Succeed())
				body, _ := runCli(tmpDir, "Hi the nameservers are {{nameservers}}.", nil, "", "", resolvConfPath, defaultNameServer, "", 0)
				Expect(body).To(Equal("Hi the nameservers are " + defaultNameServer + "."))
			})

			It("set the default nameservers if the resolv-conf file does't exist", func() {
				body, _ := runCli(tmpDir, "Hi the nameservers are {{nameservers}}.", nil, "", "", "not-existing-file.conf", defaultNameServer, "", 0)
				Expect(body).To(Equal("Hi the nameservers are " + defaultNameServer + "."))
			})
		})

		Context("Failure cases", func() {
			Describe("templating conf with non-existent include files", func() {
				It("errors with a message", func() {
					const nginxConfStr = `
	include    idontexist.conf;
`
					_, session := runCli(tmpDir, nginxConfStr, []string{"PORT=8080"}, "", "", "", "", "", 1)
					Expect(session.Err).To(gbytes.Say(fmt.Sprintf(`Could not read config file: %s/nginx.conf`, tmpDir)))
					Expect(session.Err).To(gbytes.Say(`idontexist.conf: no such file or directory`))
				})
			})
		})

	})
})
