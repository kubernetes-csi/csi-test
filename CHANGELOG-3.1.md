# Changelog since v3.1.1

## Bug Fixes
 - Updates the mock driver to create the target directory in `NodePublishVolume`, and remove it in `NodeUnpublishVolume`. ([#307](https://github.com/kubernetes-csi/csi-test/pull/307), [@msau42](https://github.com/msau42))

# Changelog since v3.1.0

## Bug Fixes

- The tests that validate the result of a `CreateVolume` operation now accept zero as the `capacity_bytes` of the returned `Volume` object, when validating this value against the requested
  capacity. ([#262](https://github.com/kubernetes-csi/csi-test/pull/262), [@NicolasT](https://github.com/NicolasT))

## Other Notable Changes

- Publishing of images on k8s.gcr.io ([#270](https://github.com/kubernetes-csi/csi-test/pull/270), [@pohly](https://github.com/pohly))
- Remove missing volume test case for NodeUnpublishVolume ([#258](https://github.com/kubernetes-csi/csi-test/pull/258), [@timoreimann](https://github.com/timoreimann))
- Update release-tools to 9084fec to support native building ([#267](https://github.com/kubernetes-csi/csi-test/pull/267), [@timoreimann](https://github.com/timoreimann))

# Changelog since v3.0.0

## New Features

- "gRPCCall" logs on standard output now include "fullError" with complete error structure. Typically, it contains "code" and "message" with gRPC error. ([#254](https://github.com/kubernetes-csi/csi-test/pull/254), [@jsafrane](https://github.com/jsafrane))
- The CSI mock driver behaviour can be tweaked with JavaScript hooks in the CSI calls. ([#251](https://github.com/kubernetes-csi/csi-test/pull/251), [@tsmetana](https://github.com/tsmetana))
- Added (optional) repeated operations to test idempotency. ([#229](https://github.com/kubernetes-csi/csi-test/pull/229), [@okartau](https://github.com/okartau))
- csi-sanity `--csi.testvolumeaccesstype=block` now runs CSI tests with raw block volumes. ([#246](https://github.com/kubernetes-csi/csi-test/pull/246), [@taaraora](https://github.com/taaraora))
- Add `NodeUnpublishVolume` test for when the volume is missing. ([#242](https://github.com/kubernetes-csi/csi-test/pull/242), [@timoreimann](https://github.com/timoreimann))
- `sanity.NewTestContext` is now exported, which simplifies writing custom tests that reuse the sanity testing infrastructure. ([#253](https://github.com/kubernetes-csi/csi-test/pull/253), [@pohly](https://github.com/pohly))
- Topology support can be enabled in the CSI mock driver with `-enable-topology`. ([#249](https://github.com/kubernetes-csi/csi-test/pull/249), [@pohly](https://github.com/pohly))

## Bug Fixes

- The CSI mock driver now returns OK for requests to detach a deleted volume, as expected by [current external-attacher](https://github.com/kubernetes-csi/external-attacher/pull/165). ([#250](https://github.com/kubernetes-csi/csi-test/pull/250), [@pohly](https://github.com/pohly))
