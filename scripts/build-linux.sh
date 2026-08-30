#!/usr/bin/env bash
set -e

# LocalValet Linux Production Build & Package Script

BASE_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
DIST_DIR="${BASE_DIR}/build/bin"
APP_NAME="localvalet"

echo "=========================================="
echo " Building LocalValet Production Release  "
echo "=========================================="

# 1. Setup Runtime Directories
"${BASE_DIR}/scripts/setup-runtime.sh"

# 2. Build Frontend
echo "==> [1/3] Building React Frontend..."
cd "${BASE_DIR}/frontend"
npm run build
cd "${BASE_DIR}"

# 3. Run Backend Test Suite
echo "==> [2/3] Running Go Test Suite..."
go test -v ./internal/...

# 4. Build Standalone Binary
echo "==> [3/3] Compiling Go Binary with optimizations..."
mkdir -p "${DIST_DIR}"
go build -ldflags="-s -w" -o "${DIST_DIR}/${APP_NAME}" .

echo "=========================================="
echo " Build Success! Artifact: ${DIST_DIR}/${APP_NAME} "
echo "=========================================="
