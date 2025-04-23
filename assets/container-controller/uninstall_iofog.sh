#!/bin/sh
set -x
set -e


CONTROLLER_LOG_DIR="iofog-controller-log"
CONTAINER_NAME="iofog-controller"
EXECUTABLE_FILE=/usr/local/bin/iofog-controller
CONTROLLER_DB=iofog-controller-db


do_uninstall_controller() {
    echo "# Removing ioFog controller..."

    # Set the appropriate systemd service file based on the Linux distribution
    if [ "$lsb_dist" = "rhel" ] || [ "$lsb_dist" = "fedora" ] || [ "$lsb_dist" = "centos" ] || [ "$lsb_dist" = "ol" ] || [ "$lsb_dist" = "sles" ] || [ "$lsb_dist" = "opensuse" ]; then
        SYSTEMD_SERVICE_FILE=/etc/containers/systemd/iofog-controller.container
        CONTAINER_RUNTIME="podman"
    else
        SYSTEMD_SERVICE_FILE=/etc/systemd/system/iofog-controller.service
        CONTAINER_RUNTIME="docker"
    fi

    # Disable and stop the systemd service
    if [ -f ${SYSTEMD_SERVICE_FILE} ]; then
        echo "Disabling and stopping the systemd service..."
        sudo systemctl stop iofog-controller.service || true
        sudo systemctl disable iofog-controller.service || true
        sudo rm -f ${SYSTEMD_SERVICE_FILE}
        sudo systemctl daemon-reload
    fi

    # Remove the container
    if ${CONTAINER_RUNTIME} ps --format '{{.Names}}' | grep -q "^${CONTAINER_NAME}$"; then
        echo "Stopping and removing the ioFog controller container..."
        ${CONTAINER_RUNTIME} stop ${CONTAINER_NAME}
        ${CONTAINER_RUNTIME} rm ${CONTAINER_NAME}
    fi

    # Remove config files
    echo "Checking if the ${CONTAINER_RUNTIME} volume exists..."

    if sudo ${CONTAINER_RUNTIME} volume inspect "${CONTROLLER_DB}" >/dev/null 2>&1; then
        echo "${CONTAINER_RUNTIME} volume '${CONTROLLER_DB}' found. Removing..."
        sudo ${CONTAINER_RUNTIME} volume rm "${CONTROLLER_DB}"
        echo "${CONTAINER_RUNTIME} volume '${CONTROLLER_DB}' has been removed."
    else
        echo "${CONTAINER_RUNTIME} volume '${CONTROLLER_DB}' does not exist. Skipping removal."
    fi

    # Remove log files
    echo "Removing log files..."
    if sudo ${CONTAINER_RUNTIME} volume inspect "${CONTROLLER_LOG_DIR}" >/dev/null 2>&1; then
        echo "${CONTAINER_RUNTIME} volume '${CONTROLLER_LOG_DIR}' found. Removing..."
        sudo ${CONTAINER_RUNTIME} volume rm "${CONTROLLER_LOG_DIR}"
        echo "${CONTAINER_RUNTIME} volume '${CONTROLLER_LOG_DIR}' has been removed."
    else
        echo "${CONTAINER_RUNTIME} volume '${CONTROLLER_LOG_DIR}' does not exist. Skipping removal."
    fi


    # Remove the executable script
    if [ -f ${EXECUTABLE_FILE} ]; then
        echo "Removing the iofog-controller executable script..."
        sudo rm -f ${EXECUTABLE_FILE}
    fi

    echo "ioFog controller uninstalled successfully!"
}

. /etc/iofog/controller/init.sh
init

do_uninstall_controller