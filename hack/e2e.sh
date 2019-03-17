#!/bin/bash

TESTARGS=$@
UDS="/tmp/e2e-csi-sanity.sock"
UDS_NODE="/tmp/e2e-csi-sanity-node.sock"
UDS_CONTROLLER="/tmp/e2e-csi-sanity-ctrl.sock"
CSI_ENDPOINTS="$CSI_ENDPOINTS ${UDS}"
CSI_MOCK_VERSION="master"

#
# $1 - endpoint for mock.
# $2 - endpoint for csi-sanity in Grpc format.
#      See https://github.com/grpc/grpc/blob/master/doc/naming.md
runTest()
{
	CSI_ENDPOINT=$1 ./bin/mock-driver &
	local pid=$!

	./cmd/csi-sanity/csi-sanity $TESTARGS --csi.endpoint=$2 --csi.testnodevolumeattachlimit; ret=$?
	kill -9 $pid

	if [ $ret -ne 0 ] ; then
		exit $ret
	fi
}

runTestWithDifferentAddresses()
{
	CSI_ENDPOINT=$1 CSI_CONTROLLER_ENDPOINT=$2 ./bin/mock-driver &
	local pid=$!

	./cmd/csi-sanity/csi-sanity $TESTARGS --csi.endpoint=$1 --csi.controllerendpoint=$2; ret=$?
	kill -9 $pid

	if [ $ret -ne 0 ] ; then
		exit $ret
	fi
}

runTestWithCreds()
{
	CSI_ENDPOINT=$1 CSI_ENABLE_CREDS=true ./bin/mock-driver &
	local pid=$!

	./cmd/csi-sanity/csi-sanity $TESTARGS --csi.endpoint=$2 --csi.secrets=mock/mocksecret.yaml --csi.testnodevolumeattachlimit; ret=$?
	kill -9 $pid

	if [ $ret -ne 0 ] ; then
		exit $ret
	fi
}

runTestAPI()
{
	CSI_ENDPOINT=$1 ./bin/mock-driver &
	local pid=$!

	GOCACHE=off go test -v ./hack/_apitest/api_test.go; ret=$?

	if [ $ret -ne 0 ] ; then
		exit $ret
	fi

	GOCACHE=off go test -v ./hack/_embedded/embedded_test.go; ret=$?
	kill -9 $pid

	if [ $ret -ne 0 ] ; then
		exit $ret
	fi
}

runTestAPIWithCustomTargetPaths()
{
	CSI_ENDPOINT=$1 ./bin/mock-driver &
	local pid=$!

	# Running a specific test to verify that the custom target paths are called
	# a deterministic number of times.
	GOCACHE=off go test -v ./hack/_apitest2/api_test.go -ginkgo.focus="NodePublishVolume"; ret=$?

	if [ $ret -ne 0 ] ; then
		exit $ret
	fi
}

runTestWithCustomTargetPaths()
{
	CSI_ENDPOINT=$1 ./bin/mock-driver &
	local pid=$!

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

	./cmd/csi-sanity/csi-sanity $TESTARGS \
		--csi.endpoint=$2 \
		--csi.mountdir="foo/target/mount" \
		--csi.stagingdir="foo/staging/mount" \
		--csi.createmountpathcmd=$creationscriptpath \
		--csi.createstagingpathcmd=$creationscriptpath \
		--csi.removemountpathcmd=$removalscriptpath \
		--csi.removestagingpathcmd=$removalscriptpath; ret=$?
	kill -9 $pid

	# Delete the script.
	rm $creationscriptpath $removalscriptpath

	if [ $ret -ne 0 ] ; then
		exit $ret
	fi
}

make

cd cmd/csi-sanity
  make clean install || exit 1
cd ../..

runTest "${UDS}" "${UDS}"
rm -f $UDS

runTestWithCreds "${UDS}" "${UDS}"
rm -f $UDS

runTestAPI "${UDS}"
rm -f $UDS

runTestWithDifferentAddresses "${UDS_NODE}" "${UDS_CONTROLLER}"
rm -f $UDS_NODE
rm -f $UDS_CONTROLLER

runTestAPIWithCustomTargetPaths "${UDS}"
rm -rf $UDS

runTestWithCustomTargetPaths "${UDS}" "${UDS}"
rm -rf $UDS

exit 0
