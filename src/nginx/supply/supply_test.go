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

	Describe("Run", func() {
	})
})
