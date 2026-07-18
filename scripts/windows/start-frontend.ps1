# Starts the standalone deployed Shihai frontend with Vite Preview.

param(
    [string]$DeployPath = "C:\shihai-deploy",
    [ValidateRange(1, 65535)]
    [int]$Port = 80,
    [switch]$DryRun
)

$ErrorActionPreference = "Stop"

$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$ScriptsDir = Split-Path -Parent $ScriptDir
$ProjectRoot = Split-Path -Parent $ScriptsDir
$FrontendProjectPath = Join-Path $ProjectRoot "frontend"
$FrontendPath = Join-Path $DeployPath "frontend"
$FrontendIndex = Join-Path $FrontendPath "index.html"
$ViteCommand = Join-Path $FrontendProjectPath "node_modules\.bin\vite.cmd"
$ViteArguments = @(
    "preview",
    "--host", "0.0.0.0",
    "--port", $Port,
    "--strictPort",
    "--outDir", $FrontendPath
)

if (-not (Test-Path $FrontendIndex)) {
    throw "Frontend files not found at $FrontendPath. Run deploy-windows.ps1 in Standalone mode first."
}
if (-not (Test-Path $ViteCommand)) {
    throw "Vite not found at $ViteCommand. Run npm install in the frontend directory first."
}

Write-Host "Starting Shihai standalone frontend..." -ForegroundColor Cyan
Write-Host "Frontend URL: http://localhost:$Port" -ForegroundColor Yellow
Write-Host "Command: $ViteCommand $($ViteArguments -join ' ')" -ForegroundColor Gray

if ($DryRun) {
    return
}

Write-Host "Press Ctrl+C to stop." -ForegroundColor Gray
Write-Host ""

Set-Location -Path $FrontendProjectPath
& $ViteCommand @ViteArguments
if ($LASTEXITCODE -ne 0) {
    throw "Vite Preview stopped with exit code $LASTEXITCODE."
}
