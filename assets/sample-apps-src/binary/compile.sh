#!/bin/bash

set -euo pipefail

BASEDIR="$(cd "$(dirname "$0")" && pwd)"
TARGET_DIR="${BASEDIR}/../../sample-apps/$(basename "${BASEDIR}")"
mkdir -p "${TARGET_DIR}"

pushd "${BASEDIR}" >/dev/null

GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
  -tags netgo \
  -ldflags="-s -w -extldflags '-static'" \
  -o "${TARGET_DIR}/binary" \
  "${BASEDIR}/main.go"

popd >/dev/null
