package brats_test

import (
	"flag"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/cloudfoundry/libbuildpack"
	"github.com/cloudfoundry/libbuildpack/bratshelper"
	"github.com/cloudfoundry/libbuildpack/cutlass"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func init() {
	flag.StringVar(&cutlass.DefaultMemory, "memory", "64M", "default memory for pushed apps")
	flag.StringVar(&cutlass.DefaultDisk, "disk", "64M", "default disk for pushed apps")
	flag.Parse()
}

var _ = SynchronizedBeforeSuite(func() []byte {
	// Run once
	return bratshelper.InitBpData(os.Getenv("CF_STACK"), ApiHasStackAssociation()).Marshal()
}, func(data []byte) {
	// Run on all nodes
	bratshelper.Data.Unmarshal(data)

	Expect(cutlass.CopyCfHome()).To(Succeed())
	cutlass.SeedRandom()
	cutlass.DefaultStdoutStderr = GinkgoWriter
})

var _ = SynchronizedAfterSuite(func() {
	// Run on all nodes
}, func() {
	// Run once
	Expect(cutlass.DeleteOrphanedRoutes()).To(Succeed())
	Expect(cutlass.DeleteBuildpack(strings.Replace(bratshelper.Data.Cached, "_buildpack", "", 1))).To(Succeed())
	Expect(cutlass.DeleteBuildpack(strings.Replace(bratshelper.Data.Uncached, "_buildpack", "", 1))).To(Succeed())
	Expect(os.Remove(bratshelper.Data.CachedFile)).To(Succeed())
})

func TestBrats(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Brats Suite")
}

func CopyBrats(version string) *cutlass.App {
	dir, err := cutlass.CopyFixture(filepath.Join(bratshelper.Data.BpDir, "fixtures", "brats"))
	Expect(err).ToNot(HaveOccurred())

	Expect(libbuildpack.NewYAML().Write(filepath.Join(dir, "buildpack.yml"), map[string]map[string]string{"nginx": {"version": version}})).To(Succeed())

	return cutlass.New(dir)
}

func ApiHasStackAssociation() bool {
	supported, err := cutlass.ApiGreaterThan("2.113.0")
	Expect(err).NotTo(HaveOccurred())
	return supported
}