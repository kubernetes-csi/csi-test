# Release notes for v4.4.0

# Changelog since v4.3.0

## Changes by Kind

### Feature

- Users will need to add ListSnapshotsSecret to a YAML file to pass secrets to the CSI Driver in from the csi-sanity -csi.secrets ${secrets} ... call. ([#358](https://github.com/kubernetes-csi/csi-test/pull/358), [@jskazinski](https://github.com/jskazinski))

## Dependencies

### Added
- github.com/antihax/optional: [v1.0.0](https://github.com/antihax/optional/tree/v1.0.0)
- github.com/cespare/xxhash/v2: [v2.1.1](https://github.com/cespare/xxhash/v2/tree/v2.1.1)
- github.com/chzyer/logex: [v1.1.10](https://github.com/chzyer/logex/tree/v1.1.10)
- github.com/chzyer/readline: [2972be2](https://github.com/chzyer/readline/tree/2972be2)
- github.com/chzyer/test: [a1ea475](https://github.com/chzyer/test/tree/a1ea475)
- github.com/cncf/xds/go: [cb28da3](https://github.com/cncf/xds/go/tree/cb28da3)
- github.com/ghodss/yaml: [v1.0.0](https://github.com/ghodss/yaml/tree/v1.0.0)
- github.com/go-task/slim-sprig: [348f09d](https://github.com/go-task/slim-sprig/tree/348f09d)
- github.com/google/pprof: [94a9f03](https://github.com/google/pprof/tree/94a9f03)
- github.com/grpc-ecosystem/grpc-gateway: [v1.16.0](https://github.com/grpc-ecosystem/grpc-gateway/tree/v1.16.0)
- github.com/ianlancetaylor/demangle: [28f6c0f](https://github.com/ianlancetaylor/demangle/tree/28f6c0f)
- github.com/onsi/ginkgo/v2: [v2.1.3](https://github.com/onsi/ginkgo/v2/tree/v2.1.3)
- github.com/rogpeppe/fastuuid: [v1.2.0](https://github.com/rogpeppe/fastuuid/tree/v1.2.0)
- github.com/yuin/goldmark: [v1.3.5](https://github.com/yuin/goldmark/tree/v1.3.5)
- go.opentelemetry.io/proto/otlp: v0.7.0
- golang.org/x/mod: v0.4.2
- gopkg.in/yaml.v3: 9f266ea

### Changed
- cloud.google.com/go: v0.26.0 → v0.34.0
- github.com/cncf/udpa/go: [efcf912 → 04548b0](https://github.com/cncf/udpa/go/compare/efcf912...04548b0)
- github.com/container-storage-interface/spec: [v1.5.0 → v1.6.0](https://github.com/container-storage-interface/spec/compare/v1.5.0...v1.6.0)
- github.com/davecgh/go-spew: [v1.1.0 → v1.1.1](https://github.com/davecgh/go-spew/compare/v1.1.0...v1.1.1)
- github.com/envoyproxy/go-control-plane: [v0.9.7 → 49ff273](https://github.com/envoyproxy/go-control-plane/compare/v0.9.7...49ff273)
- github.com/go-logr/logr: [v0.3.0 → v1.2.0](https://github.com/go-logr/logr/compare/v0.3.0...v1.2.0)
- github.com/golang/mock: [v1.4.4 → v1.6.0](https://github.com/golang/mock/compare/v1.4.4...v1.6.0)
- github.com/golang/protobuf: [v1.4.3 → v1.5.2](https://github.com/golang/protobuf/compare/v1.4.3...v1.5.2)
- github.com/google/go-cmp: [v0.5.0 → v0.5.6](https://github.com/google/go-cmp/compare/v0.5.0...v0.5.6)
- github.com/google/uuid: [v1.1.2 → v1.3.0](https://github.com/google/uuid/compare/v1.1.2...v1.3.0)
- github.com/nxadm/tail: [v1.4.5 → v1.4.8](https://github.com/nxadm/tail/compare/v1.4.5...v1.4.8)
- github.com/onsi/ginkgo: [v1.14.2 → v1.16.5](https://github.com/onsi/ginkgo/compare/v1.14.2...v1.16.5)
- github.com/onsi/gomega: [v1.10.4 → v1.19.0](https://github.com/onsi/gomega/compare/v1.10.4...v1.19.0)
- github.com/stretchr/testify: [v1.5.1 → v1.7.0](https://github.com/stretchr/testify/compare/v1.5.1...v1.7.0)
- golang.org/x/net: ac852fb → 27dd868
- golang.org/x/oauth2: d2e6202 → bf48bf1
- golang.org/x/sync: 1122301 → 036812b
- golang.org/x/sys: d4d67f9 → 1d35b9e
- golang.org/x/term: 7de9c90 → 03fcf44
- golang.org/x/text: v0.3.4 → v0.3.7
- golang.org/x/tools: 2c0ae70 → v0.1.1
- google.golang.org/grpc: v1.34.0 → v1.47.0
- google.golang.org/protobuf: v1.25.0 → v1.27.1
- k8s.io/klog/v2: v2.4.0 → v2.60.1

### Removed
_Nothing has changed._
