#!/bin/bash
# Copyright 2021 The Kubernetes Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

TESTARGS=$@
UDS="/tmp/e2e-csi-sanity.sock"
UDS_NODE="/tmp/e2e-csi-sanity-node.sock"
UDS_CONTROLLER="/tmp/e2e-csi-sanity-ctrl.sock"
# Protocol specified as for net.Listen...
TCP_SERVER="tcp://localhost:7654"
# ... and slightly differently for gRPC.
TCP_CLIENT="dns:///localhost:7654"
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

	# Create a script for custom target path check.
	echo '#!/bin/bash
if [ -f "$1" ]; then
    echo "file"
elif [ -d "$1" ]; then
    echo "directory"
elif [ -e "$1" ]; then
    echo "other"
else
    echo "not_found"
fi
' > custompathcheck.bash

	local creationscriptpath="$PWD/custompathcreation.bash"
	local removalscriptpath="$PWD/custompathremoval.bash"
	local checkscriptpath="$PWD/custompathcheck.bash"
	chmod +x $creationscriptpath $removalscriptpath $checkscriptpath

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
		--csi.removestagingpathcmd=$removalscriptpath \
		--csi.checkpathcmd=$checkscriptpath \
)

make

cd cmd/csi-sanity
  make clean install || exit 1
cd ../..

runTest "${TCP_SERVER}" "${TCP_CLIENT}" &&
runTest "${UDS}" "${UDS}" &&
runTestWithCreds "${UDS}" "${UDS}" &&
runTestAPI "${UDS}" &&
runTestWithDifferentAddresses "${UDS_NODE}" "${UDS_CONTROLLER}" &&
runTestAPIWithCustomTargetPaths "${UDS}" &&
runTestWithCustomTargetPaths "${UDS}" "${UDS}"
