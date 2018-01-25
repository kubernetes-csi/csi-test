#!/bin/bash

UDS="/tmp/e2e-csi-sanity.sock"
CSI_ENDPOINTS="tcp://127.0.0.1:9998"
CSI_ENDPOINTS="$CSI_ENDPOINTS unix://${UDS}"
CSI_ENDPOINTS="$CSI_ENDPOINTS ${UDS}"

go get -u github.com/thecodeteam/gocsi/mock
cd cmd/csi-sanity
  make clean install || exit 1
cd ../..

for endpoint in $CSI_ENDPOINTS ; do
    rm -f $UDS

	CSI_ENDPOINT=$endpoint mock > /dev/null 2>&1 &
	pid=$!

	csi-sanity $@ --csi.endpoint=$endpoint ; ret=$?
	kill -9 $pid
    rm -f $UDS

	if [ $ret -ne 0 ] ; then
		exit $ret
	fi
done

exit 0
