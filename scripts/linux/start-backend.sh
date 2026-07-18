#!/usr/bin/env bash

# Starts the deployed Shihai backend on Linux.

set -Eeuo pipefail

DEPLOY_PATH="/opt/shihai"
PORT=8080

usage() {
    cat <<'EOF'
Usage: start-backend.sh [options]

Options:
  --deploy-path PATH  Deployment directory (default: /opt/shihai)
  --port PORT         Backend HTTP port (default: 8080)
  -h, --help          Show this help
EOF
}

while [[ $# -gt 0 ]]; do
    case "$1" in
        --deploy-path)
            if [[ -z "${2:-}" ]]; then
                printf 'Missing value for --deploy-path.\n' >&2
                exit 1
            fi
            DEPLOY_PATH="$2"
            shift 2
            ;;
        --port)
            if [[ -z "${2:-}" ]]; then
                printf 'Missing value for --port.\n' >&2
                exit 1
            fi
            PORT="$2"
            shift 2
            ;;
        -h|--help)
            usage
            exit 0
            ;;
        *)
            printf 'Unknown option: %s\n' "$1" >&2
            usage >&2
            exit 1
            ;;
    esac
done

if [[ ! "$PORT" =~ ^[0-9]+$ ]] || ((PORT < 1 || PORT > 65535)); then
    printf 'Backend port must be an integer between 1 and 65535.\n' >&2
    exit 1
fi

BACKEND_EXECUTABLE="$DEPLOY_PATH/shihai-server"
if [[ ! -x "$BACKEND_EXECUTABLE" ]]; then
    printf 'Backend executable not found or not executable: %s\n' "$BACKEND_EXECUTABLE" >&2
    printf 'Run deploy-linux.sh first.\n' >&2
    exit 1
fi

printf 'Starting Shihai backend...\n'
printf 'Backend URL: http://localhost:%s\n' "$PORT"
printf 'Press Ctrl+C to stop.\n\n'

cd "$DEPLOY_PATH"
exec "$BACKEND_EXECUTABLE" -port "$PORT"
