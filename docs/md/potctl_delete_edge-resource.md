## potctl delete edge-resource

Delete an Edge Resource

### Synopsis

Delete an Edge Resource.

Only the specified version will be deleted.
Agents that this Edge Resource are attached to will be notified of the deletion.

```
potctl delete edge-resource NAME VERSION [flags]
```

### Examples

```
potctl delete edge-resource NAME VERSION
```

### Options

```
  -h, --help   help for edge-resource
```

### Options inherited from parent commands

```
      --debug              Toggle for displaying verbose output of API clients (HTTP and SSH)
  -n, --namespace string   Namespace to execute respective command within (default "default")
  -v, --verbose            Toggle for displaying verbose output of potctl
```

### SEE ALSO

* [potctl delete](potctl_delete.md)	 - Delete an existing ioFog resource

###### Auto generated by spf13/cobra on 18-Dec-2024
