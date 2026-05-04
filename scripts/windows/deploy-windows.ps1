# Shihai Poetry Platform - Windows Deployment Script
# This script builds and deploys both frontend and backend locally
# Usage: Run from project root directory

param(
    [switch]$BuildOnly,
    [switch]$SkipFrontend,
    [switch]$SkipBackend
)

$ErrorActionPreference = "Stop"

# Get script directory and project root (script is in scripts/windows/)
$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$ScriptsDir = Split-Path -Parent $ScriptDir
$ProjectRoot = Split-Path -Parent $ScriptsDir

# Configuration
$BackendPort = "8080"
$FrontendDistPath = Join-Path $ProjectRoot "frontend\dist"
$BackendBuildPath = Join-Path $ProjectRoot "backend\shihai-server.exe"
$DeployPath = "C:\shihai-deploy"

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "Shihai Poetry Platform Deployment Tool" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

# Function to check if command exists
function Test-Command($Command) {
    $exists = $null -ne (Get-Command $Command -ErrorAction SilentlyContinue)
    return $exists
}

# Check prerequisites
Write-Host "Checking prerequisites..." -ForegroundColor Yellow

if (-not $SkipFrontend) {
    if (-not (Test-Command "node")) {
        Write-Error "Node.js is not installed or not in PATH"
        exit 1
    }
    Write-Host "  [OK] Node.js found" -ForegroundColor Green
}

if (-not $SkipBackend) {
    if (-not (Test-Command "go")) {
        Write-Error "Go is not installed or not in PATH"
        exit 1
    }
    Write-Host "  [OK] Go found" -ForegroundColor Green
}

Write-Host ""

# Build Frontend
if (-not $SkipFrontend) {
    Write-Host "Building Frontend..." -ForegroundColor Yellow
    Set-Location -Path (Join-Path $ProjectRoot "frontend")
    
    Write-Host "  Installing dependencies..." -ForegroundColor Gray
    npm install
    if ($LASTEXITCODE -ne 0) {
        Write-Error "Failed to install frontend dependencies"
        exit 1
    }
    
    Write-Host "  Building production bundle..." -ForegroundColor Gray
    npm run build
    if ($LASTEXITCODE -ne 0) {
        Write-Error "Failed to build frontend"
        exit 1
    }
    
    Set-Location -Path ".."
    Write-Host "  [OK] Frontend built successfully" -ForegroundColor Green
    Write-Host ""
}

# Build Backend
if (-not $SkipBackend) {
    Write-Host "Building Backend..." -ForegroundColor Yellow
    Set-Location -Path (Join-Path $ProjectRoot "backend")
    
    Write-Host "  Downloading dependencies..." -ForegroundColor Gray
    go mod tidy
    if ($LASTEXITCODE -ne 0) {
        Write-Error "Failed to download backend dependencies"
        exit 1
    }
    
    Write-Host "  Building executable..." -ForegroundColor Gray
    go build -o shihai-server.exe cmd/server/main.go
    if ($LASTEXITCODE -ne 0) {
        Write-Error "Failed to build backend"
        exit 1
    }
    
    Set-Location -Path ".."
    Write-Host "  [OK] Backend built successfully" -ForegroundColor Green
    Write-Host ""
}

# Deploy
if (-not $BuildOnly) {
    Write-Host "Deploying to local machine..." -ForegroundColor Yellow
    
    # Create deployment directory
    if (-not (Test-Path $DeployPath)) {
        New-Item -ItemType Directory -Path $DeployPath -Force | Out-Null
        New-Item -ItemType Directory -Path "$DeployPath\frontend" -Force | Out-Null
    }
    
    # Copy frontend files
    if (-not $SkipFrontend) {
        Write-Host "  Copying frontend files..." -ForegroundColor Gray
        Copy-Item -Path "$FrontendDistPath\*" -Destination "$DeployPath\frontend\" -Recurse -Force
        Write-Host "  [OK] Frontend deployed" -ForegroundColor Green
    }
    
    # Copy backend files
    if (-not $SkipBackend) {
        Write-Host "  Copying backend executable..." -ForegroundColor Gray
        Copy-Item -Path $BackendBuildPath -Destination "$DeployPath\" -Force
        Write-Host "  [OK] Backend deployed" -ForegroundColor Green
    }
    
    Write-Host ""
    Write-Host "========================================" -ForegroundColor Cyan
    Write-Host "Deployment Complete!" -ForegroundColor Green
    Write-Host "========================================" -ForegroundColor Cyan
    Write-Host ""
    Write-Host "Deployment location: $DeployPath" -ForegroundColor White
    Write-Host ""
    Write-Host "To start the application:" -ForegroundColor Yellow
    Write-Host "  1. Start Backend:  .\scripts\windows\start-backend.ps1" -ForegroundColor White
    Write-Host "  2. Start Frontend: .\scripts\windows\start-frontend.ps1" -ForegroundColor White
    Write-Host "  3. Or start both:  .\scripts\windows\start-all.ps1" -ForegroundColor White
    Write-Host ""
} else {
    Write-Host "Build only mode - skipping deployment" -ForegroundColor Yellow
}
