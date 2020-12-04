/*
Copyright 2018 Kubernetes Authors

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
	"context"
	"flag"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"

	"github.com/kubernetes-csi/csi-test/v4/driver"
	"github.com/kubernetes-csi/csi-test/v4/internal/endpoint"
	"github.com/kubernetes-csi/csi-test/v4/internal/proxy"
	"github.com/kubernetes-csi/csi-test/v4/mock/service"
	"gopkg.in/yaml.v2"
	"k8s.io/klog/v2"
)

func init() {
	// klog verbosity guide for this package
	// Use V(2) for one time config information
	// Use V(3) for general debug information logging
	klog.InitFlags(flag.CommandLine)
}

func main() {
	var config service.Config
	var hooksFile string = ""
	flag.BoolVar(&config.DisableAttach, "disable-attach", false, "Disables RPC_PUBLISH_UNPUBLISH_VOLUME capability.")
	flag.StringVar(&config.DriverName, "name", service.Name, "CSI driver name.")
	flag.Int64Var(&config.AttachLimit, "attach-limit", 2, "number of attachable volumes on a node")
	flag.BoolVar(&config.NodeExpansionRequired, "node-expand-required", false, "Enables NodeServiceCapability_RPC_EXPAND_VOLUME capacity.")
	flag.BoolVar(&config.EnableTopology, "enable-topology", false, "Enables PluginCapability_Service_VOLUME_ACCESSIBILITY_CONSTRAINTS capability.")
	flag.BoolVar(&config.DisableControllerExpansion, "disable-controller-expansion", false, "Disables ControllerServiceCapability_RPC_EXPAND_VOLUME capability.")
	flag.BoolVar(&config.DisableOnlineExpansion, "disable-online-expansion", false, "Disables online volume expansion capability.")
	flag.BoolVar(&config.PermissiveTargetPath, "permissive-target-path", false, "Allows the CO to create PublishVolumeRequest.TargetPath, which violates the CSI spec.")
	flag.StringVar(&hooksFile, "hooks-file", "", "YAML file with hook scripts.")
	proxyEndpoint := flag.String("proxy-endpoint", "", "Instead of running the CSI driver code, just proxy connections from $CSI_ENDPOINT to the given listening socket.")
	flag.Parse()

	csiEndpoint := os.Getenv("CSI_ENDPOINT")
	controllerEndpoint := os.Getenv("CSI_CONTROLLER_ENDPOINT")
	if len(controllerEndpoint) == 0 {
		// If empty, set to the common endpoint.
		controllerEndpoint = csiEndpoint
	}

	if *proxyEndpoint != "" {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		closer, err := proxy.Run(ctx, csiEndpoint, *proxyEndpoint)
		if err != nil {
			klog.Fatalf("failed to run proxy: %v", err)
		}
		defer closer.Close()

		// Wait for signal
		sigc := make(chan os.Signal, 1)
		sigs := []os.Signal{
			syscall.SIGTERM,
			syscall.SIGHUP,
			syscall.SIGINT,
			syscall.SIGQUIT,
		}
		signal.Notify(sigc, sigs...)

		<-sigc
		return
	}

	if hooksFile != "" {
		execHooks, err := parseHooksFile(hooksFile)
		if err == nil {
			config.ExecHooks = execHooks
		} else {
			klog.Errorf("Failed to load hooks file %s: %v", hooksFile, err)
		}
	}

	// Create mock driver
	s := service.New(config)

	if csiEndpoint == controllerEndpoint {
		servers := &driver.CSIDriverServers{
			Controller: s,
			Identity:   s,
			Node:       s,
		}
		d := driver.NewCSIDriver(servers)

		// If creds is enabled, set the default creds.
		setCreds := os.Getenv("CSI_ENABLE_CREDS")
		if len(setCreds) > 0 && setCreds == "true" {
			d.SetDefaultCreds()
		}

		// Listen
		l, cleanup, err := endpoint.Listen(csiEndpoint)
		if err != nil {
			klog.Exitf("Error: Unable to listen on %s socket: %v\n",
				csiEndpoint,
				err)
		}
		defer cleanup()

		// Start server
		if err := d.Start(l); err != nil {
			klog.Exitf("Error: Unable to start mock CSI server: %v\n",
				err)
		}
		klog.Info("mock driver started")

		// Wait for signal
		sigc := make(chan os.Signal, 1)
		sigs := []os.Signal{
			syscall.SIGTERM,
			syscall.SIGHUP,
			syscall.SIGINT,
			syscall.SIGQUIT,
		}
		signal.Notify(sigc, sigs...)

		<-sigc
		d.Stop()
		klog.Info("mock driver stopped")
	} else {
		controllerServer := &driver.CSIDriverControllerServer{
			Controller: s,
			Identity:   s,
		}
		dc := driver.NewCSIDriverController(controllerServer)

		nodeServer := &driver.CSIDriverNodeServer{
			Node:     s,
			Identity: s,
		}
		dn := driver.NewCSIDriverNode(nodeServer)

		setCreds := os.Getenv("CSI_ENABLE_CREDS")
		if len(setCreds) > 0 && setCreds == "true" {
			dc.SetDefaultCreds()
			dn.SetDefaultCreds()
		}

		// Listen controller.
		l, cleanupController, err := endpoint.Listen(controllerEndpoint)
		if err != nil {
			klog.Exitf("Error: Unable to listen on %s socket: %v\n",
				controllerEndpoint,
				err)
		}
		defer cleanupController()

		// Start controller server.
		if err = dc.Start(l); err != nil {
			klog.Exitf("Error: Unable to start mock CSI controller server: %v\n",
				err)
		}
		klog.Infof("mock controller driver started")

		// Listen node.
		l, cleanupNode, err := endpoint.Listen(csiEndpoint)
		if err != nil {
			klog.Exitf("Error: Unable to listen on %s socket: %v\n",
				csiEndpoint,
				err)
		}
		defer cleanupNode()

		// Start node server.
		if err = dn.Start(l); err != nil {
			klog.Exitf("Error: Unable to start mock CSI node server: %v\n",
				err)
		}
		klog.Infof("mock node driver started")

		// Wait for signal
		sigc := make(chan os.Signal, 1)
		sigs := []os.Signal{
			syscall.SIGTERM,
			syscall.SIGHUP,
			syscall.SIGINT,
			syscall.SIGQUIT,
		}
		signal.Notify(sigc, sigs...)

		<-sigc
		dc.Stop()
		dn.Stop()
		klog.Infof("mock drivers stopped")
	}
}

func parseHooksFile(file string) (*service.Hooks, error) {
	var hooks service.Hooks

	fr, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer fr.Close()
	bytes, _ := ioutil.ReadAll(fr)
	err = yaml.UnmarshalStrict([]byte(bytes), &hooks)
	if err != nil {
		return nil, err
	}
	klog.V(2).Infof("Hooks file %s loaded\n", file)
	return &hooks, err
}
