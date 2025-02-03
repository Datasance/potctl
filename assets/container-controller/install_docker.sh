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


do_install_docker() {
	# Check that Docker 25.0.0 or greater is installed
	if command_exists docker; then
		docker_version=$(docker -v | sed 's/.*version \(.*\),.*/\1/' | tr -d '.')
		if [ "$docker_version" -ge 2600 ]; then
			echo "# Docker $docker_version already installed"
			start_docker
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
}

. /etc/iofog/controller/init.sh
init
do_install_docker
