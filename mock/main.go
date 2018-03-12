package main

import (
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
	endpoint := os.Getenv("CSI_ENDPOINT")
	if len(endpoint) == 0 {
		fmt.Println("CSI_ENDPOINT must be defined and must be a path")
		return
	}
	if strings.Contains(endpoint, ":") {
		fmt.Println("CSI_ENDPOINT must be a unix path")
		return
	}

	// Create mock driver
	s := service.New()
	servers := &driver.CSIDriverServers{
		Controller: s,
		Identity:   s,
		Node:       s,
	}
	d := driver.NewCSIDriver(servers)

	// Listen
	os.Remove(endpoint)
	l, err := net.Listen("unix", endpoint)
	if err != nil {
		fmt.Printf("Error: Unable to listen on %s socket: %v\n",
			endpoint,
			err)
		return
	}
	defer os.Remove(endpoint)

	// Start server
	if err := d.Start(l); err != nil {
		fmt.Printf("Error: Unable to start mock CSI server: %v\n",
			err)
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
}
