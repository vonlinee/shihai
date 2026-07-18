#!/usr/bin/env bash

# Starts Shihai in development, standalone deployment, or embedded deployment
# mode while managing all child processes in the current terminal.

set -Eeuo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SCRIPTS_DIR="$(dirname "$SCRIPT_DIR")"
PROJECT_ROOT="$(dirname "$SCRIPTS_DIR")"
FRONTEND_PATH="$PROJECT_ROOT/frontend"
BACKEND_PATH="$PROJECT_ROOT/backend"
DEPLOY_SCRIPT="$SCRIPT_DIR/deploy-linux.sh"
START_FRONTEND_SCRIPT="$SCRIPT_DIR/start-frontend.sh"
START_BACKEND_SCRIPT="$SCRIPT_DIR/start-backend.sh"
ORIGINAL_ARGUMENTS=("$@")

MODE="deployment"
FRONTEND_MODE="embedded"
REDEPLOY=true
DEPLOY_PATH="/opt/shihai"
STANDALONE_API_BASE_URL=""
BACKEND_PORT=8080
FRONTEND_PORT=80
DRY_RUN=false
DETACH=false
ACTION="start"
ACTION_COUNT=0
BACKGROUND_CHILD="${SHIHAI_BACKGROUND_CHILD:-false}"
CHILD_PIDS=()

usage() {
    cat <<'EOF'
Usage: start-all.sh [options]

Options:
  --mode development|deployment        Startup mode (default: deployment)
  --frontend-mode standalone|embedded  Deployment frontend mode (default: embedded)
  --redeploy                           Rebuild before deployment startup (default)
  --no-redeploy                        Reuse existing deployment artifacts
  --deploy-path PATH                   Deployment directory (default: /opt/shihai)
  --standalone-api-base-url URL        API URL compiled into the standalone frontend
  --backend-port PORT                  Backend HTTP port (default: 8080)
  --frontend-port PORT                 Standalone frontend port (default: 80)
  --detach                             Start services in the background
  --stop                               Stop the managed background instance
  --status                             Show background instance status
  --dry-run                            Print the startup plan without executing it
  -h, --help                           Show this help
EOF
}

require_value() {
    local option="$1"
    local value="${2:-}"
    if [[ -z "$value" ]]; then
        printf 'Missing value for %s.\n' "$option" >&2
        exit 1
    fi
}

command_exists() {
    command -v "$1" >/dev/null 2>&1
}

read_managed_pid() {
    local pid=""
    if [[ ! -f "$PID_FILE" ]]; then
        return 1
    fi
    IFS= read -r pid < "$PID_FILE" || true
    if [[ ! "$pid" =~ ^[0-9]+$ ]]; then
        return 1
    fi
    printf '%s\n' "$pid"
}

managed_process_is_running() {
    local pid="${1:-}"
    [[ "$pid" =~ ^[0-9]+$ ]] && kill -0 "$pid" 2>/dev/null
}

cleanup_children() {
    local pid
    for pid in "${CHILD_PIDS[@]}"; do
        if kill -0 "$pid" 2>/dev/null; then
            kill "$pid" 2>/dev/null || true
        fi
    done
    for pid in "${CHILD_PIDS[@]}"; do
        wait "$pid" 2>/dev/null || true
    done
    CHILD_PIDS=()
}

handle_signal() {
    printf '\nStopping services...\n'
    cleanup_children
    exit 130
}
trap handle_signal INT TERM

cleanup_pid_file() {
    local pid=""
    if [[ "$BACKGROUND_CHILD" != true ]]; then
        return
    fi
    if pid="$(read_managed_pid 2>/dev/null)" && [[ "$pid" == "$$" ]]; then
        rm -f -- "$PID_FILE"
    fi
}

while [[ $# -gt 0 ]]; do
    case "$1" in
        --mode)
            require_value "$1" "${2:-}"
            MODE="${2,,}"
            shift 2
            ;;
        --frontend-mode)
            require_value "$1" "${2:-}"
            FRONTEND_MODE="${2,,}"
            shift 2
            ;;
        --redeploy)
            REDEPLOY=true
            shift
            ;;
        --no-redeploy)
            REDEPLOY=false
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
        --backend-port)
            require_value "$1" "${2:-}"
            BACKEND_PORT="$2"
            shift 2
            ;;
        --frontend-port)
            require_value "$1" "${2:-}"
            FRONTEND_PORT="$2"
            shift 2
            ;;
        --detach)
            DETACH=true
            shift
            ;;
        --stop)
            ACTION="stop"
            ((ACTION_COUNT += 1))
            shift
            ;;
        --status)
            ACTION="status"
            ((ACTION_COUNT += 1))
            shift
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

if [[ "$MODE" != "development" && "$MODE" != "deployment" ]]; then
    printf 'Invalid mode: %s\n' "$MODE" >&2
    exit 1
fi
if [[ "$FRONTEND_MODE" != "standalone" && "$FRONTEND_MODE" != "embedded" ]]; then
    printf 'Invalid frontend mode: %s\n' "$FRONTEND_MODE" >&2
    exit 1
fi
if [[ "$DEPLOY_PATH" != /* || "$DEPLOY_PATH" == "/" ]]; then
    printf 'Deployment path must be an absolute path other than /.\n' >&2
    exit 1
fi
if [[ ! "$FRONTEND_PORT" =~ ^[0-9]+$ ]] || ((FRONTEND_PORT < 1 || FRONTEND_PORT > 65535)); then
    printf 'Frontend port must be an integer between 1 and 65535.\n' >&2
    exit 1
fi
if [[ ! "$BACKEND_PORT" =~ ^[0-9]+$ ]] || ((BACKEND_PORT < 1 || BACKEND_PORT > 65535)); then
    printf 'Backend port must be an integer between 1 and 65535.\n' >&2
    exit 1
fi
if [[ -z "$STANDALONE_API_BASE_URL" ]]; then
    STANDALONE_API_BASE_URL="http://localhost:$BACKEND_PORT/api"
fi
if ((ACTION_COUNT > 1)); then
    printf 'Use only one of --stop or --status.\n' >&2
    exit 1
fi
if [[ "$DETACH" == true && "$ACTION" != "start" ]]; then
    printf '%s cannot be combined with --detach.\n' "--$ACTION" >&2
    exit 1
fi

if [[ "$MODE" == "development" ]]; then
    RUNTIME_DIRECTORY="$PROJECT_ROOT/tmp/linux-runtime"
    LOG_DIRECTORY="$RUNTIME_DIRECTORY/logs"
else
    RUNTIME_DIRECTORY="$DEPLOY_PATH/run"
    LOG_DIRECTORY="$DEPLOY_PATH/logs"
fi
PID_FILE="$RUNTIME_DIRECTORY/start-all.pid"
LOG_FILE="$LOG_DIRECTORY/start-all.log"
trap cleanup_pid_file EXIT

if [[ "$ACTION" == "status" ]]; then
    MANAGED_PID=""
    if MANAGED_PID="$(read_managed_pid 2>/dev/null)" && managed_process_is_running "$MANAGED_PID"; then
        printf 'Status: running (PID %s)\n' "$MANAGED_PID"
    else
        rm -f -- "$PID_FILE"
        printf 'Status: stopped\n'
    fi
    exit 0
fi

if [[ "$ACTION" == "stop" ]]; then
    MANAGED_PID=""
    if ! MANAGED_PID="$(read_managed_pid 2>/dev/null)" || ! managed_process_is_running "$MANAGED_PID"; then
        rm -f -- "$PID_FILE"
        printf 'No running Shihai process found.\n'
        exit 0
    fi

    kill "$MANAGED_PID"
    for _ in {1..50}; do
        if ! managed_process_is_running "$MANAGED_PID"; then
            break
        fi
        sleep 0.1
    done
    if managed_process_is_running "$MANAGED_PID"; then
        printf 'Timed out while stopping Shihai process %s.\n' "$MANAGED_PID" >&2
        exit 1
    fi
    rm -f -- "$PID_FILE"
    printf 'Stopped Shihai background process %s.\n' "$MANAGED_PID"
    exit 0
fi

printf 'Mode: %s\n' "$MODE"
printf 'Backend port: %s\n' "$BACKEND_PORT"
printf 'Detach: %s\n' "$DETACH"
if [[ "$DETACH" == true ]]; then
    printf 'PID file: %s\n' "$PID_FILE"
    printf 'Log file: %s\n' "$LOG_FILE"
fi
if [[ "$MODE" == "deployment" ]]; then
    printf 'Frontend mode: %s\n' "$FRONTEND_MODE"
    printf 'Redeploy: %s\n' "$REDEPLOY"
    printf 'Deploy path: %s\n' "$DEPLOY_PATH"
    if [[ "$FRONTEND_MODE" == "standalone" ]]; then
        printf 'Standalone API base URL: %s\n' "$STANDALONE_API_BASE_URL"
    fi
fi

if [[ "$MODE" == "development" ]]; then
    printf 'Process count: 2\n'
    printf 'Process: frontend\n  Directory: %s\n  Command: VITE_BACKEND_URL=http://localhost:%s npm run dev\n' \
        "$FRONTEND_PATH" "$BACKEND_PORT"
    printf 'Process: backend\n  Directory: %s\n  Command: go run ./cmd/server -port %s\n' \
        "$BACKEND_PATH" "$BACKEND_PORT"
elif [[ "$FRONTEND_MODE" == "embedded" ]]; then
    printf 'Process count: 1\n'
    printf 'Process: backend-with-embedded-frontend\n  Directory: %s\n  Command: %s --deploy-path %s --port %s\n' \
        "$DEPLOY_PATH" "$START_BACKEND_SCRIPT" "$DEPLOY_PATH" "$BACKEND_PORT"
else
    printf 'Process count: 2\n'
    printf 'Process: frontend\n  Directory: %s\n  Command: %s --deploy-path %s --port %s\n' \
        "$DEPLOY_PATH" "$START_FRONTEND_SCRIPT" "$DEPLOY_PATH" "$FRONTEND_PORT"
    printf 'Process: backend\n  Directory: %s\n  Command: %s --deploy-path %s --port %s\n' \
        "$DEPLOY_PATH" "$START_BACKEND_SCRIPT" "$DEPLOY_PATH" "$BACKEND_PORT"
fi

if [[ "$DRY_RUN" == true ]]; then
    exit 0
fi

if [[ "$DETACH" == true && "$BACKGROUND_CHILD" != true ]]; then
    EXISTING_PID=""
    if EXISTING_PID="$(read_managed_pid 2>/dev/null)" && managed_process_is_running "$EXISTING_PID"; then
        printf 'Shihai is already running in the background with PID %s.\n' "$EXISTING_PID" >&2
        exit 1
    fi
    rm -f -- "$PID_FILE"
fi

if [[ "$BACKGROUND_CHILD" != true ]]; then
    if [[ "$MODE" == "development" ]]; then
        for command_name in node npm go; do
            if ! command_exists "$command_name"; then
                printf '%s must be installed and available in PATH.\n' "$command_name" >&2
                exit 1
            fi
        done
    else
        if [[ "$REDEPLOY" == true ]]; then
            bash "$DEPLOY_SCRIPT" \
                --frontend-mode "$FRONTEND_MODE" \
                --deploy-path "$DEPLOY_PATH" \
                --standalone-api-base-url "$STANDALONE_API_BASE_URL"
        else
            BACKEND_EXECUTABLE="$DEPLOY_PATH/shihai-server"
            MODE_FILE="$DEPLOY_PATH/frontend-mode.txt"
            if [[ ! -x "$BACKEND_EXECUTABLE" ]]; then
                printf 'Backend deployment not found at %s.\n' "$BACKEND_EXECUTABLE" >&2
                exit 1
            fi
            if [[ ! -f "$MODE_FILE" ]]; then
                printf 'Deployment mode marker not found at %s. Redeploy first.\n' "$MODE_FILE" >&2
                exit 1
            fi
            DEPLOYED_MODE="$(tr -d '[:space:]' < "$MODE_FILE")"
            if [[ "$DEPLOYED_MODE" != "$FRONTEND_MODE" ]]; then
                printf "Existing deployment uses frontend mode '%s', but '%s' was requested.\n" \
                    "$DEPLOYED_MODE" "$FRONTEND_MODE" >&2
                exit 1
            fi
            if [[ "$FRONTEND_MODE" == "standalone" && ! -f "$DEPLOY_PATH/frontend/index.html" ]]; then
                printf 'Standalone frontend deployment not found at %s/frontend/index.html.\n' "$DEPLOY_PATH" >&2
                exit 1
            fi
        fi
    fi
fi

if [[ "$DETACH" == true && "$BACKGROUND_CHILD" != true ]]; then
    if ! command_exists nohup; then
        printf 'nohup must be installed and available in PATH.\n' >&2
        exit 1
    fi
    mkdir -p "$RUNTIME_DIRECTORY" "$LOG_DIRECTORY"
    SHIHAI_BACKGROUND_CHILD=true nohup bash "$SCRIPT_DIR/start-all.sh" "${ORIGINAL_ARGUMENTS[@]}" \
        >> "$LOG_FILE" 2>&1 < /dev/null &
    MANAGED_PID=$!
    printf '%s\n' "$MANAGED_PID" > "$PID_FILE"
    sleep 1
    if ! managed_process_is_running "$MANAGED_PID"; then
        rm -f -- "$PID_FILE"
        printf 'Background startup failed. Check %s.\n' "$LOG_FILE" >&2
        exit 1
    fi
    printf 'Started Shihai in the background with PID %s.\n' "$MANAGED_PID"
    printf 'Log file: %s\n' "$LOG_FILE"
    exit 0
fi

if [[ "$MODE" == "development" ]]; then
    (
        cd "$FRONTEND_PATH"
        exec env VITE_BACKEND_URL="http://localhost:$BACKEND_PORT" npm run dev
    ) &
    CHILD_PIDS+=("$!")

    (
        cd "$BACKEND_PATH"
        exec go run ./cmd/server -port "$BACKEND_PORT"
    ) &
    CHILD_PIDS+=("$!")
elif [[ "$FRONTEND_MODE" == "embedded" ]]; then
    bash "$START_BACKEND_SCRIPT" --deploy-path "$DEPLOY_PATH" --port "$BACKEND_PORT" &
    CHILD_PIDS+=("$!")
else
    bash "$START_FRONTEND_SCRIPT" --deploy-path "$DEPLOY_PATH" --port "$FRONTEND_PORT" &
    CHILD_PIDS+=("$!")
    bash "$START_BACKEND_SCRIPT" --deploy-path "$DEPLOY_PATH" --port "$BACKEND_PORT" &
    CHILD_PIDS+=("$!")
fi

printf 'Services started. Press Ctrl+C to stop all services.\n'
set +e
wait -n
EXIT_STATUS=$?
set -e
cleanup_children
exit "$EXIT_STATUS"
