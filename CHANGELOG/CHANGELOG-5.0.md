# Release notes for v5.0.0

# Changelog since v4.4.0

## Changes by Kind

### API Changes
 - csi-sanity is now built with Ginkgo v2. All camel case flags
   (e.g. -ginkgo.randomizeAllSpecs) are replaced with kebab case flags
   (e.g. -ginkgo.randomize-all-specs). The camel case versions continue to work but
   emit a deprecation warning. The "junitfile" configuration file option was
   removed because ginkgo now handles JUnit report creation
   itself. "-csi.junitfile" continues to work, but is only an alias for
   "-ginkgo.junit-report" and should be replaced by that.
   ([#366](https://github.com/kubernetes-csi/csi-test/pull/366),
   [@pohly](https://github.com/pohly))

## Dependencies

### Added
_Nothing has changed._

### Changed
- github.com/google/go-cmp: [v0.5.6 → v0.5.8](https://github.com/google/go-cmp/compare/v0.5.6...v0.5.8)
- github.com/onsi/ginkgo/v2: [v2.1.3 → v2.1.4](https://github.com/onsi/ginkgo/v2/compare/v2.1.3...v2.1.4)
- github.com/onsi/ginkgo: [v1.16.5 → v1.16.4](https://github.com/onsi/ginkgo/compare/v1.16.5...v1.16.4)
- github.com/onsi/gomega: [v1.19.0 → v1.20.0](https://github.com/onsi/gomega/compare/v1.19.0...v1.20.0)
- github.com/yuin/goldmark: [v1.3.5 → v1.4.1](https://github.com/yuin/goldmark/compare/v1.3.5...v1.4.1)
- golang.org/x/crypto: 75b2880 → 089bfa5
- golang.org/x/mod: v0.4.2 → 9b9b3d8
- golang.org/x/net: 27dd868 → 0bcc04d
- golang.org/x/sys: 1d35b9e → a90be44
- golang.org/x/tools: v0.1.1 → v0.1.10
- google.golang.org/grpc: v1.47.0 → v1.48.0
- google.golang.org/protobuf: v1.27.1 → v1.28.0
- gopkg.in/yaml.v3: 9f266ea → v3.0.1
- k8s.io/klog/v2: v2.60.1 → v2.70.1

### Removed
_Nothing has changed._
