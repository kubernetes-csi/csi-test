package embedded

import (
	"os"
	"testing"

	"github.com/kubernetes-csi/csi-test/pkg/sanity"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestMyDriverGinkgo(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "CSI Sanity Test Suite")
}

var _ = Describe("MyCSIDriver", func() {
	Context("Config A", func() {
		config := &sanity.Config{
			TargetPath:  os.TempDir() + "/csi",
			StagingPath: os.TempDir() + "/csi",
			Address:     "/tmp/e2e-csi-sanity.sock",
		}

		BeforeEach(func() {})

		AfterEach(func() {})

		Describe("CSI Driver Test Suite", func() {
			sanity.GinkgoTest(config)
		})
	})
})
