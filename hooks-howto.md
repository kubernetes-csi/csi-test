# JavaScript Hooks in the CSI Mock Driver
The CSI mock driver can call JavaScript hooks during the execution of the CSI
method calls. These hooks can be configured in an external YAML file that can be
passed to the mock driver sidecar container as a `ConfigMap` volume.

## Example
Here's an usage example of the feature. First, prepare a YAML file defining
some JavaScript hooks. Let's call it `hooks.yaml`:

```javascript
globals: |
  count = 0;
  console.log("Globals loaded, count set to " + count);
createVolumeStart: |
  count = count + 1;
  console.log("CreateVolumeStart: The value of count is " + count);
  if (count > 2) { OK; } else { DEADLINEEXCEEDED; };
createVolumeEnd: |
  console.log("CreateVolumeEnd");
```

Create a `ConfigMap` from the YAML:

```sh
kubectl create configmap hooks-config --from-file hooks.yaml
```
Configure the CSI mock driver to use it -- these are the relevant parts of the
sidecar definition YAMLs:

```yaml
...
containers:
  ...
  - name: csi-mock-plugin
    ...
    args:
      - --hooks-file=/etc/hooks/hooks.yaml
    ...
    volumeMounts:
      ...
      - name: hooks-volume
        mountPath: /etc/hooks
...
volumes:
  - name: hooks-volume
    configMap:
      name: hooks-config
...
```

Now the driver finds the hooks definitions file and creates a JavaScript VM.
Just one instance of the VM exists through the life of the driver. It executes
the pieces of code as per the YAML file. One special case is the `globals` hook
that is executed right after loading the file and it's intended for variable
initialization, common function definitions, etc.

Each code snippet is executed at given place in the code: in this example
`createVolumeStart` at the beginning of the `CreateVolume` method in the driver.
The result of the hook snipped is evaluated: it would be the value of the last
expression executed. If the result is an integer value it gets interpreted as
gRPC code and in case it's non-zero (not `OK`), the CSI driver method would
return this as its response error. The complete list of the gRPC constants is
defined in the [hooks-const.go](./mock/service/hooks-const.go) file.

In the example: The variable `count` was initialized to zero in the `globals`
hook. Each time the `CreateVolume` method is called the hook is executed so
`count` is increased by one (it's always the same variable) and for the first
two calls the last expression of the code snipped is `DEADLINEEXCEEDED` which
causes the `CreateVolume` method to fail with the given code. Subsequent calls
would return `0` so things would work normally. Also -- the `console.out()`
allows for logging to stdout so the messages from the hooks appear among the
output from the driver itself.
