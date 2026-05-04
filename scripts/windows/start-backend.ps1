# Start Backend Server (Windows)
# Usage: Run from project root directory

$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$DeployPath = "C:\shihai-deploy"
$BackendPath = "$DeployPath\shihai-server.exe"

if (-not (Test-Path $BackendPath)) {
    Write-Error "Backend executable not found at $BackendPath"
    Write-Host "Please run .\scripts\windows\deploy-windows.ps1 first"
    exit 1
}

Write-Host "Starting Shihai Backend Server..." -ForegroundColor Cyan
Write-Host "Server will run on http://localhost:8080" -ForegroundColor Yellow
Write-Host "Press Ctrl+C to stop" -ForegroundColor Gray
Write-Host ""

Set-Location -Path $DeployPath
& $BackendPath
