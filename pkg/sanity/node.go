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
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/container-storage-interface/spec/lib/go/csi"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func isNodeCapabilitySupported(c csi.NodeClient,
	capType csi.NodeServiceCapability_RPC_Type,
) bool {

	caps, err := c.NodeGetCapabilities(
		context.Background(),
		&csi.NodeGetCapabilitiesRequest{})
	Expect(err).NotTo(HaveOccurred())
	Expect(caps).NotTo(BeNil())

	for _, cap := range caps.GetCapabilities() {
		Expect(cap.GetRpc()).NotTo(BeNil())
		if cap.GetRpc().GetType() == capType {
			return true
		}
	}
	return false
}

func isPluginCapabilitySupported(c csi.IdentityClient,
	capType csi.PluginCapability_Service_Type,
) bool {

	caps, err := c.GetPluginCapabilities(
		context.Background(),
		&csi.GetPluginCapabilitiesRequest{})
	Expect(err).NotTo(HaveOccurred())
	Expect(caps).NotTo(BeNil())
	Expect(caps.GetCapabilities()).NotTo(BeNil())

	for _, cap := range caps.GetCapabilities() {
		if cap.GetService() != nil && cap.GetService().GetType() == capType {
			return true
		}
	}
	return false
}

func isPluginCapabilityVolumeExpansionSupported(c csi.IdentityClient,
	capType csi.PluginCapability_VolumeExpansion_Type,
) bool {

	caps, err := c.GetPluginCapabilities(
		context.Background(),
		&csi.GetPluginCapabilitiesRequest{})
	Expect(err).NotTo(HaveOccurred())
	Expect(caps).NotTo(BeNil())
	Expect(caps.GetCapabilities()).NotTo(BeNil())

	for _, cap := range caps.GetCapabilities() {
		if cap.GetVolumeExpansion() != nil && cap.GetVolumeExpansion().GetType() == capType {
			return true
		}
	}
	return false
}

func isControllerPublishSupported(c csi.ControllerClient) bool {
	return isControllerCapabilitySupported(
		c,
		csi.ControllerServiceCapability_RPC_PUBLISH_UNPUBLISH_VOLUME)
}

func isControllerExpandVolumeSupported(c csi.ControllerClient) bool {
	return isControllerCapabilitySupported(c, csi.ControllerServiceCapability_RPC_EXPAND_VOLUME)
}
func isNodeStageSupported(c csi.NodeClient) bool {
	return isNodeCapabilitySupported(c, csi.NodeServiceCapability_RPC_STAGE_UNSTAGE_VOLUME)
}

func isNodeVolumeStatsSupported(c csi.NodeClient) bool {
	return isNodeCapabilitySupported(c, csi.NodeServiceCapability_RPC_GET_VOLUME_STATS)
}

func isNodeVolumeExpandSupported(c csi.NodeClient) bool {
	return isNodeCapabilitySupported(c, csi.NodeServiceCapability_RPC_EXPAND_VOLUME)

}

var _ = DescribeSanity("Node Service", func(sc *SanityContext) {
	var (
		cl *Cleanup
		c  csi.NodeClient
		s  csi.ControllerClient

		controllerPublishSupported bool
		nodeStageSupported         bool
		nodeVolumeStatsSupported   bool
	)

	BeforeEach(func() {
		c = csi.NewNodeClient(sc.Conn)
		s = csi.NewControllerClient(sc.ControllerConn)

		controllerPublishSupported = isControllerPublishSupported(s)
		nodeStageSupported = isNodeStageSupported(c)
		nodeVolumeStatsSupported = isNodeVolumeStatsSupported(c)
		cl = &Cleanup{
			Context:                    sc,
			NodeClient:                 c,
			ControllerClient:           s,
			ControllerPublishSupported: controllerPublishSupported,
			NodeStageSupported:         nodeStageSupported,
		}
	})

	AfterEach(func() {
		cl.DeleteVolumes()
	})

	Describe("NodeGetCapabilities", func() {
		It("should return appropriate capabilities", func() {
			caps, err := c.NodeGetCapabilities(
				context.Background(),
				&csi.NodeGetCapabilitiesRequest{})

			By("checking successful response")
			Expect(err).NotTo(HaveOccurred())
			Expect(caps).NotTo(BeNil())

			for _, cap := range caps.GetCapabilities() {
				Expect(cap.GetRpc()).NotTo(BeNil())

				switch cap.GetRpc().GetType() {
				case csi.NodeServiceCapability_RPC_UNKNOWN:
				case csi.NodeServiceCapability_RPC_STAGE_UNSTAGE_VOLUME:
				case csi.NodeServiceCapability_RPC_GET_VOLUME_STATS:
				case csi.NodeServiceCapability_RPC_EXPAND_VOLUME:
				default:
					Fail(fmt.Sprintf("Unknown capability: %v\n", cap.GetRpc().GetType()))
				}
			}
		})
	})

	Describe("NodeGetInfo", func() {
		var (
			i                                csi.IdentityClient
			accessibilityConstraintSupported bool
		)

		BeforeEach(func() {
			i = csi.NewIdentityClient(sc.Conn)
			accessibilityConstraintSupported = isPluginCapabilitySupported(i, csi.PluginCapability_Service_VOLUME_ACCESSIBILITY_CONSTRAINTS)
		})

		It("should return appropriate values", func() {
			ninfo, err := c.NodeGetInfo(
				context.Background(),
				&csi.NodeGetInfoRequest{})

			Expect(err).NotTo(HaveOccurred())
			Expect(ninfo).NotTo(BeNil())
			Expect(ninfo.GetNodeId()).NotTo(BeEmpty())
			Expect(ninfo.GetMaxVolumesPerNode()).NotTo(BeNumerically("<", 0))

			if accessibilityConstraintSupported {
				Expect(ninfo.GetAccessibleTopology()).NotTo(BeNil())
			}
		})
	})

	Describe("NodePublishVolume", func() {
		It("should fail when no volume id is provided", func() {
			_, err := c.NodePublishVolume(
				context.Background(),
				&csi.NodePublishVolumeRequest{
					Secrets: sc.Secrets.NodePublishVolumeSecret,
				},
			)
			Expect(err).To(HaveOccurred())

			serverError, ok := status.FromError(err)
			Expect(ok).To(BeTrue())
			Expect(serverError.Code()).To(Equal(codes.InvalidArgument))
		})

		It("should fail when no target path is provided", func() {
			_, err := c.NodePublishVolume(
				context.Background(),
				&csi.NodePublishVolumeRequest{
					VolumeId: "id",
					Secrets:  sc.Secrets.NodePublishVolumeSecret,
				},
			)
			Expect(err).To(HaveOccurred())

			serverError, ok := status.FromError(err)
			Expect(ok).To(BeTrue())
			Expect(serverError.Code()).To(Equal(codes.InvalidArgument))
		})

		It("should fail when no volume capability is provided", func() {
			_, err := c.NodePublishVolume(
				context.Background(),
				&csi.NodePublishVolumeRequest{
					VolumeId:   "id",
					TargetPath: sc.TargetPath + "/target",
					Secrets:    sc.Secrets.NodePublishVolumeSecret,
				},
			)
			Expect(err).To(HaveOccurred())

			serverError, ok := status.FromError(err)
			Expect(ok).To(BeTrue())
			Expect(serverError.Code()).To(Equal(codes.InvalidArgument))
		})
	})

	Describe("NodeUnpublishVolume", func() {
		It("should fail when no volume id is provided", func() {

			_, err := c.NodeUnpublishVolume(
				context.Background(),
				&csi.NodeUnpublishVolumeRequest{})
			Expect(err).To(HaveOccurred())

			serverError, ok := status.FromError(err)
			Expect(ok).To(BeTrue())
			Expect(serverError.Code()).To(Equal(codes.InvalidArgument))
		})

		It("should fail when no target path is provided", func() {

			_, err := c.NodeUnpublishVolume(
				context.Background(),
				&csi.NodeUnpublishVolumeRequest{
					VolumeId: "id",
				})
			Expect(err).To(HaveOccurred())

			serverError, ok := status.FromError(err)
			Expect(ok).To(BeTrue())
			Expect(serverError.Code()).To(Equal(codes.InvalidArgument))
		})
	})

	Describe("NodeStageVolume", func() {
		var (
			device string
		)

		BeforeEach(func() {
			if !nodeStageSupported {
				Skip("NodeStageVolume not supported")
			}

			device = "/dev/mock"
		})

		It("should fail when no volume id is provided", func() {
			_, err := c.NodeStageVolume(
				context.Background(),
				&csi.NodeStageVolumeRequest{
					StagingTargetPath: sc.StagingPath,
					VolumeCapability: &csi.VolumeCapability{
						AccessType: &csi.VolumeCapability_Mount{
							Mount: &csi.VolumeCapability_MountVolume{},
						},
						AccessMode: &csi.VolumeCapability_AccessMode{
							Mode: csi.VolumeCapability_AccessMode_SINGLE_NODE_WRITER,
						},
					},
					PublishContext: map[string]string{
						"device": device,
					},
					Secrets: sc.Secrets.NodeStageVolumeSecret,
				},
			)
			Expect(err).To(HaveOccurred())

			serverError, ok := status.FromError(err)
			Expect(ok).To(BeTrue())
			Expect(serverError.Code()).To(Equal(codes.InvalidArgument))
		})

		It("should fail when no staging target path is provided", func() {
			_, err := c.NodeStageVolume(
				context.Background(),
				&csi.NodeStageVolumeRequest{
					VolumeId: "id",
					VolumeCapability: &csi.VolumeCapability{
						AccessType: &csi.VolumeCapability_Mount{
							Mount: &csi.VolumeCapability_MountVolume{},
						},
						AccessMode: &csi.VolumeCapability_AccessMode{
							Mode: csi.VolumeCapability_AccessMode_SINGLE_NODE_WRITER,
						},
					},
					PublishContext: map[string]string{
						"device": device,
					},
					Secrets: sc.Secrets.NodeStageVolumeSecret,
				},
			)
			Expect(err).To(HaveOccurred())

			serverError, ok := status.FromError(err)
			Expect(ok).To(BeTrue())
			Expect(serverError.Code()).To(Equal(codes.InvalidArgument))
		})

		It("should fail when no volume capability is provided", func() {

			// Create Volume First
			By("creating a single node writer volume")
			name := UniqueString("sanity-node-stage-nocaps")

			vol, err := s.CreateVolume(
				context.Background(),
				&csi.CreateVolumeRequest{
					Name: name,
					VolumeCapabilities: []*csi.VolumeCapability{
						{
							AccessType: &csi.VolumeCapability_Mount{
								Mount: &csi.VolumeCapability_MountVolume{},
							},
							AccessMode: &csi.VolumeCapability_AccessMode{
								Mode: csi.VolumeCapability_AccessMode_SINGLE_NODE_WRITER,
							},
						},
					},
					Secrets:    sc.Secrets.CreateVolumeSecret,
					Parameters: sc.Config.TestVolumeParameters,
				},
			)
			Expect(err).NotTo(HaveOccurred())
			Expect(vol).NotTo(BeNil())
			Expect(vol.GetVolume()).NotTo(BeNil())
			Expect(vol.GetVolume().GetVolumeId()).NotTo(BeEmpty())
			cl.RegisterVolume(name, VolumeInfo{VolumeID: vol.GetVolume().GetVolumeId()})

			_, err = c.NodeStageVolume(
				context.Background(),
				&csi.NodeStageVolumeRequest{
					VolumeId:          vol.GetVolume().GetVolumeId(),
					StagingTargetPath: sc.StagingPath,
					PublishContext: map[string]string{
						"device": device,
					},
					Secrets: sc.Secrets.NodeStageVolumeSecret,
				},
			)
			Expect(err).To(HaveOccurred())

			serverError, ok := status.FromError(err)
			Expect(ok).To(BeTrue())
			Expect(serverError.Code()).To(Equal(codes.InvalidArgument))

			By("cleaning up deleting the volume")

			_, err = s.DeleteVolume(
				context.Background(),
				&csi.DeleteVolumeRequest{
					VolumeId: vol.GetVolume().GetVolumeId(),
					Secrets:  sc.Secrets.DeleteVolumeSecret,
				},
			)
			Expect(err).NotTo(HaveOccurred())
			cl.UnregisterVolume(name)
		})
	})

	Describe("NodeUnstageVolume", func() {
		BeforeEach(func() {
			if !nodeStageSupported {
				Skip("NodeUnstageVolume not supported")
			}
		})

		It("should fail when no volume id is provided", func() {

			_, err := c.NodeUnstageVolume(
				context.Background(),
				&csi.NodeUnstageVolumeRequest{
					StagingTargetPath: sc.StagingPath,
				})
			Expect(err).To(HaveOccurred())

			serverError, ok := status.FromError(err)
			Expect(ok).To(BeTrue())
			Expect(serverError.Code()).To(Equal(codes.InvalidArgument))
		})

		It("should fail when no staging target path is provided", func() {

			_, err := c.NodeUnstageVolume(
				context.Background(),
				&csi.NodeUnstageVolumeRequest{
					VolumeId: "id",
				})
			Expect(err).To(HaveOccurred())

			serverError, ok := status.FromError(err)
			Expect(ok).To(BeTrue())
			Expect(serverError.Code()).To(Equal(codes.InvalidArgument))
		})
	})

	Describe("NodeGetVolumeStats", func() {
		BeforeEach(func() {
			if !nodeVolumeStatsSupported {
				Skip("NodeGetVolume not supported")
			}
		})

		It("should fail when no volume id is provided", func() {
			_, err := c.NodeGetVolumeStats(
				context.Background(),
				&csi.NodeGetVolumeStatsRequest{
					VolumePath: "some/path",
				},
			)
			Expect(err).To(HaveOccurred())

			serverError, ok := status.FromError(err)
			Expect(ok).To(BeTrue())
			Expect(serverError.Code()).To(Equal(codes.InvalidArgument))
		})

		It("should fail when no volume path is provided", func() {
			_, err := c.NodeGetVolumeStats(
				context.Background(),
				&csi.NodeGetVolumeStatsRequest{
					VolumeId: "id",
				},
			)
			Expect(err).To(HaveOccurred())

			serverError, ok := status.FromError(err)
			Expect(ok).To(BeTrue())
			Expect(serverError.Code()).To(Equal(codes.InvalidArgument))
		})

		It("should fail when volume is not found", func() {
			_, err := c.NodeGetVolumeStats(
				context.Background(),
				&csi.NodeGetVolumeStatsRequest{
					VolumeId:   "id",
					VolumePath: "some/path",
				},
			)
			Expect(err).To(HaveOccurred())

			serverError, ok := status.FromError(err)
			Expect(ok).To(BeTrue())
			Expect(serverError.Code()).To(Equal(codes.NotFound))
		})

		It("should fail when volume does not exist on the specified path", func() {
			name := UniqueString("sanity-node-get-volume-stats")

			By("creating a single node writer volume")
			vol, err := s.CreateVolume(
				context.Background(),
				&csi.CreateVolumeRequest{
					Name: name,
					VolumeCapabilities: []*csi.VolumeCapability{
						{
							AccessType: &csi.VolumeCapability_Mount{
								Mount: &csi.VolumeCapability_MountVolume{},
							},
							AccessMode: &csi.VolumeCapability_AccessMode{
								Mode: csi.VolumeCapability_AccessMode_SINGLE_NODE_WRITER,
							},
						},
					},
					Secrets:    sc.Secrets.CreateVolumeSecret,
					Parameters: sc.Config.TestVolumeParameters,
				},
			)
			Expect(err).NotTo(HaveOccurred())
			Expect(vol).NotTo(BeNil())
			Expect(vol.GetVolume()).NotTo(BeNil())
			Expect(vol.GetVolume().GetVolumeId()).NotTo(BeEmpty())
			cl.RegisterVolume(name, VolumeInfo{VolumeID: vol.GetVolume().GetVolumeId()})

			By("getting a node id")
			nid, err := c.NodeGetInfo(
				context.Background(),
				&csi.NodeGetInfoRequest{})
			Expect(err).NotTo(HaveOccurred())
			Expect(nid).NotTo(BeNil())
			Expect(nid.GetNodeId()).NotTo(BeEmpty())

			var conpubvol *csi.ControllerPublishVolumeResponse
			if controllerPublishSupported {
				By("controller publishing volume")

				conpubvol, err = s.ControllerPublishVolume(
					context.Background(),
					&csi.ControllerPublishVolumeRequest{
						VolumeId: vol.GetVolume().GetVolumeId(),
						NodeId:   nid.GetNodeId(),
						VolumeCapability: &csi.VolumeCapability{
							AccessType: &csi.VolumeCapability_Mount{
								Mount: &csi.VolumeCapability_MountVolume{},
							},
							AccessMode: &csi.VolumeCapability_AccessMode{
								Mode: csi.VolumeCapability_AccessMode_SINGLE_NODE_WRITER,
							},
						},
						VolumeContext: vol.GetVolume().GetVolumeContext(),
						Readonly:      false,
						Secrets:       sc.Secrets.ControllerPublishVolumeSecret,
					},
				)
				Expect(err).NotTo(HaveOccurred())
				cl.RegisterVolume(name, VolumeInfo{VolumeID: vol.GetVolume().GetVolumeId(), NodeID: nid.GetNodeId()})
				Expect(conpubvol).NotTo(BeNil())
			}
			// NodeStageVolume
			if nodeStageSupported {
				By("node staging volume")
				nodestagevol, err := c.NodeStageVolume(
					context.Background(),
					&csi.NodeStageVolumeRequest{
						VolumeId: vol.GetVolume().GetVolumeId(),
						VolumeCapability: &csi.VolumeCapability{
							AccessType: &csi.VolumeCapability_Mount{
								Mount: &csi.VolumeCapability_MountVolume{},
							},
							AccessMode: &csi.VolumeCapability_AccessMode{
								Mode: csi.VolumeCapability_AccessMode_SINGLE_NODE_WRITER,
							},
						},
						StagingTargetPath: sc.StagingPath,
						VolumeContext:     vol.GetVolume().GetVolumeContext(),
						PublishContext:    conpubvol.GetPublishContext(),
						Secrets:           sc.Secrets.NodeStageVolumeSecret,
					},
				)
				Expect(err).NotTo(HaveOccurred())
				Expect(nodestagevol).NotTo(BeNil())
			}
			// NodePublishVolume
			By("publishing the volume on a node")
			var stagingPath string
			if nodeStageSupported {
				stagingPath = sc.StagingPath
			}
			nodepubvol, err := c.NodePublishVolume(
				context.Background(),
				&csi.NodePublishVolumeRequest{
					VolumeId:          vol.GetVolume().GetVolumeId(),
					TargetPath:        sc.TargetPath + "/target",
					StagingTargetPath: stagingPath,
					VolumeCapability: &csi.VolumeCapability{
						AccessType: &csi.VolumeCapability_Mount{
							Mount: &csi.VolumeCapability_MountVolume{},
						},
						AccessMode: &csi.VolumeCapability_AccessMode{
							Mode: csi.VolumeCapability_AccessMode_SINGLE_NODE_WRITER,
						},
					},
					VolumeContext:  vol.GetVolume().GetVolumeContext(),
					PublishContext: conpubvol.GetPublishContext(),
					Secrets:        sc.Secrets.NodePublishVolumeSecret,
				},
			)
			Expect(err).NotTo(HaveOccurred())
			Expect(nodepubvol).NotTo(BeNil())

			// NodeGetVolumeStats
			By("Get node volume stats")
			_, err = c.NodeGetVolumeStats(
				context.Background(),
				&csi.NodeGetVolumeStatsRequest{
					VolumeId:   vol.GetVolume().GetVolumeId(),
					VolumePath: "some/path",
				},
			)
			Expect(err).To(HaveOccurred())

			serverError, ok := status.FromError(err)
			Expect(ok).To(BeTrue())
			Expect(serverError.Code()).To(Equal(codes.NotFound))

			// NodeUnpublishVolume
			By("cleaning up calling nodeunpublish")
			nodeunpubvol, err := c.NodeUnpublishVolume(
				context.Background(),
				&csi.NodeUnpublishVolumeRequest{
					VolumeId:   vol.GetVolume().GetVolumeId(),
					TargetPath: sc.TargetPath + "/target",
				})
			Expect(err).NotTo(HaveOccurred())
			Expect(nodeunpubvol).NotTo(BeNil())

			if nodeStageSupported {
				By("cleaning up calling nodeunstage")
				nodeunstagevol, err := c.NodeUnstageVolume(
					context.Background(),
					&csi.NodeUnstageVolumeRequest{
						VolumeId:          vol.GetVolume().GetVolumeId(),
						StagingTargetPath: sc.StagingPath,
					},
				)
				Expect(err).NotTo(HaveOccurred())
				Expect(nodeunstagevol).NotTo(BeNil())
			}

			if controllerPublishSupported {
				By("cleaning up calling controllerunpublishing")

				controllerunpubvol, err := s.ControllerUnpublishVolume(
					context.Background(),
					&csi.ControllerUnpublishVolumeRequest{
						VolumeId: vol.GetVolume().GetVolumeId(),
						NodeId:   nid.GetNodeId(),
						Secrets:  sc.Secrets.ControllerUnpublishVolumeSecret,
					},
				)
				Expect(err).NotTo(HaveOccurred())
				Expect(controllerunpubvol).NotTo(BeNil())
			}

			By("cleaning up deleting the volume")

			_, err = s.DeleteVolume(
				context.Background(),
				&csi.DeleteVolumeRequest{
					VolumeId: vol.GetVolume().GetVolumeId(),
					Secrets:  sc.Secrets.DeleteVolumeSecret,
				},
			)
			Expect(err).NotTo(HaveOccurred())

		})

	})

	// CSI spec poses no specific requirements for the cluster/storage setups that a SP MUST support. To perform
	// meaningful checks the following test assumes that topology-aware provisioning on a single node setup is supported
	It("should work", func() {
		name := UniqueString("sanity-node-full")

		By("getting node information")
		ni, err := c.NodeGetInfo(
			context.Background(),
			&csi.NodeGetInfoRequest{})
		Expect(err).NotTo(HaveOccurred())
		Expect(ni).NotTo(BeNil())
		Expect(ni.GetNodeId()).NotTo(BeEmpty())

		var accReqs *csi.TopologyRequirement
		if ni.AccessibleTopology != nil {
			// Topology requirements are honored if provided by the driver
			accReqs = &csi.TopologyRequirement{
				Requisite: []*csi.Topology{ni.AccessibleTopology},
			}
		}

		// Create Volume First
		By("creating a single node writer volume")
		vol, err := s.CreateVolume(
			context.Background(),
			&csi.CreateVolumeRequest{
				Name: name,
				VolumeCapabilities: []*csi.VolumeCapability{
					{
						AccessType: &csi.VolumeCapability_Mount{
							Mount: &csi.VolumeCapability_MountVolume{},
						},
						AccessMode: &csi.VolumeCapability_AccessMode{
							Mode: csi.VolumeCapability_AccessMode_SINGLE_NODE_WRITER,
						},
					},
				},
				Secrets:                   sc.Secrets.CreateVolumeSecret,
				Parameters:                sc.Config.TestVolumeParameters,
				AccessibilityRequirements: accReqs,
			},
		)
		Expect(err).NotTo(HaveOccurred())
		Expect(vol).NotTo(BeNil())
		Expect(vol.GetVolume()).NotTo(BeNil())
		Expect(vol.GetVolume().GetVolumeId()).NotTo(BeEmpty())
		cl.RegisterVolume(name, VolumeInfo{VolumeID: vol.GetVolume().GetVolumeId()})

		var conpubvol *csi.ControllerPublishVolumeResponse
		if controllerPublishSupported {
			By("controller publishing volume")

			conpubvol, err = s.ControllerPublishVolume(
				context.Background(),
				&csi.ControllerPublishVolumeRequest{
					VolumeId: vol.GetVolume().GetVolumeId(),
					NodeId:   ni.GetNodeId(),
					VolumeCapability: &csi.VolumeCapability{
						AccessType: &csi.VolumeCapability_Mount{
							Mount: &csi.VolumeCapability_MountVolume{},
						},
						AccessMode: &csi.VolumeCapability_AccessMode{
							Mode: csi.VolumeCapability_AccessMode_SINGLE_NODE_WRITER,
						},
					},
					VolumeContext: vol.GetVolume().GetVolumeContext(),
					Readonly:      false,
					Secrets:       sc.Secrets.ControllerPublishVolumeSecret,
				},
			)
			Expect(err).NotTo(HaveOccurred())
			cl.RegisterVolume(name, VolumeInfo{VolumeID: vol.GetVolume().GetVolumeId(), NodeID: ni.GetNodeId()})
			Expect(conpubvol).NotTo(BeNil())
		}
		// NodeStageVolume
		if nodeStageSupported {
			By("node staging volume")
			nodestagevol, err := c.NodeStageVolume(
				context.Background(),
				&csi.NodeStageVolumeRequest{
					VolumeId: vol.GetVolume().GetVolumeId(),
					VolumeCapability: &csi.VolumeCapability{
						AccessType: &csi.VolumeCapability_Mount{
							Mount: &csi.VolumeCapability_MountVolume{},
						},
						AccessMode: &csi.VolumeCapability_AccessMode{
							Mode: csi.VolumeCapability_AccessMode_SINGLE_NODE_WRITER,
						},
					},
					StagingTargetPath: sc.StagingPath,
					VolumeContext:     vol.GetVolume().GetVolumeContext(),
					PublishContext:    conpubvol.GetPublishContext(),
					Secrets:           sc.Secrets.NodeStageVolumeSecret,
				},
			)
			Expect(err).NotTo(HaveOccurred())
			Expect(nodestagevol).NotTo(BeNil())
		}
		// NodePublishVolume
		By("publishing the volume on a node")
		var stagingPath string
		if nodeStageSupported {
			stagingPath = sc.StagingPath
		}
		nodepubvol, err := c.NodePublishVolume(
			context.Background(),
			&csi.NodePublishVolumeRequest{
				VolumeId:          vol.GetVolume().GetVolumeId(),
				TargetPath:        sc.TargetPath + "/target",
				StagingTargetPath: stagingPath,
				VolumeCapability: &csi.VolumeCapability{
					AccessType: &csi.VolumeCapability_Mount{
						Mount: &csi.VolumeCapability_MountVolume{},
					},
					AccessMode: &csi.VolumeCapability_AccessMode{
						Mode: csi.VolumeCapability_AccessMode_SINGLE_NODE_WRITER,
					},
				},
				VolumeContext:  vol.GetVolume().GetVolumeContext(),
				PublishContext: conpubvol.GetPublishContext(),
				Secrets:        sc.Secrets.NodePublishVolumeSecret,
			},
		)
		Expect(err).NotTo(HaveOccurred())
		Expect(nodepubvol).NotTo(BeNil())

		// NodeGetVolumeStats
		if nodeVolumeStatsSupported {
			By("Get node volume stats")
			statsResp, err := c.NodeGetVolumeStats(
				context.Background(),
				&csi.NodeGetVolumeStatsRequest{
					VolumeId:   vol.GetVolume().GetVolumeId(),
					VolumePath: sc.TargetPath + "/target",
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(statsResp.GetUsage()).ToNot(BeNil())
		}

		// NodeUnpublishVolume
		By("cleaning up calling nodeunpublish")
		nodeunpubvol, err := c.NodeUnpublishVolume(
			context.Background(),
			&csi.NodeUnpublishVolumeRequest{
				VolumeId:   vol.GetVolume().GetVolumeId(),
				TargetPath: sc.TargetPath + "/target",
			})
		Expect(err).NotTo(HaveOccurred())
		Expect(nodeunpubvol).NotTo(BeNil())

		if nodeStageSupported {
			By("cleaning up calling nodeunstage")
			nodeunstagevol, err := c.NodeUnstageVolume(
				context.Background(),
				&csi.NodeUnstageVolumeRequest{
					VolumeId:          vol.GetVolume().GetVolumeId(),
					StagingTargetPath: sc.StagingPath,
				},
			)
			Expect(err).NotTo(HaveOccurred())
			Expect(nodeunstagevol).NotTo(BeNil())
		}

		if controllerPublishSupported {
			By("cleaning up calling controllerunpublishing")

			controllerunpubvol, err := s.ControllerUnpublishVolume(
				context.Background(),
				&csi.ControllerUnpublishVolumeRequest{
					VolumeId: vol.GetVolume().GetVolumeId(),
					NodeId:   ni.GetNodeId(),
					Secrets:  sc.Secrets.ControllerUnpublishVolumeSecret,
				},
			)
			Expect(err).NotTo(HaveOccurred())
			Expect(controllerunpubvol).NotTo(BeNil())
		}

		By("cleaning up deleting the volume")

		_, err = s.DeleteVolume(
			context.Background(),
			&csi.DeleteVolumeRequest{
				VolumeId: vol.GetVolume().GetVolumeId(),
				Secrets:  sc.Secrets.DeleteVolumeSecret,
			},
		)
		Expect(err).NotTo(HaveOccurred())
	})
})

type VolumeGetter interface {
	GetVolume() *csi.Volume
}

func getTotalBytesFromUsageList(usage []*csi.VolumeUsage) int64 {
	for _, u := range usage {
		if u.Unit == csi.VolumeUsage_BYTES {
			return u.Total
		}
	}
	return 0
}

var _ = DescribeSanity("ExpandVolume [Node Server]", func(sc *SanityContext) {
	var (
		c  csi.ControllerClient
		n  csi.NodeClient
		i  csi.IdentityClient
		cl *Cleanup
	)
	BeforeEach(func() {
		i = csi.NewIdentityClient(sc.ControllerConn)
		c = csi.NewControllerClient(sc.ControllerConn)
		n = csi.NewNodeClient(sc.Conn)
		if !isNodeVolumeExpandSupported(n) {
			Skip("NodeExpandVolume not supported")
		}
		if !isPluginCapabilityVolumeExpansionSupported(i, csi.PluginCapability_VolumeExpansion_ONLINE) {
			Skip("VolumeExpansion Online is not supported")
		}
		cl = &Cleanup{
			ControllerClient: c,
			NodeClient:       n,
			Context:          sc,
		}
	})
	AfterEach(func() {
		cl.DeleteVolumes()
	})

	It("should work", func() {
		name := UniqueString("sanity-node-expand-volume")
		ni, err := n.NodeGetInfo(context.TODO(), &csi.NodeGetInfoRequest{})
		Expect(err).NotTo(HaveOccurred())
		Expect(ni).NotTo(BeNil())
		vol := createVolume(sc, cl, name, c)
		pub := maybeControllerPublishVolume(sc, cl, c, name, vol, ni.GetNodeId())
		_ = maybeNodeStageVolume(sc, n, vol, pub.GetPublishContext())
		_ = nodePublishVolume(sc, n, vol, pub.GetPublishContext())

		stat, _ := maybeGetVolumeStats(n, vol, sc.TargetPath+"/target")
		var beforeTotalBytes int64
		if stat != nil {
			beforeTotalBytes = getTotalBytesFromUsageList(stat.Usage)
		}
		_ = maybeControllerExpandVolume(sc, c, vol)
		By("Expanding the volume on a node")
		rsp, err := n.NodeExpandVolume(
			context.Background(),
			&csi.NodeExpandVolumeRequest{
				VolumeId: vol.GetVolume().GetVolumeId(),
				CapacityRange: &csi.CapacityRange{
					RequiredBytes: TestVolumeExpandSize(sc),
				},
				VolumePath: sc.TargetPath + "/target",
			},
		)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp).NotTo(BeNil())
		newStat, _ := maybeGetVolumeStats(n, vol, sc.TargetPath+"/target")
		if newStat != nil && stat != nil {
			Expect(getTotalBytesFromUsageList(newStat.Usage) > beforeTotalBytes).To(BeTrue())
		}

		_ = nodeUnpublishVolume(sc, n, vol)
		_ = maybeNodeUnstageVolume(sc, n, vol)
		_ = maybeControllerUnpublishVolume(sc, c, vol, ni.GetNodeId())
		deleteVolume(sc, c, vol)

	})

})

func createVolume(sc *SanityContext, cl *Cleanup, name string, c csi.ControllerClient) (rsp *csi.CreateVolumeResponse) {
	By("Creating the volume")
	rsp, err := c.CreateVolume(
		context.Background(),
		&csi.CreateVolumeRequest{
			Name: name,
			VolumeCapabilities: []*csi.VolumeCapability{
				{
					AccessType: &csi.VolumeCapability_Mount{
						Mount: &csi.VolumeCapability_MountVolume{},
					},
					AccessMode: &csi.VolumeCapability_AccessMode{
						Mode: csi.VolumeCapability_AccessMode_SINGLE_NODE_WRITER,
					},
				},
			},
			Secrets:    sc.Secrets.CreateVolumeSecret,
			Parameters: sc.Config.TestVolumeParameters,
		},
	)
	Expect(err).NotTo(HaveOccurred())
	Expect(rsp).NotTo(BeNil())
	Expect(rsp.GetVolume()).NotTo(BeNil())
	Expect(rsp.GetVolume().GetVolumeId()).NotTo(BeEmpty())
	cl.RegisterVolume(name, VolumeInfo{VolumeID: rsp.GetVolume().GetVolumeId()})
	return
}

func deleteVolume(sc *SanityContext, c csi.ControllerClient, vol VolumeGetter) {
	By("cleaning up deleting the volume")

	_, err := c.DeleteVolume(
		context.Background(),
		&csi.DeleteVolumeRequest{
			VolumeId: vol.GetVolume().GetVolumeId(),
			Secrets:  sc.Secrets.DeleteVolumeSecret,
		},
	)
	Expect(err).NotTo(HaveOccurred())
}

func nodePublishVolume(sc *SanityContext, n csi.NodeClient, vol VolumeGetter,
	publishContext map[string]string) (rsp *csi.NodePublishVolumeResponse) {
	// NodePublishVolume
	By("publishing the volume on a node")
	var stagingPath string
	if isNodeStageSupported(n) {
		stagingPath = sc.StagingPath
	}
	rsp, err := n.NodePublishVolume(
		context.Background(),
		&csi.NodePublishVolumeRequest{
			VolumeId:          vol.GetVolume().GetVolumeId(),
			TargetPath:        sc.TargetPath + "/target",
			StagingTargetPath: stagingPath,
			VolumeCapability: &csi.VolumeCapability{
				AccessType: &csi.VolumeCapability_Mount{
					Mount: &csi.VolumeCapability_MountVolume{},
				},
				AccessMode: &csi.VolumeCapability_AccessMode{
					Mode: csi.VolumeCapability_AccessMode_SINGLE_NODE_WRITER,
				},
			},
			VolumeContext:  vol.GetVolume().GetVolumeContext(),
			PublishContext: publishContext,
			Secrets:        sc.Secrets.NodePublishVolumeSecret,
		},
	)
	Expect(err).NotTo(HaveOccurred())
	Expect(rsp).NotTo(BeNil())
	return
}

func nodeUnpublishVolume(sc *SanityContext, n csi.NodeClient, vol VolumeGetter) (rsp *csi.NodeUnpublishVolumeResponse) {
	By("cleaning up calling nodeunpublish")
	rsp, err := n.NodeUnpublishVolume(
		context.Background(),
		&csi.NodeUnpublishVolumeRequest{
			VolumeId:   vol.GetVolume().GetVolumeId(),
			TargetPath: sc.TargetPath + "/target",
		})
	Expect(err).NotTo(HaveOccurred())
	Expect(rsp).NotTo(BeNil())
	return
}

func maybeNodeStageVolume(
	sc *SanityContext, n csi.NodeClient, vol VolumeGetter,
	publishContext map[string]string) (rsp *csi.NodeStageVolumeResponse) {
	if isNodeStageSupported(n) {

		rsp, err := n.NodeStageVolume(
			context.Background(),
			&csi.NodeStageVolumeRequest{
				VolumeId: vol.GetVolume().GetVolumeId(),
				VolumeCapability: &csi.VolumeCapability{
					AccessType: &csi.VolumeCapability_Mount{
						Mount: &csi.VolumeCapability_MountVolume{},
					},
					AccessMode: &csi.VolumeCapability_AccessMode{
						Mode: csi.VolumeCapability_AccessMode_SINGLE_NODE_WRITER,
					},
				},
				StagingTargetPath: sc.StagingPath,
				VolumeContext:     vol.GetVolume().GetVolumeContext(),
				PublishContext:    publishContext,
				Secrets:           sc.Secrets.NodeStageVolumeSecret,
			},
		)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp).NotTo(BeNil())
	}
	return
}

func maybeNodeUnstageVolume(sc *SanityContext, n csi.NodeClient, vol VolumeGetter) (rsp *csi.NodeUnstageVolumeResponse) {
	if isNodeStageSupported(n) {
		By("cleaning up calling nodeunstage")
		rsp, err := n.NodeUnstageVolume(
			context.Background(),
			&csi.NodeUnstageVolumeRequest{
				VolumeId:          vol.GetVolume().GetVolumeId(),
				StagingTargetPath: sc.StagingPath,
			},
		)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp).NotTo(BeNil())
	}
	return
}

func maybeControllerExpandVolume(sc *SanityContext, c csi.ControllerClient, vol VolumeGetter) (rsp *csi.ControllerExpandVolumeResponse) {
	if isControllerExpandVolumeSupported(c) {
		By("calling controllerexpandvolume")
		rsp, err := c.ControllerExpandVolume(context.Background(),
			&csi.ControllerExpandVolumeRequest{
				CapacityRange: &csi.CapacityRange{
					RequiredBytes: TestVolumeExpandSize(sc),
				},
				VolumeId: vol.GetVolume().GetVolumeId(),
			},
		)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp).NotTo(BeNil())
	}
	return
}

func maybeGetVolumeStats(
	c csi.NodeClient, vol VolumeGetter, targetPath string) (rsp *csi.NodeGetVolumeStatsResponse, err error) {
	if isNodeVolumeStatsSupported(c) {
		By("Get node volume stats")
		rsp, err = c.NodeGetVolumeStats(
			context.Background(),
			&csi.NodeGetVolumeStatsRequest{
				VolumeId:   vol.GetVolume().GetVolumeId(),
				VolumePath: targetPath,
			},
		)
		Expect(err).ToNot(HaveOccurred())
		Expect(rsp.GetUsage()).ToNot(BeNil())
	}
	return
}

func maybeControllerPublishVolume(
	sc *SanityContext, cl *Cleanup, c csi.ControllerClient, name string, vol *csi.CreateVolumeResponse,
	nodeID string) (rsp *csi.ControllerPublishVolumeResponse) {
	if isControllerPublishSupported(c) {
		By("controller publishing volume")
		rsp, err := c.ControllerPublishVolume(
			context.Background(),
			&csi.ControllerPublishVolumeRequest{
				VolumeId: vol.GetVolume().GetVolumeId(),
				NodeId:   nodeID,
				VolumeCapability: &csi.VolumeCapability{
					AccessType: &csi.VolumeCapability_Mount{
						Mount: &csi.VolumeCapability_MountVolume{},
					},
					AccessMode: &csi.VolumeCapability_AccessMode{
						Mode: csi.VolumeCapability_AccessMode_SINGLE_NODE_WRITER,
					},
				},
				VolumeContext: vol.GetVolume().GetVolumeContext(),
				Readonly:      false,
				Secrets:       sc.Secrets.ControllerPublishVolumeSecret,
			},
		)
		Expect(err).NotTo(HaveOccurred())
		cl.RegisterVolume(name, VolumeInfo{VolumeID: vol.GetVolume().GetVolumeId(), NodeID: nodeID})
		Expect(rsp).NotTo(BeNil())
	}
	return
}

func maybeControllerUnpublishVolume(sc *SanityContext, c csi.ControllerClient, vol VolumeGetter, nodeID string) (rsp *csi.ControllerUnpublishVolumeResponse) {
	if isControllerPublishSupported(c) {
		By("cleaning up calling controllerunpublishing")

		rsp, err := c.ControllerUnpublishVolume(
			context.Background(),
			&csi.ControllerUnpublishVolumeRequest{
				VolumeId: vol.GetVolume().GetVolumeId(),
				NodeId:   nodeID,
				Secrets:  sc.Secrets.ControllerUnpublishVolumeSecret,
			},
		)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp).NotTo(BeNil())
	}
	return
}
