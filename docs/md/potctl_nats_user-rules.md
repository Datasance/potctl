## potctl nats user-rules

NATS user rule operations

### Synopsis

List, create, update, and delete NATS user rules.

### Examples

```
potctl nats user-rules list
potctl nats nats-user-rule create-yaml -f ./user-rule.yaml
potctl nats user-rule update-yaml default-user-rule -f ./user-rule.yaml
```

### Options

```
  -h, --help   help for user-rules
```

### Options inherited from parent commands

```
      --debug              Toggle for displaying verbose output of API clients (HTTP and SSH)
  -n, --namespace string   Namespace to execute respective command within (default "default")
  -v, --verbose            Toggle for displaying verbose output of potctl
```

### SEE ALSO

* [potctl nats](potctl_nats.md)	 - Manage NATS resources
* [potctl nats user-rules create-yaml](potctl_nats_user-rules_create-yaml.md)	 - Create NATS user rule from YAML file
* [potctl nats user-rules delete](potctl_nats_user-rules_delete.md)	 - Delete NATS user rule
* [potctl nats user-rules list](potctl_nats_user-rules_list.md)	 - List NATS user rules
* [potctl nats user-rules update-yaml](potctl_nats_user-rules_update-yaml.md)	 - Update NATS user rule from YAML file


