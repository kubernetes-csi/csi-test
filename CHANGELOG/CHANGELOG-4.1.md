# Release notes for v4.1.0

# Changelog since v4.0.2

## Changes by Kind

### Feature
 - Mock driver: optionally proxy connections instead of handling them directly ([#314](https://github.com/kubernetes-csi/csi-test/pull/314), [@pohly](https://github.com/pohly))
 - csi-sanity: clean up sourced volumes prior to deleting sourcing resources ([#297](https://github.com/kubernetes-csi/csi-test/pull/297), [@timoreimann](https://github.com/timoreimann))

### Bug or Regression
 - Add kubelet pod and csi plugin directories to mock driver manifest ([#308](https://github.com/kubernetes-csi/csi-test/pull/308), [@msau42](https://github.com/msau42))
 - Log output from klog/v2 at higher log levels is now available ([#313](https://github.com/kubernetes-csi/csi-test/pull/313), [@pohly](https://github.com/pohly))
 - Mock driver: fix formatting of error message for invalid starting token ([#322](https://github.com/kubernetes-csi/csi-test/pull/322), [@pohly](https://github.com/pohly))
 - Updates the mock driver to create the target directory in `NodePublishVolume`, and remove it in `NodeUnpublishVolume`. ([#303](https://github.com/kubernetes-csi/csi-test/pull/303), [@huffmanca](https://github.com/huffmanca))

### Other (Cleanup or Flake)
 - Updated dependencies ([#315](https://github.com/kubernetes-csi/csi-test/pull/315), [@pohly](https://github.com/pohly))

## Dependencies

### Added
- github.com/nxadm/tail: [v1.4.5](https://github.com/nxadm/tail/tree/v1.4.5)
- golang.org/x/term: 7de9c90
- golang.org/x/xerrors: 5ec99f8
- google.golang.org/protobuf: v1.25.0
- k8s.io/klog/v2: v2.4.0

### Changed
- github.com/cncf/udpa/go: [269d4d4 → efcf912](https://github.com/cncf/udpa/go/compare/269d4d4...efcf912)
- github.com/envoyproxy/go-control-plane: [v0.9.4 → v0.9.7](https://github.com/envoyproxy/go-control-plane/compare/v0.9.4...v0.9.7)
- github.com/fsnotify/fsnotify: [v1.4.7 → v1.4.9](https://github.com/fsnotify/fsnotify/compare/v1.4.7...v1.4.9)
- github.com/go-logr/logr: [v0.1.0 → v0.3.0](https://github.com/go-logr/logr/compare/v0.1.0...v0.3.0)
- github.com/golang/mock: [v1.3.1 → v1.4.4](https://github.com/golang/mock/compare/v1.3.1...v1.4.4)
- github.com/golang/protobuf: [v1.3.3 → v1.4.3](https://github.com/golang/protobuf/compare/v1.3.3...v1.4.3)
- github.com/google/go-cmp: [v0.2.0 → v0.5.0](https://github.com/google/go-cmp/compare/v0.2.0...v0.5.0)
- github.com/google/uuid: [v1.1.1 → v1.1.2](https://github.com/google/uuid/compare/v1.1.1...v1.1.2)
- github.com/onsi/ginkgo: [v1.10.3 → v1.14.2](https://github.com/onsi/ginkgo/compare/v1.10.3...v1.14.2)
- github.com/onsi/gomega: [v1.7.1 → v1.10.4](https://github.com/onsi/gomega/compare/v1.7.1...v1.10.4)
- github.com/robertkrimen/otto: [c382bd3 → ef014fd](https://github.com/robertkrimen/otto/compare/c382bd3...ef014fd)
- github.com/sirupsen/logrus: [v1.4.2 → v1.7.0](https://github.com/sirupsen/logrus/compare/v1.4.2...v1.7.0)
- github.com/stretchr/objx: [v0.1.1 → v0.1.0](https://github.com/stretchr/objx/compare/v0.1.1...v0.1.0)
- github.com/stretchr/testify: [v1.2.2 → v1.5.1](https://github.com/stretchr/testify/compare/v1.2.2...v1.5.1)
- golang.org/x/crypto: c2843e0 → 75b2880
- golang.org/x/net: 2180aed → ac852fb
- golang.org/x/sys: 4c7a9d0 → d4d67f9
- golang.org/x/text: v0.3.2 → v0.3.4
- google.golang.org/genproto: 6bbd007 → f927205
- google.golang.org/grpc: v1.29.1 → v1.34.0
- gopkg.in/yaml.v2: v2.2.5 → v2.4.0

### Removed
- github.com/konsorten/go-windows-terminal-sequences: [v1.0.2](https://github.com/konsorten/go-windows-terminal-sequences/tree/v1.0.2)
- k8s.io/klog: v1.0.0
