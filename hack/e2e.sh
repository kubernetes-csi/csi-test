#!/bin/bash

CSI_ENDPOINTS="tcp://127.0.0.1:9998"
#CSI_ENDPOINTS="$CSI_ENDPOINTS /tmp/csi.sock"
#CSI_ENDPOINTS="$CSI_ENDPOINTS unix:///tmp/csi.sock"

if [ ! -x $GOPATH/bin/mock ] ; then
	go get -u github.com/thecodeteam/gocsi/mock
fi

cd cmd/csi-sanity
  make clean install || exit 1
cd ../..

for endpoint in $CSI_ENDPOINTS ; do
	CSI_ENDPOINT=$endpoint mock &
	pid=$!

	csi-sanity $@ --csi.endpoint=$endpoint ; ret=$?
	kill -9 $pid

	if [ $ret -ne 0 ] ; then
		exit $ret
	fi
done

exit 0
