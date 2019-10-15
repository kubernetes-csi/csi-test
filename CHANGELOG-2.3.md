# Changelog since v2.2.0

## New Features

- add '--csi.testvolumeexpandsize` flag ([#219](https://github.com/kubernetes-csi/csi-test/pull/219), [@wnxn](https://github.com/wnxn))
- pass ControllerExpandVolumeSecret to volume expansion tests ([#217](https://github.com/kubernetes-csi/csi-test/pull/217), [@suneeth51](https://github.com/suneeth51))


## Bug Fixes

- fix nil pointer error when Config.IDGen is not set ([#220](https://github.com/kubernetes-csi/csi-test/pull/220), [@pohly](https://github.com/pohly))
- fix cloning of first volume in mock CSI driver([#227](https://github.com/kubernetes-csi/csi-test/pull/227), [@avalluri](https://github.com/avalluri))


## Other Notable Changes

- add Go module support and thus can be built also outside of GOPATH ([#224](https://github.com/kubernetes-csi/csi-test/pull/224), [@pohly](https://github.com/pohly))
- disable unreliable ListVolume pagination test ([#226](https://github.com/kubernetes-csi/csi-test/pull/226), [@avalluri](https://github.com/avalluri))
