#!/usr/bin/env bash

# Builds and deploys standalone frontend/backend artifacts or a backend with
# embedded frontend assets on Linux.

set -Eeuo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SCRIPTS_DIR="$(dirname "$SCRIPT_DIR")"
PROJECT_ROOT="$(dirname "$SCRIPTS_DIR")"
FRONTEND_PATH="$PROJECT_ROOT/frontend"
BACKEND_PATH="$PROJECT_ROOT/backend"
FRONTEND_DIST_PATH="$FRONTEND_PATH/dist"
EMBEDDED_DIST_PATH="$BACKEND_PATH/internal/webui/dist"
BACKEND_BUILD_PATH="$BACKEND_PATH/shihai-server"

FRONTEND_MODE="standalone"
BUILD_ONLY=false
SKIP_FRONTEND=false
SKIP_BACKEND=false
DEPLOY_PATH="/opt/shihai"
STANDALONE_API_BASE_URL="http://localhost:8080/api"

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m'

ORIGINAL_DIRECTORY="$PWD"
HAD_API_BASE_URL=false
ORIGINAL_API_BASE_URL=""
EMBEDDED_FILES_PREPARED=false

if [[ -v VITE_API_BASE_URL ]]; then
    HAD_API_BASE_URL=true
    ORIGINAL_API_BASE_URL="$VITE_API_BASE_URL"
fi

usage() {
    cat <<'EOF'
Usage: deploy-linux.sh [options]

Options:
  --frontend-mode standalone|embedded  Frontend deployment form (default: standalone)
  --build-only                          Build without copying to the deployment path
  --skip-frontend                       Skip standalone frontend build and deployment
  --skip-backend                        Skip standalone backend build and deployment
  --deploy-path PATH                    Deployment directory (default: /opt/shihai)
  --standalone-api-base-url URL         API URL compiled into the standalone frontend
  -h, --help                            Show this help
EOF
}

require_value() {
    local option="$1"
    local value="${2:-}"
    if [[ -z "$value" ]]; then
        printf '%bMissing value for %s%b\n' "$RED" "$option" "$NC" >&2
        exit 1
    fi
}

command_exists() {
    command -v "$1" >/dev/null 2>&1
}

cleanup() {
    cd "$ORIGINAL_DIRECTORY"

    if [[ "$HAD_API_BASE_URL" == true ]]; then
        export VITE_API_BASE_URL="$ORIGINAL_API_BASE_URL"
    else
        unset VITE_API_BASE_URL || true
    fi

    if [[ "$EMBEDDED_FILES_PREPARED" == true && -d "$EMBEDDED_DIST_PATH" ]]; then
        rm -rf -- "$EMBEDDED_DIST_PATH"
    fi
}
trap cleanup EXIT

while [[ $# -gt 0 ]]; do
    case "$1" in
        --frontend-mode)
            require_value "$1" "${2:-}"
            FRONTEND_MODE="${2,,}"
            shift 2
            ;;
        --build-only)
            BUILD_ONLY=true
            shift
            ;;
        --skip-frontend)
            SKIP_FRONTEND=true
            shift
            ;;
        --skip-backend)
            SKIP_BACKEND=true
            shift
            ;;
        --deploy-path)
            require_value "$1" "${2:-}"
            DEPLOY_PATH="$2"
            shift 2
            ;;
        --standalone-api-base-url)
            require_value "$1" "${2:-}"
            STANDALONE_API_BASE_URL="$2"
            shift 2
            ;;
        -h|--help)
            usage
            exit 0
            ;;
        *)
            printf '%bUnknown option: %s%b\n' "$RED" "$1" "$NC" >&2
            usage >&2
            exit 1
            ;;
    esac
done

if [[ "$FRONTEND_MODE" != "standalone" && "$FRONTEND_MODE" != "embedded" ]]; then
    printf '%bInvalid frontend mode: %s%b\n' "$RED" "$FRONTEND_MODE" "$NC" >&2
    exit 1
fi
if [[ "$DEPLOY_PATH" != /* || "$DEPLOY_PATH" == "/" ]]; then
    printf '%bDeployment path must be an absolute path other than /.%b\n' "$RED" "$NC" >&2
    exit 1
fi
if [[ "$FRONTEND_MODE" == "embedded" && ("$SKIP_FRONTEND" == true || "$SKIP_BACKEND" == true) ]]; then
    printf '%bEmbedded mode requires both frontend and backend builds.%b\n' "$RED" "$NC" >&2
    exit 1
fi

printf '%b========================================%b\n' "$CYAN" "$NC"
printf '%bShihai Poetry Platform Deployment Tool%b\n' "$CYAN" "$NC"
printf '%b========================================%b\n' "$CYAN" "$NC"
printf 'Frontend mode: %s\n' "$FRONTEND_MODE"
printf 'Deployment path: %s\n\n' "$DEPLOY_PATH"

printf '%bChecking prerequisites...%b\n' "$YELLOW" "$NC"
if [[ "$SKIP_FRONTEND" == false ]]; then
    if ! command_exists node || ! command_exists npm; then
        printf '%bNode.js and npm must be installed and available in PATH.%b\n' "$RED" "$NC" >&2
        exit 1
    fi
    printf '  %b[OK]%b Node.js and npm found\n' "$GREEN" "$NC"
fi
if [[ "$SKIP_BACKEND" == false ]]; then
    if ! command_exists go; then
        printf '%bGo must be installed and available in PATH.%b\n' "$RED" "$NC" >&2
        exit 1
    fi
    printf '  %b[OK]%b Go found\n' "$GREEN" "$NC"
fi
printf '\n'

if [[ "$SKIP_FRONTEND" == false ]]; then
    printf '%bBuilding frontend...%b\n' "$YELLOW" "$NC"
    cd "$FRONTEND_PATH"

    if [[ "$FRONTEND_MODE" == "standalone" ]]; then
        export VITE_API_BASE_URL="$STANDALONE_API_BASE_URL"
    else
        export VITE_API_BASE_URL="/api"
    fi

    npm install
    npm run build
    printf '  %b[OK]%b Frontend built successfully\n\n' "$GREEN" "$NC"
fi

if [[ "$SKIP_BACKEND" == false ]]; then
    printf '%bBuilding backend...%b\n' "$YELLOW" "$NC"
    cd "$BACKEND_PATH"
    go mod tidy

    if [[ "$FRONTEND_MODE" == "embedded" ]]; then
        rm -rf -- "$EMBEDDED_DIST_PATH"
        mkdir -p "$EMBEDDED_DIST_PATH"
        EMBEDDED_FILES_PREPARED=true
        cp -a "$FRONTEND_DIST_PATH/." "$EMBEDDED_DIST_PATH/"
        go build -tags embed_frontend -o "$BACKEND_BUILD_PATH" ./cmd/server
    else
        go build -o "$BACKEND_BUILD_PATH" ./cmd/server
    fi
    printf '  %b[OK]%b Backend built successfully\n\n' "$GREEN" "$NC"
fi

if [[ "$BUILD_ONLY" == true ]]; then
    printf '%bBuild-only mode: deployment skipped.%b\n' "$YELLOW" "$NC"
    exit 0
fi

printf '%bDeploying to local machine...%b\n' "$YELLOW" "$NC"
if ! mkdir -p "$DEPLOY_PATH" 2>/dev/null || [[ ! -w "$DEPLOY_PATH" ]]; then
    if ! command_exists sudo; then
        printf '%bCannot create or write %s and sudo is unavailable.%b\n' "$RED" "$DEPLOY_PATH" "$NC" >&2
        exit 1
    fi
    sudo mkdir -p "$DEPLOY_PATH"
    sudo chown -R "$(id -u):$(id -g)" "$DEPLOY_PATH"
fi

DEPLOYED_FRONTEND_PATH="$DEPLOY_PATH/frontend"
if [[ "$FRONTEND_MODE" == "standalone" && "$SKIP_FRONTEND" == false ]]; then
    rm -rf -- "$DEPLOYED_FRONTEND_PATH"
    mkdir -p "$DEPLOYED_FRONTEND_PATH"
    cp -a "$FRONTEND_DIST_PATH/." "$DEPLOYED_FRONTEND_PATH/"
    printf '  %b[OK]%b Standalone frontend deployed\n' "$GREEN" "$NC"
elif [[ "$FRONTEND_MODE" == "embedded" ]]; then
    rm -rf -- "$DEPLOYED_FRONTEND_PATH"
fi

if [[ "$SKIP_BACKEND" == false ]]; then
    cp "$BACKEND_BUILD_PATH" "$DEPLOY_PATH/shihai-server"
    chmod +x "$DEPLOY_PATH/shihai-server"
    printf '  %b[OK]%b Backend deployed\n' "$GREEN" "$NC"

    if [[ -f "$BACKEND_PATH/config.json" ]]; then
        cp "$BACKEND_PATH/config.json" "$DEPLOY_PATH/config.json"
        printf '  %b[OK]%b Backend configuration deployed\n' "$GREEN" "$NC"
    fi
fi

if [[ "$SKIP_FRONTEND" == false && "$SKIP_BACKEND" == false ]]; then
    printf '%s\n' "$FRONTEND_MODE" > "$DEPLOY_PATH/frontend-mode.txt"
fi

printf '\n%bDeployment complete%b\n' "$GREEN" "$NC"
printf 'Location: %s\n' "$DEPLOY_PATH"
printf 'Frontend mode: %s\n' "$FRONTEND_MODE"
