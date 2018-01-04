package finalize_test

import (
	"io/ioutil"
	"nginx/finalize"
	"os"

	"bytes"

	"github.com/cloudfoundry/libbuildpack"
	"github.com/cloudfoundry/libbuildpack/ansicleaner"
	"github.com/golang/mock/gomock"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

//go:generate mockgen -source=finalize.go --destination=mocks_test.go --package=finalize_test

var _ = Describe("Compile", func() {
	var (
		err       error
		buildDir  string
		depDir    string
		finalizer *finalize.Finalizer
		logger    *libbuildpack.Logger
		mockCtrl  *gomock.Controller
		buffer    *bytes.Buffer
	)

	BeforeEach(func() {
		buildDir, err = ioutil.TempDir("", "nginx-buildpack.build.")
		Expect(err).To(BeNil())

		depDir, err = ioutil.TempDir("", "nginx-buildpack.depDir.")
		Expect(err).To(BeNil())

		buffer = new(bytes.Buffer)
		logger = libbuildpack.NewLogger(ansicleaner.New(buffer))

		mockCtrl = gomock.NewController(GinkgoT())
	})

	JustBeforeEach(func() {
		finalizer = &finalize.Finalizer{
			BuildDir: buildDir,
			DepDir:   depDir,
			Log:      logger,
		}
	})

	AfterEach(func() {
		mockCtrl.Finish()

		err = os.RemoveAll(buildDir)
		Expect(err).To(BeNil())

		err = os.RemoveAll(depDir)
		Expect(err).To(BeNil())
	})
})
