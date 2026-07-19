# Starts the Shihai frontend development server.

param(
    [ValidateRange(1, 65535)]
    [int]$BackendPort = 8080
)

$ErrorActionPreference = "Stop"

$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$ScriptsDir = Split-Path -Parent $ScriptDir
$ProjectRoot = Split-Path -Parent $ScriptsDir
$FrontendPath = Join-Path $ProjectRoot "frontend"

if (-not (Test-Path (Join-Path $FrontendPath "package.json"))) {
    throw "Frontend project not found at $FrontendPath."
}

$env:VITE_BACKEND_URL = "http://localhost:$BackendPort"

Write-Host "Starting Shihai frontend development server..." -ForegroundColor Cyan
Write-Host "Frontend URL: http://localhost:3000" -ForegroundColor Yellow
Write-Host "Backend proxy target: $env:VITE_BACKEND_URL" -ForegroundColor Yellow
Write-Host "Press Ctrl+C to stop." -ForegroundColor Gray
Write-Host ""

Set-Location -Path $FrontendPath
npm run dev
if ($LASTEXITCODE -ne 0) {
    throw "Frontend development server stopped with exit code $LASTEXITCODE."
}
