# Release notes for v5.3.0

# Changelog since v5.2.0

## Changes by Kind

### Feature

- Add csi-sanity tests for ControllerModifyVolume RPC ([#538](https://github.com/kubernetes-csi/csi-test/pull/538), [@AndrewSirenko](https://github.com/AndrewSirenko))
- Added mock interfaces for the `SnapshotMetadata` service introduced in CSI v1.10.0.
  Added a `csi.SnapshotMetadataServer` interface to `CSIDriver` and `CSIDriverController`. ([#545](https://github.com/kubernetes-csi/csi-test/pull/545), [@carlbraganza](https://github.com/carlbraganza))

### Bug or Regression

- Added check to ensure ContentSource is set when creating a volume from a snapshot. ([#497](https://github.com/kubernetes-csi/csi-test/pull/497), [@palvarez89](https://github.com/palvarez89))

### Other (Cleanup or Flake)

- Use git describe --tags for version ([#526](https://github.com/kubernetes-csi/csi-test/pull/526), [@huww98](https://github.com/huww98))
- Wait for snapshot ReadyToUse before using ([#525](https://github.com/kubernetes-csi/csi-test/pull/525), [@huww98](https://github.com/huww98))
- Replace offset args with new Ginkgo Helper([#523](https://github.com/kubernetes-csi/csi-test/pull/523), [@huww98](https://github.com/huww98))

## Dependencies

### Added
_Nothing has changed._

### Changed
- cloud.google.com/go/compute: v1.21.0 → v1.25.1
- github.com/cncf/xds/go: [e9ce688 → 8a4994d](https://github.com/cncf/xds/compare/e9ce688...8a4994d)
- github.com/container-storage-interface/spec: [v1.9.0 → v1.10.0](https://github.com/container-storage-interface/spec/compare/v1.9.0...v1.10.0)
- github.com/envoyproxy/go-control-plane: [v0.11.1 → v0.12.0](https://github.com/envoyproxy/go-control-plane/compare/v0.11.1...v0.12.0)
- github.com/envoyproxy/protoc-gen-validate: [v1.0.2 → v1.0.4](https://github.com/envoyproxy/protoc-gen-validate/compare/v1.0.2...v1.0.4)
- github.com/golang/glog: [v1.1.0 → v1.2.0](https://github.com/golang/glog/compare/v1.1.0...v1.2.0)
- github.com/golang/protobuf: [v1.5.3 → v1.5.4](https://github.com/golang/protobuf/compare/v1.5.3...v1.5.4)
- github.com/go-logr/logr: v1.3.0 → v1.4.1
- github.com/google/uuid: [v1.4.0 → v1.6.0](https://github.com/google/uuid/compare/v1.4.0...v1.6.0)
- golang.org/x/crypto: v0.14.0 → v0.21.0
- golang.org/x/net: v0.17.0 → v0.25.0
- golang.org/x/oauth2: v0.10.0 → v0.18.0
- golang.org/x/sync: v0.4.0 → v0.6.0
- golang.org/x/sys: v0.14.0 → v0.20.0
- golang.org/x/term: v0.13.0 → v0.18.0
- golang.org/x/text: v0.13.0 → v0.15.0
- google.golang.org/appengine: v1.6.7 → v1.6.8
- google.golang.org/genproto/googleapis/api: 782d3b1 → 5315273
- google.golang.org/genproto/googleapis/rpc: f966b18 → 5315273
- google.golang.org/grpc: v1.58.1 → v1.65.0
- google.golang.org/protobuf: v1.31.0 → v1.34.1
- k8s.io/klog/v2: v2.110.1 → v2.130.1

### Removed
- github.com/cncf/udpa/go: [c52dc94](https://github.com/cncf/udpa/tree/c52dc94)
- google.golang.org/genproto: 23370e0
