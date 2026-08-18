# Release notes for v5.5.0

# Changelog since v5.4.0

## Changes by Kind

### Uncategorized

- Add CSI sanity testing of ControllerServer.GetSnapshot. ([#573](https://github.com/kubernetes-csi/csi-test/pull/573), [@nixpanic](https://github.com/nixpanic))
- Adds tests for `SINGLE_NODE_SINGLE_WRITER` volume access mode behavior. Drivers advertising the `SINGLE_NODE_MULTI_WRITER` node service capability must reject concurrent publish attempts for `SINGLE_NODE_SINGLE_WRITER` volumes with `FailedPrecondition`. ([#602](https://github.com/kubernetes-csi/csi-test/pull/602), [@chrishenzie](https://github.com/chrishenzie))

## Dependencies
- since commit b0702afc614aba2d487d090c6dcd583657f9d05a

### Added
- github.com/Masterminds/semver/v3: [v3.4.0](https://github.com/Masterminds/semver/tree/v3.4.0)
- github.com/envoyproxy/go-control-plane/envoy: [v1.37.0](https://github.com/envoyproxy/go-control-plane/tree/envoy/v1.37.0)
- github.com/envoyproxy/go-control-plane/ratelimit: [v0.1.0](https://github.com/envoyproxy/go-control-plane/tree/ratelimit/v0.1.0)
- github.com/gkampitakis/ciinfo: [v0.3.2](https://github.com/gkampitakis/ciinfo/tree/v0.3.2)
- github.com/gkampitakis/go-diff: [v1.3.2](https://github.com/gkampitakis/go-diff/tree/v1.3.2)
- github.com/gkampitakis/go-snaps: [v0.5.15](https://github.com/gkampitakis/go-snaps/tree/v0.5.15)
- github.com/go-jose/go-jose/v4: [v4.1.4](https://github.com/go-jose/go-jose/tree/v4.1.4)
- github.com/goccy/go-yaml: [v1.18.0](https://github.com/goccy/go-yaml/tree/v1.18.0)
- github.com/joshdk/go-junit: [v1.0.0](https://github.com/joshdk/go-junit/tree/v1.0.0)
- github.com/kr/pretty: [v0.3.1](https://github.com/kr/pretty/tree/v0.3.1)
- github.com/kr/text: [v0.2.0](https://github.com/kr/text/tree/v0.2.0)
- github.com/maruel/natural: [v1.1.1](https://github.com/maruel/natural/tree/v1.1.1)
- github.com/mfridman/tparse: [v0.18.0](https://github.com/mfridman/tparse/tree/v0.18.0)
- github.com/rogpeppe/go-internal: [v1.13.1](https://github.com/rogpeppe/go-internal/tree/v1.13.1)
- github.com/spiffe/go-spiffe/v2: [v2.6.0](https://github.com/spiffe/go-spiffe/tree/v2.6.0)
- github.com/tidwall/gjson: [v1.18.0](https://github.com/tidwall/gjson/tree/v1.18.0)
- github.com/tidwall/match: [v1.1.1](https://github.com/tidwall/match/tree/v1.1.1)
- github.com/tidwall/pretty: [v1.2.1](https://github.com/tidwall/pretty/tree/v1.2.1)
- github.com/tidwall/sjson: [v1.2.5](https://github.com/tidwall/sjson/tree/v1.2.5)
- go.opentelemetry.io/auto/sdk: v1.2.1
- go.uber.org/mock: v0.5.2
- go.yaml.in/yaml/v3: v3.0.4
- gonum.org/v1/gonum: v0.17.0

### Changed
- cel.dev/expr: v0.16.2 → v0.25.1
- cloud.google.com/go/compute/metadata: v0.5.2 → v0.9.0
- github.com/GoogleCloudPlatform/opentelemetry-operations-go/detectors/gcp: [v1.24.2 → v1.32.0](https://github.com/GoogleCloudPlatform/opentelemetry-operations-go/compare/detectors/gcp/v1.24.2...detectors/gcp/v1.32.0)
- github.com/cncf/xds/go: [b4127c9 → dba9d58](https://github.com/cncf/xds/compare/b4127c9...dba9d58)
- github.com/envoyproxy/go-control-plane: [v0.13.1 → v0.14.0](https://github.com/envoyproxy/go-control-plane/compare/v0.13.1...v0.14.0)
- github.com/envoyproxy/protoc-gen-validate: [v1.1.0 → v1.3.3](https://github.com/envoyproxy/protoc-gen-validate/compare/v1.1.0...v1.3.3)
- github.com/go-logr/logr: [v1.4.2 → v1.4.3](https://github.com/go-logr/logr/compare/v1.4.2...v1.4.3)
- github.com/golang/glog: [v1.2.2 → v1.2.5](https://github.com/golang/glog/compare/v1.2.2...v1.2.5)
- github.com/google/go-cmp: [v0.6.0 → v0.7.0](https://github.com/google/go-cmp/compare/v0.6.0...v0.7.0)
- github.com/google/pprof: [40e02aa → 545e8a4](https://github.com/google/pprof/compare/40e02aa...545e8a4)
- github.com/ianlancetaylor/demangle: [bd984b5 → f615e6b](https://github.com/ianlancetaylor/demangle/compare/bd984b5...f615e6b)
- github.com/onsi/ginkgo/v2: [v2.22.0 → v2.32.0](https://github.com/onsi/ginkgo/compare/v2.22.0...v2.32.0)
- github.com/onsi/gomega: [v1.36.1 → v1.42.1](https://github.com/onsi/gomega/compare/v1.36.1...v1.42.1)
- github.com/stretchr/testify: [v1.8.4 → v1.9.0](https://github.com/stretchr/testify/compare/v1.8.4...v1.9.0)
- go.opentelemetry.io/contrib/detectors/gcp: v1.31.0 → v1.43.0
- go.opentelemetry.io/otel: v1.31.0 → v1.43.0
- go.opentelemetry.io/otel/metric: v1.31.0 → v1.43.0
- go.opentelemetry.io/otel/sdk: v1.31.0 → v1.43.0
- go.opentelemetry.io/otel/sdk/metric: v1.31.0 → v1.43.0
- go.opentelemetry.io/otel/trace: v1.31.0 → v1.43.0
- golang.org/x/crypto: v0.36.0 → v0.53.0
- golang.org/x/mod: v0.22.0 → v0.36.0
- golang.org/x/net: v0.38.0 → v0.56.0
- golang.org/x/oauth2: v0.23.0 → v0.36.0
- golang.org/x/sync: v0.12.0 → v0.21.0
- golang.org/x/sys: v0.31.0 → v0.46.0
- golang.org/x/telemetry: bda5523 → 42602be
- golang.org/x/term: v0.30.0 → v0.44.0
- golang.org/x/text: v0.23.0 → v0.38.0
- golang.org/x/tools: v0.28.0 → v0.45.0
- google.golang.org/genproto/googleapis/api: 796eee8 → afd174a
- google.golang.org/genproto/googleapis/rpc: 9240e9c → afd174a
- google.golang.org/grpc: v1.69.2 → v1.82.0
- google.golang.org/protobuf: v1.36.0 → v1.36.11
- k8s.io/klog/v2: v2.130.1 → v2.140.0

### Removed
- github.com/census-instrumentation/opencensus-proto: [v0.4.1](https://github.com/census-instrumentation/opencensus-proto/tree/v0.4.1)
