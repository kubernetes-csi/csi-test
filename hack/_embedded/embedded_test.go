/*
Copyright 2018 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package embedded

import (
	"testing"

	"github.com/kubernetes-csi/csi-test/v4/pkg/sanity"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var context *sanity.TestContext

func TestMyDriverGinkgo(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "CSI Sanity Test Suite")
}

// The test suite into which the sanity tests get embedded may already
// have before/after suite functions. There can only be one such
// function. Here we define empty ones because then Ginkgo
// will start complaining at runtime when invoking the embedded case
// in hack/e2e.sh if a PR adds back such functions in the sanity test
// code.
var _ = BeforeSuite(func() {})
var _ = AfterSuite(func() {
	if context != nil {
		context.Finalize()
	}
})

var _ = Describe("MyCSIDriver", func() {
	Context("Config A", func() {
		config := sanity.NewTestConfig()
		config.Address = "/tmp/e2e-csi-sanity.sock"
		config.TestNodeVolumeAttachLimit = true

		BeforeEach(func() {})

		AfterEach(func() {})

		Describe("CSI Driver Test Suite", func() {
			context = sanity.GinkgoTest(&config)
		})
	})
})
