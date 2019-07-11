#!/bin/bash

set -euo pipefail

BASEDIR="$(cd "$(dirname "$0")" && pwd)"
TARGET_DIR="${BASEDIR}/../../sample-apps/$(basename "${BASEDIR}")"

if [[ -f "${TARGET_DIR}/App.jar" ]]; then
  rm "${TARGET_DIR}/App.jar"
fi
