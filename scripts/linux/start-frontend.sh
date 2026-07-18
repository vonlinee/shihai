#!/usr/bin/env bash

# Starts the standalone deployed Shihai frontend with Vite Preview on Linux.

set -Eeuo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SCRIPTS_DIR="$(dirname "$SCRIPT_DIR")"
PROJECT_ROOT="$(dirname "$SCRIPTS_DIR")"
FRONTEND_PROJECT_PATH="$PROJECT_ROOT/frontend"
VITE_COMMAND="$FRONTEND_PROJECT_PATH/node_modules/.bin/vite"

DEPLOY_PATH="/opt/shihai"
PORT=80
DRY_RUN=false

usage() {
    cat <<'EOF'
Usage: start-frontend.sh [options]

Options:
  --deploy-path PATH  Deployment directory (default: /opt/shihai)
  --port PORT         Vite Preview port (default: 80)
  --dry-run           Print the command without starting Vite Preview
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
        --dry-run)
            DRY_RUN=true
            shift
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
    printf 'Port must be an integer between 1 and 65535.\n' >&2
    exit 1
fi

FRONTEND_PATH="$DEPLOY_PATH/frontend"
FRONTEND_INDEX="$FRONTEND_PATH/index.html"
if [[ ! -f "$FRONTEND_INDEX" ]]; then
    printf 'Frontend files not found at %s.\n' "$FRONTEND_PATH" >&2
    printf 'Run deploy-linux.sh in standalone mode first.\n' >&2
    exit 1
fi
if [[ ! -x "$VITE_COMMAND" ]]; then
    printf 'Vite not found at %s. Run npm install in the frontend directory first.\n' "$VITE_COMMAND" >&2
    exit 1
fi

VITE_ARGUMENTS=(
    preview
    --host 0.0.0.0
    --port "$PORT"
    --strictPort
    --outDir "$FRONTEND_PATH"
)

printf 'Starting Shihai standalone frontend...\n'
printf 'Frontend URL: http://localhost:%s\n' "$PORT"
printf 'Command: %s %s\n' "$VITE_COMMAND" "${VITE_ARGUMENTS[*]}"

if [[ "$DRY_RUN" == true ]]; then
    exit 0
fi

printf 'Press Ctrl+C to stop.\n\n'
cd "$FRONTEND_PROJECT_PATH"
exec "$VITE_COMMAND" "${VITE_ARGUMENTS[@]}"
