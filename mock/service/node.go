package service

import (
	"os"
	"path"
	"strconv"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"golang.org/x/net/context"

	"github.com/container-storage-interface/spec/lib/go/csi"
)

func (s *service) NodeStageVolume(
	ctx context.Context,
	req *csi.NodeStageVolumeRequest) (
	*csi.NodeStageVolumeResponse, error) {

	device, ok := req.PublishContext["device"]
	if !ok {
		if s.config.DisableAttach {
			device = "mock device"
		} else {
			return nil, status.Error(
				codes.InvalidArgument,
				"stage volume info 'device' key required")
		}
	}

	if len(req.GetVolumeId()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Volume ID cannot be empty")
	}

	if len(req.GetStagingTargetPath()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Staging Target Path cannot be empty")
	}

	if req.GetVolumeCapability() == nil {
		return nil, status.Error(codes.InvalidArgument, "Volume Capability cannot be empty")
	}

	exists, err := checkTargetExists(req.StagingTargetPath)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if !exists {
		status.Errorf(codes.Internal, "staging target path %s does not exist", req.StagingTargetPath)
	}

	s.volsRWL.Lock()
	defer s.volsRWL.Unlock()

	i, v := s.findVolNoLock("id", req.VolumeId)
	if i < 0 {
		return nil, status.Error(codes.NotFound, req.VolumeId)
	}

	// nodeStgPathKey is the key in the volume's attributes that is set to a
	// mock stage path if the volume has been published by the node
	nodeStgPathKey := path.Join(s.nodeID, req.StagingTargetPath)

	// Check to see if the volume has already been staged.
	if v.VolumeContext[nodeStgPathKey] != "" {
		// TODO: Check for the capabilities to be equal. Return "ALREADY_EXISTS"
		// if the capabilities don't match.
		return &csi.NodeStageVolumeResponse{}, nil
	}

	// Stage the volume.
	v.VolumeContext[nodeStgPathKey] = device
	s.vols[i] = v

	return &csi.NodeStageVolumeResponse{}, nil
}

func (s *service) NodeUnstageVolume(
	ctx context.Context,
	req *csi.NodeUnstageVolumeRequest) (
	*csi.NodeUnstageVolumeResponse, error) {

	if len(req.GetVolumeId()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Volume ID cannot be empty")
	}

	if len(req.GetStagingTargetPath()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Staging Target Path cannot be empty")
	}

	s.volsRWL.Lock()
	defer s.volsRWL.Unlock()

	i, v := s.findVolNoLock("id", req.VolumeId)
	if i < 0 {
		return nil, status.Error(codes.NotFound, req.VolumeId)
	}

	// nodeStgPathKey is the key in the volume's attributes that is set to a
	// mock stage path if the volume has been published by the node
	nodeStgPathKey := path.Join(s.nodeID, req.StagingTargetPath)

	// Check to see if the volume has already been unstaged.
	if v.VolumeContext[nodeStgPathKey] == "" {
		return &csi.NodeUnstageVolumeResponse{}, nil
	}

	// Unpublish the volume.
	delete(v.VolumeContext, nodeStgPathKey)
	s.vols[i] = v

	return &csi.NodeUnstageVolumeResponse{}, nil
}

func (s *service) NodePublishVolume(
	ctx context.Context,
	req *csi.NodePublishVolumeRequest) (
	*csi.NodePublishVolumeResponse, error) {

	ephemeralVolume := req.GetVolumeContext()["csi.storage.k8s.io/ephemeral"] == "true"
	device, ok := req.PublishContext["device"]
	if !ok {
		if ephemeralVolume || s.config.DisableAttach {
			device = "mock device"
		} else {
			return nil, status.Error(
				codes.InvalidArgument,
				"stage volume info 'device' key required")
		}
	}

	if len(req.GetVolumeId()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Volume ID cannot be empty")
	}

	if len(req.GetTargetPath()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Target Path cannot be empty")
	}

	if req.GetVolumeCapability() == nil {
		return nil, status.Error(codes.InvalidArgument, "Volume Capability cannot be empty")
	}

	// May happen with old (or, at this time, even the current) Kubernetes
	// although it shouldn't (https://github.com/kubernetes/kubernetes/issues/75535).
	exists, err := checkTargetExists(req.TargetPath)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if !s.config.PermissiveTargetPath && exists {
		status.Errorf(codes.Internal, "target path %s does exist", req.TargetPath)
	}

	s.volsRWL.Lock()
	defer s.volsRWL.Unlock()

	i, v := s.findVolNoLock("id", req.VolumeId)
	if i < 0 && !ephemeralVolume {
		return nil, status.Error(codes.NotFound, req.VolumeId)
	}
	if i >= 0 && ephemeralVolume {
		return nil, status.Error(codes.AlreadyExists, req.VolumeId)
	}

	// nodeMntPathKey is the key in the volume's attributes that is set to a
	// mock mount path if the volume has been published by the node
	nodeMntPathKey := path.Join(s.nodeID, req.TargetPath)

	// Check to see if the volume has already been published.
	if v.VolumeContext[nodeMntPathKey] != "" {

		// Requests marked Readonly fail due to volumes published by
		// the Mock driver supporting only RW mode.
		if req.Readonly {
			return nil, status.Error(codes.AlreadyExists, req.VolumeId)
		}

		return &csi.NodePublishVolumeResponse{}, nil
	}

	// Publish the volume.
	if ephemeralVolume {
		MockVolumes[req.VolumeId] = Volume{
			ISEphemeral: true,
		}
	} else {
		if req.GetStagingTargetPath() != "" {
			exists, err := checkTargetExists(req.GetStagingTargetPath())
			if err != nil {
				return nil, status.Error(codes.Internal, err.Error())
			}
			if !exists {
				status.Errorf(codes.Internal, "staging target path %s does not exist", req.GetStagingTargetPath())
			}
			v.VolumeContext[nodeMntPathKey] = req.GetStagingTargetPath()
		} else {
			v.VolumeContext[nodeMntPathKey] = device
		}
		s.vols[i] = v
	}

	return &csi.NodePublishVolumeResponse{}, nil
}

func (s *service) NodeUnpublishVolume(
	ctx context.Context,
	req *csi.NodeUnpublishVolumeRequest) (
	*csi.NodeUnpublishVolumeResponse, error) {

	if len(req.GetVolumeId()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Volume ID cannot be empty")
	}
	if len(req.GetTargetPath()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Target Path cannot be empty")
	}

	s.volsRWL.Lock()
	defer s.volsRWL.Unlock()

	ephemeralVolume := MockVolumes[req.VolumeId].ISEphemeral
	i, v := s.findVolNoLock("id", req.VolumeId)
	if i < 0 && !ephemeralVolume {
		return nil, status.Error(codes.NotFound, req.VolumeId)
	}

	if ephemeralVolume {
		delete(MockVolumes, req.VolumeId)
	} else {
		// nodeMntPathKey is the key in the volume's attributes that is set to a
		// mock mount path if the volume has been published by the node
		nodeMntPathKey := path.Join(s.nodeID, req.TargetPath)

		// Check to see if the volume has already been unpublished.
		if v.VolumeContext[nodeMntPathKey] == "" {
			return &csi.NodeUnpublishVolumeResponse{}, nil
		}

		// Unpublish the volume.
		delete(v.VolumeContext, nodeMntPathKey)
		s.vols[i] = v
	}

	return &csi.NodeUnpublishVolumeResponse{}, nil
}

func (s *service) NodeExpandVolume(ctx context.Context, req *csi.NodeExpandVolumeRequest) (*csi.NodeExpandVolumeResponse, error) {
	if len(req.GetVolumeId()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Volume ID cannot be empty")
	}
	if len(req.GetVolumePath()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Volume Path cannot be empty")
	}
	if req.GetCapacityRange() == nil {
		return nil, status.Error(codes.InvalidArgument, "Request capacity cannot be empty")
	}

	s.volsRWL.Lock()
	defer s.volsRWL.Unlock()

	i, v := s.findVolNoLock("id", req.VolumeId)
	if i < 0 {
		return nil, status.Error(codes.NotFound, req.VolumeId)
	}

	// TODO: NodeExpandVolume MUST be called after successful NodeStageVolume as we has STAGE_UNSTAGE_VOLUME node capacity.

	requestCapacity := req.GetCapacityRange().RequiredBytes
	resp := &csi.NodeExpandVolumeResponse{CapacityBytes: requestCapacity}

	// fsCapacityKey is the key in the volume's attributes that is set to the file system's size.
	fsCapacityKey := path.Join(s.nodeID, req.GetVolumePath(), "size")
	oldCapacityStr, exist := v.VolumeContext[fsCapacityKey]
	if exist {
		oldCapacity, err := strconv.ParseInt(oldCapacityStr, 10, 64)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
		if oldCapacity > requestCapacity {
			return nil, status.Error(codes.InvalidArgument, "cannot change file system size to a smaller size")
		}
		if oldCapacity == requestCapacity {
			// File system capacity is equal to requested size, no need to expand.
			return resp, nil
		}
	}

	// Update volume's fs capacity to requested size.
	v.VolumeContext[fsCapacityKey] = strconv.FormatInt(requestCapacity, 10)
	s.vols[i] = v

	return resp, nil
}

func (s *service) NodeGetCapabilities(
	ctx context.Context,
	req *csi.NodeGetCapabilitiesRequest) (
	*csi.NodeGetCapabilitiesResponse, error) {

	capabilities := []*csi.NodeServiceCapability{
		{
			Type: &csi.NodeServiceCapability_Rpc{
				Rpc: &csi.NodeServiceCapability_RPC{
					Type: csi.NodeServiceCapability_RPC_UNKNOWN,
				},
			},
		},
		{
			Type: &csi.NodeServiceCapability_Rpc{
				Rpc: &csi.NodeServiceCapability_RPC{
					Type: csi.NodeServiceCapability_RPC_STAGE_UNSTAGE_VOLUME,
				},
			},
		},
		{
			Type: &csi.NodeServiceCapability_Rpc{
				Rpc: &csi.NodeServiceCapability_RPC{
					Type: csi.NodeServiceCapability_RPC_GET_VOLUME_STATS,
				},
			},
		},
	}
	if s.config.NodeExpansionRequired {
		capabilities = append(capabilities, &csi.NodeServiceCapability{
			Type: &csi.NodeServiceCapability_Rpc{
				Rpc: &csi.NodeServiceCapability_RPC{
					Type: csi.NodeServiceCapability_RPC_EXPAND_VOLUME,
				},
			},
		})
	}

	return &csi.NodeGetCapabilitiesResponse{
		Capabilities: capabilities,
	}, nil
}

func (s *service) NodeGetInfo(ctx context.Context,
	req *csi.NodeGetInfoRequest) (*csi.NodeGetInfoResponse, error) {
	csiNodeResponse := &csi.NodeGetInfoResponse{
		NodeId: s.nodeID,
	}
	if s.config.AttachLimit > 0 {
		csiNodeResponse.MaxVolumesPerNode = s.config.AttachLimit
	}
	return csiNodeResponse, nil
}

func (s *service) NodeGetVolumeStats(ctx context.Context,
	req *csi.NodeGetVolumeStatsRequest) (*csi.NodeGetVolumeStatsResponse, error) {

	if len(req.GetVolumeId()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Volume ID cannot be empty")
	}

	if len(req.GetVolumePath()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Volume Path cannot be empty")
	}

	i, v := s.findVolNoLock("id", req.VolumeId)
	if i < 0 {
		return nil, status.Error(codes.NotFound, req.VolumeId)
	}

	nodeMntPathKey := path.Join(s.nodeID, req.VolumePath)

	_, exists := v.VolumeContext[nodeMntPathKey]
	if !exists {
		return nil, status.Errorf(codes.NotFound, "volume %q doest not exist on the specified path %q", req.VolumeId, req.VolumeId)
	}

	return &csi.NodeGetVolumeStatsResponse{
		Usage: []*csi.VolumeUsage{
			{
				Total: v.GetCapacityBytes(),
				Unit:  csi.VolumeUsage_BYTES,
			},
		},
	}, nil
}

// checkTargetExists checks if a given path exists.
func checkTargetExists(targetPath string) (bool, error) {
	_, err := os.Stat(targetPath)
	switch {
	case err == nil:
		return true, nil
	case os.IsNotExist(err):
		return false, nil
	default:
		return false, err
	}
}
