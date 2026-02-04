#!/bin/sh
set -x

# Check if running as root or via sudo 
if [ "$(id -u)" -ne 0 ]; then
	echo "Error: This script must be run as root or via sudo"
	echo "   Please run: sudo \$0 \$*"
	exit 1
fi

# Check can sudo without password (when re-executed or for subcommands)
if ! $(sudo ls /tmp/ > /dev/null); then
	MSG="Unable to successfully use sudo with user $USER on this host.\nUser $USER must be in sudoers group and using sudo without password must be enabled.\nPlease see docs.datasance.com documentation for more details."
	echo $MSG
	exit 1
fi
