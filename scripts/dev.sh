#!/usr/bin/env bash

set -euo pipefail

# Color formatting variables
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

echo -e "${BLUE}======================================${NC}"
echo -e "${BLUE}   andtls Development Environment     ${NC}"
echo -e "${BLUE}======================================${NC}"

# Check for Go installation
if ! command -v go >/dev/null 2>&1; then
    echo -e "${RED}Error: Go is not installed in PATH${NC}" >&2
    exit 1
fi

# Run formatting and vet analysis
echo -e "${BLUE}Formatting and verifying source code...${NC}"
go fmt ./...
go vet ./...

echo -e "${GREEN}✓ Code check passed${NC}"

# Start development run
echo -e "${BLUE}Starting andtls in development mode...${NC}"
exec go run ./cmd/andtls

