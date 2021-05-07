# Release notes for v4.2.0

*Note:* The mock driver is deprecated and will be removed in the next
release. It gets replaced by the
[csi-driver-host-path](https://github.com/kubernetes-csi/csi-driver-host-path/). Error
injection via JavaScript gets replaced with proxying CSI requests and
[Go callbacks inside E2E
tests](https://github.com/kubernetes/kubernetes/blob/5ad79eae2dcbf33df3b35c48ec993d30fbda46dd/test/e2e/storage/csi_mock_volume.go#L110).

# Changelog since v4.1.0

## Changes by Kind

### Enhancements
 - Removal of the CSI nodepublish path by the kubelet is deprecated. This must be done by the CSI plugin according to the CSI spec. ([#338](https://github.com/kubernetes-csi/csi-test/pull/338), [@dobsonj](https://github.com/dobsonj))
 - Better error messages in csi-sanity when encountering unexpected gRPC status codes. ([#334](https://github.com/kubernetes-csi/csi-test/pull/334), [@avorima](https://github.com/avorima))

### Uncategorized
 - Updated runtime (Go 1.16) and dependencies ([#330](https://github.com/kubernetes-csi/csi-test/pull/330), [@pohly](https://github.com/pohly))

## Dependencies

### Added
_Nothing has changed._

### Changed
_Nothing has changed._

### Removed
_Nothing has changed._
