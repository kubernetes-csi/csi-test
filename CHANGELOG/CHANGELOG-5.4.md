# Release notes for v5.4.0

# Changelog since v5.3.0

## Urgent Upgrade Notes

### (No, really, you MUST read this before you upgrade)

- If you are upgrading CSI spec to v1.10.0 and rely on mocking gRPC calls in your tests, you may need to use the csi-test utils package's Protobuf Matcher: `example.EXPECT().ExampleRequest(Protobuf(requestMsg)).Return(responseMsg, nil).AnyTimes()` ([#553](https://github.com/kubernetes-csi/csi-test/pull/553), [@AndrewSirenko](https://github.com/AndrewSirenko))

## Changes by Kind

### Uncategorized

- Update to github.com/container-storage-interface/spec v1.12.0 ([#577](https://github.com/kubernetes-csi/csi-test/pull/577), [@sunnylovestiramisu](https://github.com/sunnylovestiramisu))

## Dependencies
- since commit 066ea1dba922380ef517db882621207fc25286f4

### Added
- github.com/GoogleCloudPlatform/opentelemetry-operations-go/detectors/gcp: [v1.24.2](https://github.com/GoogleCloudPlatform/opentelemetry-operations-go/tree/detectors/gcp/v1.24.2)
- github.com/go-logr/stdr: [v1.2.2](https://github.com/go-logr/stdr/tree/v1.2.2)
- github.com/go-task/slim-sprig/v3: [v3.0.0](https://github.com/go-task/slim-sprig/tree/v3.0.0)
- github.com/planetscale/vtprotobuf: [0393e58](https://github.com/planetscale/vtprotobuf/tree/0393e58)
- go.opentelemetry.io/contrib/detectors/gcp: v1.31.0
- go.opentelemetry.io/otel/metric: v1.31.0
- go.opentelemetry.io/otel/sdk/metric: v1.31.0
- go.opentelemetry.io/otel/sdk: v1.31.0
- go.opentelemetry.io/otel/trace: v1.31.0
- go.opentelemetry.io/otel: v1.31.0
- golang.org/x/telemetry: bda5523

### Changed
- cel.dev/expr: v0.15.0 → v0.16.2
- cloud.google.com/go/compute/metadata: v0.3.0 → v0.5.2
- github.com/chzyer/readline: [2972be2 → v1.5.1](https://github.com/chzyer/readline/compare/2972be2...v1.5.1)
- github.com/cncf/xds/go: [555b57e → b4127c9](https://github.com/cncf/xds/compare/555b57e...b4127c9)
- github.com/container-storage-interface/spec: [v1.10.0 → v1.12.0](https://github.com/container-storage-interface/spec/compare/v1.10.0...v1.12.0)
- github.com/envoyproxy/go-control-plane: [v0.12.0 → v0.13.1](https://github.com/envoyproxy/go-control-plane/compare/v0.12.0...v0.13.1)
- github.com/envoyproxy/protoc-gen-validate: [v1.0.4 → v1.1.0](https://github.com/envoyproxy/protoc-gen-validate/compare/v1.0.4...v1.1.0)
- github.com/go-logr/logr: [v1.4.1 → v1.4.2](https://github.com/go-logr/logr/compare/v1.4.1...v1.4.2)
- github.com/golang/glog: [v1.2.1 → v1.2.2](https://github.com/golang/glog/compare/v1.2.1...v1.2.2)
- github.com/google/pprof: [94a9f03 → 40e02aa](https://github.com/google/pprof/compare/94a9f03...40e02aa)
- github.com/ianlancetaylor/demangle: [28f6c0f → bd984b5](https://github.com/ianlancetaylor/demangle/compare/28f6c0f...bd984b5)
- github.com/onsi/ginkgo/v2: [v2.13.1 → v2.22.0](https://github.com/onsi/ginkgo/compare/v2.13.1...v2.22.0)
- github.com/onsi/gomega: [v1.30.0 → v1.36.1](https://github.com/onsi/gomega/compare/v1.30.0...v1.36.1)
- github.com/stretchr/testify: [v1.6.1 → v1.8.4](https://github.com/stretchr/testify/compare/v1.6.1...v1.8.4)
- golang.org/x/crypto: v0.23.0 → v0.36.0
- golang.org/x/mod: v0.13.0 → v0.22.0
- golang.org/x/net: v0.25.0 → v0.38.0
- golang.org/x/oauth2: v0.20.0 → v0.23.0
- golang.org/x/sync: v0.7.0 → v0.12.0
- golang.org/x/sys: v0.20.0 → v0.31.0
- golang.org/x/term: v0.20.0 → v0.30.0
- golang.org/x/text: v0.15.0 → v0.23.0
- golang.org/x/tools: v0.14.0 → v0.28.0
- google.golang.org/genproto/googleapis/api: 5315273 → 796eee8
- google.golang.org/genproto/googleapis/rpc: 5315273 → 9240e9c
- google.golang.org/grpc: v1.65.0 → v1.69.2
- google.golang.org/protobuf: v1.34.1 → v1.36.0

### Removed
- github.com/chzyer/logex: [v1.1.10](https://github.com/chzyer/logex/tree/v1.1.10)
- github.com/chzyer/test: [a1ea475](https://github.com/chzyer/test/tree/a1ea475)
- github.com/go-task/slim-sprig: [52ccab3](https://github.com/go-task/slim-sprig/tree/52ccab3)
- github.com/stretchr/objx: [v0.1.0](https://github.com/stretchr/objx/tree/v0.1.0)
