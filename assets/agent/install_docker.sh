#!/bin/sh
set -x
set -e

start_docker() {
	set +e
	# check if docker is running
	if ! $sh_c "docker ps" >/dev/null 2>&1; then
		# Try init.d
		$sh_c "/etc/init.d/docker start"
		local err_code=$?
		# Try systemd
		if [ $err_code -ne 0 ]; then
			$sh_c "service docker start"
			err_code=$?
		fi
		# Try snapd
		if [ $err_code -ne 0 ]; then
			$sh_c "snap docker start"
			err_code=$?
		fi
		if [ $err_code -ne 0 ]; then
			echo "Could not start Docker daemon"
			exit 1
		fi
	fi
	set -e
}

do_configure_overlay() {
	local driver="$DOCKER_STORAGE_DRIVER"
	if [ -z "$driver" ]; then
		driver="overlay2"
	fi
	echo "# Configuring /etc/systemd/system/docker.service.d/overlay.conf..."
	if [ "$lsb_dist" = "raspbian" ] || [ "$(uname -m)" = "armv7l" ] || [ "$(uname -m)" = "aarch64" ] || [ "$(uname -m)" = "armv8" ]; then
		if [ ! -d "/etc/systemd/system/docker.service.d" ]; then
			$sh_c "mkdir -p /etc/systemd/system/docker.service.d"
		fi
		if [ ! -f "/etc/systemd/system/docker.service.d/overlay.conf" ] || ! grep -Fxq "ExecStart=/usr/bin/dockerd --storage-driver $driver -H unix:// -H tcp://127.0.0.1:2375" "/etc/systemd/system/docker.service.d/overlay.conf"; then
			$sh_c 'echo "[Service]" > /etc/systemd/system/docker.service.d/overlay.conf'
			$sh_c 'echo "ExecStart=" >> /etc/systemd/system/docker.service.d/overlay.conf'
			$sh_c "echo \"ExecStart=/usr/bin/dockerd --storage-driver $driver -H unix:// -H tcp://127.0.0.1:2375\" >> /etc/systemd/system/docker.service.d/overlay.conf"
		fi
		$sh_c "systemctl daemon-reload"
		$sh_c "service docker restart"
	fi
}

do_modify_daemon() {
	if [ ! -f /etc/docker/daemon.json ]; then
		echo "Creating /etc/docker/daemon.json..."
		sudo tee /etc/docker/daemon.json > /dev/null <<EOF
{
	"features": {
		"containerd-snapshotter": true,
		"cdi": true
	},
	"cdi-spec-dirs": ["/etc/cdi/", "/var/run/cdi"]
}
EOF
	else
		echo "/etc/docker/daemon.json already exists"
	fi
	echo "Restarting Docker daemon..."
	$sh_c "systemctl daemon-reload"
	$sh_c "service docker restart"
}

do_set_datasance_repo() {

    echo $lsb_dist
	if [ "$lsb_dist" = "fedora" ] || [ "$lsb_dist" = "centos" ]; then

		cd /etc/yum.repos.d ; curl https://downloads.datasance.com/datasance.repo -LO
		$sh_c "yum update"
	else
	$sh_c "apt update -qy"
    $sh_c "apt install -qy debian-archive-keyring"
    $sh_c "apt install -qy apt-transport-https"
	sudo wget -qO- https://downloads.datasance.com/datasance.gpg | sudo tee /etc/apt/trusted.gpg.d/datasance.gpg >/dev/null
	echo "deb [arch=all signed-by=/etc/apt/trusted.gpg.d/datasance.gpg] https://downloads.datasance.com/deb stable main" | sudo tee /etc/apt/sources.list.d/datansance.list >/dev/null
    $sh_c "apt update -qy"
	fi

}

do_install_wasm_shim() {
    echo "Detected OS: $lsb_dist"
    arch=$(uname -m)

    # Normalize architecture for consistency
    case "$arch" in
        arm64|aarch64|armv7l|armv8) arch="aarch64" ;;
        amd64|x86_64) arch="x86_64" ;;
    esac

    if [ "$lsb_dist" = "fedora" ] || [ "$lsb_dist" = "centos" ]; then
        $sh_c "yum update -y" || { echo "Failed to update yum packages"; exit 1; }

        if [ "$arch" = "aarch64" ]; then
            $sh_c "yum install -y containerd-shim-wasmedge-v1-aarch64-linux-gnu-0.7.0-1.aarch64"
            $sh_c "yum install -y containerd-shim-wasmer-v1-aarch64-linux-gnu-0.7.0-1.aarch64"
            $sh_c "yum install -y containerd-shim-wasmtime-v1-aarch64-linux-gnu-0.7.0-1.aarch64"
        elif [ "$arch" = "x86_64" ]; then
            $sh_c "yum install -y containerd-shim-wasmedge-v1-x86_64-linux-gnu-0.7.0-1.x86_64"
            $sh_c "yum install -y containerd-shim-wasmer-v1-x86_64-linux-gnu-0.7.0-1.x86_64"
            $sh_c "yum install -y containerd-shim-wasmtime-v1-x86_64-linux-gnu-0.7.0-1.x86_64"
        else
            echo "Unsupported architecture: $arch for Fedora/CentOS"
            exit 1
        fi
    elif [ "$lsb_dist" = "debian" ] || [ "$lsb_dist" = "raspbian" ] || [ "$lsb_dist" = "ubuntu" ]; then
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
            echo "Unsupported architecture: $arch for Debian/Ubuntu"
            exit 1
        fi
    else
        echo "Unsupported OS: $lsb_dist"
        exit 1
    fi
}

do_install_docker() {
	# Check that Docker 25.0.0 or greater is installed
	if command_exists docker; then
		docker_version=$(docker -v | sed 's/.*version \(.*\),.*/\1/' | tr -d '.')
		if [ "$docker_version" -ge 2600 ]; then
			echo "# Docker $docker_version already installed"
			start_docker
			do_configure_overlay
			do_modify_daemon
			return
		fi
	fi
	echo "# Installing Docker..."
	case "$dist_version" in
		"stretch")
			$sh_c "apt install -y apt-transport-https ca-certificates curl gnupg2 software-properties-common"
			curl -fsSL https://download.docker.com/linux/debian/gpg | $sh_c "apt-key add -"
			$sh_c "sudo add-apt-repository \"deb [arch=$(dpkg --print-architecture) https://download.docker.com/linux/debian $(lsb_release -cs) stable\""
			$sh_c "apt-get update -y"
			$sh_c "sudo apt install -y docker-ce"
		;;
    7|8)
      $sh_c "sudo yum install -y yum-utils || echo 'yum-utils already installed'"
      $sh_c "sudo yum-config-manager \
            --add-repo \
            https://download.docker.com/linux/centos/docker-ce.repo"
      $sh_c "sudo yum install docker-ce docker-ce-cli containerd.io -y"
    ;;
		*)
			curl -fsSL https://get.docker.com/ | sh
		;;
	esac
	
	if ! command_exists docker; then
		echo "Failed to install Docker"
		exit 1
	fi
	start_docker
	do_configure_overlay
	do_modify_daemon
}

. /etc/iofog/agent/init.sh
init
do_install_docker
do_set_datasance_repo
do_install_wasm_shim
