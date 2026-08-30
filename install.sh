#!/usr/bin/env bash

set -euo pipefail

# Color codes for terminal output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

BINARY_NAME="andtls"

echo -e "${BLUE}======================================${NC}"
echo -e "${BLUE}    Installing andtls CLI/TUI         ${NC}"
echo -e "${BLUE}======================================${NC}"

# Build the binary first
./build.sh

# Determine install destination
TARGET_DIR="${INSTALL_DIR:-}"
if [ -z "${TARGET_DIR}" ]; then
    if [ -w "/usr/local/bin" ]; then
        TARGET_DIR="/usr/local/bin"
    else
        TARGET_DIR="${HOME}/.local/bin"
    fi
fi

mkdir -p "${TARGET_DIR}" 2>/dev/null || {
    echo -e "${YELLOW}! Unable to write to ${TARGET_DIR}, attempting sudo...${NC}"
    sudo mkdir -p "${TARGET_DIR}"
    sudo cp "${BINARY_NAME}" "${TARGET_DIR}/${BINARY_NAME}"
    sudo chmod +x "${TARGET_DIR}/${BINARY_NAME}"
    echo -e "${GREEN}✓ Successfully installed ${BINARY_NAME} to: ${TARGET_DIR}/${BINARY_NAME}${NC}"
    exit 0
}

cp "${BINARY_NAME}" "${TARGET_DIR}/${BINARY_NAME}"
chmod +x "${TARGET_DIR}/${BINARY_NAME}"

echo -e "${GREEN}✓ Successfully installed ${BINARY_NAME} to: ${TARGET_DIR}/${BINARY_NAME}${NC}"

# Check if TARGET_DIR is in PATH
if [[ ":$PATH:" != *":${TARGET_DIR}:"* ]]; then
    echo -e "${YELLOW}! Note: ${TARGET_DIR} is not in your current PATH${NC}"
    echo -e "Add this to your ~/.bashrc or ~/.zshrc:"
    echo -e "  export PATH=\"\$PATH:${TARGET_DIR}\""
else
    echo -e "You can now run '${GREEN}${BINARY_NAME}${NC}' from anywhere"
fi

