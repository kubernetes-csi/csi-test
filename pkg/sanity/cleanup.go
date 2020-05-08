/*
Copyright 2018 Intel Corporation

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
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/container-storage-interface/spec/lib/go/csi"

	. "github.com/onsi/gomega"
)

// VolumeInfo keeps track of the information needed to delete a volume.
type VolumeInfo struct {
	// Node on which the volume was published, empty if none
	// or publishing is not supported.
	NodeID string

	// Volume ID assigned by CreateVolume.
	VolumeID string
}

// Cleanup keeps track of resources, in particular volumes and snapshots, that
// need to be freed when testing is done. It implements both ControllerClient
// and NodeClient and should be used as the only interaction point to either
// APIs. That way, Cleanup can ensure that resources are marked for cleanup as
// necessary.
// All methods can be called concurrently.
type Cleanup struct {
	Context *TestContext
	// ControllerClient is meant for struct-internal use only
	csi.ControllerClient
	// NodeClient is meant for struct-internal use only
	csi.NodeClient
	ControllerPublishSupported bool
	NodeStageSupported         bool

	// mutex protects access to the volumes and snapshots maps.
	mutex sync.Mutex
	// volumes maps from volume IDs to VolumeInfo structs and records if a given
	// volume must be cleaned up.
	volumes map[string]*VolumeInfo
	// snapshots is keyed by snapshot IDs and records if a given snapshot must
	// be cleaned up.
	snapshots map[string]bool
}

// ControllerClient interface wrappers

// CreateVolume proxies to a Controller service implementation and registers the
// volume for cleanup.
func (cl *Cleanup) CreateVolume(ctx context.Context, in *csi.CreateVolumeRequest, _ ...grpc.CallOption) (*csi.CreateVolumeResponse, error) {
	return cl.createVolume(ctx, in)
}

// DeleteVolume proxies to a Controller service implementation and unregisters
// the volume from cleanup.
func (cl *Cleanup) DeleteVolume(ctx context.Context, in *csi.DeleteVolumeRequest, _ ...grpc.CallOption) (*csi.DeleteVolumeResponse, error) {
	return cl.deleteVolume(ctx, in)
}

// ControllerPublishVolume proxies to a Controller service implementation and
// adds the node ID to the corresponding volume for cleanup.
func (cl *Cleanup) ControllerPublishVolume(ctx context.Context, in *csi.ControllerPublishVolumeRequest, _ ...grpc.CallOption) (*csi.ControllerPublishVolumeResponse, error) {
	return cl.controllerPublishVolume(ctx, in)
}

// CreateSnapshot proxies to a Controller service implementation and registers
// the snapshot for cleanup.
func (cl *Cleanup) CreateSnapshot(ctx context.Context, in *csi.CreateSnapshotRequest, _ ...grpc.CallOption) (*csi.CreateSnapshotResponse, error) {
	return cl.createSnapshot(ctx, in)
}

// DeleteSnapshot proxies to a Controller service implementation and unregisters
// the snapshot from cleanup.
func (cl *Cleanup) DeleteSnapshot(ctx context.Context, in *csi.DeleteSnapshotRequest, _ ...grpc.CallOption) (*csi.DeleteSnapshotResponse, error) {
	return cl.deleteSnapshot(ctx, in)
}

// MustCreateVolume is like CreateVolume but asserts that the volume was
// successfully created.
func (cl *Cleanup) MustCreateVolume(ctx context.Context, req *csi.CreateVolumeRequest) *csi.CreateVolumeResponse {
	vol, err := cl.createVolume(ctx, req)
	Expect(err).NotTo(HaveOccurred())
	Expect(vol).NotTo(BeNil())
	Expect(vol.GetVolume()).NotTo(BeNil())
	Expect(vol.GetVolume().GetVolumeId()).NotTo(BeEmpty())
	return vol
}

func (cl *Cleanup) createVolume(ctx context.Context, req *csi.CreateVolumeRequest) (*csi.CreateVolumeResponse, error) {
	vol, err := cl.ControllerClient.CreateVolume(ctx, req)
	if err == nil && vol != nil && vol.GetVolume().GetVolumeId() != "" {
		cl.registerVolume(VolumeInfo{VolumeID: vol.GetVolume().GetVolumeId()})
	}
	return vol, err
}

func (cl *Cleanup) deleteVolume(ctx context.Context, req *csi.DeleteVolumeRequest) (*csi.DeleteVolumeResponse, error) {
	vol, err := cl.ControllerClient.DeleteVolume(ctx, req)
	if err == nil {
		cl.unregisterVolume(req.VolumeId)
	}
	return vol, err
}

// MustControllerPublishVolume is like ControllerPublishVolume but asserts that
// the volume was successfully controller-published.
func (cl *Cleanup) MustControllerPublishVolume(ctx context.Context, req *csi.ControllerPublishVolumeRequest) *csi.ControllerPublishVolumeResponse {
	conpubvol, err := cl.controllerPublishVolume(ctx, req)
	Expect(err).NotTo(HaveOccurred())
	Expect(conpubvol).NotTo(BeNil())
	return conpubvol
}

func (cl *Cleanup) controllerPublishVolume(ctx context.Context, req *csi.ControllerPublishVolumeRequest) (*csi.ControllerPublishVolumeResponse, error) {
	conpubvol, err := cl.ControllerClient.ControllerPublishVolume(ctx, req)
	if err == nil && req.VolumeId != "" && req.NodeId != "" {
		cl.registerVolume(VolumeInfo{VolumeID: req.VolumeId, NodeID: req.NodeId})
	}
	return conpubvol, err
}

// registerVolume adds or updates an entry for given volume.
func (cl *Cleanup) registerVolume(info VolumeInfo) {
	Expect(info).NotTo(BeNil())
	Expect(info.VolumeID).NotTo(BeEmpty())
	cl.mutex.Lock()
	defer cl.mutex.Unlock()
	if cl.volumes == nil {
		cl.volumes = make(map[string]*VolumeInfo)
	}
	cl.volumes[info.VolumeID] = &info
}

// unregisterVolume removes the entry for the volume with the
// given ID, thus preventing all cleanup operations for it.
func (cl *Cleanup) unregisterVolume(id string) {
	cl.mutex.Lock()
	defer cl.mutex.Unlock()
	cl.unregisterVolumeNoLock(id)
}

func (cl *Cleanup) unregisterVolumeNoLock(id string) {
	Expect(id).NotTo(BeEmpty())
	if cl.volumes != nil {
		delete(cl.volumes, id)
	}
}

// MustCreateSnapshot is like CreateSnapshot but asserts that the snapshot was
// successfully created.
func (cl *Cleanup) MustCreateSnapshot(ctx context.Context, req *csi.CreateSnapshotRequest) *csi.CreateSnapshotResponse {
	snap, err := cl.createSnapshot(ctx, req)
	Expect(err).NotTo(HaveOccurred())
	Expect(snap).NotTo(BeNil())
	verifySnapshotInfo(snap.GetSnapshot())
	return snap
}

// MustCreateSnapshotFromVolumeRequest creates a volume from the given
// CreateVolumeRequest and a snapshot subsequently. It registers the volume and
// snapshot and asserts that both were created successfully.
func (cl *Cleanup) MustCreateSnapshotFromVolumeRequest(ctx context.Context, req *csi.CreateVolumeRequest, snapshotName string) (*csi.CreateSnapshotResponse, *csi.CreateVolumeResponse) {
	vol := cl.MustCreateVolume(ctx, req)
	snap := cl.MustCreateSnapshot(ctx, MakeCreateSnapshotReq(cl.Context, snapshotName, vol.Volume.VolumeId))
	return snap, vol
}

func (cl *Cleanup) createSnapshot(ctx context.Context, req *csi.CreateSnapshotRequest) (*csi.CreateSnapshotResponse, error) {
	snap, err := cl.ControllerClient.CreateSnapshot(ctx, req)
	if err == nil && snap.GetSnapshot().GetSnapshotId() != "" {
		cl.registerSnapshot(snap.Snapshot.SnapshotId)
	}
	return snap, err
}

func (cl *Cleanup) deleteSnapshot(ctx context.Context, req *csi.DeleteSnapshotRequest) (*csi.DeleteSnapshotResponse, error) {
	snap, err := cl.ControllerClient.DeleteSnapshot(ctx, req)
	if err == nil && req.SnapshotId != "" {
		cl.unregisterSnapshot(req.SnapshotId)
	}
	return snap, err
}

func (cl *Cleanup) registerSnapshot(id string) {
	cl.mutex.Lock()
	defer cl.mutex.Unlock()
	cl.registerSnapshotNoLock(id)
}

func (cl *Cleanup) registerSnapshotNoLock(id string) {
	Expect(id).NotTo(BeEmpty())
	if cl.snapshots == nil {
		cl.snapshots = make(map[string]bool)
	}
	cl.snapshots[id] = true
}

func (cl *Cleanup) unregisterSnapshot(id string) {
	cl.mutex.Lock()
	defer cl.mutex.Unlock()
	cl.unregisterSnapshotNoLock(id)
}

func (cl *Cleanup) unregisterSnapshotNoLock(id string) {
	Expect(id).NotTo(BeEmpty())
	if cl.snapshots != nil {
		delete(cl.snapshots, id)
	}
}

// Cleanup calls unpublish methods as needed and deletes all volumes and
// snapshots.
func (cl *Cleanup) Cleanup() {
	cl.mutex.Lock()
	defer cl.mutex.Unlock()
	ctx := context.Background()

	cl.deleteVolumes(ctx)
	cl.deleteSnapshots(ctx)
}

func (cl *Cleanup) deleteVolumes(ctx context.Context) {
	logger := newLogger("cleanup volumes:")
	defer logger.Assert()

	for volumeID, info := range cl.volumes {
		logger.Infof("deleting %s", volumeID)
		if cl.NodeClient != nil {
			if _, err := cl.NodeUnpublishVolume(
				ctx,
				&csi.NodeUnpublishVolumeRequest{
					VolumeId:   volumeID,
					TargetPath: cl.Context.TargetPath + "/target",
				},
			); err != nil && status.Code(err) != codes.NotFound {
				logger.Errorf(err, "NodeUnpublishVolume failed: %s", err)
			}

			if cl.NodeStageSupported {
				if _, err := cl.NodeUnstageVolume(
					ctx,
					&csi.NodeUnstageVolumeRequest{
						VolumeId:          volumeID,
						StagingTargetPath: cl.Context.StagingPath,
					},
				); err != nil && status.Code(err) != codes.NotFound {
					logger.Errorf(err, "NodeUnstageVolume failed: %s", err)
				}
			}
		}

		if cl.ControllerPublishSupported && info.NodeID != "" {
			_, err := cl.ControllerClient.ControllerUnpublishVolume(
				ctx,
				&csi.ControllerUnpublishVolumeRequest{
					VolumeId: volumeID,
					NodeId:   info.NodeID,
					Secrets:  cl.Context.Secrets.ControllerUnpublishVolumeSecret,
				},
			)
			logger.Errorf(err, "ControllerUnpublishVolume failed: %s", err)
		}

		if _, err := cl.ControllerClient.DeleteVolume(
			ctx,
			&csi.DeleteVolumeRequest{
				VolumeId: volumeID,
				Secrets:  cl.Context.Secrets.DeleteVolumeSecret,
			},
		); err != nil && status.Code(err) != codes.NotFound {
			logger.Errorf(err, "DeleteVolume failed: %s", err)
		}

		cl.unregisterVolumeNoLock(volumeID)
	}
}

func (cl *Cleanup) deleteSnapshots(ctx context.Context) {
	logger := newLogger("cleanup snapshots:")
	defer logger.Assert()

	for id := range cl.snapshots {
		logger.Infof("deleting %s", id)
		_, err := cl.ControllerClient.DeleteSnapshot(
			ctx,
			&csi.DeleteSnapshotRequest{
				SnapshotId: id,
				Secrets:    cl.Context.Secrets.DeleteSnapshotSecret,
			},
		)
		logger.Errorf(err, "DeleteSnapshot failed: %s", err)

		cl.unregisterSnapshotNoLock(id)
	}
}
