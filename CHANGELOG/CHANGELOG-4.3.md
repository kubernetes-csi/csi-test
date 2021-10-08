# Release notes for v4.3.0

# Changelog since v4.2.0

## Changes by Kind

### Failing Test
 - Avoid test failures when a CSI driver advertises the `VOLUME_MOUNT_GROUP` capability from CSI 1.5 ([#349](https://github.com/kubernetes-csi/csi-test/pull/349), [@andyzhangx](https://github.com/andyzhangx))
 - Avoid test failures when a CSI driver advertises the `SINGLE_NODE_MULTI_WRITER` capability from CSI 1.5 ([#342](https://github.com/kubernetes-csi/csi-test/pull/342), [@chrishenzie](https://github.com/chrishenzie))

### Other (Cleanup or Flake)
 - The mock driver gets removed. https://github.com/kubernetes-csi/csi-driver-host-path should be used instead. ([#351](https://github.com/kubernetes-csi/csi-test/pull/351), [@pohly](https://github.com/pohly))

## Dependencies

### Added
_Nothing has changed._

### Changed
- github.com/container-storage-interface/spec: [v1.3.0 → v1.5.0](https://github.com/container-storage-interface/spec/compare/v1.3.0...v1.5.0)
- github.com/davecgh/go-spew: [v1.1.1 → v1.1.0](https://github.com/davecgh/go-spew/compare/v1.1.1...v1.1.0)

### Removed
- github.com/robertkrimen/otto: [ef014fd](https://github.com/robertkrimen/otto/tree/ef014fd)
- github.com/sirupsen/logrus: [v1.7.0](https://github.com/sirupsen/logrus/tree/v1.7.0)
- gopkg.in/sourcemap.v1: v1.0.5
