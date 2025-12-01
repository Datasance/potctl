#!/bin/sh
set -x
set -e

AGENT_LOG_FOLDER=/var/log/iofog-agent
AGENT_BACKUP_FOLDER=/var/backups/iofog-agent
AGENT_MESSAGE_FOLDER=/var/lib/iofog-agent
AGENT_SHARE_FOLDER=/usr/share/iofog-agent
SAVED_AGENT_CONFIG_FOLDER=/tmp/agent-config-save
AGENT_CONTAINER_NAME="iofog-agent"
ETC_DIR=/etc/iofog/agent

do_check_install() {
	if command_exists iofog-agent; then
		local VERSION=$(sudo iofog-agent version | head -n1 | sed "s/ioFog//g" | tr -d ' ' | tr -d "\n")
		if [ "$VERSION" = "$agent_version" ]; then
			echo "Agent $VERSION already installed."
			exit 0
		fi
	fi
}

do_stop_iofog() {
	if command_exists iofog-agent; then
		sudo systemctl stop iofog-agent
	fi
}

# do_check_iofog_on_arm() {
#   if [ "$lsb_dist" = "raspbian" ] || [ "$(uname -m)" = "armv7l" ] || [ "$(uname -m)" = "aarch64" ] || [ "$(uname -m)" = "armv8" ]; then
#     echo "# We re on ARM ($(uname -m)) : Updating config.xml to use correct docker_url"
#     $sh_c 'sed -i -e "s|<docker_url>.*</docker_url>|<docker_url>tcp://127.0.0.1:2375/</docker_url>|g" /etc/iofog-agent/config.xml'

#     echo "# Restarting iofog-agent service"
#     $sh_c "systemctl stop iofog-agent"
#     sleep 3
#     $sh_c "systemctl start iofog-agent"
#  fi
# }

do_create_env() {
ENV_FILE_NAME=iofog-agent.env # Used as an env file in systemd

ENV_FILE="$ETC_DIR/$ENV_FILE_NAME"

# Env file (for systemd)
rm -f "$ENV_FILE"
touch "$ENV_FILE"

echo "IOFOG_AGENT_IMAGE=${agent_image}" >> "$ENV_FILE"
echo "IOFOG_AGENT_TZ=${agent_tz}" >> "$ENV_FILE"

}

do_install_iofog() {
	echo "# Installing ioFog agent..."
	
    # 1. Ensure folders exist
    for FOLDER in ${ETC_DIR} ${AGENT_LOG_FOLDER} ${AGENT_BACKUP_FOLDER} ${AGENT_MESSAGE_FOLDER} ${AGENT_SHARE_FOLDER}; do
        if [ ! -d "$FOLDER" ]; then
            echo "Creating folder: $FOLDER"
            sudo mkdir -p "$FOLDER"
            sudo chmod 775 "$FOLDER"
        fi
    done
	do_create_env

    # Check if we're using Podman (RHEL, CentOS, Fedora, etc.)
    if [ "$lsb_dist" = "rhel" ] || [ "$lsb_dist" = "centos" ] || [ "$lsb_dist" = "fedora" ] || [ "$lsb_dist" = "ol" ] || [ "$lsb_dist" = "sles" ] || [ "$lsb_dist" = "opensuse" ]; then
        echo "Using Podman for container management..."
        SYSTEMD_SERVICE_FILE=/etc/containers/systemd/iofog-agent.container

        cat <<EOF | sudo tee ${SYSTEMD_SERVICE_FILE} > /dev/null
[Unit]
Description=Datasance PoT IoFog Agent Service
After=podman.service
Requires=podman.service

[Container]
ContainerName=${AGENT_CONTAINER_NAME}
Image=${agent_image}
PodmanArgs=--privileged --stop-timeout=60
EnvironmentFile=${ETC_DIR}/iofog-agent.env
Network=host
Volume=/run/podman/podman.sock:/run/podman/podman.sock:rw
Volume=iofog-agent-config:/etc/iofog-agent:rw
Volume=/var/log/iofog-agent:/var/log/iofog-agent:rw
Volume=/var/backups/iofog-agent:/var/backups/iofog-agent:rw
Volume=/usr/share/iofog-agent:/usr/share/iofog-agent:rw
Volume=/var/lib/iofog-agent:/var/lib/iofog-agent:rw
LogDriver=journald

[Service]
Restart=always

[Install]
WantedBy=default.target
EOF

        # Reload systemd and enable the service
        sudo systemctl daemon-reload
        sudo systemctl restart podman
        sudo systemctl start iofog-agent.service

        # Create the iofog-agent executable script for Podman
        EXECUTABLE_FILE=/usr/local/bin/iofog-agent
        cat <<'EOF' | sudo tee ${EXECUTABLE_FILE} > /dev/null
#!/bin/bash
CONTAINER_NAME="iofog-agent"

# Check if the container is running
if ! podman ps --format '{{.Names}}' | grep -q "^${CONTAINER_NAME}$"; then
    echo "Error: The iofog-agent container is not running."
    exit 1
fi

# Execute the command in the container
podman exec ${CONTAINER_NAME} iofog-agent "$@"
EOF
    else
        # Docker-based installation (Debian, Ubuntu, etc.)
        echo "Using Docker for container management..."
        SYSTEMD_SERVICE_FILE=/etc/systemd/system/iofog-agent.service

        cat <<EOF | sudo tee ${SYSTEMD_SERVICE_FILE} > /dev/null
[Unit]
Description=Datasance PoT IoFog Agent Service
After=docker.service
Requires=docker.service

[Service]
Restart=always
ExecStartPre=-/usr/bin/docker rm -f ${AGENT_CONTAINER_NAME}
ExecStart=/usr/bin/docker run --rm --name ${AGENT_CONTAINER_NAME} \\
--env-file ${ETC_DIR}/iofog-agent.env \\
-v /var/run/docker.sock:/var/run/docker.sock:rw \\
-v iofog-agent-config:/etc/iofog-agent:rw \\
-v /var/log/iofog-agent:/var/log/iofog-agent:rw \\
-v /var/backups/iofog-agent:/var/backups/iofog-agent:rw \\
-v /usr/share/iofog-agent:/usr/share/iofog-agent:rw \\
-v /var/lib/iofog-agent:/var/lib/iofog-agent:rw \\
--net=host \\
--privileged \\
--stop-timeout 60 \\
--attach stdout \\
--attach stderr \\
${agent_image}
ExecStop=/usr/bin/docker stop ${AGENT_CONTAINER_NAME}

[Install]
WantedBy=default.target
EOF

        # Reload systemd and enable the service
        sudo systemctl daemon-reload
        sudo systemctl enable iofog-agent.service

        # Create the iofog-agent executable script for Docker
        EXECUTABLE_FILE=/usr/local/bin/iofog-agent
        cat <<'EOF' | sudo tee ${EXECUTABLE_FILE} > /dev/null
#!/bin/bash
CONTAINER_NAME="iofog-agent"

# Check if the container is running
if ! docker ps --format '{{.Names}}' | grep -q "^${CONTAINER_NAME}$"; then
    echo "Error: The iofog-agent container is not running."
    exit 1
fi

# Execute the command in the container
docker exec ${CONTAINER_NAME} iofog-agent "$@"
EOF
    fi

    # Make the script executable
    sudo chmod +x ${EXECUTABLE_FILE}

    echo "ioFog agent installation completed!"
}

do_start_iofog(){
	sudo systemctl start iofog-agent > /dev/null 2&>1 &
	local STATUS=""
	local ITER=0
	while [ "$STATUS" != "RUNNING" ] ; do
    ITER=$((ITER+1))
    if [ "$ITER" -gt 600 ]; then
      echo 'Timed out waiting for Agent to be RUNNING'
      exit 1;
    fi
    sleep 1
    STATUS=$(sudo iofog-agent status | cut -f2 -d: | head -n 1 | tr -d '[:space:]')
    echo "${STATUS}"
	done
	sudo iofog-agent "config -cf 10 -sf 10"
    if [ "$lsb_dist" = "rhel" ] || [ "$lsb_dist" = "centos" ] || [ "$lsb_dist" = "fedora" ] || [ "$lsb_dist" = "ol" ] || [ "$lsb_dist" = "sles" ] || [ "$lsb_dist" = "opensuse" ]; then
        sudo iofog-agent "config -c unix:///var/run/podman/podman.sock"
    fi    
}

agent_image="$1"
agent_tz="$2"
echo "Using variables"
echo "version: $agent_image"
echo "timezone: $agent_tz"
. /etc/iofog/agent/init.sh
init
do_check_install
do_stop_iofog
do_install_iofog
do_start_iofog