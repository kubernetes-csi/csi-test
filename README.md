[![Build Status](https://k8s-testgrid.appspot.com/sig-storage-csi-other#pull-kubernetes-csi-csi-test)](https://k8s-testgrid.appspot.com/sig-storage-csi-other#pull-kubernetes-csi-csi-test)
[![Docker Repository on gcr](https://console.cloud.google.com/gcr/images/k8s-artifacts-prod/US/sig-storage/mock-driver)](https://console.cloud.google.com/gcr/images/k8s-artifacts-prod/US/sig-storage/mock-driver)

# csi-test

csi-test houses packages and libraries to help test CSI client and plugins.

## For Container Orchestration Tests

CO developers can use this framework to create drivers based on the
[Golang mock](https://github.com/golang/mock) framework. Please see
[co_test.go](test/co_test.go) for an example.

### Mock driver for testing

We also provide a container called `k8s.gcr.io/sig-storage/mock-driver` which can be used as an in-memory mock driver.
You will need to setup the environment variable `CSI_ENDPOINT` for the mock driver to know where to create the unix
domain socket.

There is an [example](https://github.com/kubernetes-csi/csi-test/tree/master/mock/example) deployment
for experiment with the mock csi driver.

For more complicated test-cases see [how to use JavaScript hooks from the driver](hooks-howto.md).

## For CSI Driver Tests

To test drivers please take a look at [pkg/sanity](https://github.com/kubernetes-csi/csi-test/tree/master/pkg/sanity).
This package and [csi-sanity](https://github.com/kubernetes-csi/csi-test/tree/master/cmd/csi-sanity) are meant to test
the CSI API capability of a driver. They are meant to be an additional test to the unit, functional, and e2e tests of a
CSI driver.

### Note

* Master is for CSI v1.3.0. Please see the branches for other CSI releases.
* Building has been tested with the Go version specified in release-tools/travis.yml

## Community, discussion, contribution, and support

Learn how to engage with the Kubernetes community on the [community page](http://kubernetes.io/community/).

You can reach the maintainers of this project at:

* [Slack channel](https://kubernetes.slack.com/messages/sig-storage)

* [Mailing list](https://groups.google.com/forum/#!forum/kubernetes-sig-storage)

### Code of conduct

Participation in the Kubernetes community is governed by the [Kubernetes Code of Conduct](code-of-conduct.md).
