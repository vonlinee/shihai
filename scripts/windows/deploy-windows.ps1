# Shihai Poetry Platform - Windows Deployment Script
# Builds standalone frontend/backend artifacts or a backend with embedded frontend assets.

param(
    [ValidateSet("Standalone", "Embedded")]
    [string]$FrontendMode = "Standalone",
    [switch]$BuildOnly,
    [switch]$SkipFrontend,
    [switch]$SkipBackend,
    [string]$DeployPath = "C:\shihai-deploy",
    [string]$StandaloneApiBaseUrl = "http://localhost:8080/api"
)

$ErrorActionPreference = "Stop"

$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$ScriptsDir = Split-Path -Parent $ScriptDir
$ProjectRoot = Split-Path -Parent $ScriptsDir
$FrontendPath = Join-Path $ProjectRoot "frontend"
$BackendPath = Join-Path $ProjectRoot "backend"
$FrontendDistPath = Join-Path $FrontendPath "dist"
$EmbeddedDistPath = Join-Path $BackendPath "internal\webui\dist"
$BackendBuildPath = Join-Path $BackendPath "shihai-server.exe"
$DeployModePath = Join-Path $DeployPath "frontend-mode.txt"
$OriginalLocation = Get-Location
$OriginalApiBaseUrl = $env:VITE_API_BASE_URL
$HadOriginalApiBaseUrl = Test-Path Env:VITE_API_BASE_URL
$EmbeddedFilesPrepared = $false

function Test-Command {
    param([string]$Command)

    return $null -ne (Get-Command $Command -ErrorAction SilentlyContinue)
}

function Assert-LastCommandSucceeded {
    param([string]$Message)

    if ($LASTEXITCODE -ne 0) {
        throw $Message
    }
}

if ($FrontendMode -eq "Embedded" -and ($SkipFrontend -or $SkipBackend)) {
    throw "Embedded mode requires both frontend and backend builds. Do not use -SkipFrontend or -SkipBackend."
}

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "Shihai Poetry Platform Deployment Tool" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "Frontend mode: $FrontendMode" -ForegroundColor White
Write-Host "Deployment path: $DeployPath" -ForegroundColor White
Write-Host ""

try {
    Write-Host "Checking prerequisites..." -ForegroundColor Yellow

    if (-not $SkipFrontend) {
        if (-not (Test-Command "node") -or -not (Test-Command "npm")) {
            throw "Node.js and npm must be installed and available in PATH."
        }
        Write-Host "  [OK] Node.js and npm found" -ForegroundColor Green
    }

    if (-not $SkipBackend) {
        if (-not (Test-Command "go")) {
            throw "Go must be installed and available in PATH."
        }
        Write-Host "  [OK] Go found" -ForegroundColor Green
    }

    Write-Host ""

    if (-not $SkipFrontend) {
        Write-Host "Building frontend..." -ForegroundColor Yellow
        Set-Location -Path $FrontendPath

        if ($FrontendMode -eq "Standalone") {
            $env:VITE_API_BASE_URL = $StandaloneApiBaseUrl
        } else {
            $env:VITE_API_BASE_URL = "/api"
        }

        npm install
        Assert-LastCommandSucceeded "Failed to install frontend dependencies."

        npm run build
        Assert-LastCommandSucceeded "Failed to build frontend."

        Write-Host "  [OK] Frontend built successfully" -ForegroundColor Green
        Write-Host ""
    }

    if (-not $SkipBackend) {
        Write-Host "Building backend..." -ForegroundColor Yellow
        Set-Location -Path $BackendPath

        go mod tidy
        Assert-LastCommandSucceeded "Failed to prepare backend dependencies."

        if ($FrontendMode -eq "Embedded") {
            if (Test-Path $EmbeddedDistPath) {
                Remove-Item -LiteralPath $EmbeddedDistPath -Recurse -Force
            }
            New-Item -ItemType Directory -Path $EmbeddedDistPath -Force | Out-Null
            Copy-Item -Path (Join-Path $FrontendDistPath "*") -Destination $EmbeddedDistPath -Recurse -Force
            $EmbeddedFilesPrepared = $true

            go build -tags embed_frontend -o $BackendBuildPath ./cmd/server
        } else {
            go build -o $BackendBuildPath ./cmd/server
        }
        Assert-LastCommandSucceeded "Failed to build backend."

        Write-Host "  [OK] Backend built successfully" -ForegroundColor Green
        Write-Host ""
    }

    if ($BuildOnly) {
        Write-Host "Build-only mode: deployment skipped." -ForegroundColor Yellow
        return
    }

    Write-Host "Deploying to local machine..." -ForegroundColor Yellow
    New-Item -ItemType Directory -Path $DeployPath -Force | Out-Null

    if ($FrontendMode -eq "Standalone" -and -not $SkipFrontend) {
        $DeployedFrontendPath = Join-Path $DeployPath "frontend"
        if (Test-Path $DeployedFrontendPath) {
            Remove-Item -LiteralPath $DeployedFrontendPath -Recurse -Force
        }
        New-Item -ItemType Directory -Path $DeployedFrontendPath -Force | Out-Null
        Copy-Item -Path (Join-Path $FrontendDistPath "*") -Destination $DeployedFrontendPath -Recurse -Force
        Write-Host "  [OK] Standalone frontend deployed" -ForegroundColor Green
    } elseif ($FrontendMode -eq "Embedded") {
        $DeployedFrontendPath = Join-Path $DeployPath "frontend"
        if (Test-Path $DeployedFrontendPath) {
            Remove-Item -LiteralPath $DeployedFrontendPath -Recurse -Force
        }
    }

    if (-not $SkipBackend) {
        Copy-Item -Path $BackendBuildPath -Destination (Join-Path $DeployPath "shihai-server.exe") -Force
        Write-Host "  [OK] Backend deployed" -ForegroundColor Green

        $BackendConfigPath = Join-Path $BackendPath "config.json"
        if (Test-Path $BackendConfigPath) {
            Copy-Item -LiteralPath $BackendConfigPath -Destination (Join-Path $DeployPath "config.json") -Force
            Write-Host "  [OK] Backend configuration deployed" -ForegroundColor Green
        }
    }

    if (-not $SkipFrontend -and -not $SkipBackend) {
        Set-Content -LiteralPath $DeployModePath -Value $FrontendMode -Encoding ASCII
    }

    Write-Host ""
    Write-Host "========================================" -ForegroundColor Cyan
    Write-Host "Deployment complete" -ForegroundColor Green
    Write-Host "========================================" -ForegroundColor Cyan
    Write-Host "Location: $DeployPath" -ForegroundColor White
    Write-Host "Frontend mode: $FrontendMode" -ForegroundColor White
} finally {
    Set-Location -Path $OriginalLocation
    if ($HadOriginalApiBaseUrl) {
        $env:VITE_API_BASE_URL = $OriginalApiBaseUrl
    } else {
        Remove-Item Env:VITE_API_BASE_URL -ErrorAction SilentlyContinue
    }
    if ($EmbeddedFilesPrepared -and (Test-Path $EmbeddedDistPath)) {
        Remove-Item -LiteralPath $EmbeddedDistPath -Recurse -Force
    }
}
