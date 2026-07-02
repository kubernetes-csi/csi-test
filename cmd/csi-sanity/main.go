/*
Copyright 2021 The Kubernetes Authors.

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
package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/onsi/ginkgo/v2"
	"k8s.io/klog/v2"

	"github.com/kubernetes-csi/csi-test/v5/pkg/sanity"
)

const (
	prefix string = "csi."
)

var (
	VERSION = "(dev)"
)

func stringVar(p *string, name string, usage string) {
	flag.StringVar(p, prefix+name, *p, usage)
}

func boolVar(p *bool, name string, usage string) {
	flag.BoolVar(p, prefix+name, *p, usage)
}

func intVar(p *int, name string, usage string) {
	flag.IntVar(p, prefix+name, *p, usage)
}

func int64Var(p *int64, name string, usage string) {
	flag.Int64Var(p, prefix+name, *p, usage)
}

func durationVar(p *time.Duration, name string, usage string) {
	flag.DurationVar(p, prefix+name, *p, usage)
}

type testing struct {
	result int
}

func (t *testing) Fail() {
	t.result = 1
}

func main() {
	version := flag.Bool("version", false, "print version of this program")

	// Get configuration with defaults.
	config := sanity.NewTestConfig()

	// Support overriding the default configuration via flags.
	stringVar(&config.Address, "endpoint", "CSI endpoint")
	stringVar(&config.ControllerAddress, "controllerendpoint", "CSI controller endpoint")
	stringVar(&config.TargetPath, "mountdir", "Mount point for NodePublish")
	stringVar(&config.StagingPath, "stagingdir", "Mount point for NodeStage if staging is supported")
	stringVar(&config.CreateTargetPathCmd, "createmountpathcmd", "Command to run for target path creation")
	stringVar(&config.CreateStagingPathCmd, "createstagingpathcmd", "Command to run for staging path creation")
	durationVar(&config.CreatePathCmdTimeout, "createpathcmdtimeout", "Timeout for the commands to create target and staging paths, in seconds")
	stringVar(&config.RemoveTargetPathCmd, "removemountpathcmd", "Command to run for target path removal")
	stringVar(&config.RemoveStagingPathCmd, "removestagingpathcmd", "Command to run for staging path removal")
	durationVar(&config.RemovePathCmdTimeout, "removepathcmdtimeout", "Timeout for the commands to remove target and staging paths, in seconds")
	stringVar(&config.CheckPathCmd, "checkpathcmd", "Command to run to check a given path. It must print 'file', 'directory', 'not_found', or 'other' on stdout.")
	durationVar(&config.CheckPathCmdTimeout, "checkpathcmdtimeout", "Timeout for the command to check a given path, in seconds")
	stringVar(&config.SecretsFile, "secrets", "CSI secrets file")
	stringVar(&config.TestVolumeAccessType, "testvolumeaccesstype", "Volume capability access type, valid values are mount or block")
	int64Var(&config.TestVolumeSize, "testvolumesize", "Base volume size used for provisioned volumes")
	int64Var(&config.TestVolumeExpandSize, "testvolumeexpandsize", "Target size for expanded volumes")
	stringVar(&config.TestVolumeParametersFile, "testvolumeparameters", "YAML file of volume parameters for provisioned volumes")
	stringVar(&config.TestVolumeMutableParametersFile, "testvolumemutableparameters", "YAML file of mutable parameters for modifying volumes")
	stringVar(&config.TestSnapshotParametersFile, "testsnapshotparameters", "YAML file of snapshot parameters for provisioned snapshots")
	stringVar(&config.TestTopologyRequirementsFile, "testtopologyrequirements", "YAML file of topology requirements for provisioned volumes")
	boolVar(&config.TestNodeVolumeAttachLimit, "testnodevolumeattachlimit", "Test node volume attach limit")
	flag.Var(flag.Lookup("ginkgo.junit-report").Value, prefix+"junitfile", "JUnit XML output file where test results will be written (deprecated: use ginkgo.junit-report instead)")

	flag.Parse()
	if *version {
		fmt.Printf("Version = %s\n", VERSION)
		os.Exit(0)
	}
	if config.Address == "" {
		fmt.Printf("--%sendpoint must be provided with an CSI endpoint\n", prefix)
		os.Exit(1)
	}
	if at := strings.TrimSpace(strings.ToLower(config.TestVolumeAccessType)); !(at == "mount" || at == "block") {
		fmt.Printf("--%stestvolumeaccesstype valid values are mount or block\n", prefix)
		os.Exit(1)
	}

	klog.SetOutput(ginkgo.GinkgoWriter)
	t := testing{}
	sanity.Test(&t, config)
	os.Exit(t.result)
}
