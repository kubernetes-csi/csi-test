[![Build Status](https://travis-ci.org/kubernetes-csi/csi-test.svg?branch=master)](https://travis-ci.org/kubernetes-csi/csi-test)
# csi-test
csi-test is a golang unit testing framework for container orchestration (CO)
system and CSI driver developers.

### For Container Orchestration Unit Tests
CO developers can use this framework to create drivers based on the
[Golang mock](https://github.com/golang/mock) framework. Please see
[co_test.go](test/co_test.go) for an example.

### For CSI Driver Unit Tests
Driver developers do not need to leverage the mocking framework, and
instead just use the CSI protocol buffers golang output library. This
framework may provide little value currently, but if the need arises,
it may provide future libraries to make developement and testing of
drivers easier. Please see the example [driver_test.go](test/driver_test.go)
for more information. 

### Note

* Only Golang 1.9+ supported. See (gRPC issue)[https://github.com/grpc/grpc-go/issues/711#issuecomment-326626790]
