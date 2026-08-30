#!/usr/bin/env bash
set -e

# LocalValet Runtime Environment Setup Script

BASE_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
echo "==> Setting up LocalValet runtime environment in ${BASE_DIR}"

mkdir -p "${BASE_DIR}/runtime/certs"
mkdir -p "${BASE_DIR}/runtime/logs"
mkdir -p "${BASE_DIR}/runtime/pids"
mkdir -p "${BASE_DIR}/runtime/linux/nginx/vhosts"
mkdir -p "${BASE_DIR}/config"

# Set permissions
chmod 755 "${BASE_DIR}/runtime/certs"
chmod 755 "${BASE_DIR}/runtime/logs"
chmod 755 "${BASE_DIR}/runtime/pids"
chmod 755 "${BASE_DIR}/runtime/linux/nginx/vhosts"

echo "==> Runtime directory structure created successfully."
