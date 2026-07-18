# Starts the deployed Shihai backend.

param(
    [string]$DeployPath = "C:\shihai-deploy",
    [ValidateRange(1, 65535)]
    [int]$Port = 8080
)

$ErrorActionPreference = "Stop"
$BackendExecutable = Join-Path $DeployPath "shihai-server.exe"

if (-not (Test-Path $BackendExecutable)) {
    throw "Backend executable not found at $BackendExecutable. Run deploy-windows.ps1 first."
}

Write-Host "Starting Shihai backend..." -ForegroundColor Cyan
Write-Host "Backend URL: http://localhost:$Port" -ForegroundColor Yellow
Write-Host "Press Ctrl+C to stop." -ForegroundColor Gray
Write-Host ""

Set-Location -Path $DeployPath
& $BackendExecutable -port $Port
