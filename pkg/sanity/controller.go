/*
Copyright 2017 Kubernetes Authors.

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

package sanity

import (
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/container-storage-interface/spec/lib/go/csi"
	context "golang.org/x/net/context"
	"google.golang.org/grpc"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ControllerGetCapabilities [Controller Server]", func() {
	var (
		c    csi.ControllerClient
		conn *grpc.ClientConn
	)

	BeforeEach(func() {
		var err error
		conn, err = grpc.Dial(driverAddress, grpc.WithInsecure())
		Expect(err).ToNot(HaveOccurred())
		c = csi.NewControllerClient(conn)
	})

	AfterEach(func() {
		conn.Close()
	})

	It("should fail when no version is provided", func() {
		_, err := c.ControllerGetCapabilities(
			context.Background(),
			&csi.ControllerGetCapabilitiesRequest{})
		Expect(err).To(HaveOccurred())

		serverError, ok := status.FromError(err)
		Expect(ok).To(BeTrue())
		Expect(serverError.Code()).To(Equal(codes.InvalidArgument))
	})

	It("should return appropriate capabilities", func() {
		caps, err := c.ControllerGetCapabilities(
			context.Background(),
			&csi.ControllerGetCapabilitiesRequest{
				Version: csiClientVersion,
			})

		By("checking successful response")
		Expect(err).NotTo(HaveOccurred())
		Expect(caps).NotTo(BeNil())
		Expect(caps.GetCapabilities()).NotTo(BeNil())

		for _, cap := range caps.GetCapabilities() {
			Expect(cap.GetRpc()).NotTo(BeNil())

			switch cap.GetRpc().GetType() {
			case csi.ControllerServiceCapability_RPC_CREATE_DELETE_VOLUME:
			case csi.ControllerServiceCapability_RPC_PUBLISH_UNPUBLISH_VOLUME:
			case csi.ControllerServiceCapability_RPC_LIST_VOLUMES:
			case csi.ControllerServiceCapability_RPC_GET_CAPACITY:
			default:
				Fail(fmt.Sprintf("Unknown capability: %v\n", cap.GetRpc().GetType()))
			}
		}
	})
})