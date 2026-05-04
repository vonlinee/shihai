# Start Both Backend and Frontend (Windows)
# Usage: Run from project root directory

$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$ScriptsDir = Split-Path -Parent $ScriptDir
$ProjectRoot = Split-Path -Parent $ScriptsDir
$DeployPath = "C:\shihai-deploy"

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "Shihai Poetry Platform - Starting All" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

# Check if deployment exists
if (-not (Test-Path $DeployPath)) {
    Write-Error "Deployment not found at $DeployPath"
    Write-Host "Please run .\scripts\windows\deploy-windows.ps1 first" -ForegroundColor Yellow
    exit 1
}

# Start Backend in new window
Write-Host "Starting Backend Server..." -ForegroundColor Green
Start-Process powershell -ArgumentList "-NoExit", "-Command", "$ProjectRoot\scripts\windows\start-backend.ps1" -WindowStyle Normal

# Wait a moment for backend to start
Start-Sleep -Seconds 2

# Start Frontend in new window
Write-Host "Starting Frontend Server..." -ForegroundColor Green
Start-Process powershell -ArgumentList "-NoExit", "-Command", "$ProjectRoot\scripts\windows\start-frontend.ps1" -WindowStyle Normal

Write-Host ""
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "All services started!" -ForegroundColor Green
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "Backend: http://localhost:8080" -ForegroundColor White
Write-Host "Frontend: http://localhost" -ForegroundColor White
Write-Host ""
Write-Host "Close the PowerShell windows to stop the services" -ForegroundColor Yellow
