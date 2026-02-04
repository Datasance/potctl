#!/bin/sh
# Script to install Docker/Podman based on Linux distribution
# Sources init.sh for distribution detection

set -x
set -e

CONTAINER_ENGINE_MSG="This operating system does not support automatic container engine installation. Please install Docker 25+ or Podman 4+ on the target host and re-run, or use an airgap deployment with a pre-installed engine."

check_docker_version() {
    docker_version_num=0
    if command -v docker >/dev/null 2>&1; then
        raw=$(docker -v 2>/dev/null | sed -n 's/.*version \([0-9][0-9]*\.[0-9][0-9]*\).*/\1/p' | tr -d '.')
        [ -n "$raw" ] && docker_version_num="$raw"
    fi
    [ "$docker_version_num" -ge 2500 ] 2>/dev/null || return 1
}

check_podman_version() {
    podman_version_num=0
    if command -v podman >/dev/null 2>&1; then
        raw=$(podman --version 2>/dev/null | sed -n 's/.*version \([0-9][0-9]*\).*/\1/p')
        [ -n "$raw" ] && podman_version_num="$raw"
    fi
    [ "$podman_version_num" -ge 4 ] 2>/dev/null || return 1
}

start_docker() {
    set +e
    if $sh_c "docker ps" >/dev/null 2>&1; then
        set -e
        return 0
    fi
    err_code=1
    case "${INIT_SYSTEM:-unknown}" in
        systemd)
            $sh_c "systemctl start docker" >/dev/null 2>&1
            err_code=$?
            ;;
        sysvinit)
            $sh_c "service docker start" >/dev/null 2>&1 || $sh_c "/etc/init.d/docker start" >/dev/null 2>&1
            err_code=$?
            ;;
        openrc)
            $sh_c "rc-service docker start" >/dev/null 2>&1
            err_code=$?
            ;;
        *)
            $sh_c "/etc/init.d/docker start" >/dev/null 2>&1
            err_code=$?
            [ $err_code -ne 0 ] && $sh_c "systemctl start docker" >/dev/null 2>&1 && err_code=0
            [ $err_code -ne 0 ] && $sh_c "service docker start" >/dev/null 2>&1 && err_code=0
            [ $err_code -ne 0 ] && $sh_c "snap start docker" >/dev/null 2>&1 && err_code=0
            ;;
    esac
    set -e
    if [ $err_code -ne 0 ]; then
        echo "Could not start Docker daemon"
        exit 1
    fi
}

start_podman() {
    set +e
    case "${INIT_SYSTEM:-unknown}" in
        systemd)
            $sh_c "systemctl start podman" >/dev/null 2>&1
            $sh_c "systemctl start podman.socket" >/dev/null 2>&1
            ;;
        sysvinit)
            $sh_c "service podman start" >/dev/null 2>&1 || $sh_c "/etc/init.d/podman start" >/dev/null 2>&1
            ;;
        openrc)
            $sh_c "rc-service podman start" >/dev/null 2>&1
            ;;
        *)
            $sh_c "systemctl start podman" >/dev/null 2>&1 || true
            $sh_c "systemctl start podman.socket" >/dev/null 2>&1 || true
            $sh_c "service podman start" >/dev/null 2>&1 || true
            ;;
    esac
    set -e
}


do_modify_daemon() {
    # Skip for Podman installations
    if [ "$USE_PODMAN" = "true" ]; then
        echo "# Configuring Podman for CDI directory support..."

        # Create CDI directories
        $sh_c "mkdir -p /etc/cdi /var/run/cdi"

        # Ensure /etc/containers exists
        $sh_c "mkdir -p /etc/containers"

        # Create containers.conf if it doesn't exist
        if [ ! -f "/etc/containers/containers.conf" ]; then
            $sh_c 'cat > /etc/containers/containers.conf <<EOF
[engine]
runtime = "crun"
cdi_spec_dirs = ["/etc/cdi", "/var/run/cdi"]
EOF'
        else
            # Check if [engine] block exists
            if grep -q "^\[engine\]" /etc/containers/containers.conf; then
                # Ensure runtime is set under [engine]
                if grep -q "^runtime" /etc/containers/containers.conf; then
                    $sh_c "sed -i 's|^runtime *=.*|runtime = \"crun\"|' /etc/containers/containers.conf"
                else
                    $sh_c "sed -i '/^\[engine\]/a runtime = \"crun\"' /etc/containers/containers.conf"
                fi

                # Ensure cdi_spec_dirs is set under [engine]
                if grep -q "^cdi_spec_dirs" /etc/containers/containers.conf; then
                    $sh_c "sed -i 's|^cdi_spec_dirs *=.*|cdi_spec_dirs = [\"/etc/cdi\", \"/var/run/cdi\"]|' /etc/containers/containers.conf"
                else
                    $sh_c "sed -i '/^\[engine\]/a cdi_spec_dirs = [\"/etc/cdi\", \"/var/run/cdi\"]' /etc/containers/containers.conf"
                fi
            else
                # Append full engine block if missing
                $sh_c 'echo -e "\n[engine]\nruntime = \"crun\"\ncdi_spec_dirs = [\"/etc/cdi\", \"/var/run/cdi\"]" >> /etc/containers/containers.conf'
            fi
        fi

        case "${INIT_SYSTEM:-unknown}" in
            systemd)
                $sh_c "systemctl enable podman" 2>/dev/null || true
                $sh_c "systemctl enable podman.socket" 2>/dev/null || true
                ;;
            openrc)
                $sh_c "rc-update add podman default" 2>/dev/null || true
                ;;
            sysvinit)
                $sh_c "update-rc.d podman defaults" 2>/dev/null || $sh_c "chkconfig podman on" 2>/dev/null || true
                ;;
            *) ;;
        esac
        start_podman
        return
    fi

    # Original Docker daemon configuration
    if [ ! -f /etc/docker/daemon.json ]; then
        echo "Creating /etc/docker/daemon.json..."
        $sh_c "mkdir -p /etc/docker"
        $sh_c 'cat > /etc/docker/daemon.json << EOF
{
	"storage-driver": "overlayfs",
    "features": {
        "containerd-snapshotter": true,
        "cdi": true
    },
    "cdi-spec-dirs": ["/etc/cdi/", "/var/run/cdi"]
}
EOF'
    else
        echo "/etc/docker/daemon.json already exists"
    fi
    echo "Restarting Docker daemon..."
    case "${INIT_SYSTEM:-unknown}" in
        systemd)
            $sh_c "systemctl daemon-reload"
            $sh_c "systemctl restart docker"
            ;;
        *)
            $sh_c "systemctl daemon-reload" 2>/dev/null || true
            $sh_c "systemctl restart docker" 2>/dev/null || start_docker
            ;;
    esac
}

do_set_datasance_repo() {
    echo "# Setting up Datasance repository for $lsb_dist..."
    
    case "$lsb_dist" in
        fedora|centos|rhel|ol|sles|opensuse*)
            # RPM-based distros
            $sh_c "cd /etc/yum.repos.d && curl -s https://downloads.datasance.com/datasance.repo -LO"
            if [ "$lsb_dist" = "fedora" ] || [ "$lsb_dist" = "centos" ] || [ "$lsb_dist" = "rhel" ] || [ "$lsb_dist" = "ol" ]; then
                $sh_c "yum update -y"
            else
                $sh_c "zypper refresh"
            fi
        ;;
        debian|ubuntu|raspbian|*)
            # DEB-based distros
            $sh_c "apt update -qy"
            $sh_c "apt install -qy debian-archive-keyring apt-transport-https"
            $sh_c "wget -qO- https://downloads.datasance.com/datasance.gpg | tee /etc/apt/trusted.gpg.d/datasance.gpg >/dev/null"
            $sh_c "echo 'deb [arch=all signed-by=/etc/apt/trusted.gpg.d/datasance.gpg] https://downloads.datasance.com/deb stable main' | tee /etc/apt/sources.list.d/datasance.list >/dev/null"
            $sh_c "apt update -qy"
        ;;
    esac
}

do_install_wasm_shim() {
    echo "# Installing WebAssembly runtime support for $lsb_dist..."
    arch=$(uname -m)

    # Normalize architecture for consistency
    case "$arch" in
        arm64|aarch64|armv7l|armv8) arch="aarch64" ;;
        amd64|x86_64) arch="x86_64" ;;
    esac

    if [ "$USE_PODMAN" = "true" ]; then
        case "$lsb_dist" in
            fedora|centos|rhel|ol)
                $sh_c "yum install -y crun crun-wasm"
            ;;
            sles|opensuse*)
                $sh_c "zypper install -y crun"
                # Note: crun-wasm might not be available in SUSE repos
                # In that case, we'll need to use the standard crun
            ;;
        esac
        
    else
        # Original containerd WASM shim installation for Docker
        case "$lsb_dist" in
            debian|raspbian|ubuntu)
                $sh_c "apt update -qy" || { echo "Failed to update apt packages"; exit 1; }

                if [ "$arch" = "aarch64" ]; then
                    $sh_c "apt install -qy containerd-shim-wasmedge-v1-aarch64-linux-gnu"
                    $sh_c "apt install -qy containerd-shim-wasmer-v1-aarch64-linux-gnu"
                    $sh_c "apt install -qy containerd-shim-wasmtime-v1-aarch64-linux-gnu"
                elif [ "$arch" = "x86_64" ]; then
                    $sh_c "apt install -qy containerd-shim-wasmedge-v1-x86-64-linux-gnu"
                    $sh_c "apt install -qy containerd-shim-wasmer-v1-x86-64-linux-gnu"
                    $sh_c "apt install -qy containerd-shim-wasmtime-v1-x86-64-linux-gnu"
                else
                    echo "Unsupported architecture: $arch for $lsb_dist"
                    exit 1
                fi
            ;;
            *)
                echo "Unsupported OS: $lsb_dist"
                exit 1
            ;;
        esac
    fi
}

do_install_container_engine() {
    if [ "$PACKAGE_TYPE" = "other" ]; then
        if check_docker_version; then
            USE_PODMAN="false"
            echo "# Docker (>= 25) found; using Docker."
            start_docker
            do_modify_daemon
            return 0
        fi
        if check_podman_version; then
            USE_PODMAN="true"
            echo "# Podman (>= 4) found; using Podman."
            do_modify_daemon
            return 0
        fi
        echo "Error: $CONTAINER_ENGINE_MSG"
        exit 1
    fi

    if [ "$USE_PODMAN" = "true" ]; then
        echo "# Installing Podman and related packages..."
        case "$lsb_dist" in
            fedora|centos|rhel|ol)
                $sh_c "yum install -y podman crun podman-docker"
            ;;
            sles|opensuse*)
                $sh_c "zypper install -y podman crun podman-docker"
            ;;
        esac
        if ! check_podman_version; then
            echo "Error: Podman 4+ is required. Please upgrade Podman."
            exit 1
        fi
        do_modify_daemon
        return
    fi

    if command_exists docker; then
        docker_version=$(docker -v 2>/dev/null | sed -n 's/.*version \([0-9][0-9]*\.[0-9][0-9]*\).*/\1/p' | tr -d '.')
        if [ -n "$docker_version" ] && [ "$docker_version" -ge 2500 ] 2>/dev/null; then
            echo "# Docker already installed (>= 25)"
            start_docker
            do_modify_daemon
            return
        fi
    fi

    echo "# Installing Docker..."
    case "$lsb_dist" in
        debian|ubuntu|raspbian)
            case "$dist_version" in
                "stretch")
                    $sh_c "apt install -y apt-transport-https ca-certificates curl gnupg2 software-properties-common"
                    curl -fsSL https://download.docker.com/linux/debian/gpg | $sh_c "apt-key add -"
                    $sh_c "add-apt-repository \"deb [arch=$(dpkg --print-architecture)] https://download.docker.com/linux/debian $(lsb_release -cs) stable\""
                    $sh_c "apt update -y"
                    $sh_c "apt install -y docker-ce"
                ;;
                *)
                    curl -fsSL https://get.docker.com/ | $sh_c "sh"
                ;;
            esac
        ;;
        *)
            curl -fsSL https://get.docker.com/ | $sh_c "sh"
        ;;
    esac

    if ! command_exists docker; then
        echo "Failed to install Docker"
        exit 1
    fi
    if ! check_docker_version; then
        echo "Error: Docker 25+ is required. Please upgrade Docker."
        exit 1
    fi
    start_docker
    do_modify_daemon
}

# Check if we should use Podman based on distribution
determine_container_engine() {
    USE_PODMAN="false"
    case "$lsb_dist" in
        fedora|centos|rhel|ol|sles|opensuse*)
            USE_PODMAN="true"
            echo "# Using Podman for $lsb_dist"
        ;;
        *)
            echo "# Using Docker for $lsb_dist"
        ;;
    esac
}

# Source init.sh to get distribution info
. /etc/iofog/controller/init.sh
init

# Configure container engine based on distribution
determine_container_engine

# Install appropriate container engine
do_install_container_engine

if [ "$PACKAGE_TYPE" = "deb" ] || [ "$PACKAGE_TYPE" = "rpm" ]; then
    do_set_datasance_repo
    do_install_wasm_shim
fi

echo "# Installation completed successfully"