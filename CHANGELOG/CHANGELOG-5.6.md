# Release notes for v5.6.0

# Changelog since v5.5.0

## Changes by Kind

### Feature

- Add support for Volume health. ([#614](https://github.com/kubernetes-csi/csi-test/pull/614), [@gnufied](https://github.com/gnufied))
- Add support for Volume Snapshot Topology. ([#617](https://github.com/kubernetes-csi/csi-test/pull/617), [@mdzraf](https://github.com/mdzraf))

### Uncategorized

- The csi-sanity invalid ListVolumes starting token can now be customized with the `--csi.testinvalidlistvolumesstartingtoken` flag. ([#620](https://github.com/kubernetes-csi/csi-test/pull/620), [@dlanov](https://github.com/dlanov))

## Dependencies
- since commit 2f1e2f2aa8b77199c63e926628fc4c377bfa525d

### Added
_Nothing has changed._

### Changed
- github.com/container-storage-interface/spec: [v1.12.0 → v1.13.0](https://github.com/container-storage-interface/spec/compare/v1.12.0...v1.13.0)
- github.com/onsi/ginkgo/v2: [v2.32.0 → v2.32.1](https://github.com/onsi/ginkgo/compare/v2.32.0...v2.32.1)
- go.uber.org/mock: v0.5.2 → v0.6.0
- google.golang.org/grpc: [v1.82.0 → v1.83.1](https://github.com/grpc/grpc-go/compare/v1.82.0...v1.83.1)
- google.golang.org/protobuf: v1.36.11 → v1.36.12

### Removed
- github.com/golang/mock: [v1.6.0](https://github.com/golang/mock/tree/v1.6.0)
