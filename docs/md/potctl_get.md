## potctl get

Get information of existing resources

### Synopsis

Get information of existing resources.

Resources like Agents will require a working Controller in the namespace to display all information.

```
potctl get RESOURCE [flags]
```

### Examples

```
potctl get all
             namespaces
             controllers
             agents
             edge-resources
             application-templates
             applications
             system-applications
             microservices
             system-microservices
             catalog
             registries
             volumes
             routes
             secrets
             configmaps
             services
             volume-mounts
             certificates
```

### Options

```
      --detached   Specify command is to run against detached resources
  -h, --help       help for get
```

### Options inherited from parent commands

```
      --debug              Toggle for displaying verbose output of API clients (HTTP and SSH)
  -n, --namespace string   Namespace to execute respective command within (default "default")
  -v, --verbose            Toggle for displaying verbose output of potctl
```

### SEE ALSO

* [potctl](potctl.md)	 - 


