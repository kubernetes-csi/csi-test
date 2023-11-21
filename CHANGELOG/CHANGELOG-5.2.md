# Release notes for v5.2.0

# Changelog since v5.1.0

## Changes by Kind

### API Change

- Bumps the CSI spec dependency to v1.9.0 and regenerates the controller mock to account for the new ControllerModifyVolume RPC. ([#492](https://github.com/kubernetes-csi/csi-test/pull/492), [@ConnorJC3](https://github.com/ConnorJC3))

## Dependencies

### Added
_Nothing has changed._

### Changed
- github.com/container-storage-interface/spec: [v1.8.0 → v1.9.0](https://github.com/container-storage-interface/spec/compare/v1.8.0...v1.9.0)
- github.com/go-logr/logr: [v1.2.4 → v1.3.0](https://github.com/go-logr/logr/compare/v1.2.4...v1.3.0)
- github.com/google/go-cmp: [v0.5.9 → v0.6.0](https://github.com/google/go-cmp/compare/v0.5.9...v0.6.0)
- github.com/google/uuid: [v1.3.1 → v1.4.0](https://github.com/google/uuid/compare/v1.3.1...v1.4.0)
- github.com/onsi/ginkgo/v2: [v2.12.0 → v2.13.1](https://github.com/onsi/ginkgo/v2/compare/v2.12.0...v2.13.1)
- github.com/onsi/gomega: [v1.27.10 → v1.30.0](https://github.com/onsi/gomega/compare/v1.27.10...v1.30.0)
- golang.org/x/crypto: v0.12.0 → v0.14.0
- golang.org/x/mod: v0.12.0 → v0.13.0
- golang.org/x/net: v0.14.0 → v0.17.0
- golang.org/x/sync: v0.3.0 → v0.4.0
- golang.org/x/sys: v0.11.0 → v0.14.0
- golang.org/x/term: v0.11.0 → v0.13.0
- golang.org/x/text: v0.12.0 → v0.13.0
- golang.org/x/tools: v0.12.0 → v0.14.0
- google.golang.org/genproto/googleapis/rpc: 782d3b1 → f966b18
- google.golang.org/genproto: 782d3b1 → 23370e0
- k8s.io/klog/v2: v2.100.1 → v2.110.1

### Removed
_Nothing has changed._
