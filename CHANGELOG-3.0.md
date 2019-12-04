# Changelog since v2.3.0

## New Features

- `go get github.com/kubernetes-csi/csi-test/cmd/csi-sanity` now works ([#234](https://github.com/kubernetes-csi/csi-test/pull/234), [@pohly](https://github.com/pohly))
- Add '--csi.testsnapshotparameters' flag ([#236](https://github.com/kubernetes-csi/csi-test/pull/236), [@wnxn](https://github.com/wnxn))


## API Changes

- gRPC keepalive is no longer enabled by default and must be added to dial options in `TestConfig` if needed ([#239](https://github.com/kubernetes-csi/csi-test/pull/239), [@pohly](https://github.com/pohly))
- revised sanity API: more consistent naming, NewTestConfig must be used ([#233](https://github.com/kubernetes-csi/csi-test/pull/233), [@pohly](https://github.com/pohly))
- Update package path to v3, i.e. imports must be changed to `github.com/kubernetes-csi/csi-test/v3/pkg/sanity`. Vendoring with dep depends on https://github.com/golang/dep/pull/1963 or the workaround described in [v3/README.md](./v3/README.md). ([#232](https://github.com/kubernetes-csi/csi-test/pull/232), [@Ntr0](https://github.com/Ntr0))


## Other Notable Changes

- Update CSI Spec to v1.2.0 1.13.3 ([#230](https://github.com/kubernetes-csi/csi-test/pull/230), [@davidz627](https://github.com/davidz627))
- csi-test is now built and tested with Go 1.13.3 and uses the latest releases of all dependencies. It should still build with older releases. ([#240](https://github.com/kubernetes-csi/csi-test/pull/240), [@pohly](https://github.com/pohly))
