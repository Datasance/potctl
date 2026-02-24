## potctl nats account-rules

NATS account rule operations

### Synopsis

List, create, update, and delete NATS account rules.

### Examples

```
potctl nats account-rules list
potctl nats nats-account-rule create-yaml -f ./account-rule.yaml
potctl nats account-rule update-yaml default-account-rule -f ./account-rule.yaml
```

### Options

```
  -h, --help   help for account-rules
```

### Options inherited from parent commands

```
      --debug              Toggle for displaying verbose output of API clients (HTTP and SSH)
  -n, --namespace string   Namespace to execute respective command within (default "default")
  -v, --verbose            Toggle for displaying verbose output of potctl
```

### SEE ALSO

* [potctl nats](potctl_nats.md)	 - Manage NATS resources
* [potctl nats account-rules create-yaml](potctl_nats_account-rules_create-yaml.md)	 - Create NATS account rule from YAML file
* [potctl nats account-rules delete](potctl_nats_account-rules_delete.md)	 - Delete NATS account rule
* [potctl nats account-rules list](potctl_nats_account-rules_list.md)	 - List NATS account rules
* [potctl nats account-rules update-yaml](potctl_nats_account-rules_update-yaml.md)	 - Update NATS account rule from YAML file


