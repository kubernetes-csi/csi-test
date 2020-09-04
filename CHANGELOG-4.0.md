# Changelog since v4.0.0

### Other Notable Changes

- Restore snapshots-before-volumes cleanup order ([#289](https://github.com/kubernetes-csi/csi-test/pull/289), [@timoreimann](https://github.com/timoreimann))

# Changelog since v3.1.0

## New Features

- Replace logging to klog for mock CSI driver. ([#279](https://github.com/kubernetes-csi/csi-test/pull/279), [@Jiawei0227](https://github.com/Jiawei0227))
- Add CSI LIST_VOLUMES_PUBLISHED_NODES capability to csi-mock ([#280](https://github.com/kubernetes-csi/csi-test/pull/280), [@yuga711](https://github.com/yuga711))
- The flag --csi.testvolumeaccesstype=block is now also used by the following tests: ExpandVolume [Controller Server], CreateSnapshot [Controller Server], DeleteSnapshot [Controller Server], ListSnapshots [Controller Server], ControllerPublishVolume, Create Volume from snapshot, Clone Volume ([#269](https://github.com/kubernetes-csi/csi-test/pull/269), [@taaraora](https://github.com/taaraora))
- Fixes an issue where expand volume RPC requests were failing, since
  volume parameters were missing from the test case. ([#245](https://github.com/kubernetes-csi/csi-test/pull/245), [@utkarshmani1997](https://github.com/utkarshmani1997))


### API Changes

- Rename Cleanup to Resources and unexport cleanup (un-)registration, which is now handled implicitly and automatically. ([#261](https://github.com/kubernetes-csi/csi-test/pull/261), [@timoreimann](https://github.com/timoreimann))


### Other Notable Changes

- Build with Go 1.15 ([#283](https://github.com/kubernetes-csi/csi-test/pull/283), [@pohly](https://github.com/pohly))
- Deployment manifests and examples are added for mock CSI driver under `mock/example` ([#277](https://github.com/kubernetes-csi/csi-test/pull/277), [@Jiawei0227](https://github.com/Jiawei0227))
- Add CSI sanity tests for NodeExpandVolume call. ([#273](https://github.com/kubernetes-csi/csi-test/pull/273), [@gnufied](https://github.com/gnufied))
- Updates CSI Spec to v1.3.0. Updates the mock driver to support volume health monitoring by implementing ControllerGetVolume and updating ListVolumes and NodeVolumeStats functions. ([#268](https://github.com/kubernetes-csi/csi-test/pull/268), [@fengzixu](https://github.com/fengzixu))
