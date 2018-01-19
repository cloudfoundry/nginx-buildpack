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
		It("replaces {{.Port}} in file", func() {
			body := runCli(tmpDir, "Hi the port is {{.Port}}.", []string{"PORT=8080"})
			Expect(body).To(Equal("Hi the port is 8080."))
		})
	})
})
