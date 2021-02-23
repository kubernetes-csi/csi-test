# Mock CSI Driver
Extremely simple mock driver used to test `csi-sanity` based on `rexray/gocsi/mock`.
It can be used for testing of Container Orchestrators that implement client side
of CSI interface.

```
Usage of mock:
  -disable-attach
        Disables RPC_PUBLISH_UNPUBLISH_VOLUME capability.
  -name string
        CSI driver name. (default "io.kubernetes.storage.mock")
  -v=5
        Enable gRPC call logging.
```

Limitation about this mock CSI Driver are:
- It only supports single node.
- It requires all the components to run on the same node, i.e. the pod that uses pv created by this CSI driver should be on the same node with the driver.
- It does not persist data across restarts. All the states are in memory.

It prints all received CSI messages to stdout encoded as json, so a test can check that
CO sent the right CSI message.

Example of such output:

```
gRPCCall: {"Method":"/csi.v0.Controller/ControllerGetCapabilities","Request":{},"Response":{"capabilities":[{"Type":{"Rpc":{"type":1}}},{"Type":{"Rpc":{"type":3}}},{"Type":{"Rpc":{"type":4}}},{"Type":{"Rpc":{"type":6}}},{"Type":{"Rpc":{"type":5}}},{"Type":{"Rpc":{"type":2}}}]},"Error":""}
gRPCCall: {"Method":"/csi.v0.Controller/ControllerPublishVolume","Request":{"volume_id":"12","node_id":"some-fake-node-id","volume_capability":{"AccessType":{"Mount":{}},"access_mode":{"mode":1}}},"Response":null,"Error":"rpc error: code = NotFound desc = Not matching Node ID some-fake-node-id to Mock Node ID io.kubernetes.storage.mock"}
```

## Mock CSI Driver Example

This example requires kubernetes 1.18+

The example folder contains an example manifest of deploying CSIDriver including `csi-driver-node-registrar`, `csi-provisioner`,
`csi-resizer` and `csi-snapshotter` onto a cluster. For testing purpose, `csi-attacher` is not included. Thus, 
`"--disable-attach=true"` is set on the csi mock driver.

### Usage

To install the CSI mock driver in your cluster.
```
$ kubectl apply -f example/deploy/csi-mock-driver-rbac.yaml
$ kubectl apply -f example/deploy/csi-mock-driver-deployment.yaml
```

This will deploy a `CSIDriver` called `io.kubernetes.storage.mock`. There will be a deployment `csi-mockplugin` being 
installed in the default namespace. Correspondingly, you can find a pod with prefix `csi-mockplugin`. In the meantime,
 a `StorageClass` called `test-csi-mock` will be generated along with a `VolumeSnapshotClass` called `csi-mock-snapclass`.

K8s distribution 1.17+ is supposed to install [snapshot-controller](https://github.com/kubernetes-csi/external-snapshotter#design)
by default. However, if it does not, you need to install it manually follow the [external-snapshotter doc](https://github.com/kubernetes-csi/external-snapshotter#usage) to test the snapshot feature.


### Experiment

Recommended steps to follow:

1. Deploy `example/pvc-test.yaml` to test the provisioner. Edit the pvc spec to test csi volume resizer.
2. Deploy `example/pod-test.yaml` to test volume mounting.
3. Deploy `example/snapshot-test.yaml` to test snapshot of the pvc deployed in step 1.
4. Deploy `example/snapshot-restore-test.yaml` to test restore volume from a snapshot generated from step 3s.
