## potctl deploy

Deploy Edge Compute Network components on existing infrastructure

### Synopsis

Deploy Edge Compute Network components on existing infrastructure.
Visit iofog.org to view all YAML specifications usable with this command.

### Deployment requirements

- **Container engine:** Deployments that run the agent or controller in a container require **Docker 25+** or **Podman 4+** on the target host. On Debian/Ubuntu/Raspbian and Fedora/CentOS/RHEL/OL/SLES/openSUSE, the deploy scripts can install a supported engine automatically. On other operating systems, you must install Docker 25+ or Podman 4+ yourself; the scripts will only verify presence and version, then configure and start the engine.
- **Native agent:** The **native** (package-managed) agent is supported only on **deb/rpm-based** distributions (**Debian, Ubuntu, Raspbian, Fedora, CentOS, RHEL, OL, SLES, openSUSE**) with **systemd**. On all other OSes or init systems, use the **container agent** on that host.
- **Container agent and controller:** The container-based agent and controller support multiple init systems. The install scripts detect the init system and install an appropriate service:
  - **systemd:** systemd unit (Docker) or Quadlet unit (Podman)
  - **sysvinit, openrc, s6, runit, upstart:** init-specific scripts that run the container with the same configuration
- **Airgap:** Airgap deployments do not install a container engine. The host must have Docker 25+ or Podman 4+ already installed; the script will detect which is available, verify the version, then configure and start it.

```
potctl deploy [flags]
```

### Examples

```
deploy -f ecn.yaml
          application-template.yaml
          application.yaml
          microservice.yaml
          edge-resource.yaml
          catalog.yaml
          volume.yaml
          route.yaml
          secret.yaml
          configmap.yaml
          service.yaml
          volume-mount.yaml
```

### Options

```
  -f, --file string         YAML file containing specifications for ioFog resources to deploy
  -h, --help                help for deploy
      --no-cache            Disable caching for OfflineImage images after download
      --transfer-pool int   Maximum number of concurrent OfflineImage transfers (default 2)
```

### Options inherited from parent commands

```
      --debug              Toggle for displaying verbose output of API clients (HTTP and SSH)
  -n, --namespace string   Namespace to execute respective command within (default "default")
  -v, --verbose            Toggle for displaying verbose output of potctl
```

### SEE ALSO

* [potctl](potctl.md)	 - 


