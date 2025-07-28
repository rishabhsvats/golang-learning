package driver

import (
	"context"
	"fmt"
	"strconv"

	"github.com/container-storage-interface/spec/lib/go/csi"
	metadata "github.com/digitalocean/go-metadata"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (d *Driver) NodeStageVolume(ctx context.Context, req *csi.NodeStageVolumeRequest) (*csi.NodeStageVolumeResponse, error) {
	// make sure all the req fields are present
	if req.VolumeId == "" {
		return nil, status.Error(codes.InvalidArgument, "VolumeID must be present in the NodeStageVolumeReq")
	}

	if req.StagingTargetPath == "" {
		return nil, status.Error(codes.InvalidArgument, "StagingTargetPath must be present in the NodeSVolReq")
	}

	if req.VolumeCapability == nil {
		return nil, status.Error(codes.InvalidArgument, "VolumeCaps must be present in the NodeSVolReq")
	}
	switch req.VolumeCapability.AccessType.(type) {
	case *csi.VolumeCapability_Block:
		return &csi.NodeStageVolumeResponse{}, nil
	}

	volumeName := ""
	if val, ok := req.PublishContext[volNameKeyFromContPub]; !ok {
		return nil, status.Error(codes.InvalidArgument, "Volumename is not present in the publish context of request")
	} else {
		volumeName = val
	}

	mnt := req.VolumeCapability.GetMount()
	fsType := "ext4"
	if mnt.FsType != "" {
		fsType = mnt.FsType
	}

	// figure out the source and target
	source := getPathFromVolumeName(volumeName)
	target := req.StagingTargetPath
	return nil, nil
}

func getPathFromVolumeName(volName string) string {
	return fmt.Sprintf("/dev/disk/by-id/scsi-0DO_Volume_%s", volName)
}
func (d *Driver) NodeUnstageVolume(context.Context, *csi.NodeUnstageVolumeRequest) (*csi.NodeUnstageVolumeResponse, error) {
	return nil, nil
}
func (d *Driver) NodePublishVolume(context.Context, *csi.NodePublishVolumeRequest) (*csi.NodePublishVolumeResponse, error) {
	return nil, nil
}
func (d *Driver) NodeUnpublishVolume(context.Context, *csi.NodeUnpublishVolumeRequest) (*csi.NodeUnpublishVolumeResponse, error) {
	return nil, nil
}
func (d *Driver) NodeGetVolumeStats(context.Context, *csi.NodeGetVolumeStatsRequest) (*csi.NodeGetVolumeStatsResponse, error) {
	return nil, nil
}
func (d *Driver) NodeExpandVolume(context.Context, *csi.NodeExpandVolumeRequest) (*csi.NodeExpandVolumeResponse, error) {
	return nil, nil
}
func (d *Driver) NodeGetCapabilities(context.Context, *csi.NodeGetCapabilitiesRequest) (*csi.NodeGetCapabilitiesResponse, error) {
	fmt.Println("NodeGetCapabilities of the node service was called")

	return &csi.NodeGetCapabilitiesResponse{
		Capabilities: []*csi.NodeServiceCapability{
			{
				Type: &csi.NodeServiceCapability_Rpc{
					Rpc: &csi.NodeServiceCapability_RPC{
						Type: csi.NodeServiceCapability_RPC_STAGE_UNSTAGE_VOLUME,
					},
				},
			},
		},
	}, nil
}
func (d *Driver) NodeGetInfo(context.Context, *csi.NodeGetInfoRequest) (*csi.NodeGetInfoResponse, error) {
	mdClient := metadata.NewClient()
	id, err := mdClient.DropletID()
	if err != nil {
		return nil, status.Error(codes.Internal, "Error getting nodeID")
	}
	return &csi.NodeGetInfoResponse{
		NodeId:            strconv.Itoa(id),
		MaxVolumesPerNode: 5,
		AccessibleTopology: &csi.Topology{
			Segments: map[string]string{
				"region": "ams3",
			},
		},
	}, nil
}
