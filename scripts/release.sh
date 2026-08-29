#!/usr/bin/env bash

set -euo pipefail

# Release cross-compilation automation script
VERSION="${1:-0.1.0}"
DIST_DIR="dist"
MAIN_PKG="./cmd/andtls"
COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ")

LDFLAGS="-s -w -X main.Version=${VERSION} -X main.Commit=${COMMIT} -X main.BuildDate=${BUILD_DATE}"

# Target architectures and OS list
TARGETS=(
    "linux/amd64"
    "linux/arm64"
    "darwin/amd64"
    "darwin/arm64"
    "windows/amd64"
)

echo "Building andtls v${VERSION} distribution packages..."
rm -rf "${DIST_DIR}"
mkdir -p "${DIST_DIR}"

for target in "${TARGETS[@]}"; do
    os="${target%%/*}"
    arch="${target##*/}"
    bin_name="andtls"

    if [ "${os}" == "windows" ]; then
        bin_name="andtls.exe"
    fi

    echo "Compiling for ${os}/${arch}..."
    out_dir="${DIST_DIR}/andtls_${VERSION}_${os}_${arch}"
    mkdir -p "${out_dir}"

    GOOS="${os}" GOARCH="${arch}" go build -ldflags "${LDFLAGS}" -o "${out_dir}/${bin_name}" "${MAIN_PKG}"
    cp README.md "${out_dir}/" 2>/dev/null || true

    # Package into archive
    if [ "${os}" == "windows" ]; then
        (cd "${DIST_DIR}" && zip -q -r "andtls_${VERSION}_${os}_${arch}.zip" "andtls_${VERSION}_${os}_${arch}")
    else
        tar -czf "${DIST_DIR}/andtls_${VERSION}_${os}_${arch}.tar.gz" -C "${DIST_DIR}" "andtls_${VERSION}_${os}_${arch}"
    fi

    rm -rf "${out_dir}"
done

# Generate SHA256 checksums
echo "Generating SHA256 checksums..."
(cd "${DIST_DIR}" && sha256sum * > "checksums.txt" 2>/dev/null || shasum -a 256 * > "checksums.txt")

echo "Release build complete. Artifacts saved in ${DIST_DIR}/"
ls -la "${DIST_DIR}"

