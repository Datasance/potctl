![potctl-logo](potctl-logo.png?raw=true "potctl logo")

`potctl` is a CLI for the installation, configuration, and operation of ioFog 
[Edge Compute Networks](https://docs.datasance.com/getting-started/core-concepts) (ECNs).
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
wget -qO- https://downloads.datasance.com/datasance.gpg | sudo tee /etc/apt/trusted.gpg.d/datasance.gpg >/dev/null
echo "deb [arch=all signed-by=/etc/apt/trusted.gpg.d/datasance.gpg] https://downloads.datasance.com/deb stable main" | sudo tee /etc/apt/sources.list.d/datansance.list >/dev/null
sudo apt update
sudo apt install potctl -y
```

And similarly, the RPM package can be installed like so:
```
cd /etc/yum.repos.d ; curl https://downloads.datasance.com/datasance.repo -LO
sudo yum update
sudo yum install potctl
```

## Usage

### Documentation

The best way to learn how to use `potctl` is through the [docs.datasance.com](https://docs.datasance.com/getting-started/core-concepts) learning resources.

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
                                                   


Potctl is the CLI for Datasance PoT, an Enterprise version of Eclipse iofog. Think of it as a mix between terraform and kubectl.

Use `potctl version` to display the current version.

Find more information at: https://docs.datasance.com 


Usage:
  potctl [flags]
  potctl [command]

Available Commands:
  attach        Attach one ioFog resource to another
  completion    Generate the autocompletion script for the specified shell
  configure     Configure potctl or ioFog resources
  connect       Connect to an existing Control Plane
  create        Create a resource
  delete        Delete an existing ioFog resource
  deploy        Deploy Edge Compute Network components on existing infrastructure
  describe      Get detailed information of an existing resources
  detach        Detach one ioFog resource from another
  disconnect    Disconnect from an ioFog cluster
  exec          Connect to an Exec Session of a resource
  get           Get information of existing resources
  help          Help about any command
  legacy        Execute commands using legacy CLI
  logs          Get log contents of deployed resource
  move          Move an existing resources inside the current Namespace
  prune         prune ioFog resources
  rebuild       Rebuilds a microservice or system-microservice
  rename        Rename the iofog resources that are currently deployed
  rollback      Rollback ioFog resources
  start         Starts a resource
  stop          Stops a resource
  upgrade       Upgrade ioFog resources
  version       Get CLI application version
  view          Open ECN Viewer

Flags:
      --detached           Use/Show detached resources
      --debug              Toggle for displaying verbose output of API clients (HTTP and SSH)
  -h, --help               help for potctl
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

Go 1.19+ is a prerequisite. Install all other dependancies with:
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
