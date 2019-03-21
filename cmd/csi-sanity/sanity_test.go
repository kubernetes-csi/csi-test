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
	"flag"
	"fmt"
	"os"
	"testing"

	"github.com/kubernetes-csi/csi-test/pkg/sanity"
)

const (
	prefix string = "csi."
)

var (
	VERSION = "(dev)"
	version bool
	config  sanity.Config
)

func init() {
	flag.StringVar(&config.Address, prefix+"endpoint", "", "CSI endpoint")
	flag.StringVar(&config.ControllerAddress, prefix+"controllerendpoint", "", "CSI controller endpoint")
	flag.BoolVar(&version, prefix+"version", false, "Version of this program")
	flag.StringVar(&config.TargetPath, prefix+"mountdir", os.TempDir()+"/csi-mount", "Mount point for NodePublish")
	flag.StringVar(&config.StagingPath, prefix+"stagingdir", os.TempDir()+"/csi-staging", "Mount point for NodeStage if staging is supported")
	flag.StringVar(&config.CreateTargetPathCmd, prefix+"createmountpathcmd", "", "Command to run for target path creation")
	flag.StringVar(&config.CreateStagingPathCmd, prefix+"createstagingpathcmd", "", "Command to run for staging path creation")
	flag.IntVar(&config.CreatePathCmdTimeout, prefix+"createpathcmdtimeout", 10, "Timeout for the commands to create target and staging paths, in seconds")
	flag.StringVar(&config.RemoveTargetPathCmd, prefix+"removemountpathcmd", "", "Command to run for target path removal")
	flag.StringVar(&config.RemoveStagingPathCmd, prefix+"removestagingpathcmd", "", "Command to run for staging path removal")
	flag.IntVar(&config.RemovePathCmdTimeout, prefix+"removepathcmdtimeout", 10, "Timeout for the commands to remove target and staging paths, in seconds")
	flag.StringVar(&config.SecretsFile, prefix+"secrets", "", "CSI secrets file")
	flag.Int64Var(&config.TestVolumeSize, prefix+"testvolumesize", sanity.DefTestVolumeSize, "Base volume size used for provisioned volumes")
	flag.StringVar(&config.TestVolumeParametersFile, prefix+"testvolumeparameters", "", "YAML file of volume parameters for provisioned volumes")
	flag.BoolVar(&config.TestNodeVolumeAttachLimit, prefix+"testnodevolumeattachlimit", false, "Test node volume attach limit")
	flag.StringVar(&config.JUnitFile, prefix+"junitfile", "", "JUnit XML output file where test results will be written")
	flag.Parse()
}

func TestSanity(t *testing.T) {
	if version {
		fmt.Printf("Version = %s\n", VERSION)
		return
	}
	if len(config.Address) == 0 {
		t.Fatalf("--%sendpoint must be provided with an CSI endpoint", prefix)
	}
	sanity.Test(t, &config)
}
