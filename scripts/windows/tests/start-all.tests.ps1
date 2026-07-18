$ErrorActionPreference = "Stop"

$TestDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$WindowsScriptsDir = Split-Path -Parent $TestDir
$ProjectRoot = Split-Path -Parent (Split-Path -Parent $WindowsScriptsDir)
$StartScript = Join-Path $WindowsScriptsDir "start-all.ps1"
$StartFrontendScript = Join-Path $WindowsScriptsDir "start-frontend.ps1"

function Assert-Contains {
    param(
        [string[]]$Output,
        [string]$Expected
    )

    $Text = $Output -join "`n"
    if (-not $Text.Contains($Expected)) {
        throw "Expected output to contain '$Expected'. Actual output:`n$Text"
    }
}

$StartFrontendCommand = Get-Command $StartFrontendScript
if (-not $StartFrontendCommand.Parameters.ContainsKey("DryRun")) {
    throw "start-frontend.ps1 must support -DryRun for command verification."
}

$OriginalTerminalSession = $env:WT_SESSION
$HadTerminalSession = Test-Path Env:WT_SESSION
try {
    Remove-Item Env:WT_SESSION -ErrorAction SilentlyContinue
    $DefaultOutput = & $StartScript -DryRun 6>&1
    Assert-Contains $DefaultOutput "Mode: Deployment"
    Assert-Contains $DefaultOutput "Backend port: 8080"
    Assert-Contains $DefaultOutput "Frontend mode: Embedded"
    Assert-Contains $DefaultOutput "Terminal window: New (new)"
    Assert-Contains $DefaultOutput "Tab count: 1"
    Assert-Contains $DefaultOutput "Backend (with embedded frontend)"

    $env:WT_SESSION = "test-terminal-session"
    $AttachedOutput = & $StartScript -DryRun 6>&1
    Assert-Contains $AttachedOutput "Terminal window: Current (0)"
} finally {
    if ($HadTerminalSession) {
        $env:WT_SESSION = $OriginalTerminalSession
    } else {
        Remove-Item Env:WT_SESSION -ErrorAction SilentlyContinue
    }
}

$DevelopmentOutput = & $StartScript -Mode Development -DryRun 6>&1
Assert-Contains $DevelopmentOutput "Mode: Development"
Assert-Contains $DevelopmentOutput "Tab count: 2"
Assert-Contains $DevelopmentOutput "VITE_BACKEND_URL='http://localhost:8080'"
Assert-Contains $DevelopmentOutput "go run ./cmd/server -port 8080"

$StandaloneOutput = & $StartScript -Mode Deployment -FrontendMode Standalone -DryRun 6>&1
Assert-Contains $StandaloneOutput "Frontend mode: Standalone"
Assert-Contains $StandaloneOutput "Redeploy: True"
Assert-Contains $StandaloneOutput "Tab count: 2"

$EmbeddedOutput = & $StartScript -Mode Deployment -FrontendMode Embedded -DryRun 6>&1
Assert-Contains $EmbeddedOutput "Frontend mode: Embedded"
Assert-Contains $EmbeddedOutput "Tab count: 1"
Assert-Contains $EmbeddedOutput "Backend (with embedded frontend)"

$ReuseOutput = & $StartScript -Mode Deployment -Redeploy:$false -DryRun 6>&1
Assert-Contains $ReuseOutput "Redeploy: False"

$CustomTitleOutput = & $StartScript -Mode Development -FrontendTabTitle "Poem Web" -BackendTabTitle "Poem API" -DryRun 6>&1
Assert-Contains $CustomTitleOutput "Tab: Poem Web"
Assert-Contains $CustomTitleOutput "Tab: Poem API"

$EmbeddedTitleOutput = & $StartScript -Mode Deployment -FrontendMode Embedded -BackendTabTitle "Poem App" -DryRun 6>&1
Assert-Contains $EmbeddedTitleOutput "Tab: Poem App"

$CustomPortOutput = & $StartScript -Mode Development -BackendPort 9090 -DryRun 6>&1
Assert-Contains $CustomPortOutput "Backend port: 9090"
Assert-Contains $CustomPortOutput "VITE_BACKEND_URL='http://localhost:9090'"
Assert-Contains $CustomPortOutput "go run ./cmd/server -port 9090"

$StandalonePortOutput = & $StartScript -FrontendMode Standalone -BackendPort 9090 -DryRun 6>&1
Assert-Contains $StandalonePortOutput "Standalone API base URL: http://localhost:9090/api"
Assert-Contains $StandalonePortOutput "-Port 9090"

$CustomApiOutput = & $StartScript -FrontendMode Standalone -BackendPort 9090 -StandaloneApiBaseUrl "https://api.example.com/api" -DryRun 6>&1
Assert-Contains $CustomApiOutput "Standalone API base URL: https://api.example.com/api"

$TemporaryDeployPath = Join-Path $ProjectRoot "tmp\start-frontend-test"
try {
    $TemporaryFrontendPath = Join-Path $TemporaryDeployPath "frontend"
    New-Item -ItemType Directory -Path $TemporaryFrontendPath -Force | Out-Null
    Set-Content -LiteralPath (Join-Path $TemporaryFrontendPath "index.html") -Value "<html></html>"

    $FrontendOutput = & $StartFrontendScript -DeployPath $TemporaryDeployPath -Port 4173 -DryRun 6>&1
    Assert-Contains $FrontendOutput "vite.cmd preview"
    Assert-Contains $FrontendOutput "--port 4173"
    Assert-Contains $FrontendOutput "--outDir $TemporaryFrontendPath"
} finally {
    if (Test-Path $TemporaryDeployPath) {
        Remove-Item -LiteralPath $TemporaryDeployPath -Recurse -Force
    }
}

Write-Host "start-all.ps1 dry-run tests passed." -ForegroundColor Green
