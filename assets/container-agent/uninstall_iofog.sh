#!/bin/sh
set -x
set -e

AGENT_CONFIG_FOLDER=iofog-agent-config
AGENT_LOG_FOLDER=/var/log/iofog-agent
AGENT_BACKUP_FOLDER=/var/backups/iofog-agent
AGENT_MESSAGE_FOLDER=/var/lib/iofog-agent
SYSTEMD_SERVICE_FILE=/etc/systemd/system/iofog-agent.service
EXECUTABLE_FILE=/usr/local/bin/iofog-agent
CONTAINER_NAME="iofog-agent"

do_uninstall_iofog() {
    echo "# Removing ioFog agent..."

    # Disable and stop the systemd service
    if [ -f ${SYSTEMD_SERVICE_FILE} ]; then
        echo "Disabling and stopping the systemd service..."
        sudo systemctl stop iofog-agent.service || true
        sudo systemctl disable iofog-agent.service || true
        sudo rm -f ${SYSTEMD_SERVICE_FILE}
        sudo systemctl daemon-reload
    fi

    # Remove the Docker container
    if docker ps --format '{{.Names}}' | grep -q "^${CONTAINER_NAME}$"; then
        echo "Stopping and removing the ioFog agent container..."
        docker stop ${CONTAINER_NAME}
        docker rm ${CONTAINER_NAME}
    fi

    # Remove config files
    echo "Checking if the Docker volume exists..."

    if sudo docker volume inspect "${AGENT_CONFIG_FOLDER}" >/dev/null 2>&1; then
        echo "Docker volume '${AGENT_CONFIG_FOLDER}' found. Removing..."
        sudo docker volume rm "${AGENT_CONFIG_FOLDER}"
        echo "Docker volume '${AGENT_CONFIG_FOLDER}' has been removed."
    else
        echo "Docker volume '${AGENT_CONFIG_FOLDER}' does not exist. Skipping removal."
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
