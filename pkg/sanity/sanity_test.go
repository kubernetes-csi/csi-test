/*
Copyright 2017 Luis Pabón luis@portworx.com

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
	"testing"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/golang/mock/gomock"
	"github.com/kubernetes-csi/csi-test/driver"
	"github.com/kubernetes-csi/csi-test/utils"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestSanity(t *testing.T) {

	m := gomock.NewController(&utils.SafeGoroutineTester{})
	defer m.Finish()
	d := driver.NewMockIdentityServer(m)

	// Setup argument match
	pluginInforeq := &csi.GetPluginInfoRequest{
		Version: &csi.Version{
			Major: 0,
			Minor: 1,
			Patch: 0,
		},
	}

	// Mock Driver Matches
	d.EXPECT().
		GetSupportedVersions(gomock.Any(), &csi.GetSupportedVersionsRequest{}).
		Return(&csi.GetSupportedVersionsResponse{
			SupportedVersions: []*csi.Version{
				{
					Major: 0,
					Minor: 1,
					Patch: 0,
				},
			},
		}, nil).
		AnyTimes()
	d.EXPECT().
		GetPluginInfo(gomock.Any(), &csi.GetPluginInfoRequest{}).
		Return(nil, status.Error(codes.InvalidArgument, "Version invalid")).
		AnyTimes()
	d.EXPECT().
		GetPluginInfo(gomock.Any(), pluginInforeq).
		Return(&csi.GetPluginInfoResponse{
			Name: "org._csi-test_.mock",
		}, nil).
		AnyTimes()

	// Setup driver
	server := driver.NewMockCSIDriver(&driver.MockCSIDriverServers{
		Identity: d,
	})
	_, err := server.Nexus()
	if err != nil {
		t.Fatalf("Unable to setup mock server: %v", err)
	}
	defer server.Close()

	// Start test
	Test(t, server.Address())
}
