#! /bin/bash -e

# A Prow job can override these defaults, but this shouldn't be necessary.

# Only these tests make sense for csi-sanity.
: ${CSI_PROW_TESTS:="unit sanity"}

# What matters for csi-sanity is the version of the hostpath driver
# that we test against, not the version of Kubernetes that it runs on.
# We pick the latest stable Kubernetes here because the corresponding
# deployment has the current driver. v1.0.1 (from the current 1.13
# deployment) does not pass csi-sanity testing.
: ${CSI_PROW_KUBERNETES_VERSION:=1.14.0}

# This repo supports and wants sanity testing.
CSI_PROW_TESTS_SANITY=sanity

. release-tools/prow.sh

# Here we override "install_sanity" to use the pre-built one.
install_sanity () {
    cp -a cmd/csi-sanity/csi-sanity "${CSI_PROW_WORK}/csi-sanity"
}

main
