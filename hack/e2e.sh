#!/bin/bash

TESTARGS=$@
UDS="/tmp/e2e-csi-sanity.sock"
UDS_NODE="/tmp/e2e-csi-sanity-node.sock"
UDS_CONTROLLER="/tmp/e2e-csi-sanity-ctrl.sock"
CSI_ENDPOINTS="$CSI_ENDPOINTS ${UDS}"
CSI_MOCK_VERSION="master"

# cleanup mock_driver_pid files...
cleanup () {
    local pid="$1"
    shift
    kill -9 "$pid"
    # We don't care about the 'hack/e2e.sh: line 15: 117018 Killed                  CSI_ENDPOINT=$1 ./bin/mock-driver'
    wait "$pid" 2>/dev/null
    rm -f "$@"
}

#
# $1 - endpoint for mock.
# $2 - endpoint for csi-sanity in Grpc format.
#      See https://github.com/grpc/grpc/blob/master/doc/naming.md
runTest()
(
	CSI_ENDPOINT=$1 ./bin/mock-driver &
	local pid=$!
        trap 'cleanup $pid $1' EXIT

	./cmd/csi-sanity/csi-sanity $TESTARGS --csi.endpoint=$2 --csi.testnodevolumeattachlimit
)

runTestWithDifferentAddresses()
(
	CSI_ENDPOINT=$1 CSI_CONTROLLER_ENDPOINT=$2 ./bin/mock-driver &
	local pid=$!
        trap 'cleanup $pid $1' EXIT

	./cmd/csi-sanity/csi-sanity $TESTARGS --csi.endpoint=$1 --csi.controllerendpoint=$2
)

runTestWithCreds()
(
	CSI_ENDPOINT=$1 CSI_ENABLE_CREDS=true ./bin/mock-driver &
	local pid=$!
        trap 'cleanup $pid $1' EXIT

	./cmd/csi-sanity/csi-sanity $TESTARGS --csi.endpoint=$2 --csi.secrets=mock/mocksecret.yaml --csi.testnodevolumeattachlimit
)

runTestAPI()
(
	CSI_ENDPOINT=$1 ./bin/mock-driver &
	local pid=$!
        trap 'cleanup $pid $1' EXIT

	go test -count=1 -v ./hack/_apitest/api_test.go && \
	go test -count=1 -v ./hack/_embedded/embedded_test.go
)

runTestAPIWithCustomTargetPaths()
(
	CSI_ENDPOINT=$1 ./bin/mock-driver &
	local pid=$!
        trap 'cleanup $pid $1' EXIT

	# Running a specific test to verify that the custom target paths are called
	# a deterministic number of times.
	go test -count=1 -v ./hack/_apitest2/api_test.go -ginkgo.focus="NodePublishVolume"
)

runTestWithCustomTargetPaths()
(

	# Create a script for custom target path creation.
	echo '#!/bin/bash
targetpath="/tmp/csi/$@"
mkdir -p $targetpath
echo $targetpath
' > custompathcreation.bash

	# Create a script for custom target path removal.
	echo '#!/bin/bash
rm -rf $@
' > custompathremoval.bash

	chmod +x custompathcreation.bash custompathremoval.bash
	local creationscriptpath="$PWD/custompathcreation.bash"
	local removalscriptpath="$PWD/custompathremoval.bash"

	CSI_ENDPOINT=$1 ./bin/mock-driver &
	local pid=$!
        trap 'cleanup $pid $1; rm $creationscriptpath $removalscriptpath' EXIT

	./cmd/csi-sanity/csi-sanity $TESTARGS \
		--csi.endpoint=$2 \
		--csi.mountdir="foo/target/mount" \
		--csi.stagingdir="foo/staging/mount" \
		--csi.createmountpathcmd=$creationscriptpath \
		--csi.createstagingpathcmd=$creationscriptpath \
		--csi.removemountpathcmd=$removalscriptpath \
		--csi.removestagingpathcmd=$removalscriptpath
)

make

cd cmd/csi-sanity
  make clean install || exit 1
cd ../..

runTest "${UDS}" "${UDS}" &&
runTestWithCreds "${UDS}" "${UDS}" &&
runTestAPI "${UDS}" &&
runTestWithDifferentAddresses "${UDS_NODE}" "${UDS_CONTROLLER}" &&
runTestAPIWithCustomTargetPaths "${UDS}" &&
runTestWithCustomTargetPaths "${UDS}" "${UDS}"
