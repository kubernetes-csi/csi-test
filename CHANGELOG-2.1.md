# Changelog since v2.0.1

## New Features

- the sanity volume code can be used concurrently and avoids a potential volume leak because of an incorrect target path ([#210](https://github.com/kubernetes-csi/csi-test/pull/210), [@pohly](https://github.com/pohly))
- csi-mock driver now supports relaxing the target path check (--permissive-target-path) and ephemeral inline volumes ([#209](https://github.com/kubernetes-csi/csi-test/pull/209), [@pohly](https://github.com/pohly))


## Bug Fixes

- Remove volume listing tests that assumed tokens were integers and replace them with one that confirms pagination works when adding/deleting volumes between page requests.   ([#205](https://github.com/kubernetes-csi/csi-test/pull/205), [@Akrog](https://github.com/Akrog))


## Other Notable Changes

- Added tests for ControllerExpandVolume ([#203](https://github.com/kubernetes-csi/csi-test/pull/203), [@Ntr0](https://github.com/Ntr0))
