#!/bin/sh
set -x
set -e

AGENT_CONFIG_FOLDER=iofog-agent-config
AGENT_LOG_FOLDER=/var/log/iofog-agent
AGENT_BACKUP_FOLDER=/var/backups/iofog-agent
AGENT_MESSAGE_FOLDER=/var/lib/iofog-agent
EXECUTABLE_FILE=/usr/local/bin/iofog-agent
CONTAINER_NAME="iofog-agent"

do_uninstall_iofog() {
    echo "# Removing ioFog agent..."

    # Set the appropriate systemd service file based on the Linux distribution
    if [ "$lsb_dist" = "rhel" ] || [ "$lsb_dist" = "fedora" ] || [ "$lsb_dist" = "centos" ] || [ "$lsb_dist" = "ol" ] || [ "$lsb_dist" = "sles" ] || [ "$lsb_dist" = "opensuse" ]; then
        SYSTEMD_SERVICE_FILE=/etc/containers/systemd/iofog-agent.container
        CONTAINER_RUNTIME="podman"
    else
        SYSTEMD_SERVICE_FILE=/etc/systemd/system/iofog-agent.service
        CONTAINER_RUNTIME="docker"
    fi

    # Disable and stop the systemd service
    if [ -f ${SYSTEMD_SERVICE_FILE} ]; then
        echo "Disabling and stopping the systemd service..."
        sudo systemctl stop iofog-agent.service || true
        sudo systemctl disable iofog-agent.service || true
        sudo rm -f ${SYSTEMD_SERVICE_FILE}
        sudo systemctl daemon-reload
    fi

    # Remove the container
    if ${CONTAINER_RUNTIME} ps --format '{{.Names}}' | grep -q "^${CONTAINER_NAME}$"; then
        echo "Stopping and removing the ioFog agent container..."
        ${CONTAINER_RUNTIME} stop ${CONTAINER_NAME}
        ${CONTAINER_RUNTIME} rm ${CONTAINER_NAME}
    fi

    # Remove config files
    echo "Checking if the ${CONTAINER_RUNTIME} volume exists..."

    if sudo ${CONTAINER_RUNTIME} volume inspect "${AGENT_CONFIG_FOLDER}" >/dev/null 2>&1; then
        echo "${CONTAINER_RUNTIME} volume '${AGENT_CONFIG_FOLDER}' found. Removing..."
        sudo ${CONTAINER_RUNTIME} volume rm "${AGENT_CONFIG_FOLDER}"
        echo "${CONTAINER_RUNTIME} volume '${AGENT_CONFIG_FOLDER}' has been removed."
    else
        echo "${CONTAINER_RUNTIME} volume '${AGENT_CONFIG_FOLDER}' does not exist. Skipping removal."
    fi

    # Remove log files
    echo "Removing log files..."
    sudo rm -rf ${AGENT_LOG_FOLDER}

    # Remove backup files
    echo "Removing backup files..."
    sudo rm -rf ${AGENT_BACKUP_FOLDER}

    # Remove message files
    echo "Removing message files..."
    sudo rm -rf ${AGENT_MESSAGE_FOLDER}

    # Remove the executable script
    if [ -f ${EXECUTABLE_FILE} ]; then
        echo "Removing the iofog-agent executable script..."
        sudo rm -f ${EXECUTABLE_FILE}
    fi

    echo "ioFog agent uninstalled successfully!"
}

. /etc/iofog/agent/init.sh
init

do_uninstall_iofog
