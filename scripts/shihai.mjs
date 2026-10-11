#!/usr/bin/env node

/**
 * Cross-platform command line entry point for the Shihai project.
 *
 * This file intentionally uses only Node.js built-ins. It can therefore be
 * started from cmd.exe, PowerShell, Bash, or any other terminal with:
 *   node scripts/shihai.mjs <command> [options]
 */

import { spawn, spawnSync } from "node:child_process";
import { existsSync, mkdirSync, openSync, readFileSync, rmSync, chmodSync, cpSync, writeFileSync } from "node:fs";
import { dirname, join, parse, resolve, isAbsolute } from "node:path";
import { fileURLToPath } from "node:url";

const SCRIPT_PATH = fileURLToPath(import.meta.url);
const SCRIPTS_DIR = dirname(SCRIPT_PATH);
const PROJECT_ROOT = dirname(SCRIPTS_DIR);
const FRONTEND_PATH = join(PROJECT_ROOT, "frontend");
const BACKEND_PATH = join(PROJECT_ROOT, "backend");
const IS_WINDOWS = process.platform === "win32";
const DEFAULT_DEPLOY_PATH = IS_WINDOWS ? "C:\\shihai-deploy" : "/opt/shihai";
const BACKEND_BINARY_NAME = IS_WINDOWS ? "shihai-server.exe" : "shihai-server";
const NPM_COMMAND = IS_WINDOWS ? "npm.cmd" : "npm";
const VITE_COMMAND = join(FRONTEND_PATH, "node_modules", ".bin", IS_WINDOWS ? "vite.cmd" : "vite");
const BACKGROUND_FLAG = "--_background-child";

const HELP = `Usage:
  node scripts/shihai.mjs <command> [options]

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
  node scripts/shihai.mjs start --mode development
  node scripts/shihai.mjs start --frontend-mode standalone --frontend-port 4173
  node scripts/shihai.mjs deploy --frontend-mode embedded --build-only
  node scripts/shihai.mjs start --detach
  node scripts/shihai.mjs start --status
  node scripts/shihai.mjs start --stop
`;

function fail(message) {
  throw new Error(message);
}

function parseArgs(argv) {
  const options = {};
  const positional = [];
  const booleanOptions = new Set([
    "redeploy", "no-redeploy", "build-only", "skip-frontend", "skip-backend",
    "detach", "status", "stop", "dry-run", "_background-child", "help",
  ]);

  for (let index = 0; index < argv.length; index += 1) {
    const token = argv[index];
    if (token === "-h") {
      options.help = true;
      continue;
    }
    if (!token.startsWith("--")) {
      positional.push(token);
      continue;
    }

    const raw = token.slice(2);
    const equalsIndex = raw.indexOf("=");
    const name = equalsIndex >= 0 ? raw.slice(0, equalsIndex) : raw;
    let value = equalsIndex >= 0 ? raw.slice(equalsIndex + 1) : undefined;
    if (booleanOptions.has(name)) {
      if (value !== undefined) {
        value = value.toLowerCase() === "true";
      } else {
        value = true;
      }
    } else {
      if (value === undefined) {
        value = argv[index + 1];
        if (value === undefined || value.startsWith("--")) {
          fail(`Missing value for --${name}.`);
        }
        index += 1;
      }
    }
    options[name] = value;
  }

  return { options, positional };
}

function option(options, name, fallback) {
  return options[name] === undefined ? fallback : options[name];
}

function requirePort(value, name) {
  const port = Number(value);
  if (!Number.isInteger(port) || port < 1 || port > 65535) {
    fail(`${name} must be an integer between 1 and 65535.`);
  }
  return port;
}

function requireAbsoluteDeployPath(value) {
  if (!isAbsolute(value) || parse(resolve(value)).root === resolve(value)) {
    fail(`Deployment path must be an absolute path other than the filesystem root: ${value}`);
  }
  return resolve(value);
}

function commandExists(command) {
  const lookup = IS_WINDOWS ? "where" : "which";
  return spawnSync(lookup, [command], { stdio: "ignore" }).status === 0;
}

function assertCommands(commands) {
  const missing = commands.filter((command) => !commandExists(command));
  if (missing.length > 0) {
    fail(`${missing.join(", ")} must be installed and available in PATH.`);
  }
}

function displayCommand(command, args) {
  return [command, ...args].map((part) => {
    const text = String(part);
    return /[\s"]/.test(text) ? JSON.stringify(text) : text;
  }).join(" ");
}

function processOptions(command) {
  // Windows .cmd launchers require the native command shell in some Node.js
  // installations (notably when npm is provided through a PowerShell shim).
  return {
    windowsHide: false,
    shell: IS_WINDOWS && /\.(cmd|bat)$/i.test(command),
  };
}

function run(command, args, { cwd, env = process.env, stdio = "inherit" } = {}) {
  return new Promise((resolvePromise, reject) => {
    const child = spawn(command, args, { cwd, env, stdio, ...processOptions(command) });
    child.once("error", reject);
    child.once("exit", (code, signal) => {
      if (signal) {
        reject(new Error(`${displayCommand(command, args)} was terminated by ${signal}.`));
      } else if (code !== 0) {
        reject(new Error(`${displayCommand(command, args)} exited with code ${code}.`));
      } else {
        resolvePromise();
      }
    });
  });
}

function runSync(command, args, { cwd, env = process.env } = {}) {
  const result = spawnSync(command, args, { cwd, env, stdio: "inherit", ...processOptions(command) });
  if (result.error) throw result.error;
  if (result.status !== 0) fail(`${displayCommand(command, args)} exited with code ${result.status}.`);
}

function copyDirectory(source, destination) {
  mkdirSync(destination, { recursive: true });
  cpSync(source, destination, { recursive: true, force: true });
}

async function deploy(options) {
  const frontendMode = String(option(options, "frontend-mode", "standalone")).toLowerCase();
  const buildOnly = Boolean(option(options, "build-only", false));
  const skipFrontend = Boolean(option(options, "skip-frontend", false));
  const skipBackend = Boolean(option(options, "skip-backend", false));
  const deployPath = requireAbsoluteDeployPath(option(options, "deploy-path", DEFAULT_DEPLOY_PATH));
  const standaloneApiBaseUrl = option(options, "standalone-api-base-url", "http://localhost:8080/api");
  const frontendDistPath = join(FRONTEND_PATH, "dist");
  const embeddedDistPath = join(BACKEND_PATH, "internal", "webui", "dist");
  const backendBuildPath = join(BACKEND_PATH, BACKEND_BINARY_NAME);
  const originalApiBaseUrl = process.env.VITE_API_BASE_URL;
  const hadApiBaseUrl = Object.prototype.hasOwnProperty.call(process.env, "VITE_API_BASE_URL");
  let embeddedFilesPrepared = false;

  if (!["standalone", "embedded"].includes(frontendMode)) {
    fail(`Invalid frontend mode: ${frontendMode}`);
  }
  if (frontendMode === "embedded" && (skipFrontend || skipBackend)) {
    fail("Embedded mode requires both frontend and backend builds; do not use --skip-frontend or --skip-backend.");
  }

  console.log("========================================");
  console.log("Shihai Poetry Platform Deployment Tool");
  console.log("========================================");
  console.log(`Frontend mode: ${frontendMode}`);
  console.log(`Deployment path: ${deployPath}\n`);

  try {
    console.log("Checking prerequisites...");
    if (!skipFrontend) {
      assertCommands(["node", "npm"]);
      console.log("  [OK] Node.js and npm found");
    }
    if (!skipBackend) {
      assertCommands(["go"]);
      console.log("  [OK] Go found");
    }
    console.log("");

    if (!skipFrontend) {
      console.log("Building frontend...");
      process.env.VITE_API_BASE_URL = frontendMode === "standalone" ? standaloneApiBaseUrl : "/api";
      runSync(NPM_COMMAND, ["install"], { cwd: FRONTEND_PATH });
      runSync(NPM_COMMAND, ["run", "build"], { cwd: FRONTEND_PATH });
      console.log("  [OK] Frontend built successfully\n");
    }

    if (!skipBackend) {
      console.log("Building backend...");
      runSync("go", ["mod", "tidy"], { cwd: BACKEND_PATH });
      if (frontendMode === "embedded") {
        rmSync(embeddedDistPath, { recursive: true, force: true });
        copyDirectory(frontendDistPath, embeddedDistPath);
        embeddedFilesPrepared = true;
        runSync("go", ["build", "-tags", "embed_frontend", "-o", backendBuildPath, "./cmd/server"], { cwd: BACKEND_PATH });
      } else {
        runSync("go", ["build", "-o", backendBuildPath, "./cmd/server"], { cwd: BACKEND_PATH });
      }
      console.log("  [OK] Backend built successfully\n");
    }

    if (buildOnly) {
      console.log("Build-only mode: deployment skipped.");
      return;
    }

    console.log("Deploying to local machine...");
    mkdirSync(deployPath, { recursive: true });
    const deployedFrontendPath = join(deployPath, "frontend");
    if (frontendMode === "standalone" && !skipFrontend) {
      rmSync(deployedFrontendPath, { recursive: true, force: true });
      copyDirectory(frontendDistPath, deployedFrontendPath);
      console.log("  [OK] Standalone frontend deployed");
    } else if (frontendMode === "embedded") {
      rmSync(deployedFrontendPath, { recursive: true, force: true });
    }

    if (!skipBackend) {
      const deployedBackendPath = join(deployPath, BACKEND_BINARY_NAME);
      cpSync(backendBuildPath, deployedBackendPath, { force: true });
      if (!IS_WINDOWS) chmodSync(deployedBackendPath, 0o755);
      console.log("  [OK] Backend deployed");
      const configPath = join(BACKEND_PATH, "config.json");
      if (existsSync(configPath)) {
        cpSync(configPath, join(deployPath, "config.json"), { force: true });
        console.log("  [OK] Backend configuration deployed");
      }
    }

    if (!skipFrontend && !skipBackend) {
      writeFileSync(join(deployPath, "frontend-mode.txt"), `${frontendMode}\n`, "ascii");
    }
    console.log(`\nDeployment complete\nLocation: ${deployPath}\nFrontend mode: ${frontendMode}`);
  } finally {
    if (hadApiBaseUrl) process.env.VITE_API_BASE_URL = originalApiBaseUrl;
    else delete process.env.VITE_API_BASE_URL;
    if (embeddedFilesPrepared) rmSync(embeddedDistPath, { recursive: true, force: true });
  }
}

function resolveStartOptions(options) {
  const mode = String(option(options, "mode", "deployment")).toLowerCase();
  const frontendMode = String(option(options, "frontend-mode", "embedded")).toLowerCase();
  const deployPath = requireAbsoluteDeployPath(option(options, "deploy-path", DEFAULT_DEPLOY_PATH));
  const backendPort = requirePort(option(options, "backend-port", 8080), "Backend port");
  const frontendPort = requirePort(option(options, "frontend-port", 80), "Frontend port");
  const standaloneApiBaseUrl = option(options, "standalone-api-base-url", `http://localhost:${backendPort}/api`);
  if (!["development", "deployment"].includes(mode)) fail(`Invalid mode: ${mode}`);
  if (!["standalone", "embedded"].includes(frontendMode)) fail(`Invalid frontend mode: ${frontendMode}`);
  if (Boolean(options.detach) && (options.status || options.stop)) fail("--detach cannot be combined with --status or --stop.");
  if (options.status && options.stop) fail("Use only one of --status or --stop.");
  return {
    mode, frontendMode, deployPath, backendPort, frontendPort, standaloneApiBaseUrl,
    redeploy: options["no-redeploy"] ? false : option(options, "redeploy", true),
    detach: Boolean(option(options, "detach", false)),
    dryRun: Boolean(option(options, "dry-run", false)),
    status: Boolean(option(options, "status", false)),
    stop: Boolean(option(options, "stop", false)),
  };
}

function runtimePaths(settings) {
  const deployment = settings.mode === "deployment";
  const runtimeDirectory = deployment
    ? join(settings.deployPath, "run")
    : join(PROJECT_ROOT, "tmp", "universal-runtime");
  const logDirectory = deployment ? join(settings.deployPath, "logs") : join(runtimeDirectory, "logs");
  return {
    runtimeDirectory,
    pidFile: join(runtimeDirectory, "start-all.pid"),
    logFile: join(logDirectory, "start-all.log"),
  };
}

function readPid(pidFile) {
  if (!existsSync(pidFile)) return null;
  const value = readFileSync(pidFile, "utf8").trim();
  return /^\d+$/.test(value) ? Number(value) : null;
}

function processIsRunning(pid) {
  if (!pid) return false;
  try {
    process.kill(pid, 0);
    return true;
  } catch {
    return false;
  }
}

async function terminateProcess(pid) {
  if (!processIsRunning(pid)) return;
  if (IS_WINDOWS) {
    runSync("taskkill", ["/PID", String(pid), "/T", "/F"]);
    return;
  }
  process.kill(pid, "SIGTERM");
  for (let attempt = 0; attempt < 50 && processIsRunning(pid); attempt += 1) {
    await new Promise((resolvePromise) => setTimeout(resolvePromise, 100));
  }
  if (processIsRunning(pid)) fail(`Timed out while stopping process ${pid}.`);
}

function validateDeploymentArtifacts(settings) {
  const backendExecutable = join(settings.deployPath, BACKEND_BINARY_NAME);
  const modeFile = join(settings.deployPath, "frontend-mode.txt");
  if (!existsSync(backendExecutable)) fail(`Backend deployment not found at ${backendExecutable}. Enable redeployment or run deploy first.`);
  if (!existsSync(modeFile)) fail(`Deployment mode marker not found at ${modeFile}. Redeploy the application first.`);
  const deployedMode = readFileSync(modeFile, "utf8").trim().toLowerCase();
  if (deployedMode !== settings.frontendMode) fail(`Existing deployment uses frontend mode '${deployedMode}', but '${settings.frontendMode}' was requested.`);
  if (settings.frontendMode === "standalone" && !existsSync(join(settings.deployPath, "frontend", "index.html"))) {
    fail(`Standalone frontend deployment not found at ${join(settings.deployPath, "frontend", "index.html")}.`);
  }
}

function startupPlan(settings) {
  const backendExecutable = join(settings.deployPath, BACKEND_BINARY_NAME);
  const plan = [];
  if (settings.mode === "development") {
    plan.push({ name: "frontend", cwd: FRONTEND_PATH, command: NPM_COMMAND, args: ["run", "dev"] });
    plan.push({ name: "backend", cwd: BACKEND_PATH, command: "go", args: ["run", "./cmd/server", "-port", String(settings.backendPort)] });
  } else if (settings.frontendMode === "embedded") {
    plan.push({ name: "backend-with-embedded-frontend", cwd: settings.deployPath, command: backendExecutable, args: ["-port", String(settings.backendPort)] });
  } else {
    plan.push({ name: "frontend", cwd: FRONTEND_PATH, command: VITE_COMMAND, args: ["preview", "--host", "0.0.0.0", "--port", String(settings.frontendPort), "--strictPort", "--outDir", join(settings.deployPath, "frontend")] });
    plan.push({ name: "backend", cwd: settings.deployPath, command: backendExecutable, args: ["-port", String(settings.backendPort)] });
  }
  return plan;
}

function printStartupPlan(settings, plan, paths) {
  console.log(`Mode: ${settings.mode}`);
  console.log(`Backend port: ${settings.backendPort}`);
  console.log(`Detach: ${settings.detach}`);
  if (settings.detach) {
    console.log(`PID file: ${paths.pidFile}`);
    console.log(`Log file: ${paths.logFile}`);
  }
  if (settings.mode === "deployment") {
    console.log(`Frontend mode: ${settings.frontendMode}`);
    console.log(`Redeploy: ${settings.redeploy}`);
    console.log(`Deploy path: ${settings.deployPath}`);
    if (settings.frontendMode === "standalone") console.log(`Standalone API base URL: ${settings.standaloneApiBaseUrl}`);
  }
  const backendUrl = `http://localhost:${settings.backendPort}`;
  const frontendUrl = settings.mode === "development"
    ? "http://localhost:3000"
    : settings.frontendMode === "standalone"
      ? `http://localhost:${settings.frontendPort}`
      : backendUrl;
  console.log(`Frontend URL: ${frontendUrl}`);
  console.log(`Backend URL: ${backendUrl}`);
  console.log(`Process count: ${plan.length}`);
  for (const item of plan) {
    console.log(`Process: ${item.name}`);
    console.log(`  Directory: ${item.cwd}`);
    console.log(`  Command: ${displayCommand(item.command, item.args)}`);
  }
}

function spawnService(item, settings) {
  const env = { ...process.env };
  if (settings.mode === "development" && item.name === "frontend") {
    env.VITE_BACKEND_URL = `http://localhost:${settings.backendPort}`;
  }
  const child = spawn(item.command, item.args, { cwd: item.cwd, env, stdio: "inherit", ...processOptions(item.command) });
  return child;
}

async function runServices(settings, plan, { background = false } = {}) {
  const children = plan.map((item) => spawnService(item, settings));
  let stopping = false;
  const stopChildren = () => {
    if (stopping) return;
    stopping = true;
    for (const child of children) {
      if (child.exitCode === null) {
        if (IS_WINDOWS) spawnSync("taskkill", ["/PID", String(child.pid), "/T", "/F"], { stdio: "ignore" });
        else child.kill("SIGTERM");
      }
    }
  };
  const onSignal = () => stopChildren();
  process.once("SIGINT", onSignal);
  process.once("SIGTERM", onSignal);

  const result = await new Promise((resolvePromise) => {
    let settled = false;
    let remaining = children.length;
    let exitCode = 0;
    for (const child of children) {
      child.once("error", (error) => {
        if (!settled) {
          console.error(error.message);
          exitCode = 1;
          stopChildren();
        }
      });
      child.once("exit", (code, signal) => {
        remaining -= 1;
        if (!signal && code && exitCode === 0) exitCode = code;
        if (!settled && (stopping || remaining === 0)) {
          settled = true;
          stopChildren();
          resolvePromise(exitCode);
        }
      });
    }
  });
  process.removeListener("SIGINT", onSignal);
  process.removeListener("SIGTERM", onSignal);
  if (!background && result !== 0) fail(`A Shihai service exited with code ${result}.`);
  return result;
}

async function start(options, originalArgs) {
  const settings = resolveStartOptions(options);
  const paths = runtimePaths(settings);
  const plan = startupPlan(settings);
  printStartupPlan(settings, plan, paths);
  if (settings.dryRun) return;

  if (settings.status) {
    const pid = readPid(paths.pidFile);
    if (pid && processIsRunning(pid)) console.log(`Status: running (PID ${pid})`);
    else {
      rmSync(paths.pidFile, { force: true });
      console.log("Status: stopped");
    }
    return;
  }
  if (settings.stop) {
    const pid = readPid(paths.pidFile);
    if (!pid || !processIsRunning(pid)) {
      rmSync(paths.pidFile, { force: true });
      console.log("No running Shihai process found.");
      return;
    }
    await terminateProcess(pid);
    rmSync(paths.pidFile, { force: true });
    console.log(`Stopped Shihai background process ${pid}.`);
    return;
  }

  if (!options[BACKGROUND_FLAG]) {
    if (settings.mode === "development") {
      assertCommands(["node", "npm", "go"]);
    } else if (settings.redeploy) {
      await deploy({
        "frontend-mode": settings.frontendMode,
        "deploy-path": settings.deployPath,
        "standalone-api-base-url": settings.standaloneApiBaseUrl,
      });
    } else {
      validateDeploymentArtifacts(settings);
    }
    if (settings.frontendMode === "standalone" && !existsSync(VITE_COMMAND)) {
      fail(`Vite not found at ${VITE_COMMAND}. Run npm install in the frontend directory first.`);
    }
  }

  if (settings.detach && !options[BACKGROUND_FLAG]) {
    const existingPid = readPid(paths.pidFile);
    if (existingPid && processIsRunning(existingPid)) fail(`Shihai is already running in the background with PID ${existingPid}.`);
    rmSync(paths.pidFile, { force: true });
    mkdirSync(dirname(paths.pidFile), { recursive: true });
    mkdirSync(dirname(paths.logFile), { recursive: true });
    const logHandle = openSync(paths.logFile, "a");
    const child = spawn(process.execPath, [SCRIPT_PATH, ...originalArgs, BACKGROUND_FLAG], {
      cwd: PROJECT_ROOT,
      detached: true,
      stdio: ["ignore", logHandle, logHandle],
      windowsHide: true,
    });
    child.unref();
    writeFileSync(paths.pidFile, `${child.pid}\n`, "ascii");
    console.log(`Started Shihai in the background with PID ${child.pid}.`);
    console.log(`Log file: ${paths.logFile}`);
    return;
  }

  if (options[BACKGROUND_FLAG]) {
    await runServices(settings, plan, { background: true });
    rmSync(paths.pidFile, { force: true });
    return;
  }

  console.log("Services started. Press Ctrl+C to stop all services.");
  await runServices(settings, plan);
}

async function startBackend(options) {
  const deployPath = requireAbsoluteDeployPath(option(options, "deploy-path", DEFAULT_DEPLOY_PATH));
  const port = requirePort(option(options, "port", 8080), "Backend port");
  const executable = join(deployPath, BACKEND_BINARY_NAME);
  if (!existsSync(executable)) fail(`Backend executable not found at ${executable}. Run deploy first.`);
  console.log(`Starting Shihai backend...\nBackend URL: http://localhost:${port}\nPress Ctrl+C to stop.\n`);
  await run(executable, ["-port", String(port)], { cwd: deployPath });
}

async function startDevelopmentFrontend(options) {
  const backendPort = requirePort(option(options, "backend-port", 8080), "Backend port");
  if (!existsSync(join(FRONTEND_PATH, "package.json"))) fail(`Frontend project not found at ${FRONTEND_PATH}.`);
  const env = { ...process.env, VITE_BACKEND_URL: `http://localhost:${backendPort}` };
  console.log(`Starting Shihai frontend development server...\nFrontend URL: http://localhost:3000\nBackend proxy target: ${env.VITE_BACKEND_URL}\nPress Ctrl+C to stop.\n`);
  await run(NPM_COMMAND, ["run", "dev"], { cwd: FRONTEND_PATH, env });
}

async function startFrontend(options) {
  const deployPath = requireAbsoluteDeployPath(option(options, "deploy-path", DEFAULT_DEPLOY_PATH));
  const port = requirePort(option(options, "port", 80), "Frontend port");
  const frontendPath = join(deployPath, "frontend");
  if (!existsSync(join(frontendPath, "index.html"))) fail(`Frontend files not found at ${frontendPath}. Run deploy in standalone mode first.`);
  if (!existsSync(VITE_COMMAND)) fail(`Vite not found at ${VITE_COMMAND}. Run npm install in the frontend directory first.`);
  const args = ["preview", "--host", "0.0.0.0", "--port", String(port), "--strictPort", "--outDir", frontendPath];
  console.log(`Starting Shihai standalone frontend...\nFrontend URL: http://localhost:${port}\nCommand: ${displayCommand(VITE_COMMAND, args)}`);
  if (options["dry-run"]) return;
  console.log("Press Ctrl+C to stop.\n");
  await run(VITE_COMMAND, args, { cwd: FRONTEND_PATH });
}

async function runTests(options) {
  const target = String(option(options, "target", "all")).toLowerCase();
  if (!["all", "backend", "frontend"].includes(target)) fail(`Invalid test target: ${target}`);
  if (target === "all" || target === "backend") {
    console.log("Running backend tests...");
    runSync("go", ["test", "./..."], { cwd: BACKEND_PATH });
  }
  if (target === "all" || target === "frontend") {
    console.log("Running frontend build...");
    runSync(NPM_COMMAND, ["run", "build"], { cwd: FRONTEND_PATH });
  }
}

async function main() {
  const { options, positional } = parseArgs(process.argv.slice(2));
  const command = positional[0] || "start";
  if (options.help || command === "help") {
    console.log(HELP);
    return;
  }
  switch (command) {
    case "start":
    case "start-all":
      await start(options, process.argv.slice(2));
      break;
    case "deploy":
    case "deploy-linux":
    case "deploy-windows":
      await deploy(options);
      break;
    case "backend":
    case "start-backend":
      await startBackend(options);
      break;
    case "frontend":
    case "start-frontend":
      await startFrontend(options);
      break;
    case "development-frontend":
    case "start-development-frontend":
      await startDevelopmentFrontend(options);
      break;
    case "test":
      await runTests(options);
      break;
    default:
      fail(`Unknown command: ${command}\n\n${HELP}`);
  }
}

main().catch((error) => {
  console.error(`Error: ${error.message}`);
  process.exitCode = 1;
});
