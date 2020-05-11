package apitest

import (
	"testing"

	"github.com/kubernetes-csi/csi-test/v4/pkg/sanity"
)

func TestMyDriver(t *testing.T) {
	config := sanity.NewTestConfig()
	config.Address = "/tmp/e2e-csi-sanity.sock"
	config.TestNodeVolumeAttachLimit = true

	sanity.Test(t, config)
}
