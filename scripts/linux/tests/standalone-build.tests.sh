#!/usr/bin/env bash

set -euo pipefail

TEST_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
LINUX_SCRIPTS_DIR="$(dirname "$TEST_DIR")"
PROJECT_ROOT="$(dirname "$(dirname "$LINUX_SCRIPTS_DIR")")"
DEPLOY_SCRIPT="$LINUX_SCRIPTS_DIR/deploy-linux.sh"
FRONTEND_ASSETS_PATH="$PROJECT_ROOT/frontend/dist/assets"
EXPECTED_API_BASE_URL="http://localhost:8080/api"

bash "$DEPLOY_SCRIPT" --frontend-mode standalone --build-only --skip-backend

if ! grep -R -F --include='*.js' -q "$EXPECTED_API_BASE_URL" "$FRONTEND_ASSETS_PATH"; then
    printf 'Standalone frontend bundle does not contain API base URL %s.\n' "$EXPECTED_API_BASE_URL" >&2
    exit 1
fi

printf 'Linux standalone frontend API base URL test passed.\n'
