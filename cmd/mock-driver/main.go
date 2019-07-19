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
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/kubernetes-csi/csi-test/driver"
	"github.com/kubernetes-csi/csi-test/mock/service"
)

func main() {
	var config service.Config
	flag.BoolVar(&config.DisableAttach, "disable-attach", false, "Disables RPC_PUBLISH_UNPUBLISH_VOLUME capability.")
	flag.StringVar(&config.DriverName, "name", service.Name, "CSI driver name.")
	flag.Int64Var(&config.AttachLimit, "attach-limit", 2, "number of attachable volumes on a node")
	flag.BoolVar(&config.NodeExpansionRequired, "node-expand-required", false, "Enables NodeServiceCapability_RPC_EXPAND_VOLUME capacity.")
	flag.BoolVar(&config.DisableControllerExpansion, "disable-controller-expansion", false, "Disables ControllerServiceCapability_RPC_EXPAND_VOLUME capability.")
	flag.BoolVar(&config.DisableOnlineExpansion, "disable-online-expansion", false, "Disables online volume expansion capability.")
	flag.BoolVar(&config.PermissiveTargetPath, "permissive-target-path", false, "Allows the CO to create PublishVolumeRequest.TargetPath, which violates the CSI spec.")
	flag.Parse()

	endpoint := os.Getenv("CSI_ENDPOINT")
	if len(endpoint) == 0 {
		fmt.Println("CSI_ENDPOINT must be defined and must be a path")
		os.Exit(1)
	}
	if strings.Contains(endpoint, ":") {
		fmt.Println("CSI_ENDPOINT must be a unix path")
		os.Exit(1)
	}

	controllerEndpoint := os.Getenv("CSI_CONTROLLER_ENDPOINT")
	if len(controllerEndpoint) == 0 {
		// If empty, set to the common endpoint.
		controllerEndpoint = endpoint
	}
	if strings.Contains(controllerEndpoint, ":") {
		fmt.Println("CSI_CONTROLLER_ENDPOINT must be a unix path")
		os.Exit(1)
	}

	// Create mock driver
	s := service.New(config)

	if endpoint == controllerEndpoint {
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
		os.Remove(endpoint)
		os.Remove(controllerEndpoint)
		l, err := net.Listen("unix", endpoint)
		if err != nil {
			fmt.Printf("Error: Unable to listen on %s socket: %v\n",
				endpoint,
				err)
			os.Exit(1)
		}
		defer os.Remove(endpoint)

		// Start server
		if err := d.Start(l); err != nil {
			fmt.Printf("Error: Unable to start mock CSI server: %v\n",
				err)
			os.Exit(1)
		}
		fmt.Println("mock driver started")

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
		fmt.Println("mock driver stopped")
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
		os.Remove(controllerEndpoint)
		l, err := net.Listen("unix", controllerEndpoint)
		if err != nil {
			fmt.Printf("Error: Unable to listen on %s socket: %v\n",
				controllerEndpoint,
				err)
			os.Exit(1)
		}
		defer os.Remove(controllerEndpoint)

		// Start controller server.
		if err = dc.Start(l); err != nil {
			fmt.Printf("Error: Unable to start mock CSI controller server: %v\n",
				err)
			os.Exit(1)
		}
		fmt.Println("mock controller driver started")

		// Listen node.
		os.Remove(endpoint)
		l, err = net.Listen("unix", endpoint)
		if err != nil {
			fmt.Printf("Error: Unable to listen on %s socket: %v\n",
				endpoint,
				err)
			os.Exit(1)
		}
		defer os.Remove(endpoint)

		// Start node server.
		if err = dn.Start(l); err != nil {
			fmt.Printf("Error: Unable to start mock CSI node server: %v\n",
				err)
			os.Exit(1)
		}
		fmt.Println("mock node driver started")

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
		fmt.Println("mock drivers stopped")
	}
}
