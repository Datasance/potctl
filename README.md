![potctl-logo](potctl-logo.png?raw=true "potctl logo")

`potctl` is a CLI for the installation, configuration, and operation of ioFog 
[Edge Compute Networks](https://docs.datasance.com/#/./ioFog_3.0/getting-started/core-concepts) (ECNs).
It can be used to remotely manage multiple ECNs from a single host. It is built for ioFog users and DevOps engineers 
wanting to manage ECNs. It is modelled on existing tools such as Terraform or kubectl that can be used to automate
infrastructure-as-code.

## Install

#### Mac

Mac users can use Homebrew:

```bash
brew tap datasance/potctl
brew install potctl
```

#### Linux

The Debian package can be installed like so:
```bash
https://packagecloud.io/install/repositories/datasance/potctl/script.deb.sh | sudo bash
sudo apt install potctl
```

And similarly, the RPM package can be installed like so:
```
https://packagecloud.io/install/repositories/datasance/potctl/script.rpm.sh | sudo bash
sudo apt install potctl
```

## Usage

### Documentation

The best way to learn how to use `potctl` is through the [iofog.org](https://docs.datasance.com/#/./ioFog_3.0/getting-started/quick-start-local) learning resources.

#### Quick Start

See all potctl options

```
potctl --help
```

Current options include:

```
██████╗  ██████╗ ████████╗ ██████╗████████╗██╗     
██╔══██╗██╔═══██╗╚══██╔══╝██╔════╝╚══██╔══╝██║     
██████╔╝██║   ██║   ██║   ██║        ██║   ██║     
██╔═══╝ ██║   ██║   ██║   ██║        ██║   ██║     
██║     ╚██████╔╝   ██║   ╚██████╗   ██║   ███████╗
╚═╝      ╚═════╝    ╚═╝    ╚═════╝   ╚═╝   ╚══════╝
                                                   


Welcome to the cool new potctl Cli!

Use `potctl version` to display the current version.


Usage:
  potctl [flags]
  potctl [command]

Available Commands:
  attach        Attach an existing ioFog resource to Control Plane
  configure     Configure potctl or ioFog resources
  connect       Connect to an existing Control Plane
  create        Create a resource
  delete        Delete an existing ioFog resource
  deploy        Deploy Edge Compute Network components on existing infrastructure
  describe      Get detailed information of existing resources
  detach        Detach an existing ioFog resource from its ECN
  disconnect    Disconnect from an ioFog cluster
  get           Get information of existing resources
  help          Help about any command
  legacy        Execute commands using legacy CLI
  logs          Get log contents of deployed resource
  move          Move an existing resources inside the current Namespace
  prune         prune ioFog resources
  rename        Rename the iofog resources that are currently deployed
  start         Starts a resource
  stop          Stops a resource
  version       Get CLI application version
  view          Open ECN Viewer

Flags:
      --detached           Use/Show detached resources
  -h, --help               help for potctl
      --http-verbose       Toggle for displaying verbose output of API client
  -n, --namespace string   Namespace to execute respective command within (default "default")
  -v, --verbose            Toggle for displaying verbose output of potctl

Use "potctl [command] --help" for more information about a command.

```

### Autocomplete

If you are running BASH or ZSH, potctl comes with shell autocompletion scripts.
In order to generate those scripts, run:

```bash
potctl autocomplete bash
```
OR

```bash
potctl autocomplete zsh
```

Then follow the instructions output by the command.

Example:
```bash
potctl autocomplete bash
✔ $HOME/.iofog/completion.bash.sh generated
Run `source $HOME/.iofog/completion.bash.sh` to update your current session
Add `source $HOME/.iofog/completion.bash.sh` to your bash profile to have it saved

source $HOME/.iofog/completion.bash.sh
echo "$HOME/.iofog/completion.bash.sh" >> $HOME/.bash_profile
```

## Build from Source

This project uses go modules so it must be built from outside of your $GOPATH.

Go 1.16+ is a prerequisite. Install all other dependancies with:
```
make bootstrap
```
Make sure that your `$PATH` contains `$GOBIN`, otherwise `make build` will fail on the basis that command `rice` is not found.

See all `make` commands by running:
```
make
```

To build and install, go ahead and run:
```
make build install
potctl --help
```

potctl is installed in `/usr/local/bin`

## Running Tests

Run project unit tests:
```
make test
```

This will output a JUnit compatible file into `reports/TEST-potctl.xml` that can be imported in most CI systems.

## Embed assets

Run project build:
```
make build
```
