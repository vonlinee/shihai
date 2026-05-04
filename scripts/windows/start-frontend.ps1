# Start Frontend Server (Windows)
# This script serves the frontend using Python's built-in HTTP server
# Usage: Run from project root directory

$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$DeployPath = "C:\shihai-deploy\frontend"
$Port = "80"

if (-not (Test-Path $DeployPath)) {
    Write-Error "Frontend files not found at $DeployPath"
    Write-Host "Please run .\scripts\windows\deploy-windows.ps1 first"
    exit 1
}

# Check if Python is installed
$pythonCmd = Get-Command python -ErrorAction SilentlyContinue
if (-not $pythonCmd) {
    $pythonCmd = Get-Command python3 -ErrorAction SilentlyContinue
}

Write-Host "Starting Shihai Frontend Server..." -ForegroundColor Cyan
Write-Host "Server will run on http://localhost:$Port" -ForegroundColor Yellow
Write-Host "Press Ctrl+C to stop" -ForegroundColor Gray
Write-Host ""

Set-Location -Path $DeployPath

if ($pythonCmd) {
    # Use Python's built-in HTTP server
    & $pythonCmd.Source -m http.server $Port
} else {
    # Fallback: Try to use Node.js http-server if available
    Write-Host "Python not found, trying Node.js..." -ForegroundColor Yellow
    npx http-server -p $Port
}
