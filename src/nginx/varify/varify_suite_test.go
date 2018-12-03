package main_test

import (
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

func TestVarify(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Varify Suite")
}

var pathToCli string
var _ = BeforeSuite(func() {
	var err error
	pathToCli, err = gexec.Build("github.com/cloudfoundry/nginx-buildpack/src/nginx/varify")
	Expect(err).ToNot(HaveOccurred())
})

var _ = AfterSuite(func() {
	gexec.CleanupBuildArtifacts()
})

func runCli(tmpDir, body string, env []string, localModulePath, globalModulePath string) string {
	Expect(ioutil.WriteFile(filepath.Join(tmpDir, "nginx.conf"), []byte(body), 0644)).To(Succeed())

	command := exec.Command(pathToCli, filepath.Join(tmpDir, "nginx.conf"), localModulePath, globalModulePath)
	command.Env = env
	session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
	Expect(err).ToNot(HaveOccurred())
	Eventually(session).Should(gexec.Exit(0))

	output, err := ioutil.ReadFile(filepath.Join(tmpDir, "nginx.conf"))
	Expect(err).ToNot(HaveOccurred())

	return string(output)
}
