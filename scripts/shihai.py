#!/usr/bin/env python3
"""Cross-platform command line entry point for the Shihai project.

Only Python's standard library is used by this launcher. Project build tools
(Node.js/npm and Go) are invoked directly when a build or service requires
them.
"""

from __future__ import annotations

import os
import shlex
import shutil
import signal
import subprocess
import sys
import time
from pathlib import Path
from typing import Dict, Iterable, List, Optional, Sequence, Tuple


SCRIPT_PATH = Path(__file__).resolve()
SCRIPTS_DIR = SCRIPT_PATH.parent
PROJECT_ROOT = SCRIPTS_DIR.parent
FRONTEND_PATH = PROJECT_ROOT / "frontend"
BACKEND_PATH = PROJECT_ROOT / "backend"
IS_WINDOWS = os.name == "nt"
DEFAULT_DEPLOY_PATH = Path(r"C:\shihai-deploy") if IS_WINDOWS else Path("/opt/shihai")
BACKEND_BINARY_NAME = "shihai-server.exe" if IS_WINDOWS else "shihai-server"
NPM_COMMAND = "npm.cmd" if IS_WINDOWS else "npm"
VITE_COMMAND = FRONTEND_PATH / "node_modules" / ".bin" / ("vite.cmd" if IS_WINDOWS else "vite")
BACKGROUND_FLAG = "--_background-child"


HELP = """Usage:
  python scripts/shihai.py <command> [options]

Commands:
  start       Build (unless disabled) and start development/deployment services
  deploy      Build and deploy frontend/backend artifacts without starting them
  backend     Start the deployed backend in the current terminal
  frontend    Start the standalone deployed frontend in the current terminal
  start-development-frontend
              Start the frontend development server in the current terminal
  test        Run project checks (backend tests and/or frontend build)

Common start options:
  --mode development|deployment       Startup mode (default: deployment)
  --frontend-mode standalone|embedded Deployment frontend form (default: embedded)
  --redeploy                           Rebuild before deployment startup (default)
  --no-redeploy                        Reuse existing deployment artifacts
  --deploy-path PATH                   Deployment directory (platform default)
  --standalone-api-base-url URL        API URL compiled into standalone frontend
  --backend-port PORT                  Backend HTTP port (default: 8080)
  --frontend-port PORT                 Standalone frontend port (default: 80)
  --detach                             Start services in the background
  --status                             Show background instance status
  --stop                               Stop the managed background instance
  --dry-run                            Print the startup plan without executing it

Deploy options:
  --frontend-mode standalone|embedded Deployment frontend form (default: standalone)
  --build-only                        Build without copying to deployment path
  --skip-frontend                     Skip frontend build/deployment
  --skip-backend                      Skip backend build/deployment

Backend/frontend options:
  --deploy-path PATH                  Deployment directory
  --port PORT                         Listening port (default: 8080/80)

Test options:
  --target all|backend|frontend        Check target (default: all)

Examples:
  python scripts/shihai.py start --mode development
  python scripts/shihai.py start --frontend-mode standalone --frontend-port 4173
  python scripts/shihai.py deploy --frontend-mode embedded --build-only
  python scripts/shihai.py start --detach
  python scripts/shihai.py start --status
  python scripts/shihai.py start --stop
"""


class ScriptError(RuntimeError):
    """Expected command-line or process failure."""


def fail(message: str) -> None:
    raise ScriptError(message)


def parse_args(argv: Sequence[str]) -> Tuple[Dict[str, object], List[str]]:
    options: Dict[str, object] = {}
    positional: List[str] = []
    boolean_options = {
        "redeploy", "no-redeploy", "build-only", "skip-frontend", "skip-backend",
        "detach", "status", "stop", "dry-run", "_background-child", "help",
    }
    index = 0
    while index < len(argv):
        token = argv[index]
        if token == "-h":
            options["help"] = True
            index += 1
            continue
        if not token.startswith("--"):
            positional.append(token)
            index += 1
            continue

        raw = token[2:]
        if "=" in raw:
            name, value = raw.split("=", 1)
        else:
            name, value = raw, None
        if name in boolean_options:
            if value is None:
                value = True
            else:
                value = value.lower() == "true"
        else:
            if value is None:
                if index + 1 >= len(argv) or argv[index + 1].startswith("--"):
                    fail(f"Missing value for --{name}.")
                value = argv[index + 1]
                index += 1
        options[name] = value
        index += 1
    return options, positional


def option(options: Dict[str, object], name: str, fallback: object) -> object:
    return options[name] if name in options else fallback


def require_port(value: object, name: str) -> int:
    try:
        port = int(str(value))
    except ValueError:
        fail(f"{name} must be an integer between 1 and 65535.")
    if port < 1 or port > 65535:
        fail(f"{name} must be an integer between 1 and 65535.")
    return port


def require_absolute_deploy_path(value: object) -> Path:
    path = Path(str(value)).expanduser()
    if not path.is_absolute() or path == Path(path.anchor):
        fail(f"Deployment path must be an absolute path other than the filesystem root: {path}")
    return path.resolve()


def command_exists(command: str) -> bool:
    return shutil.which(command) is not None


def assert_commands(commands: Iterable[str]) -> None:
    missing = [command for command in commands if not command_exists(command)]
    if missing:
        fail(f"{', '.join(missing)} must be installed and available in PATH.")


def display_command(command: object, args: Sequence[object]) -> str:
    parts = [str(command), *(str(arg) for arg in args)]
    if IS_WINDOWS:
        return subprocess.list2cmdline(parts)
    return " ".join(shlex.quote(part) for part in parts)


def popen_command(command: object, args: Sequence[object], **kwargs: object) -> subprocess.Popen:
    """Start a process, handling Windows .cmd launchers consistently."""
    parts = [str(command), *(str(arg) for arg in args)]
    if IS_WINDOWS and str(command).lower().endswith((".cmd", ".bat")):
        kwargs["shell"] = True
        return subprocess.Popen(subprocess.list2cmdline(parts), **kwargs)
    return subprocess.Popen(parts, **kwargs)


def run(command: object, args: Sequence[object], cwd: Path, env: Optional[Dict[str, str]] = None) -> None:
    try:
        completed = popen_command(command, args, cwd=str(cwd), env=env, stdin=None)
        code = completed.wait()
    except OSError as error:
        fail(f"Unable to start {display_command(command, args)}: {error}")
    if code != 0:
        fail(f"{display_command(command, args)} exited with code {code}.")


def run_sync(command: object, args: Sequence[object], cwd: Path, env: Optional[Dict[str, str]] = None) -> None:
    try:
        completed = popen_command(command, args, cwd=str(cwd), env=env)
        code = completed.wait()
    except OSError as error:
        fail(f"Unable to start {display_command(command, args)}: {error}")
    if code != 0:
        fail(f"{display_command(command, args)} exited with code {code}.")


def copy_directory(source: Path, destination: Path) -> None:
    destination.mkdir(parents=True, exist_ok=True)
    for item in source.iterdir():
        target = destination / item.name
        if item.is_dir():
            shutil.copytree(item, target, dirs_exist_ok=True)
        else:
            shutil.copy2(item, target)


def deploy(options: Dict[str, object]) -> None:
    frontend_mode = str(option(options, "frontend-mode", "standalone")).lower()
    build_only = bool(option(options, "build-only", False))
    skip_frontend = bool(option(options, "skip-frontend", False))
    skip_backend = bool(option(options, "skip-backend", False))
    deploy_path = require_absolute_deploy_path(option(options, "deploy-path", DEFAULT_DEPLOY_PATH))
    standalone_api_base_url = str(option(options, "standalone-api-base-url", "http://localhost:8080/api"))
    frontend_dist_path = FRONTEND_PATH / "dist"
    embedded_dist_path = BACKEND_PATH / "internal" / "webui" / "dist"
    backend_build_path = BACKEND_PATH / BACKEND_BINARY_NAME
    original_api_base_url = os.environ.get("VITE_API_BASE_URL")
    had_api_base_url = "VITE_API_BASE_URL" in os.environ
    embedded_files_prepared = False

    if frontend_mode not in {"standalone", "embedded"}:
        fail(f"Invalid frontend mode: {frontend_mode}")
    if frontend_mode == "embedded" and (skip_frontend or skip_backend):
        fail("Embedded mode requires both frontend and backend builds; do not use --skip-frontend or --skip-backend.")

    print("========================================")
    print("Shihai Poetry Platform Deployment Tool")
    print("========================================")
    print(f"Frontend mode: {frontend_mode}")
    print(f"Deployment path: {deploy_path}\n")

    try:
        print("Checking prerequisites...")
        if not skip_frontend:
            assert_commands(["node", "npm.cmd" if IS_WINDOWS else "npm"])
            print("  [OK] Node.js and npm found")
        if not skip_backend:
            assert_commands(["go"])
            print("  [OK] Go found")
        print()

        if not skip_frontend:
            print("Building frontend...")
            env = os.environ.copy()
            env["VITE_API_BASE_URL"] = standalone_api_base_url if frontend_mode == "standalone" else "/api"
            run_sync(NPM_COMMAND, ["install"], FRONTEND_PATH, env)
            run_sync(NPM_COMMAND, ["run", "build"], FRONTEND_PATH, env)
            print("  [OK] Frontend built successfully\n")

        if not skip_backend:
            print("Building backend...")
            run_sync("go", ["mod", "tidy"], BACKEND_PATH)
            if frontend_mode == "embedded":
                if embedded_dist_path.exists():
                    shutil.rmtree(embedded_dist_path)
                copy_directory(frontend_dist_path, embedded_dist_path)
                embedded_files_prepared = True
                run_sync("go", ["build", "-tags", "embed_frontend", "-o", backend_build_path, "./cmd/server"], BACKEND_PATH)
            else:
                run_sync("go", ["build", "-o", backend_build_path, "./cmd/server"], BACKEND_PATH)
            print("  [OK] Backend built successfully\n")

        if build_only:
            print("Build-only mode: deployment skipped.")
            return

        print("Deploying to local machine...")
        deploy_path.mkdir(parents=True, exist_ok=True)
        deployed_frontend_path = deploy_path / "frontend"
        if frontend_mode == "standalone" and not skip_frontend:
            if deployed_frontend_path.exists():
                shutil.rmtree(deployed_frontend_path)
            copy_directory(frontend_dist_path, deployed_frontend_path)
            print("  [OK] Standalone frontend deployed")
        elif frontend_mode == "embedded" and deployed_frontend_path.exists():
            shutil.rmtree(deployed_frontend_path)

        if not skip_backend:
            deployed_backend_path = deploy_path / BACKEND_BINARY_NAME
            shutil.copy2(backend_build_path, deployed_backend_path)
            if not IS_WINDOWS:
                deployed_backend_path.chmod(deployed_backend_path.stat().st_mode | 0o111)
            print("  [OK] Backend deployed")
            config_path = BACKEND_PATH / "config.json"
            if config_path.exists():
                shutil.copy2(config_path, deploy_path / "config.json")
                print("  [OK] Backend configuration deployed")

        if not skip_frontend and not skip_backend:
            (deploy_path / "frontend-mode.txt").write_text(f"{frontend_mode}\n", encoding="ascii")
        print(f"\nDeployment complete\nLocation: {deploy_path}\nFrontend mode: {frontend_mode}")
    finally:
        if had_api_base_url:
            os.environ["VITE_API_BASE_URL"] = str(original_api_base_url)
        else:
            os.environ.pop("VITE_API_BASE_URL", None)
        if embedded_files_prepared and embedded_dist_path.exists():
            shutil.rmtree(embedded_dist_path)


def resolve_start_options(options: Dict[str, object]) -> Dict[str, object]:
    mode = str(option(options, "mode", "deployment")).lower()
    frontend_mode = str(option(options, "frontend-mode", "embedded")).lower()
    deploy_path = require_absolute_deploy_path(option(options, "deploy-path", DEFAULT_DEPLOY_PATH))
    backend_port = require_port(option(options, "backend-port", 8080), "Backend port")
    frontend_port = require_port(option(options, "frontend-port", 80), "Frontend port")
    standalone_api_base_url = str(option(options, "standalone-api-base-url", f"http://localhost:{backend_port}/api"))
    if mode not in {"development", "deployment"}:
        fail(f"Invalid mode: {mode}")
    if frontend_mode not in {"standalone", "embedded"}:
        fail(f"Invalid frontend mode: {frontend_mode}")
    if options.get("detach") and (options.get("status") or options.get("stop")):
        fail("--detach cannot be combined with --status or --stop.")
    if options.get("status") and options.get("stop"):
        fail("Use only one of --stop or --status.")
    return {
        "mode": mode,
        "frontend_mode": frontend_mode,
        "deploy_path": deploy_path,
        "backend_port": backend_port,
        "frontend_port": frontend_port,
        "standalone_api_base_url": standalone_api_base_url,
        "redeploy": False if options.get("no-redeploy") else option(options, "redeploy", True),
        "detach": bool(option(options, "detach", False)),
        "dry_run": bool(option(options, "dry-run", False)),
        "status": bool(option(options, "status", False)),
        "stop": bool(option(options, "stop", False)),
    }


def runtime_paths(settings: Dict[str, object]) -> Tuple[Path, Path]:
    if settings["mode"] == "deployment":
        runtime_directory = Path(settings["deploy_path"]) / "run"
        log_directory = Path(settings["deploy_path"]) / "logs"
    else:
        runtime_directory = PROJECT_ROOT / "tmp" / "universal-runtime"
        log_directory = runtime_directory / "logs"
    return runtime_directory / "start-all.pid", log_directory / "start-all.log"


def read_pid(pid_file: Path) -> Optional[int]:
    try:
        value = pid_file.read_text(encoding="ascii").strip()
    except (FileNotFoundError, OSError):
        return None
    return int(value) if value.isdigit() else None


def process_is_running(pid: Optional[int]) -> bool:
    if not pid:
        return False
    try:
        os.kill(pid, 0)
        return True
    except (OSError, ProcessLookupError):
        return False


def terminate_process(pid: int) -> None:
    if not process_is_running(pid):
        return
    if IS_WINDOWS:
        run_sync("taskkill", ["/PID", str(pid), "/T", "/F"], PROJECT_ROOT)
        return
    os.kill(pid, signal.SIGTERM)
    for _ in range(50):
        if not process_is_running(pid):
            return
        time.sleep(0.1)
    fail(f"Timed out while stopping process {pid}.")


def validate_deployment_artifacts(settings: Dict[str, object]) -> None:
    deploy_path = Path(settings["deploy_path"])
    backend_executable = deploy_path / BACKEND_BINARY_NAME
    mode_file = deploy_path / "frontend-mode.txt"
    if not backend_executable.exists():
        fail(f"Backend deployment not found at {backend_executable}. Enable redeployment or run deploy first.")
    if not mode_file.exists():
        fail(f"Deployment mode marker not found at {mode_file}. Redeploy the application first.")
    deployed_mode = mode_file.read_text(encoding="utf-8").strip().lower()
    if deployed_mode != settings["frontend_mode"]:
        fail(f"Existing deployment uses frontend mode '{deployed_mode}', but '{settings['frontend_mode']}' was requested.")
    if settings["frontend_mode"] == "standalone" and not (deploy_path / "frontend" / "index.html").exists():
        fail(f"Standalone frontend deployment not found at {deploy_path / 'frontend' / 'index.html'}.")


def startup_plan(settings: Dict[str, object]) -> List[Tuple[str, Path, object, List[object]]]:
    deploy_path = Path(settings["deploy_path"])
    backend_executable = deploy_path / BACKEND_BINARY_NAME
    if settings["mode"] == "development":
        return [
            ("frontend", FRONTEND_PATH, NPM_COMMAND, ["run", "dev"]),
            ("backend", BACKEND_PATH, "go", ["run", "./cmd/server", "-port", settings["backend_port"]]),
        ]
    if settings["frontend_mode"] == "embedded":
        return [("backend-with-embedded-frontend", deploy_path, backend_executable, ["-port", settings["backend_port"]])]
    return [
        ("frontend", FRONTEND_PATH, VITE_COMMAND, [
            "preview", "--host", "0.0.0.0", "--port", settings["frontend_port"],
            "--strictPort", "--outDir", deploy_path / "frontend",
        ]),
        ("backend", deploy_path, backend_executable, ["-port", settings["backend_port"]]),
    ]


def print_startup_plan(settings: Dict[str, object], plan: List[Tuple[str, Path, object, List[object]], pid_file: Path, log_file: Path) -> None:
    print(f"Mode: {settings['mode']}")
    print(f"Backend port: {settings['backend_port']}")
    print(f"Detach: {str(settings['detach']).lower()}")
    if settings["detach"]:
        print(f"PID file: {pid_file}")
        print(f"Log file: {log_file}")
    if settings["mode"] == "deployment":
        print(f"Frontend mode: {settings['frontend_mode']}")
        print(f"Redeploy: {str(settings['redeploy']).lower()}")
        print(f"Deploy path: {settings['deploy_path']}")
        if settings["frontend_mode"] == "standalone":
            print(f"Standalone API base URL: {settings['standalone_api_base_url']}")
    backend_url = f"http://localhost:{settings['backend_port']}"
    if settings["mode"] == "development":
        frontend_url = "http://localhost:3000"
    elif settings["frontend_mode"] == "standalone":
        frontend_url = f"http://localhost:{settings['frontend_port']}"
    else:
        frontend_url = backend_url
    print(f"Frontend URL: {frontend_url}")
    print(f"Backend URL: {backend_url}")
    print(f"Process count: {len(plan)}")
    for name, cwd, command, args in plan:
        print(f"Process: {name}")
        print(f"  Directory: {cwd}")
        print(f"  Command: {display_command(command, args)}")


def spawn_service(item: Tuple[str, Path, object, List[object]], settings: Dict[str, object]) -> subprocess.Popen:
    name, cwd, command, args = item
    env = os.environ.copy()
    if settings["mode"] == "development" and name == "frontend":
        env["VITE_BACKEND_URL"] = f"http://localhost:{settings['backend_port']}"
    return popen_command(command, args, cwd=str(cwd), env=env)


def run_services(settings: Dict[str, object], plan: List[Tuple[str, Path, object, List[object]]]) -> int:
    children: List[subprocess.Popen] = []
    stopping = False

    def stop_children(*_signal_args: object) -> None:
        nonlocal stopping
        if stopping:
            return
        stopping = True
        for child in children:
            if child.poll() is None:
                if IS_WINDOWS:
                    subprocess.run(["taskkill", "/PID", str(child.pid), "/T", "/F"], stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL, check=False)
                else:
                    child.terminate()

    previous_handlers = {}
    for signum in (signal.SIGINT, signal.SIGTERM):
        try:
            previous_handlers[signum] = signal.getsignal(signum)
            signal.signal(signum, stop_children)
        except (ValueError, OSError):
            pass

    try:
        for item in plan:
            try:
                children.append(spawn_service(item, settings))
            except OSError as error:
                print(f"Unable to start service: {error}", file=sys.stderr)
                stop_children()
                return 1
        print("Services started. Press Ctrl+C to stop all services.")
        while children:
            for child in children:
                code = child.poll()
                if code is not None:
                    children.remove(child)
                    if not stopping:
                        stop_children()
                    return code or 0
            time.sleep(0.1)
        return 0
    finally:
        stop_children()
        for child in children:
            try:
                child.wait(timeout=2)
            except subprocess.TimeoutExpired:
                child.kill()
        for signum, handler in previous_handlers.items():
            signal.signal(signum, handler)


def start(options: Dict[str, object], original_args: Sequence[str]) -> None:
    settings = resolve_start_options(options)
    pid_file, log_file = runtime_paths(settings)
    plan = startup_plan(settings)
    print_startup_plan(settings, plan, pid_file, log_file)
    if settings["dry_run"]:
        return

    if settings["status"]:
        pid = read_pid(pid_file)
        if pid and process_is_running(pid):
            print(f"Status: running (PID {pid})")
        else:
            pid_file.unlink(missing_ok=True)
            print("Status: stopped")
        return
    if settings["stop"]:
        pid = read_pid(pid_file)
        if not pid or not process_is_running(pid):
            pid_file.unlink(missing_ok=True)
            print("No running Shihai process found.")
            return
        terminate_process(pid)
        pid_file.unlink(missing_ok=True)
        print(f"Stopped Shihai background process {pid}.")
        return

    is_background_child = bool(options.get("_background-child"))
    if not is_background_child:
        if settings["mode"] == "development":
            assert_commands(["node", "npm.cmd" if IS_WINDOWS else "npm", "go"])
        elif settings["redeploy"]:
            deploy({
                "frontend-mode": settings["frontend_mode"],
                "deploy-path": settings["deploy_path"],
                "standalone-api-base-url": settings["standalone_api_base_url"],
            })
        else:
            validate_deployment_artifacts(settings)
        if settings["frontend_mode"] == "standalone" and not VITE_COMMAND.exists():
            fail(f"Vite not found at {VITE_COMMAND}. Run npm install in the frontend directory first.")

    if settings["detach"] and not is_background_child:
        existing_pid = read_pid(pid_file)
        if existing_pid and process_is_running(existing_pid):
            fail(f"Shihai is already running in the background with PID {existing_pid}.")
        pid_file.unlink(missing_ok=True)
        pid_file.parent.mkdir(parents=True, exist_ok=True)
        log_file.parent.mkdir(parents=True, exist_ok=True)
        log_handle = log_file.open("a", encoding="utf-8")
        command = [sys.executable, str(SCRIPT_PATH), *original_args, BACKGROUND_FLAG]
        kwargs = {"cwd": str(PROJECT_ROOT), "stdout": log_handle, "stderr": subprocess.STDOUT, "stdin": subprocess.DEVNULL}
        if IS_WINDOWS:
            kwargs["creationflags"] = subprocess.DETACHED_PROCESS | subprocess.CREATE_NEW_PROCESS_GROUP
        else:
            kwargs["start_new_session"] = True
        child = subprocess.Popen(command, **kwargs)
        log_handle.close()
        pid_file.write_text(f"{child.pid}\n", encoding="ascii")
        print(f"Started Shihai in the background with PID {child.pid}.")
        print(f"Log file: {log_file}")
        return

    if is_background_child:
        run_services(settings, plan)
        pid_file.unlink(missing_ok=True)
        return

    code = run_services(settings, plan)
    if code:
        fail(f"A Shihai service exited with code {code}.")


def start_backend(options: Dict[str, object]) -> None:
    deploy_path = require_absolute_deploy_path(option(options, "deploy-path", DEFAULT_DEPLOY_PATH))
    port = require_port(option(options, "port", 8080), "Backend port")
    executable = deploy_path / BACKEND_BINARY_NAME
    if not executable.exists():
        fail(f"Backend executable not found at {executable}. Run deploy first.")
    print(f"Starting Shihai backend...\nBackend URL: http://localhost:{port}\nPress Ctrl+C to stop.\n")
    run(executable, ["-port", port], deploy_path)


def start_development_frontend(options: Dict[str, object]) -> None:
    backend_port = require_port(option(options, "backend-port", 8080), "Backend port")
    if not (FRONTEND_PATH / "package.json").exists():
        fail(f"Frontend project not found at {FRONTEND_PATH}.")
    env = os.environ.copy()
    env["VITE_BACKEND_URL"] = f"http://localhost:{backend_port}"
    print(f"Starting Shihai frontend development server...\nFrontend URL: http://localhost:3000\nBackend proxy target: {env['VITE_BACKEND_URL']}\nPress Ctrl+C to stop.\n")
    run(NPM_COMMAND, ["run", "dev"], FRONTEND_PATH, env)


def start_frontend(options: Dict[str, object]) -> None:
    deploy_path = require_absolute_deploy_path(option(options, "deploy-path", DEFAULT_DEPLOY_PATH))
    port = require_port(option(options, "port", 80), "Frontend port")
    frontend_path = deploy_path / "frontend"
    if not (frontend_path / "index.html").exists():
        fail(f"Frontend files not found at {frontend_path}. Run deploy in standalone mode first.")
    if not VITE_COMMAND.exists():
        fail(f"Vite not found at {VITE_COMMAND}. Run npm install in the frontend directory first.")
    args = ["preview", "--host", "0.0.0.0", "--port", port, "--strictPort", "--outDir", frontend_path]
    print(f"Starting Shihai standalone frontend...\nFrontend URL: http://localhost:{port}\nCommand: {display_command(VITE_COMMAND, args)}")
    if options.get("dry-run"):
        return
    print("Press Ctrl+C to stop.\n")
    run(VITE_COMMAND, args, FRONTEND_PATH)


def run_tests(options: Dict[str, object]) -> None:
    target = str(option(options, "target", "all")).lower()
    if target not in {"all", "backend", "frontend"}:
        fail(f"Invalid test target: {target}")
    if target in {"all", "backend"}:
        print("Running backend tests...")
        run_sync("go", ["test", "./..."], BACKEND_PATH)
    if target in {"all", "frontend"}:
        print("Running frontend build...")
        run_sync(NPM_COMMAND, ["run", "build"], FRONTEND_PATH)


def main() -> int:
    options, positional = parse_args(sys.argv[1:])
    command = positional[0] if positional else "start"
    if options.get("help") or command == "help":
        print(HELP)
        return 0
    if command in {"start", "start-all"}:
        start(options, sys.argv[1:])
    elif command in {"deploy", "deploy-linux", "deploy-windows"}:
        deploy(options)
    elif command in {"backend", "start-backend"}:
        start_backend(options)
    elif command in {"frontend", "start-frontend"}:
        start_frontend(options)
    elif command in {"development-frontend", "start-development-frontend"}:
        start_development_frontend(options)
    elif command == "test":
        run_tests(options)
    else:
        fail(f"Unknown command: {command}\n\n{HELP}")
    return 0


if __name__ == "__main__":
    try:
        raise SystemExit(main())
    except ScriptError as error:
        print(f"Error: {error}", file=sys.stderr)
        raise SystemExit(1)
