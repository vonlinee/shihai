#!/usr/bin/env bash

set -euo pipefail

TEST_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
LINUX_SCRIPTS_DIR="$(dirname "$TEST_DIR")"
PROJECT_ROOT="$(dirname "$(dirname "$LINUX_SCRIPTS_DIR")")"
START_SCRIPT="$LINUX_SCRIPTS_DIR/start-all.sh"
START_FRONTEND_SCRIPT="$LINUX_SCRIPTS_DIR/start-frontend.sh"

assert_contains() {
    local output="$1"
    local expected="$2"

    if [[ "$output" != *"$expected"* ]]; then
        printf 'Expected output to contain %q. Actual output:\n%s\n' "$expected" "$output" >&2
        exit 1
    fi
}

default_output="$(bash "$START_SCRIPT" --dry-run)"
assert_contains "$default_output" "Mode: deployment"
assert_contains "$default_output" "Frontend mode: embedded"
assert_contains "$default_output" "Redeploy: true"
assert_contains "$default_output" "Process count: 1"

development_output="$(bash "$START_SCRIPT" --mode development --dry-run)"
assert_contains "$development_output" "Mode: development"
assert_contains "$development_output" "Process count: 2"
assert_contains "$development_output" "npm run dev"
assert_contains "$development_output" "go run ./cmd/server"

standalone_output="$(bash "$START_SCRIPT" --mode deployment --frontend-mode standalone --dry-run)"
assert_contains "$standalone_output" "Frontend mode: standalone"
assert_contains "$standalone_output" "Process count: 2"

reuse_output="$(bash "$START_SCRIPT" --no-redeploy --dry-run)"
assert_contains "$reuse_output" "Redeploy: false"

detach_output="$(bash "$START_SCRIPT" --detach --dry-run)"
assert_contains "$detach_output" "Detach: true"
assert_contains "$detach_output" "PID file: /opt/shihai/run/start-all.pid"
assert_contains "$detach_output" "Log file: /opt/shihai/logs/start-all.log"

temporary_deploy_path="$PROJECT_ROOT/tmp/linux-start-frontend-test"
temporary_detach_path="$PROJECT_ROOT/tmp/linux-detach-test"
cleanup() {
    if [[ -f "$temporary_detach_path/run/start-all.pid" ]]; then
        managed_pid="$(cat "$temporary_detach_path/run/start-all.pid" 2>/dev/null || true)"
        if [[ "$managed_pid" =~ ^[0-9]+$ ]]; then
            kill "$managed_pid" 2>/dev/null || true
        fi
    fi
    rm -rf "$temporary_deploy_path"
    rm -rf "$temporary_detach_path"
}
trap cleanup EXIT

mkdir -p "$temporary_deploy_path/frontend"
printf '<html></html>\n' > "$temporary_deploy_path/frontend/index.html"

status_output="$(bash "$START_SCRIPT" --deploy-path "$temporary_deploy_path" --status)"
assert_contains "$status_output" "Status: stopped"

mkdir -p "$temporary_deploy_path/run"
printf '999999\n' > "$temporary_deploy_path/run/start-all.pid"
stop_output="$(bash "$START_SCRIPT" --deploy-path "$temporary_deploy_path" --stop)"
assert_contains "$stop_output" "No running Shihai process found"
if [[ -f "$temporary_deploy_path/run/start-all.pid" ]]; then
    printf 'Expected --stop to remove a stale PID file.\n' >&2
    exit 1
fi

frontend_output="$(bash "$START_FRONTEND_SCRIPT" --deploy-path "$temporary_deploy_path" --port 4173 --dry-run)"
assert_contains "$frontend_output" "vite preview"
assert_contains "$frontend_output" "--port 4173"
assert_contains "$frontend_output" "--outDir $temporary_deploy_path/frontend"

mkdir -p "$temporary_detach_path"
cat > "$temporary_detach_path/shihai-server" <<'EOF'
#!/usr/bin/env bash
trap 'exit 0' INT TERM
while true; do
    sleep 1
done
EOF
chmod +x "$temporary_detach_path/shihai-server"
printf 'embedded\n' > "$temporary_detach_path/frontend-mode.txt"

detached_start_output="$(bash "$START_SCRIPT" --deploy-path "$temporary_detach_path" --no-redeploy --detach)"
assert_contains "$detached_start_output" "Started Shihai in the background with PID"

detached_status_output="$(bash "$START_SCRIPT" --deploy-path "$temporary_detach_path" --status)"
assert_contains "$detached_status_output" "Status: running"

detached_stop_output="$(bash "$START_SCRIPT" --deploy-path "$temporary_detach_path" --stop)"
assert_contains "$detached_stop_output" "Stopped Shihai background process"

detached_final_status="$(bash "$START_SCRIPT" --deploy-path "$temporary_detach_path" --status)"
assert_contains "$detached_final_status" "Status: stopped"

printf 'Linux start script dry-run tests passed.\n'
