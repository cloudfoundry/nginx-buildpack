package main_test

import (
	"io/ioutil"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("varify", func() {
	var (
		tmpDir string
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
			body := runCli(tmpDir, "Hi the port is {{port}}.", []string{"PORT=8080"})
			Expect(body).To(Equal("Hi the port is 8080."))
		})

		It("templates a load_module directive using the 'module' func", func() {
			body := runCli(tmpDir, `{{module "foo"}}`, []string{"NGINX_MODULES=/some/directory"})
			Expect(body).To(Equal("load_module /some/directory/foo.so;"))
		})

		It("templates environment variables using the 'env' func", func() {
			body := runCli(tmpDir, `The env var FOO is {{env "FOO"}}`, []string{"FOO=BAR"})
			Expect(body).To(Equal("The env var FOO is BAR"))
		})
	})
})
