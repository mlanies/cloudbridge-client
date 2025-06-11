#!/bin/bash

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Version
VERSION="1.0.0"

# Default values
CLIENT_VERSION="latest"
INSTALL_DIR="/usr/local/bin"
CONFIG_DIR="/etc/cloudbridge-client"
DATA_DIR="/var/lib/cloudbridge-client"
LOG_DIR="/var/log/cloudbridge-client"
SERVICE_DIR="/etc/systemd/system"
SERVICE_USER="cloudbridge-client"
REPO_URL="https://github.com/twogc/2GC-CloudBridge-Tunnel"

# Print usage
usage() {
    echo "CloudBridge Client Installer v${VERSION}"
    echo
    echo "Usage: $0 [options]"
    echo
    echo "Options:"
    echo "  -h, --help                 Show this help message"
    echo "  -v, --version              Show version"
    echo "  -d, --debug                Enable debug output"
    echo "  --client-version VERSION   Install specific version (default: latest)"
    echo "  --install-dir DIR          Installation directory (default: ${INSTALL_DIR})"
    echo "  --config-dir DIR           Configuration directory (default: ${CONFIG_DIR})"
    echo "  --data-dir DIR             Data directory (default: ${DATA_DIR})"
    echo "  --log-dir DIR              Log directory (default: ${LOG_DIR})"
    echo
    echo "Example:"
    echo "  $0 --client-version 1.0.0"
    echo
}

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -h|--help)
            usage
            exit 0
            ;;
        -v|--version)
            echo "CloudBridge Client Installer v${VERSION}"
            exit 0
            ;;
        -d|--debug)
            set -x
            ;;
        --client-version)
            CLIENT_VERSION="$2"
            shift
            ;;
        --install-dir)
            INSTALL_DIR="$2"
            shift
            ;;
        --config-dir)
            CONFIG_DIR="$2"
            shift
            ;;
        --data-dir)
            DATA_DIR="$2"
            shift
            ;;
        --log-dir)
            LOG_DIR="$2"
            shift
            ;;
        *)
            echo -e "${RED}Error: Unknown option $1${NC}"
            usage
            exit 1
            ;;
    esac
    shift
done

# Check if running as root
if [ "$EUID" -ne 0 ]; then
    echo -e "${RED}Error: Please run as root${NC}"
    exit 1
fi

# Function to print status
print_status() {
    echo -e "${GREEN}==>${NC} $1"
}

# Function to print error
print_error() {
    echo -e "${RED}Error:${NC} $1"
}

# Function to print warning
print_warning() {
    echo -e "${YELLOW}Warning:${NC} $1"
}

# Check system requirements
print_status "Checking system requirements..."

# Check if systemd is available
if ! command -v systemctl &> /dev/null; then
    print_error "systemd is not available. This installer requires systemd."
    exit 1
fi

# Check if curl is available
if ! command -v curl &> /dev/null; then
    print_error "curl is not available. Please install curl first."
    exit 1
fi

# Create necessary directories
print_status "Creating directories..."
mkdir -p "${CONFIG_DIR}"
mkdir -p "${DATA_DIR}"
mkdir -p "${LOG_DIR}"

# Create system user
print_status "Creating system user..."
if ! id "${SERVICE_USER}" &>/dev/null; then
    useradd -r -s /bin/false "${SERVICE_USER}" || {
        print_error "Failed to create system user"
        exit 1
    }
fi

# Download and install client
print_status "Downloading CloudBridge Client..."
if [ "${CLIENT_VERSION}" = "latest" ]; then
    LATEST_VERSION=$(curl -s "https://api.github.com/repos/twogc/2GC-CloudBridge-Tunnel/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
    CLIENT_VERSION="${LATEST_VERSION}"
fi

DOWNLOAD_URL="https://github.com/twogc/2GC-CloudBridge-Tunnel/releases/download/${CLIENT_VERSION}/cloudbridge-client"
curl -L -o "${INSTALL_DIR}/cloudbridge-client" "${DOWNLOAD_URL}" || {
    print_error "Failed to download client"
    exit 1
}

# Download systemd service file
print_status "Downloading systemd service file..."
curl -L -o "${SERVICE_DIR}/cloudbridge-client.service" "${REPO_URL}/raw/main/systemd/cloudbridge-client.service" || {
    print_error "Failed to download systemd service file"
    exit 1
}

# Download default configuration
print_status "Downloading default configuration..."
curl -L -o "${CONFIG_DIR}/config.yaml" "${REPO_URL}/raw/main/config/client-config.yaml" || {
    print_error "Failed to download configuration"
    exit 1
}

# Set permissions
print_status "Setting permissions..."
chown -R "${SERVICE_USER}:${SERVICE_USER}" "${DATA_DIR}" "${LOG_DIR}"
chmod 755 "${INSTALL_DIR}/cloudbridge-client"
chmod 644 "${SERVICE_DIR}/cloudbridge-client.service"
chmod 644 "${CONFIG_DIR}/config.yaml"

# Reload systemd
print_status "Reloading systemd..."
systemctl daemon-reload || {
    print_error "Failed to reload systemd"
    exit 1
}

# Enable and start service
print_status "Enabling and starting service..."
systemctl enable cloudbridge-client || {
    print_error "Failed to enable service"
    exit 1
}

systemctl start cloudbridge-client || {
    print_error "Failed to start service"
    exit 1
}

# Print success message
echo
echo -e "${GREEN}CloudBridge Client v${CLIENT_VERSION} installed successfully!${NC}"
echo
echo "Service status:"
systemctl status cloudbridge-client --no-pager
echo
echo "Useful commands:"
echo "  systemctl status cloudbridge-client    # Check service status"
echo "  journalctl -u cloudbridge-client -f    # View logs"
echo "  systemctl restart cloudbridge-client   # Restart service"
echo
echo "Configuration:"
echo "  Config file: ${CONFIG_DIR}/config.yaml"
echo "  Log file: ${LOG_DIR}/client.log"
echo "  Data directory: ${DATA_DIR}"
echo 