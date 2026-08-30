#!/usr/bin/env bash

set -euo pipefail

# Color codes for terminal output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

BINARY_NAME="andtls"
MAIN_PKG="./cmd/andtls"
VERSION="1.0.0"
BUILD_DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
COMMIT_HASH=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")

echo -e "${BLUE}======================================${NC}"
echo -e "${BLUE}    Building andtls v${VERSION}        ${NC}"
echo -e "${BLUE}======================================${NC}"

# Check for Go installation
if ! command -v go >/dev/null 2>&1; then
    echo -e "${RED}Error: Go compiler is not installed or not found in PATH${NC}" >&2
    exit 1
fi

GO_VERSION=$(go version | awk '{print $3}')
echo -e "${GREEN}✓${NC} Found Go version: ${GO_VERSION}"

# Check for ADB installation
if command -v adb >/dev/null 2>&1; then
    ADB_VERSION=$(adb version 2>&1 | head -n 1)
    echo -e "${GREEN}✓${NC} Found ADB: ${ADB_VERSION}"
else
    echo -e "${YELLOW}! Warning: ADB not found in PATH. Install android-tools/platform-tools to manage devices${NC}"
fi

# Build binary
echo -e "${BLUE}Compiling binary...${NC}"
LDFLAGS="-s -w -X main.Version=${VERSION} -X main.Commit=${COMMIT_HASH} -X main.BuildDate=${BUILD_DATE}"

go build -ldflags "${LDFLAGS}" -o "${BINARY_NAME}" "${MAIN_PKG}"

echo -e "${GREEN}✓ Successfully built: ./${BINARY_NAME}${NC}"
echo -e "Run with: ${GREEN}./${BINARY_NAME}${NC}"

