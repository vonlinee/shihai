import assert from "node:assert/strict";
import { spawnSync } from "node:child_process";
import { dirname, join } from "node:path";
import { fileURLToPath } from "node:url";
import test from "node:test";

const TEST_DIR = dirname(fileURLToPath(import.meta.url));
const SCRIPT = join(TEST_DIR, "..", "shihai.mjs");

function runScript(...args) {
  return spawnSync(process.execPath, [SCRIPT, ...args], {
    encoding: "utf8",
    windowsHide: true,
  });
}

test("prints help from the operating system terminal", () => {
  const result = runScript("--help");
  assert.equal(result.status, 0);
  assert.match(result.stdout, /node scripts\/shihai\.mjs <command>/);
  assert.match(result.stdout, /start-development-frontend/);
});

test("prints a cross-platform development startup plan", () => {
  const result = runScript("start", "--mode", "development", "--backend-port", "9090", "--dry-run");
  assert.equal(result.status, 0, result.stderr);
  assert.match(result.stdout, /Mode: development/);
  assert.match(result.stdout, /Backend port: 9090/);
  assert.match(result.stdout, /Frontend URL: http:\/\/localhost:3000/);
  assert.match(result.stdout, /Backend URL: http:\/\/localhost:9090/);
  assert.match(result.stdout, /Process count: 2/);
  assert.match(result.stdout, /go run \.\/cmd\/server -port 9090/);
});

test("supports standalone deployment planning without platform shells", () => {
  const deployPath = process.platform === "win32" ? "C:\\tmp\\shihai-test" : "/tmp/shihai-test";
  const result = runScript(
    "start", "--frontend-mode", "standalone", "--deploy-path", deployPath,
    "--backend-port", "9090", "--frontend-port", "4173", "--no-redeploy", "--dry-run",
  );
  assert.equal(result.status, 0, result.stderr);
  assert.match(result.stdout, /Frontend mode: standalone/);
  assert.match(result.stdout, /Process count: 2/);
  assert.match(result.stdout, /--port 4173/);
  assert.match(result.stdout, /Frontend URL: http:\/\/localhost:4173/);
  assert.match(result.stdout, /Backend URL: http:\/\/localhost:9090/);
});

test("prints the shared backend URL for embedded deployment", () => {
  const deployPath = process.platform === "win32" ? "C:\\tmp\\shihai-embedded-test" : "/tmp/shihai-embedded-test";
  const result = runScript(
    "start", "--frontend-mode", "embedded", "--deploy-path", deployPath,
    "--backend-port", "9091", "--no-redeploy", "--dry-run",
  );
  assert.equal(result.status, 0, result.stderr);
  assert.match(result.stdout, /Frontend URL: http:\/\/localhost:9091/);
  assert.match(result.stdout, /Backend URL: http:\/\/localhost:9091/);
});

test("rejects invalid ports and incompatible embedded build options", () => {
  const invalidPort = runScript("start", "--backend-port", "0", "--dry-run");
  assert.notEqual(invalidPort.status, 0);
  assert.match(invalidPort.stderr, /Backend port must be an integer/);

  const incompatible = runScript("deploy", "--frontend-mode", "embedded", "--skip-frontend");
  assert.notEqual(incompatible.status, 0);
  assert.match(incompatible.stderr, /Embedded mode requires both frontend and backend/);
});
