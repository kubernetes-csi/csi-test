#! /bin/bash -e

CSI_PROW_TESTS="unit sanity"

. release-tools/prow.sh

# Here we override "install_sanity" to use the pre-built one.
install_sanity () {
    cp -a cmd/csi-sanity/csi-sanity "${CSI_PROW_WORK}/csi-sanity"
}

main
