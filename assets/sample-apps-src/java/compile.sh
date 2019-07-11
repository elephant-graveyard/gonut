#!/bin/bash

set -euo pipefail

BASEDIR="$(cd "$(dirname "$0")" && pwd)"
TARGET_DIR="${BASEDIR}/../../sample-apps/$(basename "${BASEDIR}")"
mkdir -p "${TARGET_DIR}"

pushd "${BASEDIR}" >/dev/null
if ! javac -target 8 -source 8 App.java; then
  echo "An error occurred while compiling"
  exit 1
fi

trap 'rm ${BASEDIR}/App.class' EXIT

if ! jar cfen "${TARGET_DIR}/App.jar" App App.class; then
  echo "An error occurred while packaging"
  exit 1
fi

popd >/dev/null
