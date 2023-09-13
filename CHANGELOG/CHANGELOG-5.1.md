# Release notes for v5.1.0

# Changelog since v5.0.0

## Changes by Kind

### Feature

- Addition of initial tests for the CSI GroupController server and its VolumeGroupSnapshot service. ([#467](https://github.com/kubernetes-csi/csi-test/pull/467), [@nixpanic](https://github.com/nixpanic))
- The external-snapshotter is going to support [Volume Group Snapshots](kubernetes/enhancements#1551), so drivers will return the `GROUP_CONTROLLER_SERVICE` capability as well. ([#450](https://github.com/kubernetes-csi/csi-test/pull/450), [@nixpanic](https://github.com/nixpanic))

### Bug or Regression

- Set the grpc authority to "localhost" to mimic Kubernetes behaviour ([#468](https://github.com/kubernetes-csi/csi-test/pull/468), [@LucaDev](https://github.com/LucaDev))

### Uncategorized

- Add identity service sanity test for controller service ([#465](https://github.com/kubernetes-csi/csi-test/pull/465), [@carlory](https://github.com/carlory))

## Dependencies

### Added
- cloud.google.com/go/compute/metadata: v0.2.3
- cloud.google.com/go/compute: v1.21.0
- google.golang.org/genproto/googleapis/api: 782d3b1
- google.golang.org/genproto/googleapis/rpc: 782d3b1

### Changed
- github.com/census-instrumentation/opencensus-proto: [v0.2.1 → v0.4.1](https://github.com/census-instrumentation/opencensus-proto/compare/v0.2.1...v0.4.1)
- github.com/cespare/xxhash/v2: [v2.1.1 → v2.2.0](https://github.com/cespare/xxhash/v2/compare/v2.1.1...v2.2.0)
- github.com/cncf/udpa/go: [04548b0 → c52dc94](https://github.com/cncf/udpa/go/compare/04548b0...c52dc94)
- github.com/cncf/xds/go: [cb28da3 → e9ce688](https://github.com/cncf/xds/go/compare/cb28da3...e9ce688)
- github.com/container-storage-interface/spec: [v1.6.0 → v1.8.0](https://github.com/container-storage-interface/spec/compare/v1.6.0...v1.8.0)
- github.com/envoyproxy/go-control-plane: [49ff273 → v0.11.1](https://github.com/envoyproxy/go-control-plane/compare/49ff273...v0.11.1)
- github.com/envoyproxy/protoc-gen-validate: [v0.1.0 → v1.0.2](https://github.com/envoyproxy/protoc-gen-validate/compare/v0.1.0...v1.0.2)
- github.com/go-logr/logr: [v1.2.0 → v1.2.4](https://github.com/go-logr/logr/compare/v1.2.0...v1.2.4)
- github.com/go-task/slim-sprig: [348f09d → 52ccab3](https://github.com/go-task/slim-sprig/compare/348f09d...52ccab3)
- github.com/golang/glog: [23def4e → v1.1.0](https://github.com/golang/glog/compare/23def4e...v1.1.0)
- github.com/golang/protobuf: [v1.5.2 → v1.5.3](https://github.com/golang/protobuf/compare/v1.5.2...v1.5.3)
- github.com/google/go-cmp: [v0.5.8 → v0.5.9](https://github.com/google/go-cmp/compare/v0.5.8...v0.5.9)
- github.com/google/uuid: [v1.3.0 → v1.3.1](https://github.com/google/uuid/compare/v1.3.0...v1.3.1)
- github.com/onsi/ginkgo/v2: [v2.1.4 → v2.12.0](https://github.com/onsi/ginkgo/v2/compare/v2.1.4...v2.12.0)
- github.com/onsi/gomega: [v1.20.0 → v1.27.10](https://github.com/onsi/gomega/compare/v1.20.0...v1.27.10)
- github.com/stretchr/testify: [v1.7.0 → v1.6.1](https://github.com/stretchr/testify/compare/v1.7.0...v1.6.1)
- github.com/yuin/goldmark: [v1.4.1 → v1.4.13](https://github.com/yuin/goldmark/compare/v1.4.1...v1.4.13)
- golang.org/x/crypto: 089bfa5 → v0.12.0
- golang.org/x/mod: 9b9b3d8 → v0.12.0
- golang.org/x/net: 0bcc04d → v0.14.0
- golang.org/x/oauth2: bf48bf1 → v0.10.0
- golang.org/x/sync: 036812b → v0.3.0
- golang.org/x/sys: a90be44 → v0.11.0
- golang.org/x/term: 03fcf44 → v0.11.0
- golang.org/x/text: v0.3.7 → v0.12.0
- golang.org/x/tools: v0.1.10 → v0.12.0
- google.golang.org/appengine: v1.4.0 → v1.6.7
- google.golang.org/genproto: f927205 → 782d3b1
- google.golang.org/grpc: v1.48.0 → v1.58.0
- google.golang.org/protobuf: v1.28.0 → v1.31.0
- k8s.io/klog/v2: v2.70.1 → v2.100.1

### Removed
- cloud.google.com/go: v0.34.0
- github.com/BurntSushi/toml: [v0.3.1](https://github.com/BurntSushi/toml/tree/v0.3.1)
- github.com/antihax/optional: [v1.0.0](https://github.com/antihax/optional/tree/v1.0.0)
- github.com/client9/misspell: [v0.3.4](https://github.com/client9/misspell/tree/v0.3.4)
- github.com/fsnotify/fsnotify: [v1.4.9](https://github.com/fsnotify/fsnotify/tree/v1.4.9)
- github.com/ghodss/yaml: [v1.0.0](https://github.com/ghodss/yaml/tree/v1.0.0)
- github.com/grpc-ecosystem/grpc-gateway: [v1.16.0](https://github.com/grpc-ecosystem/grpc-gateway/tree/v1.16.0)
- github.com/hpcloud/tail: [v1.0.0](https://github.com/hpcloud/tail/tree/v1.0.0)
- github.com/nxadm/tail: [v1.4.8](https://github.com/nxadm/tail/tree/v1.4.8)
- github.com/onsi/ginkgo: [v1.16.4](https://github.com/onsi/ginkgo/tree/v1.16.4)
- github.com/prometheus/client_model: [14fe0d1](https://github.com/prometheus/client_model/tree/14fe0d1)
- github.com/rogpeppe/fastuuid: [v1.2.0](https://github.com/rogpeppe/fastuuid/tree/v1.2.0)
- go.opentelemetry.io/proto/otlp: v0.7.0
- golang.org/x/exp: 509febe
- golang.org/x/lint: d0100b6
- gopkg.in/fsnotify.v1: v1.4.7
- gopkg.in/tomb.v1: dd63297
- honnef.co/go/tools: ea95bdf
