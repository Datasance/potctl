#!/bin/sh
set -x
set -e

# INSTALL_DIR="/opt/iofog"
TMP_DIR="/tmp/iofog"
ETC_DIR="/etc/iofog/controller"
CONTROLLER_LOG_FOLDER=/var/log/iofog-controller
CONTROLLER_CONTAINER_NAME="iofog-controller"


command_exists() {
    command -v "$1" >/dev/null 2>&1
}

do_stop_iofog_controller() {
	if command_exists iofog-controller; then
		sudo service iofog-controller stop
	fi
}

do_install_iofog_controller() {

	echo "# Installing ioFog controller..."


    # 1. Ensure folders exist
    for FOLDER in ${ETC_DIR} ${CONTROLLER_LOG_FOLDER}; do
        if [ ! -d "$FOLDER" ]; then
            echo "Creating folder: $FOLDER"
            sudo mkdir -p "$FOLDER"
            sudo chmod 775 "$FOLDER"
        fi
    done

    # 4. Create systemd service file
    echo "Creating systemd service for ioFog controller..."
    SYSTEMD_SERVICE_FILE=/etc/systemd/system/iofog-controller.service

    cat <<EOF | sudo tee ${SYSTEMD_SERVICE_FILE} > /dev/null
[Unit]
Description=Datasance PoT IoFog Controller Service
After=docker.service
Requires=docker.service

[Service]
TimeoutStartSec=0
Restart=always
ExecStartPre=/usr/bin/docker pull ${controller_image} 
ExecStartPre=-/usr/bin/docker rm -f ${CONTROLLER_CONTAINER_NAME}
ExecStart=/usr/bin/docker run --rm --name ${CONTROLLER_CONTAINER_NAME} \\
-e IOFOG_CONTROLLER_IMAGE=${controller_image} \\
--env-file ${ETC_DIR}/iofog-controller.env \\
-v iofog-controller-db:/home/runner/.npm-global/lib/node_modules/@datasance/iofogcontroller/src/data/sqlite_files/:rw \\
-v iofog-controller-log:/var/log/iofog-controller:rw \\
-p 51121:51121 \\
-p 8008:8008 \\
--stop-timeout 60 \\
--attach stdout \\
--attach stderr \\
${controller_image}
ExecStop=/usr/bin/docker stop ${CONTROLLER_CONTAINER_NAME}

[Install]
WantedBy=default.target
EOF

    # Reload systemd and enable the service
    sudo systemctl daemon-reload
    sudo systemctl enable iofog-controller.service
	sudo systemctl start iofog-controller.service

    # 5. Create the /usr/local/bin/iofog-controller script
    echo "Creating iofog-controller executable script..."
    EXECUTABLE_FILE=/usr/local/bin/iofog-controller

    cat <<'EOF' | sudo tee ${EXECUTABLE_FILE} > /dev/null
#!/bin/bash
CONTAINER_NAME="iofog-controller"

# Check if the container is running
if ! docker ps --format '{{.Names}}' | grep -q "^${CONTAINER_NAME}$"; then
    echo "Error: The iofog-controller container is not running."
    exit 1
fi

# Execute the command in the container
docker exec ${CONTAINER_NAME} iofog-controller "$@"
EOF

    # Make the script executable
    sudo chmod +x ${EXECUTABLE_FILE}

    echo "ioFog controller installation completed!"

}


# main
controller_image="$1"

. /etc/iofog/controller/init.sh
init
do_stop_iofog_controller
do_install_iofog_controller
