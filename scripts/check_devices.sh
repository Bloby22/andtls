#!/usr/bin/env bash

set -euo pipefail

# Quick terminal diagnostic tool for ADB environment
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

echo "Checking ADB environment and connected devices..."

if ! command -v adb >/dev/null 2>&1; then
    echo -e "${RED}ADB is not installed or not in PATH${NC}"
    exit 1
fi

echo -e "${GREEN}✓ ADB binary located at:${NC} $(which adb)"
adb version | head -n 2

echo ""
echo "Querying connected devices:"
adb devices -l

echo ""
echo "Active ADB server port status:"
if command -v ss >/dev/null 2>&1; then
    ss -tlpn | grep 5037 || echo "ADB port 5037 not actively listening"
elif command -v netstat >/dev/null 2>&1; then
    netstat -tlpn | grep 5037 || echo "ADB port 5037 not actively listening"
fi

